package engine

import (
	"encoding/json"
	"fmt"
)

func (s *Store) appendWAL(sm Sample) error {
	if s.wal == nil || s.replaying.Load() { return nil }
	b, err := json.Marshal(sm)
	if err != nil { return err }
	s.walMu.Lock()
	defer s.walMu.Unlock()
	if _, err = s.wal.Write(append(b, '\n')); err != nil { return err }
	return s.wal.Sync()
}
func (s *Store) walError(err error) error { return fmt.Errorf("wal append: %w", err) }
