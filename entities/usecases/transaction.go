package usecases

import (
	"amartha-test/entities/transactions"
	"context"
)

type TransactionUsecase interface {
	DoReconciliation(ctx context.Context, param transactions.DoReconciliationRequest) (transactions.DoReconciliationResponse, error)
}