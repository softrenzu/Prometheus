package engine

import "math"

type onlineStats struct {
	N        int64
	Mean, M2 float64
}

func (s *onlineStats) Add(x float64) (mean, stddev, z float64, anomalous bool) {
	if s.N >= 20 {
		variance := s.M2 / float64(s.N-1)
		if variance > 0 {
			stddev = math.Sqrt(variance)
			z = math.Abs((x - s.Mean) / stddev)
			anomalous = z >= 4
		}
	}
	s.N++
	delta := x - s.Mean
	s.Mean += delta / float64(s.N)
	s.M2 += delta * (x - s.Mean)
	return s.Mean, stddev, z, anomalous
}

type storedSeries struct {
	Tenant, Metric string
	Labels         map[string]string
	Points         []Point
	Stats          onlineStats
}
