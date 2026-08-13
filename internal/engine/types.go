package engine

import (
	"sort"
	"strconv"
	"strings"
)

type Point struct {
	Timestamp int64   `json:"timestamp"`
	Value     float64 `json:"value"`
}

type Sample struct {
	Tenant    string            `json:"tenant,omitempty"`
	Metric    string            `json:"metric"`
	Labels    map[string]string `json:"labels,omitempty"`
	Timestamp int64             `json:"timestamp"`
	Value     float64           `json:"value"`
}

type Series struct {
	Tenant string            `json:"tenant"`
	Metric string            `json:"metric"`
	Labels map[string]string `json:"labels"`
	Points []Point           `json:"points"`
}

type Anomaly struct {
	Tenant    string            `json:"tenant"`
	Metric    string            `json:"metric"`
	Labels    map[string]string `json:"labels"`
	Timestamp int64             `json:"timestamp"`
	Value     float64           `json:"value"`
	Mean      float64           `json:"mean"`
	StdDev    float64           `json:"stddev"`
	ZScore    float64           `json:"zscore"`
}

func CanonicalKey(tenant, metric string, labels map[string]string) string {
	keys := make([]string, 0, len(labels))
	for k := range labels {
		if k == "__name__" {
			continue
		}
		keys = append(keys, k)
	}
	sort.Strings(keys)
	var b strings.Builder
	b.WriteString(tenant)
	b.WriteByte(0)
	b.WriteString(metric)
	for _, k := range keys {
		b.WriteByte(0)
		b.WriteString(k)
		b.WriteByte('=')
		b.WriteString(labels[k])
	}
	return b.String()
}

func CloneLabels(in map[string]string) map[string]string {
	out := make(map[string]string, len(in))
	for k, v := range in {
		out[k] = v
	}
	return out
}

func HashKey(key string, shards int) int {
	var h uint64 = 1469598103934665603
	for i := 0; i < len(key); i++ {
		h ^= uint64(key[i])
		h *= 1099511628211
	}
	return int(h % uint64(shards))
}

func LabelString(metric string, labels map[string]string) string {
	if len(labels) == 0 {
		return metric
	}
	keys := make([]string, 0, len(labels))
	for k := range labels {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	var b strings.Builder
	b.WriteString(metric)
	b.WriteByte('{')
	for i, k := range keys {
		if i > 0 {
			b.WriteByte(',')
		}
		b.WriteString(k)
		b.WriteString("=\"")
		b.WriteString(strings.ReplaceAll(labels[k], "\"", "\\\""))
		b.WriteByte('"')
	}
	b.WriteByte('}')
	return b.String()
}

func ParseInt64(s string, def int64) int64 {
	if s == "" {
		return def
	}
	v, err := strconv.ParseInt(s, 10, 64)
	if err != nil {
		return def
	}
	return v
}
