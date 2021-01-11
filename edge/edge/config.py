"""Logic for parsing configuration"""
import base64
import os
from pathlib import Path
from typing import Any, Dict, List, Union

from environs import Env, EnvError
from pydantic import BaseModel, validator, ValidationError


class ConfigError(Exception):
    """Exception to indicate that there is a configuration error"""


class ConfigMapper(BaseModel):
    """Helper class for parsing configuration from the environment"""

    identifier: str
    parser: str


def _write_cert(cert: str, file: Path) -> None:
    # Turn base64 encoded string into a certificate file
    with file.open("wb") as output_file:
        output_file.write(base64.b64decode(cert))


def _convert_cert(cls, value: Union[str,Path], field: str, values) -> Path:
    if isinstance(value, Path):
        # Short-circuit, we're not being asked to convert from base64
        return value
    # Pull out the certs dir from the queued variables
    certs_dir = values["certs_dir"].expanduser()
    certs_dir.mkdir(exist_ok=True, parents=True)
    # Establish the path of the new certificate file
    cert_path = certs_dir / cls._certs[field.name]
    # Create the file from the environment variable
    _write_cert(value, cert_path)
    # Replace the env variable content with the path to the certificate
    return cert_path


class AppConfig(BaseModel):
    """Class that holds application configuration"""

    # pylint: disable=too-many-instance-attributes
    _parsers = {
        "aws_thing_name": ConfigMapper(identifier="AWS_THING_NAME", parser="str"),
        "aws_root_cert": ConfigMapper(identifier="AWS_ROOT_CERT", parser="str"),
        "aws_thing_cert": ConfigMapper(identifier="AWS_THING_CERT", parser="str"),
        "aws_thing_key": ConfigMapper(identifier="AWS_THING_KEY", parser="str"),
        "aws_endpoint": ConfigMapper(identifier="AWS_ENDPOINT", parser="str"),
        "certs_dir": ConfigMapper(identifier="CERT_DIR", parser="path"),
        "power_poll_period": ConfigMapper(identifier="POWER_POLL_PERIOD", parser="float"),
    }
    _certs = {
        "aws_root_cert": "root-CA.crt",
        "aws_thing_key": "thing.private.key",
        "aws_thing_cert": "thing.cert.pem",
    }

    certs_dir: Path = Path("~/.detectordag/certs")
    aws_thing_name: str
    aws_root_cert: Path
    aws_thing_cert: Path
    aws_thing_key: Path
    aws_endpoint: str
    power_poll_period: float = 60.0
    keep_alive_period: int = 1200

    _convert_thing_cert = validator("aws_thing_cert", pre=True, allow_reuse=True)(_convert_cert)
    _convert_thing_key = validator("aws_thing_key", pre=True, allow_reuse=True)(_convert_cert)
    _convert_root_cert = validator("aws_root_cert", pre=True, allow_reuse=True)(_convert_cert)

    @classmethod
    def from_env(cls, dotenv: bool = True) -> "AppConfig":
        """Parse configuration from environment variables

        Returns:
            AppConfig: Application configuration
        """
        env = Env()
        if dotenv:
            env.read_env(str(Path(os.getcwd()) / ".env"))
        # Parse our variables
        parsed = {}
        for name, mapping in cls._parsers.items():
            # This may fail if env vars are not present
            try:
                # Parse a variable without a default
                parsed[name] = getattr(env, mapping.parser)(
                    mapping.identifier
                )
            except EnvError:
                # Assume we've left this variable out intentionally
                # The AppConfig constructor performs validation
                pass
        # Return a new config object
        try:
            return AppConfig(**parsed)
        except ValidationError as exc:
            # Validation of fields failed
            raise ConfigError() from exc

    @classmethod
    def variables(cls) -> List[str]:
        """Get the variables this config looks for

        Returns:
            List[str]: Identifiers of all variables searched for
        """
        return [mapper.identifier for mapper in cls._parsers.values()]
