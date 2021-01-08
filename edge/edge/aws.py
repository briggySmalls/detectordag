"""Logic for connecting to AWS IoT"""
import logging
from dataclasses import dataclass
from pathlib import Path
from types import TracebackType
from typing import Optional, Type, Callable

from AWSIoTPythonSDK.MQTTLib import AWSIoTMQTTShadowClient

from edge.data import DeviceShadowState, PowerStatus
from edge.exceptions import ConnectionFailedError

_LOGGER = logging.getLogger(__file__)
logging.getLogger("AWSIoTPythonSDK").setLevel(logging.WARNING)


@dataclass
class ClientConfig:
    """Configuration for the CloudClient"""

    device_id: str
    root_cert: Path
    thing_cert: Path
    thing_key: Path
    endpoint: str
    port: int


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
        self.client.configureCredentials(
            str(self.config.root_cert.resolve()),
            str(self.config.thing_key.resolve()),
            str(self.config.thing_cert.resolve()),
        )
        self.client.configureCredentials(
            str(self.config.root_cert.resolve()),
            str(self.config.thing_key.resolve()),
            str(self.config.thing_cert.resolve()),
        )
        # Configure connect/disconnect timeout to be 10 seconds
        self.client.configureConnectDisconnectTimeout(self._DISCONNECT_TIMEOUT)
        # Configure MQTT operation timeout to be 5 seconds
        self.client.configureMQTTOperationTimeout(self._OPERATION_TIMEOUT)
        # Create the shadow handler
        self.shadow = self.client.createShadowHandlerWithName(
            config.device_id, False
        )

    def __enter__(self) -> "CloudClient":
        # Connect
        logging.info("Connecting client...")
        if not self.client.connect():
            raise ConnectionFailedError()
        logging.info("Connected!")
        # Return this
        return self

    def __exit__(
        self,
        exc_type: Optional[Type[BaseException]],
        exc_value: Optional[BaseException],
        traceback: Optional[TracebackType],
    ) -> None:
        del exc_type, exc_value, traceback
        self.client.disconnect()

    def send_status_update(
        self,
        state: DeviceShadowState,
        callback: Callable[[DeviceShadowState], None] = None,
    ) -> None:
        """Send a messaging indicating the power status has updated

        Args:
            status (bool): New power status
        """
        payload = state.json()
        _LOGGER.info("Publishing status update: %s", payload)
        token = self.shadow.shadowUpdate(
            payload,
            lambda payload, response_status, token: self.shadow_update_handler(
                payload, response_status, token, callback
            ),
            self._OPERATION_TIMEOUT,
        )
        _LOGGER.debug("Status update returned token: %s", token)

    @staticmethod
    def shadow_update_handler(
        payload: str,
        response_status: str,
        token: str,
        callback: Callable[[DeviceShadowState], None],
    ) -> None:
        """Handle a device shadow update response

        Args:
            payload (str): Response body
            response_status (str): Response status
            token (str): Request identifier

        Raises:
            RuntimeError: Unexpected response
        """
        # Log the outcome of the update
        del token
        if response_status == "accepted":
            _LOGGER.info("Shadow update accepted: payload=%s", payload)
            # Send confirmation to the caller, if requested
            if callback is not None:
                callback(DeviceShadowState.parse_raw(payload))
        elif response_status in ["timeout", "rejected"]:
            _LOGGER.error(
                "Shadow update failed: status=%s, payload=%s",
                response_status,
                payload,
            )
        else:
            raise RuntimeError(
                f"Unexpected response_status: {response_status}"
            )
