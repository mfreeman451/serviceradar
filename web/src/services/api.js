// src/services/api.js
/**
 * API service for making authenticated requests to the backend
 */

/**
 * Makes an authenticated API request
 * @param {string} url - The URL to fetch
 * @param {Object} options - Fetch options
 * @returns {Promise} - Fetch promise
 */

export const apiRequest = async (url, options = {}) => {
    console.log('API request to (originalUrl): ', url);
    const headers = { ...options.headers || {} };
    const proxyUrl = url.startsWith('/api/') ? `/web-api${url}` : url; // Makes /api/status -> /web-api/api/status
    console.log(`API request to: ${proxyUrl}`);
    return fetch(proxyUrl, { ...options, headers });
};

// Debug helper (temporary)
function logRequest(url) {
    console.log(`API request to: ${url}`);
}

/**
 * GET request helper
 * @param {string} url - The URL to fetch
 * @returns {Promise} - Fetch promise that resolves to JSON
 */
export const get = async (url) => {
    const response = await apiRequest(url);
    if (!response.ok) {
        throw new Error(`API request failed: ${response.status}`);
    }
    return response.json();
};

/**
 * POST request helper
 * @param {string} url - The URL to post to
 * @param {Object} data - The data to send
 * @returns {Promise} - Fetch promise that resolves to JSON
 */
export const post = async (url, data) => {
    const response = await apiRequest(url, {
        method: 'POST',
        headers: {
            'Content-Type': 'application/json',
        },
        body: JSON.stringify(data),
    });
    if (!response.ok) {
        throw new Error(`API request failed: ${response.status}`);
    }
    return response.json();
};