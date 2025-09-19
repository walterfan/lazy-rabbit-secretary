import { defineStore } from 'pinia';
import { ref } from 'vue';
import type { 
  InboxItem, 
  CreateInboxItemRequest, 
  UpdateInboxItemRequest, 
  InboxListResponse,
  InboxStatsResponse 
} from '@/types/gtd';
import { handleHttpError, showErrorAlert, logError } from '@/utils/errorHandler';
import { makeAuthenticatedRequest } from '@/utils/httpInterceptor';

export const useInboxStore = defineStore('inbox', () => {
  const items = ref<InboxItem[]>([]);
  const totalCount = ref(0);
  const loading = ref(false);
  const error = ref<string | null>(null);
  const searchQuery = ref('');
  const currentPage = ref(1);
  const pageSize = ref(20);
  const filters = ref({
    status: '',
    priority: '',
    context: ''
  });

  // Fetch inbox items with pagination and filters
  const fetchItems = async (params: {
    page?: number;
    page_size?: number;
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
      
      const response = await makeAuthenticatedRequest(`/api/v1/inbox?${queryParams}`);
      
      if (!response.ok) {
        const apiError = await handleHttpError(response);
        logError(apiError, 'fetchItems');
        showErrorAlert(apiError);
        throw new Error(apiError.message);
      }
      
      const data: InboxListResponse = await response.json();
      
      items.value = data.items.map(item => ({
        ...item,
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

  // Get pending items
  const getPendingItems = async () => {
    loading.value = true;
    error.value = null;
    
    try {
      const response = await makeAuthenticatedRequest('/api/v1/inbox/pending');
      
      if (!response.ok) {
        const apiError = await handleHttpError(response);
        logError(apiError, 'getPendingItems');
        showErrorAlert(apiError);
        throw new Error(apiError.message);
      }
      
      const data = await response.json();
      return data.items.map((item: InboxItem) => ({
        ...item,
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

  // Get urgent items
  const getUrgentItems = async () => {
    loading.value = true;
    error.value = null;
    
    try {
      const response = await makeAuthenticatedRequest('/api/v1/inbox/urgent');
      
      if (!response.ok) {
        const apiError = await handleHttpError(response);
        logError(apiError, 'getUrgentItems');
        showErrorAlert(apiError);
        throw new Error(apiError.message);
      }
      
      const data = await response.json();
      return data.items.map((item: InboxItem) => ({
        ...item,
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

  // Get inbox statistics
  const getStats = async (): Promise<InboxStatsResponse> => {
    loading.value = true;
    error.value = null;
    
    try {
      const response = await makeAuthenticatedRequest('/api/v1/inbox/stats');
      
      if (!response.ok) {
        const apiError = await handleHttpError(response);
        logError(apiError, 'getStats');
        showErrorAlert(apiError);
        throw new Error(apiError.message);
      }
      
      return await response.json();
    } catch (err) {
      error.value = err instanceof Error ? err.message : 'An error occurred';
      throw err;
    } finally {
      loading.value = false;
    }
  };

  // Get a single inbox item
  const getItem = async (id: string): Promise<InboxItem> => {
    loading.value = true;
    error.value = null;
    
    try {
      const response = await makeAuthenticatedRequest(`/api/v1/inbox/${id}`);
      
      if (!response.ok) {
        const apiError = await handleHttpError(response);
        logError(apiError, 'getItem');
        showErrorAlert(apiError);
        throw new Error(apiError.message);
      }
      
      const item = await response.json();
      return {
        ...item,
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

  // Create a new inbox item
  const addItem = async (item: InboxItem): Promise<InboxItem> => {
    loading.value = true;
    error.value = null;
    
    try {
      const requestData: CreateInboxItemRequest = {
        title: item.title,
        description: item.description,
        priority: item.priority,
        tags: item.tags,
        context: item.context
      };
      
      const response = await makeAuthenticatedRequest('/api/v1/inbox', {
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

  // Update an existing inbox item
  const updateItem = async (id: string, updates: Partial<InboxItem>): Promise<InboxItem> => {
    loading.value = true;
    error.value = null;
    
    try {
      const requestData: UpdateInboxItemRequest = {};
      
      if (updates.title !== undefined) requestData.title = updates.title;
      if (updates.description !== undefined) requestData.description = updates.description;
      if (updates.priority !== undefined) requestData.priority = updates.priority;
      if (updates.status !== undefined) requestData.status = updates.status;
      if (updates.tags !== undefined) requestData.tags = updates.tags;
      if (updates.context !== undefined) requestData.context = updates.context;
      
      const response = await makeAuthenticatedRequest(`/api/v1/inbox/${id}`, {
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

  // Delete an inbox item
  const deleteItem = async (id: string): Promise<void> => {
    loading.value = true;
    error.value = null;
    
    try {
      const response = await makeAuthenticatedRequest(`/api/v1/inbox/${id}`, {
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
  const updateStatus = async (id: string, status: string): Promise<InboxItem> => {
    loading.value = true;
    error.value = null;
    
    try {
      const response = await makeAuthenticatedRequest(`/api/v1/inbox/${id}/status`, {
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
      const response = await makeAuthenticatedRequest('/api/v1/inbox/bulk/status', {
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
    filters,
    fetchItems,
    getPendingItems,
    getUrgentItems,
    getStats,
    getItem,
    addItem,
    updateItem,
    deleteItem,
    updateStatus,
    bulkUpdateStatus,
    clearItems
  };
});
