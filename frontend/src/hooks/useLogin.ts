import { useMutation, UseMutationOptions } from "@tanstack/react-query";
import AuthService from "@/services/authService";
import type { LoginResponse, LoginRequest } from "@/types/api";

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
    UseMutationOptions<LoginResponse, Error, LoginRequest>,
    "onSuccess" | "onError"
  >
) {
  // const { setUser } = useAuthStore(); // setUser is not used here anymore

  const { mutate, isPending, error, data } = useMutation<
    LoginResponse,
    Error,
    LoginRequest
  >({
    mutationFn: AuthService.login,
    onSuccess: (loginResponse, variables, context) => {
      // Core side effect: Update global auth state
      // setUser(user); // This is incorrect here as loginResponse is not a User object.
      // Token handling and user fetching should be managed by the component or another service.
      console.log("Login successful, token received:", loginResponse);
      options?.onSuccess?.(loginResponse, variables, context);
    },
    // Pass through the onError callback directly
    onError: (error, variables, context) => {
      console.error("Login failed:", error.message);
      // Call the optional callback provided by the component
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
