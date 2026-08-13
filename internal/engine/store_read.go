package engine

import "sort"

func labelsMatch(labels, matchers map[string]string) bool { for k,v := range matchers { if labels[k] != v { return false } }; return true }
func (s *Store) Select(tenant, metric string, matchers map[string]string, start, end int64) []Series {
	out := []Series{}
	for i := range s.shards { sh:=&s.shards[i]; sh.mu.RLock(); for _,ss := range sh.series { if ss.Tenant!=tenant || ss.Metric!=metric || !labelsMatch(ss.Labels,matchers) { continue }; pts:=[]Point{}; lo:=sort.Search(len(ss.Points),func(i int)bool{return ss.Points[i].Timestamp>=start}); for j:=lo;j<len(ss.Points);j++ { p:=ss.Points[j]; if end>0 && p.Timestamp>end { break }; pts=append(pts,p) }; out=append(out,Series{Tenant:tenant,Metric:metric,Labels:CloneLabels(ss.Labels),Points:pts}) }; sh.mu.RUnlock() }
	sort.Slice(out,func(i,j int)bool{return LabelString(out[i].Metric,out[i].Labels)<LabelString(out[j].Metric,out[j].Labels)}); return out
}
func (s *Store) Metrics(tenant string) []string { seen:=map[string]struct{}{}; for i:=range s.shards { sh:=&s.shards[i]; sh.mu.RLock(); for _,ss:=range sh.series { if ss.Tenant==tenant { seen[ss.Metric]=struct{}{} } }; sh.mu.RUnlock() }; out:=make([]string,0,len(seen)); for m:=range seen { out=append(out,m) }; sort.Strings(out); return out }
func (s *Store) Cardinality(tenant string) int { s.tenantMu.Lock(); defer s.tenantMu.Unlock(); return s.tenantSeries[tenant] }
func (s *Store) CardinalityLimit() int { return s.cardinalityLimit }
func (s *Store) Anomalies(tenant string,limit int) []Anomaly { if limit<=0||limit>1000{limit=100}; s.anomaliesMu.RLock(); defer s.anomaliesMu.RUnlock(); out:=[]Anomaly{}; for i:=len(s.anomalies)-1;i>=0&&len(out)<limit;i-- { if s.anomalies[i].Tenant==tenant { out=append(out,s.anomalies[i]) } }; return out }
func (s *Store) Ingested() uint64 { return s.ingested.Load() }
func (s *Store) Rejected() uint64 { return s.rejected.Load() }
