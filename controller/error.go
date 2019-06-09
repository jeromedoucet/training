package controller

import (
	"encoding/json"
	"net/http"

	"github.com/jeromedoucet/training/controller/response"
)

func renderError(statusCode int, msg string, w http.ResponseWriter) {
	res := response.Error{Message: msg}
	body, _ := json.Marshal(res)
	w.WriteHeader(statusCode)
	w.Write(body)
}
