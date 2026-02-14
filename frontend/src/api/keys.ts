import { api } from "./client";
import type { AppApiKey } from "../types";

export function listKeys(
  appId: string,
  pageIdx = 0,
): Promise<AppApiKey[]> {
  return api<AppApiKey[]>(`/apps/${appId}/keys/?page_idx=${pageIdx}`);
}

export function createKey(
  appId: string,
  name: string,
): Promise<AppApiKey> {
  return api<AppApiKey>(`/apps/${appId}/keys/`, {
    method: "POST",
    body: JSON.stringify({ name }),
  });
}

export function deleteKey(
  appId: string,
  keyId: string,
): Promise<void> {
  return api<void>(`/apps/${appId}/keys/${keyId}/`, {
    method: "DELETE",
  });
}
