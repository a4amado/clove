import { api } from "./client";
import type { Region } from "../types";

export function listRegions(appId: string): Promise<Region[]> {
  return api<Region[]>(`/apps/${appId}/regions/`);
}

export function updateRegions(
  appId: string,
  regions: Region[],
): Promise<Region[]> {
  return api<Region[]>(`/apps/${appId}/regions/`, {
    method: "PATCH",
    body: JSON.stringify({ regions }),
  });
}
