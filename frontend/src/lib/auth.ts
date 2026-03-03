/**
 * getAuthToken — reads the JWT access token from localStorage.
 * Returns null if not found or if the stored value is "null"/"undefined".
 */
export function getAuthToken(): string | null {
    if (typeof window === 'undefined') return null;
    const raw = localStorage.getItem('access_token');
    if (!raw || raw === 'null' || raw === 'undefined' || raw.trim() === '') {
        return null;
    }
    return raw;
}

/**
 * storeAuthToken — saves the JWT to localStorage.
 * Call this after any successful login/register.
 */
export function storeAuthToken(token: string | undefined | null) {
    if (!token || token === 'null' || token === 'undefined') return;
    localStorage.setItem('access_token', token);
}

/**
 * authHeaders — returns the Authorization header object if a token exists,
 * otherwise returns an empty object (lets the request proceed without it,
 * which the backend will reject with a clear 401).
 */
export function authHeaders(): Record<string, string> {
    const token = getAuthToken();
    if (!token) return {};
    return { Authorization: `Bearer ${token}` };
}

/**
 * clearAuthToken — removes the token on logout.
 */
export function clearAuthToken() {
    localStorage.removeItem('access_token');
}
