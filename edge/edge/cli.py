"""Console script for edge."""
import sys
import click
from edge.edge import run


@click.command()
def main():
    """Console script for edge."""
    run()
    return 0


if __name__ == "__main__":
    sys.exit(main())  # pragma: no cover
