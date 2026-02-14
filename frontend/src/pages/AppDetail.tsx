import { useState, useEffect, useCallback, type FormEvent } from "react";
import { useParams, Link } from "react-router-dom";
import * as keysApi from "../api/keys";
import * as regionsApi from "../api/regions";
import * as tokensApi from "../api/tokens";
import type { AppApiKey, Region, OneTimeTokenResponse } from "../types";

export default function AppDetail() {
  const { appId } = useParams<{ appId: string }>();

  if (!appId) return <p>Invalid app ID</p>;

  return (
    <div>
      <Link
        to="/"
        className="text-sm text-zinc-500 hover:text-zinc-300 transition-colors"
      >
        &larr; Back to apps
      </Link>

      <h1 className="text-2xl font-bold mt-4 mb-8">App {appId.slice(0, 8)}...</h1>

      <div className="space-y-10">
        <KeysSection appId={appId} />
        <RegionsSection appId={appId} />
        <TokensSection appId={appId} />
        <MessageSection appId={appId} />
      </div>
    </div>
  );
}

// ─── API Keys ──────────────────────────────────────────

function KeysSection({ appId }: { appId: string }) {
  const [keys, setKeys] = useState<AppApiKey[]>([]);
  const [newKeyName, setNewKeyName] = useState("");
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState<string | null>(null);

  const fetchKeys = useCallback(async () => {
    try {
      const data = await keysApi.listKeys(appId);
      setKeys(data ?? []);
    } catch (e: any) {
      setError(e.message);
    }
  }, [appId]);

  useEffect(() => {
    fetchKeys();
  }, [fetchKeys]);

  const handleCreate = async (e: FormEvent) => {
    e.preventDefault();
    if (!newKeyName.trim()) return;
    setLoading(true);
    setError(null);
    try {
      const key = await keysApi.createKey(appId, newKeyName.trim());
      setKeys((prev) => [...prev, key]);
      setNewKeyName("");
    } catch (e: any) {
      setError(e.message);
    } finally {
      setLoading(false);
    }
  };

  const handleDelete = async (keyId: string) => {
    try {
      await keysApi.deleteKey(appId, keyId);
      setKeys((prev) => prev.filter((k) => k.id !== keyId));
    } catch (e: any) {
      setError(e.message);
    }
  };

  return (
    <section>
      <h2 className="text-lg font-semibold mb-4">API Keys</h2>

      {error && (
        <div className="mb-3 p-3 bg-red-950 border border-red-800 rounded-lg text-red-300 text-sm">
          {error}
        </div>
      )}

      <form onSubmit={handleCreate} className="flex gap-2 mb-4">
        <input
          value={newKeyName}
          onChange={(e) => setNewKeyName(e.target.value)}
          placeholder="Key name"
          className="flex-1 px-3 py-2 bg-zinc-900 border border-zinc-700 rounded-lg focus:outline-none focus:border-zinc-500 text-zinc-100 text-sm"
        />
        <button
          type="submit"
          disabled={loading}
          className="px-4 py-2 bg-zinc-100 text-zinc-900 text-sm font-medium rounded-lg hover:bg-zinc-200 transition-colors disabled:opacity-50 cursor-pointer"
        >
          {loading ? "Creating..." : "Create Key"}
        </button>
      </form>

      {keys.length === 0 ? (
        <p className="text-sm text-zinc-500">No API keys yet.</p>
      ) : (
        <div className="space-y-2">
          {keys.map((key) => (
            <div
              key={key.id}
              className="flex items-center justify-between p-3 bg-zinc-900 border border-zinc-800 rounded-lg"
            >
              <div>
                <span className="font-medium text-sm">{key.key_name}</span>
                <span className="ml-3 text-xs text-zinc-500 font-mono">
                  {key.prefix}...{key.suffix}
                </span>
              </div>
              <button
                onClick={() => handleDelete(key.id)}
                className="text-xs text-red-400 hover:text-red-300 transition-colors cursor-pointer"
              >
                Delete
              </button>
            </div>
          ))}
        </div>
      )}
    </section>
  );
}

// ─── Regions ───────────────────────────────────────────

function RegionsSection({ appId }: { appId: string }) {
  const [regions, setRegions] = useState<Region[]>([]);
  const [editing, setEditing] = useState(false);
  const [editValue, setEditValue] = useState("");
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState<string | null>(null);

  const fetchRegions = useCallback(async () => {
    try {
      const data = await regionsApi.listRegions(appId);
      setRegions(data ?? []);
    } catch (e: any) {
      setError(e.message);
    }
  }, [appId]);

  useEffect(() => {
    fetchRegions();
  }, [fetchRegions]);

  const handleUpdate = async (e: FormEvent) => {
    e.preventDefault();
    setLoading(true);
    setError(null);
    try {
      const newRegions = editValue
        .split(",")
        .map((r) => r.trim())
        .filter(Boolean);
      const updated = await regionsApi.updateRegions(appId, newRegions);
      setRegions(updated ?? []);
      setEditing(false);
    } catch (e: any) {
      setError(e.message);
    } finally {
      setLoading(false);
    }
  };

  return (
    <section>
      <div className="flex items-center justify-between mb-4">
        <h2 className="text-lg font-semibold">Regions</h2>
        <button
          onClick={() => {
            setEditing(!editing);
            setEditValue(regions.join(", "));
          }}
          className="text-sm text-zinc-400 hover:text-zinc-100 transition-colors cursor-pointer"
        >
          {editing ? "Cancel" : "Edit"}
        </button>
      </div>

      {error && (
        <div className="mb-3 p-3 bg-red-950 border border-red-800 rounded-lg text-red-300 text-sm">
          {error}
        </div>
      )}

      {editing ? (
        <form onSubmit={handleUpdate} className="flex gap-2">
          <input
            value={editValue}
            onChange={(e) => setEditValue(e.target.value)}
            className="flex-1 px-3 py-2 bg-zinc-900 border border-zinc-700 rounded-lg focus:outline-none focus:border-zinc-500 text-zinc-100 text-sm"
            placeholder="us-east-1, eu-west-1"
          />
          <button
            type="submit"
            disabled={loading}
            className="px-4 py-2 bg-zinc-100 text-zinc-900 text-sm font-medium rounded-lg hover:bg-zinc-200 transition-colors disabled:opacity-50 cursor-pointer"
          >
            {loading ? "Saving..." : "Save"}
          </button>
        </form>
      ) : regions.length === 0 ? (
        <p className="text-sm text-zinc-500">No regions configured.</p>
      ) : (
        <div className="flex flex-wrap gap-2">
          {regions.map((r) => (
            <span
              key={r}
              className="px-3 py-1 bg-zinc-900 border border-zinc-800 rounded-full text-sm text-zinc-300"
            >
              {r}
            </span>
          ))}
        </div>
      )}
    </section>
  );
}

// ─── One-Time Tokens ───────────────────────────────────

function TokensSection({ appId }: { appId: string }) {
  const [channelId, setChannelId] = useState("");
  const [token, setToken] = useState<OneTimeTokenResponse | null>(null);
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState<string | null>(null);

  const handleCreate = async (e: FormEvent) => {
    e.preventDefault();
    if (!channelId.trim()) return;
    setLoading(true);
    setError(null);
    try {
      const result = await tokensApi.createOneTimeToken(appId, channelId.trim());
      setToken(result);
    } catch (e: any) {
      setError(e.message);
    } finally {
      setLoading(false);
    }
  };

  return (
    <section>
      <h2 className="text-lg font-semibold mb-4">One-Time Token</h2>

      {error && (
        <div className="mb-3 p-3 bg-red-950 border border-red-800 rounded-lg text-red-300 text-sm">
          {error}
        </div>
      )}

      <form onSubmit={handleCreate} className="flex gap-2 mb-4">
        <input
          value={channelId}
          onChange={(e) => setChannelId(e.target.value)}
          placeholder="Channel ID"
          className="flex-1 px-3 py-2 bg-zinc-900 border border-zinc-700 rounded-lg focus:outline-none focus:border-zinc-500 text-zinc-100 text-sm"
        />
        <button
          type="submit"
          disabled={loading}
          className="px-4 py-2 bg-zinc-100 text-zinc-900 text-sm font-medium rounded-lg hover:bg-zinc-200 transition-colors disabled:opacity-50 cursor-pointer"
        >
          {loading ? "Generating..." : "Generate Token"}
        </button>
      </form>

      {token && (
        <div className="p-3 bg-zinc-900 border border-zinc-800 rounded-lg space-y-2">
          <div>
            <span className="text-xs text-zinc-500">Region:</span>
            <span className="ml-2 text-sm">{token.region}</span>
          </div>
          <div>
            <span className="text-xs text-zinc-500">Token:</span>
            <code className="ml-2 text-xs bg-zinc-800 px-2 py-1 rounded break-all">
              {token.token}
            </code>
          </div>
          <p className="text-xs text-zinc-600">Expires in 1 minute.</p>
        </div>
      )}
    </section>
  );
}

// ─── Message Entry ─────────────────────────────────────

function MessageSection({ appId }: { appId: string }) {
  const [channelId, setChannelId] = useState("");
  const [payload, setPayload] = useState("");
  const [sending, setSending] = useState(false);
  const [status, setStatus] = useState<string | null>(null);
  const [error, setError] = useState<string | null>(null);

  const handleSend = async (e: FormEvent) => {
    e.preventDefault();
    if (!channelId.trim() || !payload.trim()) return;
    setSending(true);
    setError(null);
    setStatus(null);

    try {
      const res = await fetch(
        `/api/v1/apps/${appId}/entry/?channel_id=${encodeURIComponent(channelId.trim())}`,
        {
          method: "POST",
          credentials: "include",
          headers: { "Content-Type": "application/octet-stream" },
          body: new TextEncoder().encode(payload),
        },
      );

      if (!res.ok) {
        const body = await res.json().catch(() => null);
        throw new Error(body?.message || res.statusText);
      }

      setStatus("Message sent (202 Accepted)");
      setPayload("");
    } catch (e: any) {
      setError(e.message);
    } finally {
      setSending(false);
    }
  };

  return (
    <section>
      <h2 className="text-lg font-semibold mb-4">Send Message</h2>

      {error && (
        <div className="mb-3 p-3 bg-red-950 border border-red-800 rounded-lg text-red-300 text-sm">
          {error}
        </div>
      )}
      {status && (
        <div className="mb-3 p-3 bg-emerald-950 border border-emerald-800 rounded-lg text-emerald-300 text-sm">
          {status}
        </div>
      )}

      <form onSubmit={handleSend} className="space-y-3">
        <input
          value={channelId}
          onChange={(e) => setChannelId(e.target.value)}
          placeholder="Channel ID"
          className="w-full px-3 py-2 bg-zinc-900 border border-zinc-700 rounded-lg focus:outline-none focus:border-zinc-500 text-zinc-100 text-sm"
        />
        <textarea
          value={payload}
          onChange={(e) => setPayload(e.target.value)}
          placeholder="Message payload (max 32KB)"
          rows={4}
          className="w-full px-3 py-2 bg-zinc-900 border border-zinc-700 rounded-lg focus:outline-none focus:border-zinc-500 text-zinc-100 text-sm resize-none"
        />
        <button
          type="submit"
          disabled={sending}
          className="px-4 py-2 bg-zinc-100 text-zinc-900 text-sm font-medium rounded-lg hover:bg-zinc-200 transition-colors disabled:opacity-50 cursor-pointer"
        >
          {sending ? "Sending..." : "Send Message"}
        </button>
      </form>
    </section>
  );
}
