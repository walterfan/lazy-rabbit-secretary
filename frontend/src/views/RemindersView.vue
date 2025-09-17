<template>
  <div class="reminders-view">
    <!-- Page Header -->
    <div class="page-header">
      <div class="header-content">
        <h1 class="page-title">
          <i class="bi bi-bell me-3"></i>
          Reminders
        </h1>
        <p class="page-subtitle">Manage your reminders and never miss important events</p>
      </div>
      <div class="header-actions">
        <button 
          class="btn btn-primary btn-lg"
          @click="showForm = true"
        >
          <i class="bi bi-plus-lg me-2"></i>
          New Reminder
        </button>
      </div>
    </div>

    <!-- Quick Stats -->
    <div class="stats-row">
      <div class="stat-card">
        <div class="stat-icon bg-warning">
          <i class="bi bi-clock"></i>
        </div>
        <div class="stat-content">
          <div class="stat-number">{{ stats.pending }}</div>
          <div class="stat-label">Pending</div>
        </div>
      </div>
      <div class="stat-card">
        <div class="stat-icon bg-danger">
          <i class="bi bi-exclamation-triangle"></i>
        </div>
        <div class="stat-content">
          <div class="stat-number">{{ stats.overdue }}</div>
          <div class="stat-label">Overdue</div>
        </div>
      </div>
      <div class="stat-card">
        <div class="stat-icon bg-info">
          <i class="bi bi-bell"></i>
        </div>
        <div class="stat-content">
          <div class="stat-number">{{ stats.active }}</div>
          <div class="stat-label">Active</div>
        </div>
      </div>
      <div class="stat-card">
        <div class="stat-icon bg-success">
          <i class="bi bi-check-circle"></i>
        </div>
        <div class="stat-content">
          <div class="stat-number">{{ stats.completed }}</div>
          <div class="stat-label">Completed</div>
        </div>
      </div>
    </div>

    <!-- Main Content -->
    <div class="main-content">
      <!-- Form Section -->
      <div v-if="showForm" class="form-section">
        <div class="form-container">
          <div class="form-header">
            <h3 class="form-title">
              <i :class="editingReminder ? 'bi-pencil-square' : 'bi-plus-lg'" class="bi me-2"></i>
              {{ editingReminder ? 'Edit Reminder' : 'Create New Reminder' }}
            </h3>
            <button class="btn btn-outline-secondary" @click="closeForm">
              <i class="bi bi-x-lg"></i>
            </button>
          </div>
          <ReminderForm
            :reminder="editingReminder || undefined"
            :loading="reminderStore.loading"
            @submit="handleFormSubmit"
            @cancel="closeForm"
          />
        </div>
      </div>

      <!-- List Section -->
      <div class="list-section">
        <ReminderList
          :reminders="reminders"
          :search-query="searchQuery"
          :total-count="totalCount"
          :current-page="currentPage"
          :total-pages="totalPages"
          :loading="reminderStore.loading"
          @update:search-query="searchQuery = $event"
          @update:filters="handleFiltersUpdate"
          @page-change="handlePageChange"
          @view="handleView"
          @edit="handleEdit"
          @delete="handleDelete"
          @complete="handleComplete"
          @snooze="handleSnooze"
        />
      </div>
    </div>

    <!-- View Modal -->
    <div 
      v-if="viewingReminder"
      class="modal fade show d-block"
      tabindex="-1"
      @click.self="viewingReminder = null"
    >
      <div class="modal-dialog modal-lg">
        <div class="modal-content">
          <div class="modal-header">
            <h5 class="modal-title">
              <i class="bi bi-bell me-2"></i>
              Reminder Details
            </h5>
            <button 
              type="button" 
              class="btn-close" 
              @click="viewingReminder = null"
            ></button>
          </div>
          <div class="modal-body">
            <div class="reminder-details">
              <div class="detail-row">
                <span class="detail-label">Name:</span>
                <span class="detail-value fw-bold">{{ viewingReminder.name }}</span>
              </div>
              <div class="detail-row">
                <span class="detail-label">Content:</span>
                <span class="detail-value">{{ viewingReminder.content }}</span>
              </div>
              <div class="detail-row">
                <span class="detail-label">Status:</span>
                <ReminderStatusBadge :status="viewingReminder.status" />
              </div>
              <div class="detail-row">
                <span class="detail-label">Remind Time:</span>
                <span class="detail-value">{{ formatDate(viewingReminder.remind_time) }}</span>
              </div>
              <div class="detail-row">
                <span class="detail-label">Tags:</span>
                <ReminderTags :tags="viewingReminder.tags" />
              </div>
              <div class="detail-row">
                <span class="detail-label">Created:</span>
                <span class="detail-value">{{ formatDate(viewingReminder.created_at) }}</span>
              </div>
              <div class="detail-row">
                <span class="detail-label">Updated:</span>
                <span class="detail-value">{{ formatDate(viewingReminder.updated_at) }}</span>
              </div>
            </div>
          </div>
          <div class="modal-footer">
            <button 
              type="button" 
              class="btn btn-secondary" 
              @click="viewingReminder = null"
            >
              Close
            </button>
            <button 
              type="button" 
              class="btn btn-primary" 
              @click="handleEdit(viewingReminder!)"
            >
              <i class="bi bi-pencil me-2"></i>
              Edit
            </button>
          </div>
        </div>
      </div>
    </div>
    <div v-if="viewingReminder" class="modal-backdrop fade show"></div>

    <!-- Snooze Modal -->
    <div 
      v-if="snoozeReminder"
      class="modal fade show d-block"
      tabindex="-1"
      @click.self="closeSnoozeModal"
    >
      <div class="modal-dialog">
        <div class="modal-content">
          <div class="modal-header">
            <h5 class="modal-title">
              <i class="bi bi-clock me-2"></i>
              Snooze Reminder
            </h5>
            <button 
              type="button" 
              class="btn-close" 
              @click="closeSnoozeModal"
            ></button>
          </div>
          <div class="modal-body">
            <div class="mb-3">
              <label for="snooze-time" class="form-label">New Remind Time</label>
              <input
                type="datetime-local"
                class="form-control"
                id="snooze-time"
                v-model="snoozeTimeString"
                :min="minSnoozeTime"
              />
            </div>
            <div class="quick-snooze-options">
              <p class="mb-2 fw-semibold">Quick Options:</p>
              <div class="btn-group-vertical d-grid gap-2">
                <button 
                  class="btn btn-outline-primary btn-sm"
                  @click="setQuickSnooze(15)"
                >
                  <i class="bi bi-clock me-2"></i>
                  15 minutes
                </button>
                <button 
                  class="btn btn-outline-primary btn-sm"
                  @click="setQuickSnooze(60)"
                >
                  <i class="bi bi-clock me-2"></i>
                  1 hour
                </button>
                <button 
                  class="btn btn-outline-primary btn-sm"
                  @click="setQuickSnooze(1440)"
                >
                  <i class="bi bi-calendar me-2"></i>
                  Tomorrow
                </button>
              </div>
            </div>
          </div>
          <div class="modal-footer">
            <button 
              type="button" 
              class="btn btn-secondary" 
              @click="closeSnoozeModal"
            >
              Cancel
            </button>
            <button 
              type="button" 
              class="btn btn-warning" 
              @click="handleSnoozeConfirm"
              :disabled="!snoozeTimeString || snoozeLoading"
            >
              <span v-if="snoozeLoading" class="spinner-border spinner-border-sm me-2" role="status"></span>
              <i v-else class="bi bi-clock me-2"></i>
              {{ snoozeLoading ? 'Snoozing...' : 'Snooze' }}
            </button>
          </div>
        </div>
      </div>
    </div>
    <div v-if="snoozeReminder" class="modal-backdrop fade show"></div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted, watch, computed } from 'vue';
import { useReminderStore } from '@/stores/reminderStore';
import type { Reminder, CreateReminderRequest, UpdateReminderRequest } from '@/types';
import { formatDate } from '@/utils/dateUtils';
import ReminderForm from '@/components/reminders/ReminderForm.vue';
import ReminderList from '@/components/reminders/ReminderList.vue';
import ReminderStatusBadge from '@/components/reminders/ReminderStatusBadge.vue';
import ReminderTags from '@/components/reminders/ReminderTags.vue';

const reminderStore = useReminderStore();

// UI State
const showForm = ref(false);
const editingReminder = ref<Reminder | null>(null);
const viewingReminder = ref<Reminder | null>(null);

// Snooze Modal State
const snoozeReminder = ref<Reminder | null>(null);
const snoozeTimeString = ref('');
const snoozeLoading = ref(false);

// Search and Filter State
const searchQuery = ref('');
const filters = ref({
  status: '',
  tags: ''
});
const currentPage = ref(1);
const pageSize = ref(20);

// Data State
const reminders = ref<Reminder[]>([]);
const totalCount = ref(0);
const totalPages = ref(0);

// Stats
const stats = computed(() => {
  const now = new Date();
  return {
    pending: reminders.value.filter(r => r.status === 'pending').length,
    active: reminders.value.filter(r => r.status === 'active').length,
    completed: reminders.value.filter(r => r.status === 'completed').length,
    overdue: reminders.value.filter(r => 
      r.status === 'pending' && new Date(r.remind_time) < now
    ).length
  };
});

const minSnoozeTime = computed(() => {
  const now = new Date();
  now.setMinutes(now.getMinutes() + 1); // At least 1 minute in the future
  return now.toISOString().slice(0, 16);
});

// Debounced search
let searchTimeout: ReturnType<typeof setTimeout>;
watch([searchQuery, filters], () => {
  clearTimeout(searchTimeout);
  searchTimeout = setTimeout(() => {
    currentPage.value = 1;
    loadReminders();
  }, 300);
});

watch(currentPage, () => {
  loadReminders();
});

// Methods
const loadReminders = async () => {
  try {
    const response = await reminderStore.fetchReminders({
      q: searchQuery.value,
      status: filters.value.status,
      tags: filters.value.tags,
      page: currentPage.value,
      page_size: pageSize.value
    });
    
    reminders.value = response.items;
    totalCount.value = response.total;
    totalPages.value = response.total_pages;
  } catch (error) {
    console.error('Failed to load reminders:', error);
  }
};

const handleFormSubmit = async (reminderData: CreateReminderRequest | UpdateReminderRequest) => {
  try {
    if (editingReminder.value) {
      await reminderStore.updateReminder(editingReminder.value.id, reminderData as UpdateReminderRequest);
    } else {
      await reminderStore.createReminder(reminderData as CreateReminderRequest);
    }
    closeForm();
    await loadReminders();
  } catch (error) {
    console.error('Failed to save reminder:', error);
  }
};

const handleFiltersUpdate = (newFilters: any) => {
  filters.value = { ...newFilters };
};

const handlePageChange = (page: number) => {
  currentPage.value = page;
};

const handleView = (reminder: Reminder) => {
  viewingReminder.value = reminder;
};

const handleEdit = (reminder: Reminder) => {
  editingReminder.value = reminder;
  showForm.value = true;
  if (viewingReminder.value) {
    viewingReminder.value = null;
  }
};

const handleDelete = async (id: string) => {
  if (confirm('Are you sure you want to delete this reminder?')) {
    try {
      await reminderStore.deleteReminder(id);
      await loadReminders();
    } catch (error) {
      console.error('Failed to delete reminder:', error);
      alert('Failed to delete reminder. Please try again.');
    }
  }
};

const handleComplete = async (id: string) => {
  try {
    await reminderStore.markAsCompleted(id);
    await loadReminders();
  } catch (error) {
    console.error('Failed to complete reminder:', error);
    alert('Failed to complete reminder. Please try again.');
  }
};

const handleSnooze = (reminder: Reminder) => {
  snoozeReminder.value = reminder;
  const futureTime = new Date();
  futureTime.setMinutes(futureTime.getMinutes() + 15); // Default to 15 minutes
  snoozeTimeString.value = futureTime.toISOString().slice(0, 16);
};

const setQuickSnooze = (minutes: number) => {
  const futureTime = new Date();
  futureTime.setMinutes(futureTime.getMinutes() + minutes);
  snoozeTimeString.value = futureTime.toISOString().slice(0, 16);
};

const handleSnoozeConfirm = async () => {
  if (!snoozeReminder.value || !snoozeTimeString.value) return;
  
  snoozeLoading.value = true;
  try {
    const newRemindTime = new Date(snoozeTimeString.value);
    await reminderStore.snoozeReminder(snoozeReminder.value.id, newRemindTime);
    closeSnoozeModal();
    await loadReminders();
  } catch (error) {
    console.error('Failed to snooze reminder:', error);
    alert('Failed to snooze reminder. Please try again.');
  } finally {
    snoozeLoading.value = false;
  }
};

const closeForm = () => {
  showForm.value = false;
  editingReminder.value = null;
};

const closeSnoozeModal = () => {
  snoozeReminder.value = null;
  snoozeTimeString.value = '';
  snoozeLoading.value = false;
};

// Initialize
onMounted(() => {
  loadReminders();
});
</script>

<style scoped>
/* Page Layout */
.reminders-view {
  padding: 2rem;
  max-width: 1400px;
  margin: 0 auto;
}

/* Page Header */
.page-header {
  display: flex;
  justify-content: space-between;
  align-items: flex-start;
  margin-bottom: 2rem;
  padding-bottom: 1rem;
  border-bottom: 2px solid #f8f9fa;
}

.header-content h1 {
  color: #2c3e50;
  font-weight: 700;
  margin-bottom: 0.5rem;
}

.page-subtitle {
  color: #6c757d;
  font-size: 1.1rem;
  margin: 0;
}

.header-actions .btn {
  box-shadow: 0 4px 12px rgba(79, 172, 254, 0.3);
  transition: all 0.3s ease;
}

.header-actions .btn:hover {
  transform: translateY(-2px);
  box-shadow: 0 6px 20px rgba(79, 172, 254, 0.4);
}

/* Stats Row */
.stats-row {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(200px, 1fr));
  gap: 1.5rem;
  margin-bottom: 2rem;
}

.stat-card {
  background: white;
  border-radius: 12px;
  padding: 1.5rem;
  border: 1px solid #e9ecef;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.05);
  display: flex;
  align-items: center;
  gap: 1rem;
  transition: all 0.3s ease;
}

.stat-card:hover {
  transform: translateY(-2px);
  box-shadow: 0 4px 16px rgba(0, 0, 0, 0.1);
}

.stat-icon {
  width: 48px;
  height: 48px;
  border-radius: 12px;
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 1.5rem;
  color: white;
}

.stat-number {
  font-size: 2rem;
  font-weight: 700;
  color: #2c3e50;
  line-height: 1;
}

.stat-label {
  color: #6c757d;
  font-weight: 500;
  text-transform: uppercase;
  font-size: 0.875rem;
  letter-spacing: 0.5px;
}

/* Main Content */
.main-content {
  display: flex;
  flex-direction: column;
  gap: 2rem;
}

/* Form Section */
.form-section {
  background: white;
  border-radius: 16px;
  border: 1px solid #e9ecef;
  box-shadow: 0 4px 16px rgba(0, 0, 0, 0.05);
  overflow: hidden;
}

.form-container {
  padding: 0;
}

.form-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 1.5rem 2rem;
  background: linear-gradient(135deg, #f8f9fa 0%, #e9ecef 100%);
  border-bottom: 1px solid #e9ecef;
}

.form-title {
  margin: 0;
  color: #495057;
  font-weight: 600;
}

.form-container :deep(.reminder-form) {
  padding: 2rem;
}

/* Modal Styles */
.modal-content {
  border-radius: 16px;
  border: none;
  box-shadow: 0 10px 40px rgba(0, 0, 0, 0.15);
}

.modal-header {
  background: linear-gradient(135deg, #f8f9fa 0%, #e9ecef 100%);
  border-bottom: 1px solid #e9ecef;
  border-radius: 16px 16px 0 0;
}

.reminder-details {
  display: flex;
  flex-direction: column;
  gap: 1rem;
}

.detail-row {
  display: flex;
  align-items: center;
  gap: 1rem;
  padding: 0.75rem 0;
  border-bottom: 1px solid #f8f9fa;
}

.detail-row:last-child {
  border-bottom: none;
}

.detail-label {
  font-weight: 600;
  color: #495057;
  min-width: 120px;
  flex-shrink: 0;
}

.detail-value {
  color: #6c757d;
  flex: 1;
}

/* Quick Snooze Options */
.quick-snooze-options {
  margin-top: 1rem;
  padding-top: 1rem;
  border-top: 1px solid #e9ecef;
}

/* Responsive Design */
@media (max-width: 1200px) {
  .reminders-view {
    padding: 1.5rem;
  }
}

@media (max-width: 768px) {
  .reminders-view {
    padding: 1rem;
  }
  
  .page-header {
    flex-direction: column;
    gap: 1rem;
    align-items: stretch;
  }
  
  .header-actions {
    display: flex;
    justify-content: center;
  }
  
  .stats-row {
    grid-template-columns: repeat(2, 1fr);
    gap: 1rem;
  }
  
  .stat-card {
    padding: 1rem;
  }
  
  .stat-number {
    font-size: 1.5rem;
  }
  
  .form-header {
    padding: 1rem 1.5rem;
  }
  
  .form-container :deep(.reminder-form) {
    padding: 1.5rem;
  }
}

@media (max-width: 576px) {
  .stats-row {
    grid-template-columns: 1fr;
  }
  
  .stat-card {
    padding: 1rem;
    gap: 0.75rem;
  }
  
  .stat-icon {
    width: 40px;
    height: 40px;
    font-size: 1.25rem;
  }
  
  .detail-row {
    flex-direction: column;
    align-items: flex-start;
    gap: 0.5rem;
  }
  
  .detail-label {
    min-width: unset;
  }
}
</style>
