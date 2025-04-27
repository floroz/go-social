import { http, HttpResponse } from "msw";
import { User } from "@/domain/user"; // Import User type
// Removed unused SignupRequestDTO import
// import { SignupRequestDTO } from "@/domain/auth.dto";

// Define a base URL for the API (should match apiClient config)
const API_BASE_URL = "http://localhost:8080/api/v1";

// Mock user data for successful responses
const mockUser: User = {
  id: 1,
  first_name: "Mock",
  last_name: "User",
  username: "mockuser",
  email: "mock@example.com",
  created_at: new Date().toISOString(),
  updated_at: new Date().toISOString(),
  last_login: null, // Or new Date().toISOString() if needed
};

export const handlers = [
  // Mock the POST /auth/signup endpoint
  http.post(`${API_BASE_URL}/auth/signup`, async (/*{ request }*/) => {
    // Removed unused request param
    // You can optionally inspect the request body if needed
    // const body = await request.json() as { data: SignupRequestDTO };
    // console.log('MSW intercepted signup:', body.data);

    // Simulate a successful signup
    // Backend returns { data: User }
    return HttpResponse.json({ data: mockUser }, { status: 201 });

    // --- Example Error Response ---
    // return HttpResponse.json(
    //   { errors: [{ message: 'Username already exists' }] },
    //   { status: 400 }
    // );
  }),

  // TODO: Add handlers for login, posts, comments etc. later
  // http.post(`${API_BASE_URL}/auth/login`, async ({ request }) => { ... })
];
