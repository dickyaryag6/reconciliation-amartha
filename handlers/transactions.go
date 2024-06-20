package handlers

import (
	"amartha-test/entities/transactions"
	"amartha-test/entities/usecases"
	libError "amartha-test/errors"
	"amartha-test/response"
	"context"
	"errors"
	"net/http"
)

type TransactionHandler struct {
	TransactionUsecase usecases.TransactionUsecase
}

func NewTransactionHandler(handler TransactionHandler) TransactionHandler {
	return handler
}

func (handler TransactionHandler) HandleReconciliation(w http.ResponseWriter, r *http.Request) {
	ctx := context.Background()

	// to limit the file size to be no greater than 10 MB
	err := r.ParseMultipartForm(10 << 20)
	if err != nil {
		return
	}

	// get file 1 from form
	bankStatements, bankStatementsFileHeader, err := r.FormFile("bank_statements")
	if err != nil {
		if errors.Is(err, http.ErrMissingFile) {
			libError.SetBadRequestErrorForHandler(w, "File is not found")
			return
		}
		return
	}

	// get file 2 from form
	systemTransactions, systemTransactionsFileHeader, err := r.FormFile("system_transactions")
	if err != nil {
		if errors.Is(err, http.ErrMissingFile) {
			libError.SetBadRequestErrorForHandler(w, "File is not found")
			return
		}
		libError.SetInternalServerErrorForHandler(w, err)
		return
	}

	// check if file is csv
	if len(bankStatementsFileHeader.Header["Content-Type"]) > 0 && bankStatementsFileHeader.Header["Content-Type"][0] != "text/csv" {
		libError.SetBadRequestErrorForHandler(w, "File Upload is not csv")
		return
	}

	// check if file is csv
	if len(systemTransactionsFileHeader.Header["Content-Type"]) > 0 && systemTransactionsFileHeader.Header["Content-Type"][0] != "text/csv" {
		libError.SetBadRequestErrorForHandler(w, "File Upload is not csv")
		return
	}

	result, err := handler.TransactionUsecase.DoReconciliation(ctx, transactions.DoReconciliationRequest{
		SystemTransactions: systemTransactions,
		BankStatements:     bankStatements,
	})
	if err != nil {
		libError.SetError(w, err)
		return
	}

	response.SetOK(w, result)

}
