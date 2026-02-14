import { api, ApiRequestError } from "./client";
import type { User } from "../types";

export function signUp(
  email: string,
  password: string,
  confirm_password: string,
): Promise<User> {
  return api<User>("/auth/sign-up", {
    method: "POST",
    body: JSON.stringify({ email, password, confirm_password }),
  });
}

export function signIn(
  email: string,
  password: string,
): Promise<User> {
  return api<User>("/auth/sign-in", {
    method: "POST",
    body: JSON.stringify({ email, password }),
  });
}

export async function me(): Promise<User | null> {
  try {
    return await api<User>("/auth/me");
  } catch (e) {
    if (e instanceof ApiRequestError && e.status === 401) {
      return null;
    }
    return null;
  }
}
