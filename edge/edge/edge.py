"""Main module."""
import logging
from types import TracebackType
from typing import Optional, Type

from edge.aws import ClientConfig, CloudClient
from edge.config import AppConfig
from edge.data import DeviceShadowState
from edge.timer import PeriodicTimer

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
        # Prepare to periodically check for status changes
        self._timer = PeriodicTimer(config.power_poll_period, self._check_status)
        self._previous_status = None

    def __enter__(self) -> "EdgeApp":
        # Connect the MQTT client
        self._client.__enter__()
        # Configure the device
        logging.info("Configuring edge...")
        self.configure()
        logging.info("Configured!")
        # Return this instance
        return self

    def configure(self) -> None:
        """Configure the app"""
        # Send updates when power status changes
        self._device.when_activated = self._publish_update
        self._device.when_deactivated = self._publish_update
        # Check to send updates on a timer
        self._timer.start()

    def __exit__(
        self,
        exc_type: Optional[Type[BaseException]],
        exc_value: Optional[BaseException],
        traceback: Optional[TracebackType],
    ) -> None:
        # Teardown the AWS client
        self._client.__exit__(exc_type, exc_value, traceback)
        # Stop the timer
        self._timer.stop()

    def _get_status(self) -> DeviceShadowState:
        """Fetch the current device state"""
        return DeviceShadowState(status=self._device.value)

    def _check_status(self) -> None:
        """
        Check our most-recent message is still valid

        For some reason recently gpiozero's edge detection has been playing up.
        This function sends an update if the last message we sent is out-of-date.
        """
        if self._previous_status == self._get_status():
            # No change, short-circuit
            return
        # We need to send an update
        _LOGGER.info("Periodic check noticed status change")
        self._publish_update()

    def _record_status(self, status: DeviceShadowState) -> None:
        # Record the provided payload
        self._previous_status = status

    def _publish_update(self) -> None:
        """Publish an update to the cloud"""
        # Get the current status of the device
        status = self._get_status()
        # Send it
        self._client.send_status_update(status, callback=self._record_status)
        # Record what we sent
        self._previous_status = status
