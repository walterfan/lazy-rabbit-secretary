/**
 * Utility functions for handling HTTP errors with detailed messages
 */

import { useAuthStore } from '@/stores/authStore';

export interface ApiError {
  message: string;
  status: number;
  details?: string;
}

/**
 * Handle 401 Unauthorized errors by redirecting to home page
 */
export function handleUnauthorized(): void {
  const authStore = useAuthStore();
  
  // Clear auth state
  authStore.signOut();
  
  // Redirect to home page
  if (window.location.pathname !== '/') {
    window.location.href = '/';
  }
}

/**
 * Enhanced fetch function that automatically handles 401 errors
 * Use this for direct API calls that need 401 handling
 */
export async function fetchWith401Handling(url: string, options: RequestInit = {}): Promise<Response> {
  const response = await fetch(url, options);
  
  // If 401, handle it and return the response
  if (response.status === 401) {
    handleUnauthorized();
  }
  
  return response;
}

/**
 * Handles HTTP response errors and returns detailed error messages
 */
export async function handleHttpError(response: Response): Promise<ApiError> {
  let errorMessage = 'An unexpected error occurred';
  let errorDetails = '';

  try {
    const errorData = await response.json();
    errorMessage = errorData.error || errorData.message || errorMessage;
    errorDetails = errorData.details || errorData.message || '';
  } catch {
    // If response is not JSON, use status text
    errorMessage = response.statusText || errorMessage;
  }

  // Add specific messages based on status codes
  switch (response.status) {
    case 401:
      // Handle 401 by redirecting to home page
      handleUnauthorized();
      errorMessage = 'Authentication failed. Redirecting to home page.';
      errorDetails = errorDetails || 'Your session may have expired or the token is invalid.';
      break;
    case 403:
      errorMessage = 'Access denied. You do not have permission to perform this action.';
      errorDetails = errorDetails || 'Please contact your administrator if you believe this is an error.';
      break;
    case 404:
      errorMessage = 'Resource not found.';
      errorDetails = errorDetails || 'The requested item may have been deleted or moved.';
      break;
    case 409:
      errorMessage = 'Conflict. The resource already exists or is in use.';
      errorDetails = errorDetails || 'Please check if the item already exists or try again later.';
      break;
    case 422:
      errorMessage = 'Validation error. Please check your input.';
      errorDetails = errorDetails || 'Some fields may be missing or contain invalid values.';
      break;
    case 429:
      errorMessage = 'Too many requests. Please slow down.';
      errorDetails = errorDetails || 'You are making requests too quickly. Please wait a moment and try again.';
      break;
    case 500:
      errorMessage = 'Server error. Please try again later.';
      errorDetails = errorDetails || 'Something went wrong on our end. We are working to fix it.';
      break;
    case 502:
    case 503:
    case 504:
      errorMessage = 'Service temporarily unavailable.';
      errorDetails = errorDetails || 'The server is currently down for maintenance. Please try again later.';
      break;
    default:
      if (response.status >= 400 && response.status < 500) {
        errorMessage = 'Client error. Please check your request.';
      } else if (response.status >= 500) {
        errorMessage = 'Server error. Please try again later.';
      }
      break;
  }

  return {
    message: errorMessage,
    status: response.status,
    details: errorDetails
  };
}

/**
 * Creates a user-friendly error message with details
 */
export function createUserErrorMessage(error: ApiError): string {
  if (error.details) {
    return `${error.message}\n\nDetails: ${error.details}`;
  }
  return error.message;
}

/**
 * Shows an error alert to the user with detailed information
 */
export function showErrorAlert(error: ApiError): void {
  const message = createUserErrorMessage(error);
  alert(message);
}

/**
 * Logs error details to console for debugging
 */
export function logError(error: ApiError, context?: string): void {
  console.error(`[${context || 'API Error'}] Status: ${error.status}, Message: ${error.message}`);
  if (error.details) {
    console.error(`Details: ${error.details}`);
  }
}
