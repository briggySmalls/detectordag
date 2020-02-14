"""Main module."""
import os
import base64
import logging
from pathlib import Path

from time import sleep

logger = logging.getLogger(__name__)

CERT_ROOT_PATH = Path(__file__).parent / 'certs'


def payload_report(self, params, packet):
    logger.info("----- New Payload -----")
    logger.info("Topic: %s", packet.topic)
    logger.info("Message: %s", packet.payload)
    logger.info("-----------------------")


def set_cred(env_name, file_name):
    # Turn base64 encoded environmental variable into a certificate file
    env = os.getenv(env_name)
    with (CERT_ROOT_PATH / file_name).open('wb') as output_file:
        output_file.write(base64.b64decode(env))


def run():
    """Runs the application"""
    logger.debug("MQTT Thing Starting...")

    # Create certificates from environment variables
    CERT_ROOT_PATH.mkdir(parents=True, exist_ok=True)
    set_cred("AWS_ROOT_CERT", "root-CA.crt")
    set_cred("AWS_THING_CERT", "thing.cert.pem")
    set_cred("AWS_PRIVATE_CERT", "thing.private.key")

    # Configure the client
    config = ClientConfig(
        aws_endpoint = os.getenv("AWS_ENDPOINT", "data.iot.us-east-1.amazonaws.com")
        aws_port = os.getenv("AWS_PORT", 8883)
        device_uuid = os.getenv("BALENA_DEVICE_UUID")
        )
    with CloudClient(config) as client:
        # Subscribe to the desired topic and register a callback.
        client.subscribe("balena/payload_test", 1, payload_report)

        # Send messages too
        i = 0
        while True:
            i += 1
            logger.info('Publishing to "balena/payload_write_test" the value: %i', i)
            client.publish("balena/payload_write_test", i, 0)
            sleep(5)
