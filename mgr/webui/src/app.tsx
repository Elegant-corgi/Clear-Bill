import { App as AntApp, ConfigProvider, theme } from "antd";
import type { ThemeConfig } from "antd";
import { RouterProvider } from "react-router-dom";

import defaultSettings from "@config/defaultSettings";
import { AuthProvider } from "@/auth/AuthContext";
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
  components: {
    Input: {
      activeBg: "#ffffff",
      hoverBg: "#ffffff",
      addonBg: "#ffffff",
      colorBgContainer: "#ffffff",
      colorTextPlaceholder: "#a0aec0",
      activeBorderColor: defaultSettings.primaryColor,
      hoverBorderColor: "#cfd8e3",
    },
    Table: {
      colorBgContainer: "#ffffff",
      headerBg: "#f8fbff",
      headerColor: "#60708a",
      rowHoverBg: "#f9fcff",
      borderColor: "#edf2f7",
      footerBg: "#ffffff",
    },
    Select: {
      optionSelectedBg: "#eefaf4",
      optionActiveBg: "#f7fbf9",
      optionSelectedColor: "#17233d",
      selectorBg: "#ffffff",
      colorBgElevated: "#ffffff",
      colorTextPlaceholder: "#a0aec0",
      activeBorderColor: defaultSettings.primaryColor,
      hoverBorderColor: "#7fdab2",
    },
  },
};

export default function App() {
  if (typeof document !== "undefined") {
    document.title = defaultSettings.title;
  }

  return (
    <ConfigProvider theme={appTheme}>
      <AntApp>
        <AuthProvider>
          <RouterProvider router={router} />
        </AuthProvider>
      </AntApp>
    </ConfigProvider>
  );
}
