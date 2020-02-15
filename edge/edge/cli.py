"""Console script for edge."""
import logging
import sys

import click

from edge.edge import EdgeApp

_POWER_PIN = 4


@click.group()
def main():
    """Entrypoint for the edge application"""
    # Configure logging
    logging.basicConfig(level=logging.DEBUG)
    logging.info("Starting")


@main.command()
def app():
    """Run the 'production' edge software"""
    # Track power status GPIO
    from gpiozero import DigitalInputDevice  # noqa: E501, pylint: disable=import-error,import-outside-toplevel
    power_status_device = DigitalInputDevice(_POWER_PIN)
    # Start the application
    with EdgeApp(power_status_device):
        while True:
            pass


@main.command()
def mock():
    """Run the mock edge software"""
    # Create a mock device
    from edge.mocks import MockDigitalInputDevice  # noqa: E501, pylint: disable=import-outside-toplevel
    power_status_device = MockDigitalInputDevice(_POWER_PIN)
    # Run the 'real' software
    with EdgeApp(power_status_device):
        # Allow the user to toggle the power status
        while True:
            char = click.getchar()
            if char == 'p':
                # Toggle power status
                power_status_device.toggle()


if __name__ == "__main__":
    # Configure logging
    sys.exit(main())  # pragma: no cover
