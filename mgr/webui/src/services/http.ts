import appConfig from "@config/config";

interface ApiEnvelope<T> {
  success: boolean;
  data: T;
  error?: string;
}

export async function requestWithFallback<T>(path: string, fallback: T): Promise<T> {
  try {
    const response = await fetch(`${appConfig.apiPrefix}${path}`, {
      headers: {
        Accept: "application/json",
      },
    });

    if (!response.ok) {
      throw new Error(`Request failed: ${response.status}`);
    }

    const payload = (await response.json()) as ApiEnvelope<T>;

    if (!payload.success) {
      throw new Error(payload.error ?? "Unknown API error");
    }

    return payload.data;
  } catch {
    return fallback;
  }
}
