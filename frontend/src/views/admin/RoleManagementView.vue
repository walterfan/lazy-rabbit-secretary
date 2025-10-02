<template>
  <div class="container mt-4">
    <div class="row">
      <div class="col-12">
        <div class="d-flex justify-content-between align-items-center mb-4">
          <h2>
            <i class="bi bi-shield"></i> {{ $t('admin.roleManagement') }}
          </h2>
          <button class="btn btn-primary" @click="showCreateModal = true">
            <i class="bi bi-plus"></i> {{ $t('admin.addRole') }}
          </button>
        </div>

        <!-- Loading State -->
        <div v-if="roleStore.loading" class="text-center py-4">
          <div class="spinner-border" role="status">
            <span class="visually-hidden">{{ $t('common.loading') }}</span>
          </div>
        </div>

        <!-- Error State -->
        <div v-else-if="roleStore.error" class="alert alert-danger" role="alert">
          <i class="bi bi-exclamation-triangle"></i> {{ roleStore.error }}
        </div>

        <!-- Role List -->
        <div v-else class="card">
          <div class="card-header">
            <div class="row align-items-center">
              <div class="col">
                <h5 class="mb-0">{{ $t('admin.roles') }} ({{ roleStore.totalRoles }})</h5>
              </div>
              <div class="col-auto">
                <input
                  type="text"
                  class="form-control"
                  :placeholder="$t('admin.searchRoles')"
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
                    <th>{{ $t('admin.name') }}</th>
                    <th>{{ $t('admin.description') }}</th>
                    <th>{{ $t('admin.realmId') }}</th>
                    <th>{{ $t('admin.policies') }}</th>
                    <th>{{ $t('admin.created') }}</th>
                    <th>{{ $t('admin.actions') }}</th>
                  </tr>
                </thead>
                <tbody>
                  <tr v-for="role in filteredRoles" :key="role.id">
                    <td>
                      <strong>{{ role.name }}</strong>
                    </td>
                    <td>{{ role.description || $t('admin.noDescription') }}</td>
                    <td>
                      <span class="badge bg-secondary">{{ role.realm_id }}</span>
                    </td>
                    <td>
                      <div v-if="role.policies && role.policies.length > 0" class="d-flex flex-wrap gap-1">
                        <span 
                          v-for="policy in role.policies" 
                          :key="policy.id"
                          class="badge bg-info cursor-pointer"
                          @click="viewPolicy(policy.id)"
                          :title="policy.description"
                        >
                          {{ policy.name }}
                        </span>
                      </div>
                      <span v-else class="text-muted">{{ $t('admin.noPolicies') }}</span>
                    </td>
                    <td>{{ formatDate(role.created_at) }}</td>
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

            <!-- Pagination -->
            <nav v-if="roleStore.totalPages > 1" aria-label="Role pagination">
              <ul class="pagination justify-content-center">
                <li class="page-item" :class="{ disabled: roleStore.currentPage === 1 }">
                  <button class="page-link" @click="changePage(roleStore.currentPage - 1)" :disabled="roleStore.currentPage === 1">
                    {{ $t('common.previous') }}
                  </button>
                </li>
                <li 
                  v-for="page in visiblePages" 
                  :key="page" 
                  class="page-item" 
                  :class="{ active: page === roleStore.currentPage }"
                >
                  <button class="page-link" @click="changePage(page)">{{ page }}</button>
                </li>
                <li class="page-item" :class="{ disabled: roleStore.currentPage === roleStore.totalPages }">
                  <button class="page-link" @click="changePage(roleStore.currentPage + 1)" :disabled="roleStore.currentPage === roleStore.totalPages">
                    {{ $t('common.next') }}
                  </button>
                </li>
              </ul>
            </nav>
          </div>
        </div>
      </div>
    </div>

    <!-- Create/Edit Role Modal -->
    <div class="modal fade" :class="{ show: showCreateModal }" :style="{ display: showCreateModal ? 'block' : 'none' }" tabindex="-1">
      <div class="modal-dialog">
        <div class="modal-content">
          <div class="modal-header">
            <h5 class="modal-title">{{ editingRole ? $t('admin.editRole') : $t('admin.createRole') }}</h5>
            <button type="button" class="btn-close" @click="closeModal"></button>
          </div>
          <div class="modal-body">
            <form @submit.prevent="saveRole">
              <div class="mb-3">
                <label for="name" class="form-label">{{ $t('admin.roleName') }}</label>
                <input
                  type="text"
                  class="form-control"
                  id="name"
                  v-model="roleForm.name"
                  required
                />
              </div>
              <div class="mb-3">
                <label for="description" class="form-label">{{ $t('admin.description') }}</label>
                <textarea
                  class="form-control"
                  id="description"
                  v-model="roleForm.description"
                  rows="3"
                  :placeholder="$t('admin.roleDescriptionPlaceholder')"
                ></textarea>
              </div>
              <div class="mb-3">
                <label for="realmName" class="form-label">{{ $t('admin.realmName') }}</label>
                <input
                  type="text"
                  class="form-control"
                  id="realmName"
                  v-model="roleForm.realm_name"
                  required
                  :disabled="!!editingRole"
                />
                <div v-if="editingRole" class="form-text">{{ $t('admin.realmCannotBeChanged') }}</div>
              </div>
            </form>
          </div>
          <div class="modal-footer">
            <button type="button" class="btn btn-secondary" @click="closeModal">{{ $t('common.cancel') }}</button>
            <button type="button" class="btn btn-primary" @click="saveRole" :disabled="roleStore.loading">
              {{ roleStore.loading ? $t('common.saving') : $t('common.save') }}
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
            <h5 class="modal-title">{{ $t('admin.managePoliciesFor', { roleName: selectedRole?.name }) }}</h5>
            <button type="button" class="btn-close" @click="closePolicyModal"></button>
          </div>
          <div class="modal-body">
            <div class="row">
              <div class="col-md-6">
                <h6>{{ $t('admin.availablePolicies') }}</h6>
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
                      <small class="text-muted">{{ policy.description || $t('admin.noDescription') }}</small>
                    </div>
                    <button class="btn btn-sm btn-outline-primary">
                      <i class="bi bi-plus"></i>
                    </button>
                  </div>
                  <div v-if="availablePolicies.length === 0" class="list-group-item text-muted text-center">
                    {{ $t('admin.noAvailablePolicies') }}
                  </div>
                </div>
              </div>
              <div class="col-md-6">
                <h6>{{ $t('admin.assignedPolicies') }}</h6>
                <div class="list-group">
                  <div
                    v-for="policy in assignedPolicies"
                    :key="policy.id"
                    class="list-group-item d-flex justify-content-between align-items-center"
                  >
                    <div>
                      <strong>{{ policy.name }}</strong>
                      <br>
                      <small class="text-muted">{{ policy.description || $t('admin.noDescription') }}</small>
                    </div>
                    <button class="btn btn-sm btn-outline-danger" @click="removePolicy(policy.id)">
                      <i class="bi bi-dash"></i>
                    </button>
                  </div>
                  <div v-if="assignedPolicies.length === 0" class="list-group-item text-muted text-center">
                    {{ $t('admin.noAssignedPolicies') }}
                  </div>
                </div>
              </div>
            </div>
          </div>
          <div class="modal-footer">
            <button type="button" class="btn btn-secondary" @click="closePolicyModal">{{ $t('common.close') }}</button>
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
import { useI18n } from 'vue-i18n';
import { useRouter } from 'vue-router';
import { useRoleStore, type Role, type CreateRoleRequest, type UpdateRoleRequest } from '@/stores/roleStore';
import { makeAuthenticatedRequest } from '@/utils/httpInterceptor';
import { getApiUrl } from '@/utils/apiConfig';

const { t } = useI18n();
const router = useRouter();
const roleStore = useRoleStore();

// Local state
const searchTerm = ref('');
const showCreateModal = ref(false);
const showPolicyModal = ref(false);
const editingRole = ref<Role | null>(null);
const selectedRole = ref<Role | null>(null);
const availablePolicies = ref<any[]>([]);
const assignedPolicies = ref<any[]>([]);

const roleForm = ref<CreateRoleRequest>({
  realm_name: 'default',
  name: '',
  description: '',
});

// Computed
const filteredRoles = computed(() => {
  return roleStore.searchRoles(searchTerm.value);
});

const visiblePages = computed(() => {
  const pages = [];
  const current = roleStore.currentPage;
  const total = roleStore.totalPages;
  
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

const loadRoles = async (page = 1) => {
  try {
    await roleStore.getRoles({
      page,
      page_size: roleStore.pagination.page_size,
    });
  } catch (err) {
    console.error('Failed to load roles:', err);
  }
};

const changePage = (page: number) => {
  if (page >= 1 && page <= roleStore.totalPages) {
    loadRoles(page);
  }
};

const editRole = (role: Role) => {
  editingRole.value = role;
  roleForm.value = {
    realm_name: 'default', // Cannot change realm for existing roles
    name: role.name,
    description: role.description,
  };
  showCreateModal.value = true;
};

const saveRole = async () => {
  try {
    if (editingRole.value) {
      const updateData: UpdateRoleRequest = {
        name: roleForm.value.name,
        description: roleForm.value.description,
      };
      await roleStore.updateRole(editingRole.value.id, updateData);
    } else {
      await roleStore.createRole(roleForm.value);
    }
    
    await loadRoles(roleStore.currentPage);
    closeModal();
  } catch (err) {
    console.error('Failed to save role:', err);
  }
};

const deleteRole = async (roleId: string) => {
  if (!confirm(t('admin.confirmDeleteRole'))) return;
  
  try {
    await roleStore.deleteRole(roleId);
    await loadRoles(roleStore.currentPage);
  } catch (err) {
    console.error('Failed to delete role:', err);
  }
};

const managePolicies = async (role: Role) => {
  selectedRole.value = role;
  showPolicyModal.value = true;
  
  try {
    // Load all policies from the API
    const url = getApiUrl('/api/v1/admin/policies?page_size=100');
    const response = await makeAuthenticatedRequest(url, {
      method: 'GET',
    });
    
    if (!response.ok) {
      throw new Error(`Failed to fetch policies: ${response.statusText}`);
    }
    
    const data = await response.json();
    const allPolicies = data.policies || [];
    
    // Get policies already assigned to this role
    assignedPolicies.value = role.policies || [];
    
    // Filter out assigned policies from available policies
    const assignedPolicyIds = new Set(assignedPolicies.value.map(p => p.id));
    availablePolicies.value = allPolicies.filter(policy => !assignedPolicyIds.has(policy.id));
    
  } catch (error) {
    console.error('Failed to load policies:', error);
    availablePolicies.value = [];
    assignedPolicies.value = [];
  }
};

const assignPolicy = async (policy: any) => {
  if (!selectedRole.value) return;
  
  try {
    // Placeholder for policy assignment - would need to implement policy assignment API
    assignedPolicies.value.push(policy);
    availablePolicies.value = availablePolicies.value.filter(p => p.id !== policy.id);
  } catch (error) {
    console.error('Failed to assign policy:', error);
  }
};

const removePolicy = async (policyId: string) => {
  if (!selectedRole.value) return;
  
  try {
    // Placeholder for policy removal - would need to implement policy removal API
    const policy = assignedPolicies.value.find(p => p.id === policyId);
    if (policy) {
      availablePolicies.value.push(policy);
      assignedPolicies.value = assignedPolicies.value.filter(p => p.id !== policyId);
    }
  } catch (error) {
    console.error('Failed to remove policy:', error);
  }
};

const closeModal = () => {
  showCreateModal.value = false;
  editingRole.value = null;
  roleForm.value = {
    realm_name: 'default',
    name: '',
    description: '',
  };
};

const closePolicyModal = () => {
  showPolicyModal.value = false;
  selectedRole.value = null;
  availablePolicies.value = [];
  assignedPolicies.value = [];
};

const viewPolicy = (policyId: string) => {
  // Navigate to permission management page with the specific policy
  router.push(`/admin/permissions?policy=${policyId}`);
};

// Lifecycle
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

.cursor-pointer {
  cursor: pointer;
}

.cursor-pointer:hover {
  opacity: 0.8;
}
</style>