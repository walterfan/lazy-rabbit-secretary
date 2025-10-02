<template>
  <div class="daily-panel">
    <!-- Date Selector and Quick Stats -->
    <div class="daily-header">
      <div class="date-selector">
        <button class="btn btn-outline-secondary" @click="changeDate(-1)">
          <i class="bi bi-chevron-left"></i>
        </button>
        <div class="date-display">
          <h3>{{ formatDateHeader(selectedDate) }}</h3>
          <p class="date-subtitle">{{ getDayOfWeek(selectedDate) }}</p>
        </div>
        <button class="btn btn-outline-secondary" @click="changeDate(1)">
          <i class="bi bi-chevron-right"></i>
        </button>
        <button class="btn btn-primary" @click="selectedDate = new Date()">
          {{ $t('gtd.today') }}
        </button>
      </div>
      
      <div class="daily-stats" v-if="dailyStats">
        <div class="stat-item">
          <div class="stat-value">{{ dailyStats.completion_rate.toFixed(1) }}%</div>
          <div class="stat-label">{{ $t('gtd.completionRate') }}</div>
        </div>
        <div class="stat-item">
          <div class="stat-value">{{ dailyStats.completed_items }}/{{ dailyStats.total_items }}</div>
          <div class="stat-label">{{ $t('gtd.completedItems') }}</div>
        </div>
        <div class="stat-item">
          <div class="stat-value">{{ formatTime(dailyStats.total_actual_time) }}</div>
          <div class="stat-label">{{ $t('gtd.actualTime') }}</div>
        </div>
        <div class="stat-item">
          <div class="stat-value">{{ formatTime(dailyStats.total_estimated_time) }}</div>
          <div class="stat-label">{{ $t('gtd.estimatedTime') }}</div>
        </div>
      </div>
    </div>

    <!-- Quick Add Task -->
    <div class="quick-add-section">
      <div class="quick-add-card">
        <div class="quick-add-header">
          <h4>
            <i class="bi bi-plus-circle"></i>
            {{ $t('gtd.addTask') }}
          </h4>
        </div>
        <form @submit.prevent="handleQuickAdd" class="quick-add-form">
          <div class="form-row">
            <div class="form-group">
              <input
                type="text"
                v-model="quickAddTitle"
                :placeholder="$t('gtd.taskTitle')"
                class="form-control"
                required
              />
            </div>
            <div class="form-group">
              <select v-model="quickAddPriority" class="form-select">
                <option value="A">{{ $t('gtd.aLevel') }}</option>
                <option value="B+">{{ $t('gtd.bPlusLevel') }}</option>
                <option value="B">{{ $t('gtd.bLevel') }}</option>
                <option value="C">{{ $t('gtd.cLevel') }}</option>
                <option value="D">{{ $t('gtd.dLevel') }}</option>
              </select>
            </div>
            <div class="form-group">
              <input
                type="number"
                v-model="quickAddEstimatedTime"
                :placeholder="$t('gtd.estimatedTime')"
                class="form-control"
                min="0"
                max="1440"
              />
            </div>
            <button type="submit" class="btn btn-primary">
              <i class="bi bi-plus"></i>
              {{ $t('gtd.add') }}
            </button>
          </div>
        </form>
      </div>
    </div>

    <!-- Priority Sections -->
    <div class="priority-sections">
      <div
        v-for="priority in priorityOrder"
        :key="priority.value"
        class="priority-section"
        :class="priority.value"
      >
        <div class="priority-header">
          <h4>
            <i :class="priority.icon"></i>
            {{ priority.label }}
            <span class="priority-count">({{ getPriorityCount(priority.value) }})</span>
          </h4>
          <div class="priority-progress">
            <div class="progress-bar">
              <div
                class="progress-fill"
                :style="{ width: getPriorityProgress(priority.value) + '%' }"
              ></div>
            </div>
            <span class="progress-text">{{ getPriorityProgress(priority.value).toFixed(0) }}%</span>
          </div>
        </div>
        
        <div class="tasks-list">
          <div
            v-for="task in getTasksByPriority(priority.value)"
            :key="task.id"
            :class="['task-card', task.status]"
            @click="selectTask(task)"
          >
            <div class="task-header">
              <div class="task-checkbox">
                <input
                  type="checkbox"
                  :checked="task.status === 'completed'"
                  @change="toggleTaskStatus(task)"
                  @click.stop
                />
              </div>
              <div class="task-title">
                <h5>{{ task.title }}</h5>
                <p v-if="task.description" class="task-description">{{ task.description }}</p>
              </div>
              <div class="task-actions">
                <button
                  class="btn btn-sm btn-outline-primary"
                  @click.stop="editTask(task)"
                  :title="$t('gtd.edit')"
                >
                  <i class="bi bi-pencil"></i>
                </button>
                <button
                  class="btn btn-sm btn-outline-danger"
                  @click.stop="deleteTask(task.id)"
                  :title="$t('gtd.delete')"
                >
                  <i class="bi bi-trash"></i>
                </button>
              </div>
            </div>
            
            <div class="task-meta">
              <div class="task-time">
                <span class="time-estimate">
                  <i class="bi bi-clock"></i>
                  {{ formatTime(task.estimated_time) }}
                </span>
                <span v-if="task.actual_time > 0" class="time-actual">
                  <i class="bi bi-stopwatch"></i>
                  {{ formatTime(task.actual_time) }}
                </span>
              </div>
              
              <div class="task-context" v-if="task.context">
                <span class="context-badge">
                  <i :class="getContextIcon(task.context)"></i>
                  {{ task.context }}
                </span>
              </div>
              
              <div class="task-deadline" v-if="task.deadline">
                <span class="deadline-badge" :class="getDeadlineClass(task.deadline)">
                  <i class="bi bi-calendar-event"></i>
                  {{ formatTime(task.deadline) }}
                </span>
              </div>
            </div>
            
            <div class="task-status">
              <span :class="['status-badge', task.status]">
                {{ getStatusLabel(task.status) }}
              </span>
              <span v-if="task.completion_time" class="completion-time">
                完成于 {{ formatTime(task.completion_time) }}
              </span>
            </div>
          </div>
        </div>
      </div>
    </div>

    <!-- Empty State -->
    <div class="empty-state" v-if="!loading && items.length === 0">
      <div class="empty-icon">
        <i class="bi bi-calendar-check"></i>
      </div>
      <h4>{{ $t('gtd.noTasksToday') }}</h4>
      <p>{{ $t('gtd.startPlanning') }}</p>
    </div>

    <!-- Loading State -->
    <div class="loading-state" v-if="loading">
      <div class="spinner-border text-primary" role="status">
        <span class="visually-hidden">Loading...</span>
      </div>
      <p>{{ $t('gtd.loading') }}</p>
    </div>

    <!-- Task Detail Modal -->
    <div v-if="selectedTask" class="modal fade show d-block" tabindex="-1" @click.self="closeTaskModal">
      <div class="modal-dialog modal-lg">
        <div class="modal-content">
          <div class="modal-header">
            <h5 class="modal-title">{{ $t('gtd.taskDetails') }}</h5>
            <button type="button" class="btn-close" @click="closeTaskModal"></button>
          </div>
          <div class="modal-body">
            <div class="task-detail">
              <div class="detail-section">
                <h6>{{ $t('gtd.basicInfo') }}</h6>
                <div class="detail-grid">
                  <div class="detail-item">
                    <label>{{ $t('gtd.title') }}:</label>
                    <span>{{ selectedTask.title }}</span>
                  </div>
                  <div class="detail-item">
                    <label>{{ $t('gtd.priority') }}:</label>
                    <span :class="['priority-badge', selectedTask.priority]">
                      {{ selectedTask.priority }}
                    </span>
                  </div>
                  <div class="detail-item">
                    <label>{{ $t('gtd.status') }}:</label>
                    <span :class="['status-badge', selectedTask.status]">
                      {{ getStatusLabel(selectedTask.status) }}
                    </span>
                  </div>
                  <div class="detail-item">
                    <label>{{ $t('gtd.estimatedTime') }}:</label>
                    <span>{{ formatTime(selectedTask.estimated_time) }}</span>
                  </div>
                  <div class="detail-item">
                    <label>{{ $t('gtd.actualTime') }}:</label>
                    <span>{{ formatTime(selectedTask.actual_time) }}</span>
                  </div>
                  <div class="detail-item" v-if="selectedTask.deadline">
                    <label>{{ $t('gtd.deadline') }}:</label>
                    <span>{{ formatDateTime(selectedTask.deadline, locale) }}</span>
                  </div>
                </div>
              </div>
              
              <div class="detail-section" v-if="selectedTask.description">
                <h6>{{ $t('gtd.description') }}</h6>
                <p>{{ selectedTask.description }}</p>
              </div>
              
              <div class="detail-section" v-if="selectedTask.notes">
                <h6>{{ $t('gtd.executionNotes') }}</h6>
                <p>{{ selectedTask.notes }}</p>
              </div>
              
              <div class="detail-section">
                <h6>{{ $t('gtd.timeTracking') }}</h6>
                <div class="time-tracking">
                  <div class="time-input-group">
                    <label>{{ $t('gtd.recordActualTime') }}</label>
                    <div class="input-group">
                      <input
                        type="number"
                        v-model="actualTimeInput"
                        class="form-control"
                        :placeholder="$t('gtd.minutes')"
                        min="0"
                        max="1440"
                      />
                      <button
                        class="btn btn-outline-primary"
                        @click="updateActualTime"
                      >
                        {{ $t('gtd.update') }}
                      </button>
                    </div>
                  </div>
                </div>
              </div>
            </div>
          </div>
          <div class="modal-footer">
            <button type="button" class="btn btn-secondary" @click="closeTaskModal">{{ $t('gtd.close') }}</button>
            <button type="button" class="btn btn-primary" @click="editTask(selectedTask)">{{ $t('gtd.edit') }}</button>
          </div>
        </div>
      </div>
    </div>
    <div v-if="selectedTask" class="modal-backdrop fade show"></div>

    <!-- Edit Task Modal -->
    <div v-if="editingTask" class="modal fade show d-block" tabindex="-1" @click.self="closeEditModal">
      <div class="modal-dialog">
        <div class="modal-content">
          <div class="modal-header">
            <h5 class="modal-title">{{ $t('gtd.editTask') }}</h5>
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
                  <option value="A">{{ $t('gtd.aLevel') }}</option>
                  <option value="B+">{{ $t('gtd.bPlusLevel') }}</option>
                  <option value="B">{{ $t('gtd.bLevel') }}</option>
                  <option value="C">{{ $t('gtd.cLevel') }}</option>
                  <option value="D">{{ $t('gtd.dLevel') }}</option>
                </select>
              </div>
              <div class="mb-3">
                <label class="form-label">{{ $t('gtd.estimatedTime') }}</label>
                <input type="number" class="form-control" v-model="editForm.estimated_time" min="0" max="1440">
              </div>
              <div class="mb-3">
                <label class="form-label">{{ $t('gtd.deadline') }}</label>
                <input type="datetime-local" class="form-control" v-model="editForm.deadline">
              </div>
              <div class="mb-3">
                <label class="form-label">{{ $t('gtd.context') }}</label>
                <input type="text" class="form-control" v-model="editForm.context" placeholder="@home, @office, @phone">
              </div>
              <div class="mb-3">
                <label class="form-label">{{ $t('gtd.notes') }}</label>
                <textarea class="form-control" v-model="editForm.notes" rows="3"></textarea>
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
    <div v-if="editingTask" class="modal-backdrop fade show"></div>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted, watch } from 'vue';
import { useI18n } from 'vue-i18n';
import { useDailyStore } from '@/stores/dailyStore';
import type { ChecklistItem } from '@/types/gtd';
import { formatDate, formatDateTime } from '@/utils/dateUtils';

const { t, locale } = useI18n();
const emit = defineEmits<{
  (e: 'refresh-stats'): void;
}>();

const dailyStore = useDailyStore();

// UI State
const loading = ref(false);
const selectedDate = ref(new Date());
const dailyStats = ref<any>(null);
const selectedTask = ref<ChecklistItem | null>(null);
const editingTask = ref<ChecklistItem | null>(null);
const actualTimeInput = ref(0);

// Quick Add State
const quickAddTitle = ref('');
const quickAddPriority = ref('B');
const quickAddEstimatedTime = ref(30);

// Edit State
const editForm = ref({
  title: '',
  description: '',
  priority: 'B' as 'A' | 'B+' | 'B' | 'C' | 'D',
  estimated_time: 30,
  deadline: '',
  context: '',
  notes: ''
});

// Computed
const items = computed(() => dailyStore.items);

const priorityOrder = computed(() => [
  { value: 'A', label: t('gtd.aLevel'), icon: 'bi bi-exclamation-triangle' },
  { value: 'B+', label: t('gtd.bPlusLevel'), icon: 'bi bi-triangle' },
  { value: 'B', label: t('gtd.bLevel'), icon: 'bi bi-circle-fill' },
  { value: 'C', label: t('gtd.cLevel'), icon: 'bi bi-circle' },
  { value: 'D', label: t('gtd.dLevel'), icon: 'bi bi-dash-circle' }
]);

// Methods
const loadItems = async () => {
  try {
    loading.value = true;
    await dailyStore.fetchItems({
      date: selectedDate.value
    });
    await loadStats();
  } catch (error) {
    console.error('Failed to load items:', error);
  } finally {
    loading.value = false;
  }
};

const loadStats = async () => {
  try {
    dailyStats.value = await dailyStore.getStats(selectedDate.value);
  } catch (error) {
    console.error('Failed to load stats:', error);
  }
};

const changeDate = (days: number) => {
  const newDate = new Date(selectedDate.value);
  newDate.setDate(newDate.getDate() + days);
  selectedDate.value = newDate;
};

const formatDateHeader = (date: Date): string => {
  return date.toLocaleDateString(locale.value, {
    year: 'numeric',
    month: 'long',
    day: 'numeric'
  });
};

const getDayOfWeek = (date: Date): string => {
  const days = [
    t('gtd.sunday'), t('gtd.monday'), t('gtd.tuesday'), t('gtd.wednesday'), 
    t('gtd.thursday'), t('gtd.friday'), t('gtd.saturday')
  ];
  return days[date.getDay()];
};

const formatTime = (minutes: number | Date): string => {
  if (minutes instanceof Date) {
    return minutes.toLocaleTimeString(locale.value, {
      hour: '2-digit',
      minute: '2-digit'
    });
  }
  if (minutes < 60) {
    return `${minutes}m`;
  }
  const hours = Math.floor(minutes / 60);
  const remainingMinutes = minutes % 60;
  return remainingMinutes > 0 ? `${hours}h ${remainingMinutes}m` : `${hours}h`;
};

const getPriorityCount = (priority: string): number => {
  return items.value.filter(item => item.priority === priority).length;
};

const getPriorityProgress = (priority: string): number => {
  const priorityItems = items.value.filter(item => item.priority === priority);
  if (priorityItems.length === 0) return 0;
  const completed = priorityItems.filter(item => item.status === 'completed').length;
  return (completed / priorityItems.length) * 100;
};

const getTasksByPriority = (priority: string): ChecklistItem[] => {
  return items.value
    .filter(item => item.priority === priority)
    .sort((a, b) => {
      // Sort by status: completed last, then by creation time
      if (a.status === 'completed' && b.status !== 'completed') return 1;
      if (b.status === 'completed' && a.status !== 'completed') return -1;
      return new Date(a.created_at).getTime() - new Date(b.created_at).getTime();
    });
};

const handleQuickAdd = async () => {
  if (!quickAddTitle.value.trim()) return;
  
  try {
    const newTask: ChecklistItem = {
      id: crypto.randomUUID(),
      title: quickAddTitle.value.trim(),
      description: '',
      priority: quickAddPriority.value as 'A' | 'B+' | 'B' | 'C' | 'D',
      estimated_time: quickAddEstimatedTime.value,
      deadline: undefined,
      context: '',
      status: 'pending',
      completion_time: undefined,
      actual_time: 0,
      notes: '',
      inbox_item_id: undefined,
      date: selectedDate.value,
      created_by: '',
      created_at: new Date(),
      updated_at: new Date()
    };
    
    await dailyStore.addItem(newTask);
    
    // Reset form
    quickAddTitle.value = '';
    quickAddPriority.value = 'B';
    quickAddEstimatedTime.value = 30;
    
    emit('refresh-stats');
  } catch (error) {
    console.error('Failed to add task:', error);
  }
};

const selectTask = (task: ChecklistItem) => {
  selectedTask.value = task;
  actualTimeInput.value = task.actual_time;
};

const closeTaskModal = () => {
  selectedTask.value = null;
};

const editTask = (task: ChecklistItem) => {
  editingTask.value = task;
  editForm.value = {
    title: task.title,
    description: task.description,
    priority: task.priority,
    estimated_time: task.estimated_time,
    deadline: task.deadline ? task.deadline.toISOString().slice(0, 16) : '',
    context: task.context,
    notes: task.notes
  };
  closeTaskModal();
};

const closeEditModal = () => {
  editingTask.value = null;
};

const handleEditSubmit = async () => {
  if (!editingTask.value) return;
  
  try {
    const updates = {
      ...editForm.value,
      deadline: editForm.value.deadline ? new Date(editForm.value.deadline) : undefined
    };
    
    await dailyStore.updateItem(editingTask.value.id, updates);
    closeEditModal();
    emit('refresh-stats');
  } catch (error) {
    console.error('Failed to update task:', error);
  }
};

const toggleTaskStatus = async (task: ChecklistItem) => {
  try {
    const newStatus = task.status === 'completed' ? 'pending' : 'completed';
    await dailyStore.updateStatus(task.id, newStatus);
    emit('refresh-stats');
  } catch (error) {
    console.error('Failed to toggle task status:', error);
  }
};

const updateActualTime = async () => {
  if (!selectedTask.value) return;
  
  try {
    await dailyStore.updateActualTime(selectedTask.value.id, actualTimeInput.value);
    selectedTask.value.actual_time = actualTimeInput.value;
    emit('refresh-stats');
  } catch (error) {
    console.error('Failed to update actual time:', error);
  }
};

const deleteTask = async (id: string) => {
  if (!confirm(t('gtd.confirmDeleteTask'))) return;
  
  try {
    await dailyStore.deleteItem(id);
    emit('refresh-stats');
  } catch (error) {
    console.error('Failed to delete task:', error);
  }
};

// Helper functions
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
    in_progress: t('gtd.inProgress'),
    completed: t('gtd.completed'),
    cancelled: t('gtd.cancelled')
  };
  return statusMap[status] || status;
};

const getDeadlineClass = (deadline: Date) => {
  const now = new Date();
  const diffHours = (deadline.getTime() - now.getTime()) / (1000 * 60 * 60);
  
  if (diffHours < 0) return 'overdue';
  if (diffHours < 2) return 'urgent';
  if (diffHours < 24) return 'soon';
  return 'normal';
};

// Lifecycle
onMounted(() => {
  loadItems();
});

// Watch for date changes
watch(selectedDate, () => {
  loadItems();
});
</script>

<style scoped>
/* Daily Panel Styles */
.daily-panel {
  padding: 0;
}

/* Daily Header */
.daily-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 2rem;
  padding: 1.5rem;
  background: linear-gradient(135deg, #667eea, #764ba2);
  border-radius: 16px;
  color: white;
}

.date-selector {
  display: flex;
  align-items: center;
  gap: 1rem;
}

.date-display h3 {
  font-size: 1.5rem;
  font-weight: 600;
  margin-bottom: 0.25rem;
}

.date-subtitle {
  opacity: 0.9;
  margin: 0;
  font-size: 0.9rem;
}

.daily-stats {
  display: flex;
  gap: 1.5rem;
}

.stat-item {
  text-align: center;
}

.stat-value {
  font-size: 1.2rem;
  font-weight: 600;
  margin-bottom: 0.25rem;
}

.stat-label {
  font-size: 0.8rem;
  opacity: 0.9;
}

/* Quick Add Section */
.quick-add-section {
  margin-bottom: 2rem;
}

.quick-add-card {
  background: white;
  border-radius: 12px;
  padding: 1.5rem;
  box-shadow: 0 2px 10px rgba(0, 0, 0, 0.05);
  border: 1px solid #e2e8f0;
}

.quick-add-header h4 {
  font-size: 1.1rem;
  font-weight: 600;
  color: #2d3748;
  margin-bottom: 1rem;
  display: flex;
  align-items: center;
  gap: 0.5rem;
}

.form-row {
  display: grid;
  grid-template-columns: 2fr 1fr 1fr auto;
  gap: 1rem;
  align-items: end;
}

.form-group {
  display: flex;
  flex-direction: column;
}

.form-group label {
  font-size: 0.875rem;
  font-weight: 500;
  color: #4a5568;
  margin-bottom: 0.25rem;
}

/* Priority Sections */
.priority-sections {
  display: flex;
  flex-direction: column;
  gap: 2rem;
}

.priority-section {
  background: white;
  border-radius: 12px;
  padding: 1.5rem;
  box-shadow: 0 2px 10px rgba(0, 0, 0, 0.05);
  border-left: 4px solid;
}

.priority-section.A { border-left-color: #e53e3e; }
.priority-section.B\\+ { border-left-color: #ed8936; }
.priority-section.B { border-left-color: #4299e1; }
.priority-section.C { border-left-color: #48bb78; }
.priority-section.D { border-left-color: #a0aec0; }

.priority-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 1rem;
}

.priority-header h4 {
  font-size: 1.1rem;
  font-weight: 600;
  color: #2d3748;
  display: flex;
  align-items: center;
  gap: 0.5rem;
}

.priority-count {
  font-size: 0.9rem;
  color: #718096;
  font-weight: normal;
}

.priority-progress {
  display: flex;
  align-items: center;
  gap: 0.5rem;
}

.progress-bar {
  width: 100px;
  height: 8px;
  background: #e2e8f0;
  border-radius: 4px;
  overflow: hidden;
}

.progress-fill {
  height: 100%;
  background: linear-gradient(90deg, #667eea, #764ba2);
  transition: width 0.3s ease;
}

.progress-text {
  font-size: 0.8rem;
  color: #4a5568;
  font-weight: 500;
}

/* Tasks List */
.tasks-list {
  display: flex;
  flex-direction: column;
  gap: 1rem;
}

.task-card {
  background: #f8f9fa;
  border-radius: 8px;
  padding: 1rem;
  border: 2px solid transparent;
  transition: all 0.2s ease;
  cursor: pointer;
}

.task-card:hover {
  background: #e2e8f0;
  border-color: #cbd5e0;
}

.task-card.completed {
  background: #f0fff4;
  border-color: #9ae6b4;
  opacity: 0.8;
}

.task-card.completed .task-title h5 {
  text-decoration: line-through;
  color: #68d391;
}

.task-header {
  display: flex;
  align-items: flex-start;
  gap: 1rem;
  margin-bottom: 0.75rem;
}

.task-checkbox {
  margin-top: 0.25rem;
}

.task-checkbox input[type="checkbox"] {
  width: 18px;
  height: 18px;
  accent-color: #667eea;
}

.task-title {
  flex: 1;
}

.task-title h5 {
  font-size: 1rem;
  font-weight: 600;
  color: #2d3748;
  margin-bottom: 0.25rem;
  line-height: 1.4;
}

.task-description {
  font-size: 0.875rem;
  color: #718096;
  margin: 0;
  line-height: 1.4;
}

.task-actions {
  display: flex;
  gap: 0.25rem;
}

.task-meta {
  display: flex;
  flex-wrap: wrap;
  gap: 0.75rem;
  margin-bottom: 0.75rem;
}

.task-time {
  display: flex;
  gap: 0.5rem;
}

.time-estimate,
.time-actual {
  font-size: 0.8rem;
  color: #4a5568;
  display: flex;
  align-items: center;
  gap: 0.25rem;
}

.time-actual {
  color: #667eea;
  font-weight: 500;
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

.deadline-badge {
  font-size: 0.75rem;
  padding: 0.25rem 0.5rem;
  border-radius: 4px;
  display: inline-flex;
  align-items: center;
  gap: 0.25rem;
}

.deadline-badge.normal { background: #e2e8f0; color: #4a5568; }
.deadline-badge.soon { background: #fed7d7; color: #c53030; }
.deadline-badge.urgent { background: #f56565; color: white; }
.deadline-badge.overdue { background: #e53e3e; color: white; }

.task-status {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding-top: 0.75rem;
  border-top: 1px solid #e2e8f0;
}

.status-badge {
  font-size: 0.75rem;
  padding: 0.25rem 0.5rem;
  border-radius: 4px;
  font-weight: 500;
}

.status-badge.pending { background: #fed7d7; color: #c53030; }
.status-badge.in_progress { background: #bee3f8; color: #2b6cb0; }
.status-badge.completed { background: #c6f6d5; color: #2f855a; }
.status-badge.cancelled { background: #e2e8f0; color: #4a5568; }

.completion-time {
  font-size: 0.75rem;
  color: #718096;
}

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

/* Modal Styles */
.modal {
  background: rgba(0, 0, 0, 0.5);
}

.modal-dialog {
  margin-top: 2rem;
}

.task-detail {
  padding: 1rem 0;
}

.detail-section {
  margin-bottom: 2rem;
}

.detail-section h6 {
  font-weight: 600;
  color: #2d3748;
  margin-bottom: 1rem;
  padding-bottom: 0.5rem;
  border-bottom: 1px solid #e2e8f0;
}

.detail-grid {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(200px, 1fr));
  gap: 1rem;
}

.detail-item {
  display: flex;
  flex-direction: column;
  gap: 0.25rem;
}

.detail-item label {
  font-size: 0.875rem;
  font-weight: 500;
  color: #4a5568;
}

.priority-badge {
  font-size: 0.75rem;
  padding: 0.25rem 0.5rem;
  border-radius: 4px;
  font-weight: 600;
  display: inline-block;
}

.priority-badge.A { background: #fed7d7; color: #c53030; }
.priority-badge.B\\+ { background: #fbd38d; color: #c05621; }
.priority-badge.B { background: #bee3f8; color: #2b6cb0; }
.priority-badge.C { background: #c6f6d5; color: #2f855a; }
.priority-badge.D { background: #e2e8f0; color: #4a5568; }

.time-tracking {
  padding: 1rem;
  background: #f8f9fa;
  border-radius: 8px;
}

.time-input-group {
  display: flex;
  flex-direction: column;
  gap: 0.5rem;
}

.time-input-group label {
  font-size: 0.875rem;
  font-weight: 500;
  color: #4a5568;
}

/* Responsive Design */
@media (max-width: 768px) {
  .daily-header {
    flex-direction: column;
    gap: 1rem;
    text-align: center;
  }
  
  .daily-stats {
    justify-content: center;
    flex-wrap: wrap;
  }
  
  .form-row {
    grid-template-columns: 1fr;
    gap: 0.75rem;
  }
  
  .priority-header {
    flex-direction: column;
    align-items: flex-start;
    gap: 0.5rem;
  }
  
  .task-header {
    flex-direction: column;
    gap: 0.5rem;
  }
  
  .task-actions {
    align-self: flex-end;
  }
  
  .task-meta {
    flex-direction: column;
    gap: 0.5rem;
  }
  
  .detail-grid {
    grid-template-columns: 1fr;
  }
}

@media (max-width: 576px) {
  .daily-header {
    padding: 1rem;
  }
  
  .quick-add-card {
    padding: 1rem;
  }
  
  .priority-section {
    padding: 1rem;
  }
  
  .task-card {
    padding: 0.75rem;
  }
}
</style>
