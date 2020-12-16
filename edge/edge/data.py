from dataclasses import dataclass
from enum import Enum
import json


class PowerStatus(Enum):
    ON
    OFF


@dataclass
class DeviceShadowState:
    """Helper function for capturing a device shadow update"""

    status: bool

    def to_json(self) -> str:
        """Convert shadow state to an AWS shadow JSON payload

        Returns:
            str: AWS shadow JSON payload
        """
        payload = {'state': {'reported': asdict(self)}}
        return json.dumps(payload)

