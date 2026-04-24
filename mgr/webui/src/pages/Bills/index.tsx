import { useEffect, useState } from "react";

import { Card, Input, Segmented, Space, Table, Tag, Typography } from "antd";
import type { TableProps } from "antd";

import { getBills } from "@/services/clear-bill";
import type { BillRecord, BillStatus } from "@/services/clear-bill";
import { formatCurrency } from "@/utils/formatters";

type BillFilter = "all" | BillStatus;

const statusColorMap: Record<BillStatus, string> = {
  ready: "green",
  reviewing: "gold",
  overdue: "red",
  paid: "blue",
};

export function BillsPage() {
  const [records, setRecords] = useState<BillRecord[]>([]);
  const [query, setQuery] = useState("");
  const [status, setStatus] = useState<BillFilter>("all");

  useEffect(() => {
    let active = true;

    const load = async () => {
      const data = await getBills();
      if (active) {
        setRecords(data);
      }
    };

    void load();

    return () => {
      active = false;
    };
  }, []);

  const filteredRecords = records.filter((record) => {
    const matchesQuery =
      query.trim() === "" ||
      record.id.toLowerCase().includes(query.toLowerCase()) ||
      record.customerName.toLowerCase().includes(query.toLowerCase());

    const matchesStatus = status === "all" || record.status === status;

    return matchesQuery && matchesStatus;
  });

  const columns: TableProps<BillRecord>["columns"] = [
    {
      title: "Bill ID",
      dataIndex: "id",
      key: "id",
      render: (value: string) => <Typography.Text strong>{value}</Typography.Text>,
    },
    {
      title: "Customer",
      dataIndex: "customerName",
      key: "customerName",
    },
    {
      title: "Billing month",
      dataIndex: "billingMonth",
      key: "billingMonth",
    },
    {
      title: "Amount",
      dataIndex: "amount",
      key: "amount",
      align: "right",
      render: (value: number) => formatCurrency(value),
    },
    {
      title: "Channel",
      dataIndex: "channel",
      key: "channel",
      render: (value: string) => <Tag>{value}</Tag>,
    },
    {
      title: "Status",
      dataIndex: "status",
      key: "status",
      render: (value: BillStatus) => (
        <Tag className="table-chip" color={statusColorMap[value]}>
          {value}
        </Tag>
      ),
    },
    {
      title: "Updated at",
      dataIndex: "updatedAt",
      key: "updatedAt",
    },
  ];

  return (
    <div className="page-stack">
      <Card className="panel-card">
        <Space direction="vertical" size={16} style={{ width: "100%" }}>
          <div>
            <Typography.Title level={3}>Billing register</Typography.Title>
            <Typography.Text type="secondary">
              Filter the current batch, watch review bottlenecks, and spot overdue statements.
            </Typography.Text>
          </div>

          <Space wrap size={[12, 12]}>
            <Input.Search
              allowClear
              placeholder="Search by bill ID or customer"
              value={query}
              onChange={(event) => setQuery(event.target.value)}
              style={{ width: 320 }}
            />
            <Segmented
              value={status}
              onChange={(value) => setStatus(value as BillFilter)}
              options={[
                { label: "All", value: "all" },
                { label: "Ready", value: "ready" },
                { label: "Reviewing", value: "reviewing" },
                { label: "Overdue", value: "overdue" },
                { label: "Paid", value: "paid" },
              ]}
            />
          </Space>
        </Space>
      </Card>

      <Card className="panel-card" title="Current bills">
        <Table
          rowKey="id"
          columns={columns}
          dataSource={filteredRecords}
          pagination={{ pageSize: 6 }}
        />
      </Card>
    </div>
  );
}
