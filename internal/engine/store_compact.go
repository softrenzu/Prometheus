package engine

import (
	"sort"
	"time"
)
func (s *Store) Compact(retention time.Duration) int {
	if retention<=0{return 0}; cutoff:=time.Now().Add(-retention).UnixMilli(); removed:=0
	for i:=range s.shards { sh:=&s.shards[i]; sh.mu.Lock(); for key,ss:=range sh.series { idx:=sort.Search(len(ss.Points),func(i int)bool{return ss.Points[i].Timestamp>=cutoff}); if idx>0 { removed+=idx; ss.Points=append([]Point(nil),ss.Points[idx:]...) }; if len(ss.Points)==0 { delete(sh.series,key); s.tenantMu.Lock(); s.tenantSeries[ss.Tenant]--; s.tenantMu.Unlock() } }; sh.mu.Unlock() }; return removed
}
