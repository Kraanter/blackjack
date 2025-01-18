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
	str := structToString(s)

	res.Write(str)
	res.Header().Set("Content-Type", "application/json")

	res.WriteHeader(code)
}

type errorResponse struct {
	reason string
	code   int
}

func handleError(res http.ResponseWriter, err error, code int) {
	resBody := errorResponse{
		code:   code,
		reason: err.Error(),
	}

	structToString(resBody)

	res.WriteHeader(code)
}
