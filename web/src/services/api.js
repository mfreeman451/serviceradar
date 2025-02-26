// src/services/api.js
/**
 * API service for making authenticated requests to the backend
 */

// Get API key from environment variables
const getApiKey = () => {
    // For Vite applications, environment variables should be prefixed with VITE_
    return import.meta.env.VITE_API_KEY || '';
};

/**
 * Makes an authenticated API request
 * @param {string} url - The URL to fetch
 * @param {Object} options - Fetch options
 * @returns {Promise} - Fetch promise
 */
export const apiRequest = async (url, options = {}) => {
    const apiKey = getApiKey();

    // Set up headers with API key if available
    const headers = {
        ...options.headers || {},
    };

    if (apiKey) {
        headers['X-API-Key'] = apiKey;
    }

    const proxyUrl = url.replace(/^\/api\//, '/web-api/');

    // Make the request with the headers
    return fetch(proxyUrl, {
        ...options,
        headers,
    });
};

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