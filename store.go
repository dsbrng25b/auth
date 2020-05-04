package auth

import (
	"context"
	"sync"
)

type Store interface {
	Get(ctx context.Context, key string) ([]byte, error)
	Set(ctx context.Context, key string, value []byte) error
	Delete(ctx context.Context, key string) error
}

type SubjectStore interface {
	Get(ctx context.Context, key string) (*Subject, error)
	Set(ctx context.Context, key string, sub *Subject) error
	Delete(ctx context.Context, key string) error
}

type MemoryStore struct {
	sync.RWMutex
	store map[string][]byte
}

func NewMemoryStore() *MemoryStore {
	return &MemoryStore{
		store: map[string][]byte{},
	}
}

func (m *MemoryStore) Get(_ context.Context, key string) ([]byte, error) {
	m.RLock()
	v, _ := m.store[key]
	value := clone(v)
	m.RUnlock()
	return value, nil
}

func (m *MemoryStore) Set(_ context.Context, key string, value []byte) error {
	save := clone(value)
	m.Lock()
	m.store[key] = save
	m.Unlock()
	return nil
}

func (m *MemoryStore) Delete(_ context.Context, key string) error {
	m.Lock()
	delete(m.store, key)
	m.Unlock()
	return nil
}

func clone(in []byte) (out []byte) {
	if in == nil {
		return
	} else {
		out = make([]byte, len(in))
		copy(out, in)
	}
	return
}

type MemorySubjectStore struct {
	sync.RWMutex
	store map[string]Subject
}

func NewMemorySubjectStore() *MemorySubjectStore {
	return &MemorySubjectStore{
		store: map[string]Subject{},
	}
}

func (m *MemorySubjectStore) Get(_ context.Context, key string) (*Subject, error) {
	m.RLock()
	sub, ok := m.store[key]
	m.RUnlock()
	if ok {
		return &sub, nil
	}
	return nil, nil
}

func (m *MemorySubjectStore) Set(_ context.Context, key string, sub *Subject) error {
	if sub == nil {
		return nil
	}
	m.Lock()
	m.store[key] = *sub
	m.Unlock()
	return nil
}

func (m *MemorySubjectStore) Delete(_ context.Context, key string) error {
	m.Lock()
	delete(m.store, key)
	m.Unlock()
	return nil
}
