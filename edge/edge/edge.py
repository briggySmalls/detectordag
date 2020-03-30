"""Main module."""
import logging

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
        self.device = device
        # Create the client
        self.client = CloudClient(config)

    def __enter__(self) -> 'EdgeApp':
        self.client.__enter__()
        # Configure the device
        self.configure()
        # Return this instance
        return self

    def configure(self) -> None:
        """Configure the app
        """
        # Send messages when power status changes
        self.device.when_activated = self._publish_update
        self.device.when_deactivated = self._publish_update

    def __exit__(self, exc_type, exc_value, traceback) -> None:
        self.client.__exit__(exc_type, exc_value, traceback)

    def _publish_update(self, device: DigitalInputDevice) -> None:
        # Get the status
        status = bool(device.value)
        # Publish
        self.client.power_status_changed(status)
