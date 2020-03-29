import Vue from 'vue';
import VueRouter from 'vue-router';
import { Response } from 'superagent';
import Review from '../views/Review.vue';
import NotFound from '../views/NotFound.vue';
import { storage, logger, clients } from '../utils';
import { Account } from '../../lib/client';
import store from '../store';

Vue.use(VueRouter);

const routes = [
  {
    path: '/review',
    alias: '/',
    name: 'Review',
    component: Review,
    meta: {
      requiresAuth: true,
    },
  },
  {
    path: '/login',
    name: 'Login',
    // route level code-splitting
    // this generates a separate chunk (about.[hash].js) for this route
    // which is lazy-loaded when the route is visited.
    component: () => import(/* webpackChunkName: "about" */ '../views/Login.vue'),
  },
  {
    path: '*',
    name: 'NotFound',
    component: NotFound,
  },
];

// Create the router
const router = new VueRouter({
  mode: 'history',
  base: process.env.BASE_URL,
  routes,
});

// Save the account details to the store
function handleAccountResponse(error: Error, data: Account, response: Response) {
  // Handle errors
  if (error) {
    logger.debug(`Account request error: ${response.text}`);
    // Clear the token (we're assuming that's why we failed)
    storage.clear();
    // Get the user to reauthenticate
    router.push('/login');
    return;
  }
  // Save the account deatils to the store
  logger.debug('Saving account details');
  store.commit('setAccount', data);
}

// Add guards to ensure we are logged in
router.beforeEach((to, from, next) => {
  // Get auth token
  const authBundle = storage.bundle;
  // Redirect from login if we already have a token
  if (to.name === 'Login' && authBundle !== null) {
    logger.debug('Navigation to login when we already have a token, redirecting...');
    next('/review');
    return;
  }
  // Shortcircuit if we don't need to ensure we're logged in
  if (to.matched.some((record) => !record.meta.requiresAuth)) {
    next();
    return;
  }
  // Redirect to login if we don't have a token
  if (authBundle == null) {
    logger.debug('Token not available');
    next('/login');
    return;
  }
  // Check if we have account details
  if (store.state.account !== null) {
    // We've got everything we need
    next();
    return;
  }
  // Request account details
  logger.debug('Requesting account details');
  clients.accounts.getAccount(`Bearer ${authBundle.token}`, authBundle.accountId, handleAccountResponse);
  next();
});

export default router;
