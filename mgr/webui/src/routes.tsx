import { createBrowserRouter, Navigate } from "react-router-dom";

import { AppShell } from "@/layouts/AppShell";
import { LoginPage } from "@/pages/Login";
import { NotFoundPage } from "@/pages/NotFound";
import { OverviewPage } from "@/pages/Overview";
import { TenantListPage } from "@/pages/Tenants/List";
import { UserListPage } from "@/pages/Users/List";

export const router = createBrowserRouter([
  {
    path: "/",
    element: <LoginPage />,
    errorElement: <NotFoundPage />,
  },
  {
    path: "/login",
    element: <LoginPage />,
  },
  {
    element: <AppShell />,
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
