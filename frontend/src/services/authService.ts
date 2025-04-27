import { apiClient } from "@/lib/api";
import { User } from "@/domain/user";
import { LoginRequestDTO, SignupRequestDTO } from "@/domain/auth.dto";

interface UserApiResponse {
  data: User;
}

const AuthService = {
  /**
   * Sends a signup request to the backend.
   * @param signupData - The user signup data.
   * @returns A promise that resolves with the created User object.
   */
  signup: async (signupData: SignupRequestDTO): Promise<User> => {
    try {
      const response = await apiClient.post<UserApiResponse>("/auth/signup", {
        data: signupData,
      });
      // Return only the user data from the response
      return response.data.data;
    } catch (error) {
      // TODO: Improve error handling/logging
      console.error("AuthService signup error:", error);
      // Re-throw the error so React Query can handle it
      throw error;
    }
  },

  /**
   * Sends a login request to the backend.
   * @param loginData - The user login credentials.
   * @returns A promise that resolves with the logged-in User object.
   */
  login: async (loginData: LoginRequestDTO): Promise<User> => {
    try {
      const response = await apiClient.post<UserApiResponse>("/auth/login", {
        data: loginData,
      });
      // Return only the user data from the response
      return response.data.data;
    } catch (error) {
      // TODO: Improve error handling/logging
      console.error("AuthService login error:", error);
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
