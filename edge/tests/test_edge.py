from unittest.mock import Mock, call, patch

import pytest

from edge.config import AppConfig
from edge.data import DeviceShadowState
from edge.edge import EdgeApp
from edge.mocks import MockDigitalInputDevice


@pytest.fixture
def aws() -> Mock:
    with patch("edge.edge.CloudClient", autospec=True) as mock:
        yield mock.return_value


@pytest.fixture
def timer() -> None:
    with patch("edge.edge.PeriodicTimer", autospec=True) as mock:
        return mock.return_value


@pytest.fixture
def device() -> None:
    device = MockDigitalInputDevice(9)
    # Ensure the device is reading 'low'
    device.low()
    return device


@pytest.fixture
def config(tmp_path) -> None:
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


def test_setup(config, aws, timer, device) -> None:
    # Create the unit under test
    with EdgeApp(device, config):
        # Check we immediately send an update
        aws.send_status_update.assert_called_once_with(
            DeviceShadowState(status="off")
        )


def test_update(config, aws, timer, device) -> None:
    # Create the unit under test
    with EdgeApp(device, config):
        # Simulate a state change
        device.toggle()
        # Assert expected update was sent
        aws.send_status_update.assert_has_calls(
            [call(DeviceShadowState(status="on"))]
        )
