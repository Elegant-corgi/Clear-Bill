import appConfig from '@config/config';

export interface RequestOptions extends Omit<RequestInit, 'body'> {
  params?: Record<string, any>;
  data?: any;
  requestType?: 'form';
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

  if (!response.ok) {
    throw new Error(`Request failed: ${response.status}`);
  }

  return (await response.json()) as T;
}
