`<template>
  <div class="container mt-4">
    <h1>Task Management</h1>
    
    <div class="mb-4">
      <button class="btn btn-primary" @click="showForm = true" v-if="!showForm">
        Add New Task
      </button>
    </div>

    <div v-if="showForm" class="card mb-4">
      <div class="card-body">
        <h2>{{ editingTask ? 'Edit Task' : 'Add New Task' }}</h2>
        <TaskForm
          :task="editingTask"
          :submit-button-text="editingTask ? 'Update Task' : 'Add Task'"
          @submit="handleSubmit"
          @cancel="cancelForm"
        />
      </div>
    </div>

    <TaskList
      :tasks="tasks"
      :search-query="searchQuery"
      @update:search-query="updateSearchQuery"
      @edit="editTask"
      @delete="deleteTask"
    />
  </div>
</template>

<script setup lang="ts">
import { ref } from 'vue';
import { useTaskStore } from '@/stores/taskStore';
import TaskForm from '@/components/tasks/TaskForm.vue';
import TaskList from '@/components/tasks/TaskList.vue';
import type { Task } from '@/types';

const taskStore = useTaskStore();
const showForm = ref(false);
const editingTask = ref<Task | undefined>();

const { tasks, searchQuery } = taskStore;

const updateSearchQuery = (query: string) => {
  taskStore.searchQuery = query;
};

const handleSubmit = (task: Task) => {
  if (editingTask.value) {
    taskStore.updateTask(editingTask.value.id, task);
  } else {
    taskStore.addTask(task);
  }
  cancelForm();
};

const editTask = (task: Task) => {
  editingTask.value = task;
  showForm.value = true;
};

const deleteTask = (id: string) => {
  if (confirm('Are you sure you want to delete this task?')) {
    taskStore.deleteTask(id);
  }
};

const cancelForm = () => {
  showForm.value = false;
  editingTask.value = undefined;
};
</script>`