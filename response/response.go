package response

import (
	"encoding/json"
	"net/http"
)

func WriteJSONResponse(w http.ResponseWriter, status int, data interface{}) (int, error) {
	w.Header().Set("Content-Type", "application/json")
	byteData, err := json.Marshal(data)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		length, err := w.Write([]byte(`{"errors":["Internal Server Error"]}`))
		if err != nil {
			return length, err
		}
		return length, err
	}

	w.WriteHeader(status)
	return w.Write(byteData)
}

func SetOK(w http.ResponseWriter, data interface{}) (err error) {
	_, err = WriteJSONResponse(w, http.StatusOK, data)
	return
}