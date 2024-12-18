import { defineStore } from 'pinia';
import { ref } from 'vue';
import type { Task } from '@/types';

export const useTaskStore = defineStore('task', () => {
  const tasks = ref<Task[]>([]);
  const searchQuery = ref('');

  const addTask = (task: Task) => {
    tasks.value.push(task);
  };

  const updateTask = (id: string, updatedTask: Task) => {
    const index = tasks.value.findIndex(task => task.id === id);
    if (index !== -1) {
      tasks.value[index] = updatedTask;
    }
  };

  const deleteTask = (id: string) => {
    tasks.value = tasks.value.filter(task => task.id !== id);
  };

  const searchTasks = (query: string) => {
    return tasks.value.filter(task => 
      task.name.toLowerCase().includes(query.toLowerCase()) ||
      task.description.toLowerCase().includes(query.toLowerCase()) ||
      task.tags.some(tag => tag.toLowerCase().includes(query.toLowerCase()))
    );
  };

  return {
    tasks,
    searchQuery,
    addTask,
    updateTask,
    deleteTask,
    searchTasks
  };
});