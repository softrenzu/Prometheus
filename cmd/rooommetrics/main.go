package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"syscall"
	"time"

	"github.com/softrenzu/Prometheus/internal/api"
	"github.com/softrenzu/Prometheus/internal/engine"
	"github.com/softrenzu/Prometheus/internal/scrape"
)

func envInt(k string, d int) int {
	if v, err := strconv.Atoi(os.Getenv(k)); err == nil && v > 0 {
		return v
	}
	return d
}
func envDuration(k string, d time.Duration) time.Duration {
	if v, err := time.ParseDuration(os.Getenv(k)); err == nil && v > 0 {
		return v
	}
	return d
}
func csvEnv(k string) []string {
	raw := strings.TrimSpace(os.Getenv(k))
	if raw == "" {
		return nil
	}
	parts := strings.Split(raw, ",")
	out := parts[:0]
	for _, p := range parts {
		if p = strings.TrimSpace(p); p != "" {
			out = append(out, p)
		}
	}
	return out
}

func main() {
	addr := os.Getenv("ROOOM_LISTEN")
	if addr == "" {
		addr = ":9090"
	}
	store, err := engine.NewStore(envInt("ROOOM_SHARDS", 128), envInt("ROOOM_CARDINALITY_LIMIT", 1_000_000), value("ROOOM_WAL_DIR", "./data"))
	if err != nil {
		log.Fatal(err)
	}
	defer store.Close()
	srv := &api.Server{Store: store, Peers: csvEnv("ROOOM_PEERS"), ReplicationFactor: envInt("ROOOM_REPLICATION_FACTOR", 1)}
	httpSrv := &http.Server{Addr: addr, Handler: srv.Handler(), ReadHeaderTimeout: 5 * time.Second, ReadTimeout: 30 * time.Second, WriteTimeout: 30 * time.Second, IdleTimeout: 60 * time.Second}
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer cancel()
	sc := &scrape.Scraper{Store: store, Targets: csvEnv("ROOOM_SCRAPE_TARGETS"), Interval: envDuration("ROOOM_SCRAPE_INTERVAL", 15*time.Second), Tenant: value("ROOOM_SCRAPE_TENANT", "default")}
	go sc.Run(ctx)
	retention := envDuration("ROOOM_RETENTION", 15*24*time.Hour)
	go func() {
		t := time.NewTicker(10 * time.Minute)
		defer t.Stop()
		for {
			select {
			case <-ctx.Done():
				return
			case <-t.C:
				store.Compact(retention)
			}
		}
	}()
	go func() {
		<-ctx.Done()
		shutdownCtx, c := context.WithTimeout(context.Background(), 10*time.Second)
		defer c()
		_ = httpSrv.Shutdown(shutdownCtx)
	}()
	api.LogServer(addr)
	if err := httpSrv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatal(err)
	}
}
func value(k, d string) string {
	if v := os.Getenv(k); v != "" {
		return v
	}
	return d
}
