import Review from '../views/Review.vue';
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
    // route level code-splitting
    // this generates a separate chunk (about.[hash].js) for this route
    // which is lazy-loaded when the route is visited.
    component: () => import(/* webpackChunkName: "about" */ '../views/Account.vue'),
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
