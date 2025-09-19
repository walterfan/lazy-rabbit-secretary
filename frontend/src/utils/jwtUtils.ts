/**
 * JWT Token utilities for decoding and checking expiration
 */

export interface JWTPayload {
  sub: string; // subject (user ID)
  iat: number; // issued at
  exp: number; // expiration time
  [key: string]: any; // additional claims
}

/**
 * Decode JWT token without verification (client-side only)
 * Note: This is for reading claims only, not for security validation
 */
export function decodeJWT(token: string): JWTPayload | null {
  try {
    if (!token) return null;
    
    // Split token into parts
    const parts = token.split('.');
    if (parts.length !== 3) {
      console.warn('Invalid JWT token format');
      return null;
    }
    
    // Decode payload (middle part)
    const payload = parts[1];
    
    // Add padding if needed
    const paddedPayload = payload + '='.repeat((4 - payload.length % 4) % 4);
    
    // Decode base64
    const decodedPayload = atob(paddedPayload);
    
    // Parse JSON
    return JSON.parse(decodedPayload);
  } catch (error) {
    console.error('Failed to decode JWT token:', error);
    return null;
  }
}

/**
 * Check if JWT token is expired
 */
export function isTokenExpired(token: string): boolean {
  const payload = decodeJWT(token);
  if (!payload) return true;
  
  const now = Math.floor(Date.now() / 1000);
  return payload.exp <= now;
}

/**
 * Check if JWT token will expire within the given time window (in seconds)
 */
export function isTokenExpiringSoon(token: string, windowSeconds: number = 300): boolean {
  const payload = decodeJWT(token);
  if (!payload) return true;
  
  const now = Math.floor(Date.now() / 1000);
  const expirationTime = payload.exp;
  const timeUntilExpiry = expirationTime - now;
  
  return timeUntilExpiry <= windowSeconds;
}

/**
 * Get time until token expiration in seconds
 */
export function getTimeUntilExpiry(token: string): number {
  const payload = decodeJWT(token);
  if (!payload) return 0;
  
  const now = Math.floor(Date.now() / 1000);
  return Math.max(0, payload.exp - now);
}

/**
 * Get token expiration date
 */
export function getTokenExpirationDate(token: string): Date | null {
  const payload = decodeJWT(token);
  if (!payload) return null;
  
  return new Date(payload.exp * 1000);
}

/**
 * Format time remaining until expiration
 */
export function formatTimeUntilExpiry(token: string): string {
  const seconds = getTimeUntilExpiry(token);
  
  if (seconds <= 0) return 'Expired';
  
  const minutes = Math.floor(seconds / 60);
  const hours = Math.floor(minutes / 60);
  const days = Math.floor(hours / 24);
  
  if (days > 0) return `${days}d ${hours % 24}h`;
  if (hours > 0) return `${hours}h ${minutes % 60}m`;
  if (minutes > 0) return `${minutes}m ${seconds % 60}s`;
  return `${seconds}s`;
}
