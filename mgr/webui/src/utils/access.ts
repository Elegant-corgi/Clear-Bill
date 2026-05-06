export const ROLE_SYSADMIN = "sysadmin";
export const ROLE_TENANT_ADMIN = "tenant_admin";
export const ROLE_USER = "user";

export const OVERVIEW_PATH = "/overview";
export const CREDENTIAL_LIST_PATH = "/credentials/list";
export const TENANT_LIST_PATH = "/tenants/list";
export const USER_LIST_PATH = "/tenants/users";
export const ROLE_LIST_PATH = "/tenants/roles";

export function isSysadmin(user?: API.User | null) {
  return user?.role === ROLE_SYSADMIN;
}

export function isTenantAdmin(user?: API.User | null) {
  return user?.role === ROLE_TENANT_ADMIN;
}

export function isBasicUser(user?: API.User | null) {
  return user?.role === ROLE_USER;
}

export function canAccessTenantList(user?: API.User | null) {
  return isSysadmin(user);
}

export function canAccessUserList(user?: API.User | null) {
  return isSysadmin(user) || isTenantAdmin(user);
}

export function canAccessRoleList(user?: API.User | null) {
  return isSysadmin(user) || isTenantAdmin(user);
}

export function canAccessCredentialList(user?: API.User | null) {
  return Boolean(user);
}

export function roleLabel(role?: string) {
  switch (role) {
    case ROLE_SYSADMIN:
      return "超级管理员";
    case ROLE_TENANT_ADMIN:
      return "租户管理员";
    case ROLE_USER:
      return "普通用户";
    default:
      return role || "-";
  }
}

export function resolveDefaultPath() {
  return OVERVIEW_PATH;
}

export function resolveFirstAccessiblePath(user?: API.User | null) {
  return resolveDefaultPath();
}

export function canAccessPath(user: API.User | null | undefined, pathname: string) {
  if (!user) {
    return false;
  }

  if (pathname.startsWith(TENANT_LIST_PATH)) {
    return canAccessTenantList(user);
  }

  if (pathname.startsWith(USER_LIST_PATH)) {
    return canAccessUserList(user);
  }

  if (pathname.startsWith(ROLE_LIST_PATH)) {
    return canAccessRoleList(user);
  }

  if (pathname.startsWith(CREDENTIAL_LIST_PATH)) {
    return canAccessCredentialList(user);
  }

  return true;
}

export function getCurrentTitle(pathname: string) {
  if (pathname.startsWith(TENANT_LIST_PATH)) {
    return "租户列表";
  }

  if (pathname.startsWith(USER_LIST_PATH)) {
    return "用户列表";
  }

  if (pathname.startsWith(ROLE_LIST_PATH)) {
    return "角色权限";
  }

  if (pathname.startsWith(CREDENTIAL_LIST_PATH)) {
    return "凭证管理";
  }

  return "首页";
}
