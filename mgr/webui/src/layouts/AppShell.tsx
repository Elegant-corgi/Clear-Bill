import { useState } from "react";

import {
  BellOutlined,
  MenuOutlined,
  ThunderboltOutlined,
} from "@ant-design/icons";
import {
  Avatar,
  Badge,
  Breadcrumb,
  Button,
  Drawer,
  Grid,
  Layout,
  Menu,
  Space,
  Tag,
  Typography,
} from "antd";
import type { MenuProps } from "antd";
import { Outlet, useLocation, useNavigate } from "react-router-dom";

import defaultSettings from "@config/defaultSettings";
import { navigationItems } from "@/routes";

const { Content, Header, Sider } = Layout;

function getSelectedPath(pathname: string) {
  const exactMatch = navigationItems.find((item) => item.path === pathname);
  if (exactMatch) {
    return exactMatch.path;
  }

  const branchMatch = navigationItems.find(
    (item) => item.path !== "/" && pathname.startsWith(item.path),
  );

  return branchMatch?.path ?? "/overview";
}

function getBreadcrumbItems(pathname: string) {
  const current = navigationItems.find((item) => item.path === getSelectedPath(pathname));
  return [
    { title: defaultSettings.title },
    { title: current?.label ?? "Workspace" },
  ];
}

export function AppShell() {
  const screens = Grid.useBreakpoint();
  const location = useLocation();
  const navigate = useNavigate();
  const [drawerOpen, setDrawerOpen] = useState(false);

  const selectedPath = getSelectedPath(location.pathname);

  const menuItems: MenuProps["items"] = navigationItems.map((item) => ({
    key: item.path,
    icon: item.icon,
    label: item.label,
  }));

  const handleNavigate = (key: string) => {
    navigate(key);
    setDrawerOpen(false);
  };

  const sideNavigation = (
    <>
      <div className="shell__brand">
        <div className="shell__logo">
          <span className="shell__logo-mark">CB</span>
          <div className="shell__logo-text">
            <Typography.Title level={4} className="shell__logo-title">
              {defaultSettings.title}
            </Typography.Title>
            <Typography.Text className="shell__logo-subtitle">
              {defaultSettings.subtitle}
            </Typography.Text>
          </div>
        </div>
        <Typography.Paragraph className="shell__logo-subtitle">
          {defaultSettings.description}
        </Typography.Paragraph>
      </div>

      <div className="shell__nav">
        <Menu
          mode="inline"
          selectedKeys={[selectedPath]}
          items={menuItems}
          onClick={({ key }) => handleNavigate(String(key))}
        />
      </div>

      <div className="shell__sider-card">
        <Space direction="vertical" size={8}>
          <Tag color="cyan">Close week ready</Tag>
          <Typography.Title level={5} style={{ margin: 0 }}>
            Shrink review cycles, not visibility.
          </Typography.Title>
          <Typography.Text type="secondary">
            Track open bills, auto-matching, and settlement friction from one place.
          </Typography.Text>
        </Space>
      </div>
    </>
  );

  return (
    <Layout className="shell">
      {screens.lg ? (
        <Sider width={292} className="shell__sider">
          {sideNavigation}
        </Sider>
      ) : null}

      <Layout className="shell__workspace">
        <Header className="shell__header">
          <Space align="start" size={16}>
            {!screens.lg ? (
              <Button
                shape="circle"
                type="default"
                icon={<MenuOutlined />}
                onClick={() => setDrawerOpen(true)}
              />
            ) : null}

            <div className="shell__hero">
              <Breadcrumb items={getBreadcrumbItems(location.pathname)} />
              <Typography.Title level={3} className="shell__title">
                {defaultSettings.heroTitle}
              </Typography.Title>
              <Typography.Text type="secondary">
                {defaultSettings.heroDescription}
              </Typography.Text>
            </div>
          </Space>

          <Space wrap size={12}>
            <Tag icon={<ThunderboltOutlined />} color="green">
              Auto-match 92.4%
            </Tag>
            <Badge count={5} size="small">
              <Button shape="circle" icon={<BellOutlined />} />
            </Badge>
            <Avatar style={{ backgroundColor: "#0f766e" }}>CB</Avatar>
          </Space>
        </Header>

        <Content className="shell__content">
          <Outlet />
        </Content>
      </Layout>

      <Drawer
        placement="left"
        width={292}
        open={drawerOpen}
        onClose={() => setDrawerOpen(false)}
        styles={{ body: { padding: 0, background: "rgba(9, 31, 43, 0.98)" } }}
      >
        <div className="shell__sider">{sideNavigation}</div>
      </Drawer>
    </Layout>
  );
}
