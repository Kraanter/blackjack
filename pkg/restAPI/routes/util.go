package routes

import (
	"encoding/json"
	"fmt"
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

func writeStructToWriter(w io.Writer, s interface{}) {
	var str []byte

	if s == nil {
		str = []byte("{}")
	} else {
		str = structToString(s)
	}

	w.Write(str)
}

func writeStructToResponse(res http.ResponseWriter, s interface{}, code int) {
	res.Header().Set("Content-Type", "application/json")

	res.WriteHeader(code)
	writeStructToWriter(res, s)
}

func sendSSEvent(w http.ResponseWriter, event string, data interface{}) {
	fmt.Fprintf(w, "id: %s\nevent: %s\ndata: %s\n\n", event, event, structToString(data))
	flusher := w.(http.Flusher)
	if flusher != nil {
		flusher.Flush()
	}
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
