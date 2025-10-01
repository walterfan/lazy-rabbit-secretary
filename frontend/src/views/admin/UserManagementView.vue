<template>
  <div class="container mt-4">
    <div class="row">
      <div class="col-12">
        <div class="d-flex justify-content-between align-items-center mb-4">
          <h2>
            <i class="bi bi-people"></i> {{ $t('admin.userManagement') }}
          </h2>
          <button class="btn btn-primary" @click="showCreateModal = true">
            <i class="bi bi-plus"></i> {{ $t('admin.addUser') }}
          </button>
        </div>

        <!-- Loading State -->
        <div v-if="userStore.loading" class="text-center py-4">
          <div class="spinner-border" role="status">
            <span class="visually-hidden">{{ $t('common.loading') }}</span>
          </div>
        </div>

        <!-- Error State -->
        <div v-else-if="userStore.error" class="alert alert-danger" role="alert">
          <i class="bi bi-exclamation-triangle"></i> {{ userStore.error }}
        </div>

        <!-- User List -->
        <div v-else class="card">
          <div class="card-header">
            <div class="row align-items-center">
              <div class="col">
                <h5 class="mb-0">{{ $t('admin.users') }} ({{ userStore.totalUsers }})</h5>
              </div>
              <div class="col-auto">
                <input
                  type="text"
                  class="form-control"
                  :placeholder="$t('admin.searchUsers')"
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
                    <th>{{ $t('admin.username') }}</th>
                    <th>{{ $t('admin.email') }}</th>
                    <th>{{ $t('admin.realmId') }}</th>
                    <th>{{ $t('admin.status') }}</th>
                    <th>{{ $t('admin.emailConfirmed') }}</th>
                    <th>{{ $t('admin.created') }}</th>
                    <th>{{ $t('admin.actions') }}</th>
                  </tr>
                </thead>
                <tbody>
                  <tr v-for="user in filteredUsers" :key="user.id">
                    <td>{{ user.username }}</td>
                    <td>{{ user.email }}</td>
                    <td>
                      <span class="badge bg-secondary">{{ user.realm_id }}</span>
                    </td>
                    <td>
                      <span :class="getStatusBadgeClass(user.status)">
                        {{ $t(`admin.userStatus.${user.status}`) }}
                      </span>
                    </td>
                    <td>
                      <span :class="user.email_confirmed_at ? 'badge bg-success' : 'badge bg-warning'">
                        {{ user.email_confirmed_at ? $t('admin.confirmed') : $t('admin.pending') }}
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

            <!-- Pagination -->
            <nav v-if="userStore.totalPages > 1" aria-label="User pagination">
              <ul class="pagination justify-content-center">
                <li class="page-item" :class="{ disabled: userStore.currentPage === 1 }">
                  <button class="page-link" @click="changePage(userStore.currentPage - 1)" :disabled="userStore.currentPage === 1">
                    {{ $t('common.previous') }}
                  </button>
                </li>
                <li 
                  v-for="page in visiblePages" 
                  :key="page" 
                  class="page-item" 
                  :class="{ active: page === userStore.currentPage }"
                >
                  <button class="page-link" @click="changePage(page)">{{ page }}</button>
                </li>
                <li class="page-item" :class="{ disabled: userStore.currentPage === userStore.totalPages }">
                  <button class="page-link" @click="changePage(userStore.currentPage + 1)" :disabled="userStore.currentPage === userStore.totalPages">
                    {{ $t('common.next') }}
                  </button>
                </li>
              </ul>
            </nav>
          </div>
        </div>
      </div>
    </div>

    <!-- Create/Edit User Modal -->
    <div class="modal fade" :class="{ show: showCreateModal }" :style="{ display: showCreateModal ? 'block' : 'none' }" tabindex="-1">
      <div class="modal-dialog">
        <div class="modal-content">
          <div class="modal-header">
            <h5 class="modal-title">{{ editingUser ? $t('admin.editUser') : $t('admin.createUser') }}</h5>
            <button type="button" class="btn-close" @click="closeModal"></button>
          </div>
          <div class="modal-body">
            <form @submit.prevent="saveUser">
              <div class="mb-3">
                <label for="username" class="form-label">{{ $t('admin.username') }}</label>
                <input
                  type="text"
                  class="form-control"
                  id="username"
                  v-model="userForm.username"
                  required
                />
              </div>
              <div class="mb-3">
                <label for="email" class="form-label">{{ $t('admin.email') }}</label>
                <input
                  type="email"
                  class="form-control"
                  id="email"
                  v-model="userForm.email"
                  required
                />
              </div>
              <div class="mb-3">
                <label for="realmName" class="form-label">{{ $t('admin.realmName') }}</label>
                <input
                  type="text"
                  class="form-control"
                  id="realmName"
                  v-model="userForm.realm_name"
                  required
                />
              </div>
              
              <!-- Status field for editing users -->
              <div class="mb-3" v-if="editingUser">
                <label for="status" class="form-label">{{ $t('admin.status') }}</label>
                <select 
                  class="form-select" 
                  id="status" 
                  v-model="userForm.status"
                >
                  <option value="pending">{{ $t('admin.userStatus.pending') }}</option>
                  <option value="confirmed">{{ $t('admin.userStatus.confirmed') }}</option>
                  <option value="approved">{{ $t('admin.userStatus.approved') }}</option>
                  <option value="denied">{{ $t('admin.userStatus.denied') }}</option>
                  <option value="suspended">{{ $t('admin.userStatus.suspended') }}</option>
                </select>
              </div>

              <!-- Password fields -->
              <div v-if="!editingUser">
                <div class="mb-3">
                  <label for="password" class="form-label">{{ $t('admin.password') }}</label>
                  <input
                    type="password"
                    class="form-control"
                    id="password"
                    v-model="userForm.password"
                    required
                    minlength="8"
                  />
                  <div class="form-text">{{ $t('admin.passwordMinLength') }}</div>
                </div>
              </div>
              
              <!-- Password change for editing users -->
              <div v-else>
                <div class="mb-3">
                  <div class="form-check">
                    <input
                      class="form-check-input"
                      type="checkbox"
                      id="changePassword"
                      v-model="userForm.changePassword"
                    />
                    <label class="form-check-label" for="changePassword">
                      {{ $t('admin.changePassword') }}
                    </label>
                  </div>
                </div>
                <div class="mb-3" v-if="userForm.changePassword">
                  <label for="newPassword" class="form-label">{{ $t('admin.newPassword') }}</label>
                  <input
                    type="password"
                    class="form-control"
                    id="newPassword"
                    v-model="userForm.password"
                    :required="userForm.changePassword"
                    minlength="8"
                  />
                  <div class="form-text">{{ $t('admin.passwordMinLength') }}</div>
                </div>
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
                    {{ $t('admin.active') }}
                  </label>
                </div>
              </div>
            </form>
          </div>
          <div class="modal-footer">
            <button type="button" class="btn btn-secondary" @click="closeModal">{{ $t('common.cancel') }}</button>
            <button type="button" class="btn btn-primary" @click="saveUser" :disabled="userStore.loading">
              {{ userStore.loading ? $t('common.saving') : $t('common.save') }}
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
import { useI18n } from 'vue-i18n';
import { useUserStore, type User, type CreateUserRequest, type UpdateUserRequest } from '@/stores/userStore';

const { t } = useI18n();
const userStore = useUserStore();

// Local state
const searchTerm = ref('');
const showCreateModal = ref(false);
const editingUser = ref<User | null>(null);

const userForm = ref<CreateUserRequest & { status?: string; changePassword?: boolean }>({
  username: '',
  email: '',
  realm_name: 'default',
  password: '',
  is_active: true,
  status: 'pending',
  changePassword: false,
});

// Computed
const filteredUsers = computed(() => {
  return userStore.searchUsers(searchTerm.value);
});

const visiblePages = computed(() => {
  const pages = [];
  const current = userStore.currentPage;
  const total = userStore.totalPages;
  
  // Show up to 5 pages around current page
  const start = Math.max(1, current - 2);
  const end = Math.min(total, current + 2);
  
  for (let i = start; i <= end; i++) {
    pages.push(i);
  }
  
  return pages;
});

// Methods
const formatDate = (dateString: string) => {
  return format(new Date(dateString), 'MMM dd, yyyy HH:mm');
};

const getStatusBadgeClass = (status: string) => {
  switch (status) {
    case 'active':
      return 'badge bg-success';
    case 'pending':
      return 'badge bg-warning';
    case 'inactive':
      return 'badge bg-danger';
    default:
      return 'badge bg-secondary';
  }
};

const loadUsers = async (page = 1) => {
  try {
    await userStore.getUsers({
      page,
      page_size: userStore.pagination.page_size,
    });
  } catch (err) {
    console.error('Failed to load users:', err);
  }
};

const changePage = (page: number) => {
  if (page >= 1 && page <= userStore.totalPages) {
    loadUsers(page);
  }
};

const editUser = (user: User) => {
  editingUser.value = user;
  userForm.value = {
    username: user.username,
    email: user.email,
    realm_name: 'default', // You might want to get the realm name from the user
    password: '',
    is_active: user.is_active,
    status: user.status,
    changePassword: false,
  };
  showCreateModal.value = true;
};

const saveUser = async () => {
  try {
    if (editingUser.value) {
      const updateData: UpdateUserRequest = {
        username: userForm.value.username,
        email: userForm.value.email,
        is_active: userForm.value.is_active,
        status: userForm.value.status,
      };
      
      // Only include password if user wants to change it
      if (userForm.value.changePassword && userForm.value.password) {
        updateData.password = userForm.value.password;
      }
      
      await userStore.updateUser(editingUser.value.id, updateData);
    } else {
      await userStore.createUser(userForm.value);
    }
    
    await loadUsers(userStore.currentPage);
    closeModal();
  } catch (err) {
    console.error('Failed to save user:', err);
  }
};

const deleteUser = async (userId: string) => {
  if (!confirm(t('admin.confirmDeleteUser'))) return;
  
  try {
    await userStore.deleteUser(userId);
    await loadUsers(userStore.currentPage);
  } catch (err) {
    console.error('Failed to delete user:', err);
  }
};

const closeModal = () => {
  showCreateModal.value = false;
  editingUser.value = null;
  userForm.value = {
    username: '',
    email: '',
    realm_name: 'default',
    password: '',
    is_active: true,
    status: 'pending',
    changePassword: false,
  };
};

// Lifecycle
onMounted(() => {
  loadUsers();
});
</script>

<style scoped>
.modal {
  background-color: rgba(0, 0, 0, 0.5);
}
</style>
