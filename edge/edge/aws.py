"""Logic for connecting to AWS IoT"""
import logging
from dataclasses import dataclass, asdict
from pathlib import Path
import json

from AWSIoTPythonSDK.MQTTLib import AWSIoTMQTTShadowClient

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
    status: bool

    def to_json(self) -> str:
        return json.dumps(asdict(self))


class CloudClient:
    """Client for interfacing with the cloud"""
    _QOS = 1
    _DISCONNECT_TIMEOUT = 10
    _OPERATION_TIMEOUT = 5

    def __init__(self, config: ClientConfig) -> None:
        self.config = config
        # Unique ID. If another connection using the same key is opened the
        # previous one is auto closed by AWS IOT
        self.client = AWSIoTMQTTShadowClient(config.device_id)
        # Used to configure the host name and port number the underneath AWS
        # IoT MQTT Client tries to connect to.
        self.client.configureEndpoint(self.config.endpoint, self.config.port)
        # Used to configure the rootCA, private key and certificate files.
        # configureCredentials(CAFilePath, KeyPath='', CertificatePath='')
        self.client.configureCredentials(str(self.config.root_cert.resolve()),
                                         str(self.config.thing_key.resolve()),
                                         str(self.config.thing_cert.resolve()))
        # Configure connect/disconnect timeout to be 10 seconds
        self.client.configureConnectDisconnectTimeout(self._DISCONNECT_TIMEOUT)
        # Configure MQTT operation timeout to be 5 seconds
        self.client.configureMQTTOperationTimeout(self._OPERATION_TIMEOUT)
        # Create the shadow handler
        self.shadow = self.client.createShadowHandlerWithName(config.device_id, False)

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
        payload = PowerStatusChangedPayload(status=status).to_json()
        _LOGGER.info('Publishing status update: %s', payload)
        token = self.shadow.shadowUpdate(payload, self.shadow_update_handler, self._OPERATION_TIMEOUT)
        _LOGGER.debug(f"Status update returned token: {token}")

    @staticmethod
    def shadow_update_handler(payload: str, response_status: str, token: str) -> None:
        del payload, token
        _LOGGER.debug("Handling update")
        if response_status == 'accepted':
            _LOGGER.info("Shadow update accepted")
        elif response_status in ['timeout', 'rejected']:
            _LOGGER.error("Show updated failed: status={response_status}")
