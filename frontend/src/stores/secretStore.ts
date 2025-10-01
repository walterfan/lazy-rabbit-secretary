import { defineStore } from 'pinia';
import { ref } from 'vue';
import type { Secret, SecretVersion, CreateSecretRequest, UpdateSecretRequest, CreatePendingVersionRequest, DecryptVersionRequest } from '@/types';
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

  // Get secret versions
  const getSecretVersions = async (secretId: string): Promise<SecretVersion[]> => {
    loading.value = true;
    error.value = null;
    
    try {
      const response = await makeAuthenticatedRequest(`/api/v1/secrets/${secretId}/versions`);
      
      if (!response.ok) {
        const apiError = await handleHttpError(response);
        logError(apiError, 'getSecretVersions');
        showErrorAlert(apiError);
        throw new Error(apiError.message);
      }
      
      const data = await response.json();
      return data.versions || [];
    } catch (err) {
      error.value = err instanceof Error ? err.message : 'An error occurred';
      throw err;
    } finally {
      loading.value = false;
    }
  };

  // Decrypt specific version
  const decryptSecretVersion = async (secretId: string, version: number, kek?: string): Promise<string> => {
    loading.value = true;
    error.value = null;
    
    try {
      const endpoint = kek 
        ? `/api/v1/secrets/${secretId}/versions/${version}/decrypt-with-kek`
        : `/api/v1/secrets/${secretId}/versions/${version}/decrypt`;
      
      const body = kek ? { kek } : {};
      
      const response = await makeAuthenticatedRequest(endpoint, {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json'
        },
        body: JSON.stringify(body)
      });
      
      if (!response.ok) {
        const apiError = await handleHttpError(response);
        logError(apiError, 'decryptSecretVersion');
        showErrorAlert(apiError);
        throw new Error(apiError.message);
      }
      
      const { value } = await response.json();
      return value;
    } catch (err) {
      error.value = err instanceof Error ? err.message : 'An error occurred';
      throw err;
    } finally {
      loading.value = false;
    }
  };

  // Activate a specific version
  const activateSecretVersion = async (secretId: string, version: number): Promise<void> => {
    loading.value = true;
    error.value = null;
    
    try {
      const response = await makeAuthenticatedRequest(`/api/v1/secrets/${secretId}/versions/${version}/activate`, {
        method: 'POST'
      });
      
      if (!response.ok) {
        const apiError = await handleHttpError(response);
        logError(apiError, 'activateSecretVersion');
        showErrorAlert(apiError);
        throw new Error(apiError.message);
      }
      
      // Refresh the secret list to get updated version info
      await fetchSecrets({});
    } catch (err) {
      error.value = err instanceof Error ? err.message : 'An error occurred';
      throw err;
    } finally {
      loading.value = false;
    }
  };

  // Delete a specific version
  const deleteSecretVersion = async (secretId: string, version: number): Promise<void> => {
    loading.value = true;
    error.value = null;
    
    try {
      const response = await makeAuthenticatedRequest(`/api/v1/secrets/${secretId}/versions/${version}`, {
        method: 'DELETE'
      });
      
      if (!response.ok) {
        const apiError = await handleHttpError(response);
        logError(apiError, 'deleteSecretVersion');
        showErrorAlert(apiError);
        throw new Error(apiError.message);
      }
      
      // Refresh the secret list to get updated version info
      await fetchSecrets({});
    } catch (err) {
      error.value = err instanceof Error ? err.message : 'An error occurred';
      throw err;
    } finally {
      loading.value = false;
    }
  };

  // Create a pending version
  const createPendingVersion = async (secretId: string, request: CreatePendingVersionRequest): Promise<SecretVersion> => {
    loading.value = true;
    error.value = null;
    
    try {
      const response = await makeAuthenticatedRequest(`/api/v1/secrets/${secretId}/versions/pending`, {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json'
        },
        body: JSON.stringify(request)
      });
      
      if (!response.ok) {
        const apiError = await handleHttpError(response);
        logError(apiError, 'createPendingVersion');
        showErrorAlert(apiError);
        throw new Error(apiError.message);
      }
      
      const newVersion = await response.json();
      
      // Refresh the secret list to get updated version info
      await fetchSecrets({});
      
      return newVersion;
    } catch (err) {
      error.value = err instanceof Error ? err.message : 'An error occurred';
      throw err;
    } finally {
      loading.value = false;
    }
  };

  // Copy secret value from specific version
  const copySecretVersionValue = async (secretId: string, version: number, kek?: string): Promise<void> => {
    try {
      const value = await decryptSecretVersion(secretId, version, kek);
      
      // Copy to clipboard
      await navigator.clipboard.writeText(value);
      
      // Show success message
      alert('Secret value copied to clipboard');
    } catch (err) {
      console.error('Failed to copy secret version:', err);
      alert('Failed to copy secret value');
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
    getSecretVersions,
    decryptSecretVersion,
    activateSecretVersion,
    deleteSecretVersion,
    createPendingVersion,
    copySecretVersionValue,
    clearSecrets
  };
});
