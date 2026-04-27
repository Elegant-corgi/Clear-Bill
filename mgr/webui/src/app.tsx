import { App as AntApp, ConfigProvider, theme } from "antd";
import type { ThemeConfig } from "antd";
import { RouterProvider } from "react-router-dom";

import defaultSettings from "@config/defaultSettings";
import { router } from "@/routes";

const appTheme: ThemeConfig = {
  algorithm: theme.defaultAlgorithm,
  token: {
    colorPrimary: defaultSettings.primaryColor,
    colorInfo: defaultSettings.primaryColor,
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
  if (typeof document !== "undefined") {
    document.title = defaultSettings.title;
  }

  return (
    <ConfigProvider theme={appTheme}>
      <AntApp>
        <RouterProvider router={router} />
      </AntApp>
    </ConfigProvider>
  );
}
