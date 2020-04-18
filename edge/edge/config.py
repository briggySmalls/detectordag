"""Logic for parsing configuration"""
import base64
from dataclasses import dataclass
from pathlib import Path
from typing import Any, Optional, Dict, List

from environs import Env


class ConfigError(Exception):
    pass


@dataclass
class ConfigMapper:
    identifier: str
    parser: str
    default: Optional[Any] = None


@dataclass
class AppConfig:
    """Class that holds application configuration"""
    _PARSERS = {
        'aws_thing_name': ConfigMapper('AWS_THING_NAME', 'str'),
        'aws_root_cert': ConfigMapper('AWS_ROOT_CERT', 'str'),
        'aws_thing_cert': ConfigMapper('AWS_THING_CERT', 'str'),
        'aws_thing_key': ConfigMapper('AWS_THING_KEY', 'str'),
        'aws_endpoint': ConfigMapper('AWS_ENDPOINT', 'str'),
        'aws_port': ConfigMapper('AWS_PORT', 'int', default=8883),
        'certs_dir': ConfigMapper('CERT_DIR', 'path', '~/.detectordag/certs'),
    }
    _CERTS = {
        'aws_root_cert': 'root-CA.crt',
        'aws_thing_key': 'thing.private.key',
        'aws_thing_cert': 'thing.cert.pem',
    }

    aws_thing_name: str
    aws_root_cert: Path
    aws_thing_cert: Path
    aws_thing_key: Path
    aws_endpoint: str
    aws_port: int
    certs_dir: Path

    @classmethod
    def from_env(cls) -> 'AppConfig':
        """Parse configuration from environment variables

        Returns:
            AppConfig: Application configuration
        """
        env = Env()
        # Read environment variables from .env file (if present)
        env.read_env()
        # Parse our variables
        parsed = {
            name: getattr(env, mapping.parser)(mapping.identifier,
                                               mapping.default)
            for name, mapping in cls._PARSERS.items()
        }
        # Ensure we have all the expected variables
        for key, value in parsed.items():
            if not value:
                raise ConfigError(f"Env variable {key} is missing")
        # Write to a file
        cls._convert_certs(parsed)
        # Return a new config object
        return AppConfig(**parsed)

    @classmethod
    def variables(cls) -> List[str]:
        """Get the variables this config looks for

        Returns:
            List[str]: Identifiers of all variables searched for
        """
        return [mapper.identifier for mapper in cls._PARSERS.values()]

    @classmethod
    def _convert_certs(cls, parsed: Dict[str, Any]) -> None:
        # Save certs to files
        certs_dir = parsed['certs_dir'].expanduser()
        certs_dir.mkdir(exist_ok=True, parents=True)
        for cert, filename in cls._CERTS.items():
            # Establish the path of the new certificate file
            cert_path = certs_dir / filename
            # Create the file from the environment variable
            cls._write_cert(parsed[cert], cert_path)
            # Replace the env variable content with the path to the certificate
            parsed[cert] = cert_path

    @staticmethod
    def _write_cert(cert: str, file: Path) -> None:
        # Turn base64 encoded string into a certificate file
        with file.open('wb') as output_file:
            output_file.write(base64.b64decode(cert))
