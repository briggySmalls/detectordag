from enum import Enum
import json
from typing import Any, Dict

from pydantic import BaseModel
from stringcase import camelcase


class PowerStatus(str, Enum):
    ON = "on"
    OFF = "off"


class DeviceShadowState(BaseModel):
    """Helper function for capturing a device shadow update"""

    status: PowerStatus

    class Config:
        alias_generator = camelcase

    def dict(self, *args, **kwargs) -> Dict[str,Any]:
        """Serialization step"""
        # Wrap up the data into AWS IoT-like structure
        return {'state': {'reported': super().dict(*args, **kwargs)}}
