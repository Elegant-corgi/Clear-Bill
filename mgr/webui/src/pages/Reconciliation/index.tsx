import { useEffect, useState } from "react";

import { Card, Col, Progress, Row, Space, Table, Tag, Timeline, Typography } from "antd";
import type { TableProps } from "antd";

import { MetricCard } from "@/components/MetricCard";
import { getReconciliationTasks } from "@/services/clear-bill";
import type { ReconciliationStatus, ReconciliationTask } from "@/services/clear-bill";
import { formatPercent } from "@/utils/formatters";

const statusColorMap: Record<ReconciliationStatus, string> = {
  running: "processing",
  done: "success",
  attention: "warning",
};

export function ReconciliationPage() {
  const [tasks, setTasks] = useState<ReconciliationTask[]>([]);

  useEffect(() => {
    let active = true;

    const load = async () => {
      const data = await getReconciliationTasks();
      if (active) {
        setTasks(data);
      }
    };

    void load();

    return () => {
      active = false;
    };
  }, []);

  const matchedTotal = tasks.reduce((sum, task) => sum + task.matched, 0);
  const transactionTotal = tasks.reduce((sum, task) => sum + task.total, 0);
  const matchRate = transactionTotal === 0 ? 0 : (matchedTotal / transactionTotal) * 100;
  const activeFlows = tasks.filter((task) => task.status === "running").length;
  const escalations = tasks.filter((task) => task.status === "attention").length;

  const columns: TableProps<ReconciliationTask>["columns"] = [
    {
      title: "Task",
      dataIndex: "id",
      key: "id",
      render: (value: string) => <Typography.Text strong>{value}</Typography.Text>,
    },
    {
      title: "Bank",
      dataIndex: "bank",
      key: "bank",
    },
    {
      title: "Period",
      dataIndex: "period",
      key: "period",
    },
    {
      title: "Match progress",
      key: "progress",
      render: (_, record) => (
        <Progress percent={Math.round((record.matched / record.total) * 100)} size="small" />
      ),
    },
    {
      title: "Owner",
      dataIndex: "owner",
      key: "owner",
    },
    {
      title: "Status",
      dataIndex: "status",
      key: "status",
      render: (value: ReconciliationStatus) => <Tag color={statusColorMap[value]}>{value}</Tag>,
    },
  ];

  return (
    <div className="page-stack">
      <Row gutter={[20, 20]}>
        <Col xs={24} sm={12} xl={4}>
          <MetricCard
            label="Tasks"
            value={String(tasks.length)}
            delta="daily board"
            hint="Open reconciliation batches across connected settlement channels."
            tone="teal"
          />
        </Col>
        <Col xs={24} sm={12} xl={4}>
          <MetricCard
            label="Match rate"
            value={formatPercent(matchRate)}
            delta="bank feed"
            hint="Transactions linked to invoices or expected receivable records."
            tone="violet"
            percent={Math.round(matchRate)}
          />
        </Col>
        <Col xs={24} sm={12} xl={4}>
          <MetricCard
            label="Running"
            value={String(activeFlows)}
            delta="in motion"
            hint="Settlement flows currently processing in the active window."
            tone="amber"
          />
        </Col>
        <Col xs={24} sm={12} xl={4}>
          <MetricCard
            label="Escalations"
            value={String(escalations)}
            delta="needs review"
            hint="Batches with drift large enough to require operator intervention."
            tone="rose"
          />
        </Col>
      </Row>

      <Row gutter={[20, 20]}>
        <Col xs={24} xl={15}>
          <Card className="panel-card" title="Settlement board">
            <Table rowKey="id" columns={columns} dataSource={tasks} pagination={false} />
          </Card>
        </Col>

        <Col xs={24} xl={9}>
          <Card className="panel-card" title="Operating rhythm">
            <Space direction="vertical" size={18} style={{ width: "100%" }}>
              <Typography.Paragraph type="secondary">
                Run the bank import first, let the matcher sweep the straightforward records, and
                hold only the real exceptions for finance ops review.
              </Typography.Paragraph>
              <Timeline
                items={[
                  {
                    color: "blue",
                    children: "08:00 - Bank statements imported and normalized.",
                  },
                  {
                    color: "green",
                    children: "08:20 - Auto-match groups payment records to invoice candidates.",
                  },
                  {
                    color: "orange",
                    children: "09:00 - Attention queue isolates tolerance breaks and duplicate risk.",
                  },
                ]}
              />
            </Space>
          </Card>
        </Col>
      </Row>
    </div>
  );
}
