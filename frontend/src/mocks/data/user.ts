import { User } from "@/types/api";

// Mock user data for use in handlers and tests
export const mockUser: User = {
  id: 1,
  first_name: "Mock",
  last_name: "User",
  username: "mockuser",
  email: "mock@example.com",
  created_at: new Date().toISOString(),
  updated_at: new Date().toISOString(),
  last_login: null, // Or new Date().toISOString() if needed
};
