import { api } from "./client";
import type { OneTimeTokenResponse } from "../types";

export function createOneTimeToken(
  appId: string,
  channelId: string,
): Promise<OneTimeTokenResponse> {
  return api<OneTimeTokenResponse>(`/apps/${appId}/tokens`, {
    method: "POST",
    body: JSON.stringify({ channel_id: channelId }),
  });
}
