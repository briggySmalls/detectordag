"""Custom exceptions for the application"""


class DetectorDagException(Exception):
    """Base exception for the application"""


class ConnectionFailedError(DetectorDagException):
    """The MQTT client failed to connect"""
