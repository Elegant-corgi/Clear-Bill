// @ts-ignore
/* eslint-disable */
import { request } from "../request";

/** 查询用户列表 超管可查询所有用户，租户管理员仅可查询本租户用户 GET /api/v1/users */
export async function usersList(
  // 叠加生成的Param类型 (非body参数swagger默认没有生成对象)
  params: API.usersListParams,
  options?: { [key: string]: any }
) {
  return request<API.ResponseResult<API.User[]>>("/api/v1/users", {
    method: "GET",
    params: {
      ...params,
    },
    ...(options || {}),
  });
}

/** 创建用户 仅超管和租户管理员可创建用户，默认密码为 stor123; POST /api/v1/users */
export async function usersCreate(
  body: API.CreateUserReq,
  options?: { [key: string]: any }
) {
  return request<API.ResponseResult<API.CreateUserResp>>("/api/v1/users", {
    method: "POST",
    headers: {
      "Content-Type": "application/json",
    },
    data: body,
    ...(options || {}),
  });
}

/** 查询用户详情 查询指定用户详情 GET /api/v1/users/${param0} */
export async function usersGet(
  // 叠加生成的Param类型 (非body参数swagger默认没有生成对象)
  params: API.usersGetParams,
  options?: { [key: string]: any }
) {
  const { id: param0, ...queryParams } = params;
  return request<API.ResponseResult<API.User>>(`/api/v1/users/${param0}`, {
    method: "GET",
    params: { ...queryParams },
    ...(options || {}),
  });
}

/** 更新用户 更新指定用户信息 PUT /api/v1/users/${param0} */
export async function usersUpdate(
  // 叠加生成的Param类型 (非body参数swagger默认没有生成对象)
  params: API.usersUpdateParams,
  body: API.UpdateUserReq,
  options?: { [key: string]: any }
) {
  const { id: param0, ...queryParams } = params;
  return request<API.ResponseResult<API.User>>(`/api/v1/users/${param0}`, {
    method: "PUT",
    headers: {
      "Content-Type": "application/json",
    },
    params: { ...queryParams },
    data: body,
    ...(options || {}),
  });
}

/** 删除用户 删除指定用户 DELETE /api/v1/users/${param0} */
export async function usersDelete(
  // 叠加生成的Param类型 (非body参数swagger默认没有生成对象)
  params: API.usersDeleteParams,
  options?: { [key: string]: any }
) {
  const { id: param0, ...queryParams } = params;
  return request<API.ResponseResult<unknown>>(`/api/v1/users/${param0}`, {
    method: "DELETE",
    params: { ...queryParams },
    ...(options || {}),
  });
}

/** 重置用户密码 超管可修改所有用户密码，租户管理员可修改本租户所有用户密码 PUT /api/v1/users/${param0}/password */
export async function usersPasswordReset(
  // 叠加生成的Param类型 (非body参数swagger默认没有生成对象)
  params: API.usersPasswordResetParams,
  body: API.ResetPasswordReq,
  options?: { [key: string]: any }
) {
  const { id: param0, ...queryParams } = params;
  return request<API.ResponseResult<unknown>>(`/api/v1/users/${param0}/password`, {
    method: "PUT",
    headers: {
      "Content-Type": "application/json",
    },
    params: { ...queryParams },
    data: body,
    ...(options || {}),
  });
}
