"""Unit tests for datastructures"""

import json

import pytest

from edge.data import DeviceShadowState, PowerStatus

_DATA = [
    (
        DeviceShadowState(status=True),
        '{"state":{"reported":{"status":"on"}}}',
    ),
    (
        DeviceShadowState(status=False),
        '{"state":{"reported":{"status":"off"}}}',
    ),
    (
        DeviceShadowState(status=1),
        '{"state":{"reported":{"status":"on"}}}',
    ),
    (
        DeviceShadowState(status=0),
        '{"state":{"reported":{"status":"off"}}}',
    ),
    (
        DeviceShadowState(status="on"),
        '{"state":{"reported":{"status":"on"}}}',
    ),
    (
        DeviceShadowState(status="off"),
        '{"state":{"reported":{"status":"off"}}}',
    ),
    (
        DeviceShadowState(status=PowerStatus.ON),
        '{"state":{"reported":{"status":"on"}}}',
    ),
    (
        DeviceShadowState(status=PowerStatus.OFF),
        '{"state":{"reported":{"status":"off"}}}',
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
