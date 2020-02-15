"""Mocks for running program off the pi"""
import logging

_LOGGER = logging.getLogger(__file__)


class MockDigitalInputDevice:
    """A mock DigitalInputDevice"""
    def __init__(self, pin: int) -> None:
        _LOGGER.debug("Creating MockDigitalInputDevice with pin: %s", pin)
        self.when_activated = None
        self.when_deactivated = None

    def high(self) -> None:
        """Simulate reading a 'high' value
        """
        self.set_status(True)

    def low(self) -> None:
        """Simulate reading a 'low' value
        """
        self.set_status(False)

    def set_status(self, status: bool) -> None:
        """Simulate reading a new status"""
        if status:
            self.when_activated(self)  # pylint: disable=not-callable
        else:
            self.when_deactivated(self)  # pylint: disable=not-callable
