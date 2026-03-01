# Clove — Realtime as a Service

> [!WARNING]
> In production, this system requires Kafka, a globally distributed geo-replicated PostgreSQL cluster, Redis clusters, and Kubernetes ingress nodes. For the sake of prototyping, this prototype runs on a single Docker Compose setup.

Clove is a hosted real-time messaging infrastructure. It lets you add WebSocket-based message delivery to your application without building or operating the underlying pub/sub, geo-replication, or connection management yourself.

You create an **App**, your backend publishes messages to a channel, and your clients receive them instantly over a persistent WebSocket connection. Clove handles the rest.

---

## Core Concepts

### App

An App is the top-level unit in Clove. It maps to your product or one logical tenant within it. Each App has:

- A set of **regions** it operates in
- **API Keys** your backend uses to publish messages and issue client tokens
- **Allowed Origins** for WebSocket connections (CORS)
- A **plan tier** (free, standard, pro) that controls resource limits

### Channel

A Channel is a named stream within an App. It has no schema — it's just a key (`notifications`, `room:42`, `user:abc123`) that groups related messages. Clients subscribe to a channel when they open a WebSocket connection. Publishers target a channel when they send a message.

### Region

A Region is a geographic deployment of Clove infrastructure. Each region runs its own message fanout locally to keep delivery latency low. When a message enters one region, Clove replicates it to all other regions so clients everywhere receive it regardless of where the message was published.

Currently active: **dk1** (Denmark).

### Tokens

Clove uses three token types with distinct lifetimes and permissions:

| Token                    | Issued to                  | Lifetime   | Can do                                |
| ------------------------ | -------------------------- | ---------- | ------------------------------------- |
| **Session Token**        | Your team (via login)      | ~30 days   | Manage apps, create API keys          |
| **SDK Token**            | Your backend (via API key) | Long-lived | Publish messages, issue client tokens |
| **One-Time Token (OTT)** | End-user clients           | 1 minute   | Open a single WebSocket connection    |

The OTT model means your clients never hold a long-lived secret. Your backend mints a fresh token per connection; Clove invalidates it immediately on use.

---

## How It Works

### Setup (once)

1. Sign up and create an App in Clove.
2. Create an API Key — this gives you an SDK Token for your backend.
3. Select which regions your App should be active in.

### Per-connection (for each client)

1. Your client asks your backend for a connection token.
2. Your backend calls Clove with its SDK Token → Clove returns a One-Time Token and a target region.
3. Your client opens a WebSocket to Clove using that token. The connection is authenticated and attached to the specified channel.

### Publishing a message

1. Your backend sends a `POST` to the Clove entry endpoint with the channel and a payload (up to 50 KB, any format).
2. Clove acknowledges immediately (202) and publishes the message to the message queue.
3. The message is fanned out to all active regions via replication.
4. Each region delivers the message over Valkey pub/sub to every client subscribed to that channel.
5. Clients receive the payload over their WebSocket connection in real time.

```
Your Backend ──POST /entry──► Clove (dk1)
                                  │
                         ┌────────┼────────────────┐
                         │        │                 │
                    Local fanout  Replicate    Replicate
                         │      to eu1 ...   to us1 ...
                         ▼
               Clients in dk1         Clients in eu1, us1 ...
```

---

## Infrastructure

Clove is built around a small set of components, each with a clear role:

**PostgreSQL** — Source of truth for users, apps, credentials, and API keys. Designed for a globally distributed, geo-sharded deployment.

**Valkey (Redis-compatible)** — Used for two separate concerns:

- _Cache_: App metadata cached close to the API servers to avoid repeated DB reads.
- _Fanout_: Pub/sub delivery of messages to connected WebSocket clients within a region. Reads are local-first for low latency.

**RabbitMQ** — Cross-region message replication. When a message is published, it goes onto a durable queue routed by region. Each region's Clove instances consume from their own queue and fanout locally. Designed to be replaced with Kafka for higher throughput.

**MongoDB** — Prepared for message history and audit trails (not yet active).

**Mailjet** — Transactional email (signup verification).

**Sentry / Loki / Prometheus** — Error tracking, structured logging, and metrics.

---

## Scaling Design

Clove is designed to scale in two directions:

**Horizontally within a region**: Run multiple Clove instances behind a load balancer. RabbitMQ delivers each message to all instances in the region; each instance fans out to its own connected clients. A heartbeat system collects CPU and memory metrics per instance so the load balancer can route new connections to the least-loaded node.

**Across regions**: Deploy Clove to additional regions. RabbitMQ replicates messages globally using per-region routing keys. Each region maintains its own Valkey fanout, so delivery latency is always regional regardless of where the message originated.

The entry point for messages is stateless — any region can accept a publish from your backend and the message will reach clients everywhere.

---

## API Overview

All endpoints are under `/v1`. The API serves OpenAPI documentation at `/docs`.

| Group         | Endpoints                                                                       |
| ------------- | ------------------------------------------------------------------------------- |
| Auth          | `POST /auth/sign-up`, `POST /auth/sign-in`, `POST /auth/verify`, `GET /auth/me` |
| Apps          | `POST /apps`, `GET /apps/:id`                                                   |
| API Keys      | `POST /apps/:id/keys`, `GET /apps/:id/keys`, `DELETE /apps/:id/keys/:key_id`    |
| Regions       | `GET /apps/:id/regions`, `PUT /apps/:id/regions`                                |
| Client Tokens | `POST /apps/:id/tokens`                                                         |
| Message Entry | `POST /apps/:id/entry?channel_id=...`                                           |
| WebSocket     | `GET /apps/:id/ws`                                                              |

---

## Permissions

Clove uses a resource-based permission model. Each token carries a set of `{Resource}:{Operation}` permissions:

- `APP:CREATE` — Create an app
- `KEY:CREATE` — Create an API key
- `ONE_TIME_TOKEN:CREATE` — Issue a client connection token
- `DELIVERY:CREATE` — Publish a message to a channel
- `DELIVERY:READ` — Receive messages over WebSocket

This ensures SDK Tokens can publish but never manage your account, and client OTTs can only receive.

---

## Optimal Infrastructure (target state)

To run at full scale Clove needs:

1. **Global message bus** (Kafka) — Higher throughput and stronger delivery guarantees than RabbitMQ for cross-region replication.
2. **Geo-distributed, sharded PostgreSQL** — Write-local, read-local auth and app data across regions.
3. **Geo-distributed Valkey cluster** — Local-first reads for cache and fanout, sharded by app to avoid hot spots.
