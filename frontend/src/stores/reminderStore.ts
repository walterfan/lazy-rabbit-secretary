import { defineStore } from 'pinia';
import { ref } from 'vue';
import type { Reminder, CreateReminderRequest, UpdateReminderRequest } from '@/types';
import { handleHttpError, showErrorAlert, logError, type ApiError } from '@/utils/errorHandler';
import { makeAuthenticatedRequest } from '@/utils/httpInterceptor';

interface SearchResponse {
  items: Reminder[];
  total: number;
  page: number;
  page_size: number;
  total_pages: number;
}

interface ReminderListResponse {
  items: Reminder[];
}

export const useReminderStore = defineStore('reminder', () => {
  const reminders = ref<Reminder[]>([]);
  const loading = ref(false);
  const error = ref<string | null>(null);

  // Convert backend reminder to frontend format
  const convertFromBackend = (backendReminder: any): Reminder => {
    return {
      ...backendReminder,
      remind_time: new Date(backendReminder.remind_time),
      created_at: new Date(backendReminder.created_at),
      updated_at: new Date(backendReminder.updated_at)
    };
  };

  // Convert frontend reminder to backend format
  const convertToBackend = (frontendReminder: CreateReminderRequest | UpdateReminderRequest): any => {
    const converted = { ...frontendReminder };
    if ('remind_time' in converted && converted.remind_time) {
      (converted as any).remind_time = converted.remind_time.toISOString();
    }
    return converted;
  };

  // Create a new reminder
  const createReminder = async (reminderData: CreateReminderRequest): Promise<Reminder> => {
    loading.value = true;
    error.value = null;

    try {
      const backendData = convertToBackend(reminderData);
      
      const response = await makeAuthenticatedRequest('/api/v1/reminders', {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json'
        },
        body: JSON.stringify(backendData)
      });

      if (!response.ok) {
        const apiError = await handleHttpError(response);
        logError(apiError, 'createReminder');
        showErrorAlert(apiError);
        throw new Error(apiError.message);
      }

      const data = await response.json();
      const reminder = convertFromBackend(data);
      
      // Add to local store
      reminders.value.unshift(reminder);
      
      return reminder;
    } catch (err) {
      const errorMessage = err instanceof Error ? err.message : 'Failed to create reminder';
      error.value = errorMessage;
      logError(err as ApiError, 'createReminder');
      throw err;
    } finally {
      loading.value = false;
    }
  };

  // Fetch reminders with search and pagination
  const fetchReminders = async (params: {
    q?: string;
    status?: string;
    tags?: string;
    page?: number;
    page_size?: number;
  } = {}): Promise<SearchResponse> => {
    loading.value = true;
    error.value = null;

    try {
      const searchParams = new URLSearchParams();
      if (params.q) searchParams.append('q', params.q);
      if (params.status) searchParams.append('status', params.status);
      if (params.tags) searchParams.append('tags', params.tags);
      if (params.page) searchParams.append('page', params.page.toString());
      if (params.page_size) searchParams.append('page_size', params.page_size.toString());

      const response = await makeAuthenticatedRequest(`/api/v1/reminders?${searchParams.toString()}`, {
        method: 'GET'
      });

      if (!response.ok) {
        const apiError = await handleHttpError(response);
        logError(apiError, 'fetchReminders');
        showErrorAlert(apiError);
        throw new Error(apiError.message);
      }

      const data = await response.json();
      const convertedReminders = data.items.map(convertFromBackend);
      
      // Update local store
      reminders.value = convertedReminders;
      
      return {
        ...data,
        items: convertedReminders
      };
    } catch (err) {
      const errorMessage = err instanceof Error ? err.message : 'Failed to fetch reminders';
      error.value = errorMessage;
      logError(err as ApiError, 'fetchReminders');
      throw err;
    } finally {
      loading.value = false;
    }
  };

  // Get a specific reminder by ID
  const getReminder = async (id: string): Promise<Reminder> => {
    loading.value = true;
    error.value = null;

    try {
      const response = await makeAuthenticatedRequest(`/api/v1/reminders/${id}`, {
        method: 'GET'
      });

      if (!response.ok) {
        const apiError = await handleHttpError(response);
        logError(apiError, 'getReminder');
        showErrorAlert(apiError);
        throw new Error(apiError.message);
      }

      const data = await response.json();
      return convertFromBackend(data);
    } catch (err) {
      const errorMessage = err instanceof Error ? err.message : 'Failed to get reminder';
      error.value = errorMessage;
      logError(err as ApiError, 'getReminder');
      throw err;
    } finally {
      loading.value = false;
    }
  };

  // Update a reminder
  const updateReminder = async (id: string, updates: UpdateReminderRequest): Promise<Reminder> => {
    loading.value = true;
    error.value = null;

    try {
      const backendData = convertToBackend(updates);
      
      const response = await makeAuthenticatedRequest(`/api/v1/reminders/${id}`, {
        method: 'PUT',
        headers: {
          'Content-Type': 'application/json'
        },
        body: JSON.stringify(backendData)
      });

      if (!response.ok) {
        const apiError = await handleHttpError(response);
        logError(apiError, 'updateReminder');
        showErrorAlert(apiError);
        throw new Error(apiError.message);
      }

      const data = await response.json();
      const updatedReminder = convertFromBackend(data);
      
      // Update in local store
      const index = reminders.value.findIndex(r => r.id === id);
      if (index !== -1) {
        reminders.value[index] = updatedReminder;
      }
      
      return updatedReminder;
    } catch (err) {
      const errorMessage = err instanceof Error ? err.message : 'Failed to update reminder';
      error.value = errorMessage;
      logError(err as ApiError, 'updateReminder');
      throw err;
    } finally {
      loading.value = false;
    }
  };

  // Delete a reminder
  const deleteReminder = async (id: string): Promise<void> => {
    loading.value = true;
    error.value = null;

    try {
      const response = await makeAuthenticatedRequest(`/api/v1/reminders/${id}`, {
        method: 'DELETE'
      });

      if (!response.ok) {
        const apiError = await handleHttpError(response);
        logError(apiError, 'deleteReminder');
        showErrorAlert(apiError);
        throw new Error(apiError.message);
      }

      // Remove from local store
      const index = reminders.value.findIndex(r => r.id === id);
      if (index !== -1) {
        reminders.value.splice(index, 1);
      }
    } catch (err) {
      const errorMessage = err instanceof Error ? err.message : 'Failed to delete reminder';
      error.value = errorMessage;
      logError(err as ApiError, 'deleteReminder');
      throw err;
    } finally {
      loading.value = false;
    }
  };

  // Get upcoming reminders
  const getUpcomingReminders = async (limit: number = 10): Promise<Reminder[]> => {
    loading.value = true;
    error.value = null;

    try {
      const response = await makeAuthenticatedRequest(`/api/v1/reminders/upcoming?limit=${limit}`, {
        method: 'GET'
      });

      if (!response.ok) {
        const apiError = await handleHttpError(response);
        logError(apiError, 'getUpcomingReminders');
        showErrorAlert(apiError);
        throw new Error(apiError.message);
      }

      const data: ReminderListResponse = await response.json();
      return data.items.map(convertFromBackend);
    } catch (err) {
      const errorMessage = err instanceof Error ? err.message : 'Failed to get upcoming reminders';
      error.value = errorMessage;
      logError(err as ApiError, 'getUpcomingReminders');
      throw err;
    } finally {
      loading.value = false;
    }
  };

  // Get overdue reminders
  const getOverdueReminders = async (limit: number = 10): Promise<Reminder[]> => {
    loading.value = true;
    error.value = null;

    try {
      const response = await makeAuthenticatedRequest(`/api/v1/reminders/overdue?limit=${limit}`, {
        method: 'GET'
      });

      if (!response.ok) {
        const apiError = await handleHttpError(response);
        logError(apiError, 'getOverdueReminders');
        showErrorAlert(apiError);
        throw new Error(apiError.message);
      }

      const data: ReminderListResponse = await response.json();
      return data.items.map(convertFromBackend);
    } catch (err) {
      const errorMessage = err instanceof Error ? err.message : 'Failed to get overdue reminders';
      error.value = errorMessage;
      logError(err as ApiError, 'getOverdueReminders');
      throw err;
    } finally {
      loading.value = false;
    }
  };

  // Mark reminder as completed
  const markAsCompleted = async (id: string): Promise<Reminder> => {
    loading.value = true;
    error.value = null;

    try {
      const response = await makeAuthenticatedRequest(`/api/v1/reminders/${id}/complete`, {
        method: 'POST'
      });

      if (!response.ok) {
        const apiError = await handleHttpError(response);
        logError(apiError, 'markAsCompleted');
        showErrorAlert(apiError);
        throw new Error(apiError.message);
      }

      const data = await response.json();
      const updatedReminder = convertFromBackend(data);
      
      // Update in local store
      const index = reminders.value.findIndex(r => r.id === id);
      if (index !== -1) {
        reminders.value[index] = updatedReminder;
      }
      
      return updatedReminder;
    } catch (err) {
      const errorMessage = err instanceof Error ? err.message : 'Failed to mark reminder as completed';
      error.value = errorMessage;
      logError(err as ApiError, 'markAsCompleted');
      throw err;
    } finally {
      loading.value = false;
    }
  };

  // Snooze reminder
  const snoozeReminder = async (id: string, newRemindTime: Date): Promise<Reminder> => {
    loading.value = true;
    error.value = null;

    try {
      const response = await makeAuthenticatedRequest(`/api/v1/reminders/${id}/snooze`, {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json'
        },
        body: JSON.stringify({
          remind_time: newRemindTime.toISOString()
        })
      });

      if (!response.ok) {
        const apiError = await handleHttpError(response);
        logError(apiError, 'snoozeReminder');
        showErrorAlert(apiError);
        throw new Error(apiError.message);
      }

      const data = await response.json();
      const updatedReminder = convertFromBackend(data);
      
      // Update in local store
      const index = reminders.value.findIndex(r => r.id === id);
      if (index !== -1) {
        reminders.value[index] = updatedReminder;
      }
      
      return updatedReminder;
    } catch (err) {
      const errorMessage = err instanceof Error ? err.message : 'Failed to snooze reminder';
      error.value = errorMessage;
      logError(err as ApiError, 'snoozeReminder');
      throw err;
    } finally {
      loading.value = false;
    }
  };

  // Clear store
  const clearReminders = () => {
    reminders.value = [];
    error.value = null;
  };

  return {
    // State
    reminders,
    loading,
    error,

    // Actions
    createReminder,
    fetchReminders,
    getReminder,
    updateReminder,
    deleteReminder,
    getUpcomingReminders,
    getOverdueReminders,
    markAsCompleted,
    snoozeReminder,
    clearReminders
  };
});
