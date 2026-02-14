import { api } from "./client";
import type { App, AppWithKeys, Region } from "../types";

export function createApp(params: {
  app_slug: string;
  regions: Region[];
  user_id: string;
  allowed_origins: string[];
}): Promise<AppWithKeys> {
  return api<AppWithKeys>("/apps/", {
    method: "POST",
    body: JSON.stringify(params),
  });
}

export function getApp(appId: string): Promise<App> {
  return api<App>(`/apps/${appId}/`);
}
