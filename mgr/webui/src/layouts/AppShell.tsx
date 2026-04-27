import { useMemo, useState } from "react";

import {
  DoubleLeftOutlined,
  DoubleRightOutlined,
  DownOutlined,
  FileDoneOutlined,
  MenuOutlined,
} from "@ant-design/icons";
import { App, Button, Drawer, Dropdown, Grid, Layout } from "antd";
import type { MenuProps } from "antd";
import { Outlet, useLocation, useNavigate } from "react-router-dom";

import defaultSettings from "@config/defaultSettings";
import {
  homeNavigationItem,
  navigationGroups,
  navigationItems,
} from "@/routes";

const { Content, Header, Sider } = Layout;

function getSelectedPath(pathname: string) {
  const current = navigationItems.find((item) => pathname.startsWith(item.path));
  return current?.path ?? "/overview";
}

function getCurrentTitle(pathname: string) {
  return navigationItems.find((item) => item.path === getSelectedPath(pathname))?.label ?? "首页";
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
  collapsedGroups: string[];
  selectedPath: string;
  onNavigate: (path: string) => void;
  onToggleCollapse: () => void;
  onToggleGroup: (title: string) => void;
}

function Sidebar({
  collapsed,
  collapsedGroups,
  selectedPath,
  onNavigate,
  onToggleCollapse,
  onToggleGroup,
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
          className={`shell__menu-item ${selectedPath === homeNavigationItem.path ? "is-active" : ""}`}
          type="button"
          title={collapsed ? homeNavigationItem.label : undefined}
          onClick={() => onNavigate(homeNavigationItem.path)}
        >
          <span className="shell__menu-icon">{homeNavigationItem.icon}</span>
          {!collapsed ? <span className="shell__menu-label">{homeNavigationItem.label}</span> : null}
        </button>

        {navigationGroups.map((group) => {
          const groupCollapsed = collapsedGroups.includes(group.title);

          return (
            <section className="shell__menu-group" key={group.title}>
              {!collapsed ? (
                <button
                  className="shell__menu-group-title"
                  type="button"
                  onClick={() => onToggleGroup(group.title)}
                >
                  <span>{group.title}</span>
                  <DownOutlined className={groupCollapsed ? "is-folded" : ""} />
                </button>
              ) : null}

              {!groupCollapsed || collapsed ? (
                <div className="shell__menu-list">
                  {group.items.map((item) => (
                    <button
                      className={`shell__menu-item ${selectedPath === item.path ? "is-active" : ""}`}
                      key={item.path}
                      type="button"
                      title={collapsed ? item.label : undefined}
                      onClick={() => onNavigate(item.path)}
                    >
                      <span className="shell__menu-icon">{item.icon}</span>
                      {!collapsed ? <span className="shell__menu-label">{item.label}</span> : null}
                    </button>
                  ))}
                </div>
              ) : null}
            </section>
          );
        })}
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
  const [collapsedGroups, setCollapsedGroups] = useState<string[]>([]);

  const selectedPath = getSelectedPath(location.pathname);
  const currentTitle = getCurrentTitle(location.pathname);

  const handleNavigate = (path: string) => {
    navigate(path);
    setDrawerOpen(false);
  };

  const handleToggleGroup = (title: string) => {
    setCollapsedGroups((prev) =>
      prev.includes(title) ? prev.filter((item) => item !== title) : [...prev, title],
    );
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
        collapsedGroups={collapsedGroups}
        selectedPath={selectedPath}
        onNavigate={handleNavigate}
        onToggleCollapse={() => setSiderCollapsed((value) => !value)}
        onToggleGroup={handleToggleGroup}
      />
    ),
    [collapsedGroups, selectedPath, siderCollapsed],
  );

  const drawerSidebar = (
    <Sidebar
      collapsed={false}
      collapsedGroups={collapsedGroups}
      selectedPath={selectedPath}
      onNavigate={handleNavigate}
      onToggleCollapse={() => setDrawerOpen(false)}
      onToggleGroup={handleToggleGroup}
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
