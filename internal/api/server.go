package api

import (
	"encoding/json"
	"github.com/softrenzu/Prometheus/internal/engine"
	"io"
	"net/http"
	"sync/atomic"
)

type Server struct { Store *engine.Store; Peers []string; ReplicationFactor int; requests atomic.Uint64 }
func (s *Server) Handler() http.Handler { mux:=http.NewServeMux(); mux.HandleFunc("/healthz",s.health); mux.HandleFunc("/readyz",s.health); mux.HandleFunc("/metrics",s.selfMetrics); mux.HandleFunc("/api/v1/write",s.remoteWrite); mux.HandleFunc("/api/v1/import/prometheus",s.importPrometheus); mux.HandleFunc("/v1/metrics",s.otlp); mux.HandleFunc("/api/v1/query",s.query); mux.HandleFunc("/api/v1/query_range",s.queryRange); mux.HandleFunc("/api/v1/label/__name__/values",s.metricNames); mux.HandleFunc("/api/v1/status/cardinality",s.cardinality); mux.HandleFunc("/api/v1/anomalies",s.anomalies); mux.HandleFunc("/internal/replicate",s.replicate); mux.HandleFunc("/internal/query",s.internalQuery); return http.HandlerFunc(func(w http.ResponseWriter,r *http.Request){s.requests.Add(1);mux.ServeHTTP(w,r)}) }
func tenant(r *http.Request) string { if v:=r.Header.Get("X-Scope-OrgID");v!=""{return v};return "default" }
func (s *Server) health(w http.ResponseWriter,r *http.Request){w.Header().Set("Content-Type","application/json");_,_=io.WriteString(w,`{"status":"ok"}`)}
func jsonResponse(w http.ResponseWriter,status int,v any){w.Header().Set("Content-Type","application/json");w.WriteHeader(status);_=json.NewEncoder(w).Encode(v)}
func promError(w http.ResponseWriter,status int,err error){jsonResponse(w,status,map[string]any{"status":"error","errorType":"bad_data","error":err.Error()})}
func (s *Server) appendAll(samples []engine.Sample)(int,error){written:=0;for _,sm:=range samples{if err:=s.Store.Append(sm);err!=nil{return written,err};written++};return written,nil}
