import { useEffect, useState } from "react";

import {
  CheckCircleOutlined,
  ClockCircleOutlined,
  SafetyCertificateOutlined,
} from "@ant-design/icons";
import {
  Card,
  Col,
  List,
  Progress,
  Row,
  Space,
  Tag,
  Timeline,
  Typography,
} from "antd";

import { MetricCard } from "@/components/MetricCard";
import { getDashboardSummary } from "@/services/clear-bill";
import type { DashboardSummary } from "@/services/clear-bill";
import { formatCurrency, formatPercent } from "@/utils/formatters";

export function OverviewPage() {
  const [summary, setSummary] = useState<DashboardSummary | null>(null);

  useEffect(() => {
    let active = true;

    const load = async () => {
      const data = await getDashboardSummary();
      if (active) {
        setSummary(data);
      }
    };

    void load();

    return () => {
      active = false;
    };
  }, []);

  const loading = summary === null;

  return (
    <div className="page-stack">
      <section className="hero-panel">
        <Space wrap size={[8, 8]}>
          <Tag color="cyan">Finance command surface</Tag>
          <Tag color="geekblue">Ant Design 6</Tag>
        </Space>
        <Typography.Title level={2}>
          One workspace for import quality, receivable risk, and reconciliation drift.
        </Typography.Title>
        <Typography.Paragraph type="secondary">
          The structure follows the reference repository while the frontend is modernized for a
          leaner UI stack and cleaner upgrade path.
        </Typography.Paragraph>
      </section>

      <Row gutter={[20, 20]}>
        <Col xs={24} sm={12} xl={6}>
          <MetricCard
            label="Open bills"
            value={loading ? "--" : String(summary.openBills)}
            delta="+12 today"
            hint="New billing batches waiting for release or manual review."
            tone="teal"
            percent={74}
          />
        </Col>
        <Col xs={24} sm={12} xl={6}>
          <MetricCard
            label="Pending invoices"
            value={loading ? "--" : String(summary.pendingInvoices)}
            delta="5 blocked"
            hint="Invoices pending compliance, tax review, or pricing confirmation."
            tone="amber"
            percent={58}
          />
        </Col>
        <Col xs={24} sm={12} xl={6}>
          <MetricCard
            label="Overdue amount"
            value={loading ? "--" : formatCurrency(summary.overdueAmount)}
            delta="2 high-risk"
            hint="Receivables already outside policy thresholds and needing escalation."
            tone="rose"
            percent={41}
          />
        </Col>
        <Col xs={24} sm={12} xl={6}>
          <MetricCard
            label="Auto-match rate"
            value={loading ? "--" : formatPercent(summary.autoMatchedRate)}
            delta="steady"
            hint="Payment records matched to invoices without manual intervention."
            tone="violet"
            percent={loading ? 0 : Math.round(summary.autoMatchedRate)}
          />
        </Col>
      </Row>

      <Row gutter={[20, 20]}>
        <Col xs={24} xl={12}>
          <Card className="panel-card" title="Attention queue" loading={loading}>
            <List
              dataSource={summary?.attentionList ?? []}
              renderItem={(item) => (
                <List.Item>
                  <List.Item.Meta
                    title={item.title}
                    description={
                      <Space direction="vertical" size={6} style={{ width: "100%" }}>
                        <Typography.Text type="secondary">
                          Owner: {item.owner} / Due: {item.dueDate}
                        </Typography.Text>
                        <Progress percent={item.progress} size="small" />
                      </Space>
                    }
                  />
                </List.Item>
              )}
            />
          </Card>
        </Col>

        <Col xs={24} xl={12}>
          <Card className="panel-card" title="Recent activity" loading={loading}>
            <Timeline
              items={(summary?.recentActivity ?? []).map((item) => ({
                color:
                  item.status === "success"
                    ? "green"
                    : item.status === "warning"
                      ? "orange"
                      : "blue",
                children: (
                  <Space direction="vertical" size={2}>
                    <Typography.Text strong>{item.title}</Typography.Text>
                    <Typography.Text type="secondary">{item.description}</Typography.Text>
                    <Typography.Text type="secondary">{item.time}</Typography.Text>
                  </Space>
                ),
              }))}
            />
          </Card>
        </Col>
      </Row>

      <Card className="panel-card" title="Operating posture">
        <Row gutter={[20, 20]}>
          <Col xs={24} md={8}>
            <Space align="start" size={12}>
              <CheckCircleOutlined style={{ color: "#0f8b8d", fontSize: 20 }} />
              <div>
                <Typography.Text strong>Release with confidence</Typography.Text>
                <Typography.Paragraph type="secondary">
                  Gate outgoing statements through a clearer approval path before they become
                  overdue collections work.
                </Typography.Paragraph>
              </div>
            </Space>
          </Col>
          <Col xs={24} md={8}>
            <Space align="start" size={12}>
              <ClockCircleOutlined style={{ color: "#d97706", fontSize: 20 }} />
              <div>
                <Typography.Text strong>Spot close-week risk early</Typography.Text>
                <Typography.Paragraph type="secondary">
                  Surface mismatches, stale approvals, and missing invoice detail before the end of
                  period squeeze.
                </Typography.Paragraph>
              </div>
            </Space>
          </Col>
          <Col xs={24} md={8}>
            <Space align="start" size={12}>
              <SafetyCertificateOutlined style={{ color: "#5b5bd6", fontSize: 20 }} />
              <div>
                <Typography.Text strong>Keep audit trails intact</Typography.Text>
                <Typography.Paragraph type="secondary">
                  Separate UI, API, scripts, and tools with the same higher-level structure as the
                  reference system.
                </Typography.Paragraph>
              </div>
            </Space>
          </Col>
        </Row>
      </Card>
    </div>
  );
}
