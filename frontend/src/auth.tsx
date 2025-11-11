import * as React from "react";

import { useRouter } from "@tanstack/react-router";
import {
  login as chrono_login,
  logout as chrono_logout,
} from "./api/chrono/auth";
import type { LoginRequest, User } from "./types/auth";
import { useEffect } from "react";
import { getUserById } from "./api/chrono/users";

export interface AuthContext {
  isAuthenticated: () => boolean;
  login: (data: LoginRequest) => Promise<void>;
  logout: () => Promise<void>;
  user: User | null;
}

const AuthContext = React.createContext<AuthContext | null>(null);

export function AuthProvider({ children }: { children: React.ReactNode }) {
  const [user, setUser] = React.useState<User | null>(null);
  const router = useRouter();

  useEffect(() => {
    let cancelled = false;

    const fn = async () => {
      const local = localStorage.getItem("user");
      if (local) {
        const user = await getUserById(Number.parseInt(local));
        if (!cancelled) setUser(user);
        console.log(user);
      }
    };

    fn();

    return () => {
      cancelled = true;
    };
  }, []);

  const logout = React.useCallback(async () => {
    await chrono_logout();
    localStorage.removeItem("user");
    setUser(null);
    await router.invalidate();
    router.navigate({ to: "/login" });
  }, []);

  const login = React.useCallback(async (data: LoginRequest) => {
    const user = (await chrono_login(data)).data;
    localStorage.setItem("user", String(user.id));
    setUser(user);
  }, []);

  const isAuthenticated = () => !!user;

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
