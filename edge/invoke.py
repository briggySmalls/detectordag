"""Configuration file for invoke recipes"""

from pathlib import Path

ROOT_DIR = Path(__file__).parent

SRC_DIR = ROOT_DIR / "edge"

PYTHON_DIRS = [SRC_DIR]
PYTHON_FILES = [ROOT_DIR / "tasks.py", ROOT_DIR / "invoke.py"]
