import { useQueryClient } from "@tanstack/react-query";
import {
  createContext,
  useCallback,
  useContext,
  useEffect,
  useState,
  type ReactNode,
} from "react";
import { ChronoClient } from "./api/chrono/client";
import type { LoginRequest, SignupRequest, User } from "./types/auth";

let logoutFn: (() => Promise<void>) | null = null;

export function registerLogout(fn: () => Promise<void>) {
  logoutFn = fn;
}

export async function logoutOutsideReact() {
  if (logoutFn) {
    await logoutFn();
  }
}

export interface AuthContext {
  isAuthenticated: boolean;
  login: (data: LoginRequest) => Promise<void>;
  signup: (data: SignupRequest) => Promise<void>;
  logout: () => Promise<void>;
  userId: number | null;
  getUser: () => Promise<User | null>;
}

const AuthContext = createContext<AuthContext | null>(null);

export function AuthProvider({ children }: { children: ReactNode }) {
  const uid = localStorage.getItem("user");
  const [userId, setUserId] = useState<number | null>(
    uid ? Number.parseInt(uid) : null,
  );

  const queryClient = useQueryClient();

  const [isAuthenticated, setIsAuthenticated] = useState(!!uid);
  const chrono = new ChronoClient();

  const getUser = useCallback(async () => {
    if (!isAuthenticated || !userId) return null;
    const u = await queryClient.ensureQueryData({
      queryKey: ["user", userId],
      queryFn: () => chrono.users.getUserById(userId),
      staleTime: 1000 * 60 * 60 * 6, // 6h
      gcTime: 1000 * 60 * 60 * 7, // 7h
      retry: false,
    });
    return u;
  }, []);

  const logout = useCallback(async () => {
    try {
      await chrono.auth.logout();
    } finally {
      localStorage.removeItem("user");
      setUserId(null);
      setIsAuthenticated(false);
      queryClient.clear();
      window.location.href = "/login";
    }
  }, []);

  const login = useCallback(async (data: LoginRequest) => {
    const user = (await chrono.auth.login(data)).data;
    localStorage.setItem("user", user.id);
    setUserId(user.id);
    setIsAuthenticated(true);
  }, []);

  const signup = useCallback(async (data: SignupRequest) => {
    const user = (await chrono.auth.signup(data)).data;
    localStorage.setItem("user", user.id);
    setUserId(user.id);
    setIsAuthenticated(true);
  }, []);

  useEffect(() => {
    registerLogout(logout);
  }, [logout]);

  return (
    <AuthContext.Provider
      value={{
        isAuthenticated,
        userId,
        login,
        logout,
        getUser,
        signup,
      }}
    >
      {children}
    </AuthContext.Provider>
  );
}

export function useAuth() {
  const context = useContext(AuthContext);
  if (!context) {
    throw new Error("useAuth must be used within an AuthProvider");
  }
  return context;
}
