import { defineStore } from 'pinia';
import { ref } from 'vue';
import type { Task } from '@/types';
import { handleHttpError, showErrorAlert, logError } from '@/utils/errorHandler';
import { makeAuthenticatedRequest } from '@/utils/httpInterceptor';

// Map frontend priority/difficulty strings to backend integers
const priorityMap = {
  'low': 1,
  'medium': 2,
  'high': 3
};

const difficultyMap = {
  'easy': 1,
  'medium': 2,
  'hard': 3
};

// Reverse maps for converting backend to frontend
const priorityReverseMap = {
  1: 'low',
  2: 'medium',
  3: 'high'
} as const;

const difficultyReverseMap = {
  1: 'easy',
  2: 'medium',
  3: 'hard'
} as const;

interface CreateTaskRequest {
  name: string;
  description: string;
  priority?: number;
  difficulty?: number;
  schedule_time: string;
  minutes: number;
  deadline: string;
  tags: string;
}

interface UpdateTaskRequest {
  name?: string;
  description?: string;
  priority?: number;
  difficulty?: number;
  status?: string;
  schedule_time?: string;
  minutes?: number;
  deadline?: string;
  tags?: string;
}

export const useTaskStore = defineStore('task', () => {
  const tasks = ref<Task[]>([]);
  const totalCount = ref(0);
  const loading = ref(false);
  const error = ref<string | null>(null);
  const searchQuery = ref('');
  const currentPage = ref(1);
  const pageSize = ref(20);

  // Fetch tasks with search and filters
  const fetchTasks = async (params: {
    realm_id?: string;
    q?: string;
    status?: string;
    tags?: string;
    priority?: number;
    difficulty?: number;
    page?: number;
    page_size?: number;
  } = {}) => {
    loading.value = true;
    error.value = null;
    
    try {
      const queryParams = new URLSearchParams();
      if (params.realm_id) queryParams.append('realm_id', params.realm_id);
      if (params.q || searchQuery.value) queryParams.append('q', params.q || searchQuery.value);
      if (params.status) queryParams.append('status', params.status);
      if (params.tags) queryParams.append('tags', params.tags);
      if (params.priority) queryParams.append('priority', params.priority.toString());
      if (params.difficulty) queryParams.append('difficulty', params.difficulty.toString());
      queryParams.append('page', (params.page || currentPage.value).toString());
      queryParams.append('page_size', (params.page_size || pageSize.value).toString());
      
      const response = await makeAuthenticatedRequest(`/api/v1/tasks?${queryParams}`);
      
      if (!response.ok) {
        const apiError = await handleHttpError(response);
        logError(apiError, 'fetchTasks');
        showErrorAlert(apiError);
        throw new Error(apiError.message);
      }
      
      const data = await response.json();
      
      // Convert backend task format to frontend format
      tasks.value = (data.items || []).map((task: any) => ({
        ...task,
        priority: priorityReverseMap[task.priority as keyof typeof priorityReverseMap] || 'medium',
        difficulty: difficultyReverseMap[task.difficulty as keyof typeof difficultyReverseMap] || 'medium',
        deadline: new Date(task.deadline),
        schedule_time: new Date(task.schedule_time),
        start_time: task.start_time ? new Date(task.start_time) : undefined,
        end_time: task.end_time ? new Date(task.end_time) : undefined,
        tags: task.tags ? task.tags.split(',').map((t: string) => t.trim()) : []
      }));
      
      totalCount.value = data.total || 0;
      currentPage.value = data.page || 1;
    } catch (err) {
      error.value = err instanceof Error ? err.message : 'An error occurred';
      throw err;
    } finally {
      loading.value = false;
    }
  };

  // Get a single task
  const getTask = async (id: string): Promise<Task> => {
    loading.value = true;
    error.value = null;
    
    try {
      const response = await makeAuthenticatedRequest(`/api/v1/tasks/${id}`);
      
      if (!response.ok) {
        const apiError = await handleHttpError(response);
        logError(apiError, 'getTask');
        showErrorAlert(apiError);
        throw new Error(apiError.message);
      }
      
      const task = await response.json();
      
      // Convert backend format to frontend format
      return {
        ...task,
        priority: priorityReverseMap[task.priority as keyof typeof priorityReverseMap] || 'medium',
        difficulty: difficultyReverseMap[task.difficulty as keyof typeof difficultyReverseMap] || 'medium',
        deadline: new Date(task.deadline),
        schedule_time: new Date(task.schedule_time),
        start_time: task.start_time ? new Date(task.start_time) : undefined,
        end_time: task.end_time ? new Date(task.end_time) : undefined,
        tags: task.tags ? task.tags.split(',').map((t: string) => t.trim()) : []
      };
    } catch (err) {
      error.value = err instanceof Error ? err.message : 'An error occurred';
      throw err;
    } finally {
      loading.value = false;
    }
  };

  // Create a new task
  const addTask = async (task: Task): Promise<Task> => {
    loading.value = true;
    error.value = null;
    
    try {
      const requestData: CreateTaskRequest = {
        name: task.name,
        description: task.description,
        priority: priorityMap[task.priority],
        difficulty: difficultyMap[task.difficulty],
        schedule_time: task.schedule_time.toISOString(),
        minutes: task.minutes,
        deadline: task.deadline.toISOString(),
        tags: task.tags.join(',')
      };
      
      const response = await makeAuthenticatedRequest('/api/v1/tasks', {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json'
        },
        body: JSON.stringify(requestData)
      });
      
      if (!response.ok) {
        const apiError = await handleHttpError(response);
        logError(apiError, 'addTask');
        showErrorAlert(apiError);
        throw new Error(apiError.message);
      }
      
      const newTask = await response.json();
      
      // Convert and add to local store
      const convertedTask = {
        ...newTask,
        priority: priorityReverseMap[newTask.priority as keyof typeof priorityReverseMap] || 'medium',
        difficulty: difficultyReverseMap[newTask.difficulty as keyof typeof difficultyReverseMap] || 'medium',
        deadline: new Date(newTask.deadline),
        schedule_time: new Date(newTask.schedule_time),
        start_time: newTask.start_time ? new Date(newTask.start_time) : undefined,
        end_time: newTask.end_time ? new Date(newTask.end_time) : undefined,
        tags: newTask.tags ? newTask.tags.split(',').map((t: string) => t.trim()) : []
      };
      
      tasks.value.push(convertedTask);
      totalCount.value += 1;
      return convertedTask;
    } catch (err) {
      error.value = err instanceof Error ? err.message : 'An error occurred';
      throw err;
    } finally {
      loading.value = false;
    }
  };

  // Update an existing task
  const updateTask = async (id: string, updates: Partial<Task>): Promise<Task> => {
    loading.value = true;
    error.value = null;
    
    try {
      const requestData: UpdateTaskRequest = {};
      
      if (updates.name !== undefined) requestData.name = updates.name;
      if (updates.description !== undefined) requestData.description = updates.description;
      if (updates.priority !== undefined) requestData.priority = priorityMap[updates.priority];
      if (updates.difficulty !== undefined) requestData.difficulty = difficultyMap[updates.difficulty];
      if (updates.schedule_time !== undefined) requestData.schedule_time = updates.schedule_time.toISOString();
      if (updates.minutes !== undefined) requestData.minutes = updates.minutes;
      if (updates.deadline !== undefined) requestData.deadline = updates.deadline.toISOString();
      if (updates.tags !== undefined) requestData.tags = updates.tags.join(',');
      
      const response = await makeAuthenticatedRequest(`/api/v1/tasks/${id}`, {
        method: 'PUT',
        headers: {
          'Content-Type': 'application/json'
        },
        body: JSON.stringify(requestData)
      });
      
      if (!response.ok) {
        const apiError = await handleHttpError(response);
        logError(apiError, 'updateTask');
        showErrorAlert(apiError);
        throw new Error(apiError.message);
      }
      
      const updatedTask = await response.json();
      
      // Convert and update in local store
      const convertedTask = {
        ...updatedTask,
        priority: priorityReverseMap[updatedTask.priority as keyof typeof priorityReverseMap] || 'medium',
        difficulty: difficultyReverseMap[updatedTask.difficulty as keyof typeof difficultyReverseMap] || 'medium',
        deadline: new Date(updatedTask.deadline),
        schedule_time: new Date(updatedTask.schedule_time),
        start_time: updatedTask.start_time ? new Date(updatedTask.start_time) : undefined,
        end_time: updatedTask.end_time ? new Date(updatedTask.end_time) : undefined,
        tags: updatedTask.tags ? updatedTask.tags.split(',').map((t: string) => t.trim()) : []
      };
      
      const index = tasks.value.findIndex(t => t.id === id);
      if (index !== -1) {
        tasks.value[index] = convertedTask;
      }
      return convertedTask;
    } catch (err) {
      error.value = err instanceof Error ? err.message : 'An error occurred';
      throw err;
    } finally {
      loading.value = false;
    }
  };

  // Delete a task
  const deleteTask = async (id: string): Promise<void> => {
    loading.value = true;
    error.value = null;
    
    try {
      const response = await makeAuthenticatedRequest(`/api/v1/tasks/${id}`, {
        method: 'DELETE'
      });
      
      if (!response.ok) {
        const apiError = await handleHttpError(response);
        logError(apiError, 'deleteTask');
        showErrorAlert(apiError);
        throw new Error(apiError.message);
      }
      
      tasks.value = tasks.value.filter(t => t.id !== id);
      totalCount.value -= 1;
    } catch (err) {
      error.value = err instanceof Error ? err.message : 'An error occurred';
      throw err;
    } finally {
      loading.value = false;
    }
  };

  // Task status actions
  const startTask = async (id: string): Promise<Task> => {
    const response = await makeAuthenticatedRequest(`/api/v1/tasks/${id}/start`, {
      method: 'POST'
    });
    
    if (!response.ok) {
      const apiError = await handleHttpError(response);
      logError(apiError, 'startTask');
      showErrorAlert(apiError);
      throw new Error(apiError.message);
    }
    
    const updatedTask = await response.json();
    await fetchTasks(); // Refresh the list
    return updatedTask;
  };

  const completeTask = async (id: string): Promise<Task> => {
    const response = await makeAuthenticatedRequest(`/api/v1/tasks/${id}/complete`, {
      method: 'POST'
    });
    
    if (!response.ok) {
      const apiError = await handleHttpError(response);
      logError(apiError, 'completeTask');
      showErrorAlert(apiError);
      throw new Error(apiError.message);
    }
    
    const updatedTask = await response.json();
    await fetchTasks(); // Refresh the list
    return updatedTask;
  };

  const failTask = async (id: string): Promise<Task> => {
    const response = await makeAuthenticatedRequest(`/api/v1/tasks/${id}/fail`, {
      method: 'POST'
    });
    
    if (!response.ok) {
      const apiError = await handleHttpError(response);
      logError(apiError, 'failTask');
      showErrorAlert(apiError);
      throw new Error(apiError.message);
    }
    
    const updatedTask = await response.json();
    await fetchTasks(); // Refresh the list
    return updatedTask;
  };

  // Get upcoming tasks
  const getUpcomingTasks = async (limit: number = 10): Promise<Task[]> => {
    const response = await makeAuthenticatedRequest(`/api/v1/tasks/upcoming?limit=${limit}`);
    
    if (!response.ok) {
      const apiError = await handleHttpError(response);
      logError(apiError, 'getUpcomingTasks');
      showErrorAlert(apiError);
      throw new Error(apiError.message);
    }
    
    const data = await response.json();
    return data.items || [];
  };

  // Get overdue tasks
  const getOverdueTasks = async (limit: number = 10): Promise<Task[]> => {
    const response = await makeAuthenticatedRequest(`/api/v1/tasks/overdue?limit=${limit}`);
    
    if (!response.ok) {
      const apiError = await handleHttpError(response);
      logError(apiError, 'getOverdueTasks');
      showErrorAlert(apiError);
      throw new Error(apiError.message);
    }
    
    const data = await response.json();
    return data.items || [];
  };

  // Clear store
  const clearTasks = () => {
    tasks.value = [];
    totalCount.value = 0;
    error.value = null;
    searchQuery.value = '';
    currentPage.value = 1;
  };

  return {
    tasks,
    totalCount,
    loading,
    error,
    searchQuery,
    currentPage,
    pageSize,
    fetchTasks,
    getTask,
    addTask,
    updateTask,
    deleteTask,
    startTask,
    completeTask,
    failTask,
    getUpcomingTasks,
    getOverdueTasks,
    clearTasks
  };
});