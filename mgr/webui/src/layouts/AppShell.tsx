import { useMemo, useState } from "react";

import {
  DoubleLeftOutlined,
  DoubleRightOutlined,
  DownOutlined,
  FileDoneOutlined,
  HomeOutlined,
  MenuOutlined,
  TeamOutlined,
  UserOutlined,
} from "@ant-design/icons";
import { App, Button, Drawer, Dropdown, Grid, Layout } from "antd";
import type { MenuProps } from "antd";
import { Outlet, useLocation, useNavigate } from "react-router-dom";

import defaultSettings from "@config/defaultSettings";

const { Content, Header, Sider } = Layout;

const OVERVIEW_PATH = "/overview";
const TENANT_LIST_PATH = "/tenants/list";
const USER_LIST_PATH = "/tenants/users";

function getSelectedPath(pathname: string) {
  if (pathname.startsWith(TENANT_LIST_PATH)) {
    return TENANT_LIST_PATH;
  }

  if (pathname.startsWith(USER_LIST_PATH)) {
    return USER_LIST_PATH;
  }

  return OVERVIEW_PATH;
}

function getCurrentTitle(pathname: string) {
  const selectedPath = getSelectedPath(pathname);

  if (selectedPath === TENANT_LIST_PATH) {
    return "租户列表";
  }

  if (selectedPath === USER_LIST_PATH) {
    return "用户列表";
  }

  return "首页";
}

function UserAvatar() {
  return (
    <span className="shell__profile-avatar">
      <svg viewBox="0 0 36 36" aria-hidden="true">
        <circle cx="18" cy="13.5" r="5.2" />
        <path d="M8.8 28.2c1.4-5 4.6-7.5 9.2-7.5s7.8 2.5 9.2 7.5" />
      </svg>
    </span>
  );
}

interface SidebarProps {
  collapsed: boolean;
  tenantMenuCollapsed: boolean;
  selectedPath: string;
  onNavigate: (path: string) => void;
  onToggleCollapse: () => void;
  onToggleTenantMenu: () => void;
}

function Sidebar({
  collapsed,
  tenantMenuCollapsed,
  selectedPath,
  onNavigate,
  onToggleCollapse,
  onToggleTenantMenu,
}: SidebarProps) {
  return (
    <div className={`shell__sidebar-inner ${collapsed ? "is-collapsed" : ""}`}>
      <div className="shell__brand">
        <span className="shell__logo-mark">
          <FileDoneOutlined />
        </span>
        {!collapsed ? <strong>{defaultSettings.title}</strong> : null}
      </div>

      <nav className="shell__menu" aria-label="后台导航">
        <button
          className={`shell__menu-item ${selectedPath === OVERVIEW_PATH ? "is-active" : ""}`}
          type="button"
          title={collapsed ? "首页" : undefined}
          onClick={() => onNavigate(OVERVIEW_PATH)}
        >
          <span className="shell__menu-icon">
            <HomeOutlined />
          </span>
          {!collapsed ? <span className="shell__menu-label">首页</span> : null}
        </button>

        <section className="shell__menu-group">
          {!collapsed ? (
            <button className="shell__menu-group-title" type="button" onClick={onToggleTenantMenu}>
              <span>租户管理</span>
              <DownOutlined className={tenantMenuCollapsed ? "is-folded" : ""} />
            </button>
          ) : null}

          {!tenantMenuCollapsed || collapsed ? (
            <div className="shell__menu-list">
              <button
                className={`shell__menu-item ${selectedPath === TENANT_LIST_PATH ? "is-active" : ""}`}
                type="button"
                title={collapsed ? "租户列表" : undefined}
                onClick={() => onNavigate(TENANT_LIST_PATH)}
              >
                <span className="shell__menu-icon">
                  <TeamOutlined />
                </span>
                {!collapsed ? <span className="shell__menu-label">租户列表</span> : null}
              </button>

              <button
                className={`shell__menu-item ${selectedPath === USER_LIST_PATH ? "is-active" : ""}`}
                type="button"
                title={collapsed ? "用户列表" : undefined}
                onClick={() => onNavigate(USER_LIST_PATH)}
              >
                <span className="shell__menu-icon">
                  <UserOutlined />
                </span>
                {!collapsed ? <span className="shell__menu-label">用户列表</span> : null}
              </button>
            </div>
          ) : null}
        </section>
      </nav>

      <button className="shell__collapse" type="button" onClick={onToggleCollapse}>
        {collapsed ? <DoubleRightOutlined /> : <DoubleLeftOutlined />}
        {!collapsed ? <span>收起菜单</span> : null}
      </button>
    </div>
  );
}

export function AppShell() {
  const { message } = App.useApp();
  const screens = Grid.useBreakpoint();
  const location = useLocation();
  const navigate = useNavigate();
  const [drawerOpen, setDrawerOpen] = useState(false);
  const [siderCollapsed, setSiderCollapsed] = useState(false);
  const [tenantMenuCollapsed, setTenantMenuCollapsed] = useState(false);

  const selectedPath = getSelectedPath(location.pathname);
  const currentTitle = getCurrentTitle(location.pathname);

  const handleNavigate = (path: string) => {
    navigate(path);
    setDrawerOpen(false);
  };

  const accountMenuItems: MenuProps["items"] = [
    { key: "password", label: "修改密码" },
    { type: "divider" },
    { key: "logout", label: "退出登录" },
  ];

  const handleAccountMenuClick: MenuProps["onClick"] = ({ key }) => {
    if (key === "password") {
      void message.info("修改密码功能开发中");
      return;
    }

    if (key === "logout") {
      navigate("/login", { replace: true });
    }
  };

  const sidebar = useMemo(
    () => (
      <Sidebar
        collapsed={siderCollapsed}
        tenantMenuCollapsed={tenantMenuCollapsed}
        selectedPath={selectedPath}
        onNavigate={handleNavigate}
        onToggleCollapse={() => setSiderCollapsed((value) => !value)}
        onToggleTenantMenu={() => setTenantMenuCollapsed((value) => !value)}
      />
    ),
    [selectedPath, siderCollapsed, tenantMenuCollapsed],
  );

  const drawerSidebar = (
    <Sidebar
      collapsed={false}
      tenantMenuCollapsed={tenantMenuCollapsed}
      selectedPath={selectedPath}
      onNavigate={handleNavigate}
      onToggleCollapse={() => setDrawerOpen(false)}
      onToggleTenantMenu={() => setTenantMenuCollapsed((value) => !value)}
    />
  );

  return (
    <Layout className="shell">
      {screens.lg ? (
        <Sider width={siderCollapsed ? 92 : 260} className="shell__sider">
          {sidebar}
        </Sider>
      ) : null}

      <Layout className="shell__workspace">
        <Header className="shell__header">
          <div className="shell__header-title">
            {!screens.lg ? (
              <Button shape="circle" icon={<MenuOutlined />} onClick={() => setDrawerOpen(true)} />
            ) : null}
            <h1>{currentTitle}</h1>
          </div>

          <div className="shell__toolbar">
            <span className="shell__toolbar-divider" />
            <Dropdown
              menu={{ items: accountMenuItems, onClick: handleAccountMenuClick }}
              placement="bottomRight"
              trigger={["click"]}
            >
              <button className="shell__profile" type="button">
                <UserAvatar />
                <span className="shell__profile-text">
                  <strong>管理员</strong>
                  <span>超级管理员</span>
                </span>
                <DownOutlined />
              </button>
            </Dropdown>
          </div>
        </Header>

        <Content className="shell__content">
          <Outlet />
        </Content>
      </Layout>

      <Drawer
        placement="left"
        width={260}
        open={drawerOpen}
        onClose={() => setDrawerOpen(false)}
        styles={{ body: { padding: 0 } }}
      >
        {drawerSidebar}
      </Drawer>
    </Layout>
  );
}
