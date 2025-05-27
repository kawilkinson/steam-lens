"use client";
import { createContext, useContext, useState, useEffect, ReactNode, useCallback } from "react";
import { backendURL } from "../definitions/urls";

interface AuthState {
  loggedIn: boolean;
  steamID: string | null;
}

const AuthContext = createContext<{
  auth: AuthState;
  setAuth: (auth: AuthState) => void;
  checkAuth: () => Promise<void>;
}>({
  auth: { loggedIn: false, steamID: null },
  setAuth: () => {},
  checkAuth: async () => {},
});

export function useAuth() {
  return useContext(AuthContext);
}

export function AuthProvider({ children }: { children: ReactNode }) {
  const [auth, setAuth] = useState<AuthState>({ loggedIn: false, steamID: null });

  const checkAuth = useCallback(async () => {
    if (!backendURL) {
    throw new Error("Backend base URL not set up in environment variables");
    }
    const url = new URL("users/me", backendURL);
    try {
      const resp = await fetch(url.toString(), {
        method: "GET",
        credentials: "include", 
      });
      if (!resp.ok) throw new Error("Not authenticated");
      const data = await resp.json();
      setAuth({ loggedIn: true, steamID: data.user?.steam_id ?? null });
    } catch {
      setAuth({ loggedIn: false, steamID: null });
    }
  }, []);
  
  useEffect(() => {
    checkAuth();
  }, [checkAuth]);

  return (
    <AuthContext.Provider value={{ auth, setAuth, checkAuth }}>
      {children}
    </AuthContext.Provider>
  );
}
