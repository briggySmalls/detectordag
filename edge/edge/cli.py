"""Console script for edge."""
import logging
import sys

import click

from edge.edge import run


@click.command()
def main():
    """Console script for edge."""
    logging.basicConfig(level=logging.DEBUG)
    logging.info("Starting")
    run()
    return 0


if __name__ == "__main__":
    # Configure logging
    sys.exit(main())  # pragma: no cover
