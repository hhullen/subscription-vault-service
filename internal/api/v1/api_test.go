package api

import (
	"context"
	"subscription-vault-service/internal/service"
	"testing"

	gomock "github.com/golang/mock/gomock"
)

type TestAPI struct {
	ctx        context.Context
	subsMock   *MockISubscriptionService
	serverMock *MockIServer
	routerMock *MockIRouter
	loggerMock *service.MockILogger
	secretMock *MockISecretProvider
	a          *API
}

func newTestAPI(t *testing.T) *TestAPI {
	mc := gomock.NewController(t)
	ta := &TestAPI{
		ctx:        context.Background(),
		subsMock:   NewMockISubscriptionService(mc),
		serverMock: NewMockIServer(mc),
		routerMock: NewMockIRouter(mc),
		loggerMock: service.NewMockILogger(mc),
		secretMock: NewMockISecretProvider(mc),
	}

	ta.routerMock.EXPECT().Handle(gomock.Any(), gomock.Any()).MinTimes(1)

	a := buildAPI(ta.ctx, ta.subsMock, ta.loggerMock, ta.secretMock, ta.serverMock, ta.routerMock)

	ta.a = a

	return ta
}
