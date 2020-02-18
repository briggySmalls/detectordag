"""Logic for parsing configuration"""
from environs import Env
import base64
from typing import Any, Optional
from pathlib import Path
from dataclasses import dataclass


@dataclass
class ConfigMapper:
    identifier: str
    parser: str
    default: Optional[Any] = None


@dataclass
class AppConfig:
    """Class that holds application configuration"""
    _PARSERS = {
        'balena_device_id': ConfigMapper('BALENA_DEVICE_UUID', 'uuid'),
        'aws_root_cert': ConfigMapper('AWS_ROOT_CERT', 'str'),
        'aws_thing_cert': ConfigMapper('AWS_THING_CERT', 'str'),
        'aws_thing_key': ConfigMapper('AWS_THING_KEY', 'str'),
        'aws_endpoint': ConfigMapper('AWS_ENDPOINT', 'str'),
        'aws_port': ConfigMapper('AWS_PORT', 'int', default=8883),
        'certs_dir': ConfigMapper('CERT_DIR', 'path', '~/.detectordag/certs'),
    }
    _CERTS = {
        'aws_root_cert': 'root-CA.crt',
        'aws_thing_key': 'thing.cert.pem',
        'aws_thing_cert': 'thing.private.key',
    }

    balena_device_id: str
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
            name: getattr(env, mapping.parser)(mapping.identifier, mapping.default)
            for name, mapping in cls._PARSERS.items()
        }
        # Save certs to files
        parsed['certs_dir'].mkdir(exist_ok=True, parents=True)
        for cert, filename in cls._CERTS.items():
            # Establish the path of the new certificate file
            cert_path = parsed['certs_dir'] / filename
            # Create the file from the environment variable
            cls._write_cert(parsed[cert], cert_path)
            # Replace the env variable content with the path to the certificate
            parsed[cert] = cert_path
        # Return a new config object
        return AppConfig(**parsed)

    @staticmethod
    def _write_certs(cert_dir: Path, root_cert: str, thing_cert: str, thing_key: str) -> None:
        AppConfig._write_cert(root_cert, cert_dir / "root-CA.crt")
        AppConfig._write_cert(thing_cert, cert_dir / "thing.cert.pem")
        AppConfig._write_cert(thing_key, cert_dir / "thing.private.key")

    @staticmethod
    def _write_cert(cert: str, file: Path) -> None:
        # Turn base64 encoded string into a certificate file
        with file.open('wb') as output_file:
            output_file.write(base64.b64decode(cert))
