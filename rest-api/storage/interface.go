package storage

import (
	"context"
	"errors"
	"github.com/d.tovstoluh/the-one-go-mid-test/rest-api/model"
)

var (
	ErrNotFound = errors.New("not found")
)

//go:generate mockgen -source=interface.go -destination=mocks/mock.go

type TaskStorage interface {
	Save(ctx context.Context, data *model.Task) (*model.Task, error)
	Get(ctx context.Context, id string) (*model.Task, error)
	Delete(ctx context.Context, id string) error
	GetAll(ctx context.Context) ([]*model.Task, error)
}
