<template>
  <div class="container mt-4">
    <div class="row">
      <div class="col-12">
        <div class="d-flex justify-content-between align-items-center mb-4">
          <h2>
            <i class="bi bi-people"></i> User Management
          </h2>
          <button class="btn btn-primary" @click="showCreateModal = true">
            <i class="bi bi-plus"></i> Add User
          </button>
        </div>

        <!-- User List -->
        <div class="card">
          <div class="card-header">
            <div class="row align-items-center">
              <div class="col">
                <h5 class="mb-0">Users</h5>
              </div>
              <div class="col-auto">
                <input
                  type="text"
                  class="form-control"
                  placeholder="Search users..."
                  v-model="searchTerm"
                />
              </div>
            </div>
          </div>
          <div class="card-body">
            <div class="table-responsive">
              <table class="table table-hover">
                <thead>
                  <tr>
                    <th>Username</th>
                    <th>Email</th>
                    <th>Realm ID</th>
                    <th>Status</th>
                    <th>Created</th>
                    <th>Actions</th>
                  </tr>
                </thead>
                <tbody>
                  <tr v-for="user in filteredUsers" :key="user.id">
                    <td>{{ user.username }}</td>
                    <td>{{ user.email }}</td>
                    <td>{{ user.realm_id }}</td>
                    <td>
                      <span :class="user.is_active ? 'badge bg-success' : 'badge bg-danger'">
                        {{ user.is_active ? 'Active' : 'Inactive' }}
                      </span>
                    </td>
                    <td>{{ formatDate(user.created_at) }}</td>
                    <td>
                      <button class="btn btn-sm btn-outline-primary me-1" @click="editUser(user)">
                        <i class="bi bi-pencil"></i>
                      </button>
                      <button class="btn btn-sm btn-outline-danger" @click="deleteUser(user.id)">
                        <i class="bi bi-trash"></i>
                      </button>
                    </td>
                  </tr>
                </tbody>
              </table>
            </div>
          </div>
        </div>
      </div>
    </div>

    <!-- Create/Edit User Modal -->
    <div class="modal fade" :class="{ show: showCreateModal }" :style="{ display: showCreateModal ? 'block' : 'none' }" tabindex="-1">
      <div class="modal-dialog">
        <div class="modal-content">
          <div class="modal-header">
            <h5 class="modal-title">{{ editingUser ? 'Edit User' : 'Create User' }}</h5>
            <button type="button" class="btn-close" @click="closeModal"></button>
          </div>
          <div class="modal-body">
            <form @submit.prevent="saveUser">
              <div class="mb-3">
                <label for="username" class="form-label">Username</label>
                <input
                  type="text"
                  class="form-control"
                  id="username"
                  v-model="userForm.username"
                  required
                />
              </div>
              <div class="mb-3">
                <label for="email" class="form-label">Email</label>
                <input
                  type="email"
                  class="form-control"
                  id="email"
                  v-model="userForm.email"
                  required
                />
              </div>
              <div class="mb-3">
                <label for="realmId" class="form-label">Realm ID</label>
                <input
                  type="text"
                  class="form-control"
                  id="realmId"
                  v-model="userForm.realm_id"
                  required
                />
              </div>
              <div class="mb-3" v-if="!editingUser">
                <label for="password" class="form-label">Password</label>
                <input
                  type="password"
                  class="form-control"
                  id="password"
                  v-model="userForm.password"
                  :required="!editingUser"
                />
              </div>
              <div class="mb-3">
                <div class="form-check">
                  <input
                    class="form-check-input"
                    type="checkbox"
                    id="isActive"
                    v-model="userForm.is_active"
                  />
                  <label class="form-check-label" for="isActive">
                    Active
                  </label>
                </div>
              </div>
            </form>
          </div>
          <div class="modal-footer">
            <button type="button" class="btn btn-secondary" @click="closeModal">Cancel</button>
            <button type="button" class="btn btn-primary" @click="saveUser" :disabled="loading">
              {{ loading ? 'Saving...' : 'Save' }}
            </button>
          </div>
        </div>
      </div>
    </div>
    <div v-if="showCreateModal" class="modal-backdrop fade show"></div>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted } from 'vue';
import { format } from 'date-fns';

interface User {
  id: string;
  username: string;
  email: string;
  realm_id: string;
  is_active: boolean;
  created_at: string;
  updated_at: string;
}

const users = ref<User[]>([]);
const searchTerm = ref('');
const showCreateModal = ref(false);
const editingUser = ref<User | null>(null);
const loading = ref(false);

const userForm = ref({
  username: '',
  email: '',
  realm_id: 'default',
  password: '',
  is_active: true,
});

const filteredUsers = computed(() => {
  if (!searchTerm.value) return users.value;
  return users.value.filter(user =>
    user.username.toLowerCase().includes(searchTerm.value.toLowerCase()) ||
    user.email.toLowerCase().includes(searchTerm.value.toLowerCase())
  );
});

const formatDate = (dateString: string) => {
  return format(new Date(dateString), 'MMM dd, yyyy');
};

const loadUsers = async () => {
  try {
    const response = await fetch('/api/v1/admin/users', {
      headers: {
        'Authorization': `Bearer ${localStorage.getItem('access_token')}`,
      },
    });
    if (response.ok) {
      users.value = await response.json();
    }
  } catch (error) {
    console.error('Failed to load users:', error);
  }
};

const editUser = (user: User) => {
  editingUser.value = user;
  userForm.value = {
    username: user.username,
    email: user.email,
    realm_id: user.realm_id,
    password: '',
    is_active: user.is_active,
  };
  showCreateModal.value = true;
};

const saveUser = async () => {
  loading.value = true;
  try {
    const url = editingUser.value 
      ? `/api/v1/admin/users/${editingUser.value.id}`
      : '/api/v1/admin/users';
    
    const method = editingUser.value ? 'PUT' : 'POST';
    
    const response = await fetch(url, {
      method,
      headers: {
        'Content-Type': 'application/json',
        'Authorization': `Bearer ${localStorage.getItem('access_token')}`,
      },
      body: JSON.stringify(userForm.value),
    });

    if (response.ok) {
      await loadUsers();
      closeModal();
    }
  } catch (error) {
    console.error('Failed to save user:', error);
  } finally {
    loading.value = false;
  }
};

const deleteUser = async (userId: string) => {
  if (!confirm('Are you sure you want to delete this user?')) return;
  
  try {
    const response = await fetch(`/api/v1/admin/users/${userId}`, {
      method: 'DELETE',
      headers: {
        'Authorization': `Bearer ${localStorage.getItem('access_token')}`,
      },
    });
    
    if (response.ok) {
      await loadUsers();
    }
  } catch (error) {
    console.error('Failed to delete user:', error);
  }
};

const closeModal = () => {
  showCreateModal.value = false;
  editingUser.value = null;
  userForm.value = {
    username: '',
    email: '',
    realm_id: 'default',
    password: '',
    is_active: true,
  };
};

onMounted(() => {
  loadUsers();
});
</script>

<style scoped>
.modal {
  background-color: rgba(0, 0, 0, 0.5);
}
</style>
