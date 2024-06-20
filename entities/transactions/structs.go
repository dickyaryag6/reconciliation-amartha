package transactions

import "time"

type SystemTransactions struct {
	TransactionID       string    `json:"trxID" csv:"trxID"`
	Amount              string    `json:"amount" csv:"amount"`
	RealAmount          float64   `json:"-" csv:"-"`
	Type                int       `json:"type" csv:"type"`
	TransactionTime     string    `json:"transactionTime" csv:"transactionTime"`
	RealTransactionTime time.Time `json:"-" csv:"-"`
}

// SortByRealDateSystemTransaction implements sort.Interface for []BankStatements based on the RealTransactionTime field.
type SortByRealDateSystemTransaction []*SystemTransactions

func (a SortByRealDateSystemTransaction) Len() int      { return len(a) }
func (a SortByRealDateSystemTransaction) Swap(i, j int) { a[i], a[j] = a[j], a[i] }
func (a SortByRealDateSystemTransaction) Less(i, j int) bool {
	if a[i].RealTransactionTime != a[j].RealTransactionTime {
		return a[i].RealTransactionTime.Before(a[j].RealTransactionTime)
	}
	return a[i].Type < a[j].Type

}

type BankStatements struct {
	ID         string    `json:"unique_identifier" csv:"unique_identifier"` // contain bank source information, example : BCA_123, BRI_256, separated by underscore
	Amount     string    `json:"amount" csv:"amount"`                       // if negative, then type is DEBIT, else CREDIT
	RealAmount float64   `json:"-" csv:"-"`
	Date       string    `json:"date" csv:"date"`
	RealDate   time.Time `json:"-" csv:"-"`
	BankSource string    `json:"bank_source"  csv:"-"`
	Type       int       `json:"-" csv:"-"`
}

// SortByRealDateSystemTransaction implements sort.Interface for []BankStatements based on the RealDate field.
type SortByRealDateBankStatement []*BankStatements

func (a SortByRealDateBankStatement) Len() int           { return len(a) }
func (a SortByRealDateBankStatement) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a SortByRealDateBankStatement) Less(i, j int) bool { 	
	if a[i].RealDate != a[j].RealDate {
		return a[i].RealDate.Before(a[j].RealDate)
	}
	return a[i].Type < a[j].Type
}
