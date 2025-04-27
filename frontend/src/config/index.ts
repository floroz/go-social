interface AppConfig {
  apiBaseUrl: string;
}

const getConfigValue = (key: string, defaultValue?: string): string => {
  const value = import.meta.env[key];
  if (value === undefined) {
    if (defaultValue !== undefined) {
      console.warn(
        `Environment variable ${key} not set, using default value: ${defaultValue}`
      );
      return defaultValue;
    }
    // Throw an error for required variables without defaults
    throw new Error(`Missing required environment variable: ${key}`);
  }
  return String(value); // Ensure value is a string
};

const config: AppConfig = {
  // Use VITE_ prefix as required by Vite to expose to client
  apiBaseUrl: getConfigValue("VITE_API_BASE_URL", "http://localhost:8080/api"),
};

// Validate required configurations (optional but recommended)
if (!config.apiBaseUrl) {
  console.error("API Base URL configuration is missing or invalid.");
  // Potentially throw error here depending on how critical it is
}

export default config;
