package main

import "github.com/qppffod/microservice-project/types"

type Storer interface {
	Insert(types.Distance) error
}

type MemoryStore struct {
}

func NewMemoryStore() Storer {
	return &MemoryStore{}
}

func (m *MemoryStore) Insert(distance types.Distance) error {
	return nil
}
