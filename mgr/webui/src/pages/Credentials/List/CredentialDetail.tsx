import { Button, Descriptions, Divider, Space, Typography } from "antd";
import { CopyOutlined } from "@ant-design/icons";

import { formatDateTime } from "@/utils/api";

import styles from "./index.module.css";

interface CredentialDetailProps {
  credential: API.Credential | null;
  onCopy: (value?: string) => void;
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
# Header: X-Access-Key / X-Timestamp / X-Signature
`;
}

export function CredentialDetail({ credential, onCopy }: CredentialDetailProps) {
  if (!credential) {
    return null;
  }

  const isToken = credential.type === "token";

  return (
    <div className={styles.detailPanel}>
      <Descriptions bordered column={1} size="small">
        <Descriptions.Item label="名称">{credential.name}</Descriptions.Item>
        <Descriptions.Item label="类型">{isToken ? "Token" : "AK/SK"}</Descriptions.Item>
        <Descriptions.Item label="状态">{credential.status}</Descriptions.Item>
        <Descriptions.Item label="创建时间">{formatDateTime(credential.createdAt)}</Descriptions.Item>
        <Descriptions.Item label="更新时间">{formatDateTime(credential.updatedAt)}</Descriptions.Item>
        <Descriptions.Item label="最后使用">{formatDateTime(credential.lastUsedAt)}</Descriptions.Item>
      </Descriptions>

      <Divider />

      <Typography.Title level={5}>调用示例</Typography.Title>
      {isToken ? (
        <>
          <Typography.Paragraph className={styles.codeBlock}>{buildTokenExample(credential.token)}</Typography.Paragraph>
          <Space>
            <Button icon={<CopyOutlined />} onClick={() => void onCopy(buildTokenExample(credential.token))}>
              复制示例
            </Button>
            <Button type="primary" icon={<CopyOutlined />} onClick={() => void onCopy(credential.token)}>
              复制 Token
            </Button>
          </Space>
        </>
      ) : (
        <>
          <Typography.Paragraph className={styles.codeBlock}>
            {buildAkSkExample(credential.accessKey, credential.secretKey)}
          </Typography.Paragraph>
          <Space>
            <Button icon={<CopyOutlined />} onClick={() => void onCopy(buildAkSkExample(credential.accessKey, credential.secretKey))}>
              复制示例
            </Button>
            <Button icon={<CopyOutlined />} onClick={() => void onCopy(credential.accessKey)}>
              复制 AccessKey
            </Button>
            <Button icon={<CopyOutlined />} onClick={() => void onCopy(credential.secretKey)}>
              复制 SecretKey
            </Button>
          </Space>
        </>
      )}
    </div>
  );
}
