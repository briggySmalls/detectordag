"""Mocks for running program off the pi"""
import logging
from edge.ina219 import INA219

_LOGGER = logging.getLogger(__file__)


class MockPower:
    """A mock Power"""

    def __init__(self, bus: INA219) -> None:
        _LOGGER.debug("Creating MockPower with bus: %s", bus)
        self._status = 0

    def high(self) -> None:
        """Simulate reading a 'high' value"""
        _LOGGER.debug("simulating high transition")
        self.set_status(True)

    def low(self) -> None:
        """Simulate reading a 'low' value"""
        _LOGGER.debug("simulating low transition")
        self.set_status(False)

    def toggle(self) -> None:
        """Simulate the input toggling value"""
        _LOGGER.debug("simulating toggle transition")
        self.set_status(not self._status)

    def is_powered(self) -> bool:
        """Get the simulated value of the input"""
        return self._status

    def set_status(self, status: bool) -> None:
        """Simulate reading a new status"""
        self._status = status
