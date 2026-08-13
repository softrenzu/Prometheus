# RooomMetrics

Version: `0.2.0`

**RooomMetrics** is an experimental open-source metrics engine aimed at the gaps that normally require Prometheus plus extra systems: multi-tenancy, cardinality protection, peer replication, distributed reads, and online anomaly detection.

It is not a fork of Prometheus. It keeps familiar ingestion/query surfaces where practical while using an independent implementation.

## Why this exists

Prometheus is excellent at reliable single-server monitoring, but its built-in local TSDB is not clustered or replicated and distributed PromQL evaluation is intentionally not built into the core server. RooomMetrics starts from a different default: cluster-aware operation and tenant safety are first-class.

### What this MVP adds beyond a stock Prometheus server

| Capability | RooomMetrics MVP | Stock Prometheus core |
|---|---|---|
| Per-tenant active-series budget | Built in | External policy/architecture usually needed |
| Multi-tenant header isolation | Built in | Not a native multi-tenant TSDB |
| Peer write replication | Built in | HA uses multiple servers/external deduplication |
| Distributed query fan-out | Built in | Remote read evaluation remains local |
| Online anomaly stream | Built in | Rules/queries or external system |
| Adaptive range downsampling | Built in | Recording rules/external long-term systems commonly used |
| Prometheus Remote Write v1 receive | Yes | Yes |
| Prometheus text exposition import/scrape | Yes | Yes |
| OTLP/HTTP JSON metrics receive | Yes | Yes in modern Prometheus |
| Full PromQL | Partial | Yes |
| Native histograms | Roadmap | Stable in modern Prometheus |
| Production maturity | Experimental | Production proven |

The goal is not to claim that an early MVP is globally better than Prometheus. The goal is to make the architectural advantages concrete and executable, then close the compatibility and durability gaps.

## Run

```bash
go test ./...
go run ./cmd/rooommetrics
```

Or:

```bash
docker compose up --build
```

The server listens on `:9090` by default.

## Ingest

### Prometheus text

```bash
curl -X POST http://localhost:9090/api/v1/import/prometheus \
  -H 'X-Scope-OrgID: demo' \
  --data-binary 'cpu_usage{host="a"} 42.5'
```

### Prometheus Remote Write

Configure Prometheus remote write to `http://rooommetrics:9090/api/v1/write`. The MVP accepts the Prometheus Remote Write v1 protobuf message compressed with Snappy and currently stores float samples and labels.

### OTLP/HTTP JSON

Send OpenTelemetry metrics JSON to:

```text
POST /v1/metrics
```

Gauge, Sum, and Histogram count/sum points are accepted in this MVP.

### Pull scraping

```bash
ROOOM_SCRAPE_TARGETS=http://node-exporter:9100/metrics,http://app:8080/metrics \
ROOOM_SCRAPE_INTERVAL=15s \
go run ./cmd/rooommetrics
```

## Query

```bash
curl -G http://localhost:9090/api/v1/query \
  -H 'X-Scope-OrgID: demo' \
  --data-urlencode 'query=avg(cpu_usage)'
```

Supported query subset:

- selectors: `metric{label="value"}`
- aggregations: `sum`, `avg`, `min`, `max`, `count`
- `rate`
- instant and range APIs

## Cardinality guard

Every tenant has an active-series ceiling. When a new label combination would exceed it, the write is rejected before memory pressure spreads across the process.

```bash
ROOOM_CARDINALITY_LIMIT=1000000
```

Inspect it with:

```text
GET /api/v1/status/cardinality
```

## Built-in anomalies

RooomMetrics maintains streaming statistics per series and records large z-score deviations after a warm-up period.

```text
GET /api/v1/anomalies?limit=100
```

This is intentionally lightweight and deterministic. More advanced seasonal and multivariate detectors belong in later releases.

## Cluster mode

Run nodes with peer URLs and a replication factor:

```bash
ROOOM_PEERS=http://metrics-2:9090,http://metrics-3:9090 \
ROOOM_REPLICATION_FACTOR=2 \
go run ./cmd/rooommetrics
```

Writes are asynchronously replicated to peers and reads fan out across peers. This is an MVP cluster mode, not a consensus or quorum implementation yet.

## Configuration

See `config.example.env`.

## Roadmap

1. Full PromQL parser/evaluator compatibility.
2. Remote Write 2.0 with metadata, exemplars and native histograms.
3. Native exponential histograms and sketches.
4. Segment files, checksums, block compaction, mmap indexes and background WAL checkpointing.
5. Consistent-hash sharding, membership, quorum replication and repair.
6. S3-compatible object-storage tier with distributed query planning.
7. Kubernetes service discovery and Operator.
8. Alert rules compatible with Prometheus rule files.
9. Query/result cache and cost-based planner.
10. RBAC, tenant quotas and per-query resource governance.

## Commercial licensing and support

The open-source MIT license remains available. Companies that require contractual commercial terms can purchase a separate ROOOMTECH commercial license agreement together with paid maintenance and support.

ROOOMTECH provides paid technical support, implementation and integration assistance, upgrade support, security and production-readiness support, SLA options, private builds, and custom development. A standard commercial software license agreement is available.

Contact: `tasuku.yoshioka@rooomtech.com`

## License

MIT. See `LICENSE`.
