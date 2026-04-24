import { App, Button, Card, Form, Input, InputNumber, List, Select, Switch, Tag, Typography } from "antd";

interface WorkspaceSettings {
  owner: string;
  reminderChannel: string;
  billingWindow: string;
  approvalThreshold: number;
  autoReconcile: boolean;
}

const integrations = [
  { name: "Billing API", status: "healthy", tone: "green" },
  { name: "Tax invoice adapter", status: "review", tone: "gold" },
  { name: "Bank feed connector", status: "healthy", tone: "blue" },
];

export function SettingsPage() {
  const { message } = App.useApp();

  const onFinish = (values: WorkspaceSettings) => {
    void message.success(
      `Saved workspace defaults for ${values.owner} with ${values.billingWindow} close window.`,
    );
  };

  return (
    <div className="page-stack">
      <Card className="settings-card">
        <Typography.Title level={3}>Workspace defaults</Typography.Title>
        <Typography.Paragraph type="secondary">
          Configure a few guardrails for the scaffold so future modules have a sensible starting
          point.
        </Typography.Paragraph>

        <Form<WorkspaceSettings>
          layout="vertical"
          initialValues={{
            owner: "Finance Ops",
            reminderChannel: "email",
            billingWindow: "T+1 noon",
            approvalThreshold: 50000,
            autoReconcile: true,
          }}
          onFinish={onFinish}
        >
          <div className="settings-grid">
            <Form.Item label="Workspace owner" name="owner" rules={[{ required: true }]}>
              <Input placeholder="Finance Ops" />
            </Form.Item>

            <Form.Item label="Reminder channel" name="reminderChannel">
              <Select
                options={[
                  { label: "Email", value: "email" },
                  { label: "Slack", value: "slack" },
                  { label: "SMS", value: "sms" },
                ]}
              />
            </Form.Item>

            <Form.Item label="Billing close window" name="billingWindow">
              <Select
                options={[
                  { label: "Same day 18:00", value: "same-day 18:00" },
                  { label: "T+1 noon", value: "T+1 noon" },
                  { label: "T+2 morning", value: "T+2 morning" },
                ]}
              />
            </Form.Item>

            <Form.Item label="Manual approval threshold" name="approvalThreshold">
              <InputNumber min={1000} step={1000} style={{ width: "100%" }} />
            </Form.Item>
          </div>

          <Form.Item
            label="Enable automatic reconciliation"
            name="autoReconcile"
            valuePropName="checked"
          >
            <Switch />
          </Form.Item>

          <Button htmlType="submit" type="primary">
            Save defaults
          </Button>
        </Form>
      </Card>

      <Card className="panel-card" title="Integration snapshot">
        <List
          dataSource={integrations}
          renderItem={(item) => (
            <List.Item>
              <div className="integration-row">
                <div>
                  <Typography.Text strong>{item.name}</Typography.Text>
                  <br />
                  <Typography.Text type="secondary">
                    Baseline placeholder ready for future service bindings.
                  </Typography.Text>
                </div>
                <Tag color={item.tone}>{item.status}</Tag>
              </div>
            </List.Item>
          )}
        />
      </Card>
    </div>
  );
}
