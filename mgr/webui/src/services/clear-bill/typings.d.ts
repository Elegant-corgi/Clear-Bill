declare namespace API {
  interface PageParams {
    page?: number;
    pageSize?: number;
  }

  interface PageResult<T = any> {
    list: T[];
    page: number;
    pageSize: number;
    total: number;
  }

  interface ResponseResult<T = any> {
    data?: T;
    error?: string;
    errorMessage?: string;
    success?: boolean;
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

  interface Tenant {
    adminDisplayName: string;
    adminUsername: string;
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

  interface Permission {
    id: string;
    method: string;
    path: string;
    tag: string;
  }

  interface LoginReq {
    password: string;
    username: string;
  }

  interface LoginResp {
    user: User;
  }

  interface ChangeOwnPasswordReq {
    newPassword: string;
    oldPassword: string;
  }

  interface CreateAPITokenReq {
    name: string;
  }

  interface Credential {
    accessKey?: string;
    accessKeyPreview?: string;
    createdAt: string;
    expiresAt?: string;
    id: number;
    lastUsedAt?: string;
    name: string;
    parentId?: number;
    rotatedAt?: string;
    secretKey?: string;
    status: string;
    token?: string;
    tokenPreview?: string;
    type: string;
    updatedAt: string;
    userId: number;
  }

  interface CreateCredentialReq {
    name: string;
    type: string;
  }

  interface RotateCredentialReq {
    name?: string;
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

  interface ResetPasswordReq {
    newPassword: string;
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

  interface rolesDeleteParams {
    id: number;
  }

  interface rolesGetParams {
    id: number;
  }

  interface rolesListParams extends PageParams {
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

  interface tenantsListParams extends PageParams {
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

  interface usersListParams extends PageParams {
    keyword?: string;
    tenantId?: number;
  }

  interface usersPasswordResetParams {
    id: number;
  }

  interface usersUpdateParams {
    id: number;
  }

  interface credentialsDeleteParams {
    id: number;
  }

  interface credentialsGetParams {
    id: number;
  }

  interface credentialsListParams extends PageParams {
    status?: string;
    type?: string;
  }

  interface credentialsRotateParams {
    id: number;
  }
}
