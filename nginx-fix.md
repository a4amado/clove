# Nginx Fix — Root Cause & Resolution

## Context

The production setup has three relevant services:

- **`nginx`** — reverse proxy (routes `/` to frontend, `/api/` to backend)
- **`frontend`** — multi-stage Docker build: builds the React app with Node, then serves the `dist/` folder using nginx on port 80
- **`clove`** — Go backend, listens on port 8080

---

## Problems Found

### 1. `nginx` service had no port mapping

The `nginx` reverse proxy container was running but had no `ports` entry in `docker-compose.yml`, so it was completely unreachable from the host machine.

At the same time, the `frontend` service had `ports: - "80:80"`, which exposed the wrong container to the outside world.

```yaml
# BEFORE (broken)
nginx:
  image: nginx:alpine
  # no ports — nothing can reach it from outside

frontend:
  ports:
    - "80:80"   # wrong service exposed
```

**Fix:** Move the port mapping to the `nginx` service and remove it from `frontend`.

---

### 2. `nginx.conf` was mounted into the `frontend` container

The `nginx.conf` is a reverse proxy config (it proxies requests upstream). It was being volume-mounted into the `frontend` container, which already runs its own nginx to serve static files.

This overwrote the frontend's default nginx config, turning it into a broken reverse proxy instead of a static file server. The frontend was no longer serving `dist/index.html` — it was trying to proxy requests upstream and failing.

```yaml
# BEFORE (broken)
frontend:
  volumes:
    - ./nginx.conf:/etc/nginx/conf.d/default.conf:ro  # overwrites static file config
```

**Fix:** Remove the `volumes` mount from the `frontend` service entirely. The frontend container has its own correct nginx config baked into the image via the Dockerfile.

---

### 3. Wrong upstream port in `nginx.conf`

The proxy config pointed to `frontend:5173`, which is the Vite development server port. In production, the frontend container does not run Vite — it runs nginx serving the built static files on port `80`.

So every request to `/` was being proxied to a port that doesn't exist in the prod container, resulting in a 502 Bad Gateway.

```nginx
# BEFORE (broken)
location / {
    proxy_pass http://frontend:5173/;  # Vite dev port, doesn't exist in prod
}
```

```nginx
# AFTER (fixed)
location / {
    proxy_pass http://frontend:80/;  # nginx serving built static files
}
```

---

## Final State

```
Host :80
  └── nginx_proxy (reverse proxy)
        ├── /       → frontend:80  (nginx serving dist/)
        └── /api/   → clove:8080   (Go backend)
```

```yaml
# AFTER (fixed docker-compose.yml, relevant parts)
nginx:
  ports:
    - "80:80"
  depends_on:
    - frontend
    - clove

frontend:
  # no ports, no volumes — uses its own baked-in nginx config
```
