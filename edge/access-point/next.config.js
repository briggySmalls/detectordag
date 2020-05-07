require('dotenv').config()
const env = process.env.NODE_ENV;

module.exports = {
  // Provide runtime configuration
  serverRuntimeConfig: {
    apiBasepath: process.env.DETECTORDAG_API_BASEPATH,
  },
  // Provide build-time configuration
  env: {
    BALENA_DEVICE_UUID: process.env.BALENA_DEVICE_UUID,
    CERTS_DIR: process.env.CERTS_DIR,
  },
}
