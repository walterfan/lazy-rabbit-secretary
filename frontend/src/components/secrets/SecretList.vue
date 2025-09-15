`<template>
  <div class="secret-list">
    <!-- List Header -->
    <div class="list-header">
      <div class="search-filters">
        <div class="search-wrapper">
          <SearchInput
            :model-value="searchQuery"
            @update:model-value="$emit('update:searchQuery', $event)"
            placeholder="Search secrets by name, group, or path..."
            class="search-input"
          />
        </div>
        <div class="filter-group">
          <select
            class="form-select form-select-sm"
            v-model="filters.group"
            @change="$emit('update:filters', filters)"
          >
            <option value="">All Groups</option>
            <option v-for="group in availableGroups" :key="group" :value="group">
              {{ group }}
            </option>
          </select>
          <select
            class="form-select form-select-sm"
            v-model="filters.realm"
            @change="$emit('update:filters', filters)"
          >
            <option value="">All Realms</option>
            <option v-for="realm in availableRealms" :key="realm" :value="realm">
              {{ realm }}
            </option>
          </select>
        </div>
      </div>
      <div class="list-stats">
        <span class="text-muted">
          <i class="bi bi-shield-lock"></i>
          {{ totalCount }} {{ totalCount === 1 ? 'secret' : 'secrets' }}
        </span>
      </div>
    </div>

    <!-- Secrets Table -->
    <div class="table-responsive">
      <table class="table table-hover">
        <thead>
          <tr>
            <th class="col-name">
              <i class="bi bi-key text-muted me-2"></i>Name
            </th>
            <th class="col-group">
              <i class="bi bi-folder text-muted me-2"></i>Group
            </th>
            <th class="col-path">
              <i class="bi bi-signpost text-muted me-2"></i>Path
            </th>
            <th class="col-desc">
              <i class="bi bi-text-paragraph text-muted me-2"></i>Description
            </th>
            <th class="col-meta">
              <i class="bi bi-info-circle text-muted me-2"></i>Info
            </th>
            <th class="col-actions">Actions</th>
          </tr>
        </thead>
        <tbody>
          <tr v-for="secret in secrets" :key="secret.id" class="secret-row">
            <td class="col-name">
              <div class="secret-name">
                <i class="bi bi-key-fill text-primary"></i>
                <span class="fw-medium">{{ secret.name }}</span>
              </div>
            </td>
            <td class="col-group">
              <span class="badge bg-secondary">{{ secret.group }}</span>
            </td>
            <td class="col-path">
              <code class="path-text">{{ secret.path }}</code>
            </td>
            <td class="col-desc">
              <span class="text-muted small">{{ secret.desc || '-' }}</span>
            </td>
            <td class="col-meta">
              <div class="meta-info">
                <div class="meta-item">
                  <i class="bi bi-shield-check text-success"></i>
                  <span class="small">{{ secret.cipher_alg }}</span>
                </div>
                <div class="meta-item">
                  <i class="bi bi-person text-muted"></i>
                  <span class="small">{{ secret.created_by }}</span>
                </div>
                <div class="meta-item">
                  <i class="bi bi-clock text-muted"></i>
                  <span class="small">{{ formatDate(secret.updated_at) }}</span>
                </div>
              </div>
            </td>
            <td class="col-actions">
              <div class="action-buttons">
                <button
                  class="btn btn-sm btn-outline-primary"
                  @click="handleCopyValue(secret)"
                  title="Copy secret value"
                >
                  <i class="bi bi-clipboard"></i>
                </button>
                <button
                  class="btn btn-sm btn-outline-secondary"
                  @click="$emit('view', secret)"
                  title="View details"
                >
                  <i class="bi bi-eye"></i>
                </button>
                <button
                  class="btn btn-sm btn-outline-warning"
                  @click="$emit('edit', secret)"
                  title="Edit secret"
                >
                  <i class="bi bi-pencil"></i>
                </button>
                <button
                  class="btn btn-sm btn-outline-danger"
                  @click="handleDelete(secret.id)"
                  title="Delete secret"
                >
                  <i class="bi bi-trash"></i>
                </button>
              </div>
            </td>
          </tr>
        </tbody>
      </table>
    </div>

    <!-- Empty State -->
    <div v-if="secrets.length === 0" class="empty-state">
      <i class="bi bi-shield-x"></i>
      <h4>No secrets found</h4>
      <p class="text-muted">
        {{ searchQuery || filters.group || filters.realm 
          ? 'Try adjusting your search or filters' 
          : 'Create your first secret to get started' }}
      </p>
    </div>

    <!-- Pagination -->
    <div v-if="totalPages > 1" class="pagination-wrapper">
      <nav>
        <ul class="pagination">
          <li class="page-item" :class="{ disabled: currentPage === 1 }">
            <a 
              class="page-link" 
              href="#" 
              @click.prevent="$emit('update:page', currentPage - 1)"
            >
              Previous
            </a>
          </li>
          <li 
            v-for="page in displayedPages" 
            :key="page"
            class="page-item" 
            :class="{ active: page === currentPage }"
          >
            <a 
              class="page-link" 
              href="#" 
              @click.prevent="$emit('update:page', page)"
            >
              {{ page }}
            </a>
          </li>
          <li class="page-item" :class="{ disabled: currentPage === totalPages }">
            <a 
              class="page-link" 
              href="#" 
              @click.prevent="$emit('update:page', currentPage + 1)"
            >
              Next
            </a>
          </li>
        </ul>
      </nav>
    </div>
  </div>
</template>

<script setup lang="ts">
import { computed, ref } from 'vue';
import type { Secret } from '@/types';
import { formatDate } from '@/utils/dateUtils';
import { useConfirmDialog } from '@/components/common/ConfirmDialog';
import SearchInput from '@/components/common/SearchInput.vue';

const props = defineProps<{
  secrets: Secret[];
  searchQuery: string;
  filters: {
    group: string;
    realm: string;
  };
  currentPage: number;
  pageSize: number;
  totalCount: number;
}>();

const emit = defineEmits<{
  (e: 'view', secret: Secret): void;
  (e: 'edit', secret: Secret): void;
  (e: 'delete', id: string): void;
  (e: 'copy', secret: Secret): void;
  (e: 'update:searchQuery', value: string): void;
  (e: 'update:filters', filters: any): void;
  (e: 'update:page', page: number): void;
}>();

const { confirm } = useConfirmDialog();

// Extract unique groups and realms from secrets
const availableGroups = computed(() => {
  const groups = new Set(props.secrets.map(s => s.group));
  return Array.from(groups).sort();
});

const availableRealms = computed(() => {
  // This would ideally come from a separate API call
  return ['default', 'production', 'staging', 'development'];
});

// Pagination calculations
const totalPages = computed(() => Math.ceil(props.totalCount / props.pageSize));

const displayedPages = computed(() => {
  const pages: number[] = [];
  const maxPages = 5;
  
  let start = Math.max(1, props.currentPage - Math.floor(maxPages / 2));
  let end = Math.min(totalPages.value, start + maxPages - 1);
  
  if (end - start + 1 < maxPages) {
    start = Math.max(1, end - maxPages + 1);
  }
  
  for (let i = start; i <= end; i++) {
    pages.push(i);
  }
  
  return pages;
});

const handleDelete = async (id: string) => {
  if (await confirm('Are you sure you want to delete this secret? This action cannot be undone.')) {
    emit('delete', id);
  }
};

const handleCopyValue = (secret: Secret) => {
  emit('copy', secret);
};
</script>

<style scoped>
/* List Header */
.list-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 2rem;
  gap: 1rem;
  flex-wrap: wrap;
}

.search-filters {
  display: flex;
  gap: 1rem;
  flex: 1;
  flex-wrap: wrap;
}

.search-wrapper {
  flex: 1;
  min-width: 300px;
}

.filter-group {
  display: flex;
  gap: 0.5rem;
}

.filter-group .form-select {
  width: auto;
}

.list-stats {
  font-size: 0.9rem;
}

/* Table Styles */
.table {
  background: white;
  border-radius: 12px;
  overflow: hidden;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.08);
}

.table thead {
  background: linear-gradient(135deg, #f8f9fa 0%, #e9ecef 100%);
}

.table th {
  border-bottom: 2px solid #dee2e6;
  font-weight: 600;
  color: #495057;
  padding: 1rem;
  white-space: nowrap;
}

.table td {
  padding: 1rem;
  vertical-align: middle;
}

.secret-row:hover {
  background-color: #f8f9fa;
}

/* Column Widths */
.col-name { width: 20%; }
.col-group { width: 15%; }
.col-path { width: 20%; }
.col-desc { width: 20%; }
.col-meta { width: 15%; }
.col-actions { width: 10%; }

/* Secret Name */
.secret-name {
  display: flex;
  align-items: center;
  gap: 0.5rem;
}

/* Path Text */
.path-text {
  background-color: #f8f9fa;
  padding: 0.25rem 0.5rem;
  border-radius: 4px;
  font-size: 0.875rem;
  color: #495057;
}

/* Meta Info */
.meta-info {
  display: flex;
  flex-direction: column;
  gap: 0.25rem;
}

.meta-item {
  display: flex;
  align-items: center;
  gap: 0.25rem;
  font-size: 0.75rem;
  color: #6c757d;
}

.meta-item i {
  font-size: 0.875rem;
}

/* Action Buttons */
.action-buttons {
  display: flex;
  gap: 0.25rem;
  justify-content: flex-end;
}

.action-buttons .btn {
  padding: 0.25rem 0.5rem;
  font-size: 0.875rem;
}

/* Empty State */
.empty-state {
  text-align: center;
  padding: 4rem 2rem;
  color: #6c757d;
}

.empty-state i {
  font-size: 4rem;
  margin-bottom: 1rem;
  display: block;
  opacity: 0.3;
}

.empty-state h4 {
  color: #495057;
  margin-bottom: 0.5rem;
}

/* Pagination */
.pagination-wrapper {
  display: flex;
  justify-content: center;
  margin-top: 2rem;
}

.pagination {
  margin: 0;
}

/* Responsive Design */
@media (max-width: 992px) {
  .table-responsive {
    border-radius: 12px;
    box-shadow: 0 2px 8px rgba(0, 0, 0, 0.08);
  }
  
  .search-filters {
    flex-direction: column;
  }
  
  .filter-group {
    width: 100%;
  }
  
  .filter-group .form-select {
    flex: 1;
  }
  
  .action-buttons {
    flex-wrap: wrap;
  }
}

@media (max-width: 768px) {
  .list-header {
    flex-direction: column;
    align-items: stretch;
  }
  
  .list-stats {
    text-align: center;
  }
  
  .table {
    font-size: 0.875rem;
  }
  
  .table th, .table td {
    padding: 0.5rem;
  }
  
  .col-desc, .col-meta {
    display: none;
  }
}
</style>`
