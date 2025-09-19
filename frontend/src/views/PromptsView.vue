<template>
  <div class="prompts-view">
    <!-- Page Header -->
    <div class="page-header">
      <div class="header-content">
        <h1 class="page-title">
          <i class="bi bi-chat-dots"></i>
          Prompt Management
        </h1>
        <p class="page-description">
          Create and manage AI prompts with system instructions and user templates for consistent AI interactions
        </p>
      </div>
      <div class="header-actions">
        <button class="btn btn-primary" @click="showCreateForm">
          <i class="bi bi-plus-lg me-2"></i>
          Create Prompt
        </button>
      </div>
    </div>

    <!-- Main Content -->
    <div class="content-wrapper">
      <!-- Create/Edit Form -->
      <div v-if="showForm" class="form-section">
        <div class="section-header">
          <h2>{{ editingPrompt ? 'Edit Prompt' : 'Create New Prompt' }}</h2>
          <button 
            class="btn btn-sm btn-light"
            @click="closeForm"
          >
            <i class="bi bi-x-lg"></i>
          </button>
        </div>
        <PromptForm
          :prompt="editingPrompt"
          :submit-button-text="editingPrompt ? 'Update Prompt' : 'Create Prompt'"
          @submit="handleFormSubmit"
          @cancel="closeForm"
        />
      </div>

      <!-- Prompt List -->
      <div v-else>
        <PromptList
          :prompts="promptStore.prompts"
          :search-query="searchQuery"
          :current-page="currentPage"
          :page-size="pageSize"
          :total-count="promptStore.totalCount"
          :loading="promptStore.loading"
          @view="handleView"
          @edit="handleEdit"
          @delete="handleDelete"
          @update:search-query="searchQuery = $event"
          @update:page="currentPage = $event"
        />
      </div>
    </div>

    <!-- Loading Overlay -->
    <div v-if="promptStore.loading" class="loading-overlay">
      <div class="spinner-border text-primary" role="status">
        <span class="visually-hidden">Loading...</span>
      </div>
    </div>

    <!-- Prompt Details Modal -->
    <div 
      v-if="viewingPrompt"
      class="modal fade show d-block"
      tabindex="-1"
      @click.self="viewingPrompt = null"
    >
      <div class="modal-dialog modal-xl">
        <div class="modal-content">
          <div class="modal-header">
            <h5 class="modal-title">
              <i class="bi bi-chat-dots me-2"></i>
              Prompt Details
            </h5>
            <button 
              type="button" 
              class="btn-close" 
              @click="viewingPrompt = null"
            ></button>
          </div>
          <div class="modal-body">
            <div class="prompt-details">
              <div class="detail-row">
                <span class="detail-label">Name:</span>
                <span class="detail-value fw-bold">{{ viewingPrompt.name }}</span>
              </div>
              <div class="detail-row" v-if="viewingPrompt.description">
                <span class="detail-label">Description:</span>
                <span class="detail-value">{{ viewingPrompt.description }}</span>
              </div>
              <div class="detail-row" v-if="viewingPrompt.tags">
                <span class="detail-label">Tags:</span>
                <span class="detail-value">
                  <span 
                    v-for="tag in getTagsArray(viewingPrompt.tags)" 
                    :key="tag" 
                    class="badge bg-light text-dark me-1"
                  >
                    {{ tag }}
                  </span>
                </span>
              </div>
              <div class="detail-row">
                <span class="detail-label">Created By:</span>
                <span class="detail-value">{{ viewingPrompt.created_by }}</span>
              </div>
              <div class="detail-row">
                <span class="detail-label">Created:</span>
                <span class="detail-value">{{ formatDate(viewingPrompt.created_at) }}</span>
              </div>
              <div class="detail-row" v-if="viewingPrompt.updated_at !== viewingPrompt.created_at">
                <span class="detail-label">Updated:</span>
                <span class="detail-value">{{ formatDate(viewingPrompt.updated_at) }}</span>
              </div>
            </div>
            
            <div class="prompt-content">
              <div class="content-section">
                <h6 class="content-title">
                  <i class="bi bi-gear"></i> System Prompt
                </h6>
                <div class="content-display">
                  {{ viewingPrompt.system_prompt }}
                </div>
              </div>
              
              <div class="content-section">
                <h6 class="content-title">
                  <i class="bi bi-person"></i> User Prompt Template
                </h6>
                <div class="content-display">
                  {{ viewingPrompt.user_prompt }}
                </div>
              </div>
            </div>
          </div>
          <div class="modal-footer">
            <button 
              type="button" 
              class="btn btn-secondary" 
              @click="viewingPrompt = null"
            >
              Close
            </button>
            <button 
              type="button" 
              class="btn btn-primary" 
              @click="handleEdit(viewingPrompt)"
            >
              <i class="bi bi-pencil me-2"></i>
              Edit Prompt
            </button>
          </div>
        </div>
      </div>
    </div>
    <div v-if="viewingPrompt" class="modal-backdrop fade show"></div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted, watch } from 'vue';
import { usePromptStore } from '@/stores/promptStore';
import type { Prompt } from '@/types';
import { formatDate } from '@/utils/dateUtils';
import PromptForm from '@/components/prompts/PromptForm.vue';
import PromptList from '@/components/prompts/PromptList.vue';

const promptStore = usePromptStore();

// UI State
const showForm = ref(false);
const editingPrompt = ref<Prompt | undefined>(undefined);
const viewingPrompt = ref<Prompt | null>(null);

// Search and Pagination State
const searchQuery = ref('');
const currentPage = ref(1);
const pageSize = ref(20);

// Debounced search
let searchTimeout: ReturnType<typeof setTimeout>;
watch([searchQuery], () => {
  clearTimeout(searchTimeout);
  searchTimeout = setTimeout(() => {
    currentPage.value = 1;
    loadPrompts();
  }, 300);
});

watch(currentPage, () => {
  loadPrompts();
});

// Load prompts
const loadPrompts = async () => {
  try {
    if (searchQuery.value.trim()) {
      await promptStore.searchPrompts(searchQuery.value, {
        page: currentPage.value,
        page_size: pageSize.value
      });
    } else {
      await promptStore.fetchPrompts({
        page: currentPage.value,
        page_size: pageSize.value
      });
    }
  } catch (error) {
    console.error('Failed to load prompts:', error);
  }
};

// Form handlers
const showCreateForm = () => {
  editingPrompt.value = undefined;
  showForm.value = true;
};

const closeForm = () => {
  showForm.value = false;
  editingPrompt.value = undefined;
};

const handleFormSubmit = async (prompt: Prompt) => {
  try {
    if (editingPrompt.value) {
      await promptStore.updatePrompt(editingPrompt.value.id, prompt);
    } else {
      await promptStore.addPrompt(prompt);
    }
    closeForm();
    await loadPrompts();
  } catch (error) {
    console.error('Failed to save prompt:', error);
    alert('Failed to save prompt. Please try again.');
  }
};

// List action handlers
const handleView = (prompt: Prompt) => {
  viewingPrompt.value = prompt;
};

const handleEdit = (prompt: Prompt) => {
  editingPrompt.value = prompt;
  showForm.value = true;
  if (viewingPrompt.value) {
    viewingPrompt.value = null;
  }
};

const handleDelete = async (id: string) => {
  if (confirm('Are you sure you want to delete this prompt?')) {
    try {
      await promptStore.deletePrompt(id);
      await loadPrompts();
    } catch (error) {
      console.error('Failed to delete prompt:', error);
      alert('Failed to delete prompt. Please try again.');
    }
  }
};

// Helper functions
const getTagsArray = (tags: string): string[] => {
  if (!tags) return [];
  return tags.split(',').map(tag => tag.trim()).filter(tag => tag.length > 0);
};

// Initialize
onMounted(() => {
  loadPrompts();
});
</script>

<style scoped>
/* Page Layout */
.prompts-view {
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
  margin-top: 2rem;
}

.prompt-details {
  display: flex;
  flex-direction: column;
  gap: 1rem;
  margin-bottom: 2rem;
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

.prompt-content {
  background: #f8f9fa;
  border-radius: 12px;
  padding: 1.5rem;
  border: 1px solid #e9ecef;
}

.content-section {
  margin-bottom: 1.5rem;
}

.content-section:last-child {
  margin-bottom: 0;
}

.content-title {
  font-size: 1rem;
  font-weight: 600;
  color: #495057;
  margin-bottom: 0.75rem;
  display: flex;
  align-items: center;
  gap: 0.5rem;
}

.content-display {
  background: white;
  border: 1px solid #e9ecef;
  border-radius: 8px;
  padding: 1rem;
  font-family: 'Monaco', 'Menlo', 'Ubuntu Mono', monospace;
  font-size: 0.9rem;
  line-height: 1.5;
  color: #495057;
  white-space: pre-wrap;
  word-break: break-word;
  max-height: 300px;
  overflow-y: auto;
}

/* Badge Styles */
.badge {
  font-size: 0.75rem;
  padding: 0.25rem 0.5rem;
  font-weight: 500;
}

/* Responsive Design */
@media (max-width: 768px) {
  .prompts-view {
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
  
  .prompt-content {
    padding: 1rem;
  }
  
  .content-display {
    padding: 0.75rem;
    font-size: 0.8rem;
  }
}
</style>
