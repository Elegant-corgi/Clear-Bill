import { useEffect, useState } from "react";

import { Avatar, Card, Col, Progress, Row, Space, Tag, Typography } from "antd";

import { getCustomers } from "@/services/clear-bill";
import type { CustomerSnapshot } from "@/services/clear-bill";
import { formatCurrency } from "@/utils/formatters";

export function CustomersPage() {
  const [customers, setCustomers] = useState<CustomerSnapshot[]>([]);

  useEffect(() => {
    let active = true;

    const load = async () => {
      const data = await getCustomers();
      if (active) {
        setCustomers(data);
      }
    };

    void load();

    return () => {
      active = false;
    };
  }, []);

  return (
    <div className="page-stack">
      <Card className="panel-card">
        <Typography.Title level={3}>Account health</Typography.Title>
        <Typography.Text type="secondary">
          Keep commercial exposure, contact ownership, and billing health visible at the same time.
        </Typography.Text>
      </Card>

      <Row gutter={[20, 20]}>
        {customers.map((customer) => (
          <Col xs={24} lg={12} xl={8} key={customer.id}>
            <Card className="customer-card" hoverable>
              <Space direction="vertical" size={20} style={{ width: "100%" }}>
                <Space align="start" size={14}>
                  <Avatar
                    size={56}
                    className="customer-card__avatar"
                    style={{ backgroundColor: "#10324a" }}
                  >
                    {customer.name
                      .split(" ")
                      .slice(0, 2)
                      .map((part) => part.charAt(0))
                      .join("")}
                  </Avatar>
                  <div>
                    <Typography.Title level={4} style={{ margin: 0 }}>
                      {customer.name}
                    </Typography.Title>
                    <Typography.Text type="secondary">
                      Primary contact: {customer.primaryContact}
                    </Typography.Text>
                  </div>
                </Space>

                <Space wrap size={[8, 8]}>
                  <Tag color="blue">Credit {customer.creditLevel}</Tag>
                  <Tag color="green">{customer.activeContracts} active contracts</Tag>
                </Space>

                <div>
                  <Typography.Text type="secondary">Outstanding amount</Typography.Text>
                  <Typography.Title level={3} style={{ marginTop: 6, marginBottom: 0 }}>
                    {formatCurrency(customer.outstandingAmount)}
                  </Typography.Title>
                </div>

                <div>
                  <Typography.Text type="secondary">Billing health</Typography.Text>
                  <Progress percent={customer.billingHealth} />
                </div>
              </Space>
            </Card>
          </Col>
        ))}
      </Row>
    </div>
  );
}
