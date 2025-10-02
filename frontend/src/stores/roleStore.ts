import { defineStore } from 'pinia';
import { ref, computed } from 'vue';
import { makeAuthenticatedRequest } from '@/utils/httpInterceptor';
import { getApiUrl } from '@/utils/apiConfig';

// Types
export interface Policy {
  id: string;
  realm_id: string;
  name: string;
  description: string;
  version: string;
  created_by: string;
  created_at: string;
  updated_by: string;
  updated_at: string;
}

export interface Role {
  id: string;
  realm_id: string;
  name: string;
  description: string;
  created_by: string;
  created_at: string;
  updated_by: string;
  updated_at: string;
  policies?: Policy[];
}

export interface RolesResponse {
  roles: Role[];
  total: number;
  page: number;
  page_size: number;
  total_pages: number;
}

export interface CreateRoleRequest {
  realm_name: string;
  name: string;
  description?: string;
}

export interface UpdateRoleRequest {
  name?: string;
  description?: string;
}

export interface RoleFilters {
  realm_name?: string;
  page?: number;
  page_size?: number;
}

export const useRoleStore = defineStore('role', () => {
  // State
  const roles = ref<Role[]>([]);
  const currentRole = ref<Role | null>(null);
  const loading = ref(false);
  const error = ref<string | null>(null);
  const pagination = ref({
    total: 0,
    page: 1,
    page_size: 10,
    total_pages: 0,
  });

  // Computed
  const hasRoles = computed(() => roles.value.length > 0);
  const totalRoles = computed(() => pagination.value.total);
  const currentPage = computed(() => pagination.value.page);
  const totalPages = computed(() => pagination.value.total_pages);

  // Actions
  const setLoading = (isLoading: boolean) => {
    loading.value = isLoading;
  };

  const setError = (errorMessage: string | null) => {
    error.value = errorMessage;
  };

  const clearError = () => {
    error.value = null;
  };

  /**
   * Get roles with pagination and filtering
   */
  const getRoles = async (filters: RoleFilters = {}) => {
    setLoading(true);
    clearError();

    try {
      const params = new URLSearchParams();
      
      if (filters.realm_name) params.append('realm_name', filters.realm_name);
      if (filters.page) params.append('page', filters.page.toString());
      if (filters.page_size) params.append('page_size', filters.page_size.toString());

      const url = `${getApiUrl('/api/v1/admin/roles')}${params.toString() ? `?${params.toString()}` : ''}`;
      
      const response = await makeAuthenticatedRequest(url, {
        method: 'GET',
      });

      if (!response.ok) {
        throw new Error(`Failed to fetch roles: ${response.statusText}`);
      }

      const data: RolesResponse = await response.json();
      
      roles.value = data.roles;
      pagination.value = {
        total: data.total,
        page: data.page,
        page_size: data.page_size,
        total_pages: data.total_pages,
      };

      return data;
    } catch (err) {
      const errorMessage = err instanceof Error ? err.message : 'Failed to get roles';
      setError(errorMessage);
      console.error('Failed to get roles:', err);
      throw err;
    } finally {
      setLoading(false);
    }
  };

  /**
   * Get a single role by ID
   */
  const getRole = async (roleId: string) => {
    setLoading(true);
    clearError();

    try {
      const url = getApiUrl(`/api/v1/admin/roles/${roleId}`);
      
      const response = await makeAuthenticatedRequest(url, {
        method: 'GET',
      });

      if (!response.ok) {
        throw new Error(`Failed to fetch role: ${response.statusText}`);
      }

      const data = await response.json();
      currentRole.value = data.role;
      
      return data.role;
    } catch (err) {
      const errorMessage = err instanceof Error ? err.message : 'Failed to get role';
      setError(errorMessage);
      console.error('Failed to get role:', err);
      throw err;
    } finally {
      setLoading(false);
    }
  };

  /**
   * Create a new role
   */
  const createRole = async (roleData: CreateRoleRequest) => {
    setLoading(true);
    clearError();

    try {
      const url = getApiUrl('/api/v1/admin/roles');
      
      const response = await makeAuthenticatedRequest(url, {
        method: 'POST',
        body: JSON.stringify(roleData),
      });

      if (!response.ok) {
        throw new Error(`Failed to create role: ${response.statusText}`);
      }

      const data = await response.json();
      
      // Add the new role to the list
      roles.value.unshift(data.role);
      pagination.value.total += 1;
      
      return data.role;
    } catch (err) {
      const errorMessage = err instanceof Error ? err.message : 'Failed to create role';
      setError(errorMessage);
      console.error('Failed to create role:', err);
      throw err;
    } finally {
      setLoading(false);
    }
  };

  /**
   * Update an existing role
   */
  const updateRole = async (roleId: string, roleData: UpdateRoleRequest) => {
    setLoading(true);
    clearError();

    try {
      const url = getApiUrl(`/api/v1/admin/roles/${roleId}`);
      
      const response = await makeAuthenticatedRequest(url, {
        method: 'PUT',
        body: JSON.stringify(roleData),
      });

      if (!response.ok) {
        throw new Error(`Failed to update role: ${response.statusText}`);
      }

      const data = await response.json();
      
      // Update the role in the list
      const index = roles.value.findIndex(role => role.id === roleId);
      if (index !== -1) {
        roles.value[index] = data.role;
      }
      
      // Update current role if it's the same role
      if (currentRole.value?.id === roleId) {
        currentRole.value = data.role;
      }
      
      return data.role;
    } catch (err) {
      const errorMessage = err instanceof Error ? err.message : 'Failed to update role';
      setError(errorMessage);
      console.error('Failed to update role:', err);
      throw err;
    } finally {
      setLoading(false);
    }
  };

  /**
   * Delete a role
   */
  const deleteRole = async (roleId: string) => {
    setLoading(true);
    clearError();

    try {
      const url = getApiUrl(`/api/v1/admin/roles/${roleId}`);
      
      const response = await makeAuthenticatedRequest(url, {
        method: 'DELETE',
      });

      if (!response.ok) {
        throw new Error(`Failed to delete role: ${response.statusText}`);
      }

      // Remove the role from the list
      roles.value = roles.value.filter(role => role.id !== roleId);
      pagination.value.total -= 1;
      
      // Clear current role if it's the deleted role
      if (currentRole.value?.id === roleId) {
        currentRole.value = null;
      }
      
      return true;
    } catch (err) {
      const errorMessage = err instanceof Error ? err.message : 'Failed to delete role';
      setError(errorMessage);
      console.error('Failed to delete role:', err);
      throw err;
    } finally {
      setLoading(false);
    }
  };

  /**
   * Search roles by name or description
   */
  const searchRoles = (searchTerm: string) => {
    if (!searchTerm.trim()) {
      return roles.value;
    }
    
    const term = searchTerm.toLowerCase();
    return roles.value.filter(role =>
      role.name.toLowerCase().includes(term) ||
      role.description.toLowerCase().includes(term)
    );
  };

  /**
   * Get role by ID from current list
   */
  const getRoleById = (roleId: string) => {
    return roles.value.find(role => role.id === roleId);
  };

  /**
   * Clear all role data
   */
  const clearRoles = () => {
    roles.value = [];
    currentRole.value = null;
    pagination.value = {
      total: 0,
      page: 1,
      page_size: 10,
      total_pages: 0,
    };
    clearError();
  };

  /**
   * Refresh current page of roles
   */
  const refreshRoles = async () => {
    return getRoles({
      page: pagination.value.page,
      page_size: pagination.value.page_size,
    });
  };

  return {
    // State
    roles,
    currentRole,
    loading,
    error,
    pagination,
    
    // Computed
    hasRoles,
    totalRoles,
    currentPage,
    totalPages,
    
    // Actions
    getRoles,
    getRole,
    createRole,
    updateRole,
    deleteRole,
    searchRoles,
    getRoleById,
    clearRoles,
    refreshRoles,
    setLoading,
    setError,
    clearError,
  };
});
