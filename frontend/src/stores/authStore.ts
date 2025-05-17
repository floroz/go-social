import { create } from "zustand";
import { User } from "@/types/api";

interface AuthState {
  user: User | null;
  token: string | null;
  isAuthenticated: boolean;
  setUser: (user: User | null) => void;
  setToken: (token: string | null) => void;
  logout: () => void;
  initializeAuth: () => void; // For checking token on app load
}

const TOKEN_STORAGE_KEY = "authToken";

// Create the store
const useAuthStore = create<AuthState>((set, get) => ({
  user: null,
  token: null,
  isAuthenticated: false,
  setUser: (user) => {
    set({ user, isAuthenticated: !!user });
    // If logging out (user is null), ensure token is also cleared.
    // If setting a user, it's assumed setToken was called prior or will be.
    if (!user) {
      get().setToken(null); // Ensures token is cleared if user is explicitly set to null
    }
  },
  setToken: (token) => {
    if (token) {
      localStorage.setItem(TOKEN_STORAGE_KEY, token);
      set({ token, isAuthenticated: true });
    } else {
      localStorage.removeItem(TOKEN_STORAGE_KEY);
      set({ token: null, isAuthenticated: false, user: null }); // Also clear user on token removal
    }
  },
  logout: () => {
    get().setToken(null); // This will clear token, user, and isAuthenticated
  },
  initializeAuth: () => {
    const token = localStorage.getItem(TOKEN_STORAGE_KEY);
    if (token) {
      get().setToken(token);
      // TODO: Optionally, verify token with a '/auth/me' endpoint here
      // and fetch user details if token is valid, then call setUser.
      // For now, just setting isAuthenticated based on token presence.
      // If '/auth/me' is implemented, it should call setUser with the user details.
    }
  },
  // TODO: Add logic here or elsewhere to check initial auth status (partially done by initializeAuth)
  // e.g., by calling a '/auth/me' endpoint when the app loads
}));

// Initialize auth status when the store is first loaded/app starts
useAuthStore.getState().initializeAuth();

export default useAuthStore;

// Removed placeholder declare module block
