const Dotenv = require('dotenv-webpack');
const env = process.env.NODE_ENV;

module.exports = {
  // Provide runtime configuration
  serverRuntimeConfig: {
    apiBasepath: process.env.DETECTORDAG_API_BASEPATH,
  },
  // Use .env to provide built-time variables
  webpack: (config, { buildId, dev, isServer, defaultLoaders, webpack }) => {
    // Add dotenv
    config.plugins.push(new Dotenv({
      path: `./.env.${env === "production" ? "production" : "dev"}`,
    }));
    return config
  },
}
