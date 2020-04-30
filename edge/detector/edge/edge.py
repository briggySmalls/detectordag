"""Main module."""
import logging
from threading import Timer
from types import TracebackType
from typing import Optional, Type
import base64
from pathlib import Path
import tempfile
from uuid import uuid4

from edge.aws import ClientConfig, CloudClient
from edge.config import AppConfig

try:
    from gpiozero import DigitalInputDevice
except ImportError:
    from edge.mocks import MockDigitalInputDevice as DigitalInputDevice  # noqa: E501,  pylint: disable=ungrouped-imports

_LOGGER = logging.getLogger(__name__)


class EdgeApp:
    """Wrapper for the entire application"""
    def __init__(self, device: DigitalInputDevice, config: AppConfig) -> None:
        self.config = config
        # Write the root cert to a file
        root_cert_path = self._write_cert(config.aws_root_cert)
        # Prepare configuration for the client
        client_config = ClientConfig(device_id=config.aws_thing_name,
                                     endpoint=config.aws_endpoint,
                                     port=config.aws_port,
                                     root_cert=root_cert_path,
                                     thing_cert=config.aws_thing_cert_path,
                                     thing_key=config.aws_thing_key_path)
        self._device = device
        # Create the client
        self._client = CloudClient(client_config)
        # Preallocate the timer
        self._timer = None  # type: Optional[Timer]
        self._is_cancelled = False

    def __enter__(self) -> 'EdgeApp':
        self._client.__enter__()
        # Configure the device
        self.configure()
        # Return this instance
        return self

    def configure(self) -> None:
        """Configure the app
        """
        # Send messages when power status changes
        self._device.when_activated = self._publish_update
        self._device.when_deactivated = self._publish_update
        # Start the alive timer
        self._tick()

    def __exit__(self, exc_type: Optional[Type[BaseException]],
                 exc_value: Optional[BaseException],
                 traceback: Optional[TracebackType]) -> None:
        # Teardown the AWS client
        self._client.__exit__(exc_type, exc_value, traceback)
        # Cancel any running timers
        if self._timer is not None:
            self._timer.cancel()
            self._is_cancelled = True

    def _publish_update(self, device: DigitalInputDevice) -> None:
        # Get the status
        status = bool(device.value)
        # Publish
        self._client.power_status_changed(status)

    def _tick(self) -> None:
        if self._is_cancelled:
            # Short-circuit if we've cancelled already
            return
        # Publish an update
        self._publish_update(self._device)
        # Schedule another tick
        self._timer = Timer(self.config.alive_interval, self._tick)
        self._timer.start()

    @staticmethod
    def _write_cert(cert: str) -> Path:
        # Get the temporary directory
        tmp = Path(tempfile.gettempdir())
        file = tmp / f"{uuid4()}.pem"
        # Turn base64 encoded string into a certificate file
        with file.open('wb') as output_file:
            output_file.write(base64.b64decode(cert))
        return file
