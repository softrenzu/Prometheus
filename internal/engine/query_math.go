package engine

import "math"

func AggregateInstant(fn string, series []Series) []Series {
	if fn == "" || fn == "rate" {
		return series
	}
	values := make([]float64, 0, len(series))
	ts := int64(0)
	for _, s := range series {
		if len(s.Points) == 0 {
			continue
		}
		p := s.Points[len(s.Points)-1]
		values = append(values, p.Value)
		if p.Timestamp > ts {
			ts = p.Timestamp
		}
	}
	if len(values) == 0 {
		return nil
	}
	v := values[0]
	switch fn {
	case "sum":
		v = 0
		for _, x := range values {
			v += x
		}
	case "avg":
		v = 0
		for _, x := range values {
			v += x
		}
		v /= float64(len(values))
	case "min":
		for _, x := range values[1:] {
			v = math.Min(v, x)
		}
	case "max":
		for _, x := range values[1:] {
			v = math.Max(v, x)
		}
	case "count":
		v = float64(len(values))
	}
	return []Series{{Metric: fn, Labels: map[string]string{}, Points: []Point{{Timestamp: ts, Value: v}}}}
}

func Rate(s Series) Series {
	out := Series{Tenant: s.Tenant, Metric: s.Metric, Labels: CloneLabels(s.Labels)}
	for i := 1; i < len(s.Points); i++ {
		dt := float64(s.Points[i].Timestamp-s.Points[i-1].Timestamp) / 1000
		if dt <= 0 {
			continue
		}
		delta := s.Points[i].Value - s.Points[i-1].Value
		if delta < 0 {
			delta = s.Points[i].Value
		}
		out.Points = append(out.Points, Point{Timestamp: s.Points[i].Timestamp, Value: delta / dt})
	}
	return out
}
