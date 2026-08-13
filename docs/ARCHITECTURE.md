# RooomMetrics architecture

Core design: sharded in-memory series store, crash-recovery WAL, per-tenant cardinality limits, Prometheus-compatible APIs, OTLP JSON ingestion, peer replication, distributed read fan-out, adaptive range downsampling, and online anomaly detection.

This repository is an experimental MVP. Full PromQL, quorum replication, object storage, and native histograms remain roadmap items.
