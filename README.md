# Servitor

**Lightweight infrastructure observability and thermal protection for small and medium-sized businesses.**

Servitor is designed for teams that depend on on-premise and hybrid servers but do not have a dedicated 24/7 infrastructure team. Its local agent collects thermal and resource signals, buffers data during network interruptions, and sends telemetry to a focused control plane so teams can spot risk and respond before a server outage.

> **Status:** this repository now contains an MVP engineering foundation for Linux telemetry collection and API ingestion. The product goals below are targets; the full SaaS experience is not yet implemented.

## Product goals

- **Low-overhead collection:** target under 1% average CPU and 50 MB resident memory, measured across supported host classes.
- **Fast thermal warning:** target alert dispatch under five seconds after a qualifying sample when host connectivity and the delivery path are healthy.
- **Actionable operations:** bring temperature, CPU, RAM, heartbeat, and incident response into one fleet view.
- **SME-friendly deployment:** use a lightweight Linux agent with outbound telemetry and an understandable operational workflow.

These figures are engineering goals, not validated service guarantees. Sensor availability and safe thermal thresholds vary by hardware.

## MVP foundation in this repository

| Component | Current foundation |
|---|---|
| Linux agent | Go prototype sampling CPU and memory from `/proc`, temperature from supported `hwmon` and thermal-zone files, and keeping a bounded offline spool. |
| Ingestion API | Node.js, TypeScript, and Express endpoints for health/readiness and authenticated telemetry batches. Payloads are validated and batch IDs make retries idempotent. |
| Persistence | PostgreSQL migration for tenant, host, sample, batch, incident, and audit entities. Telemetry batches and samples are persisted transactionally. |
| Runtime dependencies | Docker Compose configuration for local PostgreSQL and Redis. Redis is currently used for readiness checks; queue-based alert processing is future work. |

### Implemented versus planned

**Present in the starter code:** a Linux collection prototype; a development shared-token ingestion endpoint; payload validation; PostgreSQL persistence; batch deduplication; and local database/cache configuration.

**Still to build:** per-agent enrollment, credentials and revocation; tenant/user authentication and RBAC; threshold evaluation and alert suppression; email/Web Push delivery; dashboard and WebSocket streaming; incident acknowledgement workflows; PDF reporting; and production TLS, retention, backups, monitoring, and service objectives.

The `<1% CPU`, `<50 MB RAM`, and `<5 second alert` figures must be benchmarked before they are presented as measured performance. The current starter does not send alerts and must not be represented as providing automatic hardware protection.

## Quick start (development)

Requirements: Node.js 20+, npm, Go 1.22+, and Docker Compose v2.

1. Start the local PostgreSQL and Redis dependencies from the repository root:

   ```sh
   docker compose -f deploy/compose.yaml up -d
   ```

2. Copy `services/api/.env.example` to `services/api/.env`. Set a long random development `AGENT_TOKEN`. Keep this file private; it is ignored by Git.

3. Install and start the API:

   ```sh
   cd services/api
   npm install
   npm run migrate
   npm run dev
   ```

   The API listens on port 8080. `GET /health` checks process health; `GET /ready` checks PostgreSQL and Redis; `POST /v1/telemetry` accepts authenticated batches.

4. In another terminal, build and run the Linux agent from the repository root:

   ```sh
   cd agent
   go build -o servitor-agent ./cmd/servitor-agent
   ./servitor-agent --api http://localhost:8080 --token <development-token> --host-id <host-uuid>
   ```

Use the same token as `AGENT_TOKEN`. Replace `<host-uuid>` with a UUID. HTTP on localhost is for local development only. Put the API behind TLS before any non-local use.

## Repository map

- [Product and architecture blueprint](docs/SERVITOR_PRODUCT_ARCHITECTURE_BLUEPRINT.md)
- [Linux host preparation guide](Servitor_Linux_Host_Guide.md)
- [Agent entry point](agent/cmd/servitor-agent/main.go) and [Linux collectors](agent/internal/collect/collect_linux.go)
- [API service](services/api/src/server.ts), [database migration](services/api/src/migrations/001_initial.sql), and [local environment template](services/api/.env.example)
- [Local PostgreSQL and Redis services](deploy/compose.yaml)
- The original [Ubuntu Server project PDF](Projeto%20Ubuntu%20Server.pdf) is retained as historical source material; the Servitor blueprint is the current product direction.

## Validation status

The package and TypeScript configuration files were checked, but the starter was not built or run in the authoring environment: Go and Docker were unavailable and npm dependencies were not installed. Run the quick start in a prepared Linux development environment before treating the code as operational.

## Security boundary

The current API uses one shared `AGENT_TOKEN` and assigns accepted hosts to a fixed local-development tenant. That is a development shortcut, not a production identity or tenancy model. Before customer use, implement per-agent credentials and revocation, validate host-to-tenant binding, add user authentication and authorization, enforce tenant isolation, configure TLS, and establish retention and backup controls. Never commit real credentials or expose the development API directly to the internet.

## Roadmap

The v2 direction includes explainable predictive thermal anomaly detection, customer-approved remediation hooks, and read-only discovery of container and cloud context. These are future capabilities, not part of the current MVP code.
