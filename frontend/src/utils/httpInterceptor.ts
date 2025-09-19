/**
 * HTTP Interceptor for automatic token refresh
 * Handles 401 responses by refreshing tokens and retrying requests
 * Proactively refreshes tokens before they expire
 */

import { useAuthStore } from '@/stores/authStore';
import { isTokenExpiringSoon, isTokenExpired } from '@/utils/jwtUtils';

interface RequestConfig {
  url: string;
  method?: string;
  headers?: Record<string, string>;
  body?: string;
}

interface ResponseConfig {
  ok: boolean;
  status: number;
  statusText: string;
  json: () => Promise<any>;
  text: () => Promise<string>;
}

class HttpInterceptor {
  private isRefreshing = false;
  private failedQueue: Array<{
    resolve: (value: any) => void;
    reject: (error: any) => void;
    config: RequestConfig;
  }> = [];

  /**
   * Get auth store instance (lazy initialization)
   */
  private getAuthStore() {
    return useAuthStore();
  }

  /**
   * Check if token needs refresh and refresh if necessary
   * Returns true if token is valid (either already valid or successfully refreshed)
   */
  private async ensureValidToken(): Promise<boolean> {
    const authStore = this.getAuthStore();
    
    if (!authStore.token) {
      return false;
    }

    // Check if token is already expired
    if (isTokenExpired(authStore.token)) {
      console.log('Token is expired, attempting refresh...');
      return await authStore.refreshAuth();
    }

    // Check if token is expiring soon (within 5 minutes)
    if (isTokenExpiringSoon(authStore.token, 300)) {
      console.log('Token is expiring soon, proactively refreshing...');
      return await authStore.refreshAuth();
    }

    // Token is still valid
    return true;
  }

  /**
   * Process failed requests queue after token refresh
   */
  private processQueue(error: any, token: string | null = null) {
    this.failedQueue.forEach(({ resolve, reject, config }) => {
      if (error) {
        reject(error);
      } else {
        resolve(this.makeRequest(config));
      }
    });
    
    this.failedQueue = [];
  }

  /**
   * Make HTTP request with automatic retry on 401 and proactive token refresh
   */
  async makeRequest(config: RequestConfig): Promise<ResponseConfig> {
    // Ensure we have a valid token before making the request
    const authStore = this.getAuthStore();
    if (authStore.token) {
      const tokenValid = await this.ensureValidToken();
      if (!tokenValid) {
        console.error('Failed to obtain valid token');
        this.redirectToLogin();
        throw new Error('Authentication failed. Please log in again.');
      }
    }

    // Add authorization header if token exists
    const headers: Record<string, string> = {
      'Content-Type': 'application/json',
      ...config.headers,
    };

    if (authStore.token) {
      headers['Authorization'] = `Bearer ${authStore.token}`;
    }

    //log the request url and method
    console.log(`Request: ${config.method} ${config.url}`);
    //console.log(`Headers: ${JSON.stringify(headers)}`);
    //console.log(`Body: ${config.body}`);

    const response = await fetch(config.url, {
      method: config.method || 'GET',
      headers,
      body: config.body,
    });

    // If request is successful, return response
    if (response.ok) {
      return response;
    }

    // Handle 401 Unauthorized - try to refresh token
    if (response.status === 401 && authStore.refreshToken) {
      return this.handleUnauthorized(config);
    }

    // For other errors, return the response as-is
    return response;
  }

  /**
   * Handle 401 Unauthorized responses
   */
  private async handleUnauthorized(originalConfig: RequestConfig): Promise<ResponseConfig> {
    // If already refreshing, queue this request
    if (this.isRefreshing) {
      return new Promise((resolve, reject) => {
        this.failedQueue.push({ resolve, reject, config: originalConfig });
      });
    }

    this.isRefreshing = true;

    try {
      // Try to refresh the token
      const refreshSuccess = await this.getAuthStore().refreshAuth();
      
      if (refreshSuccess) {
        // Token refreshed successfully, process queue and retry original request
        this.processQueue(null);
        return this.makeRequest(originalConfig);
      } else {
        // Refresh failed, redirect to login
        this.processQueue(new Error('Token refresh failed'));
        this.redirectToLogin();
        throw new Error('Authentication failed. Please log in again.');
      }
    } catch (error) {
      // Refresh failed, process queue with error
      this.processQueue(error);
      this.redirectToLogin();
      throw error;
    } finally {
      this.isRefreshing = false;
    }
  }

  /**
   * Redirect to login page
   */
  private redirectToLogin() {
    // Clear auth state
    this.getAuthStore().signOut();
    
    // Redirect to login page
    if (window.location.pathname !== '/login') {
      window.location.href = '/login';
    }
  }

  /**
   * Enhanced fetch function with automatic token refresh
   */
  async fetch(url: string, options: RequestInit = {}): Promise<Response> {
    const config: RequestConfig = {
      url,
      method: options.method || 'GET',
      headers: options.headers as Record<string, string> || {},
      body: options.body as string,
    };

    try {
      const response = await this.makeRequest(config);
      return response as Response;
    } catch (error) {
      throw error;
    }
  }
}

// Create singleton instance
const httpInterceptor = new HttpInterceptor();

/**
 * Enhanced fetch function with automatic token refresh
 * Use this instead of the native fetch function for API calls
 */
export const apiFetch = httpInterceptor.fetch.bind(httpInterceptor);

/**
 * Make authenticated API request with automatic token refresh
 */
export async function makeAuthenticatedRequest(
  url: string,
  options: RequestInit = {}
): Promise<Response> {
  return apiFetch(url, options);
}

/**
 * Helper function to create request config for the interceptor
 */
export function createRequestConfig(
  url: string,
  options: RequestInit = {}
): RequestConfig {
  return {
    url,
    method: options.method || 'GET',
    headers: options.headers as Record<string, string> || {},
    body: options.body as string,
  };
}
