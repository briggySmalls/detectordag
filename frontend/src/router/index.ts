import router from './router';
import { storage, logger } from '../utils';

// Add guards to ensure we are logged in
router.beforeEach((to, from, next) => {
  logger.debug('running guards...');
  // Redirect to login if we don't have a token and need one
  const authBundle = storage.bundle;
  if (authBundle == null && to.matched.some((record) => record.meta.requiresAuth)) {
    logger.debug('Token not available');
    next('/login');
    return;
  }
  // Redirect from login if we already have a token
  if (to.name === 'Login' && authBundle !== null) {
    logger.debug('Navigation to login when we already have a token, redirecting...');
    next('/review');
    return;
  }
  // Proceed with the original plan
  next();
});

export default router;
