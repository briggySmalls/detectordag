"""Console script for edge."""
import logging
import sys
from threading import Event
from typing import Any
from smbus2 import SMBus

import click

from edge.config import AppConfig
from edge.edge import EdgeApp
from edge.power import Power
from edge.ina219 import INA219
from edge.exceptions import DetectorDagException

_I2C_BUS = 1


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
    ctx.obj["config"] = AppConfig.from_env()


@main.command()
@click.pass_context
def app(ctx: Any) -> None:
    """Run the 'production' edge software"""
    # Track power status GPIO
    bus = SMBus(_I2C_BUS)
    ina219 = INA219(bus)
    power = Power(ina219)

    try:
        # Start the application
        with EdgeApp(power, ctx.obj["config"]):
            # Sleep forever without burning clock cycles
            Event().wait()
    except DetectorDagException as err:
        logging.exception(err)


@main.command()
@click.pass_context
def mock(ctx: Any) -> None:
    """Run the mock edge software"""
    # Create a mock device
    from edge.mocks import (  # noqa: E501, pylint: disable=import-outside-toplevel
        MockPower,
    )

    power_status_device = MockPower(_I2C_BUS)
    # Run the 'real' software
    with EdgeApp(power_status_device, ctx.obj["config"]):
        # Allow the user to toggle the power status
        while True:
            logging.debug("listening for input")
            char = click.getchar()
            if char == "t":
                # Toggle power status
                power_status_device.toggle()
            elif char == "h":
                # Set power status high
                power_status_device.high()
            elif char == "l":
                # Set power status low
                power_status_device.low()


if __name__ == "__main__":
    # Configure logging
    sys.exit(main())  # pragma: no cover pylint: disable=no-value-for-parameter
