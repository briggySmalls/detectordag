FROM balenalib/raspberry-pi-python:3

# Copy the project
COPY . /package/
WORKDIR /package

# Install package
RUN pip install .

# Set default command
CMD ["edge"]
