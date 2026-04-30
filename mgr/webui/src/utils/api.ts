import { ApiRequestError } from "@/services/request";

export function unwrapResponse<T = any>(response: API.ResponseResult<any>, fallbackMessage: string) {
  if (!response.success) {
    throw new ApiRequestError(response.error || fallbackMessage, 200, response.errorMessage || response.error || fallbackMessage);
  }

  return response.data as T;
}

export function getErrorMessage(error: unknown, fallbackMessage: string) {
  if (error instanceof ApiRequestError && error.displayMessage) {
    return error.displayMessage;
  }

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

  const year = date.getFullYear();
  const month = date.getMonth() + 1;
  const day = date.getDate();
  const hour = String(date.getHours()).padStart(2, "0");
  const minute = String(date.getMinutes()).padStart(2, "0");
  const second = String(date.getSeconds()).padStart(2, "0");

  return `${year}-${month}-${day} ${hour}:${minute}:${second}`;
}
