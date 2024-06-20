package httphandlers

import "net/http"

type TransactionHandler interface {
	HandleReconciliation(w http.ResponseWriter, r *http.Request) 
}