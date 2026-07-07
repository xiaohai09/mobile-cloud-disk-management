/**
 * Security utilities for frontend rendering and input sanitization.
 */

/**
 * Escape HTML special characters to prevent XSS when rendering user/admin content.
 * Note: Vue template bindings like {{ }} already escape by default.
 * Use this helper only when intentionally rendering trusted-but-raw HTML text.
 */
export function escapeHtml(value: string): string {
  return value
    .replace(/&/g, '&amp;')
    .replace(/</g, '&lt;')
    .replace(/>/g, '&gt;')
    .replace(/"/g, '&quot;')
    .replace(/'/g, '&#039;')
}

/**
 * Purge known sensitive keys from localStorage to reduce token/credential exposure.
 */
export function purgeSensitiveStorage(): void {
  const keys = [
    'user',
    'token',
    'access_token',
    'refresh_token',
    'csrf_token',
    'session',
  ]

  for (const key of keys) {
    try {
      localStorage.removeItem(key)
    } catch {
      // ignore storage access errors in restricted contexts
    }
  }
}
