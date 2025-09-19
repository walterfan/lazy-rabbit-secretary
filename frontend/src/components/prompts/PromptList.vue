<template>
  <div class="prompt-list">
    <!-- Search and Filters -->
    <div class="search-section">
      <div class="search-bar">
        <div class="input-group">
          <span class="input-group-text">
            <i class="bi bi-search"></i>
          </span>
          <input
            type="text"
            class="form-control"
            :value="searchQuery"
            @input="handleSearchChange"
            placeholder="Search prompts by name, description, or content..."
          />
        </div>
      </div>
      
      <div class="filter-controls">
        <div class="filter-group">
          <label class="form-label">Tags:</label>
          <input
            type="text"
            class="form-control form-control-sm"
            :value="tagsFilter"
            @input="handleTagsFilterChange"
            placeholder="Filter by tags"
          />
        </div>
        
        <div class="filter-group">
          <label class="form-label">Sort By:</label>
          <select class="form-select form-select-sm" :value="sortBy" @change="handleSortChange">
            <option value="name">Name</option>
            <option value="created_at">Created Date</option>
            <option value="updated_at">Updated Date</option>
          </select>
        </div>
      </div>
    </div>

    <!-- Prompts Grid -->
    <div class="prompts-grid" v-if="!loading && prompts.length > 0">
      <div 
        v-for="prompt in prompts" 
        :key="prompt.id" 
        class="prompt-card"
      >
        <div class="prompt-header">
          <h5 class="prompt-title">{{ prompt.name }}</h5>
          <div class="prompt-actions">
            <button 
              class="btn btn-sm btn-outline-primary" 
              @click="$emit('view', prompt)"
              title="View Details"
            >
              <i class="bi bi-eye"></i>
            </button>
            <button 
              class="btn btn-sm btn-outline-secondary" 
              @click="$emit('edit', prompt)"
              title="Edit Prompt"
            >
              <i class="bi bi-pencil"></i>
            </button>
            <button 
              class="btn btn-sm btn-outline-danger" 
              @click="$emit('delete', prompt.id)"
              title="Delete Prompt"
            >
              <i class="bi bi-trash"></i>
            </button>
          </div>
        </div>
        
        <div class="prompt-details">
          <div class="prompt-description" v-if="prompt.description">
            <p>{{ prompt.description }}</p>
          </div>
          
          <div class="prompt-preview">
            <div class="preview-section">
              <h6 class="preview-title">
                <i class="bi bi-gear"></i> System Prompt
              </h6>
              <div class="preview-content">
                {{ truncateText(prompt.system_prompt, 150) }}
              </div>
            </div>
            
            <div class="preview-section">
              <h6 class="preview-title">
                <i class="bi bi-person"></i> User Template
              </h6>
              <div class="preview-content">
                {{ truncateText(prompt.user_prompt, 150) }}
              </div>
            </div>
          </div>
          
          <div class="prompt-tags" v-if="prompt.tags">
            <span 
              v-for="tag in getTagsArray(prompt.tags)" 
              :key="tag" 
              class="badge bg-light text-dark me-1"
            >
              {{ tag }}
            </span>
          </div>
          
          <div class="prompt-meta">
            <small class="text-muted">
              <i class="bi bi-person"></i> {{ prompt.created_by }}
            </small>
            <small class="text-muted">
              <i class="bi bi-calendar"></i> {{ formatDate(prompt.created_at) }}
            </small>
            <small class="text-muted" v-if="prompt.updated_at !== prompt.created_at">
              <i class="bi bi-arrow-clockwise"></i> Updated {{ formatDate(prompt.updated_at) }}
            </small>
          </div>
        </div>
      </div>
    </div>

    <!-- Empty State -->
    <div class="empty-state" v-if="!loading && prompts.length === 0">
      <div class="empty-icon">
        <i class="bi bi-chat-dots"></i>
      </div>
      <h4>No prompts found</h4>
      <p>Try adjusting your search criteria or create some prompts to get started.</p>
    </div>

    <!-- Loading State -->
    <div class="loading-state" v-if="loading">
      <div class="spinner-border text-primary" role="status">
        <span class="visually-hidden">Loading...</span>
      </div>
      <p>Loading prompts...</p>
    </div>

    <!-- Pagination -->
    <div class="pagination-section" v-if="totalCount > pageSize">
      <nav aria-label="Prompts pagination">
        <ul class="pagination justify-content-center">
          <li class="page-item" :class="{ disabled: currentPage <= 1 }">
            <button 
              class="page-link" 
              @click="changePage(currentPage - 1)"
              :disabled="currentPage <= 1"
            >
              <i class="bi bi-chevron-left"></i>
            </button>
          </li>
          
          <li 
            v-for="page in visiblePages" 
            :key="page" 
            class="page-item" 
            :class="{ active: page === currentPage }"
          >
            <button 
              class="page-link" 
              @click="changePage(Number(page))"
            >
              {{ page }}
            </button>
          </li>
          
          <li class="page-item" :class="{ disabled: currentPage >= totalPages }">
            <button 
              class="page-link" 
              @click="changePage(currentPage + 1)"
              :disabled="currentPage >= totalPages"
            >
              <i class="bi bi-chevron-right"></i>
            </button>
          </li>
        </ul>
      </nav>
      
      <div class="pagination-info">
        <small class="text-muted">
          Showing {{ startItem }} to {{ endItem }} of {{ totalCount }} prompts
        </small>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, computed } from 'vue';
import type { Prompt } from '@/types';
import { formatDate } from '@/utils/dateUtils';

const props = defineProps<{
  prompts: Prompt[];
  searchQuery: string;
  currentPage: number;
  pageSize: number;
  totalCount: number;
  loading: boolean;
}>();

const emit = defineEmits<{
  (e: 'view', prompt: Prompt): void;
  (e: 'edit', prompt: Prompt): void;
  (e: 'delete', id: string): void;
  (e: 'update:search-query', value: string): void;
  (e: 'update:page', value: number): void;
}>();

// Local filter state
const tagsFilter = ref('');
const sortBy = ref('name');

// Computed properties
const totalPages = computed(() => Math.ceil(props.totalCount / props.pageSize));
const startItem = computed(() => (props.currentPage - 1) * props.pageSize + 1);
const endItem = computed(() => Math.min(props.currentPage * props.pageSize, props.totalCount));

const visiblePages = computed(() => {
  const pages = [];
  const start = Math.max(1, props.currentPage - 2);
  const end = Math.min(totalPages.value, props.currentPage + 2);
  
  for (let i = start; i <= end; i++) {
    pages.push(i);
  }
  return pages;
});

// Helper functions
const truncateText = (text: string, maxLength: number): string => {
  if (text.length <= maxLength) return text;
  return text.substring(0, maxLength) + '...';
};

const getTagsArray = (tags: string): string[] => {
  if (!tags) return [];
  return tags.split(',').map(tag => tag.trim()).filter(tag => tag.length > 0);
};

// Event handlers
const handleSearchChange = (event: Event) => {
  const target = event.target as HTMLInputElement;
  emit('update:search-query', target.value);
};

const handleTagsFilterChange = (event: Event) => {
  const target = event.target as HTMLInputElement;
  tagsFilter.value = target.value;
  // For now, we'll just emit the search query change
  // In a real implementation, you'd emit separate filter events
  emit('update:search-query', props.searchQuery);
};

const handleSortChange = (event: Event) => {
  const target = event.target as HTMLSelectElement;
  sortBy.value = target.value;
  // For now, we'll just emit the search query change
  // In a real implementation, you'd emit separate sort events
  emit('update:search-query', props.searchQuery);
};

const changePage = (page: number) => {
  if (page >= 1 && page <= totalPages.value) {
    emit('update:page', page);
  }
};
</script>

<style scoped>
/* Prompt List Container */
.prompt-list {
  padding: 1rem 0;
}

/* Search Section */
.search-section {
  background: white;
  border-radius: 12px;
  padding: 1.5rem;
  margin-bottom: 2rem;
  box-shadow: 0 2px 10px rgba(0, 0, 0, 0.05);
}

.search-bar {
  margin-bottom: 1rem;
}

.search-bar .input-group-text {
  background-color: #f8f9fa;
  border-color: #e9ecef;
}

.filter-controls {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(200px, 1fr));
  gap: 1rem;
}

.filter-group {
  display: flex;
  flex-direction: column;
  gap: 0.25rem;
}

.filter-group .form-label {
  font-size: 0.875rem;
  font-weight: 500;
  color: #6c757d;
  margin: 0;
}

/* Prompts Grid */
.prompts-grid {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(400px, 1fr));
  gap: 1.5rem;
  margin-bottom: 2rem;
}

/* Prompt Card */
.prompt-card {
  background: white;
  border-radius: 12px;
  padding: 1.5rem;
  box-shadow: 0 2px 10px rgba(0, 0, 0, 0.05);
  border: 2px solid transparent;
  transition: all 0.3s ease;
  position: relative;
}

.prompt-card:hover {
  transform: translateY(-2px);
  box-shadow: 0 4px 20px rgba(0, 0, 0, 0.1);
  border-color: #667eea;
}

/* Prompt Header */
.prompt-header {
  display: flex;
  justify-content: space-between;
  align-items: flex-start;
  margin-bottom: 1rem;
}

.prompt-title {
  font-size: 1.1rem;
  font-weight: 600;
  color: #212529;
  margin: 0;
  flex: 1;
  margin-right: 1rem;
}

.prompt-actions {
  display: flex;
  gap: 0.5rem;
  flex-shrink: 0;
}

/* Prompt Details */
.prompt-details {
  margin-bottom: 1rem;
}

.prompt-description {
  margin-bottom: 1rem;
}

.prompt-description p {
  font-size: 0.9rem;
  color: #495057;
  margin: 0;
  line-height: 1.4;
}

.prompt-preview {
  margin-bottom: 1rem;
}

.preview-section {
  margin-bottom: 0.75rem;
}

.preview-section:last-child {
  margin-bottom: 0;
}

.preview-title {
  font-size: 0.8rem;
  font-weight: 600;
  color: #6c757d;
  margin-bottom: 0.5rem;
  display: flex;
  align-items: center;
  gap: 0.25rem;
}

.preview-content {
  background: #f8f9fa;
  border: 1px solid #e9ecef;
  border-radius: 6px;
  padding: 0.75rem;
  font-family: 'Monaco', 'Menlo', 'Ubuntu Mono', monospace;
  font-size: 0.8rem;
  line-height: 1.4;
  color: #495057;
  white-space: pre-wrap;
  word-break: break-word;
}

.prompt-tags {
  margin-bottom: 1rem;
}

.prompt-meta {
  display: flex;
  flex-direction: column;
  gap: 0.25rem;
  padding-top: 0.75rem;
  border-top: 1px solid #f8f9fa;
}

.prompt-meta small {
  display: flex;
  align-items: center;
  gap: 0.25rem;
}

/* Prompt Actions */
.prompt-actions .btn {
  padding: 0.375rem 0.75rem;
  border-radius: 6px;
  transition: all 0.2s ease;
}

.prompt-actions .btn:hover {
  transform: translateY(-1px);
}

/* Empty State */
.empty-state {
  text-align: center;
  padding: 4rem 2rem;
  color: #6c757d;
}

.empty-icon {
  font-size: 4rem;
  margin-bottom: 1rem;
  opacity: 0.5;
}

.empty-state h4 {
  color: #495057;
  margin-bottom: 0.5rem;
}

/* Loading State */
.loading-state {
  text-align: center;
  padding: 4rem 2rem;
  color: #6c757d;
}

.loading-state .spinner-border {
  margin-bottom: 1rem;
}

/* Pagination */
.pagination-section {
  margin-top: 2rem;
}

.pagination {
  margin-bottom: 1rem;
}

.page-link {
  border-radius: 6px;
  margin: 0 2px;
  border: 1px solid #e9ecef;
  color: #495057;
  transition: all 0.2s ease;
}

.page-link:hover {
  background-color: #e9ecef;
  border-color: #dee2e6;
}

.page-item.active .page-link {
  background-color: #667eea;
  border-color: #667eea;
}

.page-item.disabled .page-link {
  color: #6c757d;
  background-color: #f8f9fa;
  border-color: #e9ecef;
}

.pagination-info {
  text-align: center;
}

/* Badges */
.badge {
  font-size: 0.75rem;
  padding: 0.25rem 0.5rem;
  font-weight: 500;
}

/* Responsive Design */
@media (max-width: 768px) {
  .prompts-grid {
    grid-template-columns: 1fr;
  }
  
  .filter-controls {
    grid-template-columns: 1fr;
  }
  
  .prompt-header {
    flex-direction: column;
    align-items: flex-start;
    gap: 0.5rem;
  }
  
  .prompt-title {
    margin-right: 0;
  }
  
  .prompt-actions {
    justify-content: center;
    flex-wrap: wrap;
  }
  
  .prompt-actions .btn {
    flex: 1;
    min-width: 40px;
  }
}

@media (max-width: 576px) {
  .search-section {
    padding: 1rem;
  }
  
  .prompt-card {
    padding: 1rem;
  }
  
  .prompt-actions {
    gap: 0.25rem;
  }
  
  .prompt-actions .btn {
    padding: 0.25rem 0.5rem;
    font-size: 0.875rem;
  }
}
</style>
