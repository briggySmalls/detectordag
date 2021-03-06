"""Logic for connecting to AWS IoT"""
import logging
from dataclasses import dataclass
from pathlib import Path
from types import TracebackType
from typing import Any, Awaitable, Callable, Optional, Type

from awscrt import io
from awscrt import mqtt as awsmqtt
from awsiot import mqtt_connection_builder
from awsiot.iotshadow import IotShadowClient, ShadowState, UpdateShadowRequest

from edge.data import DeviceShadowState

_LOGGER = logging.getLogger(__name__)


@dataclass
class ClientConfig:
    """Configuration for the CloudClient"""

    device_id: str
    root_cert: Path
    thing_cert: Path
    thing_key: Path
    endpoint: str
    keep_alive: int


class CloudClient:
    """Client for interfacing with the cloud"""

    _OPERATION_TIMEOUT = 5

    _mqtt: Optional[awsmqtt.Connection]
    _shadow: Optional[IotShadowClient]

    def __init__(
        self, config: ClientConfig, status_request_callback: Callable[[], None]
    ) -> None:
        # Record the configuration
        self._config = config
        self._status_request_callback = status_request_callback
        self._mqtt = None
        self._shadow = None

    def __enter__(self) -> "CloudClient":
        # Spin up resources
        _LOGGER.info("Initialising...")
        event_loop_group = io.EventLoopGroup(1)
        host_resolver = io.DefaultHostResolver(event_loop_group)
        client_bootstrap = io.ClientBootstrap(event_loop_group, host_resolver)
        # Create a connection (the shadow needs a started connection)
        _LOGGER.info("Connecting...")
        self._mqtt = self._create_mqtt_connection(client_bootstrap)
        connected_future = self._mqtt.connect()
        # Create a shadow client
        _LOGGER.info("Creating shadow client...")
        self._shadow = IotShadowClient(self._mqtt)
        # Wait for connection to be fully established.
        # Note that it's not necessary to wait, commands issued to the
        # mqtt_connection before its fully connected will simply be queued.
        # But this sample waits here so it's obvious when a connection
        # fails or succeeds.
        connected_future.result()
        _LOGGER.info("Connected!")
        self._subscribe_to_update_requests(self._mqtt)
        return self

    def __exit__(
        self,
        exc_type: Optional[Type[BaseException]],
        exc_value: Optional[BaseException],
        traceback: Optional[TracebackType],
    ) -> None:
        del exc_type, exc_value, traceback
        _LOGGER.info("Disconnecting MQTT")
        assert self._mqtt is not None
        future = self._mqtt.disconnect()
        future.result()

    def send_status_update(
        self,
        state: DeviceShadowState,
    ) -> None:
        """Send a messaging indicating the power status has updated

        Args:
            status (bool): New power status
        """
        # Construct the request
        payload = state.dict()
        _LOGGER.info("Publishing status update: %s", payload)
        request = UpdateShadowRequest(
            thing_name=self._config.device_id,
            state=ShadowState(
                reported=payload,
            ),
        )
        # Make the request
        assert self._shadow is not None
        future = self._shadow.publish_update_shadow(
            request, awsmqtt.QoS.AT_LEAST_ONCE
        )
        future.add_done_callback(self._on_status_update_published)

    def _create_mqtt_connection(
        self, client_bootstrap: io.ClientBootstrap
    ) -> awsmqtt.Connection:
        # Create the connection
        return mqtt_connection_builder.mtls_from_path(
            endpoint=self._config.endpoint,
            cert_filepath=str(self._config.thing_cert.resolve()),
            pri_key_filepath=str(self._config.thing_key.resolve()),
            client_bootstrap=client_bootstrap,
            ca_filepath=str(self._config.root_cert.resolve()),
            client_id=self._config.device_id,
            keep_alive_secs=self._config.keep_alive,
        )

    def _subscribe_to_update_requests(self, mqtt: awsmqtt.Connection) -> None:
        subscribe_future, _ = mqtt.subscribe(
            topic=self._status_request_topic,
            qos=awsmqtt.QoS.AT_LEAST_ONCE,
            callback=self._on_status_requested,
        )
        subscribe_result = subscribe_future.result()
        _LOGGER.debug(
            "Subscribed to '%s' with QOS '%s'",
            subscribe_result["topic"],
            subscribe_result["qos"],
        )

    @property
    def _status_request_topic(self) -> str:
        return f"dags/{self._config.device_id}/status/request"

    @staticmethod
    def _on_status_update_published(_: Awaitable[None]) -> None:
        _LOGGER.debug("Status update published")

    def _on_status_requested(
        self, topic: str, payload: str, **kwargs: Any
    ) -> None:
        del topic, payload, kwargs
        _LOGGER.debug("Status update requested")
        self._status_request_callback()
