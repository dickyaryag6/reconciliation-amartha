package transactions

import "mime/multipart"

type DoReconciliationRequest struct {
	SystemTransactions multipart.File
	BankStatements     multipart.File
}

