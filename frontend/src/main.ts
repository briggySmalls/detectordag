import Vue from 'vue';
import App from './App.vue';
import router from './router';
import store from './store';
import { logger } from './utils';
import { ApiClient } from '../lib/client';

Vue.config.productionTip = false;

// Configure client endpoint
ApiClient.instance.basePath = 'http://localhost:8080/api/v1';

// Add a logger to all vue instances
Vue.prototype.$logger = logger;

new Vue({
  router,
  store,
  render: (h) => h(App),
}).$mount('#app');
