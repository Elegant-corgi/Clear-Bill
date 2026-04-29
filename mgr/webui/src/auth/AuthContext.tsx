import { createContext, useContext, useEffect, useMemo, useState } from "react";

import { authLogin, authLogout, authMe, authPasswordChange } from "@/services/clear-bill/auth";
import { getErrorMessage, unwrapResponse } from "@/utils/api";

interface LoginValues {
  password: string;
  username: string;
}

interface ChangePasswordValues {
  newPassword: string;
  oldPassword: string;
}

interface AuthContextValue {
  changePassword: (values: ChangePasswordValues) => Promise<void>;
  isAuthenticated: boolean;
  loading: boolean;
  login: (values: LoginValues) => Promise<API.User>;
  logout: () => Promise<void>;
  refreshCurrentUser: () => Promise<API.User | null>;
  user: API.User | null;
}

const AuthContext = createContext<AuthContextValue | undefined>(undefined);

export function AuthProvider({ children }: { children: React.ReactNode }) {
  const [user, setUser] = useState<API.User | null>(null);
  const [loading, setLoading] = useState(true);

  const refreshCurrentUser = async () => {
    try {
      const response = await authMe();
      const currentUser = unwrapResponse(response, "获取当前用户信息失败");
      setUser(currentUser);
      return currentUser;
    } catch {
      setUser(null);
      return null;
    } finally {
      setLoading(false);
    }
  };

  useEffect(() => {
    void refreshCurrentUser();
  }, []);

  const value = useMemo<AuthContextValue>(
    () => ({
      changePassword: async (values) => {
        const response = await authPasswordChange(values);
        unwrapResponse(response, "修改密码失败");
      },
      isAuthenticated: Boolean(user),
      loading,
      login: async (values) => {
        const response = await authLogin(values);
        const data = unwrapResponse(response, "登录失败");
        setUser(data.user);
        return data.user;
      },
      logout: async () => {
        try {
          const response = await authLogout();
          unwrapResponse(response, "退出登录失败");
        } catch (error) {
          throw new Error(getErrorMessage(error, "退出登录失败"));
        } finally {
          setUser(null);
        }
      },
      refreshCurrentUser,
      user,
    }),
    [loading, user],
  );

  return <AuthContext.Provider value={value}>{children}</AuthContext.Provider>;
}

export function useAuth() {
  const context = useContext(AuthContext);
  if (!context) {
    throw new Error("useAuth must be used within AuthProvider");
  }

  return context;
}
