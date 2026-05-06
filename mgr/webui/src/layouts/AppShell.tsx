import { useEffect, useMemo, useState } from "react";

import {
  DoubleLeftOutlined,
  DoubleRightOutlined,
  DownOutlined,
  FileDoneOutlined,
  HomeOutlined,
  KeyOutlined,
  MenuOutlined,
  SafetyCertificateOutlined,
  TeamOutlined,
  UserOutlined,
} from "@ant-design/icons";
import { App, Button, Drawer, Dropdown, Form, Grid, Input, Layout, Modal, Spin } from "antd";
import type { MenuProps } from "antd";
import { Navigate, Outlet, useLocation, useNavigate } from "react-router-dom";

import { useAuth } from "@/auth/AuthContext";
import defaultSettings from "@config/defaultSettings";
import { getErrorMessage } from "@/utils/api";
import {
  canAccessCredentialList,
  canAccessPath,
  canAccessRoleList,
  canAccessTenantList,
  canAccessUserList,
  getCurrentTitle,
  CREDENTIAL_LIST_PATH,
  OVERVIEW_PATH,
  ROLE_LIST_PATH,
  roleLabel,
  TENANT_LIST_PATH,
  USER_LIST_PATH,
} from "@/utils/access";

const { Content, Header, Sider } = Layout;

interface PasswordFormValues {
  confirmPassword: string;
  newPassword: string;
  oldPassword: string;
}

interface MenuItem {
  icon: React.ReactNode;
  key: string;
  label: string;
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
  menuItems: MenuItem[];
  selectedPath: string;
  tenantMenuCollapsed: boolean;
  onNavigate: (path: string) => void;
  onToggleCollapse: () => void;
  onToggleTenantMenu: () => void;
}

function Sidebar({
  collapsed,
  menuItems,
  selectedPath,
  tenantMenuCollapsed,
  onNavigate,
  onToggleCollapse,
  onToggleTenantMenu,
}: SidebarProps) {
  const tenantMenuItems = menuItems.filter((item) => item.key !== OVERVIEW_PATH);

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

        {tenantMenuItems.length > 0 ? (
          <section className="shell__menu-group">
            {!collapsed ? (
              <button className="shell__menu-group-title" type="button" onClick={onToggleTenantMenu}>
                <span>租户管理</span>
                <DownOutlined className={tenantMenuCollapsed ? "is-folded" : ""} />
              </button>
            ) : null}

            {!tenantMenuCollapsed || collapsed ? (
              <div className="shell__menu-list">
                {tenantMenuItems.map((item) => (
                  <button
                    key={item.key}
                    className={`shell__menu-item ${selectedPath === item.key ? "is-active" : ""}`}
                    type="button"
                    title={collapsed ? item.label : undefined}
                    onClick={() => onNavigate(item.key)}
                  >
                    <span className="shell__menu-icon">{item.icon}</span>
                    {!collapsed ? <span className="shell__menu-label">{item.label}</span> : null}
                  </button>
                ))}
              </div>
            ) : null}
          </section>
        ) : null}
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
  const { changePassword, loading, logout, user } = useAuth();
  const screens = Grid.useBreakpoint();
  const location = useLocation();
  const navigate = useNavigate();
  const [passwordForm] = Form.useForm<PasswordFormValues>();
  const [drawerOpen, setDrawerOpen] = useState(false);
  const [passwordOpen, setPasswordOpen] = useState(false);
  const [passwordSubmitting, setPasswordSubmitting] = useState(false);
  const [siderCollapsed, setSiderCollapsed] = useState(false);
  const [tenantMenuCollapsed, setTenantMenuCollapsed] = useState(false);

  const menuItems = useMemo<MenuItem[]>(() => {
    const items: MenuItem[] = [
      {
        icon: <HomeOutlined />,
        key: OVERVIEW_PATH,
        label: "首页",
      },
    ];

    if (canAccessTenantList(user)) {
      items.push({
        icon: <TeamOutlined />,
        key: TENANT_LIST_PATH,
        label: "租户列表",
      });
    }

    if (canAccessUserList(user)) {
      items.push({
        icon: <UserOutlined />,
        key: USER_LIST_PATH,
        label: "用户列表",
      });
    }

    if (canAccessRoleList(user)) {
      items.push({
        icon: <SafetyCertificateOutlined />,
        key: ROLE_LIST_PATH,
        label: "角色权限",
      });
    }

    if (canAccessCredentialList(user)) {
      items.push({
        icon: <KeyOutlined />,
        key: CREDENTIAL_LIST_PATH,
        label: "凭证管理",
      });
    }

    return items;
  }, [user]);

  const selectedPath =
    menuItems.find((item) => location.pathname.startsWith(item.key) && item.key !== OVERVIEW_PATH)?.key ||
    (location.pathname.startsWith(OVERVIEW_PATH) ? OVERVIEW_PATH : menuItems[0]?.key || OVERVIEW_PATH);

  useEffect(() => {
    if (user && !canAccessPath(user, location.pathname)) {
      navigate(OVERVIEW_PATH, { replace: true });
    }
  }, [location.pathname, navigate, user]);

  if (loading) {
    return (
      <div style={{ display: "grid", minHeight: "100vh", placeItems: "center" }}>
        <Spin size="large" />
      </div>
    );
  }

  if (!user) {
    return <Navigate to="/login" replace />;
  }

  const handleNavigate = (path: string) => {
    navigate(path);
    setDrawerOpen(false);
  };

  const handleLogout = async () => {
    try {
      await logout();
      message.success("已退出登录");
      navigate("/login", { replace: true });
    } catch (error) {
      message.error(getErrorMessage(error, "退出登录失败"));
    }
  };

  const handlePasswordSubmit = async () => {
    try {
      const values = await passwordForm.validateFields();
      setPasswordSubmitting(true);
      await changePassword({
        newPassword: values.newPassword.trim(),
        oldPassword: values.oldPassword.trim(),
      });
      message.success("密码修改成功，请使用新密码重新登录");
      setPasswordOpen(false);
      passwordForm.resetFields();
      await handleLogout();
    } catch (error) {
      if (typeof error === "object" && error !== null && "errorFields" in error) {
        return;
      }

      message.error(getErrorMessage(error, "修改密码失败"));
    } finally {
      setPasswordSubmitting(false);
    }
  };

  const accountMenuItems: MenuProps["items"] = [
    { key: "password", label: "修改密码" },
    { type: "divider" },
    { key: "logout", label: "退出登录" },
  ];

  const handleAccountMenuClick: MenuProps["onClick"] = ({ key }) => {
    if (key === "password") {
      passwordForm.resetFields();
      setPasswordOpen(true);
      return;
    }

    if (key === "logout") {
      void handleLogout();
    }
  };

  const sidebar = (
    <Sidebar
      collapsed={siderCollapsed}
      menuItems={menuItems}
      tenantMenuCollapsed={tenantMenuCollapsed}
      selectedPath={selectedPath}
      onNavigate={handleNavigate}
      onToggleCollapse={() => setSiderCollapsed((value) => !value)}
      onToggleTenantMenu={() => setTenantMenuCollapsed((value) => !value)}
    />
  );

  const drawerSidebar = (
    <Sidebar
      collapsed={false}
      menuItems={menuItems}
      tenantMenuCollapsed={tenantMenuCollapsed}
      selectedPath={selectedPath}
      onNavigate={handleNavigate}
      onToggleCollapse={() => setDrawerOpen(false)}
      onToggleTenantMenu={() => setTenantMenuCollapsed((value) => !value)}
    />
  );

  return (
    <>
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
              <h1>{getCurrentTitle(location.pathname)}</h1>
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
                    <strong>{user.displayName || user.username}</strong>
                    <span>{roleLabel(user.role)}</span>
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

      <Modal
        title="修改密码"
        open={passwordOpen}
        onCancel={() => {
          setPasswordOpen(false);
          passwordForm.resetFields();
        }}
        onOk={() => void handlePasswordSubmit()}
        okText="确认修改"
        cancelText="取消"
        confirmLoading={passwordSubmitting}
      >
        <Form<PasswordFormValues> form={passwordForm} layout="vertical">
          <Form.Item<PasswordFormValues>
            label="旧密码"
            name="oldPassword"
            rules={[{ required: true, message: "请输入旧密码" }]}
          >
            <Input.Password placeholder="请输入旧密码" />
          </Form.Item>
          <Form.Item<PasswordFormValues>
            label="新密码"
            name="newPassword"
            rules={[
              { required: true, message: "请输入新密码" },
              { min: 6, message: "新密码至少 6 位" },
            ]}
          >
            <Input.Password placeholder="请输入新密码" />
          </Form.Item>
          <Form.Item<PasswordFormValues>
            label="确认新密码"
            name="confirmPassword"
            dependencies={["newPassword"]}
            rules={[
              { required: true, message: "请再次输入新密码" },
              ({ getFieldValue }) => ({
                validator(_, value) {
                  if (!value || getFieldValue("newPassword") === value) {
                    return Promise.resolve();
                  }
                  return Promise.reject(new Error("两次输入的新密码不一致"));
                },
              }),
            ]}
          >
            <Input.Password placeholder="请再次输入新密码" />
          </Form.Item>
        </Form>
      </Modal>
    </>
  );
}
