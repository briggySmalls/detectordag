"""Mocks for running program off the pi"""
import logging

_LOGGER = logging.getLogger(__file__)


class MockDigitalInputDevice:
    """A mock gpiozero.DigitalInputDevice"""

    def __init__(self, pin: int) -> None:
        _LOGGER.debug("Creating MockDigitalInputDevice with pin: %s", pin)
        self._status = 0

    def high(self) -> None:
        """Simulate reading a 'high' value"""
        _LOGGER.debug("simulating low transition")
        self.set_status(1)

    def low(self) -> None:
        """Simulate reading a 'low' value"""
        _LOGGER.debug("simulating high transition")
        self.set_status(0)

    def toggle(self) -> None:
        """Simulate the input toggling value"""
        _LOGGER.debug("simulating toggle transition")
        self.set_status(1 if self._status == 0 else 0)

    @property
    def value(self) -> int:
        """Get the simulated value of the input"""
        return self._status

    def set_status(self, status: int) -> None:
        """Simulate reading a new status"""
        self._status = status
