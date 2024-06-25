package handlers

import (
	"amartha-test/entities/transactions"
	usecaseMock "amartha-test/entities/usecases/mock"
	"bytes"
	"errors"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"net/textproto"
	"reflect"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

var (
	errMock = errors.New("")
)

func TestNewTransactionHandler(t *testing.T) {

	type args struct {
		handler TransactionHandler
	}
	tests := []struct {
		name string
		args args
		want TransactionHandler
	}{
		{
			name: "Succesful",
			args: args{
				handler: TransactionHandler{},
			},
			want: TransactionHandler{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewTransactionHandler(tt.args.handler); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewTransactionHandler() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestTransactionHandler_HandleReconciliation(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	// Mock TransactionUsecase
	mockUsecase := usecaseMock.NewMockTransactionUsecase(ctrl)

	tests := []struct {
		name         string
		mock         func()
		httpStatus   int
		r            *http.Request
		generateData func() (data bytes.Buffer, contentType string)
	}{
		{
			name: "Succesful",
			mock: func() {
				mockUsecase.EXPECT().DoReconciliation(gomock.Any(), gomock.Any()).
					Return(transactions.DoReconciliationResponse{}, nil)

			},
			httpStatus: http.StatusOK,
			generateData: func() (bytes.Buffer, string) {
				var buf bytes.Buffer
				writer := multipart.NewWriter(&buf)

				// Create form file fields
				bankStatements, err := writer.CreatePart(textproto.MIMEHeader{
					"Content-Disposition": []string{`form-data; name="bank_statements"; filename="bank.csv"`},
					"Content-Type":        []string{"text/csv"},
				})
				if err != nil {
					t.Errorf("error in creating bank_statements data")
				}
				bankStatements.Write([]byte("unique_identifier,amount,date\nBCA_12345,\"Rp1,500,000\",01/01/2024"))

				systemTransactions, err := writer.CreatePart(textproto.MIMEHeader{
					"Content-Disposition": []string{`form-data; name="system_transactions"; filename="system.csv"`},
					"Content-Type":        []string{"text/csv"},
				})
				if err != nil {
					t.Errorf("error in creating system_transactions data")
				}
				systemTransactions.Write([]byte("unique_identifier,amount,date\nSYS_12345,\"Rp1,500,000\",01/01/2024"))

				err = writer.Close()
				if err != nil {
					t.Errorf("error in writing data")
				}

				return buf, writer.FormDataContentType()
			},
		},
		{
			name: "Failed",
			mock: func() {
				mockUsecase.EXPECT().DoReconciliation(gomock.Any(), gomock.Any()).
					Return(transactions.DoReconciliationResponse{}, errMock)
			},
			generateData: func() (bytes.Buffer, string) {
				var buf bytes.Buffer
				writer := multipart.NewWriter(&buf)

				bankStatements, err := writer.CreatePart(textproto.MIMEHeader{
					"Content-Disposition": []string{`form-data; name="bank_statements"; filename="bank.csv"`},
					"Content-Type":        []string{"text/csv"},
				})
				if err != nil {
					t.Errorf("error in creating bank_statements data")
				}
				bankStatements.Write([]byte("unique_identifier,amount,date\nBCA_12345,\"Rp1,500,000\",01/01/2024"))

				systemTransactions, err := writer.CreatePart(textproto.MIMEHeader{
					"Content-Disposition": []string{`form-data; name="system_transactions"; filename="system.csv"`},
					"Content-Type":        []string{"text/csv"},
				})
				if err != nil {
					t.Errorf("error in creating system_transactions data")
				}
				systemTransactions.Write([]byte("unique_identifier,amount,date\nSYS_12345,\"Rp1,500,000\",01/01/2024"))

				err = writer.Close()
				if err != nil {
					t.Errorf("error in writing data")
				}

				return buf, writer.FormDataContentType()
			},
			httpStatus: http.StatusInternalServerError,
		},
		{
			name:       "bank_statements file is not csv",
			mock:       func() {},
			httpStatus: http.StatusBadRequest,
			generateData: func() (bytes.Buffer, string) {
				var buf bytes.Buffer
				writer := multipart.NewWriter(&buf)

				bankStatements, err := writer.CreatePart(textproto.MIMEHeader{
					"Content-Disposition": []string{`form-data; name="bank_statements"; filename="bank.txt"`},
					"Content-Type":        []string{"text/plain"},
				})
				if err != nil {
					t.Errorf("error in creating bank_statements data")
				}
				bankStatements.Write([]byte("unique_identifier,amount,date\nBCA_12345,\"Rp1,500,000\",01/01/2024"))

				systemTransactions, err := writer.CreatePart(textproto.MIMEHeader{
					"Content-Disposition": []string{`form-data; name="system_transactions"; filename="system.csv"`},
					"Content-Type":        []string{"text/csv"},
				})
				if err != nil {
					t.Errorf("error in creating system_transactions data")
				}
				systemTransactions.Write([]byte("unique_identifier,amount,date\nSYS_12345,\"Rp1,500,000\",01/01/2024"))

				err = writer.Close()
				if err != nil {
					t.Errorf("error in writing data")
				}

				return buf, writer.FormDataContentType()

			},
		},
		{
			name:       "system_transaction file is not csv",
			mock:       func() {},
			httpStatus: http.StatusBadRequest,
			generateData: func() (bytes.Buffer, string) {
				var buf bytes.Buffer
				writer := multipart.NewWriter(&buf)

				bankStatements, err := writer.CreatePart(textproto.MIMEHeader{
					"Content-Disposition": []string{`form-data; name="bank_statements"; filename="bank.csv"`},
					"Content-Type":        []string{"text/csv"},
				})
				if err != nil {
					t.Errorf("error in creating bank_statements data")
				}
				bankStatements.Write([]byte("unique_identifier,amount,date\nBCA_12345,\"Rp1,500,000\",01/01/2024"))

				systemTransactions, err := writer.CreatePart(textproto.MIMEHeader{
					"Content-Disposition": []string{`form-data; name="system_transactions"; filename="system.txt"`},
					"Content-Type":        []string{"text/plain"},
				})
				if err != nil {
					t.Errorf("error in creating system_transactions data")
				}
				systemTransactions.Write([]byte("unique_identifier,amount,date\nSYS_12345,\"Rp1,500,000\",01/01/2024"))

				err = writer.Close()
				if err != nil {
					t.Errorf("error in writing data")
				}

				return buf, writer.FormDataContentType()

			},
		},
		{
			name:       "No Bank Statement File",
			mock:       func() {},
			httpStatus: http.StatusBadRequest,
			generateData: func() (bytes.Buffer, string) {
				var buf bytes.Buffer
				writer := multipart.NewWriter(&buf)

				err := writer.Close()
				if err != nil {
					t.Errorf("error in writing data")
				}

				return buf, writer.FormDataContentType()

			},
		},
		{
			name:       "No System Transaction File",
			mock:       func() {},
			httpStatus: http.StatusBadRequest,
			generateData: func() (bytes.Buffer, string) {
				var buf bytes.Buffer
				writer := multipart.NewWriter(&buf)

				bankStatements, err := writer.CreatePart(textproto.MIMEHeader{
					"Content-Disposition": []string{`form-data; name="bank_statements"; filename="bank.csv"`},
					"Content-Type":        []string{"text/csv"},
				})
				if err != nil {
					t.Errorf("error in creating bank_statements data")
				}
				bankStatements.Write([]byte("unique_identifier,amount,date\nBCA_12345,\"Rp1,500,000\",01/01/2024"))

				err = writer.Close()
				if err != nil {
					t.Errorf("error in writing data")
				}

				return buf, writer.FormDataContentType()

			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			data, contentType := tt.generateData()

			r := httptest.NewRequest(http.MethodPost, "/reconciliation", &data)
			r.Header.Set("Content-Type", contentType)

			w := httptest.NewRecorder()
			handler := TransactionHandler{
				TransactionUsecase: mockUsecase,
			}
			tt.mock()
			handler.HandleReconciliation(w, r)
			assert.Equal(t, tt.httpStatus, w.Code)
		})
	}
}
