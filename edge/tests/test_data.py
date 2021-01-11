"""Unit tests for datastructures"""

import json

import pytest

from edge.data import DeviceShadowState, PowerStatus

_DATA = [
    (
        DeviceShadowState(status=True),
        '{"status":"on"}',
    ),
    (
        DeviceShadowState(status=False),
        '{"status":"off"}',
    ),
    (
        DeviceShadowState(status=1),
        '{"status":"on"}',
    ),
    (
        DeviceShadowState(status=0),
        '{"status":"off"}',
    ),
    (
        DeviceShadowState(status="on"),
        '{"status":"on"}',
    ),
    (
        DeviceShadowState(status="off"),
        '{"status":"off"}',
    ),
    (
        DeviceShadowState(status=PowerStatus.ON),
        '{"status":"on"}',
    ),
    (
        DeviceShadowState(status=PowerStatus.OFF),
        '{"status":"off"}',
    ),
]


@pytest.mark.parametrize("state,output", _DATA)
def test_serialize(state: DeviceShadowState, output: str) -> None:
    """Tests that we can successfully serialize status updates"""
    # Serialise the state
    assert state.dict() == json.loads(output)


@pytest.mark.parametrize("output,payload", _DATA)
def test_deserialize(payload: str, output: DeviceShadowState):
    assert DeviceShadowState.parse_raw(payload) == output
