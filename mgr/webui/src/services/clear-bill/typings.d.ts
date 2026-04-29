declare namespace API {
  interface ChangeOwnPasswordReq {
    newPassword: string;
    oldPassword: string;
  }

  interface CreateAPITokenReq {
    name: string;
  }

  interface CreateRoleReq {
    code: string;
    name: string;
    permissionIds?: string[];
    scope?: string;
    tenantId?: number;
  }

  interface CreateTenantReq {
    adminDisplayName?: string;
    adminUsername?: string;
    code: string;
    contactName?: string;
    contactPhone?: string;
    name: string;
    remark?: string;
    status?: string;
  }

  interface CreateTenantResp {
    adminUsername: string;
    initialPassword: string;
    tenant: Tenant;
  }

  interface CreateUserReq {
    displayName: string;
    role: string;
    status?: string;
    tenantId?: number;
    username: string;
  }

  interface CreateUserResp {
    initialPassword: string;
    user: User;
  }

  interface LoginReq {
    password: string;
    username: string;
  }

  interface LoginResp {
    user: User;
  }

  interface Permission {
    id: string;
    method: string;
    path: string;
    tag: string;
  }

  interface ResetPasswordReq {
    newPassword: string;
  }

  interface ResponseResult<T = any> {
    data?: T;
    error?: string;
    errorMessage?: string;
    success?: boolean;
  }

  interface Role {
    builtin: boolean;
    code: string;
    createdAt: string;
    id: number;
    name: string;
    permissionIds: string[];
    scope: string;
    tenantId?: number;
    updatedAt: string;
  }

  interface Tenant {
    code: string;
    contactName: string;
    contactPhone: string;
    createdAt: string;
    id: number;
    name: string;
    remark: string;
    status: string;
    updatedAt: string;
  }

  interface UpdateRolePermissionsReq {
    permissionIds?: string[];
  }

  interface UpdateRoleReq {
    name: string;
  }

  interface UpdateTenantReq {
    code: string;
    contactName?: string;
    contactPhone?: string;
    name: string;
    remark?: string;
    status?: string;
  }

  interface UpdateUserReq {
    displayName: string;
    role: string;
    status?: string;
    tenantId?: number;
  }

  interface User {
    createdAt: string;
    displayName: string;
    id: number;
    role: string;
    status: string;
    tenantId?: number;
    updatedAt: string;
    username: string;
  }

  interface rolesDeleteParams {
    id: number;
  }

  interface rolesGetParams {
    id: number;
  }

  interface rolesListParams {
    keyword?: string;
    scope?: string;
    tenantId?: number;
  }

  interface rolesPermissionsUpdateParams {
    id: number;
  }

  interface rolesUpdateParams {
    id: number;
  }

  interface tenantsDeleteParams {
    id: number;
  }

  interface tenantsGetParams {
    id: number;
  }

  interface tenantsListParams {
    keyword?: string;
  }

  interface tenantsUpdateParams {
    id: number;
  }

  interface usersDeleteParams {
    id: number;
  }

  interface usersGetParams {
    id: number;
  }

  interface usersListParams {
    keyword?: string;
    tenantId?: number;
  }

  interface usersPasswordResetParams {
    id: number;
  }

  interface usersUpdateParams {
    id: number;
  }
}
