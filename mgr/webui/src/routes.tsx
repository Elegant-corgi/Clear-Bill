import { createBrowserRouter, Navigate } from "react-router-dom";

import { GuestRoute, ProtectedRoute } from "@/auth/AuthGuards";
import { AppShell } from "@/layouts/AppShell";
import { LoginPage } from "@/pages/Login";
import { NotFoundPage } from "@/pages/NotFound";
import { OverviewPage } from "@/pages/Overview";
import { RoleListPage } from "@/pages/Roles/List";
import { TenantListPage } from "@/pages/Tenants/List";
import { UserListPage } from "@/pages/Users/List";

export const router = createBrowserRouter([
  {
    path: "/",
    element: (
      <GuestRoute>
        <LoginPage />
      </GuestRoute>
    ),
    errorElement: <NotFoundPage />,
  },
  {
    path: "/login",
    element: (
      <GuestRoute>
        <LoginPage />
      </GuestRoute>
    ),
  },
  {
    element: (
      <ProtectedRoute>
        <AppShell />
      </ProtectedRoute>
    ),
    errorElement: <NotFoundPage />,
    children: [
      {
        path: "overview",
        element: <OverviewPage />,
      },
      {
        path: "tenants/list",
        element: <TenantListPage />,
      },
      {
        path: "tenants/users",
        element: <UserListPage />,
      },
      {
        path: "tenants/roles",
        element: <RoleListPage />,
      },
      {
        path: "",
        element: <Navigate to="/overview" replace />,
      },
      {
        path: "*",
        element: <NotFoundPage />,
      },
    ],
  },
]);
