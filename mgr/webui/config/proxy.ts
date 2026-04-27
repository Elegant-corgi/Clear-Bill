export interface ProxyItem {
  target: string;
  changeOrigin?: boolean;
  secure?: boolean;
  ws?: boolean;
}

export type ProxyConfig = Record<string, ProxyItem>;

const defaultTarget = "http://127.0.0.1:8080";

const proxy: ProxyConfig = {
  "/api": {
    target: process.env.VITE_API_PROXY_TARGET || defaultTarget,
    changeOrigin: process.env.VITE_API_PROXY_CHANGE_ORIGIN
      ? process.env.VITE_API_PROXY_CHANGE_ORIGIN === "true"
      : true,
    secure: process.env.VITE_API_PROXY_SECURE === "true",
    ws: process.env.VITE_API_PROXY_WS === "true",
  },
};

export default proxy;
