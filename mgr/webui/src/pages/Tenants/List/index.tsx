import { useEffect, useState } from "react";

import {
  App,
  Button,
  Drawer,
  Form,
  Input,
  Popconfirm,
  Select,
  Space,
  Table,
  Tag,
  Typography,
} from "antd";
import type { TableColumnsType } from "antd";

import {
  tenantsCreate,
  tenantsDelete,
  tenantsList,
  tenantsUpdate,
} from "@/services/clear-bill/tenant";
import { formatDateTime, getErrorMessage, unwrapResponse } from "@/utils/api";

import styles from "./index.module.css";

interface SearchFormValues {
  keyword?: string;
}

interface TenantFormValues {
  adminDisplayName?: string;
  adminUsername?: string;
  code: string;
  contactName?: string;
  contactPhone?: string;
  name: string;
  remark?: string;
  status: string;
}

const DEFAULT_STATUS = "active";

const STATUS_META: Record<string, { color: string; label: string }> = {
  active: {
    color: "success",
    label: "启用",
  },
  inactive: {
    color: "error",
    label: "停用",
  },
};

const STATUS_OPTIONS = [
  { label: "启用", value: "active" },
  { label: "停用", value: "inactive" },
];

function getStatusMeta(status?: string) {
  if (!status) {
    return { color: "default", label: "-" };
  }

  return STATUS_META[status] ?? { color: "default", label: status };
}

function buildUpdatePayload(tenant: API.Tenant, status?: string): API.UpdateTenantReq {
  return {
    code: tenant.code,
    contactName: tenant.contactName || "",
    contactPhone: tenant.contactPhone || "",
    name: tenant.name,
    remark: tenant.remark || "",
    status: status ?? tenant.status ?? DEFAULT_STATUS,
  };
}

export function TenantListPage() {
  const { message, modal } = App.useApp();
  const [searchForm] = Form.useForm<SearchFormValues>();
  const [tenantForm] = Form.useForm<TenantFormValues>();
  const [tenants, setTenants] = useState<API.Tenant[]>([]);
  const [keyword, setKeyword] = useState("");
  const [loading, setLoading] = useState(false);
  const [submitting, setSubmitting] = useState(false);
  const [drawerOpen, setDrawerOpen] = useState(false);
  const [editingTenant, setEditingTenant] = useState<API.Tenant | null>(null);

  const isEditMode = Boolean(editingTenant);

  const loadTenants = async (nextKeyword = keyword) => {
    setLoading(true);
    try {
      const response = await tenantsList(nextKeyword ? { keyword: nextKeyword } : {});
      const data = unwrapResponse(response, "获取租户列表失败");
      setTenants(data ?? []);
    } catch (error) {
      message.error(getErrorMessage(error, "获取租户列表失败"));
    } finally {
      setLoading(false);
    }
  };

  useEffect(() => {
    void loadTenants("");
  }, []);

  const handleSearch = async (values: SearchFormValues) => {
    const nextKeyword = values.keyword?.trim() ?? "";
    setKeyword(nextKeyword);
    await loadTenants(nextKeyword);
  };

  const handleReset = async () => {
    searchForm.resetFields();
    setKeyword("");
    await loadTenants("");
  };

  const handleOpenCreate = () => {
    setEditingTenant(null);
    tenantForm.resetFields();
    tenantForm.setFieldsValue({
      status: DEFAULT_STATUS,
    });
    setDrawerOpen(true);
  };

  const handleOpenEdit = (tenant: API.Tenant) => {
    setEditingTenant(tenant);
    tenantForm.resetFields();
    tenantForm.setFieldsValue({
      code: tenant.code,
      contactName: tenant.contactName,
      contactPhone: tenant.contactPhone,
      name: tenant.name,
      remark: tenant.remark,
      status: tenant.status || DEFAULT_STATUS,
    });
    setDrawerOpen(true);
  };

  const handleCloseDrawer = () => {
    setDrawerOpen(false);
    setEditingTenant(null);
    tenantForm.resetFields();
  };

  const handleSubmit = async () => {
    try {
      const values = await tenantForm.validateFields();
      setSubmitting(true);

      if (editingTenant) {
        const response = await tenantsUpdate(
          { id: editingTenant.id },
          {
            code: values.code.trim(),
            contactName: values.contactName?.trim() ?? "",
            contactPhone: values.contactPhone?.trim() ?? "",
            name: values.name.trim(),
            remark: values.remark?.trim() ?? "",
            status: values.status || DEFAULT_STATUS,
          },
        );

        unwrapResponse(response, "更新租户失败");
        message.success("租户更新成功");
      } else {
        const response = await tenantsCreate({
          adminDisplayName: values.adminDisplayName?.trim() ?? "",
          adminUsername: values.adminUsername?.trim() ?? "",
          code: values.code.trim(),
          contactName: values.contactName?.trim() ?? "",
          contactPhone: values.contactPhone?.trim() ?? "",
          name: values.name.trim(),
          remark: values.remark?.trim() ?? "",
          status: values.status || DEFAULT_STATUS,
        });

        const data = unwrapResponse(response, "创建租户失败");
        message.success("租户创建成功");
        modal.success({
          title: "租户创建成功",
          content: (
            <div className={styles.successContent}>
              <p>
                租户名称：
                <strong>{data.tenant.name}</strong>
              </p>
              <p>
                管理员账号：
                <strong>{data.adminUsername}</strong>
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
      await loadTenants(keyword);
    } catch (error) {
      if (typeof error === "object" && error !== null && "errorFields" in error) {
        return;
      }

      message.error(getErrorMessage(error, editingTenant ? "更新租户失败" : "创建租户失败"));
    } finally {
      setSubmitting(false);
    }
  };

  const handleToggleStatus = async (tenant: API.Tenant, status: string) => {
    try {
      setLoading(true);
      const response = await tenantsUpdate({ id: tenant.id }, buildUpdatePayload(tenant, status));
      unwrapResponse(response, "更新租户状态失败");
      message.success(`租户已${status === "active" ? "启用" : "停用"}`);
      await loadTenants(keyword);
    } catch (error) {
      message.error(getErrorMessage(error, "更新租户状态失败"));
      setLoading(false);
    }
  };

  const handleDelete = async (tenant: API.Tenant) => {
    try {
      setLoading(true);
      const response = await tenantsDelete({ id: tenant.id });
      unwrapResponse(response, "删除租户失败");
      message.success("租户删除成功");
      await loadTenants(keyword);
    } catch (error) {
      message.error(getErrorMessage(error, "删除租户失败"));
      setLoading(false);
    }
  };

  const columns: TableColumnsType<API.Tenant> = [
    {
      title: "租户名称",
      dataIndex: "name",
      key: "name",
      width: 220,
      ellipsis: true,
    },
    {
      title: "租户编码",
      dataIndex: "code",
      key: "code",
      width: 180,
      ellipsis: true,
    },
    {
      title: "管理员账号",
      dataIndex: "adminUsername",
      key: "adminUsername",
      width: 180,
      render: (value: string) => value || "-",
    },
    {
      title: "管理员名称",
      dataIndex: "adminDisplayName",
      key: "adminDisplayName",
      width: 180,
      render: (value: string) => value || "-",
    },
    {
      title: "联系人",
      dataIndex: "contactName",
      key: "contactName",
      width: 140,
      render: (value: string) => value || "-",
    },
    {
      title: "联系电话",
      dataIndex: "contactPhone",
      key: "contactPhone",
      width: 150,
      render: (value: string) => value || "-",
    },
    {
      title: "状态",
      dataIndex: "status",
      key: "status",
      width: 110,
      render: (value: string) => {
        const status = getStatusMeta(value);
        return <Tag color={status.color}>{status.label}</Tag>;
      },
    },
    {
      title: "备注",
      dataIndex: "remark",
      key: "remark",
      width: 180,
      ellipsis: true,
      render: (value: string) => value || "-",
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
      width: 160,
      render: (_, record) => {
        const nextStatus = record.status === "active" ? "inactive" : "active";
        const actionText = nextStatus === "active" ? "启用" : "停用";

        return (
          <Space size={0} className={styles.actionGroup}>
            <Button type="link" onClick={() => handleOpenEdit(record)}>
              编辑
            </Button>
            <Popconfirm
              title={`确认${actionText}该租户吗？`}
              okText="确认"
              cancelText="取消"
              onConfirm={() => handleToggleStatus(record, nextStatus)}
            >
              <Button type="link">{actionText}</Button>
            </Popconfirm>
            <Popconfirm
              title="删除后不可恢复，确认删除该租户吗？"
              okText="删除"
              cancelText="取消"
              okButtonProps={{ danger: true }}
              onConfirm={() => handleDelete(record)}
            >
              <Button danger type="link">
                删除
              </Button>
            </Popconfirm>
          </Space>
        );
      },
    },
  ];

  return (
    <section className={styles.page}>
      <div className={styles.card}>
        <div className={styles.toolbar}>
          <Form<SearchFormValues>
            form={searchForm}
            layout="inline"
            onFinish={(values) => void handleSearch(values)}
            className={styles.searchForm}
          >
            <Form.Item<SearchFormValues> name="keyword" label="租户名称" className={styles.searchItem}>
              <Input allowClear placeholder="请输入租户名称" maxLength={128} />
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

          <Button type="primary" ghost onClick={handleOpenCreate}>
            新建租户
          </Button>
        </div>

        <Table<API.Tenant>
          rowKey="id"
          loading={loading}
          columns={columns}
          dataSource={tenants}
          className={styles.table}
          scroll={{ x: 1680 }}
          pagination={{
            pageSize: 10,
            showSizeChanger: false,
            showTotal: (total) => `共 ${total} 条`,
          }}
          locale={{
            emptyText: (
              <div className={styles.emptyState}>
                <Typography.Text type="secondary">暂无租户数据</Typography.Text>
              </div>
            ),
          }}
        />
      </div>

      <Drawer
        width={520}
        title={isEditMode ? "编辑租户" : "新建租户"}
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
        <Form<TenantFormValues>
          form={tenantForm}
          layout="vertical"
          initialValues={{ status: DEFAULT_STATUS }}
          className={styles.drawerForm}
        >
          <Form.Item<TenantFormValues>
            label="租户名称"
            name="name"
            rules={[
              { required: true, message: "请输入租户名称" },
              { max: 128, message: "租户名称不能超过 128 个字符" },
            ]}
          >
            <Input placeholder="请输入租户名称" maxLength={128} />
          </Form.Item>

          <Form.Item<TenantFormValues>
            label="租户编码"
            name="code"
            rules={[
              { required: true, message: "请输入租户编码" },
              { max: 64, message: "租户编码不能超过 64 个字符" },
            ]}
          >
            <Input placeholder="请输入租户编码" maxLength={64} />
          </Form.Item>

          {!isEditMode ? (
            <>
              <Form.Item<TenantFormValues>
                label="管理员账号"
                name="adminUsername"
                rules={[{ max: 64, message: "管理员账号不能超过 64 个字符" }]}
              >
                <Input placeholder="留空则按后端规则自动生成" maxLength={64} />
              </Form.Item>

              <Form.Item<TenantFormValues>
                label="管理员名称"
                name="adminDisplayName"
                rules={[{ max: 128, message: "管理员名称不能超过 128 个字符" }]}
              >
                <Input placeholder="留空则按后端规则自动生成" maxLength={128} />
              </Form.Item>
            </>
          ) : null}

          <Form.Item<TenantFormValues>
            label="联系人"
            name="contactName"
            rules={[{ max: 64, message: "联系人不能超过 64 个字符" }]}
          >
            <Input placeholder="请输入联系人" maxLength={64} />
          </Form.Item>

          <Form.Item<TenantFormValues>
            label="联系电话"
            name="contactPhone"
            rules={[{ max: 32, message: "联系电话不能超过 32 个字符" }]}
          >
            <Input placeholder="请输入联系电话" maxLength={32} />
          </Form.Item>

          <Form.Item<TenantFormValues> label="状态" name="status">
            <Select options={STATUS_OPTIONS} />
          </Form.Item>

          <Form.Item<TenantFormValues>
            label="备注"
            name="remark"
            rules={[{ max: 255, message: "备注不能超过 255 个字符" }]}
          >
            <Input.TextArea rows={4} placeholder="请输入备注" maxLength={255} showCount />
          </Form.Item>
        </Form>
      </Drawer>
    </section>
  );
}
