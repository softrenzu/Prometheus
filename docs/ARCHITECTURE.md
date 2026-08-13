# RooomMetrics architecture

RooomMetrics is an experimental metrics engine designed around four properties that are deliberately built into the server rather than delegated to sidecars: multi-tenancy, cardinality protection, peer replication/distributed reads, and online anomaly detection.

## Data path

1. Ingest from Prometheus Remote Write v1, Prometheus text exposition, OTLP/HTTP JSON, or static scraping.
2. Normalize labels and compute a deterministic series key.
3. Enforce a per-tenant active-series budget before admitting a new series.
4. Append to one of N in-memory shards and a crash-recovery WAL.
5. Update online statistics for anomaly detection.
6. Optionally replicate accepted samples to configured peers.

## Query path

The HTTP API follows the Prometheus `/api/v1/query` and `/api/v1/query_range` response shape for the supported query subset. Reads are fanned out to configured peers and merged. Range queries automatically increase effective resolution when the result would exceed 11,000 points per series.

Supported query syntax in this MVP:

- `metric_name`
- `metric_name{label="value"}`
- `sum(...)`, `avg(...)`, `min(...)`, `max(...)`, `count(...)`
- `rate(...)`

This is not a full PromQL implementation yet.

## Failure model

The WAL protects local accepted samples across process restart. Peer replication is asynchronous and best-effort in this MVP; it is not a consensus protocol. A production-grade release should replace static peers with membership, consistent hashing, quorum writes, segment replication, checksums, compaction, and object-storage tiering.

## Security model

Tenant isolation is selected through `X-Scope-OrgID`. This header must be set by a trusted reverse proxy in production. Authentication and authorization are intentionally outside this MVP and should be enforced at the ingress or service mesh until native auth is added.
