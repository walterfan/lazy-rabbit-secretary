`<template>
  <div class="task-list">
    <SearchInput
      :model-value="searchQuery"
      @update:model-value="$emit('update:searchQuery', $event)"
      placeholder="Search tasks..."
    />

    <div class="table-responsive">
      <table class="table table-striped">
        <thead>
          <tr>
            <th>Name</th>
            <th>Description</th>
            <th>Priority</th>
            <th>Duration</th>
            <th>Deadline</th>
            <th>Start Time</th>
            <th>End Time</th>
            <th>Tags</th>
            <th>Actions</th>
          </tr>
        </thead>
        <tbody>
          <tr v-for="task in filteredTasks" :key="task.id">
            <td>{{ task.name }}</td>
            <td>{{ task.description }}</td>
            <td>
              <TaskPriorityBadge :priority="task.priority" />
            </td>
            <td>{{ task.minutes }} min</td>
            <td>{{ formatDate(task.deadline) }}</td>
            <td>{{ formatDate(task.start_time) }}</td>
            <td>{{ formatDate(task.end_time) }}</td>
            <td>
              <TaskTags :tags="task.tags" />
            </td>
            <td>
              <button
                class="btn btn-sm btn-primary me-2"
                @click="$emit('edit', task)"
              >
                Edit
              </button>
              <button
                class="btn btn-sm btn-danger"
                @click="handleDelete(task.id)"
              >
                Delete
              </button>
            </td>
          </tr>
        </tbody>
      </table>
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
</script>`