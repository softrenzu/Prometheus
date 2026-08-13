package engine

import "strconv"

func Downsample(points []Point, stepMS int64, maxPoints int) []Point {
	if len(points) <= 1 {
		return points
	}
	if maxPoints <= 0 {
		maxPoints = 11000
	}
	span := points[len(points)-1].Timestamp - points[0].Timestamp
	if len(points) > maxPoints && span > 0 {
		auto := span / int64(maxPoints)
		if auto < 1 {
			auto = 1
		}
		if stepMS < auto {
			stepMS = auto
		}
	}
	if stepMS <= 1 {
		return points
	}
	out := make([]Point, 0, len(points))
	bucket := points[0].Timestamp / stepMS
	var sum float64
	n := 0
	var last int64
	flush := func() {
		if n > 0 {
			out = append(out, Point{Timestamp: last, Value: sum / float64(n)})
		}
	}
	for _, p := range points {
		b := p.Timestamp / stepMS
		if b != bucket {
			flush()
			bucket, sum, n = b, 0, 0
		}
		sum += p.Value
		n++
		last = p.Timestamp
	}
	flush()
	return out
}

func ParseStep(s string) int64 {
	if s == "" {
		return 0
	}
	if f, err := strconv.ParseFloat(s, 64); err == nil {
		return int64(f * 1000)
	}
	units := map[byte]int64{'s': 1000, 'm': 60000, 'h': 3600000, 'd': 86400000}
	u, ok := units[s[len(s)-1]]
	if !ok {
		return 0
	}
	f, err := strconv.ParseFloat(s[:len(s)-1], 64)
	if err != nil {
		return 0
	}
	return int64(f * float64(u))
}
