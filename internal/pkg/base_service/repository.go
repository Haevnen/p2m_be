package base_service

import "context"

type Repository[T, ID any] interface {
	Add(ctx context.Context, t *T) error
	AddAll(ctx context.Context, entity *[]T) error
	GetByID(ctx context.Context, id ID) (*T, error)
	GetSingle(ctx context.Context, params *T) *T
	GetAll(ctx context.Context) (*[]T, error)
	WhereAll(ctx context.Context, params *T) (*[]T, error)
	UpdateSingle(ctx context.Context, entity *T) error
	UpdateAll(entities *[]T, ctx context.Context) error
	DeleteSingle(id any, ctx context.Context) error
	SkipTake(skip int, take int, ctx context.Context) (*[]T, error)
	CountAll(ctx context.Context) int64
	CountWhere(params *T, ctx context.Context) int64
}
