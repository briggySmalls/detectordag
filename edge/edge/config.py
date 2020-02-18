"""Logic for parsing configuration"""
from environs import Env
from typing import Any, Optional
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
        'aws_private_cert': ConfigMapper('AWS_PRIVATE_CERT', 'str'),
        'aws_endpoint': ConfigMapper('AWS_ENDPOINT', 'str'),
        'aws_port': ConfigMapper('AWS_PORT', 'int', default=8883),
    }

    balena_device_id: str
    aws_root_cert: str
    aws_thing_cert: str
    aws_private_cert: str
    aws_endpoint: str
    aws_port: int

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
        # Return a new config object
        return AppConfig(**parsed)
