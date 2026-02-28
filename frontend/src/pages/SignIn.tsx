import { useState, type FormEvent } from "react";
import { Link, useNavigate } from "react-router-dom";
import { useAuth } from "../context/AuthContext";

export default function SignIn() {
  const { signIn, loading, error } = useAuth();
  const navigate = useNavigate();
  const [email, setEmail] = useState("");
  const [password, setPassword] = useState("");

  const handleSubmit = async (e: FormEvent) => {
    e.preventDefault();
    try {
      await signIn(email, password);
      navigate("/");
    } catch {
      // error is set in context
    }
  };

  return (
    <div className="min-h-screen flex items-center justify-center">
      <div className="w-full max-w-sm">
        <h1 className="text-2xl font-bold mb-6 text-center">Sign In:</h1>

        {error && (
          <div className="mb-4 p-3 bg-red-950 border border-red-800 rounded-lg text-red-300 text-sm">
            {error}
          </div>
        )}

        <form onSubmit={handleSubmit} className="space-y-4">
          <div>
            <label className="block text-sm text-zinc-400 mb-1">Email</label>
            <input
              type="email"
              required
              value={email}
              onChange={(e) => setEmail(e.target.value)}
              className="w-full px-3 py-2 bg-zinc-900 border border-zinc-700 rounded-lg focus:outline-none focus:border-zinc-500 text-zinc-100"
              placeholder="you@example.com"
            />
          </div>

          <div>
            <label className="block text-sm text-zinc-400 mb-1">Password</label>
            <input
              type="password"
              required
              value={password}
              onChange={(e) => setPassword(e.target.value)}
              className="w-full px-3 py-2 bg-zinc-900 border border-zinc-700 rounded-lg focus:outline-none focus:border-zinc-500 text-zinc-100"
            />
          </div>

          <button
            type="submit"
            disabled={loading}
            className="w-full py-2 bg-zinc-100 text-zinc-900 font-medium rounded-lg hover:bg-zinc-200 transition-colors disabled:opacity-50 cursor-pointer"
          >
            {loading ? "Signing in..." : "Sign In"}
          </button>
        </form>

        <p className="mt-4 text-center text-sm text-zinc-500">
          Don&apos;t have an account?{" "}
          <Link to="/signup" className="text-zinc-300 hover:text-zinc-100">
            Sign Up
          </Link>
        </p>
      </div>
    </div>
  );
}
