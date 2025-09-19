import { defineStore } from 'pinia';
import { ref, computed } from 'vue';
import { makeAuthenticatedRequest } from '@/utils/httpInterceptor';
import { tokenRefreshService } from '@/services/tokenRefreshService';

export interface User {
  id: string;
  username: string;
  email: string;
  realm_name: string;
  realm_id: string;
  is_active: boolean;
  status: 'pending' | 'approved' | 'denied' | 'suspended';
  created_by: string;
  created_at: string;
  updated_by: string;
  updated_at: string;
}

export interface AuthState {
  user: User | null;
  token: string | null;
  refreshToken: string | null;
}

export const useAuthStore = defineStore('auth', () => {
  const user = ref<User | null>(null);
  const token = ref<string | null>(null);
  const refreshToken = ref<string | null>(null);

  const isAuthenticated = computed(() => !!token.value);
  const currentUser = computed(() => user.value);

  const signIn = async (username: string, password: string, realmName: string) => {
    try {
      const response = await fetch('/api/v1/auth/login', {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
        },
        body: JSON.stringify({
          username,
          password,
          realm_name: realmName,
        }),
      });

      if (!response.ok) {
        throw new Error('Login failed');
      }

      const data = await response.json();
      
      user.value = data.user;
      token.value = data.access_token;
      refreshToken.value = data.refresh_token;

      // Store tokens in localStorage
      localStorage.setItem('access_token', data.access_token);
      localStorage.setItem('refresh_token', data.refresh_token);
      localStorage.setItem('user', JSON.stringify(data.user));

      // Start token refresh monitoring
      tokenRefreshService.start();

      return data;
    } catch (error) {
      console.error('Sign in error:', error);
      throw error;
    }
  };

  const signUp = async (userData: {
    username: string;
    email: string;
    password: string;
    realm_name: string;
  }) => {
    try {
      const response = await fetch('/api/v1/auth/register', {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
        },
        body: JSON.stringify(userData),
      });

      if (!response.ok) {
        throw new Error('Registration failed');
      }

      const data = await response.json();
      return data;
    } catch (error) {
      console.error('Sign up error:', error);
      throw error;
    }
  };

  const signOut = () => {
    // Stop token refresh monitoring
    tokenRefreshService.stop();

    user.value = null;
    token.value = null;
    refreshToken.value = null;

    // Clear localStorage
    localStorage.removeItem('access_token');
    localStorage.removeItem('refresh_token');
    localStorage.removeItem('user');
  };

  const refreshAuth = async (): Promise<boolean> => {
    if (!refreshToken.value) {
      console.warn('No refresh token available');
      return false;
    }

    try {
      const response = await fetch('/api/v1/auth/refresh', {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
        },
        body: JSON.stringify({
          refresh_token: refreshToken.value,
        }),
      });

      if (!response.ok) {
        console.error('Token refresh failed:', response.status, response.statusText);
        return false;
      }

      const data = await response.json();
      
      // Update tokens
      token.value = data.access_token;
      if (data.refresh_token) {
        refreshToken.value = data.refresh_token;
        localStorage.setItem('refresh_token', data.refresh_token);
      }
      
      localStorage.setItem('access_token', data.access_token);
      
      console.log('Token refreshed successfully');
      return true;
    } catch (error) {
      console.error('Token refresh error:', error);
      return false;
    }
  };

  const initializeAuth = () => {
    const storedToken = localStorage.getItem('access_token');
    const storedRefreshToken = localStorage.getItem('refresh_token');
    const storedUser = localStorage.getItem('user');

    if (storedToken && storedUser) {
      token.value = storedToken;
      refreshToken.value = storedRefreshToken;
      user.value = JSON.parse(storedUser);
      
      // Start token refresh monitoring if we have valid tokens
      if (storedRefreshToken) {
        tokenRefreshService.start();
      }
    }
  };

  // Registration management methods
  const getRegistrations = async (params: {
    status?: string;
    realmName?: string;
    page?: number;
    pageSize?: number;
  }) => {
    try {
      const searchParams = new URLSearchParams();
      if (params.status) searchParams.append('status', params.status);
      if (params.realmName) searchParams.append('realm_name', params.realmName);
      if (params.page) searchParams.append('page', params.page.toString());
      if (params.pageSize) searchParams.append('page_size', params.pageSize.toString());

      const response = await makeAuthenticatedRequest(`/api/v1/admin/registrations?${searchParams}`);

      if (!response.ok) {
        throw new Error('Failed to fetch registrations');
      }

      return await response.json();
    } catch (error) {
      console.error('Get registrations error:', error);
      throw error;
    }
  };

  const getRegistrationStats = async (realmName?: string) => {
    try {
      const params = realmName ? `?realm_name=${realmName}` : '';
      const response = await makeAuthenticatedRequest(`/api/v1/admin/registrations/stats${params}`);

      if (!response.ok) {
        throw new Error('Failed to fetch registration stats');
      }

      return await response.json();
    } catch (error) {
      console.error('Get registration stats error:', error);
      throw error;
    }
  };

  const approveRegistration = async (userId: string, approved: boolean, reason?: string) => {
    try {
      const response = await makeAuthenticatedRequest('/api/v1/admin/registrations/approve', {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json'
        },
        body: JSON.stringify({
          user_id: userId,
          approved,
          reason,
        })
      });

      if (!response.ok) {
        throw new Error('Failed to process registration approval');
      }

      return await response.json();
    } catch (error) {
      console.error('Approve registration error:', error);
      throw error;
    }
  };

  return {
    user,
    token,
    refreshToken,
    isAuthenticated,
    currentUser,
    signIn,
    signUp,
    signOut,
    refreshAuth,
    initializeAuth,
    getRegistrations,
    getRegistrationStats,
    approveRegistration,
  };
});
