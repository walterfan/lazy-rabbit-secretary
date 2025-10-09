import { createRouter, createWebHistory } from 'vue-router';
import BooksView from '@/views/BooksView.vue';
import BookmarksView from '@/views/BookmarksView.vue';
import TasksView from '@/views/TasksView.vue';
import RemindersView from '@/views/RemindersView.vue';
import SecretsView from '@/views/SecretsView.vue';
import PostsView from '@/views/PostsView.vue';
import BlogView from '@/views/BlogView.vue';
import BlogPostView from '@/views/BlogPostView.vue';
import HomeView from '@/views/HomeView.vue';
import HelpView from '@/views/HelpView.vue';
import CommandsView from '@/views/CommandsView.vue'
import EncodingToolsView from '@/views/EncodingToolsView.vue';
import PromptsView from '@/views/PromptsView.vue';
import AiAssistantView from '@/views/AiAssistantView.vue';
import ProfileView from '@/views/ProfileView.vue';
import GTDView from '@/views/GTDView.vue';
import WikiView from '@/views/WikiView.vue';
import WikiPageView from '@/views/WikiPageView.vue';
import WikiEditView from '@/views/WikiEditView.vue';
import DiagramView from '@/views/DiagramView.vue';
import DiagramManagementView from '@/views/DiagramManagementView.vue';
import ImageView from '@/views/ImageView.vue';
import SignInView from '@/views/SignInView.vue';
import SignUpView from '@/views/SignUpView.vue';
import UserManagementView from '@/views/admin/UserManagementView.vue';
import RoleManagementView from '@/views/admin/RoleManagementView.vue';
import PermissionManagementView from '@/views/admin/PermissionManagementView.vue';
import RegistrationManagementView from '@/views/admin/RegistrationManagementView.vue';
import RegistrationTestView from '@/views/RegistrationTestView.vue';

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
      path: '/help',
      name: 'help',
      component: HelpView
    },
    {
      path: '/commands',
      name: 'commands',
      component: CommandsView,
      meta: { requiresAuth: true }
    },
    {
      path: '/admin/permissions',
      name: 'permissions',
      component: PermissionManagementView,
      meta: { requiresAuth: true }
    },
    {
      path: '/tools/encoding',
      name: 'encoding-tools',
      component: EncodingToolsView,
      meta: { requiresAuth: true }
    },
    {
      path: '/tools/diagram',
      name: 'diagram',
      component: DiagramView,
      meta: { requiresAuth: true }
    },
    {
      path: '/tools/diagram-management',
      name: 'diagram-management',
      component: DiagramManagementView,
      meta: { requiresAuth: true }
    },
    {
      path: '/tools/images',
      name: 'images',
      component: ImageView,
      meta: { requiresAuth: true }
    },
    {
      path: '/ai-assistant',
      name: 'ai-assistant',
      component: AiAssistantView,
      meta: { requiresAuth: true }
    },
    {
      path: '/prompts',
      name: 'prompts',
      component: PromptsView,
      meta: { requiresAuth: true }
    },
    {
      path: '/profile',
      name: 'profile',
      component: ProfileView,
      meta: { requiresAuth: true }
    },
    {
      path: '/gtd',
      name: 'gtd',
      component: GTDView,
      meta: { requiresAuth: true }
    },
    {
      path: '/wiki',
      name: 'wiki',
      component: WikiView,
      meta: { requiresAuth: false }
    },
    {
      path: '/wiki/page/:slug',
      name: 'wiki-page',
      component: WikiPageView,
      props: true,
      meta: { requiresAuth: false }
    },
    {
      path: '/wiki/edit/:slug?',
      name: 'wiki-edit',
      component: WikiEditView,
      props: true,
      meta: { requiresAuth: true }
    },
    {
      path: '/books',
      name: 'books',
      component: BooksView,
      meta: { requiresAuth: true }
    },
    {
      path: '/bookmarks',
      name: 'bookmarks',
      component: BookmarksView,
      meta: { requiresAuth: true }
    },
    {
      path: '/tasks',
      name: 'tasks',
      component: TasksView
    },
    {
      path: '/reminders',
      name: 'reminders',
      component: RemindersView
    },
    {
      path: '/secrets',
      name: 'secrets',
      component: SecretsView,
      meta: { requiresAuth: true, requiresAdmin: true }
    },
    {
      path: '/posts',
      name: 'posts',
      component: PostsView,
      meta: { requiresAuth: true }
    },
    {
      path: '/blog',
      name: 'blog',
      component: BlogView
    },
    {
      path: '/blog/:slug',
      name: 'blog-post',
      component: BlogPostView,
      props: true
    },
    {
      path: '/blog/archive/:year/:month',
      name: 'blog-archive',
      component: BlogView,
      props: true
    },
    {
      path: '/signin',
      name: 'signin',
      component: SignInView
    },
    {
      path: '/signup',
      name: 'signup',
      component: SignUpView
    },
    {
      path: '/register',
      name: 'register',
      component: SignUpView
    },
    {
      path: '/admin/users',
      name: 'admin-users',
      component: UserManagementView,
      meta: { requiresAuth: true, requiresAdmin: true }
    },
    {
      path: '/admin/roles',
      name: 'admin-roles',
      component: RoleManagementView,
      meta: { requiresAuth: true, requiresAdmin: true }
    },
    {
      path: '/admin/permissions',
      name: 'admin-permissions',
      component: PermissionManagementView,
      meta: { requiresAuth: true, requiresAdmin: true }
    },
    {
      path: '/admin/registrations',
      name: 'admin-registrations',
      component: RegistrationManagementView,
      meta: { requiresAuth: true, requiresAdmin: true }
    },
    {
      path: '/test/registration',
      name: 'test-registration',
      component: RegistrationTestView
    }
  ]
});

// Navigation guard for authentication
router.beforeEach((to, from, next) => {
  const token = localStorage.getItem('access_token');
  const isAuthenticated = !!token;

  if (to.meta.requiresAuth && !isAuthenticated) {
    next('/signin');
  } else if (to.meta.requiresAdmin && isAuthenticated) {
    // For now, we'll allow all authenticated users to access admin pages
    // In a real app, you'd check the user's roles/permissions here
    next();
  } else {
    next();
  }
});

export default router;