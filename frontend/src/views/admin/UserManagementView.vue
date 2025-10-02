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
                    <th>{{ $t('admin.roles') }}</th>
                    <th>{{ $t('admin.status') }}</th>
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
                      <div v-if="user.roles && user.roles.length > 0" class="d-flex flex-wrap gap-1">
                        <span 
                          v-for="role in user.roles" 
                          :key="role.id"
                          class="badge bg-primary cursor-pointer"
                          @click="viewRole(role.id)"
                          :title="role.description"
                        >
                          {{ role.name }}
                        </span>
                      </div>
                      <span v-else class="text-muted">{{ $t('admin.noRoles') }}</span>
                    </td>
                    <td>
                      <span :class="getStatusBadgeClass(user.status)">
                        {{ $t(`admin.userStatus.${user.status || 'unknown'}`) }}
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
                <div v-if="userForm.changePassword">
                  <div class="mb-3">
                    <label for="currentPassword" class="form-label">{{ $t('admin.currentPassword') }}</label>
                    <input
                      type="password"
                      class="form-control"
                      id="currentPassword"
                      v-model="userForm.currentPassword"
                      :required="userForm.changePassword"
                    />
                    <div class="form-text">{{ $t('admin.currentPasswordRequired') }}</div>
                  </div>
                  <div class="mb-3">
                    <label for="newPassword" class="form-label">{{ $t('admin.newPassword') }}</label>
                    <input
                      type="password"
                      class="form-control"
                      id="newPassword"
                      v-model="userForm.newPassword"
                      :required="userForm.changePassword"
                      minlength="8"
                    />
                    <div class="form-text">{{ $t('admin.passwordMinLength') }}</div>
                  </div>
                </div>
              </div>

              <!-- Role Management for editing users -->
              <div class="mb-3" v-if="editingUser">
                <label class="form-label">{{ $t('admin.roles') }}</label>
                <div class="border rounded p-3">
                  <div v-if="roleStore.loading" class="text-center">
                    <div class="spinner-border spinner-border-sm" role="status">
                      <span class="visually-hidden">{{ $t('common.loading') }}</span>
                    </div>
                  </div>
                  <div v-else>
                    <div class="row">
                      <div class="col-md-6">
                        <h6>{{ $t('admin.availableRoles') }}</h6>
                        <div class="available-roles" style="max-height: 150px; overflow-y: auto;">
                          <div 
                            v-for="role in availableRoles" 
                            :key="role.id"
                            class="form-check"
                          >
                            <input
                              class="form-check-input"
                              type="checkbox"
                              :id="`role-${role.id}`"
                              :checked="isRoleSelected(role.id)"
                              @change="toggleRole(role.id, $event.target.checked)"
                            />
                            <label class="form-check-label" :for="`role-${role.id}`">
                              {{ role.name }}
                              <small class="text-muted d-block">{{ role.description }}</small>
                            </label>
                          </div>
                        </div>
                      </div>
                      <div class="col-md-6">
                        <h6>{{ $t('admin.selectedRoles') }}</h6>
                        <div class="selected-roles">
                          <div v-if="selectedRoleIds.length === 0" class="text-muted">
                            {{ $t('admin.noRolesSelected') }}
                          </div>
                          <div v-else>
                            <span 
                              v-for="roleId in selectedRoleIds" 
                              :key="roleId"
                              class="badge bg-primary me-1 mb-1"
                            >
                              {{ getRoleName(roleId) }}
                              <button 
                                type="button" 
                                class="btn-close btn-close-white ms-1" 
                                @click="toggleRole(roleId, false)"
                                style="font-size: 0.7em;"
                              ></button>
                            </span>
                          </div>
                        </div>
                      </div>
                    </div>
                  </div>
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
import { useRouter } from 'vue-router';
import { useUserStore, type User, type CreateUserRequest, type UpdateUserRequest } from '@/stores/userStore';
import { useRoleStore, type Role } from '@/stores/roleStore';

const { t } = useI18n();
const router = useRouter();
const userStore = useUserStore();
const roleStore = useRoleStore();

// Local state
const searchTerm = ref('');
const showCreateModal = ref(false);
const editingUser = ref<User | null>(null);
const selectedRoleIds = ref<string[]>([]);

const userForm = ref<CreateUserRequest & { status?: string; changePassword?: boolean; currentPassword?: string; newPassword?: string }>({
  username: '',
  email: '',
  realm_name: 'default',
  password: '',
  status: 'pending',
  changePassword: false,
  currentPassword: '',
  newPassword: '',
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

// Role management computed properties
const availableRoles = computed(() => {
  return roleStore.roles || [];
});

const isRoleSelected = (roleId: string) => {
  return selectedRoleIds.value.includes(roleId);
};

const getRoleName = (roleId: string) => {
  const role = roleStore.roles.find(r => r.id === roleId);
  return role ? role.name : roleId;
};

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

// Role management methods
const toggleRole = (roleId: string, selected: boolean) => {
  if (selected) {
    if (!selectedRoleIds.value.includes(roleId)) {
      selectedRoleIds.value.push(roleId);
    }
  } else {
    selectedRoleIds.value = selectedRoleIds.value.filter(id => id !== roleId);
  }
};

const loadRoles = async () => {
  try {
    await roleStore.getRoles({ page_size: 100 }); // Load all roles
  } catch (err) {
    console.error('Failed to load roles:', err);
  }
};

const editUser = (user: User) => {
  editingUser.value = user;
  userForm.value = {
    username: user.username,
    email: user.email,
    realm_name: 'default', // You might want to get the realm name from the user
    password: '',
    status: user.status,
    changePassword: false,
    currentPassword: '',
    newPassword: '',
  };
  
  // Load user's current roles
  selectedRoleIds.value = user.roles ? user.roles.map(role => role.id) : [];
  
  // Load all available roles
  loadRoles();
  
  showCreateModal.value = true;
};

const saveUser = async () => {
  try {
    if (editingUser.value) {
      const updateData: UpdateUserRequest = {
        username: userForm.value.username,
        email: userForm.value.email,
        status: userForm.value.status,
        role_ids: selectedRoleIds.value, // Include selected role IDs
      };
      
      // Only include password fields if user wants to change password
      if (userForm.value.changePassword && userForm.value.newPassword) {
        updateData.current_password = userForm.value.currentPassword;
        updateData.new_password = userForm.value.newPassword;
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
  selectedRoleIds.value = [];
  userForm.value = {
    username: '',
    email: '',
    realm_name: 'default',
    password: '',
    status: 'pending',
    changePassword: false,
    currentPassword: '',
    newPassword: '',
  };
};

const viewRole = (roleId: string) => {
  // Navigate to role management page with the specific role
  router.push(`/admin/roles?role=${roleId}`);
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

.cursor-pointer {
  cursor: pointer;
}

.cursor-pointer:hover {
  opacity: 0.8;
}
</style>
