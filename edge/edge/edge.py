"""Main module."""
import base64
import logging
import os
from pathlib import Path
from time import sleep

from edge.aws import ClientConfig, CloudClient
from gpiozero import DigitalInputDevice

logger = logging.getLogger(__name__)

_POWER_PIN = 4
CERT_ROOT_PATH = Path(__file__).parent / 'certs'
CERTS_PATHS = {
    'root_cert': CERT_ROOT_PATH / "root-CA.crt",
    'thing_cert': CERT_ROOT_PATH / "thing.cert.pem",
    'thing_key': CERT_ROOT_PATH / "thing.private.key",
}


def _create_certs():
    # Create certificates from environment variables
    CERT_ROOT_PATH.mkdir(parents=True, exist_ok=True)
    _set_cred("AWS_ROOT_CERT", CERTS_PATHS['root_cert'])
    _set_cred("AWS_THING_CERT", CERTS_PATHS['thing_cert'])
    _set_cred("AWS_PRIVATE_CERT", CERTS_PATHS['thing_key'])


def _set_cred(env_name: str, file: Path) -> None:
    # Turn base64 encoded environmental variable into a certificate file
    env = os.getenv(env_name)
    with file.open('wb') as output_file:
        output_file.write(base64.b64decode(env))


def _publish_update(client: CloudClient, device: DigitalInputDevice) -> None:
    # Get the status
    status = device.value
    # Publish
    client.power_status_changed(status)


def run():
    """Runs the application"""
    logger.debug("MQTT Thing Starting...")

    # Ensure certificates are available
    _create_certs()

    # Prepare configuration for the client
    config = ClientConfig(
        device_id=os.getenv("BALENA_DEVICE_UUID"),
        endpoint=os.getenv("AWS_ENDPOINT"),
        port=os.getenv("AWS_PORT", "8883"),
        **{key: str(value)
           for key, value in CERTS_PATHS.items()})

    # Create a client
    with CloudClient(config) as client:
        # Track power status GPIO
        power_status_device = DigitalInputDevice(_POWER_PIN)
        # Send messages when power status changes
        power_status_device.when_activated = lambda device: _publish_update(
            client, device)
        power_status_device.when_deactivated = lambda device: _publish_update(
            client, device)
