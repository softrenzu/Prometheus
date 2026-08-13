package engine

import (
	"bufio"
	"encoding/json"
	"os"
	"path/filepath"
	"sync"
	"sync/atomic"
)

type shard struct {
	mu sync.RWMutex
	series map[string]*storedSeries
}

type Store struct {
	shards []shard
	cardinalityLimit int
	tenantMu sync.Mutex
	tenantSeries map[string]int
	anomaliesMu sync.RWMutex
	anomalies []Anomaly
	maxAnomalies int
	walMu sync.Mutex
	wal *os.File
	replaying atomic.Bool
	ingested atomic.Uint64
	rejected atomic.Uint64
}

func NewStore(shardCount, limit int, walDir string) (*Store, error) {
	if shardCount <= 0 { shardCount = 64 }
	if limit <= 0 { limit = 1_000_000 }
	s := &Store{shards: make([]shard, shardCount), cardinalityLimit: limit, tenantSeries: map[string]int{}, maxAnomalies: 10000}
	for i := range s.shards { s.shards[i].series = map[string]*storedSeries{} }
	if walDir != "" {
		if err := os.MkdirAll(walDir, 0755); err != nil { return nil, err }
		path := filepath.Join(walDir, "current.wal")
		if err := s.replay(path); err != nil { return nil, err }
		f, err := os.OpenFile(path, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
		if err != nil { return nil, err }
		s.wal = f
	}
	return s, nil
}
func (s *Store) Close() error { if s.wal != nil { return s.wal.Close() }; return nil }
func (s *Store) replay(path string) error {
	f, err := os.Open(path)
	if os.IsNotExist(err) { return nil }
	if err != nil { return err }
	defer f.Close()
	s.replaying.Store(true); defer s.replaying.Store(false)
	sc := bufio.NewScanner(f); sc.Buffer(make([]byte, 64*1024), 4*1024*1024)
	for sc.Scan() { var sm Sample; if json.Unmarshal(sc.Bytes(), &sm) == nil { _ = s.Append(sm) } }
	return sc.Err()
}
