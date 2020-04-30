"""Tests for `edge` package."""
# pylint: disable=redefined-outer-name

from pathlib import Path
from typing import Any, List
from shutil import copyfile
from dataclasses import dataclass

import pytest

from edge.config import AppConfig, ConfigError


_TEST_DATA_DIR = Path(__file__).parent.joinpath('test_data')


@dataclass
class EnvironmentVariable:
    key: str
    value: str
    optional: bool


_VARIABLES = [
    EnvironmentVariable("AWS_ENDPOINT", "www.test.com", False),
    EnvironmentVariable("BALENA_DEVICE_UUID", "03076b52-3a66-425b-ad2d-43e925486e60", False),
    EnvironmentVariable(
        "AWS_ROOT_CERT",
        "LS0tLS2CRUdJTiBDRVJUSUZJQ0FURS0tqS0tCk1JSURRVENDQWltZ0F3SUJBZ0lUQm15Zno1bS9qQW81NHZCNGlrUG1salpieWpBTkJna3Foa2lHOXcwQkFRc0YKQURBNU1Rc3dDUVlEVlFRR0V3SlZVekVQTUEwR0ExVUVDaE1HUVcxaGVtOXVNUmt3RndZRFZRUURFeEJCYldGNgpiMjRnVW05dmRDQkRRU0F4TUI0WERURTFNRFV5TmpBd01EQXdNRm9YRFRNNE1ERXhOekF3TURBd01Gb3dPVEVMCk1Ba0dBMVVFQmhNQ1ZWTXhEekFOQmdOVkJBb1RCa0Z0WVhwdmJqRVpNQmNHQTFVRUF4TVFRVzFoZW05dUlGSnYKYjNRZ1EwRWdNVENDQVNJd0RRWUpLb1pJaHZjTkFRRUJCUUFEZ2dFUEFEQ0NBUW9DZ2dFQkFMSjRnSEhLZU5YagpjYTlIZ0ZCMGZXN1kxNGgyOUpsbzkxZ2hZUGwwaEFFdnJBSXRodE9nUTNwT3NxVFFOcm9Cdm8zYlNNZ0hGelpNCjlPNklJOGMrNnpmMXRSbjRTV2l3M3RlNWRqZ2RZWjZrL29JMnBlVktWdVJGNGZuOXRCYjZkTnFjbXpVNUwvcXcKSUZBR2JIclFnTEttK2Evc1J4bVBVRGdIM0tLSE9WajR1dFdwK1Vobk1KYnVsSGhlYjRtalVjQXdobWFoUldhNgpWT3VqdzVINVNOei8wZWd3TFgwdGRIQTExNGdrOTU3RVdXNjdjNGNYOGpKR0tMaEQrcmNkcXNxMDhwOGtEaTFMCjkzRmNYbW4vNnBVQ3l6aUtybEE0Yjl2N0xXSWJ4Y2NlVk9GMzRHZklENXlISTlZL1FDQi9JSURFZ0V3K095UW0KamdTdWJKcklxZzBDQXdFQUFhTkNNRUF3RHdZRFZSMFRBUUgvQkFVd0F3RUIvekFPQmdOVkhROEJBZjhFQkFNQwpBWVl3SFFZRFZSME9CQllFRklRWXpJVTA3THdNbEpRdUNGbWN4N0lRVGdvSU1BMEdDU3FHU0liM0RRRUJDd1VBCkE0SUJBUUNZOGpkYVFaQ2hHc1YyVVNnZ05pTU9ydVlvdTZyNGxLNUlwREIvRy93a2pVdTB5S0dYOXJieGVuREkKVTVQTUNDamptQ1hQSTZUNTNpSFRmSVVKclU2YWRUckNDMnFKZUhaRVJ4aGxiSTFCamp0L21zdjB0YWRRMXdVcwpOK2dEUzYzcFlhQUNidlh5OE1XeTdWdTMzUHFVWEhlZUU2Vi9VcTJWOHZpVE85NkxYRnZLV2xKYllLOFU5MHZ2Cm8vdWZRSlZ0TVZUOFF0UEhSaDhqcmRrUFNIQ2EyWFY0Y2RGeVF6UjFibGRad2dKY0ptQXB6eU1aRm82SVE2WFUKNU1zSSt5TVJRK2hES1hKaW9hbGRYZ2pVa0s2NDJNNFV3dEJWOG9iMnhKTkRkMlpod0xub1FkZVhlR0FEYmtweQpycVhSZmJvUW5vWnNHNHE1V1RQNDY4U1F2dkc1Ci0tLS0tRU5EIENFUlRJRklDQVRFLS0tLS0K",  # noqa: E501 pylint: disable=line-too-long
        False),
    EnvironmentVariable("AWS_PORT", "8080", True),
    EnvironmentVariable("ALIVE_INTERVAL", "3600", True),
    EnvironmentVariable("CERT_DIR", "./test/dir", True),
]


@pytest.fixture
def variables() -> List[EnvironmentVariable]:
    return _VARIABLES


def test_variables(variables: List[EnvironmentVariable]) -> None:
    """Test that config is looking for expected variables"""

    _variable_ids = [var.key for var in variables]
    assert set(_variable_ids) == set(AppConfig.variables())


def test_present(monkeypatch: Any, tmp_path: Path, variables: List[EnvironmentVariable]) -> None:
    """Test 'happy path' of all variables present

    Args:
        monkeypatch (TYPE): Fixture for configuring environment
        tmp_path (TYPE): Fixture for supplying a temporary directory
    """
    aws_endpoint = "www.test.com"
    aws_thing_name = "03076b52-3a66-425b-ad2d-43e925486e60"
    alive_interval = 3600
    certs_dir = tmp_path
    # Delete existing environment variables
    for var in variables:
        monkeypatch.delenv(var.key, raising=False)
    # Configure the environment variables
    for var in variables:
        monkeypatch.setenv(var.key, var.value)
    monkeypatch.setenv("CERT_DIR", str(certs_dir))
    # Create certificates in the folder
    aws_thing_cert_path = certs_dir.joinpath('thing.cert.pem')
    copyfile(_TEST_DATA_DIR.joinpath('thing.cert.pem'), aws_thing_cert_path)
    aws_thing_key_path = certs_dir.joinpath('thing.private.key')
    copyfile(_TEST_DATA_DIR.joinpath('thing.private.key'), aws_thing_key_path)
    # Create the config
    config = AppConfig.from_env()
    # Assert root certificate is created from environment variable
    aws_root_cert_path = tmp_path / "root-CA.crt"
    aws_root_cert_path.exists()
    # Assert values
    assert config.aws_endpoint == aws_endpoint
    assert config.certs_dir == certs_dir
    assert config.aws_root_cert == aws_root_cert_path
    assert config.aws_thing_cert == aws_thing_cert_path
    assert config.aws_thing_key == aws_thing_key_path
    assert config.aws_port == 8883
    assert config.aws_thing_name == aws_thing_name
    assert config.alive_interval == alive_interval


@pytest.mark.parametrize("to_drop", [var.key for var in _VARIABLES if not var.optional])
def test_missing_env(monkeypatch: Any, tmp_path: Path, variables: List[EnvironmentVariable], to_drop: bool) -> None:
    """Run a test, dropping each of the variables in turn"""
    # Delete existing environment variables
    for var in variables:
        monkeypatch.delenv(var.key, raising=False)
    # Now set all variables, skipping variable under test
    for var in variables:
        if var.key != to_drop:
            # Set the variable in the environment
            monkeypatch.setenv(var.key, var.value)
    # Always set the cert dir
    monkeypatch.setenv("CERT_DIR", str(tmp_path))
    # Expect an error
    with pytest.raises(ConfigError):
        AppConfig.from_env()
