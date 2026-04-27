import type { ReactNode } from "react";

import {
  AppstoreOutlined,
  AuditOutlined,
  FileTextOutlined,
  HomeOutlined,
  KeyOutlined,
  SafetyCertificateOutlined,
  SettingOutlined,
  TeamOutlined,
} from "@ant-design/icons";
import { createBrowserRouter, Navigate } from "react-router-dom";

import { AppShell } from "@/layouts/AppShell";
import { LoginPage } from "@/pages/Login";
import { NotFoundPage } from "@/pages/NotFound";
import { OverviewPage } from "@/pages/Overview";

export interface NavigationItem {
  path: string;
  label: string;
  icon: ReactNode;
}

export interface NavigationGroup {
  title: string;
  items: NavigationItem[];
}

export const homeNavigationItem: NavigationItem = {
  path: "/overview",
  label: "首页",
  icon: <HomeOutlined />,
};

export const navigationGroups: NavigationGroup[] = [
  {
    title: "租户管理",
    items: [{ path: "/tenants/list", label: "租户列表", icon: <TeamOutlined /> }],
  },
  {
    title: "角色与权限",
    items: [
      { path: "/permissions/roles", label: "角色管理", icon: <SafetyCertificateOutlined /> },
      { path: "/permissions/grants", label: "权限分配", icon: <AppstoreOutlined /> },
    ],
  },
  {
    title: "凭证管理",
    items: [
      { path: "/credentials/aksk", label: "AK/SK 管理", icon: <KeyOutlined /> },
      { path: "/credentials/tokens", label: "Token 管理", icon: <SettingOutlined /> },
    ],
  },
  {
    title: "审计日志",
    items: [
      { path: "/audit/logs", label: "操作日志", icon: <FileTextOutlined /> },
      { path: "/audit/search", label: "日志查询", icon: <AuditOutlined /> },
    ],
  },
];

export const navigationItems: NavigationItem[] = [
  homeNavigationItem,
  ...navigationGroups.flatMap((group) => group.items),
];

function DevelopingPage({ title }: { title: string }) {
  return (
    <section className="empty-page">
      <SettingOutlined />
      <h2>{title}</h2>
      <p>开发中</p>
    </section>
  );
}

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
      ...navigationItems
        .filter((item) => item.path !== "/overview")
        .map((item) => ({
          path: item.path.replace(/^\//, ""),
          element: <DevelopingPage title={item.label} />,
        })),
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
