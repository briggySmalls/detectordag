"""Main module."""
import base64
import logging
import os
from pathlib import Path

from edge.aws import ClientConfig, CloudClient

try:
    from gpiozero import DigitalInputDevice
except ImportError:
    from edge.mocks import MockDigitalInputDevice as DigitalInputDevice  # noqa: E501,  pylint: disable=ungrouped-imports

_LOGGER = logging.getLogger(__name__)

CERT_ROOT_PATH = Path(__file__).parent / 'certs'
CERTS_PATHS = {
    'root_cert': CERT_ROOT_PATH / "root-CA.crt",
    'thing_cert': CERT_ROOT_PATH / "thing.cert.pem",
    'thing_key': CERT_ROOT_PATH / "thing.private.key",
}


class EdgeApp:
    """Wrapper for the entire application"""
    def __init__(self, device: DigitalInputDevice) -> None:
        # Ensure certificates are available
        self._create_certs()
        # Prepare configuration for the client
        config = ClientConfig(device_id=os.getenv("BALENA_DEVICE_UUID"),
                              endpoint=os.getenv("AWS_ENDPOINT"),
                              port=int(os.getenv("AWS_PORT", "8883")),
                              **CERTS_PATHS)
        self.device = device
        # Create the client
        self.client = CloudClient(config)

    def __enter__(self) -> 'EdgeApp':
        self.client.__enter__()
        # Configure the device
        self.configure()
        # Return this instance
        return self

    def configure(self) -> None:
        """Configure the app
        """
        # Send messages when power status changes
        self.device.when_activated = self._publish_update
        self.device.when_deactivated = self._publish_update

    def __exit__(self, exc_type, exc_value, traceback) -> None:
        self.client.__exit__(exc_type, exc_value, traceback)

    @staticmethod
    def _create_certs():
        # Create certificates from environment variables
        CERT_ROOT_PATH.mkdir(parents=True, exist_ok=True)
        EdgeApp._set_cred("AWS_ROOT_CERT", CERTS_PATHS['root_cert'])
        EdgeApp._set_cred("AWS_THING_CERT", CERTS_PATHS['thing_cert'])
        EdgeApp._set_cred("AWS_PRIVATE_CERT", CERTS_PATHS['thing_key'])

    @staticmethod
    def _set_cred(env_name: str, file: Path) -> None:
        # Turn base64 encoded environmental variable into a certificate file
        env = os.getenv(env_name)
        with file.open('wb') as output_file:
            output_file.write(base64.b64decode(env))

    def _publish_update(self, device: DigitalInputDevice) -> None:
        del device
        # Get the status
        status = self.device.value
        # Publish
        self.client.power_status_changed(status)
