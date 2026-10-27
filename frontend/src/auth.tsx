import { useQueryClient } from "@tanstack/react-query";
import { createContext, useCallback, useContext, useEffect, useState, type ReactNode } from "react";
import { chrono } from "./api/chrono/client";
import { registerLogout } from "./lib/logout";
import type { LoginRequest, SignupRequest, User } from "./types/auth";

export interface AuthContext {
  isAuthenticated: boolean;
  login: (data: LoginRequest) => Promise<void>;
  signup: (data: SignupRequest) => Promise<void>;
  logout: () => Promise<void>;
  userId: number | null;
}

const authContext = createContext<AuthContext | null>(null);

export function AuthProvider({ children }: { children: ReactNode }) {
  const uid = localStorage.getItem("user");
  const [userId, setUserId] = useState<number | null>(uid ? Number.parseInt(uid) : null);

  const queryClient = useQueryClient();

  const [isAuthenticated, setIsAuthenticated] = useState(!!uid);

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
  }, [queryClient]);

  const login = useCallback(async (data: LoginRequest) => {
    const user = (await chrono.auth.login(data)).data as User;
    localStorage.setItem("user", String(user.id));
    setUserId(user.id);
    setIsAuthenticated(true);
  }, []);

  const signup = useCallback(async (data: SignupRequest) => {
    const user = (await chrono.auth.signup(data)).data as User;
    localStorage.setItem("user", String(user.id));
    setUserId(user.id);
    setIsAuthenticated(true);
  }, []);

  useEffect(() => {
    registerLogout(logout);
  }, [logout]);

  return (
    <authContext.Provider
      value={{
        isAuthenticated,
        userId,
        login,
        logout,
        signup,
      }}
    >
      {children}
    </authContext.Provider>
  );
}

export function useAuth() {
  const context = useContext(authContext);
  if (!context) {
    throw new Error("useAuth must be used within an AuthProvider");
  }
  return context;
}
