package usecase

import (
	"amartha-test/entities/transactions"
	"bytes"
	"context"
	"errors"
	"mime/multipart"
	"reflect"
	"testing"
	"time"

	"github.com/gocarina/gocsv"
)

var (
	errMock = errors.New("err")
)

func TestNewTransactionUsecase(t *testing.T) {
	type args struct {
		usecase TransactionUsecase
	}
	tests := []struct {
		name string
		args args
		want TransactionUsecase
	}{
		{
			name: "Succesful",
			args: args{
				usecase: TransactionUsecase{},
			},
			want: TransactionUsecase{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewTransactionUsecase(tt.args.usecase); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewTransactionUsecase() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_convertCurrencyToFloat(t *testing.T) {
	type args struct {
		currency string
	}
	tests := []struct {
		name    string
		args    args
		want    float64
		wantErr bool
	}{
		{
			name: "Succesful",
			args: args{
				currency: "Rp100,000",
			},
			want:    100000,
			wantErr: false,
		},
		{
			name: "Failed",
			args: args{
				currency: "abc",
			},
			want:    0,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := convertCurrencyToFloat(tt.args.currency)
			if (err != nil) != tt.wantErr {
				t.Errorf("convertCurrencyToFloat() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("convertCurrencyToFloat() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_unmarshalCsvToStructForBankStatements(t *testing.T) {

	csvData := `unique_identifier,amount,date 
BCA_12345,"Rp1,500,000",01/01/2024`
	var buf bytes.Buffer
	writer := multipart.NewWriter(&buf)

	part, err := writer.CreateFormFile("file", "test.csv")
	if err != nil {
		t.Fatalf("Failed to create form file: %v", err)
	}

	_, err = part.Write([]byte(csvData))
	if err != nil {
		t.Fatalf("Failed to write to form file: %v", err)
	}
	err = writer.Close()
	if err != nil {
		t.Fatalf("Failed to close writer: %v", err)
	}

	reader := multipart.NewReader(&buf, writer.Boundary())
	form, err := reader.ReadForm(1024)
	if err != nil {
		t.Fatalf("Failed to read form: %v", err)
	}
	defer form.RemoveAll()
	file := form.File["file"][0]
	openedFile, err := file.Open()
	if err != nil {
		t.Fatalf("Failed to open form file: %v", err)
	}
	defer openedFile.Close()

	multipartFile := multipart.File(openedFile)

	type args struct {
		file *multipart.File
	}
	tests := []struct {
		name       string
		args       args
		wantResult []*transactions.BankStatements
		wantErr    bool
		mock       func()
		unmock     func()
	}{
		{
			name: "Succesful",
			args: args{
				file: &multipartFile,
			},
			wantResult: []*transactions.BankStatements{
				{
					ID:     "BCA_12345",
					Amount: "Rp1,500,000",
					Date:   "01/01/2024",
				},
			},
			wantErr: false,
			mock:    func() {},
			unmock:  func() {},
		},
		{
			name: "Failed",
			args: args{
				file: nil,
			},
			wantResult: nil,
			wantErr:    true,
			mock: func() {
				gocsvUnmarshalMultipartFile = func(in *multipart.File, out interface{}) error {
					return errMock
				}
			},
			unmock: func() {
				gocsvUnmarshalMultipartFile = gocsv.UnmarshalMultipartFile
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mock()
			gotResult, err := unmarshalCsvToStructForBankStatements(tt.args.file)
			if (err != nil) != tt.wantErr {
				t.Errorf("unmarshalCsvToStructForBankStatements() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotResult, tt.wantResult) {
				t.Errorf("unmarshalCsvToStructForBankStatements() = %v, want %v", gotResult, tt.wantResult)
			}
			tt.unmock()
		})
	}
}

func Test_unmarshalCsvToStructForSystemTransactions(t *testing.T) {

	csvData := `trxID,amount,type,transactionTime
1,"Rp8,500,000",2,01/01/2024 8:45:00`
	var buf bytes.Buffer
	writer := multipart.NewWriter(&buf)

	part, err := writer.CreateFormFile("file", "test.csv")
	if err != nil {
		t.Fatalf("Failed to create form file: %v", err)
	}

	_, err = part.Write([]byte(csvData))
	if err != nil {
		t.Fatalf("Failed to write to form file: %v", err)
	}
	err = writer.Close()
	if err != nil {
		t.Fatalf("Failed to close writer: %v", err)
	}

	reader := multipart.NewReader(&buf, writer.Boundary())
	form, err := reader.ReadForm(1024)
	if err != nil {
		t.Fatalf("Failed to read form: %v", err)
	}
	defer form.RemoveAll()
	file := form.File["file"][0]
	openedFile, err := file.Open()
	if err != nil {
		t.Fatalf("Failed to open form file: %v", err)
	}
	defer openedFile.Close()

	multipartFile := multipart.File(openedFile)

	type args struct {
		file *multipart.File
	}
	tests := []struct {
		name       string
		args       args
		wantResult []*transactions.SystemTransactions
		wantErr    bool
		mock       func()
		unmock     func()
	}{
		{
			name: "Succesful",
			args: args{
				file: &multipartFile,
			},
			wantResult: []*transactions.SystemTransactions{
				{
					TransactionID:   "1",
					Amount:          "Rp8,500,000",
					Type:            2,
					TransactionTime: "01/01/2024 8:45:00",
				},
			},
			wantErr: false,
			mock:    func() {},
			unmock:  func() {},
		},
		{
			name: "Failed",
			args: args{
				file: nil,
			},
			wantResult: nil,
			wantErr:    true,
			mock: func() {
				gocsvUnmarshalMultipartFile = func(in *multipart.File, out interface{}) error {
					return errMock
				}
			},
			unmock: func() {
				gocsvUnmarshalMultipartFile = gocsv.UnmarshalMultipartFile
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mock()
			gotResult, err := unmarshalCsvToStructForSystemTransactions(tt.args.file)
			if (err != nil) != tt.wantErr {
				t.Errorf("unmarshalCsvToStructForSystemTransactions() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotResult, tt.wantResult) {
				t.Errorf("unmarshalCsvToStructForSystemTransactions() = %v, want %v", gotResult, tt.wantResult)
			}
			tt.unmock()
		})
	}
}

func Test_validateBankStatementsData(t *testing.T) {
	type args struct {
		data []*transactions.BankStatements
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
		mock    func()
	}{
		{
			name: "Succesful DEBIT",
			args: args{
				data: []*transactions.BankStatements{
					{
						ID:     "BCA_123",
						Amount: "-Rp100,000",
						Date:   "13/01/2024",
					},
				},
			},
			wantErr: false,
			mock:    func() {},
		},
		{
			name: "Succesful CREDIT",
			args: args{
				data: []*transactions.BankStatements{
					{
						ID:     "BCA_123",
						Amount: "Rp100,000",
						Date:   "13/01/2024",
					},
				},
			},
			wantErr: false,
			mock:    func() {},
		},
		{
			name: "convertCurrencyToFloat return error",
			args: args{
				data: []*transactions.BankStatements{
					{
						ID:     "BCA_123",
						Amount: "abc",
						Date:   "13/01/2024",
					},
				},
			},
			wantErr: true,
			mock:    func() {},
		},
		{
			name: "Time Parsing return error",
			args: args{
				data: []*transactions.BankStatements{
					{
						ID:     "BCA_123",
						Amount: "Rp100,000",
						Date:   "",
					},
				},
			},
			wantErr: true,
			mock:    func() {},
		},
		{
			name: "ID doesn't contain bank source",
			args: args{
				data: []*transactions.BankStatements{
					{
						ID:     "",
						Amount: "Rp100,000",
						Date:   "13/01/2024",
					},
				},
			},
			wantErr: true,
			mock:    func() {},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mock()
			if err := validateBankStatementsData(tt.args.data); (err != nil) != tt.wantErr {
				t.Errorf("validateBankStatementsData() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_validateSystemTransactionsData(t *testing.T) {
	type args struct {
		data []*transactions.SystemTransactions
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "Succesful",
			args: args{
				data: []*transactions.SystemTransactions{
					{
						TransactionID:   "1",
						Amount:          "Rp100,000",
						TransactionTime: "15/02/2024 8:20:00",
					},
				},
			},
			wantErr: false,
		},
		{
			name: "convertCurrencyToFloat return error",
			args: args{
				data: []*transactions.SystemTransactions{
					{
						TransactionID:   "1",
						Amount:          "",
						TransactionTime: "15/02/2024 8:20:00",
					},
				},
			},
			wantErr: true,
		},
		{
			name: "Time parsing return error",
			args: args{
				data: []*transactions.SystemTransactions{
					{
						TransactionID:   "1",
						Amount:          "Rp100,000",
						TransactionTime: "",
					},
				},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := validateSystemTransactionsData(tt.args.data); (err != nil) != tt.wantErr {
				t.Errorf("validateSystemTransactionsData() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_findBankStatementData(t *testing.T) {
	type args struct {
		statements []*transactions.BankStatements
		targetDate time.Time
		targetType int
	}
	tests := []struct {
		name string
		args args
		want *transactions.BankStatements
	}{
		{
			name: "Succesful",
			args: args{
				statements: []*transactions.BankStatements{
					{
						ID:       "MANDIRI_12346",
						Amount:   "Rp2,500,000",
						Date:     "13/01/2024",
						Type:     transactions.CREDIT,
						RealDate: time.Date(2024, time.Month(1), 13, 0, 0, 0, 0, time.Local),
					},
					{
						ID:       "MANDIRI_12347",
						Amount:   "Rp2,500,000",
						Date:     "15/01/2024",
						Type:     transactions.CREDIT,
						RealDate: time.Date(2024, time.Month(1), 15, 0, 0, 0, 0, time.Local),
					},
					{
						ID:       "MANDIRI_12348",
						Amount:   "Rp2,500,000",
						Date:     "16/01/2024",
						Type:     transactions.CREDIT,
						RealDate: time.Date(2024, time.Month(1), 16, 0, 0, 0, 0, time.Local),
					},
				},
				targetDate: time.Date(2024, time.Month(1), 13, 0, 0, 0, 0, time.Local),
				targetType: transactions.CREDIT,
			},
			want: &transactions.BankStatements{
				ID:       "MANDIRI_12346",
				Amount:   "Rp2,500,000",
				Date:     "13/01/2024",
				Type:     transactions.CREDIT,
				RealDate: time.Date(2024, time.Month(1), 13, 0, 0, 0, 0, time.Local),
			},
		},
		{
			name: "Not Found",
			args: args{
				statements: []*transactions.BankStatements{
					{
						ID:       "MANDIRI_12346",
						Amount:   "Rp2,500,000",
						Date:     "13/01/2024",
						Type:     transactions.CREDIT,
						RealDate: time.Date(2024, time.Month(1), 13, 0, 0, 0, 0, time.Local),
					},
				},
				targetDate: time.Date(2024, time.Month(1), 19, 0, 0, 0, 0, time.Local),
				targetType: transactions.CREDIT,
			},
			want: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := findBankStatementData(tt.args.statements, tt.args.targetDate, tt.args.targetType); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("findBankStatementData() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_findSystemTransactionData(t *testing.T) {
	type args struct {
		transaction []*transactions.SystemTransactions
		targetDate  time.Time
		targetType  int
	}
	tests := []struct {
		name string
		args args
		want *transactions.SystemTransactions
	}{
		{
			name: "Succesful",
			args: args{
				transaction: []*transactions.SystemTransactions{
					{
						TransactionID:       "1",
						Amount:              "Rp2,000,000",
						RealAmount:          2000000,
						Type:                transactions.DEBIT,
						TransactionTime:     "15/02/2024 8:20:00",
						RealTransactionTime: time.Date(2024, time.Month(2), 15, 8, 20, 0, 0, time.Local),
					},
					{
						TransactionID:       "2",
						Amount:              "Rp2,000,000",
						RealAmount:          2000000,
						Type:                transactions.DEBIT,
						TransactionTime:     "16/02/2024 8:20:00",
						RealTransactionTime: time.Date(2024, time.Month(2), 16, 8, 20, 0, 0, time.Local),
					},
					{
						TransactionID:       "3",
						Amount:              "Rp2,000,000",
						RealAmount:          2000000,
						Type:                transactions.DEBIT,
						TransactionTime:     "17/02/2024 8:20:00",
						RealTransactionTime: time.Date(2024, time.Month(2), 17, 8, 20, 0, 0, time.Local),
					},
				},
				targetDate: time.Date(2024, time.Month(2), 15, 8, 20, 0, 0, time.Local),
				targetType: transactions.DEBIT,
			},
			want: &transactions.SystemTransactions{
				TransactionID:       "1",
				Amount:              "Rp2,000,000",
				RealAmount:          2000000,
				Type:                transactions.DEBIT,
				TransactionTime:     "15/02/2024 8:20:00",
				RealTransactionTime: time.Date(2024, time.Month(2), 15, 8, 20, 0, 0, time.Local),
			},
		},
		{
			name: "Not Found",
			args: args{
				transaction: []*transactions.SystemTransactions{
					{
						TransactionID:       "1",
						Amount:              "Rp2,000,000",
						RealAmount:          2000000,
						Type:                transactions.DEBIT,
						TransactionTime:     "15/02/2024 8:20:00",
						RealTransactionTime: time.Date(2024, time.Month(2), 15, 8, 20, 0, 0, time.Local),
					},
					{
						TransactionID:       "2",
						Amount:              "Rp2,000,000",
						RealAmount:          2000000,
						Type:                transactions.DEBIT,
						TransactionTime:     "16/02/2024 8:20:00",
						RealTransactionTime: time.Date(2024, time.Month(2), 16, 8, 20, 0, 0, time.Local),
					},
				},
				targetDate: time.Date(2024, time.Month(2), 19, 8, 20, 0, 0, time.Local),
				targetType: transactions.DEBIT,
			},
			want: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := findSystemTransactionData(tt.args.transaction, tt.args.targetDate, tt.args.targetType); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("findSystemTransactionData() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestTransactionUsecase_DoReconciliation(t *testing.T) {
	type args struct {
		ctx   context.Context
		param transactions.DoReconciliationRequest
	}
	tests := []struct {
		name       string
		usecase    TransactionUsecase
		args       args
		wantResult transactions.DoReconciliationResponse
		wantErr    bool
		mock       func()
		unmock     func()
	}{
		{
			name:    "Succesful",
			usecase: TransactionUsecase{},
			args: args{
				param: transactions.DoReconciliationRequest{},
			},
			wantResult: transactions.DoReconciliationResponse{
				TransactionsProceed:  3,
				MatchedTransaction:   2,
				UnmatchedTransaction: 3,
				MissingBankStatements: map[string][]transactions.BankStatements{
					"MANDIRI": {
						{
							ID:         "MANDIRI_12346",
							Amount:     "Rp2,500,000",
							RealAmount: 2500000,
							Date:       "15/01/2024",
							RealDate:   time.Date(2024, time.Month(1), 15, 0, 0, 0, 0, time.Local),
							BankSource: "MANDIRI",
							Type:       transactions.CREDIT,
						},
						// MANDIRI_12347 Rp2,500,000 2.5e+06 19/01/2024 2024-01-19 00:00:00 +0700 WIB MANDIRI 2}
						{
							ID:         "MANDIRI_12347",
							Amount:     "Rp2,500,000",
							RealAmount: 2500000,
							Date:       "19/01/2024",
							RealDate:   time.Date(2024, time.Month(1), 19, 0, 0, 0, 0, time.Local),
							BankSource: "MANDIRI",
							Type:       transactions.CREDIT,
						},
					},
				},
				MissingSystemTransactions: []transactions.SystemTransactions{
					{
						TransactionID:       "11",
						Amount:              "Rp2,000,000",
						RealAmount:          2000000,
						Type:                transactions.CREDIT,
						TransactionTime:     "14/01/2024 08:20:00",
						RealTransactionTime: time.Date(2024, time.Month(1), 14, 0, 0, 0, 0, time.Local),
					},
				},
				TotalDiscrepancies: 4000000,
			},
			wantErr: false,
			mock: func() {
				unmarshalCsvToStructForBankStatements = func(_ *multipart.File) (result []*transactions.BankStatements, err error) {
					return []*transactions.BankStatements{
						{
							ID:         "MANDIRI_12346",
							Amount:     "Rp2,500,000",
							RealAmount: 2500000,
							Date:       "15/01/2024",
							Type:       transactions.CREDIT,
							BankSource: "MANDIRI",
							RealDate:   time.Date(2024, time.Month(1), 13, 0, 0, 0, 0, time.Local),
						},
						{
							ID:         "MANDIRI_12347",
							Amount:     "Rp2,500,000",
							RealAmount: 2500000,
							Date:       "19/01/2024",
							Type:       transactions.CREDIT,
							BankSource: "MANDIRI",
							RealDate:   time.Date(2024, time.Month(1), 19, 0, 0, 0, 0, time.Local),
						},
						{
							ID:         "MANDIRI_12348",
							Amount:     "Rp2,000,000",
							RealAmount: 2000000,
							Date:       "13/01/2024",
							Type:       transactions.CREDIT,
							BankSource: "MANDIRI",
							RealDate:   time.Date(2024, time.Month(1), 13, 0, 0, 0, 0, time.Local),
						},
						{
							ID:         "MANDIRI_12349",
							Amount:     "Rp2,000,000",
							RealAmount: 2000000,
							Date:       "20/01/2024",
							Type:       transactions.CREDIT,
							BankSource: "MANDIRI",
							RealDate:   time.Date(2024, time.Month(1), 20, 0, 0, 0, 0, time.Local),
						},
					}, nil
				}

				unmarshalCsvToStructForSystemTransactions = func(_ *multipart.File) (result []*transactions.SystemTransactions, err error) {
					return []*transactions.SystemTransactions{
						{
							TransactionID:       "10",
							Amount:              "Rp2,000,000",
							RealAmount:          2000000,
							Type:                transactions.CREDIT,
							TransactionTime:     "13/01/2024 08:20:00",
							RealTransactionTime: time.Date(2024, time.Month(1), 13, 0, 0, 0, 0, time.Local),
						},
						{
							TransactionID:       "11",
							Amount:              "Rp2,000,000",
							RealAmount:          2000000,
							Type:                transactions.CREDIT,
							TransactionTime:     "14/01/2024 08:20:00",
							RealTransactionTime: time.Date(2024, time.Month(1), 14, 0, 0, 0, 0, time.Local),
						},
						{
							TransactionID:       "12",
							Amount:              "Rp2,000,000",
							RealAmount:          2000000,
							Type:                transactions.CREDIT,
							TransactionTime:     "20/01/2024 08:20:00",
							RealTransactionTime: time.Date(2024, time.Month(1), 20, 0, 0, 0, 0, time.Local),
						},
					}, nil
				}

			},
			unmock: func() {},
		},
		{
			name:    "unmarshalCsvToStructForBankStatements return error",
			usecase: TransactionUsecase{},
			args: args{
				param: transactions.DoReconciliationRequest{},
			},
			wantResult: transactions.DoReconciliationResponse{},
			wantErr:    true,
			mock: func() {
				unmarshalCsvToStructForBankStatements = func(_ *multipart.File) (result []*transactions.BankStatements, err error) {
					return nil, errMock
				}
			},
			unmock: func() {},
		},
		{
			name:    "unmarshalCsvToStructForSystemTransactions return error",
			usecase: TransactionUsecase{},
			args: args{
				param: transactions.DoReconciliationRequest{},
			},
			wantResult: transactions.DoReconciliationResponse{},
			wantErr:    true,
			mock: func() {
				unmarshalCsvToStructForBankStatements = func(_ *multipart.File) (result []*transactions.BankStatements, err error) {
					return []*transactions.BankStatements{
						{
							ID:         "MANDIRI_12346",
							Amount:     "Rp2,500,000",
							RealAmount: 2500000,
							Date:       "15/01/2024",
							Type:       transactions.DEBIT,
							BankSource: "MANDIRI",
							RealDate:   time.Date(2024, time.Month(1), 13, 0, 0, 0, 0, time.Local),
						},
					}, nil
				}

				unmarshalCsvToStructForSystemTransactions = func(_ *multipart.File) (result []*transactions.SystemTransactions, err error) {
					return nil, errMock
				}
			},
			unmock: func() {},
		},
		{
			name:    "BankStatements data is empty",
			usecase: TransactionUsecase{},
			args: args{
				param: transactions.DoReconciliationRequest{},
			},
			wantResult: transactions.DoReconciliationResponse{},
			wantErr:    true,
			mock: func() {
				unmarshalCsvToStructForBankStatements = func(_ *multipart.File) (result []*transactions.BankStatements, err error) {
					return []*transactions.BankStatements{}, nil
				}

			},
			unmock: func() {},
		},
		{
			name:    "SystemTransactions data is empty",
			usecase: TransactionUsecase{},
			args: args{
				param: transactions.DoReconciliationRequest{},
			},
			wantResult: transactions.DoReconciliationResponse{},
			wantErr:    true,
			mock: func() {
				unmarshalCsvToStructForBankStatements = func(_ *multipart.File) (result []*transactions.BankStatements, err error) {
					return []*transactions.BankStatements{
						{
							ID:         "MANDIRI_12346",
							Amount:     "Rp2,500,000",
							RealAmount: 2500000,
							Date:       "15/01/2024",
							Type:       transactions.DEBIT,
							BankSource: "MANDIRI",
							RealDate:   time.Date(2024, time.Month(1), 13, 0, 0, 0, 0, time.Local),
						},
					}, nil
				}

				unmarshalCsvToStructForSystemTransactions = func(_ *multipart.File) (result []*transactions.SystemTransactions, err error) {
					return []*transactions.SystemTransactions{}, nil
				}
			},
			unmock: func() {},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			usecase := TransactionUsecase{}

			tt.mock()
			gotResult, err := usecase.DoReconciliation(tt.args.ctx, tt.args.param)
			if (err != nil) != tt.wantErr {
				t.Errorf("TransactionUsecase.DoReconciliation() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotResult, tt.wantResult) {
				t.Errorf("TransactionUsecase.DoReconciliation() = %v, want %v", gotResult, tt.wantResult)
			}
			tt.unmock()
		})
	}
}
