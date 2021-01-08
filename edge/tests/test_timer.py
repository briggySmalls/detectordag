from time import sleep

import pytest

from edge.timer import PeriodicTimer


def test_timer() -> None:
    # Create a timer to increment periodically
    period = 0.01
    counter = 0
    def increment() -> None:
        nonlocal counter
        counter += 1
    t = PeriodicTimer(period, increment)
    # Start
    t.start()
    # Wait a bit
    factor = 10
    sleep(period * factor)
    # Expect roughly 'factor' number of ticks
    assert pytest.approx(factor, counter)
    t.stop()

def test_stop_timer() -> None:
    # Create
    t = PeriodicTimer(0.01, lambda: print("hi"))
    # Ensure we can 'stop' it without ever starting it
    t.stop()
