import { createRouter, createWebHistory, RouteRecordRaw } from 'vue-router';
import QueryExecutor from '../views/QueryExecutor.vue';

const routes: Array<RouteRecordRaw> = [
  {
    path: '/',
    name: 'QueryExecutor',
    component: QueryExecutor,
  },
  {
    path: '/:pathMatch(.*)*',
    redirect: '/',
  },
];

const router = createRouter({
  history: createWebHistory(process.env.BASE_URL),
  routes,
});

export default router;
