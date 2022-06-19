package main

import (
	"encoding/json"
	"net/http"

	"github.com/vmihailenco/msgpack/v5"
)

type ResponseError struct {
	Error string `msgpack:"error"`
}
type ResponseSuccess struct {
	Data     interface{} `msgpack:"data" json:"data"`
	HasNext  bool        `msgpack:"has_next" json:"has_next"`
	TotalRow int64       `msgpack:"total_row" json:"total_row"`
}

func errorResponse(w http.ResponseWriter, status int, message string) {
	w.WriteHeader(status)
	msgpack.NewEncoder(w).Encode(&ResponseError{message})
}

func paginationSuccessJSONResponse(w http.ResponseWriter, data interface{}, hasNext bool, totalRow int64) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(&ResponseSuccess{data, hasNext, totalRow})
}
func paginationSuccessMsgPackResponse(w http.ResponseWriter, data interface{}, hasNext bool, totalRow int64) {
	w.Header().Set("Content-Type", "application/x-msgpack")
	w.WriteHeader(http.StatusOK)
	msgpack.NewEncoder(w).Encode(&ResponseSuccess{data, hasNext, totalRow})
}
