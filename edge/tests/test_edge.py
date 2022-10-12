"""Tests for edge application"""
# pylint: disable=redefined-outer-name

from pathlib import Path
from typing import Generator, cast
from unittest.mock import Mock, patch

import pytest

from edge.config import AppConfig
from edge.data import DeviceShadowState
from edge.edge import EdgeApp
from edge.mocks import MockPower


@pytest.fixture
def aws() -> Generator[Mock, None, None]:
    """Mock our mqtt client wrapper"""
    with patch("edge.edge.CloudClient", autospec=True) as mock:
        yield mock.return_value


@pytest.fixture
def timer() -> Mock:
    """Mock our periodic timer"""
    with patch("edge.edge.PeriodicTimer", autospec=True) as mock:
        return cast(Mock, mock.return_value)


@pytest.fixture
def power() -> MockPower:
    """Mock digital power"""
    with patch("edge.ina219.INA219", autospec=True) as mock:
        power = MockPower(mock)
        # Ensure the power is reading 'low'
        power.low()
        return power


@pytest.fixture
def config(tmp_path: Path) -> AppConfig:
    """Create a configuration for the tests"""
    return AppConfig(
        aws_thing_name="",
        aws_root_cert=tmp_path,
        aws_thing_cert=tmp_path,
        aws_thing_key=tmp_path,
        aws_endpoint="",
        aws_port=1,
        certs_dir=tmp_path,
        power_poll_period=10.0,
    )


def test_setup(config: AppConfig, aws: Mock, power: Mock) -> None:
    """Confirm we can start the application"""
    # Create the unit under test
    with EdgeApp(power, config):
        # Check we immediately send an update
        aws.send_status_update.assert_called_once_with(
            DeviceShadowState(status="off")
        )
