import type { ApiError } from "../types";

const BASE = "/api/v1";

export class ApiRequestError extends Error {
  status: number;
  body: ApiError;

  constructor(status: number, body: ApiError) {
    super(body.message || body.code);
    this.status = status;
    this.body = body;
  }
}

export async function api<T>(
  path: string,
  options: RequestInit = {},
): Promise<T> {
  const res = await fetch(`${BASE}${path}`, {
    credentials: "include",
    headers: {
      "Content-Type": "application/json",
      ...options.headers,
    },
    ...options,
  });

  if (!res.ok) {
    let body: ApiError;
    try {
      body = await res.json();
    } catch {
      body = {
        request_id: "",
        code: "UNKNOWN",
        message: res.statusText,
        status_code: res.status,
      };
    }
    throw new ApiRequestError(res.status, body);
  }

  if (res.status === 202 || res.status === 204) {
    return undefined as T;
  }

  return res.json();
}

export async function apiRaw(
  path: string,
  options: RequestInit = {},
): Promise<Response> {
  return fetch(`${BASE}${path}`, {
    credentials: "include",
    ...options,
  });
}
