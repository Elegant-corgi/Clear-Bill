// @ts-ignore
/* eslint-disable */
import { request } from "../request";

/** 查询权限列表 返回基于 Swagger 接口注解生成的权限列表 GET /api/v1/permissions */
export async function permissionsList(options?: { [key: string]: any }) {
  return request<any>("/api/v1/permissions", {
    method: "GET",
    ...(options || {}),
  });
}

/** 查询角色列表 查询可见范围内的角色列表 GET /api/v1/roles */
export async function rolesList(
  // 叠加生成的Param类型 (非body参数swagger默认没有生成对象)
  params: API.rolesListParams,
  options?: { [key: string]: any }
) {
  return request<any>("/api/v1/roles", {
    method: "GET",
    params: {
      ...params,
    },
    ...(options || {}),
  });
}

/** 创建角色 创建角色并分配权限 POST /api/v1/roles */
export async function rolesCreate(
  body: API.CreateRoleReq,
  options?: { [key: string]: any }
) {
  return request<any>("/api/v1/roles", {
    method: "POST",
    headers: {
      "Content-Type": "application/json",
    },
    data: body,
    ...(options || {}),
  });
}

/** 查询角色详情 查询指定角色信息 GET /api/v1/roles/${param0} */
export async function rolesGet(
  // 叠加生成的Param类型 (非body参数swagger默认没有生成对象)
  params: API.rolesGetParams,
  options?: { [key: string]: any }
) {
  const { id: param0, ...queryParams } = params;
  return request<any>(`/api/v1/roles/${param0}`, {
    method: "GET",
    params: { ...queryParams },
    ...(options || {}),
  });
}

/** 更新角色 更新指定角色的基础信息 PUT /api/v1/roles/${param0} */
export async function rolesUpdate(
  // 叠加生成的Param类型 (非body参数swagger默认没有生成对象)
  params: API.rolesUpdateParams,
  body: API.UpdateRoleReq,
  options?: { [key: string]: any }
) {
  const { id: param0, ...queryParams } = params;
  return request<any>(`/api/v1/roles/${param0}`, {
    method: "PUT",
    headers: {
      "Content-Type": "application/json",
    },
    params: { ...queryParams },
    data: body,
    ...(options || {}),
  });
}

/** 删除角色 删除指定角色 DELETE /api/v1/roles/${param0} */
export async function rolesDelete(
  // 叠加生成的Param类型 (非body参数swagger默认没有生成对象)
  params: API.rolesDeleteParams,
  options?: { [key: string]: any }
) {
  const { id: param0, ...queryParams } = params;
  return request<any>(`/api/v1/roles/${param0}`, {
    method: "DELETE",
    params: { ...queryParams },
    ...(options || {}),
  });
}

/** 更新角色权限 覆盖指定角色的权限列表 PUT /api/v1/roles/${param0}/permissions */
export async function rolesPermissionsUpdate(
  // 叠加生成的Param类型 (非body参数swagger默认没有生成对象)
  params: API.rolesPermissionsUpdateParams,
  body: API.UpdateRolePermissionsReq,
  options?: { [key: string]: any }
) {
  const { id: param0, ...queryParams } = params;
  return request<any>(`/api/v1/roles/${param0}/permissions`, {
    method: "PUT",
    headers: {
      "Content-Type": "application/json",
    },
    params: { ...queryParams },
    data: body,
    ...(options || {}),
  });
}
