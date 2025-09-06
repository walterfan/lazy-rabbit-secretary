<template>
  <div class="container mt-4">
    <div class="row">
      <div class="col-12">
        <div class="d-flex justify-content-between align-items-center mb-4">
          <h2>
            <i class="bi bi-shield"></i> Role Management
          </h2>
          <button class="btn btn-primary" @click="showCreateModal = true">
            <i class="bi bi-plus"></i> Add Role
          </button>
        </div>

        <!-- Role List -->
        <div class="card">
          <div class="card-header">
            <div class="row align-items-center">
              <div class="col">
                <h5 class="mb-0">Roles</h5>
              </div>
              <div class="col-auto">
                <input
                  type="text"
                  class="form-control"
                  placeholder="Search roles..."
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
                    <th>Name</th>
                    <th>Description</th>
                    <th>Realm ID</th>
                    <th>Created</th>
                    <th>Actions</th>
                  </tr>
                </thead>
                <tbody>
                  <tr v-for="role in filteredRoles" :key="role.id">
                    <td>{{ role.name }}</td>
                    <td>{{ role.description }}</td>
                    <td>{{ role.realm_id }}</td>
                    <td>{{ formatDate(role.created_time) }}</td>
                    <td>
                      <button class="btn btn-sm btn-outline-primary me-1" @click="editRole(role)">
                        <i class="bi bi-pencil"></i>
                      </button>
                      <button class="btn btn-sm btn-outline-info me-1" @click="managePolicies(role)">
                        <i class="bi bi-key"></i>
                      </button>
                      <button class="btn btn-sm btn-outline-danger" @click="deleteRole(role.id)">
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

    <!-- Create/Edit Role Modal -->
    <div class="modal fade" :class="{ show: showCreateModal }" :style="{ display: showCreateModal ? 'block' : 'none' }" tabindex="-1">
      <div class="modal-dialog">
        <div class="modal-content">
          <div class="modal-header">
            <h5 class="modal-title">{{ editingRole ? 'Edit Role' : 'Create Role' }}</h5>
            <button type="button" class="btn-close" @click="closeModal"></button>
          </div>
          <div class="modal-body">
            <form @submit.prevent="saveRole">
              <div class="mb-3">
                <label for="name" class="form-label">Role Name</label>
                <input
                  type="text"
                  class="form-control"
                  id="name"
                  v-model="roleForm.name"
                  required
                />
              </div>
              <div class="mb-3">
                <label for="description" class="form-label">Description</label>
                <textarea
                  class="form-control"
                  id="description"
                  v-model="roleForm.description"
                  rows="3"
                ></textarea>
              </div>
              <div class="mb-3">
                <label for="realmId" class="form-label">Realm ID</label>
                <input
                  type="text"
                  class="form-control"
                  id="realmId"
                  v-model="roleForm.realm_id"
                  required
                />
              </div>
            </form>
          </div>
          <div class="modal-footer">
            <button type="button" class="btn btn-secondary" @click="closeModal">Cancel</button>
            <button type="button" class="btn btn-primary" @click="saveRole" :disabled="loading">
              {{ loading ? 'Saving...' : 'Save' }}
            </button>
          </div>
        </div>
      </div>
    </div>

    <!-- Policy Management Modal -->
    <div class="modal fade" :class="{ show: showPolicyModal }" :style="{ display: showPolicyModal ? 'block' : 'none' }" tabindex="-1">
      <div class="modal-dialog modal-lg">
        <div class="modal-content">
          <div class="modal-header">
            <h5 class="modal-title">Manage Policies for {{ selectedRole?.name }}</h5>
            <button type="button" class="btn-close" @click="closePolicyModal"></button>
          </div>
          <div class="modal-body">
            <div class="row">
              <div class="col-md-6">
                <h6>Available Policies</h6>
                <div class="list-group">
                  <div
                    v-for="policy in availablePolicies"
                    :key="policy.id"
                    class="list-group-item list-group-item-action d-flex justify-content-between align-items-center"
                    @click="assignPolicy(policy)"
                  >
                    <div>
                      <strong>{{ policy.name }}</strong>
                      <br>
                      <small class="text-muted">{{ policy.description }}</small>
                    </div>
                    <button class="btn btn-sm btn-outline-primary">
                      <i class="bi bi-plus"></i>
                    </button>
                  </div>
                </div>
              </div>
              <div class="col-md-6">
                <h6>Assigned Policies</h6>
                <div class="list-group">
                  <div
                    v-for="policy in assignedPolicies"
                    :key="policy.id"
                    class="list-group-item d-flex justify-content-between align-items-center"
                  >
                    <div>
                      <strong>{{ policy.name }}</strong>
                      <br>
                      <small class="text-muted">{{ policy.description }}</small>
                    </div>
                    <button class="btn btn-sm btn-outline-danger" @click="removePolicy(policy.id)">
                      <i class="bi bi-dash"></i>
                    </button>
                  </div>
                </div>
              </div>
            </div>
          </div>
          <div class="modal-footer">
            <button type="button" class="btn btn-secondary" @click="closePolicyModal">Close</button>
          </div>
        </div>
      </div>
    </div>

    <div v-if="showCreateModal || showPolicyModal" class="modal-backdrop fade show"></div>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted } from 'vue';
import { format } from 'date-fns';

interface Role {
  id: string;
  name: string;
  description: string;
  realm_id: string;
  created_time: string;
  updated_time: string;
}

interface Policy {
  id: string;
  name: string;
  description: string;
  realm_id: string;
}

const roles = ref<Role[]>([]);
const searchTerm = ref('');
const showCreateModal = ref(false);
const showPolicyModal = ref(false);
const editingRole = ref<Role | null>(null);
const selectedRole = ref<Role | null>(null);
const availablePolicies = ref<Policy[]>([]);
const assignedPolicies = ref<Policy[]>([]);
const loading = ref(false);

const roleForm = ref({
  name: '',
  description: '',
  realm_id: 'default',
});

const filteredRoles = computed(() => {
  if (!searchTerm.value) return roles.value;
  return roles.value.filter(role =>
    role.name.toLowerCase().includes(searchTerm.value.toLowerCase()) ||
    role.description.toLowerCase().includes(searchTerm.value.toLowerCase())
  );
});

const formatDate = (dateString: string) => {
  return format(new Date(dateString), 'MMM dd, yyyy');
};

const loadRoles = async () => {
  try {
    const response = await fetch('/api/v1/admin/roles', {
      headers: {
        'Authorization': `Bearer ${localStorage.getItem('access_token')}`,
      },
    });
    if (response.ok) {
      roles.value = await response.json();
    }
  } catch (error) {
    console.error('Failed to load roles:', error);
  }
};

const editRole = (role: Role) => {
  editingRole.value = role;
  roleForm.value = {
    name: role.name,
    description: role.description,
    realm_id: role.realm_id,
  };
  showCreateModal.value = true;
};

const saveRole = async () => {
  loading.value = true;
  try {
    const url = editingRole.value 
      ? `/api/v1/admin/roles/${editingRole.value.id}`
      : '/api/v1/admin/roles';
    
    const method = editingRole.value ? 'PUT' : 'POST';
    
    const response = await fetch(url, {
      method,
      headers: {
        'Content-Type': 'application/json',
        'Authorization': `Bearer ${localStorage.getItem('access_token')}`,
      },
      body: JSON.stringify(roleForm.value),
    });

    if (response.ok) {
      await loadRoles();
      closeModal();
    }
  } catch (error) {
    console.error('Failed to save role:', error);
  } finally {
    loading.value = false;
  }
};

const deleteRole = async (roleId: string) => {
  if (!confirm('Are you sure you want to delete this role?')) return;
  
  try {
    const response = await fetch(`/api/v1/admin/roles/${roleId}`, {
      method: 'DELETE',
      headers: {
        'Authorization': `Bearer ${localStorage.getItem('access_token')}`,
      },
    });
    
    if (response.ok) {
      await loadRoles();
    }
  } catch (error) {
    console.error('Failed to delete role:', error);
  }
};

const managePolicies = async (role: Role) => {
  selectedRole.value = role;
  showPolicyModal.value = true;
  
  // Load available policies and assigned policies
  try {
    const [availableResponse, assignedResponse] = await Promise.all([
      fetch('/api/v1/admin/policies', {
        headers: {
          'Authorization': `Bearer ${localStorage.getItem('access_token')}`,
        },
      }),
      fetch(`/api/v1/admin/roles/${role.id}/policies`, {
        headers: {
          'Authorization': `Bearer ${localStorage.getItem('access_token')}`,
        },
      }),
    ]);

    if (availableResponse.ok) {
      availablePolicies.value = await availableResponse.json();
    }
    if (assignedResponse.ok) {
      assignedPolicies.value = await assignedResponse.json();
    }
  } catch (error) {
    console.error('Failed to load policies:', error);
  }
};

const assignPolicy = async (policy: Policy) => {
  if (!selectedRole.value) return;
  
  try {
    const response = await fetch(`/api/v1/admin/roles/${selectedRole.value.id}/policies`, {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
        'Authorization': `Bearer ${localStorage.getItem('access_token')}`,
      },
      body: JSON.stringify({ policy_id: policy.id }),
    });

    if (response.ok) {
      assignedPolicies.value.push(policy);
      availablePolicies.value = availablePolicies.value.filter(p => p.id !== policy.id);
    }
  } catch (error) {
    console.error('Failed to assign policy:', error);
  }
};

const removePolicy = async (policyId: string) => {
  if (!selectedRole.value) return;
  
  try {
    const response = await fetch(`/api/v1/admin/roles/${selectedRole.value.id}/policies/${policyId}`, {
      method: 'DELETE',
      headers: {
        'Authorization': `Bearer ${localStorage.getItem('access_token')}`,
      },
    });

    if (response.ok) {
      const policy = assignedPolicies.value.find(p => p.id === policyId);
      if (policy) {
        availablePolicies.value.push(policy);
        assignedPolicies.value = assignedPolicies.value.filter(p => p.id !== policyId);
      }
    }
  } catch (error) {
    console.error('Failed to remove policy:', error);
  }
};

const closeModal = () => {
  showCreateModal.value = false;
  editingRole.value = null;
  roleForm.value = {
    name: '',
    description: '',
    realm_id: 'default',
  };
};

const closePolicyModal = () => {
  showPolicyModal.value = false;
  selectedRole.value = null;
  availablePolicies.value = [];
  assignedPolicies.value = [];
};

onMounted(() => {
  loadRoles();
});
</script>

<style scoped>
.modal {
  background-color: rgba(0, 0, 0, 0.5);
}

.list-group-item {
  cursor: pointer;
}

.list-group-item:hover {
  background-color: #f8f9fa;
}
</style>
