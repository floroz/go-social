import React from "react";
import { BrowserRouter } from "react-router";
import { QueryClient, QueryClientProvider } from '@tanstack/react-query';

// Create a new QueryClient instance for each test wrapper instance
// to ensure test isolation.
const createTestQueryClient = () => new QueryClient({
  defaultOptions: {
    queries: {
      retry: false, // Disable retries for tests
    },
  },
});

const PlaywrightTestWrapper: React.FC<{ children: React.ReactNode }> = ({ children }) => {
  // Create a fresh client for each mount
  const queryClient = createTestQueryClient();

  return (
    <QueryClientProvider client={queryClient}>
      <BrowserRouter>
        {children}
      </BrowserRouter>
    </QueryClientProvider>
  );
};

export default PlaywrightTestWrapper;
