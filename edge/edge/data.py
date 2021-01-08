"""Module for basic data structures"""

from enum import Enum
from typing import Any, Dict

from pydantic import BaseModel, validator
from stringcase import camelcase


class PowerStatus(str, Enum):
    """Enumeration of different power statuses"""

    ON = "on"
    OFF = "off"


class DeviceShadowState(BaseModel):  # pylint: disable=too-few-public-methods
    """Helper function for capturing a device shadow update"""

    status: PowerStatus

    class Config:  # pylint: disable=too-few-public-methods
        """Configuration for the pydantic model"""

        alias_generator = camelcase

    @validator('status', pre=True)
    def to_status(cls, status: bool) -> PowerStatus:
        return PowerStatus.ON if status else PowerStatus.OFF

    def dict(self, *args: Any, **kwargs: Any) -> Dict[str, Any]:
        """Serialization step"""
        # Wrap up the data into AWS IoT-like structure
        return {"state": {"reported": super().dict(*args, **kwargs)}}
