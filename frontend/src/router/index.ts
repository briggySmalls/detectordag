import router from './router';
import { storage, logger } from '../utils';

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
  // Proceed with the original plan
  next();
});

export default router;
