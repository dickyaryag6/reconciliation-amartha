package errors

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
)

func TestNewBadRequestError(t *testing.T) {
	type args struct {
		errValue string
	}
	tests := []struct {
		name string
		args args
		want *ErrorMessage
	}{
		{
			name: "Succesful",
			args: args{
				errValue: "error",
			},
			want: &ErrorMessage{
				ErrorDescription: "error",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewBadRequestError(tt.args.errValue); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewBadRequestError() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSetBadRequestErrorForHandler(t *testing.T) {
	type args struct {
		w        http.ResponseWriter
		errValue string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "Succesful",
			args: args{
				w:        httptest.NewRecorder(),
				errValue: "error",
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := SetBadRequestErrorForHandler(tt.args.w, tt.args.errValue); (err != nil) != tt.wantErr {
				t.Errorf("SetBadRequestErrorForHandler() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestSetError(t *testing.T) {
	type args struct {
		w        http.ResponseWriter
		errValue interface{}
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "Succesful Error Message",
			args: args{
				w:        httptest.NewRecorder(),
				errValue: NewBadRequestError("error"),
			},
			wantErr: false,
		},
		{
			name: "Succesful Error",
			args: args{
				w:        httptest.NewRecorder(),
				errValue: errors.New("error"),
			},
			wantErr: false,
		},
		{
			name: "Error Not Recognized",
			args: args{
				w: httptest.NewRecorder(),
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := SetError(tt.args.w, tt.args.errValue); (err != nil) != tt.wantErr {
				t.Errorf("SetError() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestErrorMessage_Error(t *testing.T) {
	type fields struct {
		ErrorDescription string
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		{
			name:   "Succesful",
			fields: fields{},
			want:   "",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := &ErrorMessage{
				ErrorDescription: tt.fields.ErrorDescription,
			}
			if got := e.Error(); got != tt.want {
				t.Errorf("ErrorMessage.Error() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNewErrorMessage(t *testing.T) {
	type args struct {
		status int
		err    error
	}
	tests := []struct {
		name string
		args args
		want *ErrorMessage
	}{
		{
			name: "Succesful",
			args: args{
				status: http.StatusOK,
				err:    errors.New("error"),
			},
			want: &ErrorMessage{
				ErrorDescription: "error",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewErrorMessage(tt.args.status, tt.args.err); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewErrorMessage() = %v, want %v", got, tt.want)
			}
		})
	}
}
