`<template>
  <div class="task-list">
    <!-- Search and Filter Section -->
    <div class="list-header">
      <div class="search-wrapper">
        <SearchInput
          :model-value="searchQuery"
          @update:model-value="$emit('update:searchQuery', $event)"
          placeholder="Search tasks by name, description, or tags..."
          class="search-input"
        />
      </div>
      <div class="list-stats">
        <span class="text-muted">
          <i class="bi bi-list-check"></i>
          {{ filteredTasks.length }} {{ filteredTasks.length === 1 ? 'task' : 'tasks' }}
        </span>
      </div>
    </div>

    <!-- Tasks Grid -->
    <div class="tasks-grid">
      <TransitionGroup name="task-card" tag="div" class="row g-4">
        <div 
          v-for="task in filteredTasks" 
          :key="task.id"
          class="col-12 col-md-6 col-xl-4"
        >
          <div class="task-card">
            <!-- Card Header -->
            <div class="card-header">
              <h5 class="task-name">{{ task.name }}</h5>
              <div class="task-badges">
                <TaskPriorityBadge :priority="task.priority" />
                <span 
                  class="badge difficulty-badge"
                  :class="`difficulty-${task.difficulty}`"
                >
                  <i :class="getDifficultyIcon(task.difficulty)"></i>
                  {{ task.difficulty }}
                </span>
              </div>
            </div>

            <!-- Card Body -->
            <div class="card-body">
              <p class="task-description">{{ task.description }}</p>
              
              <!-- Task Details Grid -->
              <div class="task-details">
                <div class="detail-item">
                  <i class="bi bi-clock text-muted"></i>
                  <span class="detail-label">Duration:</span>
                  <span class="detail-value">{{ task.minutes }} min</span>
                </div>
                
                <div class="detail-item">
                  <i class="bi bi-calendar-week text-primary"></i>
                  <span class="detail-label">Scheduled:</span>
                  <span class="detail-value">{{ formatDate(task.schedule_time) }}</span>
                </div>
                
                <div class="detail-item">
                  <i class="bi bi-calendar-check text-danger"></i>
                  <span class="detail-label">Deadline:</span>
                  <span class="detail-value">{{ formatDate(task.deadline) }}</span>
                </div>
                
                <div v-if="task.start_time" class="detail-item">
                  <i class="bi bi-play-circle text-success"></i>
                  <span class="detail-label">Started:</span>
                  <span class="detail-value">{{ formatDate(task.start_time) }}</span>
                </div>
                
                <div v-if="task.end_time" class="detail-item">
                  <i class="bi bi-stop-circle text-secondary"></i>
                  <span class="detail-label">Ended:</span>
                  <span class="detail-value">{{ formatDate(task.end_time) }}</span>
                </div>
              </div>

              <!-- Tags -->
              <div v-if="task.tags.length > 0" class="task-tags">
                <TaskTags :tags="task.tags" />
              </div>
            </div>

            <!-- Card Footer -->
            <div class="card-footer">
              <button
                class="btn btn-sm btn-outline-primary"
                @click="$emit('edit', task)"
              >
                <i class="bi bi-pencil me-1"></i>
                Edit
              </button>
              <button
                class="btn btn-sm btn-outline-danger"
                @click="handleDelete(task.id)"
              >
                <i class="bi bi-trash me-1"></i>
                Delete
              </button>
            </div>
          </div>
        </div>
      </TransitionGroup>
    </div>

    <!-- Empty State -->
    <div v-if="filteredTasks.length === 0" class="empty-state">
      <i class="bi bi-inbox"></i>
      <h4>No tasks found</h4>
      <p class="text-muted">
        {{ searchQuery ? 'Try adjusting your search criteria' : 'Create your first task to get started' }}
      </p>
    </div>
  </div>
</template>

<script setup lang="ts">
import { computed } from 'vue';
import type { Task } from '@/types';
import { formatDate } from '@/utils/dateUtils';
import { useConfirmDialog } from '@/components/common/ConfirmDialog';
import SearchInput from '@/components/common/SearchInput.vue';
import TaskPriorityBadge from './TaskPriorityBadge.vue';
import TaskTags from './TaskTags.vue';

const props = defineProps<{
  tasks: Task[];
  searchQuery: string;
}>();

const emit = defineEmits<{
  (e: 'edit', task: Task): void;
  (e: 'delete', id: string): void;
  (e: 'update:searchQuery', value: string): void;
}>();

const { confirm } = useConfirmDialog();

const filteredTasks = computed(() => {
  if (!props.searchQuery) return props.tasks;
  
  return props.tasks.filter(task => 
    task.name.toLowerCase().includes(props.searchQuery.toLowerCase()) ||
    task.description.toLowerCase().includes(props.searchQuery.toLowerCase()) ||
    task.tags.some(tag => tag.toLowerCase().includes(props.searchQuery.toLowerCase()))
  );
});

const handleDelete = async (id: string) => {
  if (await confirm('Are you sure you want to delete this task?')) {
    emit('delete', id);
  }
};

const getDifficultyIcon = (difficulty: string) => {
  switch (difficulty) {
    case 'easy':
      return 'bi bi-emoji-smile';
    case 'medium':
      return 'bi bi-emoji-neutral';
    case 'hard':
      return 'bi bi-emoji-dizzy';
    default:
      return 'bi bi-emoji-neutral';
  }
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

.search-wrapper {
  flex: 1;
  max-width: 500px;
}

.search-input {
  width: 100%;
}

.list-stats {
  font-size: 0.9rem;
}

/* Tasks Grid */
.tasks-grid {
  min-height: 200px;
}

/* Task Card */
.task-card {
  background: white;
  border-radius: 12px;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.08);
  transition: all 0.3s ease;
  height: 100%;
  display: flex;
  flex-direction: column;
  overflow: hidden;
}

.task-card:hover {
  transform: translateY(-4px);
  box-shadow: 0 8px 24px rgba(0, 0, 0, 0.12);
}

/* Card Header */
.card-header {
  padding: 1.25rem;
  background: linear-gradient(135deg, #f8f9fa 0%, #e9ecef 100%);
  border-bottom: 1px solid #dee2e6;
}

.task-name {
  font-size: 1.1rem;
  font-weight: 600;
  margin: 0 0 0.5rem 0;
  color: #495057;
  display: -webkit-box;
  line-clamp: 2;
  -webkit-line-clamp: 2;
  -webkit-box-orient: vertical;
  overflow: hidden;
}

.task-badges {
  display: flex;
  gap: 0.5rem;
  flex-wrap: wrap;
}

/* Difficulty Badge */
.difficulty-badge {
  font-size: 0.75rem;
  padding: 0.25rem 0.5rem;
  font-weight: 500;
  display: inline-flex;
  align-items: center;
  gap: 0.25rem;
}

.difficulty-easy {
  background-color: #cff4fc;
  color: #055160;
}

.difficulty-medium {
  background-color: #e0cffc;
  color: #432874;
}

.difficulty-hard {
  background-color: #f8d7da;
  color: #842029;
}

/* Card Body */
.card-body {
  padding: 1.25rem;
  flex: 1;
  display: flex;
  flex-direction: column;
}

.task-description {
  color: #6c757d;
  font-size: 0.9rem;
  margin-bottom: 1rem;
  display: -webkit-box;
  line-clamp: 3;
  -webkit-line-clamp: 3;
  -webkit-box-orient: vertical;
  overflow: hidden;
}

/* Task Details */
.task-details {
  display: flex;
  flex-direction: column;
  gap: 0.5rem;
  margin-bottom: 1rem;
}

.detail-item {
  display: flex;
  align-items: center;
  gap: 0.5rem;
  font-size: 0.875rem;
}

.detail-item i {
  width: 20px;
  text-align: center;
}

.detail-label {
  color: #6c757d;
  min-width: 70px;
}

.detail-value {
  color: #495057;
  font-weight: 500;
}

/* Task Tags */
.task-tags {
  margin-top: auto;
  padding-top: 0.5rem;
}

/* Card Footer */
.card-footer {
  padding: 1rem 1.25rem;
  background-color: #f8f9fa;
  border-top: 1px solid #dee2e6;
  display: flex;
  gap: 0.5rem;
  justify-content: flex-end;
}

.card-footer .btn {
  font-size: 0.875rem;
  padding: 0.375rem 0.75rem;
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

/* Animations */
.task-card-enter-active,
.task-card-leave-active {
  transition: all 0.3s ease;
}

.task-card-enter-from {
  opacity: 0;
  transform: scale(0.9);
}

.task-card-leave-to {
  opacity: 0;
  transform: scale(0.9);
}

.task-card-move {
  transition: transform 0.3s ease;
}

/* Responsive Design */
@media (max-width: 768px) {
  .list-header {
    flex-direction: column;
    align-items: stretch;
  }
  
  .search-wrapper {
    max-width: 100%;
  }
  
  .list-stats {
    text-align: center;
  }
  
  .task-card {
    margin-bottom: 1rem;
  }
}
</style>`