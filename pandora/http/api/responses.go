package api

import (
	"encoding/hex"
	"encoding/json"
	"fmt"
	"net/http"

	"golang.org/x/crypto/blake2b"
)

type ErrorResponseBody struct {
	// Note: StatusCode is the only snake cased json response we have in our entire api
	Status         string `json:"status"`
	HttpStatusCode int    `json:"httpStatusCode"`
	Message        string `json:"message"`
}

type CollectionMeta struct {
	Offset int `json:"offset"`
	Limit  int `json:"limit"`
	Total  int `json:"total"`
}

type CollectionResponse struct {
	Data []interface{}  `json:"data"`
	Meta CollectionMeta `json:"meta"`
}

var (
	UnauthorizedResponseBody = ErrorResponseBody{
		Status:         "error",
		HttpStatusCode: http.StatusUnauthorized,
		Message:        "Must be logged on to access this endpoint",
	}

	ForbiddenResponseBody = ErrorResponseBody{
		Status:         "error",
		HttpStatusCode: http.StatusForbidden,
		Message:        "You do not have access to this endpoint, this request has been logged",
	}

	MalformedAuthorizationResponseBody = ErrorResponseBody{
		Status:         "error",
		HttpStatusCode: http.StatusBadRequest,
		Message:        "An invalid authorization header was provided, expected Authorization: Bearer <token>",
	}

	InternalServerErrorResponseBody = ErrorResponseBody{
		Status:         "error",
		HttpStatusCode: http.StatusInternalServerError,
		Message:        "An internal server error occurred, please try again later",
	}
)

type ValidationErrorMap map[string]interface{}
type ValidationMetaMap map[string]interface{}

func WriteValidationErrorResponse(w http.ResponseWriter, errorMap ValidationErrorMap, metaMap ValidationMetaMap) {
	w.Header().Set("Content-Type", "application/json")

	payload := struct {
		Status         string             `json:"status"`
		HttpStatusCode int                `json:"httpStatusCode"`
		Message        string             `json:"message"`
		Errors         ValidationErrorMap `json:"errors"`
		Meta           ValidationMetaMap  `json:"meta"`
	}{
		Status:         "error",
		HttpStatusCode: http.StatusUnauthorized,
		Message:        "Valiidation Errors",
		Errors:         errorMap,
		Meta:           metaMap,
	}

	body, err := json.Marshal(payload)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusUnprocessableEntity)
	w.Write(body)
}

func WriteForbiddenError(w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/json")

	body, err := json.Marshal(&ForbiddenResponseBody)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusForbidden)
	w.Write(body)
}

func WriteOKResponse(w http.ResponseWriter, payload interface{}) {
	w.Header().Set("Content-Type", "application/json")

	body, err := json.Marshal(payload)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(fmt.Sprintf("Failed to create response: %s", err)))
		return
	}

	w.Write(body)
}

func WriteOkCachedResponse(w http.ResponseWriter, r *http.Request, payload interface{}) {

	body, err := json.Marshal(payload)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(fmt.Sprintf("Failed to create response: %s", err)))
		return
	}

	checksum := blake2b.Sum256(body)
	eTagHeader := fmt.Sprintf("W/\"%s\"", hex.EncodeToString(checksum[:]))

	w.Header().Set("ETag", eTagHeader)

	if r.Header.Get("If-None-Match") == eTagHeader {
		w.WriteHeader(http.StatusNotModified)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(body)
}

func WriteCreatedResponse(w http.ResponseWriter, payload interface{}) {
	w.Header().Set("Content-Type", "application/json")

	body, err := json.Marshal(payload)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(fmt.Sprintf("Failed to create response: %s", err)))
		return
	}

	w.WriteHeader(http.StatusCreated)
	w.Write(body)
}

func WriteNotFoundResponse(w http.ResponseWriter, detail string, args ...interface{}) {
	w.Header().Set("Content-Type", "application/json")

	payload := struct {
		Status         string `json:"status"`
		HttpStatusCode int    `json:"httpStatusCode"`
		Message        string `json:"message"`
	}{
		Status:         "error",
		HttpStatusCode: http.StatusNotFound,
		Message:        "Not Found: " + detail,
	}

	body, err := json.Marshal(&payload)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(fmt.Sprintf("Failed to create response: %s", err)))
		return
	}

	w.WriteHeader(http.StatusNotFound)
	w.Write(body)
}

func WriteInternalServerErrorResponse(w http.ResponseWriter, detail string, args ...interface{}) {
	w.Header().Set("Content-Type", "application/json")

	payload := struct {
		Status         string `json:"status"`
		HttpStatusCode int    `json:"httpStatusCode"`
		Message        string `json:"message"`
	}{
		Status:         "error",
		HttpStatusCode: http.StatusInternalServerError,
		Message:        "Internal Server Error: " + detail,
	}

	body, err := json.Marshal(&payload)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(fmt.Sprintf("Failed to create response: %s", err)))
		return
	}

	w.WriteHeader(http.StatusInternalServerError)
	w.Write(body)
}

func WriteBadRequestResponse(w http.ResponseWriter, detail string, args ...interface{}) {
	w.Header().Set("Content-Type", "application/json")

	payload := struct {
		Status         string `json:"status"`
		HttpStatusCode int    `json:"httpStatusCode"`
		Message        string `json:"message"`
	}{
		Status:         "error",
		HttpStatusCode: http.StatusBadRequest,
		Message:        "Bad request: " + detail,
	}

	body, err := json.Marshal(&payload)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(fmt.Sprintf("Failed to create response: %s", err)))
		return
	}

	w.WriteHeader(http.StatusBadRequest)
	w.Write(body)
}

func WriteUnauthorizedResponse(w http.ResponseWriter, detail string) {
	w.Header().Set("Content-Type", "application/json")

	payload := struct {
		Status         string `json:"status"`
		HttpStatusCode int    `json:"httpStatusCode"`
		Message        string `json:"message"`
	}{
		Status:         "error",
		HttpStatusCode: http.StatusUnauthorized,
		Message:        "Bad request: " + detail,
	}

	body, err := json.Marshal(&payload)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(fmt.Sprintf("Failed to create response: %s", err)))
		return
	}

	w.WriteHeader(http.StatusUnauthorized)
	w.Write(body)
}
