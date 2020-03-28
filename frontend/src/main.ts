import Vue from 'vue';
import App from './App.vue';
import router from './router';
import store from './store';
import { logger } from './utils';

Vue.config.productionTip = false;

// Add a logger to all vue instances
Vue.prototype.$logger = logger;

new Vue({
  router,
  store,
  render: (h) => h(App),
}).$mount('#app');
