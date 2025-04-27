import apiClient from "@/lib/api";
import { SignupRequest, LoginRequest, User, LoginResponse } from "@/types/api";

const AuthService = {
  /**
   * Sends a signup request to the backend.
   * @param signupData - The user signup data.
   * @returns A promise that resolves with the created User object.
   */
  signup: async (signupData: SignupRequest): Promise<User> => {
    // eslint-disable-next-line no-useless-catch
    try {
      const response = await apiClient.post<User>(
        "/v1/auth/signup",
        signupData
      );
      return response.data; // Return the user data directly
    } catch (error) {
      // TODO: Improve error handling/logging
      // console.error("AuthService signup error:", error);
      // Re-throw the error so React Query can handle it
      throw error;
    }
  },

  /**
   * Sends a login request to the backend.
   * @param loginData - The user login credentials.
   * @returns A promise that resolves with the login response containing the token.
   */
  login: async (loginData: LoginRequest): Promise<LoginResponse> => {
    // eslint-disable-next-line no-useless-catch
    try {
      const response = await apiClient.post<LoginResponse>(
        "/v1/auth/login",
        loginData
      );
      return response.data; // Return the full login response (e.g., { token: "..." })
    } catch (error) {
      // TODO: Improve error handling/logging
      // console.error("AuthService login error:", error);
      // Re-throw the error so React Query can handle it
      throw error;
    }
  },

  // TODO: Add logout function if needed (e.g., call /auth/logout endpoint)
  // logout: async (): Promise<void> => { ... }

  // TODO: Add function to check current auth status (e.g., call /auth/me)
  // checkAuth: async (): Promise<User | null> => { ... }
};

export default AuthService;
