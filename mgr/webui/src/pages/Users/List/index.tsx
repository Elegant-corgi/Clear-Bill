import { useEffect, useMemo, useState } from "react";

import {
  App,
  Button,
  Drawer,
  Form,
  Input,
  Modal,
  Popconfirm,
  Select,
  Space,
  Table,
  Tag,
  Typography,
} from "antd";
import type { TableColumnsType } from "antd";
import { PlusOutlined } from "@ant-design/icons";

import { useAuth } from "@/auth/AuthContext";
import { rolesList } from "@/services/clear-bill/role";
import { tenantsList } from "@/services/clear-bill/tenant";
import {
  usersCreate,
  usersDelete,
  usersList,
  usersPasswordReset,
  usersUpdate,
} from "@/services/clear-bill/user";
import { formatDateTime, getErrorMessage, unwrapResponse } from "@/utils/api";
import { isSysadmin, isTenantAdmin } from "@/utils/access";

import styles from "./index.module.css";

interface SearchFormValues {
  keyword?: string;
}

interface UserFormValues {
  displayName: string;
  role: string;
  status: string;
  tenantId?: number;
  username?: string;
}

interface ResetPasswordValues {
  newPassword: string;
}

const STATUS_OPTIONS = [
  { label: "启用", value: "active" },
  { label: "禁用", value: "disabled" },
];

function statusMeta(status?: string) {
  if (status === "disabled") {
    return { color: "error", label: "禁用" };
  }

  return { color: "success", label: "启用" };
}

export function UserListPage() {
  const { message, modal } = App.useApp();
  const { user: currentUser } = useAuth();
  const [searchForm] = Form.useForm<SearchFormValues>();
  const [userForm] = Form.useForm<UserFormValues>();
  const [resetPasswordForm] = Form.useForm<ResetPasswordValues>();
  const [users, setUsers] = useState<API.User[]>([]);
  const [roles, setRoles] = useState<API.Role[]>([]);
  const [tenants, setTenants] = useState<API.Tenant[]>([]);
  const [keyword, setKeyword] = useState("");
  const [loading, setLoading] = useState(false);
  const [submitting, setSubmitting] = useState(false);
  const [resetPasswordSubmitting, setResetPasswordSubmitting] = useState(false);
  const [drawerOpen, setDrawerOpen] = useState(false);
  const [resetPasswordOpen, setResetPasswordOpen] = useState(false);
  const [editingUser, setEditingUser] = useState<API.User | null>(null);
  const [passwordTarget, setPasswordTarget] = useState<API.User | null>(null);
  const selectedRoleCode = Form.useWatch("role", userForm);

  const sysadmin = isSysadmin(currentUser);
  const canCreateUser = isTenantAdmin(currentUser);
  const isEditMode = Boolean(editingUser);

  const tenantMap = useMemo(
    () =>
      tenants.reduce<Record<number, string>>((acc, item) => {
        acc[item.id] = item.name;
        return acc;
      }, {}),
    [tenants],
  );

  const roleOptions = useMemo(
    () =>
      roles.map((item) => ({
        label: item.name,
        value: item.code,
      })),
    [roles],
  );

  const selectedRole = useMemo(
    () => roles.find((item) => item.code === selectedRoleCode),
    [roles, selectedRoleCode],
  );

  const loadBaseData = async () => {
    try {
      const [roleResponse, tenantResponse] = await Promise.all([
        rolesList({}),
        sysadmin ? tenantsList({}) : Promise.resolve({ success: true, data: [] as API.Tenant[] }),
      ]);
      setRoles(unwrapResponse(roleResponse, "获取角色列表失败") ?? []);
      setTenants(unwrapResponse(tenantResponse, "获取租户列表失败") ?? []);
    } catch (error) {
      message.error(getErrorMessage(error, "初始化用户页数据失败"));
    }
  };

  const loadUsers = async (nextKeyword = keyword) => {
    setLoading(true);
    try {
      const response = await usersList(nextKeyword ? { keyword: nextKeyword } : {});
      setUsers(unwrapResponse(response, "获取用户列表失败") ?? []);
    } catch (error) {
      message.error(getErrorMessage(error, "获取用户列表失败"));
    } finally {
      setLoading(false);
    }
  };

  useEffect(() => {
    void Promise.all([loadBaseData(), loadUsers("")]);
  }, []);

  const handleSearch = async (values: SearchFormValues) => {
    const nextKeyword = values.keyword?.trim() ?? "";
    setKeyword(nextKeyword);
    await loadUsers(nextKeyword);
  };

  const handleReset = async () => {
    searchForm.resetFields();
    setKeyword("");
    await loadUsers("");
  };

  const handleOpenCreate = () => {
    setEditingUser(null);
    userForm.resetFields();
    userForm.setFieldsValue({
      status: "active",
      tenantId: currentUser?.tenantId,
    });
    setDrawerOpen(true);
  };

  const handleOpenEdit = (target: API.User) => {
    setEditingUser(target);
    userForm.resetFields();
    userForm.setFieldsValue({
      displayName: target.displayName,
      role: target.role,
      status: target.status,
      tenantId: target.tenantId,
      username: target.username,
    });
    setDrawerOpen(true);
  };

  const handleCloseDrawer = () => {
    setDrawerOpen(false);
    setEditingUser(null);
    userForm.resetFields();
  };

  const handleSubmit = async () => {
    try {
      const values = await userForm.validateFields();
      setSubmitting(true);

      if (editingUser) {
        const response = await usersUpdate(
          { id: editingUser.id },
          {
            displayName: values.displayName.trim(),
            role: values.role,
            status: values.status,
            tenantId: selectedRole?.scope === "tenant" ? values.tenantId : undefined,
          },
        );
        unwrapResponse(response, "更新用户失败");
        message.success("用户更新成功");
      } else {
        const response = await usersCreate({
          displayName: values.displayName.trim(),
          role: values.role,
          status: values.status,
          tenantId: selectedRole?.scope === "tenant" ? values.tenantId : undefined,
          username: values.username?.trim() ?? "",
        });
        const data = unwrapResponse(response, "创建用户失败");
        message.success("用户创建成功");
        modal.success({
          title: "用户创建成功",
          content: (
            <div className={styles.successContent}>
              <p>
                登录账号：
                <strong>{data.user.username}</strong>
              </p>
              <p>
                初始密码：
                <strong>{data.initialPassword}</strong>
              </p>
            </div>
          ),
          okText: "知道了",
        });
      }

      handleCloseDrawer();
      await loadUsers(keyword);
    } catch (error) {
      if (typeof error === "object" && error !== null && "errorFields" in error) {
        return;
      }

      message.error(getErrorMessage(error, editingUser ? "更新用户失败" : "创建用户失败"));
    } finally {
      setSubmitting(false);
    }
  };

  const handleDelete = async (target: API.User) => {
    try {
      setLoading(true);
      const response = await usersDelete({ id: target.id });
      unwrapResponse(response, "删除用户失败");
      message.success("用户删除成功");
      await loadUsers(keyword);
    } catch (error) {
      message.error(getErrorMessage(error, "删除用户失败"));
      setLoading(false);
    }
  };

  const handleOpenResetPassword = (target: API.User) => {
    setPasswordTarget(target);
    resetPasswordForm.resetFields();
    setResetPasswordOpen(true);
  };

  const handleResetPassword = async () => {
    if (!passwordTarget) {
      return;
    }

    try {
      const values = await resetPasswordForm.validateFields();
      setResetPasswordSubmitting(true);
      const response = await usersPasswordReset(
        { id: passwordTarget.id },
        { newPassword: values.newPassword.trim() },
      );
      unwrapResponse(response, "重置密码失败");
      message.success("密码重置成功");
      setResetPasswordOpen(false);
      setPasswordTarget(null);
      resetPasswordForm.resetFields();
    } catch (error) {
      if (typeof error === "object" && error !== null && "errorFields" in error) {
        return;
      }
      message.error(getErrorMessage(error, "重置密码失败"));
    } finally {
      setResetPasswordSubmitting(false);
    }
  };

  const columns: TableColumnsType<API.User> = [
    {
      title: "登录账号",
      dataIndex: "username",
      key: "username",
      width: 180,
    },
    {
      title: "用户名称",
      dataIndex: "displayName",
      key: "displayName",
      width: 180,
    },
    {
      title: "角色",
      dataIndex: "role",
      key: "role",
      width: 180,
      render: (value: string) => roles.find((item) => item.code === value)?.name || value,
    },
    ...(sysadmin
      ? [
          {
            title: "所属租户",
            dataIndex: "tenantId",
            key: "tenantId",
            width: 180,
            render: (value?: number) => (value ? tenantMap[value] || `租户#${value}` : "-"),
          },
        ]
      : []),
    {
      title: "状态",
      dataIndex: "status",
      key: "status",
      width: 110,
      render: (value: string) => {
        const status = statusMeta(value);
        return <Tag color={status.color}>{status.label}</Tag>;
      },
    },
    {
      title: "创建时间",
      dataIndex: "createdAt",
      key: "createdAt",
      width: 180,
      render: (value: string) => formatDateTime(value),
    },
    {
      title: "更新时间",
      dataIndex: "updatedAt",
      key: "updatedAt",
      width: 180,
      render: (value: string) => formatDateTime(value),
    },
    {
      title: "操作",
      key: "actions",
      fixed: "right",
      width: 200,
      render: (_, record) => (
        <Space size={0} className={styles.actionGroup}>
          {!sysadmin ? (
            <Button type="link" className="ui-action-link ui-action-edit" onClick={() => handleOpenEdit(record)}>
              编辑
            </Button>
          ) : null}
          <Button type="link" className="ui-action-link ui-action-reset" onClick={() => handleOpenResetPassword(record)}>
            重置密码
          </Button>
          <Popconfirm
            title="删除后不可恢复，确认删除该用户吗？"
            okText="删除"
            cancelText="取消"
            okButtonProps={{ danger: true }}
            onConfirm={() => handleDelete(record)}
          >
            <Button type="link" className="ui-action-link ui-action-delete">
              删除
            </Button>
          </Popconfirm>
        </Space>
      ),
    },
  ];

  return (
    <section className={styles.page}>
      <div className={styles.card}>
        <div className={styles.toolbar}>
          <Form<SearchFormValues>
            form={searchForm}
            layout="inline"
            className={styles.searchForm}
            onFinish={(values) => void handleSearch(values)}
          >
            <Form.Item<SearchFormValues> name="keyword" label="账号/名称" className={styles.searchItem}>
              <Input allowClear placeholder="请输入账号或名称" />
            </Form.Item>
            <Form.Item className={styles.actions}>
              <Space wrap>
                <Button type="primary" htmlType="submit" loading={loading}>
                  查询
                </Button>
                <Button onClick={() => void handleReset()} disabled={loading}>
                  重置
                </Button>
              </Space>
            </Form.Item>
          </Form>

          {canCreateUser ? (
            <Button type="primary" ghost icon={<PlusOutlined />} onClick={handleOpenCreate}>
              新建用户
            </Button>
          ) : null}
        </div>

        <Table<API.User>
          rowKey="id"
          loading={loading}
          columns={columns}
          dataSource={users}
          className={styles.table}
          scroll={{ x: sysadmin ? 1400 : 1180 }}
          pagination={{
            pageSize: 10,
            showSizeChanger: false,
            showTotal: (total) => `共 ${total} 条`,
          }}
          locale={{
            emptyText: (
              <div className={styles.emptyState}>
                <Typography.Text type="secondary">暂无用户数据</Typography.Text>
              </div>
            ),
          }}
        />
      </div>

      <Drawer
        width={520}
        title={isEditMode ? "编辑用户" : "新建用户"}
        open={drawerOpen}
        onClose={handleCloseDrawer}
        destroyOnClose
        footer={
          <Space>
            <Button onClick={handleCloseDrawer}>取消</Button>
            <Button type="primary" loading={submitting} onClick={() => void handleSubmit()}>
              保存
            </Button>
          </Space>
        }
        styles={{
          footer: {
            display: "flex",
            justifyContent: "flex-end",
          },
        }}
      >
        <Form<UserFormValues>
          form={userForm}
          layout="vertical"
          initialValues={{
            status: "active",
            tenantId: currentUser?.tenantId,
          }}
          className={styles.drawerForm}
        >
          {!isEditMode ? (
            <Form.Item<UserFormValues>
              label="登录账号"
              name="username"
              rules={[
                { required: true, message: "请输入登录账号" },
                { max: 64, message: "登录账号不能超过 64 个字符" },
              ]}
            >
              <Input placeholder="请输入登录账号" maxLength={64} />
            </Form.Item>
          ) : (
            <Form.Item<UserFormValues> label="登录账号" name="username">
              <Input disabled />
            </Form.Item>
          )}

          <Form.Item<UserFormValues>
            label="用户名称"
            name="displayName"
            rules={[
              { required: true, message: "请输入用户名称" },
              { max: 128, message: "用户名称不能超过 128 个字符" },
            ]}
          >
            <Input placeholder="请输入用户名称" maxLength={128} />
          </Form.Item>

          <Form.Item<UserFormValues>
            label="角色"
            name="role"
            rules={[{ required: true, message: "请选择角色" }]}
          >
            <Select options={roleOptions} placeholder="请选择角色" />
          </Form.Item>

          {selectedRole?.scope === "tenant" && sysadmin ? (
            <Form.Item<UserFormValues>
              label="所属租户"
              name="tenantId"
              rules={[{ required: true, message: "请选择所属租户" }]}
            >
              <Select
                options={tenants.map((item) => ({ label: item.name, value: item.id }))}
                placeholder="请选择所属租户"
              />
            </Form.Item>
          ) : null}

          <Form.Item<UserFormValues> label="状态" name="status">
            <Select options={STATUS_OPTIONS} />
          </Form.Item>
        </Form>
      </Drawer>

      <Modal
        title={passwordTarget ? `重置密码 - ${passwordTarget.displayName}` : "重置密码"}
        open={resetPasswordOpen}
        onCancel={() => {
          setResetPasswordOpen(false);
          setPasswordTarget(null);
          resetPasswordForm.resetFields();
        }}
        onOk={() => void handleResetPassword()}
        okText="确认重置"
        cancelText="取消"
        confirmLoading={resetPasswordSubmitting}
      >
        <Form<ResetPasswordValues> form={resetPasswordForm} layout="vertical">
          <Form.Item<ResetPasswordValues>
            label="新密码"
            name="newPassword"
            rules={[
              { required: true, message: "请输入新密码" },
              { min: 6, message: "新密码至少 6 位" },
            ]}
          >
            <Input.Password placeholder="请输入新密码" />
          </Form.Item>
        </Form>
      </Modal>
    </section>
  );
}
