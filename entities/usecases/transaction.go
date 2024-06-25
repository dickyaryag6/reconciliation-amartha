package usecases

import (
	"amartha-test/entities/transactions"
	"context"
)


//go:generate mockgen -destination mock/mock_transactions.go -source=transaction.go TransactionUsecase

type TransactionUsecase interface {
	DoReconciliation(ctx context.Context, param transactions.DoReconciliationRequest) (transactions.DoReconciliationResponse, error)
}