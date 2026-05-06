import { useEffect, useMemo, useState } from "react";

import {
  App,
  Button,
  Drawer,
  Divider,
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
import { CopyOutlined, PlusOutlined } from "@ant-design/icons";
import { useNavigate, useParams } from "react-router-dom";

import { CredentialDetail } from "./CredentialDetail";
import {
  credentialsCreate,
  credentialsDelete,
  credentialsGet,
  credentialsList,
  credentialsRotate,
} from "@/services/clear-bill/credential";
import { formatDateTime, getErrorMessage, unwrapResponse } from "@/utils/api";
import { PAGE_SIZE_OPTIONS, getPageAfterDelete, useTablePagination } from "@/utils/pagination";

import styles from "./index.module.css";

interface SearchFormValues {
  status?: string;
  type?: string;
}

interface CredentialFormValues {
  name: string;
  type: string;
}

interface RotateFormValues {
  name?: string;
}

const TYPE_OPTIONS = [
  { label: "Token", value: "token" },
  { label: "AK/SK", value: "aksk" },
];

const STATUS_OPTIONS = [
  { label: "全部", value: "" },
  { label: "启用", value: "active" },
  { label: "已轮转", value: "rotated" },
  { label: "已删除", value: "revoked" },
];

function getTypeLabel(type?: string) {
  if (type === "aksk") return "AK/SK";
  if (type === "token") return "Token";
  return type || "-";
}

function getStatusMeta(status?: string) {
  switch (status) {
    case "active":
      return { color: "success", label: "启用" };
    case "rotated":
      return { color: "processing", label: "已轮转" };
    case "revoked":
      return { color: "default", label: "已删除" };
    default:
      return { color: "default", label: status || "-" };
  }
}

function buildTokenExample(token?: string) {
  return `curl -H "X-API-Token: ${token || "<token>"}" \\
  http://127.0.0.1:8080/api/v1/bills`;
}

function buildAkSkExample(accessKey?: string, secretKey?: string) {
  return `# 直连方式
curl -H "X-Access-Key: ${accessKey || "<accessKey>"}" \\
  -H "X-Secret-Key: ${secretKey || "<secretKey>"}" \\
  http://127.0.0.1:8080/api/v1/bills

# 签名方式
# Header: X-Access-Key / X-Timestamp / X-Signature`;
}

export function CredentialListPage() {
  const { message, modal } = App.useApp();
  const navigate = useNavigate();
  const params = useParams();
  const [searchForm] = Form.useForm<SearchFormValues>();
  const [credentialForm] = Form.useForm<CredentialFormValues>();
  const [rotateForm] = Form.useForm<RotateFormValues>();
  const [items, setItems] = useState<API.Credential[]>([]);
  const [typeFilter, setTypeFilter] = useState("");
  const [statusFilter, setStatusFilter] = useState("");
  const [loading, setLoading] = useState(false);
  const [submitting, setSubmitting] = useState(false);
  const [drawerOpen, setDrawerOpen] = useState(false);
  const [rotateOpen, setRotateOpen] = useState(false);
  const [detailOpen, setDetailOpen] = useState(false);
  const [rotatingItem, setRotatingItem] = useState<API.Credential | null>(null);
  const [detailItem, setDetailItem] = useState<API.Credential | null>(null);
  const { pagination, resetPage, updatePageData, handleTableChange } = useTablePagination();

  const currentDetailId = useMemo(() => {
    const id = Number(params.id);
    return Number.isFinite(id) && id > 0 ? id : 0;
  }, [params.id]);

  const copyText = async (value?: string) => {
    if (!value) return;
    await navigator.clipboard.writeText(value);
    message.success("复制成功");
  };

  const loadItems = async (options?: { type?: string; status?: string; page?: number; pageSize?: number }) => {
    setLoading(true);
    try {
      const nextType = options?.type ?? typeFilter;
      const nextStatus = options?.status ?? statusFilter;
      const response = await credentialsList({
        ...(nextType ? { type: nextType } : {}),
        ...(nextStatus ? { status: nextStatus } : {}),
        page: options?.page ?? pagination.page,
        pageSize: options?.pageSize ?? pagination.pageSize,
      });
      const data = unwrapResponse<API.PageResult<API.Credential>>(response, "获取凭证列表失败");
      setItems(data.list ?? []);
      updatePageData(data);
    } catch (error) {
      message.error(getErrorMessage(error, "获取凭证列表失败"));
    } finally {
      setLoading(false);
    }
  };

  const openDetail = async (credential: API.Credential) => {
    setDetailItem(credential);
    setDetailOpen(true);
    void navigate(`/credentials/${credential.id}`);
    try {
      const response = await credentialsGet({ id: credential.id });
      const data = unwrapResponse<API.Credential>(response, "获取凭证详情失败");
      setDetailItem(data);
    } catch (error) {
      message.error(getErrorMessage(error, "获取凭证详情失败"));
    }
  };

  const openDetailById = async (id: number) => {
    const matched = items.find((item) => item.id === id);
    if (matched) {
      await openDetail(matched);
    }
  };

  useEffect(() => {
    void loadItems({ page: 1 });
  }, []);

  useEffect(() => {
    if (currentDetailId > 0) {
      void openDetailById(currentDetailId);
    }
  }, [currentDetailId]);

  const handleSearch = async (values: SearchFormValues) => {
    const nextType = values.type?.trim() ?? "";
    const nextStatus = values.status?.trim() ?? "";
    setTypeFilter(nextType);
    setStatusFilter(nextStatus);
    resetPage();
    await loadItems({ type: nextType, status: nextStatus, page: 1 });
  };

  const handleReset = async () => {
    searchForm.resetFields();
    setTypeFilter("");
    setStatusFilter("");
    resetPage();
    await loadItems({ type: "", status: "", page: 1 });
  };

  const handleOpenCreate = () => {
    credentialForm.resetFields();
    credentialForm.setFieldsValue({ type: "token" });
    setDrawerOpen(true);
  };

  const handleCloseCreate = () => {
    setDrawerOpen(false);
    credentialForm.resetFields();
  };

  const handleSubmit = async () => {
    try {
      const values = await credentialForm.validateFields();
      setSubmitting(true);
      const response = await credentialsCreate({
        name: values.name.trim(),
        type: values.type,
      });
      const data = unwrapResponse<API.Credential>(response, "创建凭证失败");
      message.success("凭证创建成功");
      await openDetail(data);
      modal.success({
        title: "凭证创建成功",
        content: (
          <div className={styles.successContent}>
            <p>
              <span className={styles.successLabel}>类型：</span>
              <span className={styles.successValue}>{getTypeLabel(data.type)}</span>
            </p>
            <p>
              <span className={styles.successLabel}>名称：</span>
              <span className={styles.successValue}>{data.name}</span>
            </p>
            <p>
              <span className={styles.successLabel}>说明：</span>
              <span className={styles.successValue}>已自动打开详情页，可直接复制调用示例</span>
            </p>
          </div>
        ),
        okText: "知道了",
      });
      handleCloseCreate();
      await loadItems();
    } catch (error) {
      if (typeof error === "object" && error !== null && "errorFields" in error) {
        return;
      }
      message.error(getErrorMessage(error, "创建凭证失败"));
    } finally {
      setSubmitting(false);
    }
  };

  const handleOpenRotate = (item: API.Credential) => {
    setRotatingItem(item);
    rotateForm.resetFields();
    rotateForm.setFieldsValue({ name: item.name });
    setRotateOpen(true);
  };

  const handleRotate = async () => {
    if (!rotatingItem) return;

    try {
      const values = await rotateForm.validateFields();
      setSubmitting(true);
      const response = await credentialsRotate({ id: rotatingItem.id }, { name: values.name?.trim() });
      const data = unwrapResponse<API.Credential>(response, "轮转凭证失败");
      message.success("凭证轮转成功");
      await openDetail(data);
      modal.success({
        title: "凭证轮转成功",
        content: (
          <div className={styles.successContent}>
            <p>
              <span className={styles.successLabel}>名称：</span>
              <span className={styles.successValue}>{data.name}</span>
            </p>
            <p>
              <span className={styles.successLabel}>说明：</span>
              <span className={styles.successValue}>新凭证已生成，详情页里可以复制新调用参数</span>
            </p>
          </div>
        ),
        okText: "知道了",
      });
      setRotateOpen(false);
      setRotatingItem(null);
      rotateForm.resetFields();
      await loadItems();
    } catch (error) {
      if (typeof error === "object" && error !== null && "errorFields" in error) {
        return;
      }
      message.error(getErrorMessage(error, "轮转凭证失败"));
    } finally {
      setSubmitting(false);
    }
  };

  const handleDelete = async (item: API.Credential) => {
    try {
      setLoading(true);
      const response = await credentialsDelete({ id: item.id });
      unwrapResponse(response, "删除凭证失败");
      message.success("凭证删除成功");
      if (detailItem?.id === item.id) {
        setDetailOpen(false);
        setDetailItem(null);
        void navigate("/credentials/list");
      }
      const nextPage = getPageAfterDelete(pagination.total - 1, pagination.page, pagination.pageSize);
      await loadItems({ page: nextPage });
    } catch (error) {
      message.error(getErrorMessage(error, "删除凭证失败"));
      setLoading(false);
    }
  };

  const columns: TableColumnsType<API.Credential> = [
    { title: "名称", dataIndex: "name", key: "name", width: 180 },
    { title: "类型", dataIndex: "type", key: "type", width: 120, render: (value: string) => getTypeLabel(value) },
    {
      title: "状态",
      dataIndex: "status",
      key: "status",
      width: 120,
      render: (value: string) => {
        const status = getStatusMeta(value);
        return <Tag color={status.color}>{status.label}</Tag>;
      },
    },
    { title: "AccessKey", dataIndex: "accessKeyPreview", key: "accessKeyPreview", width: 220, render: (value: string) => value || "-" },
    { title: "Token", dataIndex: "tokenPreview", key: "tokenPreview", width: 220, render: (value: string) => value || "-" },
    { title: "最后使用", dataIndex: "lastUsedAt", key: "lastUsedAt", width: 180, render: (value: string) => formatDateTime(value) },
    { title: "创建时间", dataIndex: "createdAt", key: "createdAt", width: 180, render: (value: string) => formatDateTime(value) },
    { title: "更新时间", dataIndex: "updatedAt", key: "updatedAt", width: 180, render: (value: string) => formatDateTime(value) },
    {
      title: "操作",
      key: "actions",
      fixed: "right",
      width: 220,
      render: (_, record) => (
        <Space size={0} className={styles.actionGroup}>
          <Button type="link" className="ui-action-link ui-action-edit" onClick={() => void openDetail(record)}>
            详情
          </Button>
          <Button type="link" className="ui-action-link ui-action-edit" onClick={() => handleOpenRotate(record)}>
            轮转
          </Button>
          <Popconfirm
            title="删除后不可恢复，确认删除该凭证吗？"
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

  const detailExample = detailItem
    ? detailItem.type === "token"
      ? buildTokenExample(detailItem.token)
      : buildAkSkExample(detailItem.accessKey, detailItem.secretKey)
    : "";

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
            <Form.Item<SearchFormValues> name="type" label="类型" className={styles.searchItem}>
              <Select allowClear options={TYPE_OPTIONS} placeholder="全部类型" />
            </Form.Item>
            <Form.Item<SearchFormValues> name="status" label="状态" className={styles.searchItem}>
              <Select allowClear options={STATUS_OPTIONS} placeholder="全部状态" />
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
            创建凭证
          </Button>
        </div>

        <Table<API.Credential>
          rowKey="id"
          loading={loading}
          columns={columns}
          dataSource={items}
          className={styles.table}
          scroll={{ x: 1560 }}
          pagination={{
            current: pagination.page,
            pageSize: pagination.pageSize,
            total: pagination.total,
            showSizeChanger: true,
            pageSizeOptions: PAGE_SIZE_OPTIONS.map(String),
            showTotal: (total) => `共 ${total} 条`,
            onChange: (page, pageSize) => {
              handleTableChange(page, pageSize);
              void loadItems({ page, pageSize });
            },
          }}
          locale={{
            emptyText: (
              <div className={styles.emptyState}>
                <Typography.Text type="secondary">暂无凭证数据</Typography.Text>
              </div>
            ),
          }}
        />
      </div>

      <Drawer
        width={520}
        title="创建凭证"
        open={drawerOpen}
        onClose={handleCloseCreate}
        destroyOnClose
        footer={
          <Space>
            <Button onClick={handleCloseCreate}>取消</Button>
            <Button type="primary" loading={submitting} onClick={() => void handleSubmit()}>
              保存
            </Button>
          </Space>
        }
        styles={{ footer: { display: "flex", justifyContent: "flex-end" } }}
      >
        <Form<CredentialFormValues> form={credentialForm} layout="vertical" initialValues={{ type: "token" }} className={styles.drawerForm}>
          <Form.Item<CredentialFormValues>
            label="凭证名称"
            name="name"
            rules={[
              { required: true, message: "请输入凭证名称" },
              { max: 128, message: "凭证名称不能超过 128 个字符" },
            ]}
          >
            <Input placeholder="请输入凭证名称" maxLength={128} />
          </Form.Item>

          <Form.Item<CredentialFormValues>
            label="凭证类型"
            name="type"
            rules={[{ required: true, message: "请选择凭证类型" }]}
          >
            <Select options={TYPE_OPTIONS} />
          </Form.Item>
        </Form>
      </Drawer>

      <Drawer
        width={720}
        title={detailItem ? `凭证详情 - ${detailItem.name}` : "凭证详情"}
        open={detailOpen}
        onClose={() => {
          setDetailOpen(false);
          setDetailItem(null);
          void navigate("/credentials/list");
        }}
        destroyOnClose
      >
        <CredentialDetail
          credential={detailItem}
          onCopy={(value) => void copyText(value)}
        />
        {detailItem ? (
          <>
            <Divider />
            <Typography.Title level={5}>完整调用示例</Typography.Title>
            <Typography.Paragraph className={styles.codeBlock}>{detailExample}</Typography.Paragraph>
            <Button icon={<CopyOutlined />} onClick={() => void copyText(detailExample)}>
              复制示例
            </Button>
          </>
        ) : null}
      </Drawer>

      <Modal
        title={rotatingItem ? `轮转凭证 - ${rotatingItem.name}` : "轮转凭证"}
        open={rotateOpen}
        onCancel={() => {
          setRotateOpen(false);
          setRotatingItem(null);
          rotateForm.resetFields();
        }}
        onOk={() => void handleRotate()}
        okText="确认轮转"
        cancelText="取消"
        confirmLoading={submitting}
      >
        <Form<RotateFormValues> form={rotateForm} layout="vertical">
          <Form.Item<RotateFormValues> label="新名称" name="name">
            <Input placeholder="可留空，默认沿用原名称" maxLength={128} />
          </Form.Item>
          <Typography.Text type="secondary">轮转会生成一条新的凭证记录，旧凭证将标记为已轮转。</Typography.Text>
        </Form>
      </Modal>
    </section>
  );
}
