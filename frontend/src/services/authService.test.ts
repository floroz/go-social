import { describe, it, expect } from "vitest";
import { http, HttpResponse } from "msw";
import { server } from "@/mocks/server";
import AuthService from "./authService";
import { SignupRequestDTO, LoginRequestDTO } from "@/domain/auth.dto";
import { User } from "@/domain/user";
import config from "@/config";

const API_BASE_URL = config.apiBaseUrl;

// Mock user data for successful responses
const mockUser: User = {
  id: 1,
  first_name: "Mock",
  last_name: "User",
  username: "mockuser",
  email: "mock@example.com",
  created_at: new Date().toISOString(),
  updated_at: new Date().toISOString(),
  last_login: null,
};

describe("AuthService", () => {
  describe("signup", () => {
    it("should return user data on successful signup", async () => {
      const signupData: SignupRequestDTO = {
        first_name: "Test",
        last_name: "User",
        username: "testuser",
        email: "test@example.com",
        password: "password123",
      };

      // MSW handler in handlers.ts already mocks a successful response
      const user = await AuthService.signup(signupData);

      // Compare relevant fields, ignoring exact timestamps
      expect(user.id).toEqual(mockUser.id);
      expect(user.first_name).toEqual(mockUser.first_name);
      expect(user.last_name).toEqual(mockUser.last_name);
      expect(user.username).toEqual(mockUser.username);
      expect(user.email).toEqual(mockUser.email);
      expect(user.last_login).toEqual(mockUser.last_login);
      // Optionally check that timestamps are valid date strings
      expect(user.created_at).toEqual(expect.any(String));
      expect(user.updated_at).toEqual(expect.any(String));
    });

    it("should throw an error on signup failure (e.g., 400)", async () => {
      const signupData: SignupRequestDTO = {
        /* ... */
      } as SignupRequestDTO; // Data doesn't matter for error test
      const errorMessage = "Username already exists";

      // Override the default handler for this specific test
      server.use(
        http.post(`${API_BASE_URL}/auth/signup`, () => {
          return HttpResponse.json(
            { errors: [{ message: errorMessage }] },
            { status: 400 }
          );
        })
      );

      await expect(AuthService.signup(signupData)).rejects.toThrow();
    });

    it("should throw an error on network or server error (e.g., 500)", async () => {
      const signupData: SignupRequestDTO = {
        /* ... */
      } as SignupRequestDTO;

      server.use(
        http.post(`${API_BASE_URL}/auth/signup`, () => {
          return new HttpResponse(null, {
            status: 500,
            statusText: "Internal Server Error",
          });
        })
      );

      await expect(AuthService.signup(signupData)).rejects.toThrow();
    });
  });

  describe("login", () => {
    // Mock successful login handler
    const loginSuccessHandler = http.post(`${API_BASE_URL}/auth/login`, () => {
      return HttpResponse.json({ data: mockUser }, { status: 200 });
    });

    // Mock failed login handler (401 Unauthorized)
    const loginFailHandler = http.post(`${API_BASE_URL}/auth/login`, () => {
      return HttpResponse.json(
        { errors: [{ message: "Invalid credentials" }] },
        { status: 401 }
      );
    });

    it("should return user data on successful login", async () => {
      const loginData: LoginRequestDTO = {
        email: "test@example.com",
        password: "password123",
      };

      server.use(loginSuccessHandler); // Use the success handler for this test
      const user = await AuthService.login(loginData);

      // Compare relevant fields, ignoring exact timestamps
      expect(user.id).toEqual(mockUser.id);
      expect(user.first_name).toEqual(mockUser.first_name);
      expect(user.last_name).toEqual(mockUser.last_name);
      expect(user.username).toEqual(mockUser.username);
      expect(user.email).toEqual(mockUser.email);
      expect(user.last_login).toEqual(mockUser.last_login);
      expect(user.created_at).toEqual(expect.any(String));
      expect(user.updated_at).toEqual(expect.any(String));
    });

    it("should throw an error on login failure (e.g., 401)", async () => {
      const loginData: LoginRequestDTO = {
        /* ... */
      } as LoginRequestDTO;

      server.use(loginFailHandler); // Use the failure handler for this test
      await expect(AuthService.login(loginData)).rejects.toThrow();
    });

    it("should throw an error on network or server error (e.g., 500)", async () => {
      const loginData: LoginRequestDTO = {
        /* ... */
      } as LoginRequestDTO;

      server.use(
        http.post(`${API_BASE_URL}/auth/login`, () => {
          return new HttpResponse(null, { status: 500 });
        })
      );
      await expect(AuthService.login(loginData)).rejects.toThrow();
    });
  });
});
