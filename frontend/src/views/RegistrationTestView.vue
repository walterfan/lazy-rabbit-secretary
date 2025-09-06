<template>
  <div class="container mt-4">
    <div class="row">
      <div class="col-12">
        <h2>Registration Test Page</h2>
        <p class="text-muted">This page demonstrates the new registration approval workflow.</p>

        <div class="row">
          <div class="col-md-6">
            <div class="card">
              <div class="card-header">
                <h5>Test Registration</h5>
              </div>
              <div class="card-body">
                <form @submit.prevent="testRegister">
                  <div class="mb-3">
                    <label class="form-label">Username</label>
                    <input 
                      v-model="testUser.username" 
                      type="text" 
                      class="form-control"
                      required
                    >
                  </div>
                  <div class="mb-3">
                    <label class="form-label">Email</label>
                    <input 
                      v-model="testUser.email" 
                      type="email" 
                      class="form-control"
                      required
                    >
                  </div>
                  <div class="mb-3">
                    <label class="form-label">Password</label>
                    <input 
                      v-model="testUser.password" 
                      type="password" 
                      class="form-control"
                      required
                      minlength="8"
                    >
                  </div>
                  <button type="submit" class="btn btn-primary" :disabled="registering">
                    {{ registering ? 'Registering...' : 'Register Test User' }}
                  </button>
                </form>
              </div>
            </div>
          </div>

          <div class="col-md-6">
            <div class="card">
              <div class="card-header">
                <h5>Registration Status</h5>
              </div>
              <div class="card-body">
                <div v-if="registrationResult" class="alert" :class="{
                  'alert-success': registrationResult.success,
                  'alert-danger': !registrationResult.success
                }">
                  <strong>{{ registrationResult.success ? 'Success!' : 'Error!' }}</strong>
                  <p class="mb-0">{{ registrationResult.message }}</p>
                  <div v-if="registrationResult.user" class="mt-2">
                    <small>
                      <strong>User ID:</strong> {{ registrationResult.user.id }}<br>
                      <strong>Status:</strong> {{ registrationResult.user.status }}<br>
                      <strong>Active:</strong> {{ registrationResult.user.is_active ? 'Yes' : 'No' }}
                    </small>
                  </div>
                </div>
                <div v-else class="text-muted">
                  No registration attempts yet.
                </div>
              </div>
            </div>
          </div>
        </div>

        <div class="row mt-4">
          <div class="col-12">
            <div class="card">
              <div class="card-header">
                <h5>Workflow Steps</h5>
              </div>
              <div class="card-body">
                <ol>
                  <li><strong>Register a new user</strong> - Use the form above to create a test user</li>
                  <li><strong>Check pending status</strong> - The user will be created with status "pending" and is_active = false</li>
                  <li><strong>Login as admin</strong> - Use admin/admin123 credentials</li>
                  <li><strong>Go to Registration Management</strong> - Navigate to Admin > Registration Management</li>
                  <li><strong>Approve or deny</strong> - Process the pending registration</li>
                  <li><strong>Test login</strong> - Try to login with the approved user</li>
                </ol>
              </div>
            </div>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref } from 'vue'
import { useAuthStore } from '@/stores/authStore'

const authStore = useAuthStore()

const registering = ref(false)
const testUser = ref({
  username: `testuser_${Date.now()}`,
  email: `test_${Date.now()}@example.com`,
  password: 'password123'
})

const registrationResult = ref<{
  success: boolean
  message: string
  user?: any
} | null>(null)

const testRegister = async () => {
  registering.value = true
  registrationResult.value = null
  
  try {
    const result = await authStore.signUp({
      username: testUser.value.username,
      email: testUser.value.email,
      password: testUser.value.password,
      realm_name: 'default'
    })
    
    registrationResult.value = {
      success: true,
      message: 'User registered successfully! Status is pending approval.',
      user: result.user
    }
    
    // Generate new test user for next registration
    testUser.value = {
      username: `testuser_${Date.now()}`,
      email: `test_${Date.now()}@example.com`,
      password: 'password123'
    }
  } catch (error) {
    registrationResult.value = {
      success: false,
      message: error instanceof Error ? error.message : 'Registration failed'
    }
  } finally {
    registering.value = false
  }
}
</script>

<style scoped>
.card {
  height: 100%;
}

.alert {
  margin-bottom: 0;
}
</style>

