"""Logic for connecting to AWS IoT"""
import logging
from dataclasses import dataclass
from pathlib import Path

from AWSIoTPythonSDK.MQTTLib import AWSIoTMQTTClient

_LOGGER = logging.getLogger(__file__)


@dataclass
class ClientConfig:
    """Configuration for the CloudClient"""

    device_id: str
    root_cert: Path
    thing_cert: Path
    thing_key: Path
    endpoint: str
    port: int
    draining_frequency: int = 2
    disconnect_timeout: int = 10
    operation_timeout: int = 5


class CloudClient:
    """Client for interfacing with the cloud"""
    _QOS = 0
    _POWER_STATUS_TOPIC = 'detectordag/power_status_changed'

    def __init__(self, config: ClientConfig) -> None:
        self.config = config
        # Unique ID. If another connection using the same key is opened the
        # previous one is auto closed by AWS IOT
        self.client = AWSIoTMQTTClient(config.device_id)
        # Used to configure the host name and port number the underneath AWS
        # IoT MQTT Client tries to connect to.
        self.client.configureEndpoint(self.config.endpoint, self.config.port)
        # Used to configure the rootCA, private key and certificate files.
        # configureCredentials(CAFilePath, KeyPath='', CertificatePath='')
        self.client.configureCredentials(str(self.config.root_cert.resolve()),
                                         str(self.config.thing_key.resolve()),
                                         str(self.config.thing_cert.resolve()))
        # Configure the offline queue for publish requests to be 20 in size and
        # drop the oldest
        self.client.configureOfflinePublishQueueing(-1)
        # Used to configure the draining speed to clear up the queued requests
        # when the connection is back. (frequencyInHz)
        self.client.configureDrainingFrequency(self.config.draining_frequency)
        # Configure connect/disconnect timeout to be 10 seconds
        self.client.configureConnectDisconnectTimeout(
            self.config.disconnect_timeout)
        # Configure MQTT operation timeout to be 5 seconds
        self.client.configureMQTTOperationTimeout(
            self.config.operation_timeout)

    def __enter__(self) -> 'CloudClient':
        # Connect
        self.client.connect()

    def __exit__(self, exc_type, exc_value, traceback) -> None:
        del exc_type, exc_value, traceback
        self.client.disconnect()

    def status_update(self, status: bool) -> None:
        """Send a messaging indicating the power status has updated

        Args:
            status (bool): New power status
        """
        _LOGGER.info('Publishing to "%s" the value: %i',
                     self._POWER_STATUS_TOPIC, status)
        self.client.publish(self._POWER_STATUS_TOPIC, status, self._QOS)
