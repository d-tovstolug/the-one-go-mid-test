package storage

import (
	"context"
	"github.com/d.tovstoluh/the-one-go-mid-test/rest-api/model"
	"github.com/google/uuid"
	"sync"
)

type memoryTaskStorage struct {
	mtx     sync.RWMutex
	storage map[string]*model.Task
}

func NewMemoryTaskStorage() TaskStorage {
	return &memoryTaskStorage{
		mtx:     sync.RWMutex{},
		storage: make(map[string]*model.Task),
	}
}

func (m *memoryTaskStorage) GetAll(_ context.Context) ([]*model.Task, error) {
	m.mtx.RLock()
	defer m.mtx.RUnlock()

	res := make([]*model.Task, 0)
	for _, task := range m.storage {
		res = append(res, task)
	}
	return res, nil
}

func (m *memoryTaskStorage) Save(_ context.Context, data *model.Task) (*model.Task, error) {
	if data.ID == "" {
		data.ID = uuid.NewString()
	}

	m.mtx.Lock()
	defer m.mtx.Unlock()
	m.storage[data.ID] = data
	return data, nil
}

func (m *memoryTaskStorage) Get(_ context.Context, id string) (*model.Task, error) {
	var data *model.Task
	m.mtx.RLock()
	defer m.mtx.RUnlock()

	data = m.storage[id]
	if data == nil {
		return nil, ErrNotFound
	}
	return data, nil
}

func (m *memoryTaskStorage) Delete(_ context.Context, id string) error {
	m.mtx.Lock()
	defer m.mtx.Unlock()
	delete(m.storage, id)
	return nil
}
