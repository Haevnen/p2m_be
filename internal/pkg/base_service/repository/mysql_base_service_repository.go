package repository

import (
	"context"

	"gorm.io/gorm"
)

type MysqlBaseServiceRepository[T, ID any] struct {
	*gorm.DB
}

func NewMysqlBaseServiceRepository[T, ID any](db *gorm.DB) *MysqlBaseServiceRepository[T, ID] {
	return &MysqlBaseServiceRepository[T, ID]{DB: db}
}

func (r *MysqlBaseServiceRepository[T, ID]) Add(ctx context.Context, t *T) error {
	return r.DB.WithContext(ctx).Create(t).Error
}

func (r *MysqlBaseServiceRepository[T, ID]) AddAll(ctx context.Context, entity *[]T) error {
	return r.DB.WithContext(ctx).Create(&entity).Error
}

func (r *MysqlBaseServiceRepository[T, ID]) GetByID(ctx context.Context, id ID) (*T, error) {
	var entity T
	err := r.DB.WithContext(ctx).Model(&entity).Where("id = ?", id).FirstOrInit(&entity).Error
	if err != nil {
		return nil, err
	}
	return &entity, nil
}

func (r *MysqlBaseServiceRepository[T, ID]) GetSingle(ctx context.Context, params *T) *T {
	var entity T
	r.DB.WithContext(ctx).Where(&params).FirstOrInit(&entity)
	return &entity
}

func (r *MysqlBaseServiceRepository[T, ID]) GetAll(ctx context.Context) (*[]T, error) {
	var entities []T
	err := r.DB.WithContext(ctx).Find(&entities).Error
	if err != nil {
		return nil, err
	}
	return &entities, nil
}

func (r *MysqlBaseServiceRepository[T, ID]) WhereAll(ctx context.Context, params *T) (*[]T, error) {
	var entities []T
	err := r.DB.WithContext(ctx).Where(&params).Find(&entities).Error
	if err != nil {
		return nil, err
	}
	return &entities, nil
}

func (r *MysqlBaseServiceRepository[T, ID]) UpdateSingle(ctx context.Context, entity *T) error {
	return r.DB.WithContext(ctx).Save(&entity).Error
}

func (r *MysqlBaseServiceRepository[T, ID]) UpdateAll(entities *[]T, ctx context.Context) error {
	return r.DB.WithContext(ctx).Save(&entities).Error
}

func (r *MysqlBaseServiceRepository[T, ID]) DeleteSingle(id any, ctx context.Context) error {
	var entity T
	return r.DB.WithContext(ctx).Model(&entity).Delete("id = ?", id).Error
}

func (r *MysqlBaseServiceRepository[T, ID]) SkipTake(skip int, take int, ctx context.Context) (*[]T, error) {
	var entities []T
	err := r.DB.WithContext(ctx).Offset(skip).Limit(take).Find(&entities).Error
	if err != nil {
		return nil, err
	}
	return &entities, nil
}

func (r *MysqlBaseServiceRepository[T, ID]) CountAll(ctx context.Context) int64 {
	var entity T
	var count int64
	r.DB.WithContext(ctx).Model(&entity).Count(&count)
	return count
}

func (r *MysqlBaseServiceRepository[T, ID]) CountWhere(params *T, ctx context.Context) int64 {
	var entity T
	var count int64
	r.DB.WithContext(ctx).Model(&entity).Where(&params).Count(&count)
	return count
}
