import { App as AntApp, ConfigProvider, theme } from "antd";
import type { ThemeConfig } from "antd";
import { RouterProvider } from "react-router-dom";

import { router } from "@/routes";

const appTheme: ThemeConfig = {
  algorithm: theme.defaultAlgorithm,
  token: {
    colorPrimary: "#20c787",
    colorInfo: "#3f8cff",
    colorSuccess: "#20c787",
    colorWarning: "#ffb545",
    colorError: "#ff4d4f",
    borderRadius: 8,
    colorBgBase: "#f5f7fb",
    colorTextBase: "#17233d",
    fontFamily: '"Avenir Next", "PingFang SC", "Microsoft YaHei", sans-serif',
  },
};

export default function App() {
  return (
    <ConfigProvider theme={appTheme}>
      <AntApp>
        <RouterProvider router={router} />
      </AntApp>
    </ConfigProvider>
  );
}
