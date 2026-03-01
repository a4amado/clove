import { useState, type FormEvent } from "react";
import { useNavigate } from "react-router-dom";
import { useAuth } from "../context/AuthContext";
import type { AppWithKeys } from "../types";
import { useV1AppListApps } from "../api/generated.ts";

export default function Dashboard() {
  const { user } = useAuth();
  const navigate = useNavigate();
  const [apps, setApps] = useState<AppWithKeys[]>([]);
  const appsQuery = useV1AppListApps()

  const [showCreate, setShowCreate] = useState(false);
  const [slug, setSlug] = useState("");
  const [regions, setRegions] = useState("");
  const [origins, setOrigins] = useState("");
  const [creating, setCreating] = useState(false);
  const [error, setError] = useState<string | null>(null);

  const handleCreate = async (e: FormEvent) => {
    e.preventDefault();
    if (!user) return;
    setCreating(true);
    setError(null);

    try {
      const result = await appsApi.createApp({
        app_slug: slug,
        regions: regions
          .split(",")
          .map((r) => r.trim())
          .filter(Boolean),
        user_id: user.id,
        allowed_origins: origins
          .split(",")
          .map((o) => o.trim())
          .filter(Boolean),
      });
      setApps((prev) => [...prev, result]);
      setSlug("");
      setRegions("");
      setOrigins("");
      setShowCreate(false);
    } catch (e: any) {
      setError(e.message || "Failed to create app");
    } finally {
      setCreating(false);
    }
  };

  return (
    <div>
      <div className="flex items-center justify-between mb-8">
        <h1 className="text-2xl font-bold">Apps</h1>
        <button
          onClick={() => setShowCreate(!showCreate)}
          className="px-4 py-2 bg-zinc-100 text-zinc-900 text-sm font-medium rounded-lg hover:bg-zinc-200 transition-colors cursor-pointer"
        >
          {showCreate ? "Cancel" : "New App"}
        </button>
      </div>

      {error && (
        <div className="mb-4 p-3 bg-red-950 border border-red-800 rounded-lg text-red-300 text-sm">
          {error}
        </div>
      )}

      {showCreate && (
        <form
          onSubmit={handleCreate}
          className="mb-8 p-4 bg-zinc-900 border border-zinc-800 rounded-lg space-y-4"
        >
          <div>
            <label className="block text-sm text-zinc-400 mb-1">App Slug</label>
            <input
              required
              value={slug}
              onChange={(e) => setSlug(e.target.value)}
              className="w-full px-3 py-2 bg-zinc-950 border border-zinc-700 rounded-lg focus:outline-none focus:border-zinc-500 text-zinc-100"
              placeholder="my-app"
            />
          </div>
          <div>
            <label className="block text-sm text-zinc-400 mb-1">
              Regions (comma-separated)
            </label>
            <input
              required
              value={regions}
              onChange={(e) => setRegions(e.target.value)}
              className="w-full px-3 py-2 bg-zinc-950 border border-zinc-700 rounded-lg focus:outline-none focus:border-zinc-500 text-zinc-100"
              placeholder="us-east-1, eu-west-1"
            />
          </div>
          <div>
            <label className="block text-sm text-zinc-400 mb-1">
              Allowed Origins (comma-separated)
            </label>
            <input
              value={origins}
              onChange={(e) => setOrigins(e.target.value)}
              className="w-full px-3 py-2 bg-zinc-950 border border-zinc-700 rounded-lg focus:outline-none focus:border-zinc-500 text-zinc-100"
              placeholder="https://example.com, http://localhost:3000"
            />
          </div>
          <button
            type="submit"
            disabled={creating}
            className="px-4 py-2 bg-zinc-100 text-zinc-900 text-sm font-medium rounded-lg hover:bg-zinc-200 transition-colors disabled:opacity-50 cursor-pointer"
          >
            {creating ? "Creating..." : "Create App"}
          </button>
        </form>
      )}

      {apps.length === 0 ? (
        <div className="text-center py-16 text-zinc-500">
          <p className="text-lg">No apps yet</p>
          <p className="text-sm mt-1">Create your first app to get started.</p>
        </div>
      ) : (
        <div className="space-y-3">
          {apps.map((a) => (
            <button
              key={a.app.id}
              onClick={() => navigate(`/apps/${a.app.id}`)}
              className="w-full text-left p-4 bg-zinc-900 border border-zinc-800 rounded-lg hover:border-zinc-700 transition-colors cursor-pointer"
            >
              <div className="flex items-center justify-between">
                <div>
                  <h3 className="font-medium">{a.app.app_slug}</h3>
                  <p className="text-sm text-zinc-500 mt-1">
                    {a.app.regions?.join(", ") || "No regions"} &middot;{" "}
                    {a.keys?.length || 0} key(s)
                  </p>
                </div>
                <span className="text-xs px-2 py-1 bg-zinc-800 rounded text-zinc-400">
                  {a.app.app_type}
                </span>
              </div>
            </button>
          ))}
        </div>
      )}
    </div>
  );
}
