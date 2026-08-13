package scrape

import (
	"context"
	"log"
	"net/http"
	"strings"
	"time"
	"github.com/softrenzu/Prometheus/internal/engine"
	"github.com/softrenzu/Prometheus/internal/ingest"
)
type Scraper struct{Store *engine.Store;Targets []string;Interval time.Duration;Tenant string;Client *http.Client}
func (s *Scraper) Run(ctx context.Context){if len(s.Targets)==0{return};if s.Interval<=0{s.Interval=15*time.Second};if s.Tenant==""{s.Tenant="default"};if s.Client==nil{s.Client=&http.Client{Timeout:10*time.Second}};ticker:=time.NewTicker(s.Interval);defer ticker.Stop();s.once(ctx);for{select{case<-ctx.Done():return;case<-ticker.C:s.once(ctx)}}}
func (s *Scraper) once(ctx context.Context){for _,target:=range s.Targets{target=strings.TrimSpace(target);if target==""{continue};req,err:=http.NewRequestWithContext(ctx,"GET",target,nil);if err!=nil{continue};resp,err:=s.Client.Do(req);if err!=nil{log.Printf("scrape %s: %v",target,err);continue};samples,err:=ingest.ParsePrometheusText(resp.Body,s.Tenant);resp.Body.Close();if err!=nil{log.Printf("scrape parse %s: %v",target,err);continue};for _,sm:=range samples{if err:=s.Store.Append(sm);err!=nil{log.Printf("scrape append %s: %v",target,err);break}}}}
