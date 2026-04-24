import { App as AntApp, ConfigProvider, theme } from "antd";
import type { ThemeConfig } from "antd";
import { RouterProvider } from "react-router-dom";

import { router } from "@/routes";

const appTheme: ThemeConfig = {
  algorithm: theme.defaultAlgorithm,
  token: {
    colorPrimary: "#0f766e",
    colorInfo: "#0f766e",
    colorSuccess: "#0f8b8d",
    colorWarning: "#d97706",
    colorError: "#dc2626",
    borderRadius: 22,
    colorBgBase: "#f7f7f2",
    colorTextBase: "#132238",
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
