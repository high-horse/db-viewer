import { createRouter, createWebHashHistory, RouteRecordRaw } from 'vue-router';

import WelcomePage from './pages/WelcomePage.vue';
import WorkSpacePage from './pages/WorkSpacePage.vue';
import AuthLayout from './layouts/AuthLayout.vue';
import WorkspaceLayout from './layouts/WorkspaceLayout.vue';

const routes: RouteRecordRaw[] = [
  {
    path: '/',
    component: AuthLayout,
    children: [
      {
        path: '',
        component: WelcomePage,
        name: 'Welcome',
      },
    ],
  },
  {
    path: '/workspace',
    component: WorkspaceLayout,
    name: 'WorkSpace',
    // children: [
    //   {
    //     path: '',
    //     component: WorkSpacePage,
    //     name: 'WorkSpace',
    //   },
    // ],
  },
];

const router = createRouter({
  history: createWebHashHistory(),
  routes,
});

export default router;
