
const routes = [
  {
    path: '/',
    component: () => import('layouts/MyLayout.vue'),
    children: [
      {
        path: '',
        component: () => import('pages/Index.vue'),
      },
      {
        path: 'about',
        component: () => import('pages/About.vue'),
      },
      {
        path: 'user/create',
        component: () => import('pages/UserCreate.vue'),
      },
      {
        path: 'user/login',
        component: () => import('pages/UserLogin.vue'),
      },
      {
        path: 'user/notifications',
        component: () => import('pages/UserNotifications.vue'),
      },
      {
        path: 'post/:creatorAddress-:wallAddress-:postHash',
        component: () => import('pages/UserPost.vue'),
      },
      {
        path: 'user/settings',
        component: () => import('pages/UserSettings.vue'),
      },
      {
        path: 'user/:address',
        component: () => import('pages/UserUser.vue'),
      },
    ],
  },
];

// Always leave this as last one
if (process.env.MODE !== 'ssr') {
  routes.push({
    path: '*',
    component: () => import('layouts/MyLayout.vue'),
    children: [
      { path: '', component: () => import('pages/Error404.vue') },
    ],
  });
}

export default routes;
