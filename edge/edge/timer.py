"""Simple periodic timer"""

from threading import Timer
from typing import Callable


class PeriodicTimer:
    """Simple periodic timer"""

    def __init__(self, period: float, callback: Callable[[], None]) -> None:
        self._callback = callback
        self._period = period
        self._timer = None

    def start(self) -> None:
        """Starts the periodic execution of the callback function"""
        self._tick()

    def _tick(self) -> None:
        # Execute the callback
        self._callback()
        # Enqueue to run again
        self._timer = Timer(self._period, self._tick)
        self._timer.start()

    def stop(self) -> None:
        if self._timer is not None:
            self._timer.cancel()
