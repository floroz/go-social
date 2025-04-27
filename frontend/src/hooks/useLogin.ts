import { useMutation, UseMutationOptions } from "@tanstack/react-query";
import AuthService from "@/services/authService";
import useAuthStore from "@/stores/authStore";
import { User } from "@/domain/user";
import { LoginRequestDTO } from "@/domain/auth.dto";

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
    UseMutationOptions<User, Error, LoginRequestDTO>,
    "onSuccess" | "onError"
  >
) {
  const { setUser } = useAuthStore();

  const { mutate, isPending, error, data } = useMutation<
    User,
    Error,
    LoginRequestDTO
  >({
    mutationFn: AuthService.login,
    onSuccess: (user, variables, context) => {
      // Core side effect: Update global auth state
      setUser(user);
      console.log("Login successful:", user);
      options?.onSuccess?.(user, variables, context);
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
