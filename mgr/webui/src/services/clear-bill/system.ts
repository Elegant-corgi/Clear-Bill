// @ts-ignore
/* eslint-disable */
import { request } from "../request";

/** 健康检查 返回服务健康状态与版本信息 GET /api/v1/health */
export async function healthGet(options?: { [key: string]: any }) {
  return request<any>("/api/v1/health", {
    method: "GET",
    ...(options || {}),
  });
}
