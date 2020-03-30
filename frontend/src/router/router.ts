import Vue from 'vue';
import VueRouter from 'vue-router';
import routes from './routes';

Vue.use(VueRouter);

// Create the router
const router = new VueRouter({
  mode: 'history',
  base: process.env.BASE_URL,
  routes,
});

// Export the router
export default router;
