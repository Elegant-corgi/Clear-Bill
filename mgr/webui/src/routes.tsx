import type { ReactNode } from "react";

import {
  AuditOutlined,
  DollarCircleOutlined,
  RadarChartOutlined,
  SettingOutlined,
  TeamOutlined,
} from "@ant-design/icons";
import { createBrowserRouter } from "react-router-dom";

import { AppShell } from "@/layouts/AppShell";
import { BillsPage } from "@/pages/Bills";
import { CustomersPage } from "@/pages/Customers";
import { LoginPage } from "@/pages/Login";
import { NotFoundPage } from "@/pages/NotFound";
import { OverviewPage } from "@/pages/Overview";
import { ReconciliationPage } from "@/pages/Reconciliation";
import { SettingsPage } from "@/pages/Settings";

export interface NavigationItem {
  path: string;
  label: string;
  subtitle: string;
  icon: ReactNode;
}

export const navigationItems: NavigationItem[] = [
  {
    path: "/overview",
    label: "Overview",
    subtitle: "Billing command surface",
    icon: <RadarChartOutlined />,
  },
  {
    path: "/bills",
    label: "Bills",
    subtitle: "Track outgoing statements",
    icon: <DollarCircleOutlined />,
  },
  {
    path: "/customers",
    label: "Customers",
    subtitle: "Keep account health visible",
    icon: <TeamOutlined />,
  },
  {
    path: "/reconciliation",
    label: "Reconciliation",
    subtitle: "Watch settlement drift",
    icon: <AuditOutlined />,
  },
  {
    path: "/settings",
    label: "Settings",
    subtitle: "Workspace defaults",
    icon: <SettingOutlined />,
  },
];

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
        path: "bills",
        element: <BillsPage />,
      },
      {
        path: "customers",
        element: <CustomersPage />,
      },
      {
        path: "reconciliation",
        element: <ReconciliationPage />,
      },
      {
        path: "settings",
        element: <SettingsPage />,
      },
      {
        path: "*",
        element: <NotFoundPage />,
      },
    ],
  },
]);
