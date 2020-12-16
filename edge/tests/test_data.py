"""Unit tests for datastructures"""

import json

import pytest

from edge.data import DeviceShadowState, PowerStatus


@pytest.mark.parametrize(
    "state,output",
    [
        (
            DeviceShadowState(status=PowerStatus.ON),
            '{"state":{"reported":{"status":"on"}}}',
        ),
        (
            DeviceShadowState(status=PowerStatus.OFF),
            '{"state":{"reported":{"status":"off"}}}',
        ),
    ],
)
def test_serialize(state: DeviceShadowState, output: str) -> None:
    """Tests that we can successfully serialize status updates"""
    # Serialise the state
    assert state.dict() == json.loads(output)