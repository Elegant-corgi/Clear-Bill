import appConfig from '@config/config';

export interface RequestOptions extends Omit<RequestInit, 'body'> {
  params?: Record<string, any>;
  data?: any;
  requestType?: 'form';
}

export class ApiRequestError extends Error {
  displayMessage?: string;
  status: number;

  constructor(message: string, status: number, displayMessage?: string) {
    super(message);
    this.name = 'ApiRequestError';
    this.status = status;
    this.displayMessage = displayMessage;
  }
}

function buildURL(path: string, params?: Record<string, any>) {
  const url = new URL(`${appConfig.apiOrigin}${path}`, 'http://placeholder.local');
  if (params) {
    Object.entries(params).forEach(([key, value]) => {
      if (value === undefined || value === null) {
        return;
      }
      if (Array.isArray(value)) {
        value.forEach((item) => url.searchParams.append(key, String(item)));
        return;
      }
      url.searchParams.set(key, String(value));
    });
  }
  return `${url.pathname}${url.search}`;
}

function parseResponsePayload(text: string) {
  if (!text) {
    return null;
  }

  try {
    return JSON.parse(text) as unknown;
  } catch {
    return null;
  }
}

export async function request<T>(path: string, options: RequestOptions = {}): Promise<T> {
  const { params, data, headers, requestType, ...rest } = options;

  let body: BodyInit | undefined;
  const nextHeaders = new Headers(headers || {});

  if (data !== undefined) {
    if (requestType === 'form' && data instanceof FormData) {
      body = data;
    } else {
      nextHeaders.set('Content-Type', 'application/json');
      body = JSON.stringify(data);
    }
  }

  const response = await fetch(buildURL(path, params), {
    credentials: 'include',
    ...rest,
    headers: nextHeaders,
    body,
  });
  const text = await response.text();
  const payload = parseResponsePayload(text);

  if (!response.ok) {
    const rawError =
      payload && typeof payload === 'object' && typeof payload.error === 'string'
        ? payload.error
        : `Request failed: ${response.status}`;
    const displayMessage =
      payload && typeof payload === 'object' && typeof payload.errorMessage === 'string'
        ? payload.errorMessage
        : undefined;
    throw new ApiRequestError(rawError, response.status, displayMessage);
  }

  return payload as T;
}
