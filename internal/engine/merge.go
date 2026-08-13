package engine

import "sort"

// MergeSeries deduplicates series and timestamp-identical points. When replicas
// contain the same timestamp, the later input wins. This makes peer fan-out
// safe when the same sample exists on multiple replicas.
func MergeSeries(in []Series) []Series {
	byKey := make(map[string]*Series, len(in))
	for _, ss := range in {
		key := CanonicalKey(ss.Tenant, ss.Metric, ss.Labels)
		dst := byKey[key]
		if dst == nil {
			clone := Series{Tenant: ss.Tenant, Metric: ss.Metric, Labels: CloneLabels(ss.Labels)}
			dst = &clone
			byKey[key] = dst
		}
		points := make(map[int64]float64, len(dst.Points)+len(ss.Points))
		for _, p := range dst.Points { points[p.Timestamp] = p.Value }
		for _, p := range ss.Points { points[p.Timestamp] = p.Value }
		dst.Points = dst.Points[:0]
		for ts, v := range points { dst.Points = append(dst.Points, Point{Timestamp: ts, Value: v}) }
		sort.Slice(dst.Points, func(i, j int) bool { return dst.Points[i].Timestamp < dst.Points[j].Timestamp })
	}
	out := make([]Series, 0, len(byKey))
	for _, ss := range byKey { out = append(out, *ss) }
	sort.Slice(out, func(i, j int) bool { return LabelString(out[i].Metric, out[i].Labels) < LabelString(out[j].Metric, out[j].Labels) })
	return out
}

func AggregateRange(fn string, in []Series) []Series {
	if fn == "" || fn == "rate" { return in }
	type agg struct { sum, min, max float64; n int }
	byTS := map[int64]agg{}
	for _, ss := range in {
		for _, p := range ss.Points {
			a := byTS[p.Timestamp]
			if a.n == 0 { a.min, a.max = p.Value, p.Value } else { if p.Value < a.min { a.min = p.Value }; if p.Value > a.max { a.max = p.Value } }
			a.sum += p.Value; a.n++; byTS[p.Timestamp] = a
		}
	}
	pts := make([]Point, 0, len(byTS))
	for ts, a := range byTS {
		var v float64
		switch fn { case "sum": v=a.sum; case "avg": v=a.sum/float64(a.n); case "min": v=a.min; case "max": v=a.max; case "count": v=float64(a.n); default: continue }
		pts = append(pts, Point{Timestamp: ts, Value: v})
	}
	sort.Slice(pts, func(i, j int) bool { return pts[i].Timestamp < pts[j].Timestamp })
	return []Series{{Metric: fn, Labels: map[string]string{}, Points: pts}}
}
