import { useMutation, UseMutationOptions } from "@tanstack/react-query";
import AuthService from "@/services/authService";
import useAuthStore from "@/stores/authStore";
import apiClient from "@/lib/api";
import type {
  LoginRequest,
  User,
  LoginSuccessResponse as LoginMutationResponse,
  GetUserProfileSuccessResponse as UserProfileResponse,
} from "@/types/api";

/**
 * Custom hook for handling user login.
 * Encapsulates the mutation logic using React Query.
 * Updates global auth state on success.
 * Allows the calling component to specify additional side effects via callbacks.
 *
 * @param options - Optional callbacks for onSuccess and onError.
 */
export function useLogin(
  options?: Pick<
    UseMutationOptions<LoginMutationResponse, Error, LoginRequest>,
    "onSuccess" | "onError"
  >
) {
  const { setToken, setUser } = useAuthStore();

  const { mutate, isPending, error, data } = useMutation<
    LoginMutationResponse,
    Error,
    LoginRequest
  >({
    mutationFn: AuthService.login,
    onSuccess: async (loginMutationResponse, variables, context) => {
      console.log(
        "Login mutation successful, token received:",
        loginMutationResponse.data.token
      );
      setToken(loginMutationResponse.data.token);

      try {
        // After token is set, interceptor will use it for this call
        const userProfileApiResponse = await apiClient.get<UserProfileResponse>(
          "/v1/users"
        );
        const userDetails: User = userProfileApiResponse.data.data;
        setUser(userDetails);
        console.log(
          "User profile fetched and auth state updated:",
          userDetails
        );

        // Call the original onSuccess from the component if it exists
        // Pass the original loginMutationResponse as it might be expected by the component
        options?.onSuccess?.(loginMutationResponse, variables, context);
      } catch (fetchUserError) {
        console.error(
          "Failed to fetch user profile after login:",
          fetchUserError
        );
        // If fetching user fails, it's a partial success/failure.
        // We have a token, but no user details.
        // For now, clear the token and treat as login failure for simplicity.
        // Alternatively, could leave token and try fetching user later, or show a specific error.
        setToken(null); // This also clears user and isAuthenticated

        // Call the original onError from the component if it exists
        // It's tricky what to pass as 'error' here.
        // For now, creating a new Error object.
        const effectiveError =
          fetchUserError instanceof Error
            ? fetchUserError
            : new Error("Failed to fetch user profile after login.");
        options?.onError?.(effectiveError, variables, context);
      }
    },
    onError: (error, variables, context) => {
      console.error("Login mutation failed:", error.message);
      setToken(null); // Ensure token is cleared on login mutation failure
      options?.onError?.(error, variables, context);
    },
  });

  return {
    login: mutate,
    isLoading: isPending,
    error,
    data,
  };
}
