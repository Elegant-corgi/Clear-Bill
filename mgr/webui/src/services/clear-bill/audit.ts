// @ts-ignore
/* eslint-disable */
import { request } from "../request";

/** 查询审计日志 支持按用户、操作类型和时间段筛选审计日志 GET /api/v1/audit-logs */
export async function auditLogsList(
  params: API.auditLogsListParams,
  options?: { [key: string]: any },
) {
  return request<any>("/api/v1/audit-logs", {
    method: "GET",
    params: {
      ...params,
    },
    ...(options || {}),
  });
}
