export function unwrapResponse<T>(response: API.ResponseResult<T>, fallbackMessage: string) {
  if (!response.success) {
    throw new Error(response.error || fallbackMessage);
  }

  return response.data as T;
}

export function getErrorMessage(error: unknown, fallbackMessage: string) {
  if (error instanceof Error && error.message && !error.message.startsWith("Request failed:")) {
    return error.message;
  }

  return fallbackMessage;
}

export function formatDateTime(value?: string) {
  if (!value) {
    return "-";
  }

  const date = new Date(value);
  if (Number.isNaN(date.getTime())) {
    return value;
  }

  return new Intl.DateTimeFormat("zh-CN", {
    year: "numeric",
    month: "2-digit",
    day: "2-digit",
    hour: "2-digit",
    minute: "2-digit",
    second: "2-digit",
    hour12: false,
  }).format(date);
}
