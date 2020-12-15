"""
Tasks for maintaining the project.

Execute 'invoke --list' for guidance on using Invoke
"""
import platform
import shutil
import webbrowser
from pathlib import Path

from invoke import task

ROOT_DIR = Path(__file__).parent
SETUP_FILE = ROOT_DIR.joinpath("setup.py")
TEST_DIR = ROOT_DIR.joinpath("tests")
SOURCE_DIR = ROOT_DIR.joinpath("edge")
TOX_DIR = ROOT_DIR.joinpath(".tox")
COVERAGE_FILE = ROOT_DIR.joinpath(".coverage")
COVERAGE_DIR = ROOT_DIR.joinpath("htmlcov")
COVERAGE_REPORT = COVERAGE_DIR.joinpath("index.html")
DOCS_DIR = ROOT_DIR.joinpath("docs")
DOCS_BUILD_DIR = DOCS_DIR.joinpath("_build")
DOCS_INDEX = DOCS_BUILD_DIR.joinpath("index.html")
PYTHON_DIRS = [str(d) for d in [SOURCE_DIR, TEST_DIR]]


def _run(c, command, *args, **kwargs):
    """Helper function for running commands"""
    pty = platform.system() != "Windows"
    c.run(command, *args, pty=pty, **kwargs)


def _format_dirs_string(c):
    """Get the directories/files to format"""
    return " ".join((str(f) for f in (*c.PYTHON_FILES, *c.PYTHON_DIRS)))


@task(
    help=dict(check="Checks if source is formatted without applying changes"),
)
def format_black(c, check=False):
    """
    Format code
    """
    black_options = "--line-length 79"
    if check:
        black_options += " --check --diff"
    _run(c, f"black {black_options} {_format_dirs_string(c)}")


@task(
    help=dict(
        check="Checks if imports are formatted without applying changes"
    ),
)
def format_isort(c, check=False):
    """
    Format imports
    """
    # Create options
    isort_options = "--profile black --line-length 79"
    if check:
        isort_options += " --check-only --diff"
    # Run the ting
    _run(c, f"isort {isort_options} {_format_dirs_string(c)}")


@task(
    help=dict(check="Checks if source is formatted without applying changes"),
)
def format(c, check=False):
    """
    Do some clever formatting
    """
    # Format code using black
    format_black(c, check)
    # Run isort
    format_isort(c, check)


@task
def lint_flake8(c):
    """
    Lint code with flake8
    """
    _run("flake8 {}".format(" ".join(PYTHON_DIRS)))


@task
def lint_pylint(c):
    """
    Lint code with pylint
    """
    _run("pylint {}".format(" ".join(PYTHON_DIRS)))


@task
def lint_mypy(c):
    """
    Lint code with mypy
    """
    _run("mypy --strict --allow-untyped-decorators")


@task(lint_flake8, lint_pylint, lint_mypy)
def lint(c):
    """
    Run all linting
    """
    pass


@task
def test(c):
    """
    Run tests
    """
    pty = platform.system() != "Windows"
    _run("python {} test".format(SETUP_FILE), pty=pty)


@task(help={"publish": "Publish the result via coveralls"})
def coverage(c, publish=False):
    """
    Create coverage report
    """
    c.run("coverage run --source {} -m pytest".format(SOURCE_DIR))
    c.run("coverage report")
    if publish:
        # Publish the results via coveralls
        c.run("coveralls")
    else:
        # Build a local report
        c.run("coverage html")
        webbrowser.open(COVERAGE_REPORT.as_uri())


@task
def docs(c):
    """
    Generate documentation
    """
    c.run("sphinx-build -b html {} {}".format(DOCS_DIR, DOCS_BUILD_DIR))
    webbrowser.open(DOCS_INDEX.as_uri())


@task
def clean_docs(c):
    """
    Clean up files from documentation builds
    """
    c.run("rm -fr {}".format(DOCS_BUILD_DIR))
