package engine

import (
	"errors"
	"strings"
)

type Query struct {
	Func     string
	Metric   string
	Matchers map[string]string
}

func ParseQuery(expr string) (Query, error) {
	expr = strings.TrimSpace(expr)
	q := Query{Matchers: map[string]string{}}
	for _, fn := range []string{"sum", "avg", "min", "max", "rate", "count"} {
		prefix := fn + "("
		if strings.HasPrefix(expr, prefix) && strings.HasSuffix(expr, ")") {
			q.Func = fn
			expr = strings.TrimSpace(expr[len(prefix) : len(expr)-1])
			break
		}
	}
	if i := strings.IndexByte(expr, '{'); i >= 0 {
		if !strings.HasSuffix(expr, "}") {
			return q, errors.New("invalid selector")
		}
		q.Metric = strings.TrimSpace(expr[:i])
		body := strings.TrimSpace(expr[i+1 : len(expr)-1])
		if body != "" {
			for _, p := range splitCSV(body) {
				kv := strings.SplitN(p, "=", 2)
				if len(kv) != 2 {
					return q, errors.New("only equality label matchers are supported")
				}
				q.Matchers[strings.TrimSpace(kv[0])] = strings.Trim(strings.TrimSpace(kv[1]), "\"")
			}
		}
	} else {
		q.Metric = strings.TrimSpace(expr)
	}
	if q.Metric == "" {
		return q, errors.New("metric is required")
	}
	return q, nil
}

func splitCSV(s string) []string {
	var out []string
	start := 0
	quoted := false
	for i := 0; i < len(s); i++ {
		if s[i] == '"' {
			quoted = !quoted
		}
		if s[i] == ',' && !quoted {
			out = append(out, strings.TrimSpace(s[start:i]))
			start = i + 1
		}
	}
	return append(out, strings.TrimSpace(s[start:]))
}
