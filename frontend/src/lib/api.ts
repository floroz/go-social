import axios from "axios";
import config from "@/config"; // Import the centralized config
import useAuthStore from "@/stores/authStore";

const apiClient = axios.create({
  baseURL: config.apiBaseUrl,
  // We can configure interceptors here later for auth tokens, etc.
  headers: {
    "Content-Type": "application/json",
  },
  // Important for cookies to be sent/received across domains during development
  // if frontend and backend are on different ports (e.g., 5173 vs 8080)
  withCredentials: true,
});

// Request interceptor to add JWT token to headers
apiClient.interceptors.request.use(
  (config) => {
    // Import useAuthStore here to avoid circular dependencies if api.ts is imported by store.ts
    // This is a common pattern for accessing Zustand store outside React components.
    const token = useAuthStore.getState().token;
    if (token) {
      config.headers.Authorization = `Bearer ${token}`;
    }
    return config;
  },
  (error) => {
    return Promise.reject(error);
  }
);

// Example response interceptor (can be uncommented and expanded later)
/*
apiClient.interceptors.response.use(response => {
  return response;
}, error => {
  // Handle errors globally, e.g., check for 401 and try token refresh
  return Promise.reject(error);
});
*/

export default apiClient;
