import { defineStore } from 'pinia';
import { ref } from 'vue';
import type { Secret, CreateSecretRequest, UpdateSecretRequest } from '@/types';
import { handleHttpError, showErrorAlert, logError } from '@/utils/errorHandler';
import { makeAuthenticatedRequest } from '@/utils/httpInterceptor';

export const useSecretStore = defineStore('secret', () => {
  const secrets = ref<Secret[]>([]);
  const totalCount = ref(0);
  const loading = ref(false);
  const error = ref<string | null>(null);

  // Fetch secrets with search and filters
  const fetchSecrets = async (params: {
    realm_id?: string;
    q?: string;
    group?: string;
    path?: string;
    page?: number;
    page_size?: number;
  }) => {
    loading.value = true;
    error.value = null;
    
    try {
      const queryParams = new URLSearchParams();
      if (params.realm_id) queryParams.append('realm_id', params.realm_id);
      if (params.q) queryParams.append('q', params.q);
      if (params.group) queryParams.append('group', params.group);
      if (params.path) queryParams.append('path', params.path);
      if (params.page) queryParams.append('page', params.page.toString());
      if (params.page_size) queryParams.append('page_size', params.page_size.toString());
      
      const response = await makeAuthenticatedRequest(`/api/v1/secrets?${queryParams}`);
      
      if (!response.ok) {
        const apiError = await handleHttpError(response);
        logError(apiError, 'fetchSecrets');
        showErrorAlert(apiError);
        throw new Error(apiError.message);
      }
      
      const data = await response.json();
      secrets.value = data.items || [];
      totalCount.value = data.total || 0;
    } catch (err) {
      error.value = err instanceof Error ? err.message : 'An error occurred';
      throw err;
    } finally {
      loading.value = false;
    }
  };

  // Get a single secret
  const getSecret = async (id: string): Promise<Secret> => {
    loading.value = true;
    error.value = null;
    
    try {
      const response = await fetch(`/api/v1/secrets/${id}`, {
        headers: {
          'Authorization': `Bearer ${localStorage.getItem('access_token')}`
        }
      });
      
      if (!response.ok) {
        throw new Error('Failed to fetch secret');
      }
      
      return await response.json();
    } catch (err) {
      error.value = err instanceof Error ? err.message : 'An error occurred';
      throw err;
    } finally {
      loading.value = false;
    }
  };

  // Create a new secret
  const createSecret = async (secretData: CreateSecretRequest): Promise<Secret> => {
    loading.value = true;
    error.value = null;
    
    try {
      const response = await makeAuthenticatedRequest('/api/v1/secrets', {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json'
        },
        body: JSON.stringify(secretData)
      });
      
      if (!response.ok) {
        const apiError = await handleHttpError(response);
        logError(apiError, 'createSecret');
        showErrorAlert(apiError);
        throw new Error(apiError.message);
      }
      
      const newSecret = await response.json();
      secrets.value.push(newSecret);
      totalCount.value += 1;
      return newSecret;
    } catch (err) {
      error.value = err instanceof Error ? err.message : 'An error occurred';
      throw err;
    } finally {
      loading.value = false;
    }
  };

  // Update an existing secret
  const updateSecret = async (id: string, updates: UpdateSecretRequest): Promise<Secret> => {
    loading.value = true;
    error.value = null;
    
    try {
      const response = await makeAuthenticatedRequest(`/api/v1/secrets/${id}`, {
        method: 'PUT',
        headers: {
          'Content-Type': 'application/json'
        },
        body: JSON.stringify(updates)
      });

      if (!response.ok) {
        const apiError = await handleHttpError(response);
        logError(apiError, 'updateSecret');
        showErrorAlert(apiError);
        throw new Error(apiError.message);
      }
      
      const updatedSecret = await response.json();
      const index = secrets.value.findIndex(s => s.id === id);
      if (index !== -1) {
        secrets.value[index] = updatedSecret;
      }
      return updatedSecret;
    } catch (err) {
      error.value = err instanceof Error ? err.message : 'An error occurred';
      throw err;
    } finally {
      loading.value = false;
    }
  };

  // Delete a secret
  const deleteSecret = async (id: string): Promise<void> => {
    loading.value = true;
    error.value = null;
    
    try {
      const response = await makeAuthenticatedRequest(`/api/v1/secrets/${id}`, {
        method: 'DELETE'
      });
      
      if (!response.ok) {
        const apiError = await handleHttpError(response);
        logError(apiError, 'deleteSecret');
        showErrorAlert(apiError);
        throw new Error(apiError.message);
      }
      
      secrets.value = secrets.value.filter(s => s.id !== id);
      totalCount.value -= 1;
    } catch (err) {
      error.value = err instanceof Error ? err.message : 'An error occurred';
      throw err;
    } finally {
      loading.value = false;
    }
  };

  // Decrypt and copy secret value
  const copySecretValue = async (secret: Secret): Promise<void> => {
    // In a real implementation, this would call a special endpoint
    // to decrypt the secret value server-side
    try {
      const response = await makeAuthenticatedRequest(`/api/v1/secrets/${secret.id}/decrypt`, {
        method: 'POST'
      });
      
      if (!response.ok) {
        const apiError = await handleHttpError(response);
        logError(apiError, 'copySecretValue');
        showErrorAlert(apiError);
        throw new Error(apiError.message);
      }
      
      const { value } = await response.json();
      
      // Copy to clipboard
      await navigator.clipboard.writeText(value);
      
      // Show success message (you might want to use a toast library)
      alert('Secret value copied to clipboard');
    } catch (err) {
      console.error('Failed to copy secret:', err);
      alert('Failed to copy secret value');
    }
  };

  // Decrypt and copy secret value with custom KEK
  const copySecretValueWithKEK = async (secret: Secret, kek: string): Promise<void> => {
    try {
      const response = await makeAuthenticatedRequest(`/api/v1/secrets/${secret.id}/decrypt-with-kek`, {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json'
        },
        body: JSON.stringify({ kek })
      });
      
      if (!response.ok) {
        const apiError = await handleHttpError(response);
        logError(apiError, 'copySecretValueWithKEK');
        showErrorAlert(apiError);
        throw new Error(apiError.message);
      }
      
      const { value } = await response.json();
      
      // Copy to clipboard
      await navigator.clipboard.writeText(value);
      
      // Show success message
      alert('Secret value copied to clipboard');
    } catch (err) {
      console.error('Failed to copy secret with KEK:', err);
      alert('Failed to copy secret value. Please check your KEK.');
    }
  };

  // Clear store
  const clearSecrets = () => {
    secrets.value = [];
    totalCount.value = 0;
    error.value = null;
  };

  return {
    secrets,
    totalCount,
    loading,
    error,
    fetchSecrets,
    getSecret,
    createSecret,
    updateSecret,
    deleteSecret,
    copySecretValue,
    copySecretValueWithKEK,
    clearSecrets
  };
});
