package db_wrap

import (
	"context"

	"github.com/Haevnen/p2m_be/pkg/logger"
)

// TxManager manage transaction.
// Expected to inject in the interactor to execute transaction
type TxManager interface {
	// TransactionExec start transaction, and commit or rollback
	TransactionExec(ctx context.Context, fn func(ctx context.Context) error) error
}

type txKeyType string

const (
	txManagerTransactionKey txKeyType = "p2m-tx-manager-key"
)

type txManager struct {
	db DBGetter
}

// NewTxManager construct TxManager
func NewTxManager(db DBGetter) TxManager {
	return &txManager{db: db}
}

func (t *txManager) TransactionExec(ctx context.Context, fn func(ctx context.Context) error) error {
	tx := t.db.Get(ctx).Begin()
	txCtx := context.WithValue(ctx, txManagerTransactionKey, tx)

	err := fn(txCtx)
	if err == nil {
		return tx.Commit().Error
	}

	txErr := tx.Rollback().Error
	if txErr != nil {
		logger.Error(txErr.Error())
	}

	return err
}
