package api

import("fmt";"log";"net/http";"strconv")
func(s *Server)metricNames(w http.ResponseWriter,r *http.Request){jsonResponse(w,200,map[string]any{"status":"success","data":s.Store.Metrics(tenant(r))})}
func(s *Server)cardinality(w http.ResponseWriter,r *http.Request){jsonResponse(w,200,map[string]any{"status":"success","data":map[string]any{"tenant":tenant(r),"active_series":s.Store.Cardinality(tenant(r)),"limit":s.Store.CardinalityLimit()}})}
func(s *Server)anomalies(w http.ResponseWriter,r *http.Request){limit,_:=strconv.Atoi(r.URL.Query().Get("limit"));jsonResponse(w,200,map[string]any{"status":"success","data":s.Store.Anomalies(tenant(r),limit)})}
func(s *Server)selfMetrics(w http.ResponseWriter,r *http.Request){w.Header().Set("Content-Type","text/plain; version=0.0.4");fmt.Fprintf(w,"# HELP rooommetrics_ingested_samples_total Total samples accepted.\n# TYPE rooommetrics_ingested_samples_total counter\nrooommetrics_ingested_samples_total %d\n# HELP rooommetrics_rejected_samples_total Total samples rejected.\n# TYPE rooommetrics_rejected_samples_total counter\nrooommetrics_rejected_samples_total %d\n# HELP rooommetrics_http_requests_total Total HTTP requests.\n# TYPE rooommetrics_http_requests_total counter\nrooommetrics_http_requests_total %d\n",s.Store.Ingested(),s.Store.Rejected(),s.requests.Load())}
func LogServer(addr string){log.Printf("RooomMetrics listening on %s",addr)}
