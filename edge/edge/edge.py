"""Main module."""
import os
import base64
import logging
from pathlib import Path

from time import sleep
from AWSIoTPythonSDK.MQTTLib import AWSIoTMQTTClient

logger = logging.getLogger(__name__)

CERT_ROOT_PATH = Path(__file__).parent / 'certs'


def payload_report(self, params, packet):
    logger.info("----- New Payload -----")
    logger.info("Topic: ", packet.topic)
    logger.info("Message: ", packet.payload)
    logger.info("-----------------------")


def set_cred(env_name, file_name):
    # Turn base64 encoded environmental variable into a certificate file
    env = os.getenv(env_name)
    with (CERT_ROOT_PATH / file_name).open('wb') as output_file:
        output_file.write(base64.b64decode(env))


def run():
    """Runs the application"""
    logger.debug("MQTT Thing Starting...")

    # Configure the client
    client = setup_mqtt()

    # Subscribe to the desired topic and register a callback.
    client.subscribe("balena/payload_test", 1, payload_report)

    # Send messages too
    i = 0
    while True:
        i += 1
        print('Publishing to "balena/payload_write_test" the value: ', i)
        client.publish("balena/payload_write_test", i, 0)
        sleep(5)


def setup_mqtt():
    """
    Configure MQTT
    """
    aws_endpoint = os.getenv("AWS_ENDPOINT", "data.iot.us-east-1.amazonaws.com")
    aws_port = os.getenv("AWS_PORT", 8883)
    device_uuid = os.getenv("BALENA_DEVICE_UUID")

    # Save credential files
    CERT_ROOT_PATH.mkdir(parents=True, exist_ok=True)
    set_cred("AWS_ROOT_CERT", "root-CA.crt")
    set_cred("AWS_THING_CERT", "thing.cert.pem")
    set_cred("AWS_PRIVATE_CERT", "thing.private.key")

    # Unique ID. If another connection using the same key is opened the previous one is auto closed by AWS IOT
    client = AWSIoTMQTTClient(device_uuid)
    #Used to configure the host name and port number the underneath AWS IoT MQTT Client tries to connect to.
    client.configureEndpoint(aws_endpoint, aws_port)
    # Used to configure the rootCA, private key and certificate files. configureCredentials(CAFilePath, KeyPath='', CertificatePath='')
    client.configureCredentials(
        (CERT_ROOT_PATH / "root-CA.crt").resolve(),
        (CERT_ROOT_PATH / "thing.private.key").resolve(),
        (CERT_ROOT_PATH / "thing.cert.pem").resolve())
    # Configure the offline queue for publish requests to be 20 in size and drop the oldest
    client.configureOfflinePublishQueueing(-1)
    # Used to configure the draining speed to clear up the queued requests when the connection is back. (frequencyInHz)
    client.configureDrainingFrequency(2)
    # Configure connect/disconnect timeout to be 10 seconds
    client.configureConnectDisconnectTimeout(10)
    # Configure MQTT operation timeout to be 5 seconds
    client.configureMQTTOperationTimeout(5)
    # Connect to AWS IoT with default keepalive set to 600 seconds
    client.connect()

    return client
