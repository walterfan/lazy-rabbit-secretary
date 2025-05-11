import { createRouter, createWebHistory } from 'vue-router';
import BooksView from '@/views/BooksView.vue';
import TasksView from '@/views/TasksView.vue';
import HomeView from '@/views/HomeView.vue';

const router = createRouter({
  history: createWebHistory(),
  routes: [
    {
      path: '/',
      redirect: '/home'
    },
    {
      path: '/home',
      name: 'home',
      component: HomeView
    },
    {
      path: '/books',
      name: 'books',
      component: BooksView
    },
    {
      path: '/tasks',
      name: 'tasks',
      component: TasksView
    }
  ]
});

export default router;