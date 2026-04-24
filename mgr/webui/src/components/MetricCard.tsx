import { Card, Flex, Progress, Tag, Typography } from "antd";

type Tone = "teal" | "amber" | "violet" | "rose";

const toneColorMap: Record<Tone, string> = {
  teal: "#0f8b8d",
  amber: "#d97706",
  violet: "#5b5bd6",
  rose: "#d9485f",
};

export interface MetricCardProps {
  label: string;
  value: string;
  delta: string;
  hint: string;
  tone: Tone;
  percent?: number;
}

export function MetricCard({
  label,
  value,
  delta,
  hint,
  tone,
  percent,
}: MetricCardProps) {
  const color = toneColorMap[tone];

  return (
    <Card className="metric-card" hoverable>
      <Flex vertical gap={20}>
        <Flex justify="space-between" align="center" wrap="wrap" gap={12}>
          <Tag style={{ color, background: `${color}12`, borderColor: "transparent" }}>
            {label}
          </Tag>
          <Typography.Text style={{ color }}>{delta}</Typography.Text>
        </Flex>

        <div>
          <Typography.Title level={2} className="metric-card__value">
            {value}
          </Typography.Title>
          <Typography.Text type="secondary" className="metric-card__hint">
            {hint}
          </Typography.Text>
        </div>

        {typeof percent === "number" ? <Progress percent={percent} showInfo={false} /> : null}
      </Flex>
    </Card>
  );
}
