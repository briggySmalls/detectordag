"""Console script for edge."""
import logging
import sys
from typing import Any

import click

from edge.config import AppConfig
from edge.edge import EdgeApp

_POWER_PIN = 4


@click.group()
@click.pass_context
def main(ctx: Any) -> None:
    """Entrypoint for the edge application"""
    # ensure that ctx.obj exists and is a dict (in case `cli()` is called
    # by means other than the `if` block below
    ctx.ensure_object(dict)
    # Configure logging
    logging.basicConfig(level=logging.DEBUG)
    logging.info("Starting")
    # Parse config
    ctx.obj['config'] = AppConfig.from_env()


@main.command()
@click.pass_context
def app(ctx: Any) -> None:
    """Run the 'production' edge software"""
    # Track power status GPIO
    from gpiozero import DigitalInputDevice  # noqa: E501, pylint: disable=import-error,import-outside-toplevel
    power_status_device = DigitalInputDevice(_POWER_PIN, bounce_time=0.2)
    # Start the application
    with EdgeApp(power_status_device, ctx.obj['config']):
        while True:
            pass


@main.command()
@click.pass_context
def mock(ctx: Any) -> None:
    """Run the mock edge software"""
    # Create a mock device
    from edge.mocks import MockDigitalInputDevice  # noqa: E501, pylint: disable=import-outside-toplevel
    power_status_device = MockDigitalInputDevice(_POWER_PIN)
    # Run the 'real' software
    with EdgeApp(power_status_device, ctx.obj['config']):
        # Allow the user to toggle the power status
        while True:
            char = click.getchar()
            if char == 't':
                # Toggle power status
                power_status_device.toggle()
            elif char == 'h':
                # Set power status high
                power_status_device.high()
            elif char == 'l':
                # Set power status low
                power_status_device.low()


if __name__ == "__main__":
    # Configure logging
    sys.exit(main())  # pragma: no cover pylint: disable=no-value-for-parameter
