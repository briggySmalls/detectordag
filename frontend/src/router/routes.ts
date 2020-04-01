import Review from '../views/Review.vue';
import Login from '../views/Login.vue';
import Account from '../views/Account.vue';
import NotFound from '../views/NotFound.vue';

export default [
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
    path: '/account',
    alias: '/',
    name: 'Account',
    component: Account,
    meta: {
      requiresAuth: true,
    },
  },
  {
    path: '/login',
    name: 'Login',
    component: Login,
  },
  {
    path: '*',
    name: 'NotFound',
    component: NotFound,
  },
];
