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
                <span class="detail-label">Encryption:</span>
                <span class="detail-value">
                  <i class="bi bi-shield-check text-success me-1"></i>
                  {{ viewingSecret.cipher_alg }}
                </span>
              </div>
              <div class="detail-row">
                <span class="detail-label">KEK Version:</span>
                <span class="detail-value">{{ viewingSecret.kek_version }}</span>
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
                  KEK Version: {{ selectedSecretForKEK?.kek_version }} (Custom KEK)
                </span>
                <span v-else>
                  <i class="bi bi-info-circle me-1"></i>
                  KEK Version: {{ selectedSecretForKEK?.kek_version }} (System Default). Only enter a custom KEK if this secret was specifically encrypted with one.
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
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted, watch, computed } from 'vue';
import { useSecretStore } from '@/stores/secretStore';
import type { Secret, CreateSecretRequest, UpdateSecretRequest } from '@/types';
import { formatDate } from '@/utils/dateUtils';
import SecretForm from '@/components/secrets/SecretForm.vue';
import SecretList from '@/components/secrets/SecretList.vue';

const secretStore = useSecretStore();

// UI State
const showForm = ref(false);
const editingSecret = ref<Secret | null>(null);
const viewingSecret = ref<Secret | null>(null);

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
  return selectedSecretForKEK.value?.kek_version === 999;
});

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
  // If kek_version is 999, it means the secret was encrypted with a custom KEK
  // and we must prompt the user for the custom KEK
  if (secret.kek_version === 999) {
    openKEKModal(secret);
    return;
  }
  
  // For other versions, try to decrypt with system default KEK
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
