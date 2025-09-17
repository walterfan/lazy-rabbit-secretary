<template>
  <div class="reminder-list">
    <!-- List Header -->
    <div class="list-header">
      <div class="search-filters">
        <div class="search-wrapper">
          <SearchInput
            :model-value="searchQuery"
            @update:model-value="$emit('update:searchQuery', $event)"
            placeholder="Search reminders by name, content, or tags..."
            class="search-input"
          />
        </div>
        <div class="filter-group">
          <select
            class="form-select form-select-sm"
            v-model="filters.status"
            @change="$emit('update:filters', filters)"
          >
            <option value="">All Status</option>
            <option value="pending">Pending</option>
            <option value="active">Active</option>
            <option value="completed">Completed</option>
            <option value="cancelled">Cancelled</option>
          </select>
          <select
            class="form-select form-select-sm"
            v-model="filters.tags"
            @change="$emit('update:filters', filters)"
          >
            <option value="">All Tags</option>
            <option v-for="tag in availableTags" :key="tag" :value="tag">
              {{ tag }}
            </option>
          </select>
        </div>
      </div>
      <div class="list-stats">
        <span class="text-muted">
          <i class="bi bi-bell"></i>
          {{ totalCount }} {{ totalCount === 1 ? 'reminder' : 'reminders' }}
        </span>
      </div>
    </div>

    <!-- Reminders Table -->
    <div class="table-responsive">
      <table class="table table-hover">
        <thead>
          <tr>
            <th class="col-name">
              <i class="bi bi-card-heading me-1"></i>
              Name
            </th>
            <th class="col-status">
              <i class="bi bi-flag me-1"></i>
              Status
            </th>
            <th class="col-remind-time">
              <i class="bi bi-clock me-1"></i>
              Remind Time
            </th>
            <th class="col-tags">
              <i class="bi bi-tags me-1"></i>
              Tags
            </th>
            <th class="col-notifications">
              <i class="bi bi-send me-1"></i>
              Notifications
            </th>
            <th class="col-actions">Actions</th>
          </tr>
        </thead>
        <tbody>
          <tr v-for="reminder in reminders" :key="reminder.id" class="reminder-row">
            <td class="col-name">
              <div class="reminder-name">
                <i class="bi bi-bell-fill text-primary"></i>
                <div class="name-content">
                  <span class="fw-medium">{{ reminder.name }}</span>
                  <small class="text-muted d-block content-preview">{{ reminder.content }}</small>
                </div>
              </div>
            </td>
            <td class="col-status">
              <ReminderStatusBadge :status="reminder.status" />
            </td>
            <td class="col-remind-time">
              <div class="time-info">
                <span class="time-display" :class="getTimeClass(reminder)">
                  {{ formatDateTime(reminder.remind_time) }}
                </span>
                <small class="text-muted d-block">
                  {{ getTimeRelative(reminder.remind_time) }}
                </small>
              </div>
            </td>
            <td class="col-tags">
              <ReminderTags :tags="reminder.tags" />
            </td>
            <td class="col-notifications">
              <div class="notification-info">
                <div class="notification-methods">
                  <span
                    v-for="method in getNotificationMethods(reminder.remind_methods || 'email')"
                    :key="method"
                    class="badge"
                    :class="getMethodBadgeClass(method)"
                  >
                    <i :class="getMethodIcon(method)" class="me-1"></i>
                    {{ method.toUpperCase() }}
                    <span v-if="method === 'email' && (!reminder.remind_methods || reminder.remind_methods === 'email')" class="ms-1">
                      (default)
                    </span>
                  </span>
                </div>
                <small v-if="reminder.remind_targets && (reminder.remind_methods || 'email')" class="text-muted d-block mt-1">
                  <i class="bi bi-bullseye"></i> {{ getTargetCount(reminder.remind_targets) }} target(s)
                </small>
              </div>
            </td>
            <td class="col-actions">
              <div class="action-buttons">
                <button
                  v-if="reminder.status === 'pending'"
                  class="btn btn-sm btn-success me-1"
                  @click="$emit('complete', reminder.id)"
                  title="Mark as completed"
                >
                  <i class="bi bi-check-lg"></i>
                </button>
                <button
                  v-if="reminder.status === 'pending'"
                  class="btn btn-sm btn-warning me-1"
                  @click="$emit('snooze', reminder)"
                  title="Snooze reminder"
                >
                  <i class="bi bi-clock"></i>
                </button>
                <button
                  class="btn btn-sm btn-outline-primary me-1"
                  @click="$emit('view', reminder)"
                  title="View details"
                >
                  <i class="bi bi-eye"></i>
                </button>
                <button
                  class="btn btn-sm btn-outline-secondary me-1"
                  @click="$emit('edit', reminder)"
                  title="Edit reminder"
                >
                  <i class="bi bi-pencil"></i>
                </button>
                <button
                  class="btn btn-sm btn-outline-danger"
                  @click="$emit('delete', reminder.id)"
                  title="Delete reminder"
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
    <div v-if="reminders.length === 0" class="empty-state">
      <div class="empty-icon">
        <i class="bi bi-bell-slash"></i>
      </div>
      <h5>No reminders found</h5>
      <p class="text-muted">
        {{ searchQuery || filters.status || filters.tags 
           ? 'Try adjusting your search or filters' 
           : 'Create your first reminder to get started' }}
      </p>
    </div>

    <!-- Pagination -->
    <div v-if="totalPages > 1" class="pagination-wrapper">
      <nav aria-label="Reminders pagination">
        <ul class="pagination pagination-sm justify-content-center">
          <li class="page-item" :class="{ disabled: currentPage === 1 }">
            <button
              class="page-link"
              @click="$emit('page-change', currentPage - 1)"
              :disabled="currentPage === 1"
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
              @click="$emit('page-change', page)"
            >
              {{ page }}
            </button>
          </li>
          
          <li class="page-item" :class="{ disabled: currentPage === totalPages }">
            <button
              class="page-link"
              @click="$emit('page-change', currentPage + 1)"
              :disabled="currentPage === totalPages"
            >
              <i class="bi bi-chevron-right"></i>
            </button>
          </li>
        </ul>
      </nav>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted } from 'vue';
import type { Reminder } from '@/types';
import { formatDate } from '@/utils/dateUtils';
import SearchInput from '@/components/common/SearchInput.vue';
import ReminderStatusBadge from './ReminderStatusBadge.vue';
import ReminderTags from './ReminderTags.vue';

// Props
interface Props {
  reminders: Reminder[];
  searchQuery: string;
  totalCount: number;
  currentPage: number;
  totalPages: number;
  loading?: boolean;
}

const props = withDefaults(defineProps<Props>(), {
  loading: false
});

// Emits
const emit = defineEmits<{
  (e: 'update:searchQuery', query: string): void;
  (e: 'update:filters', filters: any): void;
  (e: 'page-change', page: number): void;
  (e: 'view', reminder: Reminder): void;
  (e: 'edit', reminder: Reminder): void;
  (e: 'delete', id: string): void;
  (e: 'complete', id: string): void;
  (e: 'snooze', reminder: Reminder): void;
}>();

// Local state
const filters = ref({
  status: '',
  tags: ''
});

// Computed properties
const availableTags = computed(() => {
  const tagSet = new Set<string>();
  props.reminders.forEach(reminder => {
    if (reminder.tags) {
      reminder.tags.split(',').forEach(tag => {
        const trimmedTag = tag.trim();
        if (trimmedTag) {
          tagSet.add(trimmedTag);
        }
      });
    }
  });
  return Array.from(tagSet).sort();
});

const visiblePages = computed(() => {
  const pages = [];
  const maxVisible = 5;
  const halfVisible = Math.floor(maxVisible / 2);
  
  let startPage = Math.max(1, props.currentPage - halfVisible);
  let endPage = Math.min(props.totalPages, startPage + maxVisible - 1);
  
  // Adjust if we're near the end
  if (endPage - startPage < maxVisible - 1) {
    startPage = Math.max(1, endPage - maxVisible + 1);
  }
  
  for (let i = startPage; i <= endPage; i++) {
    pages.push(i);
  }
  
  return pages;
});

// Helper functions
const formatDateTime = (date: Date): string => {
  return new Date(date).toLocaleString('en-US', {
    year: 'numeric',
    month: 'short',
    day: 'numeric',
    hour: '2-digit',
    minute: '2-digit',
    hour12: true
  });
};

const getTimeRelative = (date: Date): string => {
  const now = new Date();
  const reminderDate = new Date(date);
  const diffMs = reminderDate.getTime() - now.getTime();
  const diffMinutes = Math.floor(diffMs / (1000 * 60));
  const diffHours = Math.floor(diffMs / (1000 * 60 * 60));
  const diffDays = Math.floor(diffMs / (1000 * 60 * 60 * 24));

  if (diffMs < 0) {
    const pastMinutes = Math.abs(diffMinutes);
    const pastHours = Math.abs(diffHours);
    const pastDays = Math.abs(diffDays);
    
    if (pastMinutes < 60) {
      return `${pastMinutes} min ago`;
    } else if (pastHours < 24) {
      return `${pastHours} hour${pastHours > 1 ? 's' : ''} ago`;
    } else {
      return `${pastDays} day${pastDays > 1 ? 's' : ''} ago`;
    }
  } else {
    if (diffMinutes < 60) {
      return `in ${diffMinutes} min`;
    } else if (diffHours < 24) {
      return `in ${diffHours} hour${diffHours > 1 ? 's' : ''}`;
    } else {
      return `in ${diffDays} day${diffDays > 1 ? 's' : ''}`;
    }
  }
};

const getTimeClass = (reminder: Reminder): string => {
  const now = new Date();
  const reminderDate = new Date(reminder.remind_time);
  const diffMs = reminderDate.getTime() - now.getTime();
  const diffHours = diffMs / (1000 * 60 * 60);

  if (reminder.status === 'completed') {
    return 'text-success';
  } else if (reminder.status === 'cancelled') {
    return 'text-muted';
  } else if (diffMs < 0) {
    return 'text-danger fw-bold'; // Overdue
  } else if (diffHours < 1) {
    return 'text-warning fw-bold'; // Due soon
  } else if (diffHours < 24) {
    return 'text-info'; // Due today
  } else {
    return 'text-dark'; // Future
  }
};

// Notification helper functions
const getNotificationMethods = (methods: string): string[] => {
  if (!methods) return [];
  return methods.split(',').map(m => m.trim()).filter(m => m);
};

const getMethodBadgeClass = (method: string): string => {
  const classes: Record<string, string> = {
    email: 'bg-primary',
    webhook: 'bg-success'
  };
  return classes[method.toLowerCase()] || 'bg-secondary';
};

const getMethodIcon = (method: string): string => {
  const icons: Record<string, string> = {
    email: 'bi-envelope',
    webhook: 'bi-link-45deg'
  };
  return icons[method.toLowerCase()] || 'bi-bell';
};

const getTargetCount = (targets: string): number => {
  if (!targets) return 0;
  // Simple comma-separated count, could be enhanced for JSON parsing
  return targets.split(',').filter(t => t.trim()).length;
};
</script>

<style scoped>
/* List Header */
.list-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 1.5rem;
  padding: 1rem;
  background: white;
  border-radius: 12px;
  border: 1px solid #e9ecef;
  box-shadow: 0 2px 4px rgba(0, 0, 0, 0.05);
}

.search-filters {
  display: flex;
  align-items: center;
  gap: 1rem;
  flex: 1;
}

.search-wrapper {
  flex: 1;
  max-width: 400px;
}

.filter-group {
  display: flex;
  gap: 0.75rem;
}

.filter-group .form-select {
  min-width: 120px;
}

.list-stats {
  display: flex;
  align-items: center;
  gap: 0.5rem;
  font-weight: 500;
}

/* Table Styles */
.table-responsive {
  background: white;
  border-radius: 12px;
  border: 1px solid #e9ecef;
  box-shadow: 0 2px 4px rgba(0, 0, 0, 0.05);
  overflow: hidden;
}

.table {
  margin-bottom: 0;
}

.table thead th {
  background: #f8f9fa;
  border-bottom: 2px solid #e9ecef;
  font-weight: 600;
  color: #495057;
  padding: 1rem 0.75rem;
}

.table tbody tr {
  transition: all 0.2s ease;
}

.table tbody tr:hover {
  background-color: #f8f9fa;
}

.table td {
  padding: 1rem 0.75rem;
  vertical-align: middle;
  border-bottom: 1px solid #f1f3f4;
}

/* Column Styles */
.col-name {
  width: 35%;
}

.col-status {
  width: 12%;
}

.col-remind-time {
  width: 20%;
}

.col-tags {
  width: 18%;
}

.col-actions {
  width: 15%;
}

.reminder-name {
  display: flex;
  align-items: flex-start;
  gap: 0.75rem;
}

.name-content {
  flex: 1;
  min-width: 0;
}

.content-preview {
  display: -webkit-box;
  -webkit-line-clamp: 2;
  line-clamp: 2;
  -webkit-box-orient: vertical;
  overflow: hidden;
  text-overflow: ellipsis;
  line-height: 1.3;
  margin-top: 0.25rem;
}

.time-info {
  display: flex;
  flex-direction: column;
  gap: 0.25rem;
}

.time-display {
  font-weight: 500;
}

/* Action Buttons */
.action-buttons {
  display: flex;
  align-items: center;
  gap: 0.25rem;
  flex-wrap: wrap;
}

.action-buttons .btn {
  border-radius: 6px;
  padding: 0.375rem 0.5rem;
  line-height: 1;
  transition: all 0.2s ease;
}

.action-buttons .btn:hover {
  transform: translateY(-1px);
}

/* Empty State */
.empty-state {
  text-align: center;
  padding: 3rem 2rem;
  background: white;
  border-radius: 12px;
  border: 1px solid #e9ecef;
}

.empty-icon {
  font-size: 4rem;
  color: #6c757d;
  margin-bottom: 1rem;
}

.empty-state h5 {
  color: #495057;
  margin-bottom: 0.5rem;
}

/* Pagination */
.pagination-wrapper {
  margin-top: 1.5rem;
  padding: 1rem;
  background: white;
  border-radius: 12px;
  border: 1px solid #e9ecef;
}

.pagination .page-link {
  border: none;
  color: #6c757d;
  padding: 0.5rem 0.75rem;
  margin: 0 0.125rem;
  border-radius: 6px;
  transition: all 0.2s ease;
}

.pagination .page-link:hover {
  background-color: #e9ecef;
  color: #495057;
}

.pagination .page-item.active .page-link {
  background: linear-gradient(45deg, #4facfe 0%, #00f2fe 100%);
  border: none;
  color: white;
}

.pagination .page-item.disabled .page-link {
  color: #adb5bd;
  background-color: transparent;
}

/* Responsive Design */
@media (max-width: 1200px) {
  .col-tags {
    display: none;
  }
  .col-actions {
    width: 20%;
  }
}

@media (max-width: 992px) {
  .col-remind-time {
    width: 25%;
  }
  .col-actions {
    width: 25%;
  }
}

@media (max-width: 768px) {
  .list-header {
    flex-direction: column;
    gap: 1rem;
    align-items: stretch;
  }
  
  .search-filters {
    flex-direction: column;
    gap: 0.75rem;
  }
  
  .filter-group {
    justify-content: stretch;
  }
  
  .filter-group .form-select {
    flex: 1;
    min-width: unset;
  }
  
  .table-responsive {
    font-size: 0.875rem;
  }
  
  .action-buttons .btn {
    padding: 0.25rem 0.375rem;
    font-size: 0.75rem;
  }
  
  .content-preview {
    display: none;
  }
}

/* Notification Info */
.notification-info {
  min-width: 120px;
}

.notification-methods {
  display: flex;
  flex-wrap: wrap;
  gap: 0.25rem;
}

.notification-methods .badge {
  font-size: 0.7rem;
  padding: 0.25rem 0.5rem;
  border-radius: 4px;
}

@media (max-width: 576px) {
  .col-status, .col-notifications {
    display: none;
  }
  .col-name {
    width: 50%;
  }
  .col-remind-time {
    width: 30%;
  }
  .col-actions {
    width: 20%;
  }
}
</style>
