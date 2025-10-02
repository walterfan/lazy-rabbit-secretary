<template>
  <div class="container mt-4">
    <div class="d-flex justify-content-between align-items-center mb-4">
      <h2 class="mb-0">Permission Management</h2>
      <div class="badge bg-primary">v1.0</div>
    </div>
    
    <hr class="mb-4" />
    
    <!-- Permission Summary -->
    <div class="row mb-4">
      <div class="col-12">
        <div class="card">
          <div class="card-header bg-info text-white">
            <h5 class="mb-0"><i class="fas fa-user-shield"></i> Current User Permissions</h5>
          </div>
          <div class="card-body">
            <div v-if="loading" class="text-center">
              <div class="spinner-border" role="status">
                <span class="visually-hidden">Loading...</span>
              </div>
            </div>
            <div v-else-if="permissionSummary" class="row">
              <div class="col-md-6">
                <h6>User Information</h6>
                <p><strong>User ID:</strong> {{ permissionSummary.user_id }}</p>
                <p><strong>Realm ID:</strong> {{ permissionSummary.realm_id }}</p>
                <p><strong>Roles:</strong> {{ permissionSummary.roles.join(', ') }}</p>
              </div>
              <div class="col-md-6">
                <h6>Resource Access</h6>
                <div class="table-responsive">
                  <table class="table table-sm">
                    <thead>
                      <tr>
                        <th>Resource</th>
                        <th>Level</th>
                        <th>Source</th>
                      </tr>
                    </thead>
                    <tbody>
                      <tr v-for="(access, resource) in permissionSummary.resources" :key="resource">
                        <td>{{ resource }}</td>
                        <td>
                          <span :class="getLevelBadgeClass(access.level)" class="badge">
                            {{ access.level }}
                          </span>
                        </td>
                        <td>{{ access.source }}</td>
                      </tr>
                    </tbody>
                  </table>
                </div>
              </div>
            </div>
            <div v-else class="text-muted">No permission data available</div>
          </div>
        </div>
      </div>
    </div>

    <!-- Permission Management Tabs -->
    <ul class="nav nav-tabs" id="permissionTabs" role="tablist">
      <li class="nav-item" role="presentation">
        <button class="nav-link active" id="user-permissions-tab" data-bs-toggle="tab" data-bs-target="#user-permissions" type="button" role="tab">
          <i class="fas fa-user"></i> User Permissions
        </button>
      </li>
      <li class="nav-item" role="presentation">
        <button class="nav-link" id="role-permissions-tab" data-bs-toggle="tab" data-bs-target="#role-permissions" type="button" role="tab">
          <i class="fas fa-users"></i> Role Permissions
        </button>
      </li>
      <li class="nav-item" role="presentation">
        <button class="nav-link" id="permission-check-tab" data-bs-toggle="tab" data-bs-target="#permission-check" type="button" role="tab">
          <i class="fas fa-search"></i> Permission Check
        </button>
      </li>
    </ul>

    <div class="tab-content" id="permissionTabsContent">
      <!-- User Permissions Tab -->
      <div class="tab-pane fade show active" id="user-permissions" role="tabpanel">
        <div class="card mt-3">
          <div class="card-header d-flex justify-content-between align-items-center">
            <h5 class="mb-0">User Permissions</h5>
            <button class="btn btn-primary btn-sm" @click="showCreateUserPermissionModal = true">
              <i class="fas fa-plus"></i> Add Permission
            </button>
          </div>
          <div class="card-body">
            <div v-if="loadingUserPermissions" class="text-center">
              <div class="spinner-border" role="status">
                <span class="visually-hidden">Loading...</span>
              </div>
            </div>
            <div v-else class="table-responsive">
              <table class="table table-striped">
                <thead>
                  <tr>
                    <th>User ID</th>
                    <th>Resource</th>
                    <th>Actions</th>
                    <th>Level</th>
                    <th>Expires</th>
                    <th>Actions</th>
                  </tr>
                </thead>
                <tbody>
                  <tr v-for="permission in userPermissions" :key="permission.id">
                    <td>{{ permission.user_id }}</td>
                    <td>{{ permission.resource }}</td>
                    <td>
                      <span v-for="action in getPermissionActions(permission)" :key="action" class="badge bg-secondary me-1">
                        {{ action }}
                      </span>
                    </td>
                    <td>
                      <span :class="getLevelBadgeClass(permission.level)" class="badge">
                        {{ permission.level }}
                      </span>
                    </td>
                    <td>{{ formatDate(permission.expires_at) }}</td>
                    <td>
                      <button class="btn btn-sm btn-outline-primary me-1" @click="editUserPermission(permission)">
                        <i class="fas fa-edit"></i>
                      </button>
                      <button class="btn btn-sm btn-outline-danger" @click="deleteUserPermission(permission.id)">
                        <i class="fas fa-trash"></i>
                      </button>
                    </td>
                  </tr>
                </tbody>
              </table>
            </div>
          </div>
        </div>
      </div>

      <!-- Role Permissions Tab -->
      <div class="tab-pane fade" id="role-permissions" role="tabpanel">
        <div class="card mt-3">
          <div class="card-header d-flex justify-content-between align-items-center">
            <h5 class="mb-0">Role Permissions</h5>
            <button class="btn btn-primary btn-sm" @click="showCreateRolePermissionModal = true">
              <i class="fas fa-plus"></i> Add Permission
            </button>
          </div>
          <div class="card-body">
            <div v-if="loadingRolePermissions" class="text-center">
              <div class="spinner-border" role="status">
                <span class="visually-hidden">Loading...</span>
              </div>
            </div>
            <div v-else class="table-responsive">
              <table class="table table-striped">
                <thead>
                  <tr>
                    <th>Role ID</th>
                    <th>Resource</th>
                    <th>Actions</th>
                    <th>Level</th>
                    <th>Expires</th>
                    <th>Actions</th>
                  </tr>
                </thead>
                <tbody>
                  <tr v-for="permission in rolePermissions" :key="permission.id">
                    <td>{{ permission.role_id }}</td>
                    <td>{{ permission.resource }}</td>
                    <td>
                      <span v-for="action in getPermissionActions(permission)" :key="action" class="badge bg-secondary me-1">
                        {{ action }}
                      </span>
                    </td>
                    <td>
                      <span :class="getLevelBadgeClass(permission.level)" class="badge">
                        {{ permission.level }}
                      </span>
                    </td>
                    <td>{{ formatDate(permission.expires_at) }}</td>
                    <td>
                      <button class="btn btn-sm btn-outline-primary me-1" @click="editRolePermission(permission)">
                        <i class="fas fa-edit"></i>
                      </button>
                      <button class="btn btn-sm btn-outline-danger" @click="deleteRolePermission(permission.id)">
                        <i class="fas fa-trash"></i>
                      </button>
                    </td>
                  </tr>
                </tbody>
              </table>
            </div>
          </div>
        </div>
      </div>

      <!-- Permission Check Tab -->
      <div class="tab-pane fade" id="permission-check" role="tabpanel">
        <div class="card mt-3">
          <div class="card-header">
            <h5 class="mb-0">Permission Check</h5>
          </div>
          <div class="card-body">
            <form @submit.prevent="checkPermission">
              <div class="row">
                <div class="col-md-6">
                  <div class="mb-3">
                    <label for="checkAction" class="form-label">Action</label>
                    <select id="checkAction" v-model="permissionCheck.action" class="form-select" required>
                      <option value="">Select an action...</option>
                      <option v-for="action in availableActions" :key="action" :value="action">
                        {{ action }}
                      </option>
                    </select>
                  </div>
                </div>
                <div class="col-md-6">
                  <div class="mb-3">
                    <label for="checkResource" class="form-label">Resource</label>
                    <select id="checkResource" v-model="permissionCheck.resource" class="form-select" required>
                      <option value="">Select a resource...</option>
                      <option v-for="resource in availableResources" :key="resource" :value="resource">
                        {{ resource }}
                      </option>
                    </select>
                  </div>
                </div>
              </div>
              <button type="submit" class="btn btn-primary" :disabled="checkingPermission">
                <span v-if="checkingPermission" class="spinner-border spinner-border-sm me-2"></span>
                Check Permission
              </button>
            </form>

            <div v-if="permissionCheckResult" class="mt-4">
              <div class="alert" :class="permissionCheckResult.allowed ? 'alert-success' : 'alert-danger'">
                <h6>Permission Check Result</h6>
                <p><strong>Allowed:</strong> {{ permissionCheckResult.allowed ? 'Yes' : 'No' }}</p>
                <p v-if="permissionCheckResult.level"><strong>Level:</strong> {{ permissionCheckResult.level }}</p>
                <p v-if="permissionCheckResult.source"><strong>Source:</strong> {{ permissionCheckResult.source }}</p>
                <p v-if="permissionCheckResult.reason"><strong>Reason:</strong> {{ permissionCheckResult.reason }}</p>
              </div>
            </div>
          </div>
        </div>
      </div>
    </div>

    <!-- Create User Permission Modal -->
    <div class="modal fade" id="createUserPermissionModal" tabindex="-1" v-if="showCreateUserPermissionModal">
      <div class="modal-dialog">
        <div class="modal-content">
          <div class="modal-header">
            <h5 class="modal-title">Create User Permission</h5>
            <button type="button" class="btn-close" @click="showCreateUserPermissionModal = false"></button>
          </div>
          <div class="modal-body">
            <form @submit.prevent="createUserPermission">
              <div class="mb-3">
                <label for="userPermissionUserId" class="form-label">User ID</label>
                <input type="text" id="userPermissionUserId" v-model="newUserPermission.user_id" class="form-control" required>
              </div>
              <div class="mb-3">
                <label for="userPermissionResource" class="form-label">Resource</label>
                <select id="userPermissionResource" v-model="newUserPermission.resource" class="form-select" required>
                  <option value="">Select a resource...</option>
                  <option v-for="resource in availableResources" :key="resource" :value="resource">
                    {{ resource }}
                  </option>
                </select>
              </div>
              <div class="mb-3">
                <label for="userPermissionLevel" class="form-label">Level</label>
                <select id="userPermissionLevel" v-model="newUserPermission.level" class="form-select" required>
                  <option value="">Select a level...</option>
                  <option v-for="level in availableLevels" :key="level" :value="level">
                    {{ level }}
                  </option>
                </select>
              </div>
              <div class="mb-3">
                <label class="form-label">Actions</label>
                <div class="row">
                  <div class="col-md-6" v-for="action in availableActions" :key="action">
                    <div class="form-check">
                      <input type="checkbox" :id="`userAction_${action}`" :value="action" v-model="newUserPermission.actions" class="form-check-input">
                      <label :for="`userAction_${action}`" class="form-check-label">{{ action }}</label>
                    </div>
                  </div>
                </div>
              </div>
            </form>
          </div>
          <div class="modal-footer">
            <button type="button" class="btn btn-secondary" @click="showCreateUserPermissionModal = false">Cancel</button>
            <button type="button" class="btn btn-primary" @click="createUserPermission" :disabled="creatingUserPermission">
              <span v-if="creatingUserPermission" class="spinner-border spinner-border-sm me-2"></span>
              Create Permission
            </button>
          </div>
        </div>
      </div>
    </div>

    <!-- Create Role Permission Modal -->
    <div class="modal fade" id="createRolePermissionModal" tabindex="-1" v-if="showCreateRolePermissionModal">
      <div class="modal-dialog">
        <div class="modal-content">
          <div class="modal-header">
            <h5 class="modal-title">Create Role Permission</h5>
            <button type="button" class="btn-close" @click="showCreateRolePermissionModal = false"></button>
          </div>
          <div class="modal-body">
            <form @submit.prevent="createRolePermission">
              <div class="mb-3">
                <label for="rolePermissionRoleId" class="form-label">Role ID</label>
                <input type="text" id="rolePermissionRoleId" v-model="newRolePermission.role_id" class="form-control" required>
              </div>
              <div class="mb-3">
                <label for="rolePermissionResource" class="form-label">Resource</label>
                <select id="rolePermissionResource" v-model="newRolePermission.resource" class="form-select" required>
                  <option value="">Select a resource...</option>
                  <option v-for="resource in availableResources" :key="resource" :value="resource">
                    {{ resource }}
                  </option>
                </select>
              </div>
              <div class="mb-3">
                <label for="rolePermissionLevel" class="form-label">Level</label>
                <select id="rolePermissionLevel" v-model="newRolePermission.level" class="form-select" required>
                  <option value="">Select a level...</option>
                  <option v-for="level in availableLevels" :key="level" :value="level">
                    {{ level }}
                  </option>
                </select>
              </div>
              <div class="mb-3">
                <label class="form-label">Actions</label>
                <div class="row">
                  <div class="col-md-6" v-for="action in availableActions" :key="action">
                    <div class="form-check">
                      <input type="checkbox" :id="`roleAction_${action}`" :value="action" v-model="newRolePermission.actions" class="form-check-input">
                      <label :for="`roleAction_${action}`" class="form-check-label">{{ action }}</label>
                    </div>
                  </div>
                </div>
              </div>
            </form>
          </div>
          <div class="modal-footer">
            <button type="button" class="btn btn-secondary" @click="showCreateRolePermissionModal = false">Cancel</button>
            <button type="button" class="btn btn-primary" @click="createRolePermission" :disabled="creatingRolePermission">
              <span v-if="creatingRolePermission" class="spinner-border spinner-border-sm me-2"></span>
              Create Permission
            </button>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue';
import { useAuthStore } from '@/stores/authStore';

const authStore = useAuthStore();

// Helper function to get headers with authentication
const getHeaders = () => {
  const headers: Record<string, string> = {
    'Content-Type': 'application/json',
  };
  
  if (authStore.token) {
    headers.Authorization = `Bearer ${authStore.token}`;
  }
  
  return headers;
};

// State
const loading = ref(false);
const loadingUserPermissions = ref(false);
const loadingRolePermissions = ref(false);
const checkingPermission = ref(false);
const creatingUserPermission = ref(false);
const creatingRolePermission = ref(false);

const permissionSummary = ref<any>(null);
const userPermissions = ref<any[]>([]);
const rolePermissions = ref<any[]>([]);
const availableActions = ref<string[]>([]);
const availableResources = ref<string[]>([]);
const availableLevels = ref<string[]>([]);

// Modals
const showCreateUserPermissionModal = ref(false);
const showCreateRolePermissionModal = ref(false);

// Forms
const permissionCheck = ref({
  action: '',
  resource: ''
});

const newUserPermission = ref({
  user_id: '',
  resource: '',
  level: '',
  actions: [] as string[]
});

const newRolePermission = ref({
  role_id: '',
  resource: '',
  level: '',
  actions: [] as string[]
});

const permissionCheckResult = ref<any>(null);

// API calls
const fetchPermissionSummary = async () => {
  loading.value = true;
  try {
    const apiBase = import.meta.env.VITE_API_BASE_URL;
    const response = await fetch(`${apiBase}/api/v1/user-permissions/`, {
      headers: getHeaders()
    });
    const data = await response.json();
    permissionSummary.value = data;
  } catch (error) {
    console.error('Failed to fetch permission summary:', error);
  } finally {
    loading.value = false;
  }
};

const fetchAvailableOptions = async () => {
  try {
    const apiBase = import.meta.env.VITE_API_BASE_URL;
    
    const [actionsRes, resourcesRes, levelsRes] = await Promise.all([
      fetch(`${apiBase}/api/v1/permissions/actions`, { headers: getHeaders() }),
      fetch(`${apiBase}/api/v1/permissions/resources`, { headers: getHeaders() }),
      fetch(`${apiBase}/api/v1/permissions/levels`, { headers: getHeaders() })
    ]);

    const [actionsData, resourcesData, levelsData] = await Promise.all([
      actionsRes.json(),
      resourcesRes.json(),
      levelsRes.json()
    ]);

    availableActions.value = actionsData.actions || [];
    availableResources.value = resourcesData.resources || [];
    availableLevels.value = levelsData.levels || [];
  } catch (error) {
    console.error('Failed to fetch available options:', error);
  }
};

const checkPermission = async () => {
  checkingPermission.value = true;
  try {
    const apiBase = import.meta.env.VITE_API_BASE_URL;
    const response = await fetch(`${apiBase}/api/v1/user-permissions/check?action=${permissionCheck.value.action}&resource=${permissionCheck.value.resource}`, {
      headers: getHeaders()
    });
    const data = await response.json();
    permissionCheckResult.value = data;
  } catch (error) {
    console.error('Failed to check permission:', error);
  } finally {
    checkingPermission.value = false;
  }
};

const createUserPermission = async () => {
  creatingUserPermission.value = true;
  try {
    const apiBase = import.meta.env.VITE_API_BASE_URL;
    const response = await fetch(`${apiBase}/api/v1/permissions/users`, {
      method: 'POST',
      headers: getHeaders(),
      body: JSON.stringify(newUserPermission.value)
    });

    if (response.ok) {
      showCreateUserPermissionModal.value = false;
      newUserPermission.value = { user_id: '', resource: '', level: '', actions: [] };
      // Refresh permissions list
    } else {
      const error = await response.json();
      console.error('Failed to create user permission:', error);
    }
  } catch (error) {
    console.error('Failed to create user permission:', error);
  } finally {
    creatingUserPermission.value = false;
  }
};

const createRolePermission = async () => {
  creatingRolePermission.value = true;
  try {
    const apiBase = import.meta.env.VITE_API_BASE_URL;
    const response = await fetch(`${apiBase}/api/v1/permissions/roles`, {
      method: 'POST',
      headers: getHeaders(),
      body: JSON.stringify(newRolePermission.value)
    });

    if (response.ok) {
      showCreateRolePermissionModal.value = false;
      newRolePermission.value = { role_id: '', resource: '', level: '', actions: [] };
      // Refresh permissions list
    } else {
      const error = await response.json();
      console.error('Failed to create role permission:', error);
    }
  } catch (error) {
    console.error('Failed to create role permission:', error);
  } finally {
    creatingRolePermission.value = false;
  }
};

// Helper functions
const getLevelBadgeClass = (level: string) => {
  switch (level) {
    case 'readonly':
      return 'bg-secondary';
    case 'readwrite':
      return 'bg-warning';
    case 'full':
      return 'bg-success';
    case 'custom':
      return 'bg-info';
    default:
      return 'bg-light text-dark';
  }
};

const getPermissionActions = (permission: any) => {
  try {
    return JSON.parse(permission.actions);
  } catch {
    return [];
  }
};

const formatDate = (dateString: string | null) => {
  if (!dateString) return 'Never';
  return new Date(dateString).toLocaleDateString();
};

const editUserPermission = (permission: any) => {
  // TODO: Implement edit functionality
  console.log('Edit user permission:', permission);
};

const deleteUserPermission = async (permissionId: string) => {
  if (confirm('Are you sure you want to delete this permission?')) {
    try {
      const apiBase = import.meta.env.VITE_API_BASE_URL;
      const response = await fetch(`${apiBase}/api/v1/permissions/users/${permissionId}`, {
        method: 'DELETE',
        headers: getHeaders()
      });

      if (response.ok) {
        // Refresh permissions list
      } else {
        const error = await response.json();
        console.error('Failed to delete user permission:', error);
      }
    } catch (error) {
      console.error('Failed to delete user permission:', error);
    }
  }
};

const editRolePermission = (permission: any) => {
  // TODO: Implement edit functionality
  console.log('Edit role permission:', permission);
};

const deleteRolePermission = async (permissionId: string) => {
  if (confirm('Are you sure you want to delete this permission?')) {
    try {
      const apiBase = import.meta.env.VITE_API_BASE_URL;
      const response = await fetch(`${apiBase}/api/v1/permissions/roles/${permissionId}`, {
        method: 'DELETE',
        headers: getHeaders()
      });

      if (response.ok) {
        // Refresh permissions list
      } else {
        const error = await response.json();
        console.error('Failed to delete role permission:', error);
      }
    } catch (error) {
      console.error('Failed to delete role permission:', error);
    }
  }
};

// Initialize
onMounted(() => {
  fetchPermissionSummary();
  fetchAvailableOptions();
});
</script>

<style scoped>
.card {
  border-radius: 10px;
  border: none;
  box-shadow: 0 2px 10px rgba(0,0,0,0.1);
}

.card-header {
  border-radius: 10px 10px 0 0 !important;
  border-bottom: none;
}

.nav-tabs .nav-link {
  border-radius: 8px 8px 0 0;
  border: none;
  color: #6c757d;
}

.nav-tabs .nav-link.active {
  background-color: #0d6efd;
  color: white;
}

.table th {
  border-top: none;
  font-weight: 600;
  color: #495057;
}

.badge {
  font-size: 0.75rem;
}

.modal-content {
  border-radius: 10px;
  border: none;
  box-shadow: 0 10px 30px rgba(0,0,0,0.2);
}

.form-control, .form-select {
  border-radius: 8px;
  border: 1px solid #dee2e6;
}

.form-control:focus, .form-select:focus {
  border-color: #0d6efd;
  box-shadow: 0 0 0 0.2rem rgba(13, 110, 253, 0.25);
}

.btn {
  border-radius: 8px;
}

.alert {
  border-radius: 8px;
  border: none;
}

.spinner-border-sm {
  width: 1rem;
  height: 1rem;
}
</style>
