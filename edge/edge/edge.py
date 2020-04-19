"""Main module."""
import logging
from threading import Timer

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
        # Prepare configuration for the client
        config = ClientConfig(device_id=config.aws_thing_name,
                              endpoint=config.aws_endpoint,
                              port=config.aws_port,
                              root_cert=config.aws_root_cert,
                              thing_cert=config.aws_thing_cert,
                              thing_key=config.aws_thing_key)
        self._device = device
        # Create the client
        self._client = CloudClient(config)
        # Preallocate the timer
        self._timer = None
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

    def __exit__(self, exc_type, exc_value, traceback) -> None:
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
