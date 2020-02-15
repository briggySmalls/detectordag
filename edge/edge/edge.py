"""Main module."""
import base64
import logging
import os
from pathlib import Path
from time import sleep

from edge.aws import ClientConfig, CloudClient

logger = logging.getLogger(__name__)

CERT_ROOT_PATH = Path(__file__).parent / 'certs'
CERTS_PATHS = {
    'root_cert': CERT_ROOT_PATH / "root-CA.crt",
    'thing_cert': CERT_ROOT_PATH / "thing.cert.pem",
    'thing_key': CERT_ROOT_PATH / "thing.private.key",
}


def payload_report(self, params, packet):
    logger.info("----- New Payload -----")
    logger.info("Topic: %s", packet.topic)
    logger.info("Message: %s", packet.payload)
    logger.info("-----------------------")


def set_cred(env_name: str, file: Path) -> None:
    # Turn base64 encoded environmental variable into a certificate file
    env = os.getenv(env_name)
    with file.open('wb') as output_file:
        output_file.write(base64.b64decode(env))


def run():
    """Runs the application"""
    logger.debug("MQTT Thing Starting...")

    # Create certificates from environment variables
    CERT_ROOT_PATH.mkdir(parents=True, exist_ok=True)
    set_cred("AWS_ROOT_CERT", CERTS_PATHS['root_cert'])
    set_cred("AWS_THING_CERT", CERTS_PATHS['thing_cert'])
    set_cred("AWS_PRIVATE_CERT", CERTS_PATHS['thing_key'])

    # Configure the client
    config = ClientConfig(
        device_id=os.getenv("BALENA_DEVICE_UUID"),
        endpoint=os.getenv("AWS_ENDPOINT", "data.iot.us-east-1.amazonaws.com"),
        port=os.getenv("AWS_PORT", "8883"),
        **{key: str(value)
           for key, value in CERTS_PATHS.items()})
    with CloudClient(config) as client:
        # Subscribe to the desired topic and register a callback.
        client.subscribe("balena/payload_test", 1, payload_report)

        # Send messages too
        i = 0
        while True:
            i += 1
            logger.info(
                'Publishing to "balena/payload_write_test" the value: %i', i)
            client.publish("balena/payload_write_test", i, 0)
            sleep(5)
