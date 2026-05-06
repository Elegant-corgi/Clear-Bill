// @ts-ignore
/* eslint-disable */
import { request } from "../request";

export async function credentialsList(
  params: API.credentialsListParams,
  options?: { [key: string]: any },
) {
  return request<any>("/api/v1/credentials", {
    method: "GET",
    params: {
      ...params,
    },
    ...(options || {}),
  });
}

export async function credentialsCreate(body: API.CreateCredentialReq, options?: { [key: string]: any }) {
  return request<any>("/api/v1/credentials", {
    method: "POST",
    headers: {
      "Content-Type": "application/json",
    },
    data: body,
    ...(options || {}),
  });
}

export async function credentialsGet(
  params: API.credentialsGetParams,
  options?: { [key: string]: any },
) {
  const { id: param0, ...queryParams } = params;
  return request<any>(`/api/v1/credentials/${param0}`, {
    method: "GET",
    params: { ...queryParams },
    ...(options || {}),
  });
}

export async function credentialsRotate(
  params: API.credentialsRotateParams,
  body: API.RotateCredentialReq,
  options?: { [key: string]: any },
) {
  const { id: param0, ...queryParams } = params;
  return request<any>(`/api/v1/credentials/${param0}/rotate`, {
    method: "POST",
    headers: {
      "Content-Type": "application/json",
    },
    params: { ...queryParams },
    data: body,
    ...(options || {}),
  });
}

export async function credentialsDelete(
  params: API.credentialsDeleteParams,
  options?: { [key: string]: any },
) {
  const { id: param0, ...queryParams } = params;
  return request<any>(`/api/v1/credentials/${param0}`, {
    method: "DELETE",
    params: { ...queryParams },
    ...(options || {}),
  });
}
