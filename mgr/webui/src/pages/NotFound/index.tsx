import { useNavigate } from "react-router-dom";

import { Button, Result } from "antd";

export function NotFoundPage() {
  const navigate = useNavigate();

  return (
    <Result
      status="404"
      title="Route not found"
      subTitle="The requested page is outside the current Clear Bill scaffold."
      extra={
        <Button type="primary" onClick={() => navigate("/")}>
          Back to overview
        </Button>
      }
    />
  );
}
