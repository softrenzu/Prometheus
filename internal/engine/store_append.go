package engine

import (
	"errors"
	"sort"
	"time"
)

var ErrCardinalityLimit = errors.New("tenant cardinality limit exceeded")

func (s *Store) Append(sm Sample) error {
	if sm.Tenant == "" { sm.Tenant = "default" }
	if sm.Metric == "" { return errors.New("metric is required") }
	if sm.Timestamp == 0 { sm.Timestamp = time.Now().UnixMilli() }
	if sm.Labels == nil { sm.Labels = map[string]string{} }
	key := CanonicalKey(sm.Tenant, sm.Metric, sm.Labels)
	sh := &s.shards[HashKey(key, len(s.shards))]
	sh.mu.Lock()
	ss, exists := sh.series[key]
	if !exists {
		s.tenantMu.Lock()
		if s.tenantSeries[sm.Tenant] >= s.cardinalityLimit { s.tenantMu.Unlock(); sh.mu.Unlock(); s.rejected.Add(1); return ErrCardinalityLimit }
		s.tenantSeries[sm.Tenant]++
		s.tenantMu.Unlock()
		ss = &storedSeries{Tenant: sm.Tenant, Metric: sm.Metric, Labels: CloneLabels(sm.Labels)}
		sh.series[key] = ss
	}
	if n := len(ss.Points); n > 0 && sm.Timestamp < ss.Points[n-1].Timestamp {
		idx := sort.Search(n, func(i int) bool { return ss.Points[i].Timestamp >= sm.Timestamp })
		ss.Points = append(ss.Points, Point{}); copy(ss.Points[idx+1:], ss.Points[idx:]); ss.Points[idx] = Point{Timestamp: sm.Timestamp, Value: sm.Value}
	} else { ss.Points = append(ss.Points, Point{Timestamp: sm.Timestamp, Value: sm.Value}) }
	mean, stddev, z, anom := ss.Stats.Add(sm.Value); sh.mu.Unlock()
	if anom {
		s.anomaliesMu.Lock(); s.anomalies = append(s.anomalies, Anomaly{Tenant: sm.Tenant, Metric: sm.Metric, Labels: CloneLabels(sm.Labels), Timestamp: sm.Timestamp, Value: sm.Value, Mean: mean, StdDev: stddev, ZScore: z})
		if len(s.anomalies) > s.maxAnomalies { s.anomalies = append([]Anomaly(nil), s.anomalies[len(s.anomalies)-s.maxAnomalies:]...) }
		s.anomaliesMu.Unlock()
	}
	if err := s.appendWAL(sm); err != nil { return s.walError(err) }
	s.ingested.Add(1); return nil
}
