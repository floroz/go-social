import { useMutation, UseMutationOptions } from "@tanstack/react-query";
import AuthService from "@/services/authService";
import useAuthStore from "@/stores/authStore";
import { User } from "@/domain/user";
import { SignupRequestDTO } from "@/domain/auth.dto";

/**
 * Custom hook for handling user signup.
 * Encapsulates the mutation logic using React Query.
 * Allows the calling component to specify side effects via callbacks.
 *
 * @param options - Optional callbacks for onSuccess and onError.
 */
export function useSignup(
  options?: Pick<
    UseMutationOptions<User, Error, SignupRequestDTO>,
    "onSuccess" | "onError"
  >
) {
  const { setUser } = useAuthStore(); // Get setUser from the store

  const { mutate, isPending, error, data } = useMutation<
    User,
    Error,
    SignupRequestDTO
  >({
    mutationFn: AuthService.signup,
    onSuccess: (user, variables, context) => {
      // Core side effect: Update global auth state
      setUser(user);
      options?.onSuccess?.(user, variables, context);
    },
    onError: options?.onError,
  });

  return {
    signup: mutate,
    isLoading: isPending,
    error,
    data,
  };
}
