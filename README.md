# RooomMetrics — Cluster-Aware Metrics Engine

Version: `0.3.0`

RooomMetrics is a source-available metrics engine focused on multi-tenancy, cardinality protection, peer replication, distributed reads, and online anomaly detection. It is an independent implementation and is not a fork of Prometheus.

## Core features

- Per-tenant active-series budgets
- Multi-tenant isolation with `X-Scope-OrgID`
- Peer write replication and distributed query fan-out
- Online anomaly stream and adaptive downsampling
- Prometheus Remote Write v1 receive
- Prometheus text exposition import and scraping
- OTLP/HTTP JSON metrics receive
- Partial PromQL-compatible queries
- WAL-backed local persistence

## Run

```bash
go test ./...
go run ./cmd/rooommetrics
```

Or run `docker compose up --build`. The server listens on `:9090` by default.

Cluster example:

```bash
ROOOM_PEERS=http://metrics-2:9090,http://metrics-3:9090 \
ROOOM_REPLICATION_FACTOR=2 \
go run ./cmd/rooommetrics
```

The current cluster implementation is an MVP and is not yet a consensus/quorum system.

## Licensing and support

Version `0.3.0` and later use the ROOOMTECH licensing terms in `LICENSE`. A separate commercial software license agreement and paid maintenance, support, implementation, integration, upgrades, security support, SLA options, private builds, and custom development are available.

Contact: `support@rooomtech.com`

Earlier releases retain their published license terms. Third-party software retains its own licenses.
