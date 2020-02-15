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
        self._change_status(True)

    def low(self) -> None:
        self._change_status(False)

    def _change_status(self, status: bool) -> None:
        if status:
            self.when_activated(self)
        else:
            self.when_deactivated(self)
