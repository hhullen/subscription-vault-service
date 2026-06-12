package api

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	ds "subscription-vault-service/internal/datastruct"
	"subscription-vault-service/internal/service"

	gomock "github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
)

type TestAPI struct {
	ctx         context.Context
	subsMock    *MockISubscriptionService
	serverMock  *MockIServer
	routerMock  *MockIRouter
	loggerMock  *service.MockILogger
	handlerMock *MockHandler
	a           *API
}

func newTestAPI(t *testing.T) *TestAPI {
	mc := gomock.NewController(t)
	ta := &TestAPI{
		ctx:         context.Background(),
		subsMock:    NewMockISubscriptionService(mc),
		serverMock:  NewMockIServer(mc),
		routerMock:  NewMockIRouter(mc),
		loggerMock:  service.NewMockILogger(mc),
		handlerMock: NewMockHandler(mc),
	}

	ta.routerMock.EXPECT().Handle(gomock.Any(), gomock.Any()).MinTimes(1)

	a := buildAPI(ta.ctx, ta.subsMock, ta.loggerMock, ta.serverMock, ta.routerMock)

	ta.a = a

	return ta
}

func TestStartListening(t *testing.T) {
	t.Parallel()

	t.Run("Ok", func(t *testing.T) {
		t.Parallel()

		ta := newTestAPI(t)

		ta.serverMock.EXPECT().ListenAndServe().Return(nil)
		ta.loggerMock.EXPECT().InfoKV(gomock.Any(), gomock.All())

		err := ta.a.StartListening()
		require.Nil(t, err)
	})

	t.Run("error", func(t *testing.T) {
		t.Parallel()

		ta := newTestAPI(t)

		ta.serverMock.EXPECT().ListenAndServe().Return(ds.ErrTest)
		ta.loggerMock.EXPECT().InfoKV(gomock.Any(), gomock.All())

		err := ta.a.StartListening()
		require.NotNil(t, err)
	})
}

func TestStop(t *testing.T) {
	t.Parallel()

	t.Run("Ok", func(t *testing.T) {
		t.Parallel()

		ta := newTestAPI(t)

		ta.loggerMock.EXPECT().InfoKV(gomock.Any(), gomock.All())
		ta.serverMock.EXPECT().Shutdown(gomock.Any()).Return(nil)

		err := ta.a.Stop()
		require.Nil(t, err)
	})

	t.Run("Ok", func(t *testing.T) {
		t.Parallel()

		ta := newTestAPI(t)

		ta.loggerMock.EXPECT().InfoKV(gomock.Any(), gomock.All())
		ta.serverMock.EXPECT().Shutdown(gomock.Any()).Return(ds.ErrTest)

		err := ta.a.Stop()
		require.NotNil(t, err)
	})
}

func TestServeHTTP(t *testing.T) {
	t.Parallel()

	t.Run("Ok", func(t *testing.T) {
		t.Parallel()

		ta := newTestAPI(t)

		ta.routerMock.EXPECT().ServeHTTP(gomock.Any(), gomock.Any())

		ta.a.ServeHTTP(httptest.NewRecorder(), &http.Request{})

	})

}

func TestMainMiddleware(t *testing.T) {
	t.Parallel()

	t.Run("Ok", func(t *testing.T) {
		t.Parallel()

		ta := newTestAPI(t)

		ta.handlerMock.EXPECT().ServeHTTP(gomock.Any(), gomock.Any())
		ta.loggerMock.EXPECT().InfoKV(gomock.Any(), gomock.All())

		r := httptest.NewRequest(http.MethodGet, "/test", strings.NewReader("test"))
		w := httptest.NewRecorder()
		mainMiddleware(ta.handlerMock, ta.loggerMock).ServeHTTP(w, r)
	})
}
