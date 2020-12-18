# frontend

Vuejs application for the detector dag software.

## Project setup
```
yarn install
```

### Compiles and hot-reloads for development

You will need a development API running for this to work.

You have two options for this:

1. Running a mocked API that returns very basic responses
2. Running a real API locally that actually interfaces with AWS services

Once this is running, you can start your development server with one of the following commands:

```
# Mocked server
yarn run serve
# Local API
yarn run serve --mode development-sam
```

#### Mocked API

Run a mocked API from the OpenAPI spec:

```
mage generate:spec mockApi
```

#### Local API

Run a local instance of the API:

```
sam build
sam local start-api
```

### Compiles and minifies for production
```
yarn run build
```

### Run your tests
```
yarn run test
```

### Lints and fixes files
```
yarn run lint
```

### Run your end-to-end tests
```
yarn run test:e2e
```

### Run your unit tests
```
yarn run test:unit
```

### Deploy built application
```
yarn deploy
```

### Customize configuration
See [Configuration Reference](https://cli.vuejs.org/config/).
