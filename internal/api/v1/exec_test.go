package api

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	ds "subscription-vault-service/internal/datastruct"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
)

type TestRequest struct {
	Login    string `json:"login" schema:"login" validate:"required" example:"VaKadyk359"`
	Password string `json:"password" schema:"password" validate:"required" example:"Xldf32Q"`
}

type TestResponse struct {
	ds.Status
}

func TestExtractJsonBody(t *testing.T) {
	t.Parallel()

	t.Run("Ok", func(t *testing.T) {
		t.Parallel()

		v := TestRequest{
			Login:    "login",
			Password: "password",
		}

		body, _ := json.Marshal(v)
		r := httptest.NewRequest(http.MethodPost, "/test", bytes.NewReader(body))

		res := TestRequest{}
		err := extractJsonBody(r, &res)
		require.Nil(t, err)

		require.Equal(t, v.Login, res.Login)
		require.Equal(t, v.Password, res.Password)
	})

	t.Run("error", func(t *testing.T) {
		t.Parallel()

		r := httptest.NewRequest(http.MethodPost, "/test", strings.NewReader("{asdfae}"))

		res := TestRequest{}
		err := extractJsonBody(r, &res)
		require.NotNil(t, err)

	})
}

func TestWriteJsonResponse(t *testing.T) {
	t.Parallel()

	t.Run("Ok", func(t *testing.T) {
		t.Parallel()

		v := TestResponse{
			Status: ds.Status{Message: "ok"},
		}

		w := httptest.NewRecorder()
		err := writeJsonResponse(w, &v)
		require.Nil(t, err)

		res := TestResponse{}
		err = json.Unmarshal(w.Body.Bytes(), &res)
		require.Nil(t, err)

		require.Equal(t, v.GetStatus(), res.GetStatus())
	})
}

func TestExtractSchemaQuery(t *testing.T) {
	t.Parallel()

	t.Run("Ok", func(t *testing.T) {
		t.Parallel()

		r := httptest.NewRequest(http.MethodPost, "/test", strings.NewReader(""))
		q := r.URL.Query()
		q.Add("login", "abra")
		q.Add("password", "10.01.2001")
		r.URL.RawQuery = q.Encode()

		v := TestRequest{}
		err := extractSchemaQuery(r, &v)
		require.Nil(t, err)

		require.Equal(t, v.Login, "abra")
		require.Equal(t, v.Password, "10.01.2001")

	})
}

func TestStructValidator(t *testing.T) {
	t.Parallel()

	t.Run("Ok", func(t *testing.T) {
		t.Parallel()

		v := TestRequest{
			Login:    "login",
			Password: "password",
		}

		err := structValidator(v)
		require.Nil(t, err)
	})

	t.Run("error", func(t *testing.T) {
		t.Parallel()

		v := TestRequest{
			Login: "login",
		}

		err := structValidator(v)
		require.NotNil(t, err)
	})
}

func TestExec(t *testing.T) {
	t.Parallel()

	t.Run("Ok", func(t *testing.T) {
		t.Parallel()

		ta := newTestAPI(t)

		v := TestRequest{
			Login:    "login",
			Password: "password",
		}

		body, _ := json.Marshal(v)

		r := httptest.NewRequest(http.MethodPost, "/test", bytes.NewReader(body))
		w := httptest.NewRecorder()

		Exec(ExecArgs[TestRequest, TestResponse]{
			api: ta.a,
			serviceFunc: func(_ *TestRequest) *TestResponse {
				return &TestResponse{
					Status: ds.Status{Message: ds.StatusSuccess},
				}
			},
			requestExtractor: extractJsonBody[*TestRequest],
			responseWriter:   writeJsonResponse[*TestResponse],
			httpRequest:      r,
			httpResponse:     w,
		})

		respS := TestResponse{}
		require.Nil(t, json.Unmarshal(w.Body.Bytes(), &respS))
		require.Equal(t, respS.GetStatus(), ds.StatusSuccess)
	})

	t.Run("extracting error", func(t *testing.T) {
		t.Parallel()

		ta := newTestAPI(t)

		v := TestRequest{
			Login:    "login",
			Password: "password",
		}

		body, _ := json.Marshal(v)

		r := httptest.NewRequest(http.MethodPost, "/test", bytes.NewReader(body))
		w := httptest.NewRecorder()

		ta.loggerMock.EXPECT().ErrorKV(gomock.Any(), gomock.All())

		Exec(ExecArgs[TestRequest, TestResponse]{
			api: ta.a,
			serviceFunc: func(lr *TestRequest) *TestResponse {
				return &TestResponse{
					Status: ds.Status{Message: ds.StatusSuccess},
				}
			},
			requestExtractor: func(r *http.Request, v *TestRequest) error {
				return ds.ErrTest
			},
			responseWriter: writeJsonResponse[*TestResponse],
			httpRequest:    r,
			httpResponse:   w,
		})

		require.Equal(t, w.Code, http.StatusBadRequest)
		respS := TestResponse{}
		require.Nil(t, json.Unmarshal(w.Body.Bytes(), &respS))
		require.Equal(t, respS.GetStatus(), ds.StatusFailedExctractingRequest)
	})

	t.Run("validating error", func(t *testing.T) {
		t.Parallel()

		ta := newTestAPI(t)

		v := TestRequest{
			Login:    "login",
			Password: "password",
		}

		body, _ := json.Marshal(v)

		r := httptest.NewRequest(http.MethodPost, "/test", bytes.NewReader(body))
		w := httptest.NewRecorder()

		ta.loggerMock.EXPECT().ErrorKV(gomock.Any(), gomock.All())

		Exec(ExecArgs[TestRequest, TestResponse]{
			api: ta.a,
			serviceFunc: func(lr *TestRequest) *TestResponse {
				return &TestResponse{
					Status: ds.Status{Message: ds.StatusSuccess},
				}
			},
			requestExtractor: extractJsonBody[*TestRequest],
			responseWriter:   writeJsonResponse[*TestResponse],
			httpRequest:      r,
			httpResponse:     w,
			validator: func(s *TestRequest) error {
				return ds.ErrTest
			},
		})

		require.Equal(t, w.Code, http.StatusBadRequest)
		respS := TestResponse{}
		require.Nil(t, json.Unmarshal(w.Body.Bytes(), &respS))
		require.Equal(t, respS.GetStatus(), ds.StatusFailedValidatingRequest)
	})

	t.Run("service func error", func(t *testing.T) {
		t.Parallel()

		ta := newTestAPI(t)

		v := TestRequest{
			Login:    "login",
			Password: "password",
		}

		body, _ := json.Marshal(v)

		r := httptest.NewRequest(http.MethodPost, "/test", bytes.NewReader(body))
		w := httptest.NewRecorder()

		ta.loggerMock.EXPECT().ErrorKV(gomock.Any(), gomock.All())

		Exec(ExecArgs[TestRequest, TestResponse]{
			api: ta.a,
			serviceFunc: func(lr *TestRequest) *TestResponse {
				return nil
			},
			requestExtractor: extractJsonBody[*TestRequest],
			responseWriter:   writeJsonResponse[*TestResponse],
			httpRequest:      r,
			httpResponse:     w,
		})

		require.Equal(t, w.Code, http.StatusInternalServerError)
		respS := TestResponse{}
		require.Nil(t, json.Unmarshal(w.Body.Bytes(), &respS))
		require.Equal(t, respS.GetStatus(), ds.StatusServiceError)
	})

	t.Run("response writer error", func(t *testing.T) {
		t.Parallel()

		ta := newTestAPI(t)

		v := TestRequest{
			Login:    "login",
			Password: "password",
		}

		body, _ := json.Marshal(v)

		r := httptest.NewRequest(http.MethodPost, "/test", bytes.NewReader(body))
		w := httptest.NewRecorder()

		ta.loggerMock.EXPECT().ErrorKV(gomock.Any(), gomock.All())

		Exec(ExecArgs[TestRequest, TestResponse]{
			api: ta.a,
			serviceFunc: func(lr *TestRequest) *TestResponse {
				return &TestResponse{
					Status: ds.Status{Message: ds.StatusSuccess},
				}
			},
			requestExtractor: extractJsonBody[*TestRequest],
			responseWriter: func(w http.ResponseWriter, v *TestResponse) error {
				return ds.ErrTest
			},
			httpRequest:  r,
			httpResponse: w,
		})

		require.Equal(t, w.Code, http.StatusInternalServerError)

	})
}
