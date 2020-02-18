"""Logic for connecting to AWS IoT"""
import logging
from dataclasses import dataclass, asdict
from pathlib import Path
import json

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


@dataclass
class PowerStatusChangedPayload:
    version = '0.1'
    status: bool

    def to_json(self) -> str:
        return json.dumps(asdict(self))


class CloudClient:
    """Client for interfacing with the cloud"""
    _QOS = 0
    _POWER_STATUS_TOPIC = 'detectordag/power_status_changed'
    _DRAINING_FREQUENCY = 2
    _DISCONNECT_TIMEOUT = 10
    _OPERATION_TIMEOUT = 5

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
        self.client.configureDrainingFrequency(self._DRAINING_FREQUENCY)
        # Configure connect/disconnect timeout to be 10 seconds
        self.client.configureConnectDisconnectTimeout(self._DISCONNECT_TIMEOUT)
        # Configure MQTT operation timeout to be 5 seconds
        self.client.configureMQTTOperationTimeout(self._OPERATION_TIMEOUT)

    def __enter__(self) -> 'CloudClient':
        # Connect
        self.client.connect()
        # Return this
        return self

    def __exit__(self, exc_type, exc_value, traceback) -> None:
        del exc_type, exc_value, traceback
        self.client.disconnect()

    def power_status_changed(self, status: bool) -> None:
        """Send a messaging indicating the power status has updated

        Args:
            status (bool): New power status
        """
        _LOGGER.info('Publishing to "%s" the value: %i',
                     self._POWER_STATUS_TOPIC, status)
        payload = PowerStatusChangedPayload(status=status)
        self.client.publish(self._POWER_STATUS_TOPIC, payload.to_json(), self._QOS)
