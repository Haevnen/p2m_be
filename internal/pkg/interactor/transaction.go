package interactor

import (
	"context"

	"github.com/Haevnen/p2m_be/internal/pkg/dal"
	"github.com/Haevnen/p2m_be/pkg/logger"
)

const (
	txTransactionKey = "tx-transaction-key"
)

type TxManager struct{}

func NewTxManager() *TxManager {
	return &TxManager{}
}

func (t *TxManager) TransactionExec(ctx context.Context, fn func(context.Context) error) error {
	tx := dal.Q.Begin()
	txCtx := context.WithValue(ctx, txTransactionKey, tx)

	err := fn(txCtx)
	if err == nil {
		return tx.Commit()
	}

	txErr := tx.Rollback()
	if txErr != nil {
		logger.Errorf("Transaction error: %v", txErr)
	}
	return err
}
