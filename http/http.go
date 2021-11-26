package http

import (
	"encoding/json"
	"log"
	"net/http"

	dq "doublequote"
)

type Response struct {
	Data interface{} `json:"data,omitempty"`
}

func ValidationError(w http.ResponseWriter, errs map[string][]string) {
	var invalids []dq.InvalidParam

	for field, msg := range errs {
		invalids = append(invalids, dq.InvalidParam{
			Name:   field,
			Reason: msg[0],
		})
	}

	p := dq.NewProblem("Validation failed.")
	p.InvalidParams = invalids

	w.WriteHeader(ErrorStatusCode(dq.EINVALID))
	json.NewEncoder(w).Encode(&p)
}

func Error(w http.ResponseWriter, r *http.Request, err error) {
	// Extract error code & message.
	code, message := dq.ErrorCode(err), dq.ErrorMessage(err)

	// Log & report internal errors.
	if code == dq.EINTERNAL {
		dq.ReportError(r.Context(), err, r)
		LogError(r, err)
	}

	w.WriteHeader(ErrorStatusCode(code))
	json.NewEncoder(w).Encode(dq.NewProblem(message))
}

func LogError(r *http.Request, err error) {
	log.Printf("http: ERROR: %s %s: %s", r.Method, r.URL.Path, err)
}

// Application error codes to HTTP status codes.
var codes = map[string]int{
	dq.ECONFLICT:       http.StatusConflict,
	dq.EINVALID:        http.StatusBadRequest,
	dq.ENOTFOUND:       http.StatusNotFound,
	dq.ENOTIMPLEMENTED: http.StatusNotImplemented,
	dq.EUNAUTHORIZED:   http.StatusUnauthorized,
	dq.EINTERNAL:       http.StatusInternalServerError,
}

// ErrorStatusCode returns the associated HTTP status code for a WTF error code.
func ErrorStatusCode(code string) int {
	if v, ok := codes[code]; ok {
		return v
	}
	return http.StatusInternalServerError
}

func sendJSON(w http.ResponseWriter, r *http.Request, resp interface{}) {
	err := json.NewEncoder(w).Encode(&Response{Data: resp})
	if err != nil {
		Error(w, r, err)
	}
}
