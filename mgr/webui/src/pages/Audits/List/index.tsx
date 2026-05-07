import { useEffect, useMemo, useState } from "react";

import { App, Button, DatePicker, Descriptions, Drawer, Form, Input, Pagination, Space, Tag, Timeline, Typography } from "antd";
import type { DescriptionsProps } from "antd";
import type { CSSProperties } from "react";

import { auditLogsList } from "@/services/clear-bill/audit";
import { formatDateTime, getErrorMessage, unwrapResponse } from "@/utils/api";
import { PAGE_SIZE_OPTIONS, useTablePagination } from "@/utils/pagination";

import styles from "./index.module.css";

interface SearchFormValues {
  endTime?: unknown;
  operation?: string;
  startTime?: unknown;
  user?: string;
}

function formatQueryTime(value?: unknown) {
  if (!value || typeof value !== "object" || !("format" in value)) {
    return undefined;
  }

  return String((value as { format: (pattern: string) => string }).format("YYYY-MM-DD HH:mm:ss"));
}

function getResultMeta(result?: string) {
  if (!result) {
    return {
      accent: "#64748b",
      bg: "rgba(100, 116, 139, 0.08)",
      color: "default",
      label: "-",
      tone: "default" as const,
    };
  }
  if (result.startsWith("failed")) {
    return {
      accent: "#ef4444",
      bg: "rgba(239, 68, 68, 0.08)",
      color: "error",
      label: "失败",
      tone: "error" as const,
    };
  }
  return {
    accent: "#22c55e",
    bg: "rgba(34, 197, 94, 0.08)",
    color: "success",
    label: "成功",
    tone: "success" as const,
  };
}

function buildDetailItems(item: API.AuditLog | null): DescriptionsProps["items"] {
  if (!item) {
    return [];
  }

  const result = getResultMeta(item.result);
  return [
    { key: "user", label: "用户", children: item.user },
    { key: "operation", label: "操作", children: item.operation },
    { key: "occurredAt", label: "时间", children: formatDateTime(item.occurredAt) },
    { key: "resource", label: "资源", children: item.resource },
    { key: "result", label: "结果", children: <Tag color={result.color}>{result.label}</Tag> },
  ];
}

export function AuditListPage() {
  const { message } = App.useApp();
  const [searchForm] = Form.useForm<SearchFormValues>();
  const [items, setItems] = useState<API.AuditLog[]>([]);
  const [loading, setLoading] = useState(false);
  const [detailOpen, setDetailOpen] = useState(false);
  const [detailItem, setDetailItem] = useState<API.AuditLog | null>(null);
  const [filters, setFilters] = useState<{
    endTime?: string;
    operation?: string;
    startTime?: string;
    user?: string;
  }>({});
  const { pagination, resetPage, updatePageData, handleTableChange } = useTablePagination();

  const loadItems = async (options?: {
    endTime?: string;
    operation?: string;
    page?: number;
    pageSize?: number;
    startTime?: string;
    user?: string;
  }) => {
    setLoading(true);
    try {
      const nextFilters = {
        user: options?.user ?? filters.user,
        operation: options?.operation ?? filters.operation,
        startTime: options?.startTime ?? filters.startTime,
        endTime: options?.endTime ?? filters.endTime,
      };
      const response = await auditLogsList({
        ...(nextFilters.user ? { user: nextFilters.user } : {}),
        ...(nextFilters.operation ? { operation: nextFilters.operation } : {}),
        ...(nextFilters.startTime ? { startTime: nextFilters.startTime } : {}),
        ...(nextFilters.endTime ? { endTime: nextFilters.endTime } : {}),
        page: options?.page ?? pagination.page,
        pageSize: options?.pageSize ?? pagination.pageSize,
      });
      const data = unwrapResponse<API.PageResult<API.AuditLog>>(response, "获取审计日志失败");
      setItems(data.list ?? []);
      updatePageData(data);
    } catch (error) {
      message.error(getErrorMessage(error, "获取审计日志失败"));
    } finally {
      setLoading(false);
    }
  };

  useEffect(() => {
    void loadItems({ page: 1 });
  }, []);

  const handleSearch = async (values: SearchFormValues) => {
    const nextFilters = {
      user: values.user?.trim() ?? "",
      operation: values.operation?.trim() ?? "",
      startTime: formatQueryTime(values.startTime),
      endTime: formatQueryTime(values.endTime),
    };
    setFilters(nextFilters);
    resetPage();
    await loadItems({ ...nextFilters, page: 1 });
  };

  const handleReset = async () => {
    searchForm.resetFields();
    setFilters({});
    resetPage();
    await loadItems({ page: 1, user: "", operation: "", startTime: "", endTime: "" });
  };

  const openDetail = (item: API.AuditLog) => {
    setDetailItem(item);
    setDetailOpen(true);
  };

  const timelineItems = useMemo(() => {
    return items.map((item) => {
      const result = getResultMeta(item.result);
      const cardStyle = {
        ["--audit-accent" as string]: result.accent,
        ["--audit-bg" as string]: result.bg,
      } as CSSProperties;
      const toneClass =
        result.tone === "error" ? styles.error : result.tone === "success" ? styles.success : styles.default;

      return {
        color: result.tone === "error" ? "red" : result.tone === "success" ? "green" : "gray",
        dot: <span className={`${styles.timelineDot} ${toneClass}`} />,
        children: (
          <button className={`${styles.entryCard} ${toneClass}`} style={cardStyle} type="button" onClick={() => openDetail(item)}>
            <span className={styles.entryStripe} />
            <div className={styles.entryHeader}>
              <div>
                <Typography.Text className={styles.entryUser}>{item.user}</Typography.Text>
                <Typography.Title level={5} className={styles.entryOperation}>
                  {item.operation}
                </Typography.Title>
              </div>
              <Tag color={result.color} className={styles.entryTag}>
                {result.label}
              </Tag>
            </div>

            <div className={styles.entryMeta}>
              <span>{formatDateTime(item.occurredAt)}</span>
              <span>{item.resource}</span>
            </div>

            <Typography.Paragraph className={styles.entryResource} ellipsis={{ rows: 2 }}>
              {item.resource}
            </Typography.Paragraph>

            <div className={styles.entryFooter}>
              <Button type="link" className={styles.detailButton}>
                查看详情
              </Button>
            </div>
          </button>
        ),
      };
    });
  }, [items]);

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
            <Form.Item<SearchFormValues> name="user" label="用户" className={styles.searchItem}>
              <Input allowClear placeholder="请输入用户名" />
            </Form.Item>
            <Form.Item<SearchFormValues> name="operation" label="操作" className={styles.searchItem}>
              <Input allowClear placeholder="请输入操作类型" />
            </Form.Item>
            <Form.Item<SearchFormValues> name="startTime" label="开始时间" className={styles.searchItem}>
              <DatePicker showTime className={styles.datePicker} />
            </Form.Item>
            <Form.Item<SearchFormValues> name="endTime" label="结束时间" className={styles.searchItem}>
              <DatePicker showTime className={styles.datePicker} />
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
        </div>

        <div className={styles.summaryBar}>
          <Typography.Text type="secondary">共 {pagination.total} 条记录，点击条目可查看详细信息</Typography.Text>
        </div>

        {timelineItems.length > 0 ? (
          <Timeline className={styles.timeline} items={timelineItems} />
        ) : (
          <div className={styles.emptyState}>
            <Typography.Text type="secondary">暂无审计日志数据</Typography.Text>
          </div>
        )}

        <div className={styles.paginationBar}>
          <Pagination
            current={pagination.page}
            pageSize={pagination.pageSize}
            total={pagination.total}
            showSizeChanger
            pageSizeOptions={PAGE_SIZE_OPTIONS.map(String)}
            showTotal={(total) => `共 ${total} 条`}
            onChange={(page, pageSize) => {
              handleTableChange(page, pageSize);
              void loadItems({ page, pageSize });
            }}
          />
        </div>
      </div>

      <Drawer
        width={640}
        title={detailItem ? `审计详情 - ${detailItem.operation}` : "审计详情"}
        open={detailOpen}
        onClose={() => {
          setDetailOpen(false);
          setDetailItem(null);
        }}
        destroyOnClose
      >
        {detailItem ? (
          <Space direction="vertical" size={20} className={styles.drawerBody}>
            <Descriptions items={buildDetailItems(detailItem)} bordered column={1} size="small" />

            <div className={styles.detailBlock}>
              <Typography.Title level={5}>资源</Typography.Title>
              <Typography.Paragraph className={styles.detailText} copyable>
                {detailItem.resource}
              </Typography.Paragraph>
            </div>

            <div className={styles.detailBlock}>
              <Typography.Title level={5}>结果</Typography.Title>
              <Tag color={getResultMeta(detailItem.result).color}>{getResultMeta(detailItem.result).label}</Tag>
            </div>
          </Space>
        ) : null}
      </Drawer>
    </section>
  );
}
