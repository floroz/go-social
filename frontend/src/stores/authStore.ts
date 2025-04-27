import { create } from "zustand";
import { User } from "@/types/api";

interface AuthState {
  user: User | null;
  isAuthenticated: boolean;
  setUser: (user: User | null) => void;
  logout: () => void;
}

// Create the store
const useAuthStore = create<AuthState>((set) => ({
  user: null,
  isAuthenticated: false,
  setUser: (user) => set({ user, isAuthenticated: !!user }),
  logout: () => set({ user: null, isAuthenticated: false }),
  // TODO: Add logic here or elsewhere to check initial auth status
  // e.g., by calling a '/auth/me' endpoint when the app loads
}));

export default useAuthStore;

// Removed placeholder declare module block
