import { create } from "zustand";
import { User } from "@/domain/user"; // Assuming we'll create a domain types file

// Define the shape of the store's state
interface AuthState {
  user: User | null;
  isAuthenticated: boolean;
  setUser: (user: User | null) => void;
  logout: () => void;
  // We might add loading/error states later if needed for initial auth check
}

// Create the store
const useAuthStore = create<AuthState>((set) => ({
  user: null,
  isAuthenticated: false, // Initially not authenticated
  setUser: (user) => set({ user, isAuthenticated: !!user }),
  logout: () => set({ user: null, isAuthenticated: false }),
  // TODO: Add logic here or elsewhere to check initial auth status
  // e.g., by calling a '/auth/me' endpoint when the app loads
}));

export default useAuthStore;

// Removed placeholder declare module block
