import { Spin } from "antd";
import { Navigate, useLocation } from "react-router-dom";

import { useAuth } from "@/auth/AuthContext";
import { resolveFirstAccessiblePath } from "@/utils/access";

function FullscreenSpin() {
  return (
    <div
      style={{
        display: "grid",
        minHeight: "100vh",
        placeItems: "center",
      }}
    >
      <Spin size="large" />
    </div>
  );
}

export function ProtectedRoute({ children }: { children: React.ReactNode }) {
  const { isAuthenticated, loading } = useAuth();
  const location = useLocation();

  if (loading) {
    return <FullscreenSpin />;
  }

  if (!isAuthenticated) {
    return <Navigate to="/login" replace state={{ from: location.pathname }} />;
  }

  return <>{children}</>;
}

export function GuestRoute({ children }: { children: React.ReactNode }) {
  const { loading, user } = useAuth();

  if (loading) {
    return <FullscreenSpin />;
  }

  if (user) {
    return <Navigate to={resolveFirstAccessiblePath(user)} replace />;
  }

  return <>{children}</>;
}
