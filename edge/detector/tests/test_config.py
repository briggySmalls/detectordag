"""Tests for `edge` package."""
# pylint: disable=redefined-outer-name
from shutil import copyfile
from pathlib import Path
from typing import Any, Dict
from dataclasses import dataclass

import pytest

from edge.config import AppConfig, ConfigError


_TEST_DATA_DIR = Path(__file__).parent.joinpath('test_data')


@dataclass
class EnvVarProps:
    value: Any
    optional: bool


@pytest.fixture
def variables(monkeypatch: Any, tmp_path: Path) -> Dict[str, Any]:
    # Establish the paths for the thing certificates
    certs_dir = tmp_path
    aws_thing_cert_path = certs_dir.joinpath('thing.cert.pem')
    aws_thing_key_path = certs_dir.joinpath('thing.private.key')
    # Copy some test certificates to the paths
    copyfile(_TEST_DATA_DIR.joinpath('thing.cert.pem'), aws_thing_cert_path)
    copyfile(_TEST_DATA_DIR.joinpath('thing.private.key'), aws_thing_key_path)
    # Create the variables
    variables = {
        "AWS_ENDPOINT": "www.test.com",
        "BALENA_DEVICE_UUID": "03076b52-3a66-425b-ad2d-43e925486e60",
        "AWS_ROOT_CERT": "LS0tLS2CRUdJTiBDRVJUSUZJQ0FURS0tqS0tCk1JSURRVENDQWltZ0F3SUJBZ0lUQm15Zno1bS9qQW81NHZCNGlrUG1salpieWpBTkJna3Foa2lHOXcwQkFRc0YKQURBNU1Rc3dDUVlEVlFRR0V3SlZVekVQTUEwR0ExVUVDaE1HUVcxaGVtOXVNUmt3RndZRFZRUURFeEJCYldGNgpiMjRnVW05dmRDQkRRU0F4TUI0WERURTFNRFV5TmpBd01EQXdNRm9YRFRNNE1ERXhOekF3TURBd01Gb3dPVEVMCk1Ba0dBMVVFQmhNQ1ZWTXhEekFOQmdOVkJBb1RCa0Z0WVhwdmJqRVpNQmNHQTFVRUF4TVFRVzFoZW05dUlGSnYKYjNRZ1EwRWdNVENDQVNJd0RRWUpLb1pJaHZjTkFRRUJCUUFEZ2dFUEFEQ0NBUW9DZ2dFQkFMSjRnSEhLZU5YagpjYTlIZ0ZCMGZXN1kxNGgyOUpsbzkxZ2hZUGwwaEFFdnJBSXRodE9nUTNwT3NxVFFOcm9Cdm8zYlNNZ0hGelpNCjlPNklJOGMrNnpmMXRSbjRTV2l3M3RlNWRqZ2RZWjZrL29JMnBlVktWdVJGNGZuOXRCYjZkTnFjbXpVNUwvcXcKSUZBR2JIclFnTEttK2Evc1J4bVBVRGdIM0tLSE9WajR1dFdwK1Vobk1KYnVsSGhlYjRtalVjQXdobWFoUldhNgpWT3VqdzVINVNOei8wZWd3TFgwdGRIQTExNGdrOTU3RVdXNjdjNGNYOGpKR0tMaEQrcmNkcXNxMDhwOGtEaTFMCjkzRmNYbW4vNnBVQ3l6aUtybEE0Yjl2N0xXSWJ4Y2NlVk9GMzRHZklENXlISTlZL1FDQi9JSURFZ0V3K095UW0KamdTdWJKcklxZzBDQXdFQUFhTkNNRUF3RHdZRFZSMFRBUUgvQkFVd0F3RUIvekFPQmdOVkhROEJBZjhFQkFNQwpBWVl3SFFZRFZSME9CQllFRklRWXpJVTA3THdNbEpRdUNGbWN4N0lRVGdvSU1BMEdDU3FHU0liM0RRRUJDd1VBCkE0SUJBUUNZOGpkYVFaQ2hHc1YyVVNnZ05pTU9ydVlvdTZyNGxLNUlwREIvRy93a2pVdTB5S0dYOXJieGVuREkKVTVQTUNDamptQ1hQSTZUNTNpSFRmSVVKclU2YWRUckNDMnFKZUhaRVJ4aGxiSTFCamp0L21zdjB0YWRRMXdVcwpOK2dEUzYzcFlhQUNidlh5OE1XeTdWdTMzUHFVWEhlZUU2Vi9VcTJWOHZpVE85NkxYRnZLV2xKYllLOFU5MHZ2Cm8vdWZRSlZ0TVZUOFF0UEhSaDhqcmRrUFNIQ2EyWFY0Y2RGeVF6UjFibGRad2dKY0ptQXB6eU1aRm82SVE2WFUKNU1zSSt5TVJRK2hES1hKaW9hbGRYZ2pVa0s2NDJNNFV3dEJWOG9iMnhKTkRkMlpod0xub1FkZVhlR0FEYmtweQpycVhSZmJvUW5vWnNHNHE1V1RQNDY4U1F2dkc1Ci0tLS0tRU5EIENFUlRJRklDQVRFLS0tLS0K",  # noqa: E501 pylint: disable=line-too-long
        "AWS_THING_CERT_PATH": aws_thing_cert_path.resolve(),
        "AWS_THING_KEY_PATH": aws_thing_key_path.resolve(),
        "AWS_PORT": 8080,
        "ALIVE_INTERVAL": 3600,
    }
    # Delete existing environment variables
    for key in variables.keys():
        monkeypatch.delenv(key, raising=False)
    return variables


def test_present(monkeypatch: Any, variables: Dict[str, EnvVarProps]) -> None:
    """Test 'happy path' of all variables present

    Args:
        monkeypatch (TYPE): Fixture for configuring environment
        tmp_path (TYPE): Fixture for supplying a temporary directory
    """
    # Configure the environment variables
    for name, value in variables.items():
        monkeypatch.setenv(name, str(value))
    # Create the config
    config = AppConfig.from_env()
    # Assert values
    assert config.aws_endpoint == variables["AWS_ENDPOINT"]
    assert config.aws_root_cert == variables["AWS_ROOT_CERT"]
    assert config.aws_thing_cert == variables["AWS_THING_CERT_PATH"]
    assert config.aws_thing_key == variables["AWS_THING_KEY_PATH"]
    assert config.aws_port == variables["AWS_PORT"]
    assert config.aws_thing_name == variables["AWS_THING_NAME"]
    assert config.alive_interval == variables["ALIVE_INTERVAL"]


@pytest.mark.parametrize("to_drop", [
    "AWS_ENDPOINT",
    "AWS_ROOT_CERT",
    "AWS_THING_CERT_PATH",
    "AWS_THING_KEY_PATH",
    "BALENA_DEVICE_UUID",
])
def test_missing_env(monkeypatch: Any, variables: Dict[str, EnvVarProps], to_drop: str) -> None:
    """Run a test, dropping each of the mandatory variables in turn"""
    # Set all variables, dropping the variable under test
    for key, value in variables.items():
        if key != to_drop:
            # Set the variable in the environment
            monkeypatch.setenv(key, str(value))
    # Expect an error
    with pytest.raises(ConfigError):
        AppConfig.from_env()
