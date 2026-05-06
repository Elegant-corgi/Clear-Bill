import { useEffect, useMemo, useState } from "react";

import {
  App,
  Button,
  Checkbox,
  Drawer,
  Form,
  Input,
  Popconfirm,
  Radio,
  Select,
  Space,
  Table,
  Tag,
  Typography,
} from "antd";
import type { TableColumnsType } from "antd";
import { PlusOutlined } from "@ant-design/icons";

import { useAuth } from "@/auth/AuthContext";
import {
  permissionsList,
  rolesCreate,
  rolesDelete,
  rolesList,
  rolesPermissionsUpdate,
  rolesUpdate,
} from "@/services/clear-bill/role";
import { tenantsList } from "@/services/clear-bill/tenant";
import { formatDateTime, getErrorMessage, unwrapResponse } from "@/utils/api";
import { isSysadmin } from "@/utils/access";
import { PAGE_SIZE_OPTIONS, getPageAfterDelete, useTablePagination } from "@/utils/pagination";

import styles from "./index.module.css";

interface SearchFormValues {
  keyword?: string;
}

interface RoleFormValues {
  code?: string;
  name: string;
  permissionIds: string[];
  scope: "system" | "tenant";
  tenantId?: number;
}

function createEmptyPageResult<T>(): API.PageResult<T> {
  return {
    list: [],
    page: 1,
    pageSize: 1000,
    total: 0,
  };
}

export function RoleListPage() {
  const { message } = App.useApp();
  const { user: currentUser } = useAuth();
  const [searchForm] = Form.useForm<SearchFormValues>();
  const [roleForm] = Form.useForm<RoleFormValues>();
  const [roles, setRoles] = useState<API.Role[]>([]);
  const [permissions, setPermissions] = useState<API.Permission[]>([]);
  const [tenants, setTenants] = useState<API.Tenant[]>([]);
  const [keyword, setKeyword] = useState("");
  const [loading, setLoading] = useState(false);
  const [submitting, setSubmitting] = useState(false);
  const [drawerOpen, setDrawerOpen] = useState(false);
  const [editingRole, setEditingRole] = useState<API.Role | null>(null);
  const { pagination, resetPage, updatePageData, handleTableChange } = useTablePagination();
  const scopeValue = Form.useWatch("scope", roleForm);

  const sysadmin = isSysadmin(currentUser);
  const isEditMode = Boolean(editingRole);

  const groupedPermissions = useMemo(() => {
    return permissions.reduce<Record<string, API.Permission[]>>((acc, permission) => {
      const groupName = permission.tag || "未分组";
      if (!acc[groupName]) {
        acc[groupName] = [];
      }
      acc[groupName].push(permission);
      return acc;
    }, {});
  }, [permissions]);

  const tenantMap = useMemo(
    () =>
      tenants.reduce<Record<number, string>>((acc, item) => {
        acc[item.id] = item.name;
        return acc;
      }, {}),
    [tenants],
  );

  const loadBaseData = async () => {
    try {
      const [permissionResponse, tenantResponse] = await Promise.all([
        permissionsList(),
        sysadmin ? tenantsList({ page: 1, pageSize: 1000 }) : Promise.resolve({ success: true, data: createEmptyPageResult<API.Tenant>() }),
      ]);
      setPermissions(unwrapResponse(permissionResponse, "获取权限列表失败") ?? []);
      setTenants(unwrapResponse(tenantResponse, "获取租户列表失败").list ?? []);
    } catch (error) {
      message.error(getErrorMessage(error, "初始化角色页数据失败"));
    }
  };

  const loadRoles = async (options?: { keyword?: string; page?: number; pageSize?: number }) => {
    setLoading(true);
    try {
      const nextKeyword = options?.keyword ?? keyword;
      const response = await rolesList({
        ...(nextKeyword ? { keyword: nextKeyword } : {}),
        page: options?.page ?? pagination.page,
        pageSize: options?.pageSize ?? pagination.pageSize,
      });
      const data = unwrapResponse(response, "获取角色列表失败");
      setRoles(data.list ?? []);
      updatePageData(data);
    } catch (error) {
      message.error(getErrorMessage(error, "获取角色列表失败"));
    } finally {
      setLoading(false);
    }
  };

  useEffect(() => {
    void Promise.all([loadBaseData(), loadRoles({ keyword: "", page: 1 })]);
  }, []);

  const handleSearch = async (values: SearchFormValues) => {
    const nextKeyword = values.keyword?.trim() ?? "";
    setKeyword(nextKeyword);
    resetPage();
    await loadRoles({ keyword: nextKeyword, page: 1 });
  };

  const handleReset = async () => {
    searchForm.resetFields();
    setKeyword("");
    resetPage();
    await loadRoles({ keyword: "", page: 1 });
  };

  const handleOpenCreate = () => {
    setEditingRole(null);
    roleForm.resetFields();
    roleForm.setFieldsValue({
      permissionIds: [],
      scope: "tenant",
      tenantId: currentUser?.tenantId,
    });
    setDrawerOpen(true);
  };

  const handleOpenEdit = (target: API.Role) => {
    setEditingRole(target);
    roleForm.resetFields();
    roleForm.setFieldsValue({
      code: target.code,
      name: target.name,
      permissionIds: target.permissionIds,
      scope: target.scope as "system" | "tenant",
      tenantId: target.tenantId,
    });
    setDrawerOpen(true);
  };

  const handleCloseDrawer = () => {
    setDrawerOpen(false);
    setEditingRole(null);
    roleForm.resetFields();
  };

  const handleSubmit = async () => {
    try {
      const values = await roleForm.validateFields();
      setSubmitting(true);

      if (editingRole) {
        const updateResponse = await rolesUpdate(
          { id: editingRole.id },
          {
            name: values.name.trim(),
          },
        );
        unwrapResponse(updateResponse, "更新角色失败");

        const permissionResponse = await rolesPermissionsUpdate(
          { id: editingRole.id },
          {
            permissionIds: values.permissionIds,
          },
        );
        unwrapResponse(permissionResponse, "更新角色权限失败");
        message.success("角色更新成功");
      } else {
        const response = await rolesCreate({
          code: values.code?.trim() ?? "",
          name: values.name.trim(),
          permissionIds: values.permissionIds,
          scope: sysadmin ? values.scope : "tenant",
          tenantId: sysadmin && values.scope === "tenant" ? values.tenantId : currentUser?.tenantId,
        });
        unwrapResponse(response, "创建角色失败");
        message.success("角色创建成功");
      }

      handleCloseDrawer();
      await loadRoles();
    } catch (error) {
      if (typeof error === "object" && error !== null && "errorFields" in error) {
        return;
      }

      message.error(getErrorMessage(error, editingRole ? "更新角色失败" : "创建角色失败"));
    } finally {
      setSubmitting(false);
    }
  };

  const handleDelete = async (target: API.Role) => {
    try {
      setLoading(true);
      const response = await rolesDelete({ id: target.id });
      unwrapResponse(response, "删除角色失败");
      message.success("角色删除成功");
      const nextPage = getPageAfterDelete(pagination.total - 1, pagination.page, pagination.pageSize);
      await loadRoles({ page: nextPage });
    } catch (error) {
      message.error(getErrorMessage(error, "删除角色失败"));
      setLoading(false);
    }
  };

  const columns: TableColumnsType<API.Role> = [
    {
      title: "角色名称",
      dataIndex: "name",
      key: "name",
      width: 180,
    },
    {
      title: "角色编码",
      dataIndex: "code",
      key: "code",
      width: 180,
    },
    {
      title: "作用范围",
      dataIndex: "scope",
      key: "scope",
      width: 120,
      render: (value: string) => (value === "system" ? "系统级" : "租户级"),
    },
    {
      title: "所属租户",
      dataIndex: "tenantId",
      key: "tenantId",
      width: 180,
      render: (value?: number) => (value ? tenantMap[value] || `租户#${value}` : "-"),
    },
    {
      title: "类型",
      dataIndex: "builtin",
      key: "builtin",
      width: 120,
      render: (value: boolean) => <Tag color={value ? "blue" : "purple"}>{value ? "内置" : "自定义"}</Tag>,
    },
    {
      title: "权限数",
      dataIndex: "permissionIds",
      key: "permissionIds",
      width: 100,
      render: (value: string[]) => value.length,
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
      width: 180,
      render: (_, record) =>
        record.builtin ? (
          <Typography.Text type="secondary">内置角色不可修改</Typography.Text>
        ) : (
          <Space size={4} wrap>
            <Button type="link" className="ui-action-link ui-action-edit" onClick={() => handleOpenEdit(record)}>
              编辑
            </Button>
            <Popconfirm
              title="删除角色后不可恢复，确认删除吗？"
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
            <Form.Item<SearchFormValues> name="keyword" label="角色名称" className={styles.searchItem}>
              <Input allowClear placeholder="请输入角色名称或编码" />
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

          <Button type="primary" ghost icon={<PlusOutlined />} onClick={handleOpenCreate}>
            新建角色
          </Button>
        </div>

        <Table<API.Role>
          rowKey="id"
          loading={loading}
          columns={columns}
          dataSource={roles}
          className={styles.table}
          scroll={{ x: 1260 }}
          pagination={{
            current: pagination.page,
            pageSize: pagination.pageSize,
            total: pagination.total,
            showSizeChanger: true,
            pageSizeOptions: PAGE_SIZE_OPTIONS.map(String),
            showTotal: (total) => `共 ${total} 条`,
            onChange: (page, pageSize) => {
              handleTableChange(page, pageSize);
              void loadRoles({ page, pageSize });
            },
          }}
          locale={{
            emptyText: (
              <div className={styles.emptyState}>
                <Typography.Text type="secondary">暂无角色数据</Typography.Text>
              </div>
            ),
          }}
        />
      </div>

      <Drawer
        width={680}
        title={isEditMode ? "编辑角色" : "新建角色"}
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
        <Form<RoleFormValues>
          form={roleForm}
          layout="vertical"
          initialValues={{
            permissionIds: [],
            scope: "tenant",
            tenantId: currentUser?.tenantId,
          }}
          className={styles.drawerForm}
        >
          <Form.Item<RoleFormValues>
            label="角色名称"
            name="name"
            rules={[
              { required: true, message: "请输入角色名称" },
              { max: 128, message: "角色名称不能超过 128 个字符" },
            ]}
          >
            <Input placeholder="请输入角色名称" maxLength={128} />
          </Form.Item>

          <Form.Item<RoleFormValues>
            label="角色编码"
            name="code"
            rules={
              isEditMode
                ? []
                : [
                    { required: true, message: "请输入角色编码" },
                    { max: 64, message: "角色编码不能超过 64 个字符" },
                  ]
            }
          >
            <Input placeholder="请输入角色编码" maxLength={64} disabled={isEditMode} />
          </Form.Item>

          {sysadmin ? (
            <Form.Item<RoleFormValues> label="作用范围" name="scope">
              <Radio.Group
                options={[
                  { label: "租户级", value: "tenant" },
                  { label: "系统级", value: "system" },
                ]}
                disabled={isEditMode}
              />
            </Form.Item>
          ) : null}

          {scopeValue === "tenant" ? (
            <Form.Item<RoleFormValues>
              label="所属租户"
              name="tenantId"
              rules={[{ required: true, message: "请选择所属租户" }]}
            >
              <Select
                disabled={!sysadmin || isEditMode}
                options={
                  sysadmin
                    ? tenants.map((item) => ({ label: item.name, value: item.id }))
                    : currentUser?.tenantId
                      ? [{ label: `当前租户#${currentUser.tenantId}`, value: currentUser.tenantId }]
                      : []
                }
                placeholder="请选择所属租户"
              />
            </Form.Item>
          ) : null}

          <Form.Item<RoleFormValues>
            label="权限分配"
            name="permissionIds"
            rules={[{ required: true, message: "请至少选择一个权限" }]}
          >
            <Checkbox.Group className={styles.permissionGroup}>
              {Object.entries(groupedPermissions).map(([groupName, items]) => (
                <div key={groupName} className={styles.permissionSection}>
                  <Typography.Title level={5}>{groupName}</Typography.Title>
                  <Space direction="vertical" size={8} className={styles.permissionList}>
                    {items.map((item) => (
                      <Checkbox key={item.id} value={item.id}>
                        <span className={styles.permissionItem}>
                          <strong>{item.id}</strong>
                          <span>{`${item.method} ${item.path}`}</span>
                        </span>
                      </Checkbox>
                    ))}
                  </Space>
                </div>
              ))}
            </Checkbox.Group>
          </Form.Item>
        </Form>
      </Drawer>
    </section>
  );
}
