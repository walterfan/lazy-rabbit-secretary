import { defineStore } from 'pinia';
import { ref } from 'vue';
import type { 
  DailyChecklistItem, 
  CreateDailyItemRequest, 
  UpdateDailyItemRequest, 
  DailyListResponse,
  DailyStatsResponse 
} from '@/types/gtd';
import { handleHttpError, showErrorAlert, logError } from '@/utils/errorHandler';
import { makeAuthenticatedRequest } from '@/utils/httpInterceptor';

export const useDailyStore = defineStore('daily', () => {
  const items = ref<DailyChecklistItem[]>([]);
  const totalCount = ref(0);
  const loading = ref(false);
  const error = ref<string | null>(null);
  const searchQuery = ref('');
  const currentPage = ref(1);
  const pageSize = ref(20);
  const selectedDate = ref(new Date());
  const filters = ref({
    status: '',
    priority: '',
    context: ''
  });

  // Fetch daily items with pagination and filters
  const fetchItems = async (params: {
    page?: number;
    page_size?: number;
    date?: Date;
    status?: string;
    priority?: string;
    context?: string;
    q?: string;
  } = {}) => {
    loading.value = true;
    error.value = null;
    
    try {
      const queryParams = new URLSearchParams();
      queryParams.append('page', (params.page || currentPage.value).toString());
      queryParams.append('page_size', (params.page_size || pageSize.value).toString());
      
      const date = params.date || selectedDate.value;
      queryParams.append('date', date.toISOString().split('T')[0]);
      
      if (params.status || filters.value.status) {
        queryParams.append('status', params.status || filters.value.status);
      }
      if (params.priority || filters.value.priority) {
        queryParams.append('priority', params.priority || filters.value.priority);
      }
      if (params.context || filters.value.context) {
        queryParams.append('context', params.context || filters.value.context);
      }
      if (params.q || searchQuery.value) {
        queryParams.append('q', params.q || searchQuery.value);
      }
      
      const response = await makeAuthenticatedRequest(`/api/v1/daily?${queryParams}`);
      
      if (!response.ok) {
        const apiError = await handleHttpError(response);
        logError(apiError, 'fetchItems');
        showErrorAlert(apiError);
        throw new Error(apiError.message);
      }
      
      const data: DailyListResponse = await response.json();
      
      items.value = data.items.map(item => ({
        ...item,
        date: new Date(item.date),
        deadline: item.deadline ? new Date(item.deadline) : undefined,
        completion_time: item.completion_time ? new Date(item.completion_time) : undefined,
        created_at: new Date(item.created_at),
        updated_at: new Date(item.updated_at)
      }));
      
      totalCount.value = data.total;
      currentPage.value = data.page;
    } catch (err) {
      error.value = err instanceof Error ? err.message : 'An error occurred';
      throw err;
    } finally {
      loading.value = false;
    }
  };

  // Get items for specific date
  const getItemsByDate = async (date: Date) => {
    loading.value = true;
    error.value = null;
    
    try {
      const dateStr = date.toISOString().split('T')[0];
      const response = await makeAuthenticatedRequest(`/api/v1/daily/date/${dateStr}`);
      
      if (!response.ok) {
        const apiError = await handleHttpError(response);
        logError(apiError, 'getItemsByDate');
        showErrorAlert(apiError);
        throw new Error(apiError.message);
      }
      
      const data = await response.json();
      return data.items.map((item: DailyChecklistItem) => ({
        ...item,
        date: new Date(item.date),
        deadline: item.deadline ? new Date(item.deadline) : undefined,
        completion_time: item.completion_time ? new Date(item.completion_time) : undefined,
        created_at: new Date(item.created_at),
        updated_at: new Date(item.updated_at)
      }));
    } catch (err) {
      error.value = err instanceof Error ? err.message : 'An error occurred';
      throw err;
    } finally {
      loading.value = false;
    }
  };

  // Get daily statistics
  const getStats = async (date: Date): Promise<DailyStatsResponse> => {
    loading.value = true;
    error.value = null;
    
    try {
      const dateStr = date.toISOString().split('T')[0];
      const response = await makeAuthenticatedRequest(`/api/v1/daily/stats/${dateStr}`);
      
      if (!response.ok) {
        const apiError = await handleHttpError(response);
        logError(apiError, 'getStats');
        showErrorAlert(apiError);
        throw new Error(apiError.message);
      }
      
      const stats = await response.json();
      return {
        ...stats,
        date: new Date(stats.date)
      };
    } catch (err) {
      error.value = err instanceof Error ? err.message : 'An error occurred';
      throw err;
    } finally {
      loading.value = false;
    }
  };

  // Get a single daily item
  const getItem = async (id: string): Promise<DailyChecklistItem> => {
    loading.value = true;
    error.value = null;
    
    try {
      const response = await makeAuthenticatedRequest(`/api/v1/daily/${id}`);
      
      if (!response.ok) {
        const apiError = await handleHttpError(response);
        logError(apiError, 'getItem');
        showErrorAlert(apiError);
        throw new Error(apiError.message);
      }
      
      const item = await response.json();
      return {
        ...item,
        date: new Date(item.date),
        deadline: item.deadline ? new Date(item.deadline) : undefined,
        completion_time: item.completion_time ? new Date(item.completion_time) : undefined,
        created_at: new Date(item.created_at),
        updated_at: new Date(item.updated_at)
      };
    } catch (err) {
      error.value = err instanceof Error ? err.message : 'An error occurred';
      throw err;
    } finally {
      loading.value = false;
    }
  };

  // Create a new daily item
  const addItem = async (item: DailyChecklistItem): Promise<DailyChecklistItem> => {
    loading.value = true;
    error.value = null;
    
    try {
      const requestData: CreateDailyItemRequest = {
        title: item.title,
        description: item.description,
        priority: item.priority,
        estimated_time: item.estimated_time,
        deadline: item.deadline,
        context: item.context,
        notes: item.notes,
        inbox_item_id: item.inbox_item_id,
        date: item.date
      };
      
      const response = await makeAuthenticatedRequest('/api/v1/daily', {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json'
        },
        body: JSON.stringify(requestData)
      });
      
      if (!response.ok) {
        const apiError = await handleHttpError(response);
        logError(apiError, 'addItem');
        showErrorAlert(apiError);
        throw new Error(apiError.message);
      }
      
      const newItem = await response.json();
      const convertedItem = {
        ...newItem,
        date: new Date(newItem.date),
        deadline: newItem.deadline ? new Date(newItem.deadline) : undefined,
        completion_time: newItem.completion_time ? new Date(newItem.completion_time) : undefined,
        created_at: new Date(newItem.created_at),
        updated_at: new Date(newItem.updated_at)
      };
      
      items.value.unshift(convertedItem);
      totalCount.value += 1;
      return convertedItem;
    } catch (err) {
      error.value = err instanceof Error ? err.message : 'An error occurred';
      throw err;
    } finally {
      loading.value = false;
    }
  };

  // Update an existing daily item
  const updateItem = async (id: string, updates: Partial<DailyChecklistItem>): Promise<DailyChecklistItem> => {
    loading.value = true;
    error.value = null;
    
    try {
      const requestData: UpdateDailyItemRequest = {};
      
      if (updates.title !== undefined) requestData.title = updates.title;
      if (updates.description !== undefined) requestData.description = updates.description;
      if (updates.priority !== undefined) requestData.priority = updates.priority;
      if (updates.estimated_time !== undefined) requestData.estimated_time = updates.estimated_time;
      if (updates.deadline !== undefined) requestData.deadline = updates.deadline;
      if (updates.context !== undefined) requestData.context = updates.context;
      if (updates.status !== undefined) requestData.status = updates.status;
      if (updates.actual_time !== undefined) requestData.actual_time = updates.actual_time;
      if (updates.notes !== undefined) requestData.notes = updates.notes;
      
      const response = await makeAuthenticatedRequest(`/api/v1/daily/${id}`, {
        method: 'PUT',
        headers: {
          'Content-Type': 'application/json'
        },
        body: JSON.stringify(requestData)
      });
      
      if (!response.ok) {
        const apiError = await handleHttpError(response);
        logError(apiError, 'updateItem');
        showErrorAlert(apiError);
        throw new Error(apiError.message);
      }
      
      const updatedItem = await response.json();
      const convertedItem = {
        ...updatedItem,
        date: new Date(updatedItem.date),
        deadline: updatedItem.deadline ? new Date(updatedItem.deadline) : undefined,
        completion_time: updatedItem.completion_time ? new Date(updatedItem.completion_time) : undefined,
        created_at: new Date(updatedItem.created_at),
        updated_at: new Date(updatedItem.updated_at)
      };
      
      const index = items.value.findIndex(item => item.id === id);
      if (index !== -1) {
        items.value[index] = convertedItem;
      }
      return convertedItem;
    } catch (err) {
      error.value = err instanceof Error ? err.message : 'An error occurred';
      throw err;
    } finally {
      loading.value = false;
    }
  };

  // Delete a daily item
  const deleteItem = async (id: string): Promise<void> => {
    loading.value = true;
    error.value = null;
    
    try {
      const response = await makeAuthenticatedRequest(`/api/v1/daily/${id}`, {
        method: 'DELETE'
      });
      
      if (!response.ok) {
        const apiError = await handleHttpError(response);
        logError(apiError, 'deleteItem');
        showErrorAlert(apiError);
        throw new Error(apiError.message);
      }
      
      items.value = items.value.filter(item => item.id !== id);
      totalCount.value -= 1;
    } catch (err) {
      error.value = err instanceof Error ? err.message : 'An error occurred';
      throw err;
    } finally {
      loading.value = false;
    }
  };

  // Update item status
  const updateStatus = async (id: string, status: string): Promise<DailyChecklistItem> => {
    loading.value = true;
    error.value = null;
    
    try {
      const response = await makeAuthenticatedRequest(`/api/v1/daily/${id}/status`, {
        method: 'PUT',
        headers: {
          'Content-Type': 'application/json'
        },
        body: JSON.stringify({ status })
      });
      
      if (!response.ok) {
        const apiError = await handleHttpError(response);
        logError(apiError, 'updateStatus');
        showErrorAlert(apiError);
        throw new Error(apiError.message);
      }
      
      const updatedItem = await response.json();
      const convertedItem = {
        ...updatedItem,
        date: new Date(updatedItem.date),
        deadline: updatedItem.deadline ? new Date(updatedItem.deadline) : undefined,
        completion_time: updatedItem.completion_time ? new Date(updatedItem.completion_time) : undefined,
        created_at: new Date(updatedItem.created_at),
        updated_at: new Date(updatedItem.updated_at)
      };
      
      const index = items.value.findIndex(item => item.id === id);
      if (index !== -1) {
        items.value[index] = convertedItem;
      }
      return convertedItem;
    } catch (err) {
      error.value = err instanceof Error ? err.message : 'An error occurred';
      throw err;
    } finally {
      loading.value = false;
    }
  };

  // Update actual time
  const updateActualTime = async (id: string, actualTime: number): Promise<DailyChecklistItem> => {
    loading.value = true;
    error.value = null;
    
    try {
      const response = await makeAuthenticatedRequest(`/api/v1/daily/${id}/time`, {
        method: 'PUT',
        headers: {
          'Content-Type': 'application/json'
        },
        body: JSON.stringify({ actual_time: actualTime })
      });
      
      if (!response.ok) {
        const apiError = await handleHttpError(response);
        logError(apiError, 'updateActualTime');
        showErrorAlert(apiError);
        throw new Error(apiError.message);
      }
      
      const updatedItem = await response.json();
      const convertedItem = {
        ...updatedItem,
        date: new Date(updatedItem.date),
        deadline: updatedItem.deadline ? new Date(updatedItem.deadline) : undefined,
        completion_time: updatedItem.completion_time ? new Date(updatedItem.completion_time) : undefined,
        created_at: new Date(updatedItem.created_at),
        updated_at: new Date(updatedItem.updated_at)
      };
      
      const index = items.value.findIndex(item => item.id === id);
      if (index !== -1) {
        items.value[index] = convertedItem;
      }
      return convertedItem;
    } catch (err) {
      error.value = err instanceof Error ? err.message : 'An error occurred';
      throw err;
    } finally {
      loading.value = false;
    }
  };

  // Bulk update status
  const bulkUpdateStatus = async (ids: string[], status: string): Promise<void> => {
    loading.value = true;
    error.value = null;
    
    try {
      const response = await makeAuthenticatedRequest('/api/v1/daily/bulk/status', {
        method: 'PUT',
        headers: {
          'Content-Type': 'application/json'
        },
        body: JSON.stringify({ ids, status })
      });
      
      if (!response.ok) {
        const apiError = await handleHttpError(response);
        logError(apiError, 'bulkUpdateStatus');
        showErrorAlert(apiError);
        throw new Error(apiError.message);
      }
      
      // Update local items
      items.value.forEach(item => {
        if (ids.includes(item.id)) {
          item.status = status as any;
          if (status === 'completed') {
            item.completion_time = new Date();
          }
          item.updated_at = new Date();
        }
      });
    } catch (err) {
      error.value = err instanceof Error ? err.message : 'An error occurred';
      throw err;
    } finally {
      loading.value = false;
    }
  };

  // Clear store
  const clearItems = () => {
    items.value = [];
    totalCount.value = 0;
    error.value = null;
    searchQuery.value = '';
    currentPage.value = 1;
    selectedDate.value = new Date();
    filters.value = {
      status: '',
      priority: '',
      context: ''
    };
  };

  return {
    items,
    totalCount,
    loading,
    error,
    searchQuery,
    currentPage,
    pageSize,
    selectedDate,
    filters,
    fetchItems,
    getItemsByDate,
    getStats,
    getItem,
    addItem,
    updateItem,
    deleteItem,
    updateStatus,
    updateActualTime,
    bulkUpdateStatus,
    clearItems
  };
});
