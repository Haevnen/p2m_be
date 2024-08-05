package usecase

import (
	"context"

	"github.com/Haevnen/p2m_be/internal/pkg/base_service"
)

type BaseServiceUseCase[T, ID any] struct {
	BaseRepo base_service.Repository[T, ID]
}

func NewBaseServiceUseCase[T, ID any](rb base_service.Repository[T, ID]) *BaseServiceUseCase[T, ID] {
	return &BaseServiceUseCase[T, ID]{BaseRepo: rb}
}

func (u *BaseServiceUseCase[T, ID]) Add(ctx context.Context, t *T) error {
	return u.BaseRepo.Add(ctx, t)
}

func (u *BaseServiceUseCase[T, ID]) AddAll(ctx context.Context, entity *[]T) error {
	return u.BaseRepo.AddAll(ctx, entity)
}

func (u *BaseServiceUseCase[T, ID]) GetByID(ctx context.Context, id ID) (*T, error) {
	return u.BaseRepo.GetByID(ctx, id)
}

func (u *BaseServiceUseCase[T, ID]) GetSingle(ctx context.Context, params *T) *T {
	return u.BaseRepo.GetSingle(ctx, params)
}

func (u *BaseServiceUseCase[T, ID]) GetAll(ctx context.Context) (*[]T, error) {
	return u.BaseRepo.GetAll(ctx)
}

func (u *BaseServiceUseCase[T, ID]) WhereAll(ctx context.Context, params *T) (*[]T, error) {
	return u.BaseRepo.WhereAll(ctx, params)
}

func (u *BaseServiceUseCase[T, ID]) UpdateSingle(ctx context.Context, entity *T) error {
	return u.BaseRepo.UpdateSingle(ctx, entity)
}

func (u *BaseServiceUseCase[T, ID]) UpdateAll(entities *[]T, ctx context.Context) error {
	return u.BaseRepo.UpdateAll(entities, ctx)
}

func (u *BaseServiceUseCase[T, ID]) DeleteSingle(id any, ctx context.Context) error {
	return u.BaseRepo.DeleteSingle(id, ctx)
}

func (u *BaseServiceUseCase[T, ID]) SkipTake(skip int, take int, ctx context.Context) (*[]T, error) {
	return u.BaseRepo.SkipTake(skip, take, ctx)
}

func (u *BaseServiceUseCase[T, ID]) CountAll(ctx context.Context) int64 {
	return u.BaseRepo.CountAll(ctx)
}

func (u *BaseServiceUseCase[T, ID]) CountWhere(params *T, ctx context.Context) int64 {
	return u.BaseRepo.CountWhere(params, ctx)
}
