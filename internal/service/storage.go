package service

import (
	"context"
)

type Repository interface {
	Find(context.Context, *CreateEntity) (*[]Entity, error)
	Create(context.Context, *CreateEntity) error
	Check(context.Context, *CreateEntity) (bool, error)
}