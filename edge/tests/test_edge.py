"""Tests for edge application"""
# pylint: disable=redefined-outer-name

from unittest.mock import Mock, patch
from typing import Generator

import pytest
from pathlib import Path

from edge.config import AppConfig
from edge.data import DeviceShadowState
from edge.edge import EdgeApp
from edge.mocks import MockDigitalInputDevice


@pytest.fixture
def aws() -> Generator[Mock, None, None]:
    """Mock our mqtt client wrapper"""
    with patch("edge.edge.CloudClient", autospec=True) as mock:
        yield mock.return_value


@pytest.fixture
def timer() -> Mock:
    """Mock our periodic timer"""
    with patch("edge.edge.PeriodicTimer", autospec=True) as mock:
        return mock.return_value


@pytest.fixture
def device() -> MockDigitalInputDevice:
    """Mock digital device"""
    device = MockDigitalInputDevice(9)
    # Ensure the device is reading 'low'
    device.low()
    return device


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


def test_setup(config: AppConfig, aws: Mock, device: Mock) -> None:
    """Confirm we can start the application"""
    # Create the unit under test
    with EdgeApp(device, config):
        # Check we immediately send an update
        aws.send_status_update.assert_called_once_with(
            DeviceShadowState(status="off")
        )
