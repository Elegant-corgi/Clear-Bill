// @ts-ignore
/* eslint-disable */
import { request } from "../request";

/** 用户登录 使用用户名和密码登录，设置后台 session POST /api/v1/auth/login */
export async function authLogin(
  body: API.LoginReq,
  options?: { [key: string]: any }
) {
  return request<API.ResponseResult<API.LoginResp>>("/api/v1/auth/login", {
    method: "POST",
    headers: {
      "Content-Type": "application/json",
    },
    data: body,
    ...(options || {}),
  });
}

/** 用户登出 清除当前后台 session POST /api/v1/auth/logout */
export async function authLogout(options?: { [key: string]: any }) {
  return request<API.ResponseResult<unknown>>("/api/v1/auth/logout", {
    method: "POST",
    ...(options || {}),
  });
}

/** 当前用户信息 返回当前登录用户信息 GET /api/v1/auth/me */
export async function authMe(options?: { [key: string]: any }) {
  return request<API.ResponseResult<API.User>>("/api/v1/auth/me", {
    method: "GET",
    ...(options || {}),
  });
}

/** 修改自己的密码 当前登录用户修改自己的密码 PUT /api/v1/auth/password */
export async function authPasswordChange(
  body: API.ChangeOwnPasswordReq,
  options?: { [key: string]: any }
) {
  return request<API.ResponseResult<unknown>>("/api/v1/auth/password", {
    method: "PUT",
    headers: {
      "Content-Type": "application/json",
    },
    data: body,
    ...(options || {}),
  });
}

/** 创建 API token 为当前已登录用户签发一个用于外部 API 调用的 token POST /api/v1/auth/tokens */
export async function authTokenCreate(
  body: API.CreateAPITokenReq,
  options?: { [key: string]: any }
) {
  return request<API.ResponseResult<{ name: string; token: string }>>("/api/v1/auth/tokens", {
    method: "POST",
    headers: {
      "Content-Type": "application/json",
    },
    data: body,
    ...(options || {}),
  });
}
