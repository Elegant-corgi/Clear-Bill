// @ts-ignore
/* eslint-disable */
import { request } from "../request";

/** 查询租户列表 根据关键字查询租户列表 GET /api/v1/tenants */
export async function tenantsList(
  // 叠加生成的Param类型 (非body参数swagger默认没有生成对象)
  params: API.tenantsListParams,
  options?: { [key: string]: any }
) {
  return request<API.ResponseResult<API.Tenant[]>>("/api/v1/tenants", {
    method: "GET",
    params: {
      ...params,
    },
    ...(options || {}),
  });
}

/** 创建租户 创建一条新的租户记录 POST /api/v1/tenants */
export async function tenantsCreate(
  body: API.CreateTenantReq,
  options?: { [key: string]: any }
) {
  return request<API.ResponseResult<API.CreateTenantResp>>("/api/v1/tenants", {
    method: "POST",
    headers: {
      "Content-Type": "application/json",
    },
    data: body,
    ...(options || {}),
  });
}

/** 查询租户详情 根据租户ID查询租户详情 GET /api/v1/tenants/${param0} */
export async function tenantsGet(
  // 叠加生成的Param类型 (非body参数swagger默认没有生成对象)
  params: API.tenantsGetParams,
  options?: { [key: string]: any }
) {
  const { id: param0, ...queryParams } = params;
  return request<API.ResponseResult<API.Tenant>>(`/api/v1/tenants/${param0}`, {
    method: "GET",
    params: { ...queryParams },
    ...(options || {}),
  });
}

/** 更新租户 根据租户ID更新租户信息 PUT /api/v1/tenants/${param0} */
export async function tenantsUpdate(
  // 叠加生成的Param类型 (非body参数swagger默认没有生成对象)
  params: API.tenantsUpdateParams,
  body: API.UpdateTenantReq,
  options?: { [key: string]: any }
) {
  const { id: param0, ...queryParams } = params;
  return request<API.ResponseResult<API.Tenant>>(`/api/v1/tenants/${param0}`, {
    method: "PUT",
    headers: {
      "Content-Type": "application/json",
    },
    params: { ...queryParams },
    data: body,
    ...(options || {}),
  });
}

/** 删除租户 根据租户ID删除租户 DELETE /api/v1/tenants/${param0} */
export async function tenantsDelete(
  // 叠加生成的Param类型 (非body参数swagger默认没有生成对象)
  params: API.tenantsDeleteParams,
  options?: { [key: string]: any }
) {
  const { id: param0, ...queryParams } = params;
  return request<API.ResponseResult<unknown>>(`/api/v1/tenants/${param0}`, {
    method: "DELETE",
    params: { ...queryParams },
    ...(options || {}),
  });
}
