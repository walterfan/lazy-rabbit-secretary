<template>
  <div class="container mt-4">
    <div class="row">
      <div class="col-12">
        <div class="d-flex justify-content-between align-items-center mb-4">
          <h2>
            <i class="bi bi-key"></i> Permission Management
          </h2>
          <button class="btn btn-primary" @click="showCreateModal = true">
            <i class="bi bi-plus"></i> Add Policy
          </button>
        </div>

        <!-- Policy List -->
        <div class="card">
          <div class="card-header">
            <div class="row align-items-center">
              <div class="col">
                <h5 class="mb-0">Policies</h5>
              </div>
              <div class="col-auto">
                <input
                  type="text"
                  class="form-control"
                  placeholder="Search policies..."
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
                    <th>Version</th>
                    <th>Statements</th>
                    <th>Created</th>
                    <th>Actions</th>
                  </tr>
                </thead>
                <tbody>
                  <tr v-for="policy in filteredPolicies" :key="policy.id">
                    <td>{{ policy.name }}</td>
                    <td>{{ policy.description }}</td>
                    <td>{{ policy.realm_id }}</td>
                    <td>{{ policy.version }}</td>
                    <td>
                      <span class="badge bg-info">{{ policy.statement_count || 0 }}</span>
                    </td>
                    <td>{{ formatDate(policy.created_at) }}</td>
                    <td>
                      <button class="btn btn-sm btn-outline-primary me-1" @click="editPolicy(policy)">
                        <i class="bi bi-pencil"></i>
                      </button>
                      <button class="btn btn-sm btn-outline-info me-1" @click="manageStatements(policy)">
                        <i class="bi bi-list-ul"></i>
                      </button>
                      <button class="btn btn-sm btn-outline-danger" @click="deletePolicy(policy.id)">
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

    <!-- Create/Edit Policy Modal -->
    <div class="modal fade" :class="{ show: showCreateModal }" :style="{ display: showCreateModal ? 'block' : 'none' }" tabindex="-1">
      <div class="modal-dialog">
        <div class="modal-content">
          <div class="modal-header">
            <h5 class="modal-title">{{ editingPolicy ? 'Edit Policy' : 'Create Policy' }}</h5>
            <button type="button" class="btn-close" @click="closeModal"></button>
          </div>
          <div class="modal-body">
            <form @submit.prevent="savePolicy">
              <div class="mb-3">
                <label for="name" class="form-label">Policy Name</label>
                <input
                  type="text"
                  class="form-control"
                  id="name"
                  v-model="policyForm.name"
                  required
                />
              </div>
              <div class="mb-3">
                <label for="description" class="form-label">Description</label>
                <textarea
                  class="form-control"
                  id="description"
                  v-model="policyForm.description"
                  rows="3"
                ></textarea>
              </div>
              <div class="mb-3">
                <label for="realmId" class="form-label">Realm ID</label>
                <input
                  type="text"
                  class="form-control"
                  id="realmId"
                  v-model="policyForm.realm_id"
                  required
                />
              </div>
              <div class="mb-3">
                <label for="version" class="form-label">Version</label>
                <input
                  type="text"
                  class="form-control"
                  id="version"
                  v-model="policyForm.version"
                  placeholder="2012-10-17"
                />
              </div>
            </form>
          </div>
          <div class="modal-footer">
            <button type="button" class="btn btn-secondary" @click="closeModal">Cancel</button>
            <button type="button" class="btn btn-primary" @click="savePolicy" :disabled="loading">
              {{ loading ? 'Saving...' : 'Save' }}
            </button>
          </div>
        </div>
      </div>
    </div>

    <!-- Statement Management Modal -->
    <div class="modal fade" :class="{ show: showStatementModal }" :style="{ display: showStatementModal ? 'block' : 'none' }" tabindex="-1">
      <div class="modal-dialog modal-xl">
        <div class="modal-content">
          <div class="modal-header">
            <h5 class="modal-title">Manage Statements for {{ selectedPolicy?.name }}</h5>
            <button type="button" class="btn-close" @click="closeStatementModal"></button>
          </div>
          <div class="modal-body">
            <div class="d-flex justify-content-between align-items-center mb-3">
              <h6>Policy Statements</h6>
              <button class="btn btn-sm btn-primary" @click="showAddStatementModal = true">
                <i class="bi bi-plus"></i> Add Statement
              </button>
            </div>
            
            <div class="table-responsive">
              <table class="table table-sm">
                <thead>
                  <tr>
                    <th>SID</th>
                    <th>Effect</th>
                    <th>Actions</th>
                    <th>Resources</th>
                    <th>Conditions</th>
                    <th>Actions</th>
                  </tr>
                </thead>
                <tbody>
                  <tr v-for="statement in statements" :key="statement.id">
                    <td>{{ statement.sid || '-' }}</td>
                    <td>
                      <span :class="statement.effect === 'Allow' ? 'badge bg-success' : 'badge bg-danger'">
                        {{ statement.effect }}
                      </span>
                    </td>
                    <td>
                      <small class="text-muted">{{ formatArray(statement.actions) }}</small>
                    </td>
                    <td>
                      <small class="text-muted">{{ formatArray(statement.resources) }}</small>
                    </td>
                    <td>
                      <small class="text-muted">{{ formatConditions(statement.conditions) }}</small>
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
                </tbody>
              </table>
            </div>
          </div>
          <div class="modal-footer">
            <button type="button" class="btn btn-secondary" @click="closeStatementModal">Close</button>
          </div>
        </div>
      </div>
    </div>

    <!-- Add/Edit Statement Modal -->
    <div class="modal fade" :class="{ show: showAddStatementModal }" :style="{ display: showAddStatementModal ? 'block' : 'none' }" tabindex="-1">
      <div class="modal-dialog modal-lg">
        <div class="modal-content">
          <div class="modal-header">
            <h5 class="modal-title">{{ editingStatement ? 'Edit Statement' : 'Add Statement' }}</h5>
            <button type="button" class="btn-close" @click="closeAddStatementModal"></button>
          </div>
          <div class="modal-body">
            <form @submit.prevent="saveStatement">
              <div class="row">
                <div class="col-md-6">
                  <div class="mb-3">
                    <label for="sid" class="form-label">Statement ID (SID)</label>
                    <input
                      type="text"
                      class="form-control"
                      id="sid"
                      v-model="statementForm.sid"
                      placeholder="Optional"
                    />
                  </div>
                </div>
                <div class="col-md-6">
                  <div class="mb-3">
                    <label for="effect" class="form-label">Effect</label>
                    <select class="form-select" id="effect" v-model="statementForm.effect" required>
                      <option value="Allow">Allow</option>
                      <option value="Deny">Deny</option>
                    </select>
                  </div>
                </div>
              </div>
              <div class="mb-3">
                <label for="actions" class="form-label">Actions (one per line)</label>
                <textarea
                  class="form-control"
                  id="actions"
                  v-model="statementForm.actionsText"
                  rows="3"
                  placeholder="read:project&#10;write:project&#10;delete:project"
                  required
                ></textarea>
              </div>
              <div class="mb-3">
                <label for="resources" class="form-label">Resources (one per line)</label>
                <textarea
                  class="form-control"
                  id="resources"
                  v-model="statementForm.resourcesText"
                  rows="3"
                  placeholder="project:*&#10;user:${user:id}"
                  required
                ></textarea>
              </div>
              <div class="mb-3">
                <label for="conditions" class="form-label">Conditions (JSON)</label>
                <textarea
                  class="form-control"
                  id="conditions"
                  v-model="statementForm.conditionsText"
                  rows="4"
                  placeholder='{&#10;  "StringEquals": {&#10;    "project:owner": "${user:id}"&#10;  }&#10;}'
                ></textarea>
              </div>
            </form>
          </div>
          <div class="modal-footer">
            <button type="button" class="btn btn-secondary" @click="closeAddStatementModal">Cancel</button>
            <button type="button" class="btn btn-primary" @click="saveStatement" :disabled="loading">
              {{ loading ? 'Saving...' : 'Save' }}
            </button>
          </div>
        </div>
      </div>
    </div>

    <div v-if="showCreateModal || showStatementModal || showAddStatementModal" class="modal-backdrop fade show"></div>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted } from 'vue';
import { format } from 'date-fns';

interface Policy {
  id: string;
  name: string;
  description: string;
  realm_id: string;
  version: string;
  created_at: string;
  updated_at: string;
  statement_count?: number;
}

interface Statement {
  id: string;
  policy_id: string;
  sid: string;
  effect: string;
  actions: string[];
  resources: string[];
  conditions: Record<string, any>;
}

const policies = ref<Policy[]>([]);
const statements = ref<Statement[]>([]);
const searchTerm = ref('');
const showCreateModal = ref(false);
const showStatementModal = ref(false);
const showAddStatementModal = ref(false);
const editingPolicy = ref<Policy | null>(null);
const selectedPolicy = ref<Policy | null>(null);
const editingStatement = ref<Statement | null>(null);
const loading = ref(false);

const policyForm = ref({
  name: '',
  description: '',
  realm_id: 'default',
  version: '2012-10-17',
});

const statementForm = ref({
  sid: '',
  effect: 'Allow',
  actionsText: '',
  resourcesText: '',
  conditionsText: '',
});

const filteredPolicies = computed(() => {
  if (!searchTerm.value) return policies.value;
  return policies.value.filter(policy =>
    policy.name.toLowerCase().includes(searchTerm.value.toLowerCase()) ||
    policy.description.toLowerCase().includes(searchTerm.value.toLowerCase())
  );
});

const formatDate = (dateString: string) => {
  return format(new Date(dateString), 'MMM dd, yyyy');
};

const formatArray = (arr: string[]) => {
  return arr.join(', ');
};

const formatConditions = (conditions: Record<string, any>) => {
  if (!conditions || Object.keys(conditions).length === 0) return '-';
  return JSON.stringify(conditions);
};

const loadPolicies = async () => {
  try {
    const response = await fetch('/api/v1/admin/policies', {
      headers: {
        'Authorization': `Bearer ${localStorage.getItem('access_token')}`,
      },
    });
    if (response.ok) {
      policies.value = await response.json();
    }
  } catch (error) {
    console.error('Failed to load policies:', error);
  }
};

const editPolicy = (policy: Policy) => {
  editingPolicy.value = policy;
  policyForm.value = {
    name: policy.name,
    description: policy.description,
    realm_id: policy.realm_id,
    version: policy.version,
  };
  showCreateModal.value = true;
};

const savePolicy = async () => {
  loading.value = true;
  try {
    const url = editingPolicy.value 
      ? `/api/v1/admin/policies/${editingPolicy.value.id}`
      : '/api/v1/admin/policies';
    
    const method = editingPolicy.value ? 'PUT' : 'POST';
    
    const response = await fetch(url, {
      method,
      headers: {
        'Content-Type': 'application/json',
        'Authorization': `Bearer ${localStorage.getItem('access_token')}`,
      },
      body: JSON.stringify(policyForm.value),
    });

    if (response.ok) {
      await loadPolicies();
      closeModal();
    }
  } catch (error) {
    console.error('Failed to save policy:', error);
  } finally {
    loading.value = false;
  }
};

const deletePolicy = async (policyId: string) => {
  if (!confirm('Are you sure you want to delete this policy?')) return;
  
  try {
    const response = await fetch(`/api/v1/admin/policies/${policyId}`, {
      method: 'DELETE',
      headers: {
        'Authorization': `Bearer ${localStorage.getItem('access_token')}`,
      },
    });
    
    if (response.ok) {
      await loadPolicies();
    }
  } catch (error) {
    console.error('Failed to delete policy:', error);
  }
};

const manageStatements = async (policy: Policy) => {
  selectedPolicy.value = policy;
  showStatementModal.value = true;
  
  try {
    const response = await fetch(`/api/v1/admin/policies/${policy.id}/statements`, {
      headers: {
        'Authorization': `Bearer ${localStorage.getItem('access_token')}`,
      },
    });
    
    if (response.ok) {
      statements.value = await response.json();
    }
  } catch (error) {
    console.error('Failed to load statements:', error);
  }
};

const editStatement = (statement: Statement) => {
  editingStatement.value = statement;
  statementForm.value = {
    sid: statement.sid || '',
    effect: statement.effect,
    actionsText: statement.actions.join('\n'),
    resourcesText: statement.resources.join('\n'),
    conditionsText: JSON.stringify(statement.conditions, null, 2),
  };
  showAddStatementModal.value = true;
};

const saveStatement = async () => {
  if (!selectedPolicy.value) return;
  
  loading.value = true;
  try {
    const actions = statementForm.value.actionsText.split('\n').filter(a => a.trim());
    const resources = statementForm.value.resourcesText.split('\n').filter(r => r.trim());
    let conditions = {};
    
    if (statementForm.value.conditionsText.trim()) {
      try {
        conditions = JSON.parse(statementForm.value.conditionsText);
      } catch (e) {
        alert('Invalid JSON in conditions field');
        return;
      }
    }
    
    const statementData = {
      sid: statementForm.value.sid || undefined,
      effect: statementForm.value.effect,
      actions,
      resources,
      conditions,
    };
    
    const url = editingStatement.value 
      ? `/api/v1/admin/policies/${selectedPolicy.value.id}/statements/${editingStatement.value.id}`
      : `/api/v1/admin/policies/${selectedPolicy.value.id}/statements`;
    
    const method = editingStatement.value ? 'PUT' : 'POST';
    
    const response = await fetch(url, {
      method,
      headers: {
        'Content-Type': 'application/json',
        'Authorization': `Bearer ${localStorage.getItem('access_token')}`,
      },
      body: JSON.stringify(statementData),
    });

    if (response.ok) {
      await manageStatements(selectedPolicy.value);
      closeAddStatementModal();
    }
  } catch (error) {
    console.error('Failed to save statement:', error);
  } finally {
    loading.value = false;
  }
};

const deleteStatement = async (statementId: string) => {
  if (!selectedPolicy.value || !confirm('Are you sure you want to delete this statement?')) return;
  
  try {
    const response = await fetch(`/api/v1/admin/policies/${selectedPolicy.value.id}/statements/${statementId}`, {
      method: 'DELETE',
      headers: {
        'Authorization': `Bearer ${localStorage.getItem('access_token')}`,
      },
    });
    
    if (response.ok) {
      await manageStatements(selectedPolicy.value);
    }
  } catch (error) {
    console.error('Failed to delete statement:', error);
  }
};

const closeModal = () => {
  showCreateModal.value = false;
  editingPolicy.value = null;
  policyForm.value = {
    name: '',
    description: '',
    realm_id: 'default',
    version: '2012-10-17',
  };
};

const closeStatementModal = () => {
  showStatementModal.value = false;
  selectedPolicy.value = null;
  statements.value = [];
};

const closeAddStatementModal = () => {
  showAddStatementModal.value = false;
  editingStatement.value = null;
  statementForm.value = {
    sid: '',
    effect: 'Allow',
    actionsText: '',
    resourcesText: '',
    conditionsText: '',
  };
};

onMounted(() => {
  loadPolicies();
});
</script>

<style scoped>
.modal {
  background-color: rgba(0, 0, 0, 0.5);
}

.table-sm td {
  font-size: 0.875rem;
}
</style>
