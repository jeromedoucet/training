package response

import (
	"encoding/json"
	"net/http"
)

type Error struct {
	Message string `json:"message"`
}

func RenderError(statusCode int, msg string, w http.ResponseWriter) {
	res := Error{Message: msg}
	body, _ := json.Marshal(res)
	w.WriteHeader(statusCode)
	w.Write(body)
}
