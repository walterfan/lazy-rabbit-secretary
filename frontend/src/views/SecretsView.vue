`<template>
  <div class="secrets-view">
    <!-- Page Header -->
    <div class="page-header">
      <div class="header-content">
        <h1 class="page-title">
          <i class="bi bi-shield-lock-fill"></i>
          Secrets Management
        </h1>
        <p class="page-description">
          Securely store and manage sensitive information like API keys, passwords, and certificates
        </p>
      </div>
      <div class="header-actions">
        <button class="btn btn-primary" @click="showCreateForm">
          <i class="bi bi-plus-lg me-2"></i>
          New Secret
        </button>
      </div>
    </div>

    <!-- Main Content -->
    <div class="content-wrapper">
      <!-- Create/Edit Form -->
      <div v-if="showForm" class="form-section">
        <div class="section-header">
          <h2>{{ editingSecret ? 'Edit Secret' : 'Create New Secret' }}</h2>
          <button 
            class="btn btn-sm btn-light"
            @click="closeForm"
          >
            <i class="bi bi-x-lg"></i>
          </button>
        </div>
        <SecretForm
          :secret="editingSecret || undefined"
          :is-edit-mode="!!editingSecret"
          @submit="handleFormSubmit"
          @cancel="closeForm"
        />
      </div>

      <!-- Secret List -->
      <div v-else>
        <SecretList
          :secrets="secretStore.secrets"
          :search-query="searchQuery"
          :filters="filters"
          :current-page="currentPage"
          :page-size="pageSize"
          :total-count="secretStore.totalCount"
          @view="handleView"
          @edit="handleEdit"
          @delete="handleDelete"
          @copy="handleCopy"
          @copy-with-kek="openKEKModal"
          @view-versions="handleViewVersions"
          @update:searchQuery="searchQuery = $event"
          @update:filters="filters = $event"
          @update:page="currentPage = $event"
        />
      </div>
    </div>

    <!-- Loading Overlay -->
    <div v-if="secretStore.loading" class="loading-overlay">
      <div class="spinner-border text-primary" role="status">
        <span class="visually-hidden">Loading...</span>
      </div>
    </div>

    <!-- Secret Details Modal -->
    <div 
      v-if="viewingSecret"
      class="modal fade show d-block"
      tabindex="-1"
      @click.self="viewingSecret = null"
    >
      <div class="modal-dialog modal-lg">
        <div class="modal-content">
          <div class="modal-header">
            <h5 class="modal-title">
              <i class="bi bi-shield-check me-2"></i>
              Secret Details
            </h5>
            <button 
              type="button" 
              class="btn-close" 
              @click="viewingSecret = null"
            ></button>
          </div>
          <div class="modal-body">
            <div class="secret-details">
              <div class="detail-row">
                <span class="detail-label">Name:</span>
                <span class="detail-value fw-bold">{{ viewingSecret.name }}</span>
              </div>
              <div class="detail-row">
                <span class="detail-label">Group:</span>
                <span class="detail-value">
                  <span class="badge bg-secondary">{{ viewingSecret.group }}</span>
                </span>
              </div>
              <div class="detail-row">
                <span class="detail-label">Path:</span>
                <span class="detail-value">
                  <code>{{ viewingSecret.path }}</code>
                </span>
              </div>
              <div class="detail-row">
                <span class="detail-label">Description:</span>
                <span class="detail-value">{{ viewingSecret.desc || '-' }}</span>
              </div>
              <div class="detail-row">
                <span class="detail-label">Current Version:</span>
                <span class="detail-value">
                  <span class="badge bg-primary">{{ viewingSecret.current_version }}</span>
                </span>
              </div>
              <div class="detail-row">
                <span class="detail-label">Previous Version:</span>
                <span class="detail-value">
                  <span v-if="viewingSecret.previous_version > 0" class="badge bg-secondary">{{ viewingSecret.previous_version }}</span>
                  <span v-else class="text-muted">-</span>
                </span>
              </div>
              <div class="detail-row">
                <span class="detail-label">Pending Version:</span>
                <span class="detail-value">
                  <span v-if="viewingSecret.pending_version > 0" class="badge bg-warning">{{ viewingSecret.pending_version }}</span>
                  <span v-else class="text-muted">-</span>
                </span>
              </div>
              <div class="detail-row">
                <span class="detail-label">Total Versions:</span>
                <span class="detail-value">
                  <span class="badge bg-info">{{ viewingSecret.max_version }}</span>
                </span>
              </div>
              <div class="detail-row">
                <span class="detail-label">Created By:</span>
                <span class="detail-value">{{ viewingSecret.created_by }}</span>
              </div>
              <div class="detail-row">
                <span class="detail-label">Created At:</span>
                <span class="detail-value">{{ formatDate(viewingSecret.created_at) }}</span>
              </div>
              <div class="detail-row">
                <span class="detail-label">Updated At:</span>
                <span class="detail-value">{{ formatDate(viewingSecret.updated_at) }}</span>
              </div>
            </div>
          </div>
          <div class="modal-footer">
            <button 
              type="button" 
              class="btn btn-secondary" 
              @click="viewingSecret = null"
            >
              Close
            </button>
            <button 
              v-if="viewingSecret.max_version > 1"
              type="button" 
              class="btn btn-info" 
              @click="handleViewVersions(viewingSecret)"
            >
              <i class="bi bi-layers me-2"></i>
              View Versions
            </button>
            <button 
              type="button" 
              class="btn btn-primary" 
              @click="handleCopy(viewingSecret)"
            >
              <i class="bi bi-clipboard me-2"></i>
              Copy Value
            </button>
          </div>
        </div>
      </div>
    </div>
    <div v-if="viewingSecret" class="modal-backdrop fade show"></div>

    <!-- KEK Input Modal -->
    <div 
      v-if="showKEKModal"
      class="modal fade show d-block"
      tabindex="-1"
      @click.self="closeKEKModal"
    >
      <div class="modal-dialog">
        <div class="modal-content">
          <div class="modal-header">
            <h5 class="modal-title">
              <i class="bi bi-key me-2"></i>
              Decrypt Secret
              <span v-if="isKEKRequired" class="badge bg-warning ms-2">KEK Required</span>
              <span v-else class="badge bg-info ms-2">KEK Optional</span>
            </h5>
            <button 
              type="button" 
              class="btn-close" 
              @click="closeKEKModal"
            ></button>
          </div>
          <div class="modal-body">
            <!-- Required KEK Alert for version 999 -->
            <div v-if="isKEKRequired" class="alert alert-warning" role="alert">
              <i class="bi bi-exclamation-triangle me-2"></i>
              <strong>Custom KEK Required:</strong> 
              This secret was encrypted with a custom KEK password. You must provide the exact password used during encryption to decrypt it.
            </div>
            
            <!-- Optional KEK Alert for other versions -->
            <div v-else class="alert alert-info" role="alert">
              <i class="bi bi-info-circle me-2"></i>
              <strong>Optional Custom KEK:</strong> 
              This secret was encrypted with the system default KEK. Only enter a custom KEK if you specifically used one during encryption.
            </div>

            <div class="mb-3">
              <label for="kek-input" class="form-label">
                <i class="bi bi-shield-lock me-1"></i>
                Custom KEK Password/Phrase
                <span v-if="isKEKRequired" class="text-danger">*</span>
                <span v-else class="text-muted">(Optional)</span>
              </label>
              <input
                type="password"
                class="form-control"
                id="kek-input"
                v-model="kekInput"
                :placeholder="isKEKRequired ? 'Enter the custom KEK password used for encryption' : 'Leave empty to use system default KEK'"
                maxlength="256"
                :class="{ 'is-invalid': kekError }"
                :required="isKEKRequired"
              />
              <div class="invalid-feedback" v-if="kekError">
                {{ kekError }}
              </div>
              <small class="text-muted">
                <span v-if="isKEKRequired">
                  <i class="bi bi-info-circle me-1"></i>
                  Custom KEK Required
                </span>
                <span v-else>
                  <i class="bi bi-info-circle me-1"></i>
                  System Default KEK. Only enter a custom KEK if this secret was specifically encrypted with one.
                </span>
              </small>
            </div>
          </div>
          <div class="modal-footer">
            <button 
              type="button" 
              class="btn btn-secondary" 
              @click="closeKEKModal"
            >
              Cancel
            </button>
            <button 
              type="button" 
              class="btn btn-primary" 
              @click="handleKEKDecrypt"
              :disabled="kekLoading"
            >
              <span v-if="kekLoading" class="spinner-border spinner-border-sm me-2" role="status"></span>
              <i v-else class="bi bi-unlock me-2"></i>
              {{ kekLoading ? 'Decrypting...' : 'Decrypt & Copy' }}
            </button>
          </div>
        </div>
      </div>
    </div>
    <div v-if="showKEKModal" class="modal-backdrop fade show"></div>

    <!-- Versions Modal -->
    <div 
      v-if="viewingVersions"
      class="modal fade show d-block"
      tabindex="-1"
      @click.self="viewingVersions = null"
    >
      <div class="modal-dialog modal-xl">
        <div class="modal-content">
          <div class="modal-header">
            <h5 class="modal-title">
              <i class="bi bi-layers me-2"></i>
              Secret Versions: {{ viewingVersions.name }}
            </h5>
            <button 
              type="button" 
              class="btn-close" 
              @click="viewingVersions = null"
            ></button>
          </div>
          <div class="modal-body">
            <div v-if="versionsLoading" class="text-center py-4">
              <div class="spinner-border text-primary" role="status">
                <span class="visually-hidden">Loading...</span>
              </div>
            </div>
            <div v-else>
              <div class="versions-header mb-3">
                <div class="row">
                  <div class="col-md-6">
                    <h6>Version Summary</h6>
                    <div class="version-stats">
                      <span class="badge bg-primary me-2">Current: v{{ viewingVersions.current_version }}</span>
                      <span v-if="viewingVersions.previous_version > 0" class="badge bg-secondary me-2">Previous: v{{ viewingVersions.previous_version }}</span>
                      <span v-if="viewingVersions.pending_version > 0" class="badge bg-warning me-2">Pending: v{{ viewingVersions.pending_version }}</span>
                      <span class="badge bg-info">Total: {{ viewingVersions.max_version }}</span>
                    </div>
                  </div>
                  <div class="col-md-6 text-end">
                    <button 
                      class="btn btn-primary btn-sm"
                      @click="showCreatePendingForm"
                    >
                      <i class="bi bi-plus-lg me-1"></i>
                      Create Pending Version
                    </button>
                  </div>
                </div>
              </div>
              
              <div class="table-responsive">
                <table class="table table-hover">
                  <thead>
                    <tr>
                      <th>Version</th>
                      <th>Status</th>
                      <th>KEK Version</th>
                      <th>Created By</th>
                      <th>Created At</th>
                      <th>Actions</th>
                    </tr>
                  </thead>
                  <tbody>
                    <tr v-for="version in secretVersions" :key="version.id" class="version-row">
                      <td>
                        <span class="fw-medium">v{{ version.version }}</span>
                        <span v-if="version.version === viewingVersions.current_version" class="badge bg-primary ms-2">Current</span>
                        <span v-if="version.version === viewingVersions.pending_version" class="badge bg-warning ms-2">Pending</span>
                      </td>
                      <td>
                        <span :class="`badge bg-${getStatusColor(version.status)}`">
                          {{ version.status }}
                        </span>
                      </td>
                      <td>
                        <span v-if="version.kek_version === 999" class="badge bg-warning">Custom KEK</span>
                        <span v-else class="badge bg-info">v{{ version.kek_version }}</span>
                      </td>
                      <td>{{ version.created_by }}</td>
                      <td>{{ formatDate(version.created_at) }}</td>
                      <td>
                        <div class="action-buttons">
                          <button
                            class="btn btn-sm btn-outline-primary"
                            @click="handleCopyVersion(version)"
                            title="Copy version value"
                          >
                            <i class="bi bi-clipboard"></i>
                          </button>
                          <button
                            class="btn btn-sm btn-outline-info"
                            @click="handleCopyVersionWithKEK(version)"
                            title="Copy with custom KEK"
                          >
                            <i class="bi bi-key"></i>
                          </button>
                          <button
                            v-if="version.status === 'pending'"
                            class="btn btn-sm btn-outline-success"
                            @click="handleActivateVersion(version.version)"
                            title="Activate version"
                          >
                            <i class="bi bi-check-circle"></i>
                          </button>
                          <button
                            v-if="version.version !== viewingVersions.current_version"
                            class="btn btn-sm btn-outline-danger"
                            @click="handleDeleteVersion(version.version)"
                            title="Delete version"
                          >
                            <i class="bi bi-trash"></i>
                          </button>
                        </div>
                      </td>
                    </tr>
                  </tbody>
                </table>
              </div>
            </div>
          </div>
          <div class="modal-footer">
            <button 
              type="button" 
              class="btn btn-secondary" 
              @click="viewingVersions = null"
            >
              Close
            </button>
          </div>
        </div>
      </div>
    </div>
    <div v-if="viewingVersions" class="modal-backdrop fade show"></div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted, watch, computed } from 'vue';
import { useSecretStore } from '@/stores/secretStore';
import type { Secret, SecretVersion, CreateSecretRequest, UpdateSecretRequest } from '@/types';
import { formatDate } from '@/utils/dateUtils';
import SecretForm from '@/components/secrets/SecretForm.vue';
import SecretList from '@/components/secrets/SecretList.vue';

const secretStore = useSecretStore();

// UI State
const showForm = ref(false);
const editingSecret = ref<Secret | null>(null);
const viewingSecret = ref<Secret | null>(null);
const viewingVersions = ref<Secret | null>(null);
const secretVersions = ref<SecretVersion[]>([]);
const versionsLoading = ref(false);

// KEK Modal State
const showKEKModal = ref(false);
const kekInput = ref('');
const kekError = ref('');
const kekLoading = ref(false);
const selectedSecretForKEK = ref<Secret | null>(null);

// Search and Filter State
const searchQuery = ref('');
const filters = ref({
  group: '',
  realm: ''
});
const currentPage = ref(1);
const pageSize = ref(20);

// Computed properties
const isKEKRequired = computed(() => {
  // For versioned secrets, we need to check the current version's KEK version
  // This will be handled in the KEK modal when we have the version info
  return false; // Default to false, will be determined by version
});

// Helper functions
const getStatusColor = (status: string) => {
  switch (status) {
    case 'active': return 'success';
    case 'pending': return 'warning';
    case 'deprecated': return 'secondary';
    default: return 'light';
  }
};

// Debounced search
let searchTimeout: ReturnType<typeof setTimeout>;
watch([searchQuery, filters], () => {
  clearTimeout(searchTimeout);
  searchTimeout = setTimeout(() => {
    currentPage.value = 1;
    loadSecrets();
  }, 300);
});

watch(currentPage, () => {
  loadSecrets();
});

// KEK Validation - KEK is optional, can be empty
const isValidKEK = computed(() => {
  return true; // Always valid since KEK is optional
});

// Load secrets
const loadSecrets = async () => {
  try {
    await secretStore.fetchSecrets({
      q: searchQuery.value,
      group: filters.value.group,
      realm_id: filters.value.realm,
      page: currentPage.value,
      page_size: pageSize.value
    });
  } catch (error) {
    console.error('Failed to load secrets:', error);
  }
};

// Form handlers
const showCreateForm = () => {
  editingSecret.value = null;
  showForm.value = true;
};

const closeForm = () => {
  showForm.value = false;
  editingSecret.value = null;
};

const handleFormSubmit = async (secretData: CreateSecretRequest | UpdateSecretRequest) => {
  console.log('handleFormSubmit called with:', secretData);
  console.log('editingSecret:', editingSecret.value);
  
  try {
    if (editingSecret.value) {
      console.log('Updating secret with ID:', editingSecret.value.id);
      // For editing, use UpdateSecretRequest
      await secretStore.updateSecret(editingSecret.value.id, secretData as UpdateSecretRequest);
      console.log('Secret updated successfully');
    } else {
      console.log('Creating new secret');
      // For creating, use CreateSecretRequest
      await secretStore.createSecret(secretData as CreateSecretRequest);
      console.log('Secret created successfully');
    }
    closeForm();
    await loadSecrets();
  } catch (error) {
    console.error('Failed to save secret:', error);
    alert('Failed to save secret. Please try again.');
  }
};

// List action handlers
const handleView = (secret: Secret) => {
  viewingSecret.value = secret;
};

const handleEdit = (secret: Secret) => {
  editingSecret.value = secret;
  showForm.value = true;
};

const handleDelete = async (id: string) => {
  try {
    await secretStore.deleteSecret(id);
    await loadSecrets();
  } catch (error) {
    console.error('Failed to delete secret:', error);
    alert('Failed to delete secret. Please try again.');
  }
};

const handleCopy = async (secret: Secret) => {
  // For versioned secrets, we always try the current version first
  try {
    await secretStore.copySecretValue(secret);
    if (viewingSecret.value) {
      viewingSecret.value = null;
    }
  } catch (error) {
    console.error('Failed to copy secret:', error);
    // If default KEK fails, offer the KEK modal as fallback
    openKEKModal(secret);
  }
};

const handleViewVersions = async (secret: Secret) => {
  viewingVersions.value = secret;
  versionsLoading.value = true;
  
  try {
    secretVersions.value = await secretStore.getSecretVersions(secret.id);
  } catch (error) {
    console.error('Failed to load secret versions:', error);
    alert('Failed to load secret versions');
  } finally {
    versionsLoading.value = false;
  }
};

// KEK Modal handlers
const openKEKModal = (secret: Secret) => {
  selectedSecretForKEK.value = secret;
  kekInput.value = '';
  kekError.value = '';
  showKEKModal.value = true;
};

const closeKEKModal = () => {
  showKEKModal.value = false;
  selectedSecretForKEK.value = null;
  kekInput.value = '';
  kekError.value = '';
};

const handleKEKDecrypt = async () => {
  if (!selectedSecretForKEK.value) return;
  
  // Validate KEK input for version 999 (custom KEK required)
  if (isKEKRequired.value && kekInput.value.trim().length === 0) {
    kekError.value = 'Custom KEK is required for this secret (KEK version 999)';
    return;
  }
  
  kekLoading.value = true;
  kekError.value = '';
  
  try {
    if (kekInput.value.trim().length > 0) {
      // User provided a custom KEK
      await secretStore.copySecretValueWithKEK(selectedSecretForKEK.value, kekInput.value);
    } else {
      // No custom KEK provided, use default system KEK (only allowed for non-999 versions)
      await secretStore.copySecretValue(selectedSecretForKEK.value);
    }
    closeKEKModal();
  } catch (error) {
    if (isKEKRequired.value) {
      kekError.value = 'Invalid custom KEK. Please check that you entered the exact password used during encryption.';
    } else if (kekInput.value.trim().length > 0) {
      kekError.value = 'Invalid KEK or decryption failed. Please check your KEK.';
    } else {
      kekError.value = 'Decryption failed. This secret may have been encrypted with a custom KEK.';
    }
  } finally {
    kekLoading.value = false;
  }
};

// Version management handlers
const handleCopyVersion = async (version: SecretVersion) => {
  if (!viewingVersions.value) return;
  
  try {
    await secretStore.copySecretVersionValue(viewingVersions.value.id, version.version);
  } catch (error) {
    console.error('Failed to copy version:', error);
    // If default KEK fails, offer the KEK modal as fallback
    openKEKModalForVersion(version);
  }
};

const handleCopyVersionWithKEK = (version: SecretVersion) => {
  openKEKModalForVersion(version);
};

const openKEKModalForVersion = (version: SecretVersion) => {
  // Create a temporary secret object for the KEK modal
  const tempSecret = {
    ...viewingVersions.value!,
    kek_version: version.kek_version
  };
  openKEKModal(tempSecret);
};

const handleActivateVersion = async (version: number) => {
  if (!viewingVersions.value) return;
  
  try {
    await secretStore.activateSecretVersion(viewingVersions.value.id, version);
    // Refresh versions list
    await handleViewVersions(viewingVersions.value);
    // Refresh main secrets list
    await loadSecrets();
  } catch (error) {
    console.error('Failed to activate version:', error);
    alert('Failed to activate version');
  }
};

const handleDeleteVersion = async (version: number) => {
  if (!viewingVersions.value) return;
  
  if (!confirm(`Are you sure you want to delete version ${version}? This action cannot be undone.`)) {
    return;
  }
  
  try {
    await secretStore.deleteSecretVersion(viewingVersions.value.id, version);
    // Refresh versions list
    await handleViewVersions(viewingVersions.value);
    // Refresh main secrets list
    await loadSecrets();
  } catch (error) {
    console.error('Failed to delete version:', error);
    alert('Failed to delete version');
  }
};

const showCreatePendingForm = () => {
  // TODO: Implement pending version creation form
  alert('Create pending version form - to be implemented');
};

// Initialize
onMounted(() => {
  loadSecrets();
});
</script>

<style scoped>
/* Page Layout */
.secrets-view {
  padding: 2rem;
  max-width: 1400px;
  margin: 0 auto;
}

/* Page Header */
.page-header {
  display: flex;
  justify-content: space-between;
  align-items: start;
  margin-bottom: 2rem;
  padding-bottom: 2rem;
  border-bottom: 2px solid #e9ecef;
}

.header-content {
  flex: 1;
}

.page-title {
  font-size: 2rem;
  font-weight: 600;
  color: #212529;
  margin-bottom: 0.5rem;
  display: flex;
  align-items: center;
  gap: 0.75rem;
}

.page-title i {
  color: #667eea;
}

.page-description {
  color: #6c757d;
  margin: 0;
}

.header-actions {
  display: flex;
  gap: 0.5rem;
}

/* Content Wrapper */
.content-wrapper {
  position: relative;
}

/* Form Section */
.form-section {
  background: white;
  border-radius: 16px;
  padding: 2rem;
  margin-bottom: 2rem;
  box-shadow: 0 4px 20px rgba(0, 0, 0, 0.08);
}

.section-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 2rem;
  padding-bottom: 1rem;
  border-bottom: 1px solid #e9ecef;
}

.section-header h2 {
  font-size: 1.5rem;
  font-weight: 600;
  margin: 0;
}

/* Loading Overlay */
.loading-overlay {
  position: fixed;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  background: rgba(255, 255, 255, 0.9);
  display: flex;
  align-items: center;
  justify-content: center;
  z-index: 1000;
}

/* Modal Styles */
.modal {
  background: rgba(0, 0, 0, 0.5);
}

.modal-dialog {
  margin-top: 5rem;
}

.secret-details {
  display: flex;
  flex-direction: column;
  gap: 1rem;
}

.detail-row {
  display: flex;
  align-items: center;
  gap: 1rem;
  padding: 0.5rem 0;
  border-bottom: 1px solid #f8f9fa;
}

.detail-row:last-child {
  border-bottom: none;
}

.detail-label {
  font-weight: 600;
  color: #6c757d;
  min-width: 120px;
}

.detail-value {
  color: #212529;
}

/* Responsive Design */
@media (max-width: 768px) {
  .secrets-view {
    padding: 1rem;
  }
  
  .page-header {
    flex-direction: column;
    gap: 1rem;
  }
  
  .header-actions {
    width: 100%;
  }
  
  .header-actions .btn {
    flex: 1;
  }
  
  .form-section {
    padding: 1rem;
  }
  
  .detail-row {
    flex-direction: column;
    align-items: start;
    gap: 0.25rem;
  }
}
</style>`
