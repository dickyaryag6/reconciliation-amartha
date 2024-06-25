package errors

import (
	"amartha-test/response"
	"net/http"
)

type ErrorMessage struct {
	ErrorDescription string `json:"error_description"`
}

func NewErrorMessage(status int, err error) *ErrorMessage {
	return &ErrorMessage{
		ErrorDescription: err.Error(),
	}
}

func (e *ErrorMessage) Error() string {
	return e.ErrorDescription
}

func SetInternalServerErrorForHandler(w http.ResponseWriter, errValue error) (err error) {

	_, err = response.WriteJSONResponse(w, http.StatusInternalServerError, &ErrorMessage{
		ErrorDescription: errValue.Error(),
	})

	return
}

func SetError(w http.ResponseWriter, errValue interface{}) (err error) {
	if errType, ok := errValue.(*ErrorMessage); ok {
		return SetBadRequestErrorForHandler(w, errType.ErrorDescription)
	} else if errType, ok := errValue.(error); ok {
		return SetInternalServerErrorForHandler(w, errType)
	}

	return

}

func SetBadRequestErrorForHandler(w http.ResponseWriter, errValue string) (err error) {
	_, err = response.WriteJSONResponse(w, http.StatusBadRequest, &ErrorMessage{
		ErrorDescription: errValue,
	})

	return
}

func NewBadRequestError(errValue string) *ErrorMessage {
	return &ErrorMessage{
		ErrorDescription: errValue,
	}
}
