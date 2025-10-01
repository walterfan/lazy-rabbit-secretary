<template>
  <div class="container mt-4">
    <div class="row">
      <div class="col-12">
        <div class="d-flex justify-content-between align-items-center mb-4">
          <h2>
            <i class="bi bi-key"></i> {{ $t('admin.permissionManagement') }}
          </h2>
          <button class="btn btn-primary" @click="showCreateModal = true">
            <i class="bi bi-plus"></i> {{ $t('admin.addPolicy') }}
          </button>
        </div>

        <!-- Loading State -->
        <div v-if="policyStore.loading" class="text-center py-4">
          <div class="spinner-border" role="status">
            <span class="visually-hidden">{{ $t('common.loading') }}</span>
          </div>
        </div>

        <!-- Error State -->
        <div v-else-if="policyStore.error" class="alert alert-danger" role="alert">
          <i class="bi bi-exclamation-triangle"></i> {{ policyStore.error }}
        </div>

        <!-- Policy List -->
        <div v-else class="card">
          <div class="card-header">
            <div class="row align-items-center">
              <div class="col">
                <h5 class="mb-0">{{ $t('admin.policies') }} ({{ policyStore.totalPolicies }})</h5>
              </div>
              <div class="col-auto">
                <input
                  type="text"
                  class="form-control"
                  :placeholder="$t('admin.searchPolicies')"
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
                    <th>{{ $t('admin.version') }}</th>
                    <th>{{ $t('admin.statements') }}</th>
                    <th>{{ $t('admin.created') }}</th>
                    <th>{{ $t('admin.actions') }}</th>
                  </tr>
                </thead>
                <tbody>
                  <tr v-for="policy in filteredPolicies" :key="policy.id">
                    <td>
                      <strong>{{ policy.name }}</strong>
                    </td>
                    <td>{{ policy.description || $t('admin.noDescription') }}</td>
                    <td>
                      <span class="badge bg-secondary">{{ policy.realm_id }}</span>
                    </td>
                    <td>
                      <span class="badge bg-info">{{ policy.version }}</span>
                    </td>
                    <td>
                      <button 
                        class="btn btn-sm btn-outline-secondary" 
                        @click="viewStatements(policy)"
                      >
                        <i class="bi bi-list"></i> {{ $t('admin.viewStatements') }}
                      </button>
                    </td>
                    <td>{{ formatDate(policy.created_at) }}</td>
                    <td>
                      <button class="btn btn-sm btn-outline-primary me-1" @click="editPolicy(policy)">
                        <i class="bi bi-pencil"></i>
                      </button>
                      <button class="btn btn-sm btn-outline-info me-1" @click="manageStatements(policy)">
                        <i class="bi bi-plus-square"></i>
                      </button>
                      <button class="btn btn-sm btn-outline-danger" @click="deletePolicy(policy.id)">
                        <i class="bi bi-trash"></i>
                      </button>
                    </td>
                  </tr>
                </tbody>
              </table>
            </div>

            <!-- Pagination -->
            <nav v-if="policyStore.totalPages > 1" aria-label="Policy pagination">
              <ul class="pagination justify-content-center">
                <li class="page-item" :class="{ disabled: policyStore.currentPage === 1 }">
                  <button class="page-link" @click="changePage(policyStore.currentPage - 1)" :disabled="policyStore.currentPage === 1">
                    {{ $t('common.previous') }}
                  </button>
                </li>
                <li 
                  v-for="page in visiblePages" 
                  :key="page" 
                  class="page-item" 
                  :class="{ active: page === policyStore.currentPage }"
                >
                  <button class="page-link" @click="changePage(page)">{{ page }}</button>
                </li>
                <li class="page-item" :class="{ disabled: policyStore.currentPage === policyStore.totalPages }">
                  <button class="page-link" @click="changePage(policyStore.currentPage + 1)" :disabled="policyStore.currentPage === policyStore.totalPages">
                    {{ $t('common.next') }}
                  </button>
                </li>
              </ul>
            </nav>
          </div>
        </div>
      </div>
    </div>

    <!-- Create/Edit Policy Modal -->
    <div class="modal fade" :class="{ show: showCreateModal }" :style="{ display: showCreateModal ? 'block' : 'none' }" tabindex="-1">
      <div class="modal-dialog">
        <div class="modal-content">
          <div class="modal-header">
            <h5 class="modal-title">{{ editingPolicy ? $t('admin.editPolicy') : $t('admin.createPolicy') }}</h5>
            <button type="button" class="btn-close" @click="closeModal"></button>
          </div>
          <div class="modal-body">
            <form @submit.prevent="savePolicy">
              <div class="mb-3">
                <label for="name" class="form-label">{{ $t('admin.policyName') }}</label>
                <input
                  type="text"
                  class="form-control"
                  id="name"
                  v-model="policyForm.name"
                  required
                />
              </div>
              <div class="mb-3">
                <label for="description" class="form-label">{{ $t('admin.description') }}</label>
                <textarea
                  class="form-control"
                  id="description"
                  v-model="policyForm.description"
                  rows="3"
                  :placeholder="$t('admin.policyDescriptionPlaceholder')"
                ></textarea>
              </div>
              <div class="mb-3">
                <label for="realmName" class="form-label">{{ $t('admin.realmName') }}</label>
                <input
                  type="text"
                  class="form-control"
                  id="realmName"
                  v-model="policyForm.realm_name"
                  required
                  :disabled="!!editingPolicy"
                />
                <div v-if="editingPolicy" class="form-text">{{ $t('admin.realmCannotBeChanged') }}</div>
              </div>
              <div class="mb-3">
                <label for="version" class="form-label">{{ $t('admin.version') }}</label>
                <input
                  type="text"
                  class="form-control"
                  id="version"
                  v-model="policyForm.version"
                  placeholder="2012-10-17"
                />
                <div class="form-text">{{ $t('admin.versionHelpText') }}</div>
              </div>
            </form>
          </div>
          <div class="modal-footer">
            <button type="button" class="btn btn-secondary" @click="closeModal">{{ $t('common.cancel') }}</button>
            <button type="button" class="btn btn-primary" @click="savePolicy" :disabled="policyStore.loading">
              {{ policyStore.loading ? $t('common.saving') : $t('common.save') }}
            </button>
          </div>
        </div>
      </div>
    </div>

    <!-- Statements Management Modal -->
    <div class="modal fade" :class="{ show: showStatementsModal }" :style="{ display: showStatementsModal ? 'block' : 'none' }" tabindex="-1">
      <div class="modal-dialog modal-xl">
        <div class="modal-content">
          <div class="modal-header">
            <h5 class="modal-title">{{ $t('admin.manageStatementsFor', { policyName: selectedPolicy?.name }) }}</h5>
            <button type="button" class="btn-close" @click="closeStatementsModal"></button>
          </div>
          <div class="modal-body">
            <div class="row">
              <div class="col-12">
                <div class="d-flex justify-content-between align-items-center mb-3">
                  <h6>{{ $t('admin.statements') }}</h6>
                  <button class="btn btn-sm btn-primary" @click="showStatementForm = true">
                    <i class="bi bi-plus"></i> {{ $t('admin.addStatement') }}
                  </button>
                </div>
                
                <!-- Statement Form -->
                <div v-if="showStatementForm" class="card mb-3">
                  <div class="card-header">
                    <h6 class="mb-0">{{ editingStatement ? $t('admin.editStatement') : $t('admin.createStatement') }}</h6>
                  </div>
                  <div class="card-body">
                    <form @submit.prevent="saveStatement">
                      <div class="row">
                        <div class="col-md-6">
                          <div class="mb-3">
                            <label for="sid" class="form-label">{{ $t('admin.statementId') }}</label>
                            <input
                              type="text"
                              class="form-control"
                              id="sid"
                              v-model="statementForm.sid"
                              :placeholder="$t('admin.statementIdPlaceholder')"
                            />
                          </div>
                        </div>
                        <div class="col-md-6">
                          <div class="mb-3">
                            <label for="effect" class="form-label">{{ $t('admin.effect') }}</label>
                            <select class="form-select" id="effect" v-model="statementForm.effect" required>
                              <option value="Allow">{{ $t('admin.allow') }}</option>
                              <option value="Deny">{{ $t('admin.deny') }}</option>
                            </select>
                          </div>
                        </div>
                      </div>
                      <div class="mb-3">
                        <label for="actions" class="form-label">{{ $t('admin.actions') }}</label>
                        <textarea
                          class="form-control"
                          id="actions"
                          v-model="actionsText"
                          rows="3"
                          :placeholder="$t('admin.actionsPlaceholder')"
                          required
                        ></textarea>
                        <div class="form-text">{{ $t('admin.actionsHelpText') }}</div>
                      </div>
                      <div class="mb-3">
                        <label for="resources" class="form-label">{{ $t('admin.resources') }}</label>
                        <textarea
                          class="form-control"
                          id="resources"
                          v-model="resourcesText"
                          rows="3"
                          :placeholder="$t('admin.resourcesPlaceholder')"
                          required
                        ></textarea>
                        <div class="form-text">{{ $t('admin.resourcesHelpText') }}</div>
                      </div>
                      <div class="mb-3">
                        <label for="conditions" class="form-label">{{ $t('admin.conditions') }}</label>
                        <textarea
                          class="form-control"
                          id="conditions"
                          v-model="conditionsText"
                          rows="3"
                          :placeholder="$t('admin.conditionsPlaceholder')"
                        ></textarea>
                        <div class="form-text">{{ $t('admin.conditionsHelpText') }}</div>
                      </div>
                      <div class="d-flex gap-2">
                        <button type="button" class="btn btn-secondary" @click="cancelStatementForm">
                          {{ $t('common.cancel') }}
                        </button>
                        <button type="submit" class="btn btn-primary" :disabled="policyStore.loading">
                          {{ policyStore.loading ? $t('common.saving') : $t('common.save') }}
                        </button>
                      </div>
                    </form>
                  </div>
                </div>

                <!-- Statements List -->
                <div class="table-responsive">
                  <table class="table table-sm">
                    <thead>
                      <tr>
                        <th>{{ $t('admin.statementId') }}</th>
                        <th>{{ $t('admin.effect') }}</th>
                        <th>{{ $t('admin.actions') }}</th>
                        <th>{{ $t('admin.resources') }}</th>
                        <th>{{ $t('admin.conditions') }}</th>
                        <th>{{ $t('admin.actions') }}</th>
                      </tr>
                    </thead>
                    <tbody>
                      <tr v-for="statement in policyStore.statements" :key="statement.id">
                        <td>{{ statement.sid || $t('admin.noId') }}</td>
                        <td>
                          <span :class="statement.effect === 'Allow' ? 'badge bg-success' : 'badge bg-danger'">
                            {{ statement.effect }}
                          </span>
                        </td>
                        <td>
                          <div class="text-truncate" style="max-width: 200px;" :title="statement.actions">
                            {{ formatJsonArray(statement.actions) }}
                          </div>
                        </td>
                        <td>
                          <div class="text-truncate" style="max-width: 200px;" :title="statement.resources">
                            {{ formatJsonArray(statement.resources) }}
                          </div>
                        </td>
                        <td>
                          <div class="text-truncate" style="max-width: 150px;" :title="statement.conditions">
                            {{ statement.conditions || $t('admin.none') }}
                          </div>
                        </td>
                        <td>
                          <button class="btn btn-sm btn-outline-primary me-1" @click="editStatement(statement)">
                            <i class="bi bi-pencil"></i>
                          </button>
                          <button class="btn btn-sm btn-outline-danger" @click="deleteStatement(statement.id)">
                            <i class="bi bi-trash"></i>
                          </button>
                        </td>
                      </tr>
                      <tr v-if="policyStore.statements.length === 0">
                        <td colspan="6" class="text-center text-muted">{{ $t('admin.noStatements') }}</td>
                      </tr>
                    </tbody>
                  </table>
                </div>
              </div>
            </div>
          </div>
          <div class="modal-footer">
            <button type="button" class="btn btn-secondary" @click="closeStatementsModal">{{ $t('common.close') }}</button>
          </div>
        </div>
      </div>
    </div>

    <div v-if="showCreateModal || showStatementsModal" class="modal-backdrop fade show"></div>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted } from 'vue';
import { format } from 'date-fns';
import { useI18n } from 'vue-i18n';
import { 
  usePolicyStore, 
  type Policy, 
  type Statement,
  type CreatePolicyRequest, 
  type UpdatePolicyRequest,
  type CreateStatementRequest,
  type UpdateStatementRequest
} from '@/stores/policyStore';

const { t } = useI18n();
const policyStore = usePolicyStore();

// Local state
const searchTerm = ref('');
const showCreateModal = ref(false);
const showStatementsModal = ref(false);
const showStatementForm = ref(false);
const editingPolicy = ref<Policy | null>(null);
const selectedPolicy = ref<Policy | null>(null);
const editingStatement = ref<Statement | null>(null);

const policyForm = ref<CreatePolicyRequest>({
  realm_name: 'default',
  name: '',
  description: '',
  version: '2012-10-17',
});

const statementForm = ref({
  sid: '',
  effect: 'Allow' as 'Allow' | 'Deny',
  actions: [] as string[],
  resources: [] as string[],
  conditions: {} as Record<string, any>,
});

// Text representations for form inputs
const actionsText = ref('');
const resourcesText = ref('');
const conditionsText = ref('');

// Computed
const filteredPolicies = computed(() => {
  return policyStore.searchPolicies(searchTerm.value);
});

const visiblePages = computed(() => {
  const pages = [];
  const current = policyStore.currentPage;
  const total = policyStore.totalPages;
  
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

const formatJsonArray = (jsonString: string) => {
  try {
    const array = JSON.parse(jsonString);
    return Array.isArray(array) ? array.join(', ') : jsonString;
  } catch {
    return jsonString;
  }
};

const loadPolicies = async (page = 1) => {
  try {
    await policyStore.getPolicies({
      page,
      page_size: policyStore.pagination.page_size,
    });
  } catch (err) {
    console.error('Failed to load policies:', err);
  }
};

const changePage = (page: number) => {
  if (page >= 1 && page <= policyStore.totalPages) {
    loadPolicies(page);
  }
};

const editPolicy = (policy: Policy) => {
  editingPolicy.value = policy;
  policyForm.value = {
    realm_name: 'default', // Cannot change realm for existing policies
    name: policy.name,
    description: policy.description,
    version: policy.version,
  };
  showCreateModal.value = true;
};

const savePolicy = async () => {
  try {
    if (editingPolicy.value) {
      const updateData: UpdatePolicyRequest = {
        name: policyForm.value.name,
        description: policyForm.value.description,
        version: policyForm.value.version,
      };
      await policyStore.updatePolicy(editingPolicy.value.id, updateData);
    } else {
      await policyStore.createPolicy(policyForm.value);
    }
    
    await loadPolicies(policyStore.currentPage);
    closeModal();
  } catch (err) {
    console.error('Failed to save policy:', err);
  }
};

const deletePolicy = async (policyId: string) => {
  if (!confirm(t('admin.confirmDeletePolicy'))) return;
  
  try {
    await policyStore.deletePolicy(policyId);
    await loadPolicies(policyStore.currentPage);
  } catch (err) {
    console.error('Failed to delete policy:', err);
  }
};

const viewStatements = async (policy: Policy) => {
  try {
    await policyStore.getPolicy(policy.id);
    selectedPolicy.value = policy;
    showStatementsModal.value = true;
  } catch (err) {
    console.error('Failed to load statements:', err);
  }
};

const manageStatements = async (policy: Policy) => {
  try {
    await policyStore.getPolicy(policy.id);
    selectedPolicy.value = policy;
    showStatementsModal.value = true;
    showStatementForm.value = true;
  } catch (err) {
    console.error('Failed to load statements:', err);
  }
};

const editStatement = (statement: Statement) => {
  editingStatement.value = statement;
  statementForm.value = {
    sid: statement.sid,
    effect: statement.effect,
    actions: policyStore.parseStatementActions(statement.actions),
    resources: policyStore.parseStatementResources(statement.resources),
    conditions: policyStore.parseStatementConditions(statement.conditions),
  };
  
  // Update text representations
  actionsText.value = statementForm.value.actions.join('\n');
  resourcesText.value = statementForm.value.resources.join('\n');
  conditionsText.value = Object.keys(statementForm.value.conditions).length > 0 
    ? JSON.stringify(statementForm.value.conditions, null, 2) 
    : '';
    
  showStatementForm.value = true;
};

const saveStatement = async () => {
  if (!selectedPolicy.value) return;
  
  try {
    // Parse text inputs into arrays/objects
    const actions = actionsText.value.split('\n').map(s => s.trim()).filter(s => s);
    const resources = resourcesText.value.split('\n').map(s => s.trim()).filter(s => s);
    let conditions = {};
    
    if (conditionsText.value.trim()) {
      try {
        conditions = JSON.parse(conditionsText.value);
      } catch {
        alert(t('admin.invalidJsonConditions'));
        return;
      }
    }
    
    if (editingStatement.value) {
      const updateData: UpdateStatementRequest = {
        sid: statementForm.value.sid,
        effect: statementForm.value.effect,
        actions,
        resources,
        conditions: Object.keys(conditions).length > 0 ? conditions : undefined,
      };
      await policyStore.updateStatement(editingStatement.value.id, updateData);
    } else {
      const createData: CreateStatementRequest = {
        policy_id: selectedPolicy.value.id,
        sid: statementForm.value.sid,
        effect: statementForm.value.effect,
        actions,
        resources,
        conditions: Object.keys(conditions).length > 0 ? conditions : undefined,
      };
      await policyStore.createStatement(createData);
    }
    
    cancelStatementForm();
  } catch (err) {
    console.error('Failed to save statement:', err);
  }
};

const deleteStatement = async (statementId: string) => {
  if (!confirm(t('admin.confirmDeleteStatement'))) return;
  
  try {
    await policyStore.deleteStatement(statementId);
  } catch (err) {
    console.error('Failed to delete statement:', err);
  }
};

const cancelStatementForm = () => {
  showStatementForm.value = false;
  editingStatement.value = null;
  statementForm.value = {
    sid: '',
    effect: 'Allow',
    actions: [],
    resources: [],
    conditions: {},
  };
  actionsText.value = '';
  resourcesText.value = '';
  conditionsText.value = '';
};

const closeModal = () => {
  showCreateModal.value = false;
  editingPolicy.value = null;
  policyForm.value = {
    realm_name: 'default',
    name: '',
    description: '',
    version: '2012-10-17',
  };
};

const closeStatementsModal = () => {
  showStatementsModal.value = false;
  selectedPolicy.value = null;
  cancelStatementForm();
};

// Lifecycle
onMounted(() => {
  loadPolicies();
});
</script>

<style scoped>
.modal {
  background-color: rgba(0, 0, 0, 0.5);
}

.text-truncate {
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}
</style>