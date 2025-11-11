import * as React from "react";

import { useRouter } from "@tanstack/react-router";
import {
  login as chrono_login,
  logout as chrono_logout,
} from "./api/chrono/auth";
import type { LoginRequest, User } from "./types/auth";

export interface AuthContext {
  isAuthenticated: boolean;
  login: (data: LoginRequest) => Promise<void>;
  logout: () => Promise<void>;
  user: User | null;
}

const AuthContext = React.createContext<AuthContext | null>(null);

export function AuthProvider({ children }: { children: React.ReactNode }) {
  const local = localStorage.getItem("user");
  const parsed = local ? JSON.parse(local) : null;
  const [user, setUser] = React.useState<User | null>(parsed);
  const isAuthenticated = !!user;
  const router = useRouter();

  const logout = React.useCallback(async () => {
    await chrono_logout();
    localStorage.removeItem("user");
    setUser(null);
    await router.invalidate();
    router.navigate({ to: "/login" });
  }, []);

  const login = React.useCallback(async (data: LoginRequest) => {
    const user = (await chrono_login(data)).data;
    localStorage.setItem("user", JSON.stringify(user));
    setUser(user);
  }, []);

  return (
    <AuthContext.Provider value={{ isAuthenticated, user, login, logout }}>
      {children}
    </AuthContext.Provider>
  );
}

export function useAuth() {
  const context = React.useContext(AuthContext);
  if (!context) {
    throw new Error("useAuth must be used within an AuthProvider");
  }
  return context;
}
