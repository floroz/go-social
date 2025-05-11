import { describe, it, expect } from "vitest";
import { http, HttpResponse } from "msw";
import { server } from "@/mocks/server";
import AuthService from "./authService";
import { SignupRequest, LoginRequest, LoginSuccessResponse } from "@/types/api";
import { mockUser } from "@/mocks/data/user";
import config from "@/config";

const API_BASE_URL = config.apiBaseUrl;

describe("AuthService", () => {
  describe("signup", () => {
    it("should return user data on successful signup", async () => {
      const signupData: SignupRequest = {
        first_name: "Test",
        last_name: "User",
        username: "testuser",
        email: "test@example.com",
        password: "password123",
      };

      const user = await AuthService.signup(signupData);

      // Compare relevant fields
      expect(user.id).toEqual(mockUser.id);
      expect(user.first_name).toEqual(mockUser.first_name);
      expect(user.last_name).toEqual(mockUser.last_name);
      expect(user.username).toEqual(mockUser.username);
      expect(user.email).toEqual(mockUser.email);
      expect(user.last_login).toEqual(mockUser.last_login);
      expect(user.created_at).toEqual(mockUser.created_at);
      expect(user.updated_at).toEqual(mockUser.updated_at);
    });

    it("should throw an error on signup failure (e.g., 409 Conflict)", async () => {
      const signupData: SignupRequest = {
        first_name: "Test",
        last_name: "Fail",
        username: "failuser",
        email: "fail@example.com",
        password: "password123",
      };
      const errorMessage = "Username already exists";
      const mockErrorResponse = { error: errorMessage };

      server.use(
        http.post(`${API_BASE_URL}/v1/auth/signup`, () => {
          return HttpResponse.json(mockErrorResponse, { status: 409 });
        })
      );

      await expect(AuthService.signup(signupData)).rejects.toThrow();
    });

    it("should throw an error on network or server error (e.g., 500)", async () => {
      const signupData: SignupRequest = {
        first_name: "Test",
        last_name: "Error",
        username: "erroruser",
        email: "error@example.com",
        password: "password123",
      };

      server.use(
        http.post(`${API_BASE_URL}/v1/auth/signup`, () => {
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
    const mockLoginResponse = { token: "mock-jwt-token" };
    const loginSuccessHandler = http.post(
      `${API_BASE_URL}/v1/auth/login`,
      () => {
        return HttpResponse.json(
          { data: mockLoginResponse } satisfies LoginSuccessResponse,
          { status: 200 }
        );
      }
    );

    const loginFailHandler = http.post(`${API_BASE_URL}/v1/auth/login`, () => {
      const mockError = { error: "Invalid credentials" };
      return HttpResponse.json(mockError, { status: 401 });
    });

    it("should return LoginResponse on successful login", async () => {
      const loginData: LoginRequest = {
        email: "test@example.com",
        password: "password123",
      };

      server.use(loginSuccessHandler);
      const response = await AuthService.login(loginData);

      expect(response.data.token).toEqual(mockLoginResponse.token);
    });

    it("should throw an error on login failure (e.g., 401)", async () => {
      const loginData: LoginRequest = {
        email: "wrong@example.com",
        password: "wrongpassword",
      };

      server.use(loginFailHandler);
      await expect(AuthService.login(loginData)).rejects.toThrow();
    });

    it("should throw an error on network or server error (e.g., 500)", async () => {
      const loginData: LoginRequest = {
        email: "error@example.com",
        password: "errorpassword",
      };

      server.use(
        http.post(`${API_BASE_URL}/v1/auth/login`, () => {
          return new HttpResponse(null, { status: 500 });
        })
      );
      await expect(AuthService.login(loginData)).rejects.toThrow();
    });
  });
});
