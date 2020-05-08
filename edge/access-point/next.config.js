require('dotenv').config()
const env = process.env.NODE_ENV;

module.exports = {
  // Provide runtime configuration
  serverRuntimeConfig: {
    apiBasepath: process.env.DETECTORDAG_API_BASEPATH,
    deviceUUID: process.env.BALENA_DEVICE_UUID,
  },
  // Provide build-time configuration
  env: {
    CERTS_DIR: process.env.CERTS_DIR,
  },
}
