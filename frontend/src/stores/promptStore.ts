import { defineStore } from 'pinia';
import { ref } from 'vue';
import type { Prompt, CreatePromptRequest, UpdatePromptRequest, PromptListResponse } from '@/types';
import { handleHttpError, showErrorAlert, logError } from '@/utils/errorHandler';
import { makeAuthenticatedRequest } from '@/utils/httpInterceptor';
import { getApiUrl } from '@/utils/apiConfig';

export const usePromptStore = defineStore('prompt', () => {
  const prompts = ref<Prompt[]>([]);
  const totalCount = ref(0);
  const loading = ref(false);
  const error = ref<string | null>(null);
  const searchQuery = ref('');
  const currentPage = ref(1);
  const pageSize = ref(20);

  // Fetch prompts with pagination
  const fetchPrompts = async (params: {
    page?: number;
    page_size?: number;
  } = {}) => {
    loading.value = true;
    error.value = null;
    
    try {
      const queryParams = new URLSearchParams();
      queryParams.append('page', (params.page || currentPage.value).toString());
      queryParams.append('limit', (params.page_size || pageSize.value).toString());
      
      const response = await makeAuthenticatedRequest(getApiUrl(`/api/v1/prompts?${queryParams}`));
      
      if (!response.ok) {
        const apiError = await handleHttpError(response);
        logError(apiError, 'fetchPrompts');
        showErrorAlert(apiError);
        throw new Error(apiError.message);
      }
      
      const data = await response.json();
      
      // Convert backend prompt format to frontend format
      prompts.value = (data.prompts || []).map((prompt: any) => ({
        ...prompt,
        created_at: new Date(prompt.created_at),
        updated_at: new Date(prompt.updated_at)
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

  // Search prompts
  const searchPrompts = async (query: string, params: {
    page?: number;
    page_size?: number;
  } = {}) => {
    loading.value = true;
    error.value = null;
    
    try {
      const queryParams = new URLSearchParams();
      queryParams.append('q', query);
      queryParams.append('page', (params.page || currentPage.value).toString());
      queryParams.append('limit', (params.page_size || pageSize.value).toString());
      
      const response = await makeAuthenticatedRequest(getApiUrl(`/api/v1/prompts/search?${queryParams}`));
      
      if (!response.ok) {
        const apiError = await handleHttpError(response);
        logError(apiError, 'searchPrompts');
        showErrorAlert(apiError);
        throw new Error(apiError.message);
      }
      
      const data = await response.json();
      
      // Convert backend prompt format to frontend format
      prompts.value = (data.prompts || []).map((prompt: any) => ({
        ...prompt,
        created_at: new Date(prompt.created_at),
        updated_at: new Date(prompt.updated_at)
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

  // Get prompts by tags
  const getPromptsByTags = async (tags: string[], params: {
    page?: number;
    page_size?: number;
  } = {}) => {
    loading.value = true;
    error.value = null;
    
    try {
      const queryParams = new URLSearchParams();
      queryParams.append('page', (params.page || currentPage.value).toString());
      queryParams.append('limit', (params.page_size || pageSize.value).toString());
      
      const tagsParam = tags.join(',');
      const response = await makeAuthenticatedRequest(getApiUrl(`/api/v1/prompts/tags/${encodeURIComponent(tagsParam)}?${queryParams}`));
      
      if (!response.ok) {
        const apiError = await handleHttpError(response);
        logError(apiError, 'getPromptsByTags');
        showErrorAlert(apiError);
        throw new Error(apiError.message);
      }
      
      const data = await response.json();
      
      // Convert backend prompt format to frontend format
      prompts.value = (data.prompts || []).map((prompt: any) => ({
        ...prompt,
        created_at: new Date(prompt.created_at),
        updated_at: new Date(prompt.updated_at)
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

  // Get a single prompt
  const getPrompt = async (id: string): Promise<Prompt> => {
    loading.value = true;
    error.value = null;
    
    try {
      const response = await makeAuthenticatedRequest(getApiUrl(`/api/v1/prompts/${id}`));

      if (!response.ok) {
        const apiError = await handleHttpError(response);
        logError(apiError, 'getPrompt');
        showErrorAlert(apiError);
        throw new Error(apiError.message);
      }
      
      const prompt = await response.json();
      
      // Convert backend format to frontend format
      return {
        ...prompt,
        created_at: new Date(prompt.created_at),
        updated_at: new Date(prompt.updated_at)
      };
    } catch (err) {
      error.value = err instanceof Error ? err.message : 'An error occurred';
      throw err;
    } finally {
      loading.value = false;
    }
  };

  // Create a new prompt
  const addPrompt = async (prompt: Prompt): Promise<Prompt> => {
    loading.value = true;
    error.value = null;
    
    try {
      const requestData: CreatePromptRequest = {
        name: prompt.name,
        description: prompt.description,
        system_prompt: prompt.system_prompt,
        user_prompt: prompt.user_prompt,
        tags: prompt.tags
      };
      
      const response = await makeAuthenticatedRequest(getApiUrl('/api/v1/prompts'), {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json'
        },
        body: JSON.stringify(requestData)
      });
      
      if (!response.ok) {
        const apiError = await handleHttpError(response);
        logError(apiError, 'addPrompt');
        showErrorAlert(apiError);
        throw new Error(apiError.message);
      }
      
      const newPrompt = await response.json();
      
      // Convert and add to local store
      const convertedPrompt = {
        ...newPrompt,
        created_at: new Date(newPrompt.created_at),
        updated_at: new Date(newPrompt.updated_at)
      };
      
      prompts.value.push(convertedPrompt);
      totalCount.value += 1;
      return convertedPrompt;
    } catch (err) {
      error.value = err instanceof Error ? err.message : 'An error occurred';
      throw err;
    } finally {
      loading.value = false;
    }
  };

  // Update an existing prompt
  const updatePrompt = async (id: string, updates: Partial<Prompt>): Promise<Prompt> => {
    loading.value = true;
    error.value = null;
    
    try {
      const requestData: UpdatePromptRequest = {};
      
      if (updates.name !== undefined) requestData.name = updates.name;
      if (updates.description !== undefined) requestData.description = updates.description;
      if (updates.system_prompt !== undefined) requestData.system_prompt = updates.system_prompt;
      if (updates.user_prompt !== undefined) requestData.user_prompt = updates.user_prompt;
      if (updates.tags !== undefined) requestData.tags = updates.tags;
      
      const response = await makeAuthenticatedRequest(getApiUrl(`/api/v1/prompts/${id}`), {
        method: 'PUT',
        headers: {
          'Content-Type': 'application/json'
        },
        body: JSON.stringify(requestData)
      });
      
      if (!response.ok) {
        const apiError = await handleHttpError(response);
        logError(apiError, 'updatePrompt');
        showErrorAlert(apiError);
        throw new Error(apiError.message);
      }
      
      const updatedPrompt = await response.json();
      
      // Convert and update in local store
      const convertedPrompt = {
        ...updatedPrompt,
        created_at: new Date(updatedPrompt.created_at),
        updated_at: new Date(updatedPrompt.updated_at)
      };
      
      const index = prompts.value.findIndex(p => p.id === id);
      if (index !== -1) {
        prompts.value[index] = convertedPrompt;
      }
      return convertedPrompt;
    } catch (err) {
      error.value = err instanceof Error ? err.message : 'An error occurred';
      throw err;
    } finally {
      loading.value = false;
    }
  };

  // Delete a prompt
  const deletePrompt = async (id: string): Promise<void> => {
    loading.value = true;
    error.value = null;
    
    try {
      const response = await makeAuthenticatedRequest(getApiUrl(`/api/v1/prompts/${id}`), {
        method: 'DELETE'
      });
      
      if (!response.ok) {
        const apiError = await handleHttpError(response);
        logError(apiError, 'deletePrompt');
        showErrorAlert(apiError);
        throw new Error(apiError.message);
      }
      
      prompts.value = prompts.value.filter(p => p.id !== id);
      totalCount.value -= 1;
    } catch (err) {
      error.value = err instanceof Error ? err.message : 'An error occurred';
      throw err;
    } finally {
      loading.value = false;
    }
  };

  // Clear store
  const clearPrompts = () => {
    prompts.value = [];
    totalCount.value = 0;
    error.value = null;
    searchQuery.value = '';
    currentPage.value = 1;
  };

  return {
    prompts,
    totalCount,
    loading,
    error,
    searchQuery,
    currentPage,
    pageSize,
    fetchPrompts,
    searchPrompts,
    getPromptsByTags,
    getPrompt,
    addPrompt,
    updatePrompt,
    deletePrompt,
    clearPrompts
  };
});
