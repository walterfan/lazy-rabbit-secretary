import { createApp } from 'vue';
import { createPinia } from 'pinia';
import App from './App.vue';
import router from './router';
import i18n from './i18n';
import { useAuthStore } from './stores/authStore';
import { tokenRefreshService } from './services/tokenRefreshService';

// Import Bootstrap JS and CSS
import 'bootstrap/dist/js/bootstrap.bundle.min.js';

// Import Highlight.js CSS for syntax highlighting
import 'highlight.js/styles/github.css';

const app = createApp(App);
const pinia = createPinia();

app.use(pinia);
app.use(router);
app.use(i18n);

// Initialize auth store
const authStore = useAuthStore();
authStore.initializeAuth();

// Start token refresh monitoring if user is authenticated
if (authStore.isAuthenticated) {
  tokenRefreshService.start();
}

app.mount('#app');