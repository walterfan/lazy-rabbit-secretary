import { defineStore } from 'pinia';
import { ref, computed } from 'vue';
import { makeAuthenticatedRequest } from '@/utils/httpInterceptor';
import { getApiUrl } from '@/utils/apiConfig';

// Types
export interface User {
  id: string;
  realm_id: string;
  username: string;
  email: string;
  is_active: boolean;
  status: string;
  email_confirmed_at: string | null;
  created_by: string;
  created_at: string;
  updated_by: string;
  updated_at: string;
}

export interface UsersResponse {
  users: User[];
  total: number;
  page: number;
  page_size: number;
  total_pages: number;
}

export interface CreateUserRequest {
  username: string;
  email: string;
  realm_name: string;
  password: string;
  is_active?: boolean;
}

export interface UpdateUserRequest {
  username?: string;
  email?: string;
  is_active?: boolean;
  status?: string;
  password?: string; // For password changes
}

export interface UserFilters {
  realm_name?: string;
  status?: string;
  page?: number;
  page_size?: number;
}

export interface RegistrationStats {
  total_pending: number;
  total_approved: number;
  total_denied: number;
  total_users: number;
}

export const useUserStore = defineStore('user', () => {
  // State
  const users = ref<User[]>([]);
  const currentUser = ref<User | null>(null);
  const loading = ref(false);
  const error = ref<string | null>(null);
  const pagination = ref({
    total: 0,
    page: 1,
    page_size: 10,
    total_pages: 0,
  });

  // Computed
  const hasUsers = computed(() => users.value.length > 0);
  const totalUsers = computed(() => pagination.value.total);
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
   * Get users with pagination and filtering
   */
  const getUsers = async (filters: UserFilters = {}) => {
    setLoading(true);
    clearError();

    try {
      const params = new URLSearchParams();
      
      if (filters.realm_name) params.append('realm_name', filters.realm_name);
      if (filters.status) params.append('status', filters.status);
      if (filters.page) params.append('page', filters.page.toString());
      if (filters.page_size) params.append('page_size', filters.page_size.toString());

      const url = `${getApiUrl('/api/v1/admin/users')}${params.toString() ? `?${params.toString()}` : ''}`;
      
      const response = await makeAuthenticatedRequest(url, {
        method: 'GET',
      });

      if (!response.ok) {
        throw new Error(`Failed to fetch users: ${response.statusText}`);
      }

      const data: UsersResponse = await response.json();
      
      users.value = data.users;
      pagination.value = {
        total: data.total,
        page: data.page,
        page_size: data.page_size,
        total_pages: data.total_pages,
      };

      return data;
    } catch (err) {
      const errorMessage = err instanceof Error ? err.message : 'Failed to get users';
      setError(errorMessage);
      console.error('Failed to get users:', err);
      throw err;
    } finally {
      setLoading(false);
    }
  };

  /**
   * Get a single user by ID
   */
  const getUser = async (userId: string) => {
    setLoading(true);
    clearError();

    try {
      const url = getApiUrl(`/api/v1/admin/users/${userId}`);
      
      const response = await makeAuthenticatedRequest(url, {
        method: 'GET',
      });

      if (!response.ok) {
        throw new Error(`Failed to fetch user: ${response.statusText}`);
      }

      const data = await response.json();
      currentUser.value = data.user;
      
      return data.user;
    } catch (err) {
      const errorMessage = err instanceof Error ? err.message : 'Failed to get user';
      setError(errorMessage);
      console.error('Failed to get user:', err);
      throw err;
    } finally {
      setLoading(false);
    }
  };

  /**
   * Create a new user
   */
  const createUser = async (userData: CreateUserRequest) => {
    setLoading(true);
    clearError();

    try {
      const url = getApiUrl('/api/v1/admin/users');
      
      const response = await makeAuthenticatedRequest(url, {
        method: 'POST',
        body: JSON.stringify(userData),
      });

      if (!response.ok) {
        throw new Error(`Failed to create user: ${response.statusText}`);
      }

      const data = await response.json();
      
      // Add the new user to the list
      users.value.unshift(data.user);
      pagination.value.total += 1;
      
      return data.user;
    } catch (err) {
      const errorMessage = err instanceof Error ? err.message : 'Failed to create user';
      setError(errorMessage);
      console.error('Failed to create user:', err);
      throw err;
    } finally {
      setLoading(false);
    }
  };

  /**
   * Update an existing user
   */
  const updateUser = async (userId: string, userData: UpdateUserRequest) => {
    setLoading(true);
    clearError();

    try {
      const url = getApiUrl(`/api/v1/admin/users/${userId}`);
      
      const response = await makeAuthenticatedRequest(url, {
        method: 'PUT',
        body: JSON.stringify(userData),
      });

      if (!response.ok) {
        throw new Error(`Failed to update user: ${response.statusText}`);
      }

      const data = await response.json();
      
      // Update the user in the list
      const index = users.value.findIndex(user => user.id === userId);
      if (index !== -1) {
        users.value[index] = data.user;
      }
      
      // Update current user if it's the same user
      if (currentUser.value?.id === userId) {
        currentUser.value = data.user;
      }
      
      return data.user;
    } catch (err) {
      const errorMessage = err instanceof Error ? err.message : 'Failed to update user';
      setError(errorMessage);
      console.error('Failed to update user:', err);
      throw err;
    } finally {
      setLoading(false);
    }
  };

  /**
   * Delete a user
   */
  const deleteUser = async (userId: string) => {
    setLoading(true);
    clearError();

    try {
      const url = getApiUrl(`/api/v1/admin/users/${userId}`);
      
      const response = await makeAuthenticatedRequest(url, {
        method: 'DELETE',
      });

      if (!response.ok) {
        throw new Error(`Failed to delete user: ${response.statusText}`);
      }

      // Remove the user from the list
      users.value = users.value.filter(user => user.id !== userId);
      pagination.value.total -= 1;
      
      // Clear current user if it's the deleted user
      if (currentUser.value?.id === userId) {
        currentUser.value = null;
      }
      
      return true;
    } catch (err) {
      const errorMessage = err instanceof Error ? err.message : 'Failed to delete user';
      setError(errorMessage);
      console.error('Failed to delete user:', err);
      throw err;
    } finally {
      setLoading(false);
    }
  };

  /**
   * Get pending user registrations
   */
  const getPendingRegistrations = async (filters: UserFilters = {}) => {
    setLoading(true);
    clearError();

    try {
      const params = new URLSearchParams();
      
      if (filters.realm_name) params.append('realm_name', filters.realm_name);
      if (filters.status) params.append('status', filters.status);
      if (filters.page) params.append('page', filters.page.toString());
      if (filters.page_size) params.append('page_size', filters.page_size.toString());

      const url = `${getApiUrl('/api/v1/admin/registrations')}${params.toString() ? `?${params.toString()}` : ''}`;
      
      const response = await makeAuthenticatedRequest(url, {
        method: 'GET',
      });

      if (!response.ok) {
        throw new Error(`Failed to fetch registrations: ${response.statusText}`);
      }

      const data: UsersResponse = await response.json();
      
      users.value = data.users;
      pagination.value = {
        total: data.total,
        page: data.page,
        page_size: data.page_size,
        total_pages: data.total_pages,
      };

      return data;
    } catch (err) {
      const errorMessage = err instanceof Error ? err.message : 'Failed to get registrations';
      setError(errorMessage);
      console.error('Failed to get registrations:', err);
      throw err;
    } finally {
      setLoading(false);
    }
  };

  /**
   * Approve or deny user registration
   */
  const approveRegistration = async (userId: string, approved: boolean, reason?: string) => {
    setLoading(true);
    clearError();

    try {
      const url = getApiUrl('/api/v1/admin/registrations/approve');
      
      const response = await makeAuthenticatedRequest(url, {
        method: 'POST',
        body: JSON.stringify({
          user_id: userId,
          approved,
          reason,
        }),
      });

      if (!response.ok) {
        throw new Error(`Failed to process registration: ${response.statusText}`);
      }

      const data = await response.json();
      
      // Update the user in the list
      const index = users.value.findIndex(user => user.id === userId);
      if (index !== -1) {
        users.value[index] = data.user;
      }
      
      return data.user;
    } catch (err) {
      const errorMessage = err instanceof Error ? err.message : 'Failed to approve registration';
      setError(errorMessage);
      console.error('Failed to approve registration:', err);
      throw err;
    } finally {
      setLoading(false);
    }
  };

  /**
   * Get registration statistics
   */
  const getRegistrationStats = async (realmName?: string) => {
    setLoading(true);
    clearError();

    try {
      const params = new URLSearchParams();
      if (realmName) params.append('realm_name', realmName);

      const url = `${getApiUrl('/api/v1/admin/registrations/stats')}${params.toString() ? `?${params.toString()}` : ''}`;
      
      const response = await makeAuthenticatedRequest(url, {
        method: 'GET',
      });

      if (!response.ok) {
        throw new Error(`Failed to fetch registration stats: ${response.statusText}`);
      }

      const data = await response.json();
      return data.stats;
    } catch (err) {
      const errorMessage = err instanceof Error ? err.message : 'Failed to get registration stats';
      setError(errorMessage);
      console.error('Failed to get registration stats:', err);
      throw err;
    } finally {
      setLoading(false);
    }
  };

  /**
   * Search users by username or email
   */
  const searchUsers = (searchTerm: string) => {
    if (!searchTerm.trim()) {
      return users.value;
    }
    
    const term = searchTerm.toLowerCase();
    return users.value.filter(user =>
      user.username.toLowerCase().includes(term) ||
      user.email.toLowerCase().includes(term)
    );
  };

  /**
   * Get user by ID from current list
   */
  const getUserById = (userId: string) => {
    return users.value.find(user => user.id === userId);
  };

  /**
   * Clear all user data
   */
  const clearUsers = () => {
    users.value = [];
    currentUser.value = null;
    pagination.value = {
      total: 0,
      page: 1,
      page_size: 10,
      total_pages: 0,
    };
    clearError();
  };

  /**
   * Refresh current page of users
   */
  const refreshUsers = async () => {
    return getUsers({
      page: pagination.value.page,
      page_size: pagination.value.page_size,
    });
  };

  return {
    // State
    users,
    currentUser,
    loading,
    error,
    pagination,
    
    // Computed
    hasUsers,
    totalUsers,
    currentPage,
    totalPages,
    
    // Actions
    getUsers,
    getUser,
    createUser,
    updateUser,
    deleteUser,
    getPendingRegistrations,
    approveRegistration,
    getRegistrationStats,
    searchUsers,
    getUserById,
    clearUsers,
    refreshUsers,
    setLoading,
    setError,
    clearError,
  };
});
