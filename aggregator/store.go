package main

import "github.com/qppffod/microservice-project/types"

type Storer interface {
	Insert(types.Distance) error
}

type MemoryStore struct {
	data map[int]float64
}

func NewMemoryStore() Storer {
	return &MemoryStore{
		data: make(map[int]float64),
	}
}

func (m *MemoryStore) Insert(distance types.Distance) error {
	m.data[distance.OBUID] += distance.Value
	return nil
}
