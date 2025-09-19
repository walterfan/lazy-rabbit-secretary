/**
 * Token Refresh Service
 * Handles automatic token refresh based on expiration time
 */

import { useAuthStore } from '@/stores/authStore';
import { getTimeUntilExpiry, isTokenExpired, formatTimeUntilExpiry } from '@/utils/jwtUtils';

export class TokenRefreshService {
  private refreshInterval: number | null = null;
  private readonly REFRESH_WINDOW_SECONDS = 300; // 5 minutes before expiry
  private readonly MIN_REFRESH_INTERVAL = 60; // Minimum 1 minute between checks

  /**
   * Start automatic token refresh monitoring
   */
  start(): void {
    this.stop(); // Stop any existing interval
    
    const authStore = useAuthStore();
    
    if (!authStore.token || !authStore.refreshToken) {
      console.log('No tokens available for refresh monitoring');
      return;
    }

    // Check token status immediately
    this.checkAndRefreshToken();

    // Set up interval for periodic checks
    this.refreshInterval = window.setInterval(() => {
      this.checkAndRefreshToken();
    }, this.MIN_REFRESH_INTERVAL * 1000);

    console.log('Token refresh monitoring started');
  }

  /**
   * Stop automatic token refresh monitoring
   */
  stop(): void {
    if (this.refreshInterval) {
      clearInterval(this.refreshInterval);
      this.refreshInterval = null;
      console.log('Token refresh monitoring stopped');
    }
  }

  /**
   * Check token status and refresh if necessary
   */
  private async checkAndRefreshToken(): Promise<void> {
    const authStore = useAuthStore();
    
    if (!authStore.token || !authStore.refreshToken) {
      this.stop();
      return;
    }

    // Check if token is expired
    if (isTokenExpired(authStore.token)) {
      console.log('Token is expired, attempting refresh...');
      const success = await authStore.refreshAuth();
      if (!success) {
        console.error('Failed to refresh expired token');
        this.stop();
        return;
      }
      console.log('Token refreshed successfully');
      return;
    }

    // Check if token is expiring soon
    const timeUntilExpiry = getTimeUntilExpiry(authStore.token);
    
    if (timeUntilExpiry <= this.REFRESH_WINDOW_SECONDS) {
      console.log(`Token expires in ${formatTimeUntilExpiry(authStore.token)}, refreshing...`);
      const success = await authStore.refreshAuth();
      if (!success) {
        console.error('Failed to refresh token');
        this.stop();
        return;
      }
      console.log('Token refreshed proactively');
    } else {
      // Log remaining time for debugging
      console.log(`Token expires in ${formatTimeUntilExpiry(authStore.token)}`);
    }
  }

  /**
   * Get current token status
   */
  getTokenStatus(): {
    isValid: boolean;
    timeUntilExpiry: number;
    formattedTime: string;
    needsRefresh: boolean;
  } {
    const authStore = useAuthStore();
    
    if (!authStore.token) {
      return {
        isValid: false,
        timeUntilExpiry: 0,
        formattedTime: 'No token',
        needsRefresh: false
      };
    }

    const timeUntilExpiry = getTimeUntilExpiry(authStore.token);
    const needsRefresh = timeUntilExpiry <= this.REFRESH_WINDOW_SECONDS;

    return {
      isValid: !isTokenExpired(authStore.token),
      timeUntilExpiry,
      formattedTime: formatTimeUntilExpiry(authStore.token),
      needsRefresh
    };
  }

  /**
   * Force refresh token
   */
  async forceRefresh(): Promise<boolean> {
    const authStore = useAuthStore();
    
    if (!authStore.refreshToken) {
      console.error('No refresh token available');
      return false;
    }

    console.log('Force refreshing token...');
    const success = await authStore.refreshAuth();
    
    if (success) {
      console.log('Token force refresh successful');
    } else {
      console.error('Token force refresh failed');
    }
    
    return success;
  }
}

// Create singleton instance
export const tokenRefreshService = new TokenRefreshService();
