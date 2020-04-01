module.exports = {
  // Set the title
  chainWebpack: config => {
    config
    .plugin('html')
    .tap(args => {
      args[0].title = 'Detector Dag'
      return args
    })
  },

  configureWebpack: {
    devtool: 'source-map',
  },

  devServer: {
    proxy: {
      '/api': {
        target: 'http://localhost:3000',
        pathRewrite: {
          '^/api': '/', // remove base path
        },
      },
    },
  },

  pluginOptions: {
    s3Deploy: {
      registry: undefined,
      awsProfile: 'default',
      overrideEndpoint: false,
      region: 'eu-west-2',
      bucket: 'detectordag-frontend',
      createBucket: true,
      staticHosting: true,
      staticIndexPage: 'index.html',
      staticErrorPage: 'index.html',
      assetPath: 'dist',
      assetMatch: '**',
      deployPath: '/',
      acl: 'public-read',
      pwa: false,
      enableCloudfront: true,
      cloudfrontId: 'E2PZG09FJQNGCG',
      pluginVersion: '4.0.0-rc3',
      uploadConcurrency: 5
    }
  }
};
