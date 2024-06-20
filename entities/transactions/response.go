package transactions

type DoReconciliationResponse struct {
	TransactionsProceed       int `json:"transaction_proceed"`
	MatchedTransaction        int `json:"matched_transaction"`
	UnmatchedTransaction      int `json:"unmatched_transaction"`
	MissingBankStatements     map[string][]BankStatements `json:"missing_bank_statements"`
	MissingSystemTransactions []SystemTransactions `json:"missing_system_transactions"`
	TotalDiscrepancies        float64 `json:"total_discripencies"`
}
