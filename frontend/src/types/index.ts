export interface User {
  id: string;
  email: string;
  email_verifycode?: string;
  created_at?: string;
  updated_at?: string;
}

export interface App {
  id: string;
  app_slug: string;
  app_type: string;
  user_id: string;
  regions: Region[];
  allowed_origins: string[];
  created_at?: string;
  updated_at?: string;
}

export interface AppApiKey {
  id: string;
  app_id: string;
  key_name: string;
  prefix: string;
  suffix: string;
  created_at?: string;
  updated_at?: string;
}

export interface AppWithKeys {
  app: App;
  keys: AppApiKey[];
}

export type Region = string;

export interface OneTimeTokenResponse {
  token: string;
  region: string;
}

export interface ApiError {
  request_id: string;
  code: string;
  message: string;
  status_code: number;
}

export interface WebSocketMessage {
  channel: string;
  payload: string;
}
