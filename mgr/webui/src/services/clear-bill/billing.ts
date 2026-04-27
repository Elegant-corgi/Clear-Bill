// @ts-ignore
/* eslint-disable */
import { request } from "../request";

/** 查询账单列表 返回账单列表数据 GET /api/v1/bills */
export async function billsList(options?: { [key: string]: any }) {
  return request<any>("/api/v1/bills", {
    method: "GET",
    ...(options || {}),
  });
}

/** 查询客户列表 返回客户列表数据 GET /api/v1/customers */
export async function customersList(options?: { [key: string]: any }) {
  return request<any>("/api/v1/customers", {
    method: "GET",
    ...(options || {}),
  });
}

/** 查询仪表盘概览 返回账单工作台概览数据 GET /api/v1/dashboard/summary */
export async function dashboardSummaryGet(options?: { [key: string]: any }) {
  return request<any>("/api/v1/dashboard/summary", {
    method: "GET",
    ...(options || {}),
  });
}

/** 查询对账任务列表 返回对账任务列表数据 GET /api/v1/reconciliations */
export async function reconciliationsList(options?: { [key: string]: any }) {
  return request<any>("/api/v1/reconciliations", {
    method: "GET",
    ...(options || {}),
  });
}
