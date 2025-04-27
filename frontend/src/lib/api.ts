import axios from "axios";

// TODO: Make base URL configurable via environment variables
const API_BASE_URL = "http://localhost:8080/api/v1";

const apiClient = axios.create({
  baseURL: API_BASE_URL,
  // We can configure interceptors here later for auth tokens, etc.
  headers: {
    "Content-Type": "application/json",
  },
  // Important for cookies to be sent/received across domains during development
  // if frontend and backend are on different ports (e.g., 5173 vs 8080)
  withCredentials: true,
});

// Example interceptor (can be uncommented and expanded later)
/*
apiClient.interceptors.request.use(config => {
  // Maybe add auth token here if not using HttpOnly cookies
  return config;
}, error => {
  return Promise.reject(error);
});

apiClient.interceptors.response.use(response => {
  return response;
}, error => {
  // Handle errors globally, e.g., check for 401 and try token refresh
  return Promise.reject(error);
});
*/

export { apiClient };
