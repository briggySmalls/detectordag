"""Tests for the timer class"""

from time import sleep

import pytest

from edge.timer import PeriodicTimer


def test_timer() -> None:
    """Test running a timer over multiple ticks"""
    # Create a timer to increment periodically
    period = 0.01
    counter = 0

    def increment() -> None:
        nonlocal counter
        counter += 1

    timer = PeriodicTimer(period, increment)
    # Start
    timer.start()
    # Wait a bit
    factor = 10
    sleep(period * factor)
    # Expect roughly 'factor' number of ticks
    assert pytest.approx(factor, counter)
    timer.stop()


def test_stop_timer() -> None:
    """Test stopping the timer"""
    # Create
    timer = PeriodicTimer(0.01, lambda: print("hi"))
    # Ensure we can 'stop' it without ever starting it
    timer.stop()
