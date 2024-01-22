package common

import (
	"encoding/json"
	"net/http"
)

type H struct {
	Code    int
	Message string
	Data    interface{}
	Rows    interface{}
	Total   interface{}
}

func Resp(w http.ResponseWriter, code int, data interface{}, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	h := H{
		Code:    code,
		Data:    data,
		Message: message,
	}
	ret, err := json.Marshal(h)
	if err != nil {

	}
	_, _ = w.Write(ret)
}

func RespList(w http.ResponseWriter, code int, data interface{}, message string, total interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	h := H{
		Code:    code,
		Data:    data,
		Message: message,
		Total:   total,
	}
	ret, err := json.Marshal(h)
	if err != nil {

	}
	_, _ = w.Write(ret)
}

func RespFail(w http.ResponseWriter, data string, message string) {
	//Resp(w, -1, data, message)

	Resp(w, -1, nil, message)
}
func RespOk(w http.ResponseWriter, data interface{}, message string) {
	Resp(w, 0, data, message)
}
func RespOkList(w http.ResponseWriter, data interface{}, message string, total interface{}) {
	RespList(w, 0, data, message, total)
}
