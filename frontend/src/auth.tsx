import {
  createContext,
  useCallback,
  useContext,
  useEffect,
  useState,
  type ReactNode,
} from "react";
import { ChronoClient } from "./api/chrono/client";
import type { LoginRequest, User } from "./types/auth";

export interface AuthContext {
  isAuthenticated: boolean;
  login: (data: LoginRequest) => Promise<void>;
  logout: () => Promise<void>;
  userId: number | null;
  setUser: (data: User) => void;
  user: User | null;
}

const AuthContext = createContext<AuthContext | null>(null);

export function AuthProvider({ children }: { children: ReactNode }) {
  const uid = localStorage.getItem("user");
  const [userId, setUserId] = useState<number | null>(
    uid ? Number.parseInt(uid) : null
  );

  const [user, setUser] = useState<User | null>(null);
  const [isAuthenticated, setIsAuthenticated] = useState(!!uid);
  const chrono = new ChronoClient();

  useEffect(() => {
    if (!isAuthenticated || !userId) return;
    if (isAuthenticated && userId && !user)
      chrono.users.getUserById(userId).then((v) => setUser(v));
  }, [isAuthenticated, userId]);

  const logout = useCallback(async () => {
    await chrono.auth.logout();
    localStorage.removeItem("user");
    setUserId(null);
    setUser(null);
    setIsAuthenticated(false);
  }, []);

  const login = useCallback(async (data: LoginRequest) => {
    const user = (await chrono.auth.login(data)).data;
    localStorage.setItem("user", user.id);
    setUserId(user.id);
    setUser(user);
    setIsAuthenticated(true);
  }, []);

  return (
    <AuthContext.Provider
      value={{ isAuthenticated, userId, login, logout, user, setUser }}
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
