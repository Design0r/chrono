import {
  createContext,
  useCallback,
  useContext,
  useState,
  type ReactNode,
} from "react";
import { ChronoClient } from "./api/chrono/client";
import type { LoginRequest } from "./types/auth";

export interface AuthContext {
  isAuthenticated: boolean;
  login: (data: LoginRequest) => Promise<void>;
  logout: () => Promise<void>;
  userId: number | null;
}

const AuthContext = createContext<AuthContext | null>(null);

export function AuthProvider({ children }: { children: ReactNode }) {
  const uid = localStorage.getItem("user");
  const [userId, setUserId] = useState<number | null>(
    uid ? Number.parseInt(uid) : null
  );
  const [isAuthenticated, setIsAuthenticated] = useState(!!uid);
  const chrono = new ChronoClient();

  console.log(isAuthenticated);

  const logout = useCallback(async () => {
    await chrono.auth.logout();
    localStorage.removeItem("user");
    setUserId(null);
    setIsAuthenticated(false);
  }, []);

  const login = useCallback(async (data: LoginRequest) => {
    const user = (await chrono.auth.login(data)).data;
    localStorage.setItem("user", user.id);
    setUserId(user.id);
    setIsAuthenticated(true);
  }, []);

  return (
    <AuthContext.Provider value={{ isAuthenticated, userId, login, logout }}>
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
