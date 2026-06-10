package api

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"strconv"

	ds "subscription-vault-service/internal/datastruct"
	"subscription-vault-service/internal/supports"
)

const (
	contentTypeKey = "Content-Type"
	contentLenKey  = "Content-Length"
	appJSONValue   = "application/json"
)

type ExecArgs[ReqT any, RespT IWithStatus] struct {
	api              *API
	serviceFunc      func(*ReqT) *RespT
	requestExtractor func(r *http.Request, v *ReqT) error
	responseWriter   func(w http.ResponseWriter, v *RespT) error
	httpRequest      *http.Request
	httpResponse     http.ResponseWriter
	validator        func(s *ReqT) error
}

func getStatusCode(s string) int {
	switch s {
	case ds.StatusUserNotFound:
		return http.StatusNotFound
	case ds.StatusResurceNotFound:
		return http.StatusNotFound
	case ds.StatusServiceError:
		return http.StatusInternalServerError
	case ds.StatusResourceAlreadyExists:
		return http.StatusConflict
	case ds.StatusFailedExctractingRequest:
		return http.StatusBadRequest
	case ds.StatusFailedValidatingRequest:
		return http.StatusBadRequest
	}

	return http.StatusOK
}

func extractJsonBody[ReqT any](r *http.Request, v ReqT) error {
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&v); err != nil && err != io.EOF {
		return err
	}
	return nil
}

func writeJsonResponse[RespT IWithStatus](w http.ResponseWriter, resp RespT) error {
	var buf bytes.Buffer
	if err := json.NewEncoder(&buf).Encode(&resp); err != nil {
		return err
	}

	w.Header().Set(contentLenKey, strconv.Itoa(len(buf.Bytes())))
	w.Header().Set(contentTypeKey, appJSONValue)

	code := getStatusCode(resp.GetStatus())
	w.WriteHeader(code)
	_, err := w.Write(buf.Bytes())
	if err != nil {
		return err
	}

	return nil
}

func extractSchemaQuery[ReqT any](r *http.Request, v ReqT) error {
	if err := schemaDecoder.Decode(v, r.URL.Query()); err != nil && err != io.EOF {
		return err
	}

	return nil
}

func structValidator[ReqT any](s ReqT) error {
	return supports.StructValidator().Struct(s)
}

func Exec[ReqT any, RespT IWithStatus](a ExecArgs[ReqT, RespT]) {
	var req ReqT

	if err := a.requestExtractor(a.httpRequest, &req); err != nil {
		const msg = "failed extracting request"
		a.api.logger.ErrorKV(msg, "error", err.Error())

		resp := ds.Status{Message: ds.StatusFailedExctractingRequest}
		err = writeJsonResponse(a.httpResponse, resp)
		if err != nil {
			a.api.logger.ErrorKV("failed write response",
				"error", err.Error(), "response", resp)
		}

		return
	}

	if a.validator == nil {
		a.validator = structValidator
	}

	if err := a.validator(&req); err != nil {
		const msg = "failed validating request"
		a.api.logger.ErrorKV(msg, "error", err.Error(), "request", req)

		resp := ds.Status{Message: ds.StatusFailedValidatingRequest}
		err = writeJsonResponse(a.httpResponse, resp)
		if err != nil {
			a.api.logger.ErrorKV("failed write response",
				"error", err.Error(), "response", resp)
		}
		return
	}

	resp := a.serviceFunc(&req)
	if resp == nil {
		const msg = "failed execute request on service"
		a.api.logger.ErrorKV(msg, "error", "service return no response", "request", req)
		resp := ds.Status{Message: ds.StatusServiceError}
		err := writeJsonResponse(a.httpResponse, resp)
		if err != nil {

		}
		return
	}

	if err := a.responseWriter(a.httpResponse, resp); err != nil {
		const msg = "failed writing response"
		http.Error(a.httpResponse, msg, http.StatusInternalServerError)
		a.api.logger.ErrorKV(msg, "error", err.Error(), "request", req)
	}
}
