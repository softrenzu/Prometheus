package api

import("bytes";"encoding/json";"github.com/softrenzu/Prometheus/internal/engine";"net/http";"strings";"time")
func(s *Server)replicateBody(tenant string,samples []engine.Sample){if len(s.Peers)==0||s.ReplicationFactor<=1||len(samples)==0{return};b,_:=json.Marshal(samples);limit:=s.ReplicationFactor-1;if limit>len(s.Peers){limit=len(s.Peers)};client:=&http.Client{Timeout:3*time.Second};for i:=0;i<limit;i++{req,err:=http.NewRequest("POST",strings.TrimRight(s.Peers[i],"/")+"/internal/replicate",bytes.NewReader(b));if err!=nil{continue};req.Header.Set("Content-Type","application/json");req.Header.Set("X-Scope-OrgID",tenant);req.Header.Set("X-Rooom-Replicated","1");resp,err:=client.Do(req);if err==nil{resp.Body.Close()}}}
