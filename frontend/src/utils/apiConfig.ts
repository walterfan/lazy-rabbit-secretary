/**
 * API Configuration Utility
 * Centralized configuration for API endpoints and base URLs
 */

/**
 * Get the API base URL from environment variables or runtime config
 * Falls back to development defaults if not configured
 */
export function getApiBaseUrl(): string {
  // Check for explicit base URL first (build-time)
  if (import.meta.env.VITE_API_BASE_URL) {
    return import.meta.env.VITE_API_BASE_URL;
  }

  // Check for runtime configuration (can be set after build)
  const runtimeConfig = (window as any).__API_CONFIG__;
  if (runtimeConfig?.baseUrl) {
    return runtimeConfig.baseUrl;
  }

  // Build URL from components if available
  const protocol = import.meta.env.VITE_API_PROTOCOL || 'https';
  const host = import.meta.env.VITE_API_HOST || 'lazy-rabbit-studio.top';
  const port = import.meta.env.VITE_API_PORT || '443';

  // Don't include port for standard ports
  if ((protocol === 'https' && port === '443') || (protocol === 'http' && port === '80')) {
    return `${protocol}://${host}`;
  }

  return `${protocol}://${host}:${port}`;
}

/**
 * Get the full API URL for a given endpoint
 * @param endpoint - The API endpoint (e.g., '/api/v1/users')
 * @returns The complete API URL
 */
export function getApiUrl(endpoint: string): string {
  const baseUrl = getApiBaseUrl();
  
  // Ensure endpoint starts with /
  const normalizedEndpoint = endpoint.startsWith('/') ? endpoint : `/${endpoint}`;
  
  return `${baseUrl}${normalizedEndpoint}`;
}

/**
 * Get API configuration for different environments
 */
export const apiConfig = {
  baseUrl: getApiBaseUrl(),
  timeout: 30000, // 30 seconds
  retryAttempts: 3,
  retryDelay: 1000, // 1 second
} as const;

/**
 * Environment-specific configuration
 */
export const isDevelopment = import.meta.env.DEV;
export const isProduction = import.meta.env.PROD;
export const isTest = import.meta.env.MODE === 'test';

/**
 * Log API configuration in development
 */
if (isDevelopment) {
  console.log('API Configuration:', {
    baseUrl: apiConfig.baseUrl,
    environment: import.meta.env.MODE,
    isDevelopment,
    isProduction,
    runtimeConfig: (window as any).__API_CONFIG__,
  });
}
