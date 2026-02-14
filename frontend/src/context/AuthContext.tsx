import {
  createContext,
  useContext,
  useState,
  useCallback,
  useEffect,
  type ReactNode,
} from "react";
import type { User } from "../types";
import * as authApi from "../api/auth";

interface AuthState {
  user: User | null;
  loading: boolean;
  error: string | null;
  ready: boolean;
}

interface AuthContextValue extends AuthState {
  signUp: (email: string, password: string, confirmPassword: string) => Promise<void>;
  signIn: (email: string, password: string) => Promise<void>;
  logout: () => void;
}

const AuthContext = createContext<AuthContextValue | null>(null);

export function AuthProvider({ children }: { children: ReactNode }) {
  const [state, setState] = useState<AuthState>({
    user: null,
    loading: false,
    error: null,
    ready: false,
  });

  useEffect(() => {
    authApi.me().then((user) => {
      setState({ user, loading: false, error: null, ready: true });
    });
  }, []);

  const signUp = useCallback(
    async (email: string, password: string, confirmPassword: string) => {
      setState((s) => ({ ...s, loading: true, error: null }));
      try {
        const user = await authApi.signUp(email, password, confirmPassword);
        setState({ user, loading: false, error: null, ready: true });
      } catch (e: any) {
        setState((s) => ({
          ...s,
          loading: false,
          error: e.message || "Sign up failed",
        }));
        throw e;
      }
    },
    [],
  );

  const signIn = useCallback(async (email: string, password: string) => {
    setState((s) => ({ ...s, loading: true, error: null }));
    try {
      const user = await authApi.signIn(email, password);
      setState({ user, loading: false, error: null, ready: true });
    } catch (e: any) {
      setState((s) => ({
        ...s,
        loading: false,
        error: e.message || "Sign in failed",
      }));
      throw e;
    }
  }, []);

  const logout = useCallback(() => {
    setState({ user: null, loading: false, error: null, ready: true });
  }, []);

  if (!state.ready) {
    return (
      <div className="min-h-screen flex items-center justify-center text-zinc-500">
        Loading...
      </div>
    );
  }

  return (
    <AuthContext.Provider value={{ ...state, signUp, signIn, logout }}>
      {children}
    </AuthContext.Provider>
  );
}

export function useAuth() {
  const ctx = useContext(AuthContext);
  if (!ctx) throw new Error("useAuth must be used within AuthProvider");
  return ctx;
}
