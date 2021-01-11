# Start with a builder image, this can have pip on it
FROM balenalib/raspberry-pi-python:3.7-buster-build AS builder

# Install wheels from piwheels (awscrt from source is ridiculous)
RUN printf "[global]\nextra-index-url=https://www.piwheels.org/simple\n" > /etc/pip.conf

# Install poetry
ENV POETRY_VERSION=1.0.10
RUN curl -sSL https://raw.githubusercontent.com/python-poetry/poetry/master/get-poetry.py | python
RUN ~/.poetry/bin/poetry --version

# First just install the packages
WORKDIR /app
COPY pyproject.toml poetry.lock /app/
RUN ~/.poetry/bin/poetry install --no-dev

# Copy the rest of the package source
COPY . /app

# Project build
RUN ~/.poetry/bin/poetry build -f wheel

# Build the real image
FROM balenalib/raspberry-pi-python:3.7-buster-run

# Install wheels from piwheels (awscrt from source is ridiculous)
RUN printf "[global]\nextra-index-url=https://www.piwheels.org/simple\n" > /etc/pip.conf

# Install GPIO stuff
RUN apt-get update && apt-get install python3-gpiozero python3-rpi.gpio

# Copy and install the build
COPY --from=builder /app/dist/*.whl /tmp
RUN pip install /tmp/*.whl

# Set default command
CMD ["edge", "app"]
