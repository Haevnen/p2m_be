package db_wrap

import (
	"context"

	"gorm.io/gorm"
)

type DBGetter interface {
	Get(ctx context.Context) *gorm.DB
	GetTxOrDB(ctx context.Context) *gorm.DB
}

// NewDBGetter is constructor for DBGetter
func NewDBGetter(db *gorm.DB) DBGetter {
	return &gormDb{db}
}

// gormDb is wrapper for gorm.DB
type gormDb struct {
	db *gorm.DB
}

func (d *gormDb) Get(ctx context.Context) *gorm.DB {
	return d.db.WithContext(ctx)
}

func (d *gormDb) GetTxOrDB(ctx context.Context) *gorm.DB {
	tx, ok := ctx.Value(txManagerTransactionKey).(*gorm.DB)
	if ok {
		return tx.WithContext(ctx)
	}
	return d.db.WithContext(ctx)
}
