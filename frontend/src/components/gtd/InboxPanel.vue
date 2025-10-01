<template>
  <div class="inbox-panel">
    <!-- Quick Add Section -->
    <div class="quick-add-section">
      <div class="quick-add-card">
        <div class="quick-add-header">
          <h4>
            <i class="bi bi-lightning"></i>
            {{ $t('gtd.quickCaptureDesc') }}
          </h4>
        </div>
        <form @submit.prevent="handleQuickAdd" class="quick-add-form">
          <div class="input-group">
            <input
              type="text"
              v-model="quickAddTitle"
              :placeholder="$t('gtd.taskTitle')"
              class="form-control"
              ref="quickAddInput"
              required
            />
            <button type="submit" class="btn btn-primary">
              <i class="bi bi-plus-lg"></i>
              {{ $t('gtd.add') }}
            </button>
          </div>
          <div class="quick-options" v-if="showQuickOptions">
            <div class="priority-selector">
              <label>{{ $t('gtd.priority') }}:</label>
              <div class="priority-buttons">
                <button
                  v-for="priority in priorities"
                  :key="priority.value"
                  type="button"
                  :class="['priority-btn', priority.value, { active: quickAddPriority === priority.value }]"
                  @click="quickAddPriority = priority.value"
                >
                  <i :class="priority.icon"></i>
                  {{ priority.label }}
                </button>
              </div>
            </div>
            <div class="context-selector">
              <label>{{ $t('gtd.context') }}:</label>
              <div class="context-buttons">
                <button
                  v-for="context in contexts"
                  :key="context.value"
                  type="button"
                  :class="['context-btn', { active: quickAddContext === context.value }]"
                  @click="quickAddContext = context.value"
                >
                  <i :class="context.icon"></i>
                  {{ context.label }}
                </button>
              </div>
            </div>
          </div>
          <button
            type="button"
            class="btn btn-link btn-sm"
            @click="showQuickOptions = !showQuickOptions"
          >
            {{ showQuickOptions ? $t('gtd.simplify') : $t('gtd.moreOptions') }}
          </button>
        </form>
      </div>
    </div>

    <!-- Filter and Search -->
    <div class="filter-section">
      <div class="search-bar">
        <div class="input-group">
          <span class="input-group-text">
            <i class="bi bi-search"></i>
          </span>
            <input
              type="text"
              class="form-control"
              v-model="searchQuery"
              :placeholder="$t('gtd.searchInbox')"
              @input="debouncedSearch"
            />
        </div>
      </div>
      
      <div class="filter-controls">
        <div class="filter-group">
          <label>{{ $t('gtd.status') }}:</label>
          <select v-model="statusFilter" @change="applyFilters" class="form-select">
            <option value="">{{ $t('gtd.all') }}</option>
            <option value="pending">{{ $t('gtd.pending') }}</option>
            <option value="processing">{{ $t('gtd.processing') }}</option>
            <option value="completed">{{ $t('gtd.completed') }}</option>
            <option value="archived">{{ $t('gtd.archived') }}</option>
          </select>
        </div>
        
        <div class="filter-group">
          <label>{{ $t('gtd.priority') }}:</label>
          <select v-model="priorityFilter" @change="applyFilters" class="form-select">
            <option value="">{{ $t('gtd.all') }}</option>
            <option value="urgent">{{ $t('gtd.urgent') }}</option>
            <option value="high">{{ $t('gtd.high') }}</option>
            <option value="normal">{{ $t('gtd.normal') }}</option>
            <option value="low">{{ $t('gtd.low') }}</option>
          </select>
        </div>
        
        <div class="filter-group">
          <label>{{ $t('gtd.context') }}:</label>
          <select v-model="contextFilter" @change="applyFilters" class="form-select">
            <option value="">{{ $t('gtd.all') }}</option>
            <option value="@home">@{{ $t('gtd.home') }}</option>
            <option value="@office">@{{ $t('gtd.office') }}</option>
            <option value="@phone">@{{ $t('gtd.phone') }}</option>
            <option value="@computer">@{{ $t('gtd.computer') }}</option>
          </select>
        </div>
      </div>
    </div>

    <!-- Items Grid -->
    <div class="items-grid" v-if="!loading && items.length > 0">
      <div
        v-for="item in items"
        :key="item.id"
        :class="['inbox-card', item.priority, item.status]"
      >
        <div class="card-header">
          <div class="priority-indicator">
            <i :class="getPriorityIcon(item.priority)"></i>
          </div>
          <div class="card-actions">
              <button
                class="btn btn-sm btn-outline-primary"
                @click="editItem(item)"
                :title="$t('gtd.edit')"
              >
              <i class="bi bi-pencil"></i>
            </button>
            <button
              class="btn btn-sm btn-outline-success"
              @click="processItem(item)"
              :title="$t('gtd.process')"
            >
              <i class="bi bi-arrow-right"></i>
            </button>
            <button
              class="btn btn-sm btn-outline-danger"
              @click="deleteItem(item.id)"
              :title="$t('gtd.delete')"
            >
              <i class="bi bi-trash"></i>
            </button>
          </div>
        </div>
        
        <div class="card-content">
          <h5 class="item-title">{{ item.title }}</h5>
          <p class="item-description" v-if="item.description">{{ item.description }}</p>
          
          <div class="item-meta">
            <div class="meta-tags" v-if="item.tags">
              <span
                v-for="tag in getTagsArray(item.tags)"
                :key="tag"
                class="tag"
              >
                {{ tag }}
              </span>
            </div>
            
            <div class="meta-context" v-if="item.context">
              <span class="context-badge">
                <i :class="getContextIcon(item.context)"></i>
                {{ item.context }}
              </span>
            </div>
          </div>
        </div>
        
        <div class="card-footer">
          <div class="item-status">
            <span :class="['status-badge', item.status]">
              {{ getStatusLabel(item.status) }}
            </span>
          </div>
          <div class="item-time">
            <small class="text-muted">
              {{ formatDate(item.created_at, locale) }}
            </small>
          </div>
        </div>
      </div>
    </div>

    <!-- Empty State -->
    <div class="empty-state" v-if="!loading && items.length === 0">
      <div class="empty-icon">
        <i class="bi bi-inbox"></i>
      </div>
      <h4>{{ $t('gtd.inboxEmpty') }}</h4>
      <p>{{ $t('gtd.startRecording') }}</p>
    </div>

    <!-- Loading State -->
    <div class="loading-state" v-if="loading">
      <div class="spinner-border text-primary" role="status">
        <span class="visually-hidden">Loading...</span>
      </div>
      <p>{{ $t('gtd.loading') }}</p>
    </div>

    <!-- Pagination -->
    <div class="pagination-section" v-if="totalCount > pageSize">
      <nav aria-label="Inbox pagination">
        <ul class="pagination justify-content-center">
          <li class="page-item" :class="{ disabled: currentPage <= 1 }">
            <button class="page-link" @click="changePage(currentPage - 1)" :disabled="currentPage <= 1">
              <i class="bi bi-chevron-left"></i>
            </button>
          </li>
          
          <li
            v-for="page in visiblePages"
            :key="page"
            class="page-item"
            :class="{ active: page === currentPage }"
          >
            <button class="page-link" @click="changePage(page)">
              {{ page }}
            </button>
          </li>
          
          <li class="page-item" :class="{ disabled: currentPage >= totalPages }">
            <button class="page-link" @click="changePage(currentPage + 1)" :disabled="currentPage >= totalPages">
              <i class="bi bi-chevron-right"></i>
            </button>
          </li>
        </ul>
      </nav>
    </div>

    <!-- Edit Modal -->
    <div v-if="editingItem" class="modal fade show d-block" tabindex="-1" @click.self="closeEditModal">
      <div class="modal-dialog">
        <div class="modal-content">
          <div class="modal-header">
            <h5 class="modal-title">{{ $t('gtd.editInboxItem') }}</h5>
            <button type="button" class="btn-close" @click="closeEditModal"></button>
          </div>
          <div class="modal-body">
            <form @submit.prevent="handleEditSubmit">
              <div class="mb-3">
                <label class="form-label">{{ $t('gtd.title') }}</label>
                <input type="text" class="form-control" v-model="editForm.title" required>
              </div>
              <div class="mb-3">
                <label class="form-label">{{ $t('gtd.description') }}</label>
                <textarea class="form-control" v-model="editForm.description" rows="3"></textarea>
              </div>
              <div class="mb-3">
                <label class="form-label">{{ $t('gtd.priority') }}</label>
                <select class="form-select" v-model="editForm.priority">
                  <option value="low">{{ $t('gtd.low') }}</option>
                  <option value="normal">{{ $t('gtd.normal') }}</option>
                  <option value="high">{{ $t('gtd.high') }}</option>
                  <option value="urgent">{{ $t('gtd.urgent') }}</option>
                </select>
              </div>
              <div class="mb-3">
                <label class="form-label">{{ $t('gtd.status') }}</label>
                <select class="form-select" v-model="editForm.status">
                  <option value="pending">{{ $t('gtd.pending') }}</option>
                  <option value="processing">{{ $t('gtd.processing') }}</option>
                  <option value="completed">{{ $t('gtd.completed') }}</option>
                  <option value="archived">{{ $t('gtd.archived') }}</option>
                </select>
              </div>
              <div class="mb-3">
                <label class="form-label">{{ $t('gtd.tags') }}</label>
                <input type="text" class="form-control" v-model="editForm.tags" :placeholder="$t('gtd.tagsPlaceholder')">
              </div>
              <div class="mb-3">
                <label class="form-label">{{ $t('gtd.context') }}</label>
                <input type="text" class="form-control" v-model="editForm.context" :placeholder="$t('gtd.contextPlaceholder')">
              </div>
            </form>
          </div>
          <div class="modal-footer">
            <button type="button" class="btn btn-secondary" @click="closeEditModal">{{ $t('gtd.cancel') }}</button>
            <button type="button" class="btn btn-primary" @click="handleEditSubmit">{{ $t('gtd.save') }}</button>
          </div>
        </div>
      </div>
    </div>
    <div v-if="editingItem" class="modal-backdrop fade show"></div>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted, watch, nextTick } from 'vue';
import { useI18n } from 'vue-i18n';
import { useInboxStore } from '@/stores/inboxStore';
import type { InboxItem } from '@/types/gtd';
import { formatDate } from '@/utils/dateUtils';

const { t, locale } = useI18n();
const emit = defineEmits<{
  (e: 'refresh-stats'): void;
}>();

const inboxStore = useInboxStore();

// UI State
const loading = ref(false);
const searchQuery = ref('');
const statusFilter = ref('');
const priorityFilter = ref('');
const contextFilter = ref('');
const showQuickOptions = ref(false);

// Quick Add State
const quickAddTitle = ref('');
const quickAddPriority = ref('normal');
const quickAddContext = ref('');

// Edit State
const editingItem = ref<InboxItem | null>(null);
const editForm = ref({
  title: '',
  description: '',
  priority: 'normal' as 'low' | 'normal' | 'high' | 'urgent',
  status: 'pending' as 'pending' | 'processing' | 'completed' | 'archived',
  tags: '',
  context: ''
});

// Computed
const items = computed(() => inboxStore.items);
const totalCount = computed(() => inboxStore.totalCount);
const currentPage = computed({
  get: () => inboxStore.currentPage,
  set: (value: number) => { inboxStore.currentPage = value; }
});
const pageSize = computed(() => inboxStore.pageSize);
const totalPages = computed(() => Math.ceil(totalCount.value / pageSize.value));

const visiblePages = computed(() => {
  const pages = [];
  const start = Math.max(1, currentPage.value - 2);
  const end = Math.min(totalPages.value, currentPage.value + 2);
  
  for (let i = start; i <= end; i++) {
    pages.push(i);
  }
  return pages;
});

// Constants
const priorities = computed(() => [
  { value: 'low', label: t('gtd.low'), icon: 'bi bi-circle' },
  { value: 'normal', label: t('gtd.normal'), icon: 'bi bi-circle-fill' },
  { value: 'high', label: t('gtd.high'), icon: 'bi bi-triangle' },
  { value: 'urgent', label: t('gtd.urgent'), icon: 'bi bi-exclamation-triangle' }
]);

const contexts = computed(() => [
  { value: '@home', label: t('gtd.home'), icon: 'bi bi-house' },
  { value: '@office', label: t('gtd.office'), icon: 'bi bi-building' },
  { value: '@phone', label: t('gtd.phone'), icon: 'bi bi-telephone' },
  { value: '@computer', label: t('gtd.computer'), icon: 'bi bi-laptop' }
]);

// Methods
const loadItems = async () => {
  try {
    loading.value = true;
    await inboxStore.fetchItems({
      page: currentPage.value,
      page_size: pageSize.value,
      status: statusFilter.value,
      priority: priorityFilter.value,
      context: contextFilter.value,
      q: searchQuery.value
    });
  } catch (error) {
    console.error('Failed to load items:', error);
  } finally {
    loading.value = false;
  }
};

const debouncedSearch = debounce(() => {
  currentPage.value = 1;
  loadItems();
}, 300);

const applyFilters = () => {
  currentPage.value = 1;
  loadItems();
};

const changePage = (page: number) => {
  if (page >= 1 && page <= totalPages.value) {
    currentPage.value = page;
    loadItems();
  }
};

const handleQuickAdd = async () => {
  if (!quickAddTitle.value.trim()) return;
  
  try {
    const newItem: InboxItem = {
      id: crypto.randomUUID(),
      title: quickAddTitle.value.trim(),
      description: '',
      priority: quickAddPriority.value as 'low' | 'normal' | 'high' | 'urgent',
      status: 'pending' as 'pending' | 'processing' | 'completed' | 'archived',
      tags: '',
      context: quickAddContext.value,
      created_by: '',
      created_at: new Date(),
      updated_at: new Date()
    };
    
    await inboxStore.addItem(newItem);
    
    // Reset form
    quickAddTitle.value = '';
    quickAddPriority.value = 'normal';
    quickAddContext.value = '';
    showQuickOptions.value = false;
    
    // Focus input
    await nextTick();
    (document.querySelector('.quick-add-form input') as HTMLInputElement)?.focus();
    
    emit('refresh-stats');
  } catch (error) {
    console.error('Failed to add item:', error);
  }
};

const editItem = (item: InboxItem) => {
  editingItem.value = item;
  editForm.value = {
    title: item.title,
    description: item.description,
    priority: item.priority,
    status: item.status,
    tags: item.tags,
    context: item.context
  };
};

const closeEditModal = () => {
  editingItem.value = null;
};

const handleEditSubmit = async () => {
  if (!editingItem.value) return;
  
  try {
    await inboxStore.updateItem(editingItem.value.id, editForm.value);
    closeEditModal();
    emit('refresh-stats');
  } catch (error) {
    console.error('Failed to update item:', error);
  }
};

const processItem = async (item: InboxItem) => {
  try {
    await inboxStore.updateStatus(item.id, 'processing');
    emit('refresh-stats');
  } catch (error) {
    console.error('Failed to process item:', error);
  }
};

const deleteItem = async (id: string) => {
  if (!confirm(t('gtd.confirmDelete'))) return;
  
  try {
    await inboxStore.deleteItem(id);
    emit('refresh-stats');
  } catch (error) {
    console.error('Failed to delete item:', error);
  }
};

// Helper functions
const getPriorityIcon = (priority: string) => {
  const priorityMap: Record<string, string> = {
    low: 'bi bi-circle',
    normal: 'bi bi-circle-fill',
    high: 'bi bi-triangle',
    urgent: 'bi bi-exclamation-triangle'
  };
  return priorityMap[priority] || 'bi bi-circle';
};

const getContextIcon = (context: string) => {
  const contextMap: Record<string, string> = {
    '@home': 'bi bi-house',
    '@office': 'bi bi-building',
    '@phone': 'bi bi-telephone',
    '@computer': 'bi bi-laptop'
  };
  return contextMap[context] || 'bi bi-geo-alt';
};

const getStatusLabel = (status: string) => {
  const statusMap: Record<string, string> = {
    pending: t('gtd.pending'),
    processing: t('gtd.processing'),
    completed: t('gtd.completed'),
    archived: t('gtd.archived')
  };
  return statusMap[status] || status;
};

const getTagsArray = (tags: string): string[] => {
  if (!tags) return [];
  return tags.split(',').map(tag => tag.trim()).filter(tag => tag.length > 0);
};

// Debounce utility
function debounce(func: Function, wait: number) {
  let timeout: ReturnType<typeof setTimeout>;
  return function executedFunction(...args: any[]) {
    const later = () => {
      clearTimeout(timeout);
      func(...args);
    };
    clearTimeout(timeout);
    timeout = setTimeout(later, wait);
  };
}

// Lifecycle
onMounted(() => {
  loadItems();
});

// Watch for filter changes
watch([statusFilter, priorityFilter, contextFilter], () => {
  applyFilters();
});
</script>

<style scoped>
/* Inbox Panel Styles */
.inbox-panel {
  padding: 0;
}

/* Quick Add Section */
.quick-add-section {
  margin-bottom: 2rem;
}

.quick-add-card {
  background: linear-gradient(135deg, #667eea, #764ba2);
  border-radius: 16px;
  padding: 2rem;
  color: white;
  box-shadow: 0 8px 25px rgba(102, 126, 234, 0.3);
}

.quick-add-header h3 {
  font-size: 1.5rem;
  font-weight: 600;
  margin-bottom: 0.5rem;
  display: flex;
  align-items: center;
  gap: 0.5rem;
}

.quick-add-header p {
  opacity: 0.9;
  margin-bottom: 1.5rem;
}

.quick-add-form .input-group {
  margin-bottom: 1rem;
}

.quick-add-form .form-control {
  border: none;
  border-radius: 8px;
  padding: 0.75rem 1rem;
  font-size: 1rem;
}

.quick-add-form .btn-primary {
  background: rgba(255, 255, 255, 0.2);
  border: 1px solid rgba(255, 255, 255, 0.3);
  border-radius: 8px;
  padding: 0.75rem 1.5rem;
  font-weight: 500;
}

.quick-add-form .btn-primary:hover {
  background: rgba(255, 255, 255, 0.3);
}

.quick-options {
  background: rgba(255, 255, 255, 0.1);
  border-radius: 12px;
  padding: 1rem;
  margin-bottom: 1rem;
}

.priority-selector,
.context-selector {
  margin-bottom: 1rem;
}

.priority-selector label,
.context-selector label {
  display: block;
  margin-bottom: 0.5rem;
  font-weight: 500;
}

.priority-buttons,
.context-buttons {
  display: flex;
  gap: 0.5rem;
  flex-wrap: wrap;
}

.priority-btn,
.context-btn {
  border: 1px solid rgba(255, 255, 255, 0.3);
  background: rgba(255, 255, 255, 0.1);
  color: white;
  border-radius: 6px;
  padding: 0.5rem 1rem;
  font-size: 0.875rem;
  transition: all 0.2s ease;
}

.priority-btn:hover,
.context-btn:hover {
  background: rgba(255, 255, 255, 0.2);
}

.priority-btn.active,
.context-btn.active {
  background: rgba(255, 255, 255, 0.3);
  border-color: rgba(255, 255, 255, 0.5);
}

.priority-btn.low.active { background: #48bb78; }
.priority-btn.normal.active { background: #4299e1; }
.priority-btn.high.active { background: #ed8936; }
.priority-btn.urgent.active { background: #f56565; }

/* Filter Section */
.filter-section {
  background: #f8f9fa;
  border-radius: 12px;
  padding: 1.5rem;
  margin-bottom: 2rem;
}

.search-bar {
  margin-bottom: 1rem;
}

.filter-controls {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(150px, 1fr));
  gap: 1rem;
}

.filter-group label {
  font-size: 0.875rem;
  font-weight: 500;
  color: #495057;
  margin-bottom: 0.25rem;
  display: block;
}

/* Items Grid */
.items-grid {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(350px, 1fr));
  gap: 1.5rem;
  margin-bottom: 2rem;
}

.inbox-card {
  background: white;
  border-radius: 12px;
  padding: 1.5rem;
  box-shadow: 0 2px 10px rgba(0, 0, 0, 0.05);
  border: 2px solid transparent;
  transition: all 0.3s ease;
  position: relative;
  overflow: hidden;
}

.inbox-card:hover {
  transform: translateY(-2px);
  box-shadow: 0 4px 20px rgba(0, 0, 0, 0.1);
}

.inbox-card.low { border-left: 4px solid #48bb78; }
.inbox-card.normal { border-left: 4px solid #4299e1; }
.inbox-card.high { border-left: 4px solid #ed8936; }
.inbox-card.urgent { border-left: 4px solid #f56565; }

.card-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 1rem;
}

.priority-indicator {
  font-size: 1.2rem;
}

.priority-indicator .bi-circle { color: #48bb78; }
.priority-indicator .bi-circle-fill { color: #4299e1; }
.priority-indicator .bi-triangle { color: #ed8936; }
.priority-indicator .bi-exclamation-triangle { color: #f56565; }

.card-actions {
  display: flex;
  gap: 0.25rem;
}

.card-content {
  margin-bottom: 1rem;
}

.item-title {
  font-size: 1.1rem;
  font-weight: 600;
  color: #2d3748;
  margin-bottom: 0.5rem;
  line-height: 1.4;
}

.item-description {
  color: #718096;
  font-size: 0.9rem;
  line-height: 1.5;
  margin-bottom: 1rem;
}

.item-meta {
  display: flex;
  flex-direction: column;
  gap: 0.5rem;
}

.meta-tags {
  display: flex;
  flex-wrap: wrap;
  gap: 0.25rem;
}

.tag {
  background: #e2e8f0;
  color: #4a5568;
  font-size: 0.75rem;
  padding: 0.25rem 0.5rem;
  border-radius: 4px;
}

.context-badge {
  background: #667eea;
  color: white;
  font-size: 0.75rem;
  padding: 0.25rem 0.5rem;
  border-radius: 4px;
  display: inline-flex;
  align-items: center;
  gap: 0.25rem;
}

.card-footer {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding-top: 1rem;
  border-top: 1px solid #e2e8f0;
}

.status-badge {
  font-size: 0.75rem;
  padding: 0.25rem 0.5rem;
  border-radius: 4px;
  font-weight: 500;
}

.status-badge.pending { background: #fed7d7; color: #c53030; }
.status-badge.processing { background: #bee3f8; color: #2b6cb0; }
.status-badge.completed { background: #c6f6d5; color: #2f855a; }
.status-badge.archived { background: #e2e8f0; color: #4a5568; }

/* Empty State */
.empty-state {
  text-align: center;
  padding: 4rem 2rem;
  color: #718096;
}

.empty-icon {
  font-size: 4rem;
  margin-bottom: 1rem;
  opacity: 0.5;
}

.empty-state h4 {
  color: #2d3748;
  margin-bottom: 0.5rem;
}

/* Loading State */
.loading-state {
  text-align: center;
  padding: 4rem 2rem;
  color: #718096;
}

.loading-state .spinner-border {
  margin-bottom: 1rem;
}

/* Pagination */
.pagination-section {
  margin-top: 2rem;
}

.pagination {
  margin-bottom: 0;
}

.page-link {
  border-radius: 6px;
  margin: 0 2px;
  border: 1px solid #e2e8f0;
  color: #4a5568;
  transition: all 0.2s ease;
}

.page-link:hover {
  background-color: #f7fafc;
  border-color: #cbd5e0;
}

.page-item.active .page-link {
  background-color: #667eea;
  border-color: #667eea;
}

.page-item.disabled .page-link {
  color: #a0aec0;
  background-color: #f7fafc;
  border-color: #e2e8f0;
}

/* Modal Styles */
.modal {
  background: rgba(0, 0, 0, 0.5);
}

.modal-dialog {
  margin-top: 2rem;
}

/* Responsive Design */
@media (max-width: 768px) {
  .items-grid {
    grid-template-columns: 1fr;
  }
  
  .filter-controls {
    grid-template-columns: 1fr;
  }
  
  .priority-buttons,
  .context-buttons {
    justify-content: center;
  }
  
  .card-actions {
    flex-wrap: wrap;
  }
}

@media (max-width: 576px) {
  .quick-add-card {
    padding: 1.5rem;
  }
  
  .inbox-card {
    padding: 1rem;
  }
  
  .card-actions .btn {
    padding: 0.25rem 0.5rem;
    font-size: 0.875rem;
  }
}
</style>
