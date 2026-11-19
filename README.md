# Servitor

**Lightweight infrastructure observability and thermal protection for small and medium-sized businesses.**

Servitor is a product concept for teams that operate on-premise or hybrid servers without a dedicated 24/7 infrastructure team. A lightweight local agent is intended to collect temperature, CPU, and memory telemetry, buffer data during network interruptions, and send timely risk signals to a secure SaaS control plane.

The goal is practical: help a small IT team spot overheating and server health issues early, coordinate a response, and reduce avoidable downtime.

> Servitor's footprint and alert-latency figures are engineering targets to validate on supported hardware. The current repository documents the product and an Ubuntu-based lab foundation; it does not yet implement the complete SaaS platform.

## Product promise

- **Low-overhead collection:** target less than 1% average CPU and 50 MB resident memory for the agent, measured and published by supported host class.
- **Fast thermal warning:** target alert dispatch in under five seconds after a qualifying sample when the host, network, and notification path are healthy.
- **Clear operations:** one fleet view for thermal state, CPU, memory, heartbeat, and incident response.
- **Designed for SMEs:** outbound-only host connectivity, straightforward Linux deployment, and actionable alerts rather than a complex observability stack.

## Product and architecture

The [Servitor product and architecture blueprint](docs/SERVITOR_PRODUCT_ARCHITECTURE_BLUEPRINT.md) documents the pitch, capabilities, event-driven design, recommended stack, security model, reliability approach, and v2 roadmap.

The proposed architecture uses a lightweight Linux agent, a TypeScript/Node.js ingestion service behind NGINX, PostgreSQL for durable records, Redis for ephemeral live state and suppression, and a React dashboard with WebSocket updates. Predictive thermal anomaly detection, safe automation, and environment discovery are roadmap items, not current features.

## Ubuntu lab foundation

Ubuntu Server remains the initial host environment for the agent and lab. It is the deployment platform, while **Servitor** is the product.

The [Servitor Linux host preparation guide](Servitor_Linux_Host_Guide.md) contains the basic Ubuntu update, SSH, UFW, and Fail2Ban commands used to prepare a lab host. The earlier academic server design PDF is retained at [Projeto Ubuntu Server.pdf](Projeto%20Ubuntu%20Server.pdf) as historical source material; use the blueprint above as the current product direction.

## Current status

The repository currently contains the original Ubuntu Server lab material and product documentation. Product capabilities such as fleet telemetry ingestion, threshold evaluation, deduplicated Web Push/email alerts, incident workflows, RBAC, and PDF reporting are planned implementation slices. Do not interpret the roadmap as shipped functionality.

## Author

Lucas Vinhal Ferreira — Computer Engineering student at UniCEUB.

[GitHub](https://github.com/romunac) · [LinkedIn](https://www.linkedin.com/in/lucas-vinhal-24334a2b3)
