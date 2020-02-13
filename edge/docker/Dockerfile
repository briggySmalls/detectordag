# We first, for layer caching, generate a requirements.txt
FROM python:3
# Install pipenv
RUN pip install pipenv
# Copy pipfiles
COPY Pipfile Pipfile.lock /package/
WORKDIR /package/
# Create a requirements.txt
RUN pipenv lock -r

# Next we begin the real image, first installing requirements
FROM balenalib/%%BALENA_MACHINE_NAME%%-python:3
# Copy the requirements file
COPY --from=0 /package/requirements.txt /app/
WORKDIR /app/
# Install dependencies
RUN pip install -r requirements.txt
# Copy the rest of the project
COPY . /package/
WORKDIR /package
# Install package
RUN pip install .
# Set default command
CMD ["edge"]
