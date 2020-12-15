"""
Tasks for maintaining the project.

Execute 'invoke --list' for guidance on using Invoke
"""
import platform
import webbrowser
from pathlib import Path

from invoke import task  # pylint: disable=no-name-in-module

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


def _run(ctx, command, *args, **kwargs):
    """Helper function for running commands"""
    pty = platform.system() != "Windows"
    ctx.run(command, *args, pty=pty, **kwargs)


def _format_dirs_string(ctx):
    """Get the directories/files to format"""
    return " ".join((str(f) for f in (*ctx.PYTHON_FILES, *ctx.PYTHON_DIRS)))


@task(
    help=dict(check="Checks if source is formatted without applying changes"),
)
def format_black(ctx, check=False):
    """
    Format code
    """
    black_options = "--line-length 79"
    if check:
        black_options += " --check --diff"
    _run(ctx, f"black {black_options} {_format_dirs_string(ctx)}")


@task(
    help=dict(
        check="Checks if imports are formatted without applying changes"
    ),
)
def format_isort(ctx, check=False):
    """
    Format imports
    """
    # Create options
    isort_options = "--profile black --line-length 79"
    if check:
        isort_options += " --check-only --diff"
    # Run the ting
    _run(ctx, f"isort {isort_options} {_format_dirs_string(ctx)}")


@task(
    help=dict(check="Checks if source is formatted without applying changes"),
)  # pylint: disable=redefined-builtin
def format(ctx, check=False):
    """
    Do some clever formatting
    """
    # Format code using black
    format_black(ctx, check)
    # Run isort
    format_isort(ctx, check)


@task
def lint_flake8(ctx):
    """
    Lint code with flake8
    """
    _run(ctx, "flake8 {}".format(_format_dirs_string(ctx)))


@task
def lint_pylint(ctx):
    """
    Lint code with pylint
    """
    _run(ctx, "pylint {}".format(_format_dirs_string(ctx)))


@task
def lint_mypy(ctx):
    """
    Lint code with mypy
    """
    _run(ctx, "mypy --strict --allow-untyped-decorators")


@task(lint_flake8, lint_pylint, lint_mypy)
def lint(_):
    """
    Run all linting
    """


@task
def test(ctx):
    """
    Run tests
    """
    _run(ctx, "python {} test".format(SETUP_FILE))


@task(help={"publish": "Publish the result via coveralls"})
def coverage(ctx, publish=False):
    """
    Create coverage report
    """
    ctx.run("coverage run --source {} -m pytest".format(SOURCE_DIR))
    ctx.run("coverage report")
    if publish:
        # Publish the results via coveralls
        ctx.run("coveralls")
    else:
        # Build a local report
        ctx.run("coverage html")
        webbrowser.open(COVERAGE_REPORT.as_uri())
