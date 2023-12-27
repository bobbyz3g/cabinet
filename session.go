package cabinet

import (
	"io"
	"sync"
)

type Code string

type Translator struct {
	Name   string
	Reader io.Reader
	Done   chan struct{}
}

type Sessions struct {
	mu sync.RWMutex
	ts map[Code]*Translator
}

func (s *Sessions) Push(code Code, t *Translator) {
	s.mu.Lock()
	defer s.mu.Unlock()
	if s.ts == nil {
		s.ts = make(map[Code]*Translator)
	}
	s.ts[code] = t
}

func (s *Sessions) Pop(code Code) (*Translator, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	t, ok := s.ts[code]
	return t, ok
}
