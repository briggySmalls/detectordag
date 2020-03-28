import Vue from 'vue';
import winston from 'winston';
import 'setimmediate'; // Polyfill for winston browser console
import App from './App.vue';
import router from './router';
import store from './store';

Vue.config.productionTip = false;

// Add a logger to all vue instances
Vue.prototype.$logger = winston.createLogger({
  level: (process.env.NODE_ENV !== 'production') ? 'debug' : 'error',
  format: winston.format.simple(),
  transports: [
    new winston.transports.Console(),
  ],
});

new Vue({
  router,
  store,
  render: (h) => h(App),
}).$mount('#app');
