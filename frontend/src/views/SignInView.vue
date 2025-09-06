<template>
  <div class="container mt-5">
    <div class="row justify-content-center">
      <div class="col-md-6 col-lg-4">
        <div class="card shadow">
          <div class="card-header bg-primary text-white text-center">
            <h4 class="mb-0">
              <i class="bi bi-box-arrow-in-right"></i> Sign In
            </h4>
          </div>
          <div class="card-body">
            <form @submit.prevent="handleSignIn">
              <div class="mb-3">
                <label for="realmName" class="form-label">Realm Name</label>
                <input
                  type="text"
                  class="form-control"
                  id="realmName"
                  v-model="form.realmName"
                  required
                  placeholder="Enter realm name"
                />
              </div>
              <div class="mb-3">
                <label for="username" class="form-label">Username</label>
                <input
                  type="text"
                  class="form-control"
                  id="username"
                  v-model="form.username"
                  required
                  placeholder="Enter username"
                />
              </div>
              <div class="mb-3">
                <label for="password" class="form-label">Password</label>
                <input
                  type="password"
                  class="form-control"
                  id="password"
                  v-model="form.password"
                  required
                  placeholder="Enter password"
                />
              </div>
              <div class="d-grid">
                <button type="submit" class="btn btn-primary" :disabled="loading">
                  <span v-if="loading" class="spinner-border spinner-border-sm me-2"></span>
                  {{ loading ? 'Signing In...' : 'Sign In' }}
                </button>
              </div>
            </form>
            <div class="mt-3 text-center">
              <p class="mb-0">
                Don't have an account?
                <router-link to="/signup" class="text-decoration-none">Sign Up</router-link>
              </p>
            </div>
          </div>
        </div>
        <div v-if="error" class="alert alert-danger mt-3" role="alert">
          <i class="bi bi-exclamation-triangle"></i> {{ error }}
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref } from 'vue';
import { useRouter } from 'vue-router';
import { useAuthStore } from '@/stores/authStore';

const router = useRouter();
const authStore = useAuthStore();

const form = ref({
  realmName: 'default',
  username: '',
  password: '',
});

const loading = ref(false);
const error = ref('');

const handleSignIn = async () => {
  loading.value = true;
  error.value = '';

  try {
    await authStore.signIn(form.value.username, form.value.password, form.value.realmName);
    router.push('/home');
  } catch (err) {
    error.value = err instanceof Error ? err.message : 'Sign in failed';
  } finally {
    loading.value = false;
  }
};
</script>

<style scoped>
.card {
  border: none;
  border-radius: 10px;
}

.card-header {
  border-radius: 10px 10px 0 0 !important;
}

.form-control:focus {
  border-color: #0d6efd;
  box-shadow: 0 0 0 0.2rem rgba(13, 110, 253, 0.25);
}
</style>
