package routes

import (
	"encoding/json"
	"io"
	"net/http"
)

func bodyToStruct(body io.ReadCloser, s interface{}) error {
	bodyBytes, err := io.ReadAll(body)
	if err != nil {
		return err
	}

	err = json.Unmarshal(bodyBytes, s)
	if err != nil {
		return err
	}

	return nil
}

func structToString(s interface{}) []byte {
	bodyBytes, err := json.Marshal(s)
	if err != nil {
		return nil
	}

	return bodyBytes
}

func writeStructToResponse(res http.ResponseWriter, s interface{}, code int) {
	var str []byte

	if s == nil {
		str = []byte("{}")
	} else {
		str = structToString(s)
	}

	res.Header().Set("Content-Type", "application/json")

	res.WriteHeader(code)
	res.Write(str)
}

type errorResponse struct {
	Reason string `json:"reason"`
	Code   int    `json:"code"`
}

func handleError(res http.ResponseWriter, errmsg string, code int) {
	resBody := errorResponse{
		Code:   code,
		Reason: errmsg,
	}

	writeStructToResponse(res, resBody, code)
}

func handleUnauthenticated(res http.ResponseWriter) {
	handleError(res, "User is not authenticated", http.StatusUnauthorized)
}
