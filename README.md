[![detectordag](https://circleci.com/gh/briggySmalls/detectordag.svg?style=shield)](https://circleci.com/gh/briggySmalls/detectordag)

# detectordag

Power outage detector made with ♥ by a dag

# Installation

## Mage

The project generally uses [mage](https://github.com/magefile/mage) for project build commands, which retrospectively was a pain-in-the-arse, but oh well...

```bash
go get -u -d github.com/magefile/mage

go run bootstrap.go
```

## AWS CLI

```bash
# Install the general AWS CLI
brew install awscli
# Install the SAM CLI
brew tap aws/tap
brew install aws-sam-cli
```

You will need to download credentials and save them to `~/.aws/credentials` to interact with AWS.

# Provisioning

## Cloud resources

Most AWS resources are provisioned using the CloudFormation `template.yml`.

```bash
# Build the lambda functions
sam build
# Deploy to AWS
sam deploy
```

## Things

New devices need to be provisioned on a device-by-device basis.

A device exists in two places:
- BalenaCloud for fleet management
- AWS IoT for data-tracking

TODO: document provisioning with mage commands

Some notes:

- IoT things should be of Thing Type "detectordag" in order to ensure they have the expected attributes (e.g. account ID)
- IoT things should be a member of Thing Group "detectordag" in order to give them the correct permissions for updating their shadow
