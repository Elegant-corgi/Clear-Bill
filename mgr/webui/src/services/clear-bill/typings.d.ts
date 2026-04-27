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

  interface CreateUserReq {
    displayName: string;
    role: string;
    status?: string;
    tenantId?: number;
    username: string;
  }

  interface LoginReq {
    password: string;
    username: string;
  }

  interface ResetPasswordReq {
    newPassword: string;
  }

  interface ResponseResult {
    data?: any;
    error?: string;
    success?: boolean;
  }

  interface rolesDeleteParams {
    /** 角色ID */
    id: number;
  }

  interface rolesGetParams {
    /** 角色ID */
    id: number;
  }

  interface rolesListParams {
    /** 关键字 */
    keyword?: string;
    /** 角色作用域 */
    scope?: string;
    /** 租户ID */
    tenantId?: number;
  }

  interface rolesPermissionsUpdateParams {
    /** 角色ID */
    id: number;
  }

  interface rolesUpdateParams {
    /** 角色ID */
    id: number;
  }

  interface tenantsDeleteParams {
    /** 租户ID */
    id: number;
  }

  interface tenantsGetParams {
    /** 租户ID */
    id: number;
  }

  interface tenantsListParams {
    /** 关键字 */
    keyword?: string;
  }

  interface tenantsUpdateParams {
    /** 租户ID */
    id: number;
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

  interface usersDeleteParams {
    /** 用户ID */
    id: number;
  }

  interface usersGetParams {
    /** 用户ID */
    id: number;
  }

  interface usersListParams {
    /** 关键字 */
    keyword?: string;
    /** 租户ID */
    tenantId?: number;
  }

  interface usersPasswordResetParams {
    /** 用户ID */
    id: number;
  }

  interface usersUpdateParams {
    /** 用户ID */
    id: number;
  }
}
