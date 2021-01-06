"""Main module."""
import logging
from types import TracebackType
from typing import Optional, Type

from edge.aws import ClientConfig, CloudClient
from edge.config import AppConfig

try:
    from gpiozero import DigitalInputDevice
except ImportError:
    from edge.mocks import (  # noqa: E501,  pylint: disable=ungrouped-imports
        MockDigitalInputDevice as DigitalInputDevice,
    )

_LOGGER = logging.getLogger(__name__)


class EdgeApp:
    """Wrapper for the entire application"""

    def __init__(self, device: DigitalInputDevice, config: AppConfig) -> None:
        self.config = config
        # Prepare configuration for the client
        client_config = ClientConfig(
            device_id=config.aws_thing_name,
            endpoint=config.aws_endpoint,
            port=config.aws_port,
            root_cert=config.aws_root_cert,
            thing_cert=config.aws_thing_cert,
            thing_key=config.aws_thing_key,
        )
        self._device = device
        # Create the client
        self._client = CloudClient(client_config)

    def __enter__(self) -> "EdgeApp":
        # Connect the MQTT client
        self._client.__enter__()
        # Configure the device
        logging.info("Configuring edge...")
        self.configure()
        logging.info("Configured!")
        # Send the current status
        self._publish_update(self._device)
        # Return this instance
        return self

    def configure(self) -> None:
        """Configure the app"""
        # Send messages when power status changes
        self._device.when_activated = self._publish_update
        self._device.when_deactivated = self._publish_update

    def __exit__(
        self,
        exc_type: Optional[Type[BaseException]],
        exc_value: Optional[BaseException],
        traceback: Optional[TracebackType],
    ) -> None:
        # Teardown the AWS client
        self._client.__exit__(exc_type, exc_value, traceback)

    def _publish_update(self, device: DigitalInputDevice) -> None:
        # Get the status
        status = bool(device.value)
        # Publish
        self._client.power_status_changed(status)
