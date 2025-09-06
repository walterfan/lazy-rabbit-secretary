<template>
  <div class="container mt-5">
    <div class="row justify-content-center">
      <div class="col-md-6 col-lg-4">
        <div class="card shadow">
          <div class="card-header bg-success text-white text-center">
            <h4 class="mb-0">
              <i class="bi bi-person-plus"></i> Sign Up
            </h4>
          </div>
          <div class="card-body">
            <form @submit.prevent="handleSignUp">
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
                <label for="email" class="form-label">Email</label>
                <input
                  type="email"
                  class="form-control"
                  id="email"
                  v-model="form.email"
                  required
                  placeholder="Enter email"
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
              <div class="mb-3">
                <label for="confirmPassword" class="form-label">Confirm Password</label>
                <input
                  type="password"
                  class="form-control"
                  id="confirmPassword"
                  v-model="form.confirmPassword"
                  required
                  placeholder="Confirm password"
                />
              </div>
              <div class="d-grid">
                <button type="submit" class="btn btn-success" :disabled="loading || !passwordsMatch">
                  <span v-if="loading" class="spinner-border spinner-border-sm me-2"></span>
                  {{ loading ? 'Creating Account...' : 'Sign Up' }}
                </button>
              </div>
            </form>
            <div class="mt-3 text-center">
              <p class="mb-0">
                Already have an account?
                <router-link to="/signin" class="text-decoration-none">Sign In</router-link>
              </p>
            </div>
          </div>
        </div>
        <div v-if="error" class="alert alert-danger mt-3" role="alert">
          <i class="bi bi-exclamation-triangle"></i> {{ error }}
        </div>
        <div v-if="!passwordsMatch && form.confirmPassword" class="alert alert-warning mt-3" role="alert">
          <i class="bi bi-exclamation-triangle"></i> Passwords do not match
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, computed } from 'vue';
import { useRouter } from 'vue-router';
import { useAuthStore } from '@/stores/authStore';

const router = useRouter();
const authStore = useAuthStore();

const form = ref({
  realmName: 'default',
  username: '',
  email: '',
  password: '',
  confirmPassword: '',
});

const loading = ref(false);
const error = ref('');

const passwordsMatch = computed(() => form.value.password === form.value.confirmPassword);

const handleSignUp = async () => {
  if (!passwordsMatch.value) {
    error.value = 'Passwords do not match';
    return;
  }

  loading.value = true;
  error.value = '';

  try {
    await authStore.signUp({
      username: form.value.username,
      email: form.value.email,
      password: form.value.password,
      realm_name: form.value.realmName,
    });
    
    // Redirect to sign in page after successful registration
    router.push('/signin');
  } catch (err) {
    error.value = err instanceof Error ? err.message : 'Registration failed';
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
  border-color: #198754;
  box-shadow: 0 0 0 0.2rem rgba(25, 135, 84, 0.25);
}
</style>
