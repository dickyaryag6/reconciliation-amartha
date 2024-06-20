package usecase

import (
	"amartha-test/entities/transactions"
	libError "amartha-test/errors"
	"context"
	"fmt"
	"mime/multipart"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/gocarina/gocsv"
)

const (
	dateTimeFormat = "02/01/2006 15:04:05"
	dateFormat     = "02/01/2006"
)

type TransactionUsecase struct {
}

func NewTransactionUsecase(usecase TransactionUsecase) TransactionUsecase {
	return usecase
}

// ConvertCurrencyToFloat converts a currency string like "-Rp1,000,000", "Rp1,000,000", "$1,000,000" or "€1,000,000" to float64
func convertCurrencyToFloat(currency string) (float64, error) {
	// Define a regex pattern to match any currency symbols and whitespace
	regexPattern := regexp.MustCompile(`[^\d.-]`)

	// Remove all matches of the pattern
	cleanedCurrency := regexPattern.ReplaceAllString(currency, "")

	// Convert the resulting string to float64
	value, err := strconv.ParseFloat(cleanedCurrency, 64)
	if err != nil {
		return 0, err
	}

	return value, nil
}

func unmarshalCsvToStructForBankStatements(file *multipart.File) (result []*transactions.BankStatements, err error) {
	err = gocsv.UnmarshalMultipartFile(file, &result)
	if err != nil {
		return nil, err
	}
	return
}


func unmarshalCsvToStructForSystemTransactions(file *multipart.File) (result []*transactions.SystemTransactions, err error) {
	err = gocsv.UnmarshalMultipartFile(file, &result)
	if err != nil {
		return nil, err
	}
	return
}

// function to validate bank statement data
func validateBankStatementsData(data []*transactions.BankStatements) (err error) {
	for index, d := range data {
		// convert string with currency to real amount
		data[index].RealAmount, err = convertCurrencyToFloat(d.Amount)

		if err != nil {
			return libError.NewBadRequestError("amount format in bank statements data is invalid")
		}

		var timeParsed time.Time

		// convert time in string to time format
		timeParsed, err = time.Parse(dateFormat, d.Date)
		if err != nil {
			return libError.NewBadRequestError(fmt.Sprintf("date format in bank statements data is invalid, use this format %s", dateFormat))
		}
		data[index].RealDate = time.Date(timeParsed.Year(), timeParsed.Month(), timeParsed.Day(), 0, 0, 0, 0, time.Local)

		// get transaction type
		if data[index].RealAmount < 0 {
			data[index].Type = transactions.DEBIT
		} else {
			data[index].Type = transactions.CREDIT
		}

		dataIDData := strings.Split(d.ID, "_")
		if len(dataIDData) > 0 {
			data[index].BankSource = dataIDData[0]
		} else {
			return libError.NewBadRequestError("unique_identifier data in bank statements is invalid")
		}

	}

	return
}

// function to validate system transaction data
func validateSystemTransactionsData(data []*transactions.SystemTransactions) (err error) {
	for index, d := range data {
		// convert string with currency to real amount
		data[index].RealAmount, err = convertCurrencyToFloat(d.Amount)
		if err != nil {
			return libError.NewBadRequestError("amount format in system transaction data is invalid")
		}

		var timeParsed time.Time

		// convert time in string to time format
		timeParsed, err = time.Parse(dateTimeFormat, d.TransactionTime)
		if err != nil {
			return libError.NewBadRequestError(fmt.Sprintf("date format in system transaction data is invalid, use this format %s", dateTimeFormat))
		}
		data[index].RealTransactionTime = time.Date(timeParsed.Year(), timeParsed.Month(), timeParsed.Day(), 0, 0, 0, 0, time.Local)

	}

	return
}

// findBankStatementData performs a binary search on the sorted bank statements array with the same date and type
func findBankStatementData(statements []*transactions.BankStatements, targetDate time.Time, targetType int) *transactions.BankStatements {
	left, right := 0, len(statements)-1

	for left <= right {
		mid := (left + right) / 2
		midStatement := statements[mid]

		if midStatement.RealDate.Equal(targetDate) && midStatement.Type == targetType {
			return midStatement
		}

		if midStatement.RealDate.Before(targetDate) || (midStatement.RealDate.Equal(targetDate) && midStatement.Type < targetType) {
			left = mid + 1
		} else {
			right = mid - 1
		}
	}

	return nil
}

// findSystemTransactionData performs a binary search on the sorted system transaction array with the same date and type
func findSystemTransactionData(transaction []*transactions.SystemTransactions, targetDate time.Time, targetType int) *transactions.SystemTransactions {
	left, right := 0, len(transaction)-1

	for left <= right {
		mid := (left + right) / 2
		midStatement := transaction[mid]

		if midStatement.RealTransactionTime.Equal(targetDate) && midStatement.Type == targetType {
			return midStatement
		}

		if midStatement.RealTransactionTime.Before(targetDate) || (midStatement.RealTransactionTime.Equal(targetDate) && midStatement.Type < targetType) {
			left = mid + 1
		} else {
			right = mid - 1
		}
	}

	return nil
}

// usecase function to do reconciliation
func (usecase TransactionUsecase) DoReconciliation(ctx context.Context, param transactions.DoReconciliationRequest) (result transactions.DoReconciliationResponse, err error) {

	// bank statements
	bankStatementsData, err := unmarshalCsvToStructForBankStatements(&param.BankStatements)
	if err != nil {
		return result, err
	}
	if len(bankStatementsData) <= 0 {
		return result, libError.NewBadRequestError("bank statements data is empty")
	}
	err = validateBankStatementsData(bankStatementsData)
	if err != nil {

		return result, err
	}

	sort.Sort(transactions.SortByRealDateBankStatement(bankStatementsData))

	// system transaction
	systemTransactionsData, err := unmarshalCsvToStructForSystemTransactions(&param.SystemTransactions)
	if err != nil {
		return result, err
	}
	if len(systemTransactionsData) <= 0 {
		return result, libError.NewBadRequestError("system transactions data is empty")
	}
	err = validateSystemTransactionsData(systemTransactionsData)
	if err != nil {
		return result, err
	}

	sort.Sort(transactions.SortByRealDateSystemTransaction(systemTransactionsData))

	var totalMatchedBankStatements float64 = 0
	var totalMatchedSystemTransactions float64 = 0

	// map for grouping missing bank statements data to each bank group
	missingBankStatements := make(map[string][]transactions.BankStatements)


	matchedBankStatements := map[string]bool{}
	unmatchedBankStatements := 0
	unmatchedSystemTransactions := 0

	// compare bank statements to system transaction
	for _, bankStatement := range bankStatementsData {
		systemTransactionData := findSystemTransactionData(systemTransactionsData, bankStatement.RealDate, bankStatement.Type)
		if systemTransactionData == nil {
			unmatchedBankStatements += 1
			if _, ok := missingBankStatements[bankStatement.BankSource]; !ok {
				missingBankStatements[bankStatement.BankSource] = []transactions.BankStatements{
					*bankStatement,
				}
			} else {
				missingBankStatements[bankStatement.BankSource] = append(missingBankStatements[bankStatement.BankSource], 
					*bankStatement,
				)
			}

		} else {
			totalMatchedBankStatements += bankStatement.RealAmount
			result.MatchedTransaction += 1

			if _, ok := matchedBankStatements[bankStatement.ID]; !ok {
				matchedBankStatements[bankStatement.ID] = true
			}
		}
	}

	// compare system transactions to bank statements
	for _, systemTransaction := range systemTransactionsData {
		result.TransactionsProceed += 1

		bankStatementData := findBankStatementData(bankStatementsData, systemTransaction.RealTransactionTime, systemTransaction.Type)
		if bankStatementData == nil {
			unmatchedSystemTransactions += 1
			result.MissingSystemTransactions = append(result.MissingSystemTransactions, *systemTransaction)
		} else {
			if _, ok := matchedBankStatements[bankStatementData.ID]; ok {
				continue
			}
			totalMatchedSystemTransactions += systemTransaction.RealAmount
			result.MatchedTransaction += 1
		}
	}

	// combine unmatched from bank statements and system transactions
	result.UnmatchedTransaction = unmatchedBankStatements+unmatchedSystemTransactions

	result.TotalDiscrepancies = totalMatchedBankStatements - totalMatchedSystemTransactions
	if result.TotalDiscrepancies < 0 {
		result.TotalDiscrepancies *= -1
	}

	result.MissingBankStatements = missingBankStatements

	return
}
