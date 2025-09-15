`<template>
  <div class="tasks-view">
    <!-- Page Header -->
    <div class="page-header">
      <div class="header-content">
        <h1 class="page-title">
          <i class="bi bi-check2-square"></i>
          Task Management
        </h1>
        <p class="page-description">
          Organize and track your tasks with priority levels, deadlines, and progress monitoring
        </p>
      </div>
      <div class="header-actions">
        <button class="btn btn-primary" @click="showCreateForm">
          <i class="bi bi-plus-lg me-2"></i>
          New Task
        </button>
      </div>
    </div>

    <!-- Main Content -->
    <div class="content-wrapper">
      <!-- Create/Edit Form -->
      <div v-if="showForm" class="form-section">
        <div class="section-header">
          <h2>{{ editingTask ? 'Edit Task' : 'Create New Task' }}</h2>
          <button 
            class="btn btn-sm btn-light"
            @click="closeForm"
          >
            <i class="bi bi-x-lg"></i>
          </button>
        </div>
        <TaskForm
          :task="editingTask"
          :submit-button-text="editingTask ? 'Update Task' : 'Add Task'"
          @submit="handleFormSubmit"
          @cancel="closeForm"
        />
      </div>

      <!-- Task List -->
      <div v-else>
        <TaskList
          :tasks="taskStore.tasks"
          :search-query="searchQuery"
          :current-page="currentPage"
          :page-size="pageSize"
          :total-count="taskStore.totalCount"
          @view="handleView"
          @edit="handleEdit"
          @delete="handleDelete"
          @start="handleStart"
          @complete="handleComplete"
          @fail="handleFail"
          @update:search-query="searchQuery = $event"
          @update:page="currentPage = $event"
        />
      </div>
    </div>

    <!-- Loading Overlay -->
    <div v-if="taskStore.loading" class="loading-overlay">
      <div class="spinner-border text-primary" role="status">
        <span class="visually-hidden">Loading...</span>
      </div>
    </div>

    <!-- Task Details Modal -->
    <div 
      v-if="viewingTask"
      class="modal fade show d-block"
      tabindex="-1"
      @click.self="viewingTask = null"
    >
      <div class="modal-dialog modal-lg">
        <div class="modal-content">
          <div class="modal-header">
            <h5 class="modal-title">
              <i class="bi bi-check2-square me-2"></i>
              Task Details
            </h5>
            <button 
              type="button" 
              class="btn-close" 
              @click="viewingTask = null"
            ></button>
          </div>
          <div class="modal-body">
            <div class="task-details">
              <div class="detail-row">
                <span class="detail-label">Name:</span>
                <span class="detail-value fw-bold">{{ viewingTask.name }}</span>
              </div>
              <div class="detail-row">
                <span class="detail-label">Description:</span>
                <span class="detail-value">{{ viewingTask.description || '-' }}</span>
              </div>
              <div class="detail-row">
                <span class="detail-label">Priority:</span>
                <span class="detail-value">
                  <span :class="getPriorityBadgeClass(viewingTask.priority)">
                    {{ viewingTask.priority }}
                  </span>
                </span>
              </div>
              <div class="detail-row">
                <span class="detail-label">Difficulty:</span>
                <span class="detail-value">
                  <span :class="getDifficultyBadgeClass(viewingTask.difficulty)">
                    {{ viewingTask.difficulty }}
                  </span>
                </span>
              </div>
              <div class="detail-row">
                <span class="detail-label">Status:</span>
                <span class="detail-value">
                  <span :class="getStatusBadgeClass(viewingTask.status)">
                    {{ viewingTask.status }}
                  </span>
                </span>
              </div>
              <div class="detail-row">
                <span class="detail-label">Duration:</span>
                <span class="detail-value">{{ viewingTask.minutes }} minutes</span>
              </div>
              <div class="detail-row">
                <span class="detail-label">Schedule Time:</span>
                <span class="detail-value">{{ formatDate(viewingTask.schedule_time) }}</span>
              </div>
              <div class="detail-row">
                <span class="detail-label">Deadline:</span>
                <span class="detail-value">{{ formatDate(viewingTask.deadline) }}</span>
              </div>
              <div class="detail-row" v-if="viewingTask.start_time">
                <span class="detail-label">Started:</span>
                <span class="detail-value">{{ formatDate(viewingTask.start_time) }}</span>
              </div>
              <div class="detail-row" v-if="viewingTask.end_time">
                <span class="detail-label">Completed:</span>
                <span class="detail-value">{{ formatDate(viewingTask.end_time) }}</span>
              </div>
              <div class="detail-row" v-if="viewingTask.tags && viewingTask.tags.length > 0">
                <span class="detail-label">Tags:</span>
                <span class="detail-value">
                  <span 
                    v-for="tag in viewingTask.tags" 
                    :key="tag" 
                    class="badge bg-light text-dark me-1"
                  >
                    {{ tag }}
                  </span>
                </span>
              </div>
            </div>
          </div>
          <div class="modal-footer">
            <button 
              type="button" 
              class="btn btn-secondary" 
              @click="viewingTask = null"
            >
              Close
            </button>
            <button 
              v-if="viewingTask.status === 'pending'"
              type="button" 
              class="btn btn-success" 
              @click="handleStart(viewingTask.id)"
            >
              <i class="bi bi-play me-2"></i>
              Start Task
            </button>
            <button 
              v-if="viewingTask.status === 'running'"
              type="button" 
              class="btn btn-primary" 
              @click="handleComplete(viewingTask.id)"
            >
              <i class="bi bi-check me-2"></i>
              Complete
            </button>
            <button 
              v-if="viewingTask.status === 'running'"
              type="button" 
              class="btn btn-warning" 
              @click="handleFail(viewingTask.id)"
            >
              <i class="bi bi-x me-2"></i>
              Mark Failed
            </button>
          </div>
        </div>
      </div>
    </div>
    <div v-if="viewingTask" class="modal-backdrop fade show"></div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted, watch } from 'vue';
import { useTaskStore } from '@/stores/taskStore';
import type { Task } from '@/types';
import { formatDate } from '@/utils/dateUtils';
import TaskForm from '@/components/tasks/TaskForm.vue';
import TaskList from '@/components/tasks/TaskList.vue';

const taskStore = useTaskStore();

// UI State
const showForm = ref(false);
const editingTask = ref<Task | undefined>(undefined);
const viewingTask = ref<Task | null>(null);

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
    loadTasks();
  }, 300);
});

watch(currentPage, () => {
  loadTasks();
});

// Load tasks
const loadTasks = async () => {
  try {
    await taskStore.fetchTasks({
      q: searchQuery.value,
      page: currentPage.value,
      page_size: pageSize.value
    });
  } catch (error) {
    console.error('Failed to load tasks:', error);
  }
};

// Form handlers
const showCreateForm = () => {
  editingTask.value = undefined;
  showForm.value = true;
};

const closeForm = () => {
  showForm.value = false;
  editingTask.value = undefined;
};

const handleFormSubmit = async (task: Task) => {
  try {
    if (editingTask.value) {
      await taskStore.updateTask(editingTask.value.id, task);
    } else {
      await taskStore.addTask(task);
    }
    closeForm();
    await loadTasks();
  } catch (error) {
    console.error('Failed to save task:', error);
    alert('Failed to save task. Please try again.');
  }
};

// List action handlers
const handleView = (task: Task) => {
  viewingTask.value = task;
};

const handleEdit = (task: Task) => {
  editingTask.value = task;
  showForm.value = true;
};

const handleDelete = async (id: string) => {
  if (confirm('Are you sure you want to delete this task?')) {
    try {
      await taskStore.deleteTask(id);
      await loadTasks();
    } catch (error) {
      console.error('Failed to delete task:', error);
      alert('Failed to delete task. Please try again.');
    }
  }
};

const handleStart = async (id: string) => {
  try {
    await taskStore.startTask(id);
    await loadTasks();
    if (viewingTask.value) {
      viewingTask.value = null;
    }
  } catch (error) {
    console.error('Failed to start task:', error);
  }
};

const handleComplete = async (id: string) => {
  try {
    await taskStore.completeTask(id);
    await loadTasks();
    if (viewingTask.value) {
      viewingTask.value = null;
    }
  } catch (error) {
    console.error('Failed to complete task:', error);
  }
};

const handleFail = async (id: string) => {
  try {
    await taskStore.failTask(id);
    await loadTasks();
    if (viewingTask.value) {
      viewingTask.value = null;
    }
  } catch (error) {
    console.error('Failed to mark task as failed:', error);
  }
};

// Badge helper functions
const getPriorityBadgeClass = (priority: string) => {
  const classes = {
    'low': 'badge bg-success',
    'medium': 'badge bg-warning',
    'high': 'badge bg-danger'
  };
  return classes[priority as keyof typeof classes] || 'badge bg-secondary';
};

const getDifficultyBadgeClass = (difficulty: string) => {
  const classes = {
    'easy': 'badge bg-success',
    'medium': 'badge bg-warning',
    'hard': 'badge bg-danger'
  };
  return classes[difficulty as keyof typeof classes] || 'badge bg-secondary';
};

const getStatusBadgeClass = (status: string) => {
  const classes = {
    'pending': 'badge bg-secondary',
    'running': 'badge bg-primary',
    'completed': 'badge bg-success',
    'failed': 'badge bg-danger'
  };
  return classes[status as keyof typeof classes] || 'badge bg-secondary';
};

// Initialize
onMounted(() => {
  loadTasks();
});
</script>

<style scoped>
/* Page Layout */
.tasks-view {
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

.task-details {
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

/* Badge Styles */
.badge {
  font-size: 0.75rem;
  padding: 0.25rem 0.5rem;
  font-weight: 500;
}

/* Responsive Design */
@media (max-width: 768px) {
  .tasks-view {
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