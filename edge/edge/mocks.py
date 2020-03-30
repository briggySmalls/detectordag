"""Mocks for running program off the pi"""
import logging

_LOGGER = logging.getLogger(__file__)


class MockDigitalInputDevice:
    """A mock DigitalInputDevice"""
    def __init__(self, pin: int) -> None:
        _LOGGER.debug("Creating MockDigitalInputDevice with pin: %s", pin)
        self.when_activated = None
        self.when_deactivated = None
        self._status = False

    def high(self) -> None:
        """Simulate reading a 'high' value
        """
        self.set_status(1)

    def low(self) -> None:
        """Simulate reading a 'low' value
        """
        self.set_status(0)

    def toggle(self) -> None:
        """Simulate the input toggling value
        """
        self.set_status(not self._status)

    @property
    def value(self) -> int:
        """Get the simulated value of the input"""
        return self._status

    def set_status(self, status: int) -> None:
        """Simulate reading a new status"""
        self._status = status
        if status:
            self.when_activated(self)  # pylint: disable=not-callable
        else:
            self.when_deactivated(self)  # pylint: disable=not-callable
