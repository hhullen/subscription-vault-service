package service

import (
	"context"
	"testing"
	"time"

	ds "subscription-vault-service/internal/datastruct"

	gomock "github.com/golang/mock/gomock"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
)

type TestService struct {
	ctx             context.Context
	subsStorageMock *MockISubscriptionsStorage
	loggerMock      *MockILogger
	service         *Service
}

func newTestService(t *testing.T) *TestService {
	mc := gomock.NewController(t)
	ts := &TestService{
		ctx:             context.Background(),
		subsStorageMock: NewMockISubscriptionsStorage(mc),
		loggerMock:      NewMockILogger(mc),
	}

	service := NewService(ts.ctx, ts.subsStorageMock, ts.loggerMock)
	ts.service = service

	return ts
}

func TestNewService(t *testing.T) {
	t.Parallel()

	t.Run("Ok", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		mc := gomock.NewController(t)
		subsStorageMock := NewMockISubscriptionsStorage(mc)
		loggerMock := NewMockILogger(mc)

		service := NewService(ctx, subsStorageMock, loggerMock)

		require.NotNil(t, service)
		require.Equal(t, ctx, service.ctx)
		require.Equal(t, subsStorageMock, service.subsStore)
		require.Equal(t, loggerMock, service.log)
	})
}

func TestBuildService(t *testing.T) {
	t.Parallel()

	t.Run("Ok", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		mc := gomock.NewController(t)
		subsStorageMock := NewMockISubscriptionsStorage(mc)
		loggerMock := NewMockILogger(mc)

		service := buildService(ctx, subsStorageMock, loggerMock)

		require.NotNil(t, service)
		require.Equal(t, ctx, service.ctx)
		require.Equal(t, subsStorageMock, service.subsStore)
		require.Equal(t, loggerMock, service.log)
	})
}

func TestAddSubscription(t *testing.T) {
	t.Parallel()

	t.Run("Ok", func(t *testing.T) {
		t.Parallel()

		ts := newTestService(t)

		req := &ds.CreateSubscriptionRequest{
			Subscription: ds.Subscription{
				Price: 100,
				SubscriptionID: ds.SubscriptionID{
					ServiceName: "test-service",
					UserUid:     uuid.New(),
				},
				Period: ds.Period{
					From: ds.DateType(time.Now()),
					To:   ds.DateType(time.Now().AddDate(0, 1, 0)),
				},
			},
		}

		expectedResp := &ds.CreateSubscriptionResponse{
			Status: ds.Status{Message: ds.StatusSuccess},
		}

		ts.subsStorageMock.EXPECT().AddSubscription(req).Return(expectedResp, nil)

		resp := ts.service.AddSubscription(req)

		require.NotNil(t, resp)
		require.Equal(t, expectedResp.Status.Message, resp.Status.Message)
	})

	t.Run("Error", func(t *testing.T) {
		t.Parallel()

		ts := newTestService(t)

		req := &ds.CreateSubscriptionRequest{
			Subscription: ds.Subscription{
				Price: 100,
				SubscriptionID: ds.SubscriptionID{
					ServiceName: "test-service",
					UserUid:     uuid.New(),
				},
				Period: ds.Period{
					From: ds.DateType(time.Now()),
					To:   ds.DateType(time.Now().AddDate(0, 1, 0)),
				},
			},
		}

		ts.subsStorageMock.EXPECT().AddSubscription(req).Return(nil, ds.ErrTest)
		ts.loggerMock.EXPECT().ErrorKV("AddSubscription.AddSubscription", "error", ds.ErrTest.Error())

		resp := ts.service.AddSubscription(req)

		require.Nil(t, resp)
	})
}

func TestGetSubscription(t *testing.T) {
	t.Parallel()

	t.Run("Ok", func(t *testing.T) {
		t.Parallel()

		ts := newTestService(t)

		userUid := uuid.New()
		req := &ds.GetSubscriptionRequest{
			SubscriptionID: ds.SubscriptionID{
				ServiceName: "test-service",
				UserUid:     userUid,
			},
		}

		expectedResp := &ds.GetSubscriptionResponse{
			Subscription: &ds.Subscription{
				Price: 100,
				SubscriptionID: ds.SubscriptionID{
					ServiceName: "test-service",
					UserUid:     userUid,
				},
				Period: ds.Period{
					From: ds.DateType(time.Now()),
					To:   ds.DateType(time.Now().AddDate(0, 1, 0)),
				},
			},
			Status: ds.Status{Message: ds.StatusSuccess},
		}

		ts.subsStorageMock.EXPECT().GetSubscription(req).Return(expectedResp, nil)

		resp := ts.service.GetSubscription(req)

		require.NotNil(t, resp)
		require.Equal(t, expectedResp.Status.Message, resp.Status.Message)
		require.Equal(t, expectedResp.Subscription.Price, resp.Subscription.Price)
	})

	t.Run("Error", func(t *testing.T) {
		t.Parallel()

		ts := newTestService(t)

		userUid := uuid.New()
		req := &ds.GetSubscriptionRequest{
			SubscriptionID: ds.SubscriptionID{
				ServiceName: "test-service",
				UserUid:     userUid,
			},
		}

		ts.subsStorageMock.EXPECT().GetSubscription(req).Return(nil, ds.ErrTest)
		ts.loggerMock.EXPECT().ErrorKV("GetSubscription.GetSubscription", "error", ds.ErrTest.Error())

		resp := ts.service.GetSubscription(req)

		require.Nil(t, resp)
	})

	t.Run("Not Found", func(t *testing.T) {
		t.Parallel()

		ts := newTestService(t)

		userUid := uuid.New()
		req := &ds.GetSubscriptionRequest{
			SubscriptionID: ds.SubscriptionID{
				ServiceName: "test-service",
				UserUid:     userUid,
			},
		}

		expectedResp := &ds.GetSubscriptionResponse{
			Status: ds.Status{Message: ds.StatusResurceNotFound},
		}

		ts.subsStorageMock.EXPECT().GetSubscription(req).Return(expectedResp, nil)

		resp := ts.service.GetSubscription(req)

		require.NotNil(t, resp)
		require.Equal(t, ds.StatusResurceNotFound, resp.Status.Message)
	})
}

func TestUpdateSubscription(t *testing.T) {
	t.Parallel()

	t.Run("Ok", func(t *testing.T) {
		t.Parallel()

		ts := newTestService(t)

		userUid := uuid.New()
		req := &ds.UpdateSubscriptionRequest{
			Subscription: ds.Subscription{
				Price: 150,
				SubscriptionID: ds.SubscriptionID{
					ServiceName: "test-service",
					UserUid:     userUid,
				},
				Period: ds.Period{
					From: ds.DateType(time.Now()),
					To:   ds.DateType(time.Now().AddDate(0, 1, 0)),
				},
			},
		}

		expectedResp := &ds.UpdateSubscriptionResponse{
			Status: ds.Status{Message: ds.StatusSuccess},
		}

		ts.subsStorageMock.EXPECT().UpdateSubscription(req).Return(expectedResp, nil)

		resp := ts.service.UpdateSubscription(req)

		require.NotNil(t, resp)
		require.Equal(t, ds.StatusSuccess, resp.Status.Message)
	})

	t.Run("Error", func(t *testing.T) {
		t.Parallel()

		ts := newTestService(t)

		userUid := uuid.New()
		req := &ds.UpdateSubscriptionRequest{
			Subscription: ds.Subscription{
				Price: 150,
				SubscriptionID: ds.SubscriptionID{
					ServiceName: "test-service",
					UserUid:     userUid,
				},
				Period: ds.Period{
					From: ds.DateType(time.Now()),
					To:   ds.DateType(time.Now().AddDate(0, 1, 0)),
				},
			},
		}

		ts.subsStorageMock.EXPECT().UpdateSubscription(req).Return(nil, ds.ErrTest)
		ts.loggerMock.EXPECT().ErrorKV("UpdateSubscription.UpdateSubscription", "error", ds.ErrTest.Error())

		resp := ts.service.UpdateSubscription(req)

		require.Nil(t, resp)
	})

	t.Run("Not Found", func(t *testing.T) {
		t.Parallel()

		ts := newTestService(t)

		userUid := uuid.New()
		req := &ds.UpdateSubscriptionRequest{
			Subscription: ds.Subscription{
				Price: 150,
				SubscriptionID: ds.SubscriptionID{
					ServiceName: "test-service",
					UserUid:     userUid,
				},
				Period: ds.Period{
					From: ds.DateType(time.Now()),
					To:   ds.DateType(time.Now().AddDate(0, 1, 0)),
				},
			},
		}

		expectedResp := &ds.UpdateSubscriptionResponse{
			Status: ds.Status{Message: ds.StatusResurceNotFound},
		}

		ts.subsStorageMock.EXPECT().UpdateSubscription(req).Return(expectedResp, nil)

		resp := ts.service.UpdateSubscription(req)

		require.NotNil(t, resp)
		require.Equal(t, ds.StatusResurceNotFound, resp.Status.Message)
	})
}

func TestDeleteSubscription(t *testing.T) {
	t.Parallel()

	t.Run("Ok", func(t *testing.T) {
		t.Parallel()

		ts := newTestService(t)

		userUid := uuid.New()
		req := &ds.DeleteSubscriptionRequest{
			SubscriptionID: ds.SubscriptionID{
				ServiceName: "test-service",
				UserUid:     userUid,
			},
		}

		expectedResp := &ds.DeleteSubscriptionResponse{
			Status: ds.Status{Message: ds.StatusSuccess},
		}

		ts.subsStorageMock.EXPECT().DeleteSubscription(req).Return(expectedResp, nil)

		resp := ts.service.DeleteSubscription(req)

		require.NotNil(t, resp)
		require.Equal(t, ds.StatusSuccess, resp.Status.Message)
	})

	t.Run("Error", func(t *testing.T) {
		t.Parallel()

		ts := newTestService(t)

		userUid := uuid.New()
		req := &ds.DeleteSubscriptionRequest{
			SubscriptionID: ds.SubscriptionID{
				ServiceName: "test-service",
				UserUid:     userUid,
			},
		}

		ts.subsStorageMock.EXPECT().DeleteSubscription(req).Return(nil, ds.ErrTest)
		ts.loggerMock.EXPECT().ErrorKV("DeleteSubscription.DeleteSubscription", "error", ds.ErrTest.Error())

		resp := ts.service.DeleteSubscription(req)

		require.Nil(t, resp)
	})

	t.Run("Not Found", func(t *testing.T) {
		t.Parallel()

		ts := newTestService(t)

		userUid := uuid.New()
		req := &ds.DeleteSubscriptionRequest{
			SubscriptionID: ds.SubscriptionID{
				ServiceName: "test-service",
				UserUid:     userUid,
			},
		}

		expectedResp := &ds.DeleteSubscriptionResponse{
			Status: ds.Status{Message: ds.StatusResurceNotFound},
		}

		ts.subsStorageMock.EXPECT().DeleteSubscription(req).Return(expectedResp, nil)

		resp := ts.service.DeleteSubscription(req)

		require.NotNil(t, resp)
		require.Equal(t, ds.StatusResurceNotFound, resp.Status.Message)
	})
}

func TestListSubscriptions(t *testing.T) {
	t.Parallel()

	t.Run("Ok", func(t *testing.T) {
		t.Parallel()

		ts := newTestService(t)

		userUid := uuid.New()
		req := &ds.ListSubscriptionsRequest{
			UserUid: userUid,
		}

		subscriptions := []ds.Subscription{
			{
				Price: 100,
				SubscriptionID: ds.SubscriptionID{
					ServiceName: "service1",
					UserUid:     userUid,
				},
				Period: ds.Period{
					From: ds.DateType(time.Now()),
					To:   ds.DateType(time.Now().AddDate(0, 1, 0)),
				},
			},
			{
				Price: 50,
				SubscriptionID: ds.SubscriptionID{
					ServiceName: "service2",
					UserUid:     userUid,
				},
				Period: ds.Period{
					From: ds.DateType(time.Now()),
					To:   ds.DateType(time.Now().AddDate(0, 1, 0)),
				},
			},
		}

		expectedResp := &ds.ListSubscriptionsResponse{
			Subscriptions: subscriptions,
			Status:        ds.Status{Message: ds.StatusSuccess},
		}

		ts.subsStorageMock.EXPECT().ListSubscriptions(req).Return(expectedResp, nil)

		resp := ts.service.ListSubscriptions(req)

		require.NotNil(t, resp)
		require.Equal(t, ds.StatusSuccess, resp.Status.Message)
		require.Len(t, resp.Subscriptions, 2)
		require.Equal(t, subscriptions[0].SubscriptionID.ServiceName, resp.Subscriptions[0].SubscriptionID.ServiceName)
		require.Equal(t, subscriptions[1].SubscriptionID.ServiceName, resp.Subscriptions[1].SubscriptionID.ServiceName)
	})

	t.Run("Empty List", func(t *testing.T) {
		t.Parallel()

		ts := newTestService(t)

		userUid := uuid.New()
		req := &ds.ListSubscriptionsRequest{
			UserUid: userUid,
		}

		expectedResp := &ds.ListSubscriptionsResponse{
			Subscriptions: []ds.Subscription{},
			Status:        ds.Status{Message: ds.StatusSuccess},
		}

		ts.subsStorageMock.EXPECT().ListSubscriptions(req).Return(expectedResp, nil)

		resp := ts.service.ListSubscriptions(req)

		require.NotNil(t, resp)
		require.Equal(t, ds.StatusSuccess, resp.Status.Message)
		require.Len(t, resp.Subscriptions, 0)
	})

	t.Run("Error", func(t *testing.T) {
		t.Parallel()

		ts := newTestService(t)

		userUid := uuid.New()
		req := &ds.ListSubscriptionsRequest{
			UserUid: userUid,
		}

		ts.subsStorageMock.EXPECT().ListSubscriptions(req).Return(nil, ds.ErrTest)
		ts.loggerMock.EXPECT().ErrorKV("ListSubscriptions.ListSubscriptions", "error", ds.ErrTest.Error())

		resp := ts.service.ListSubscriptions(req)

		require.Nil(t, resp)
	})
}

func TestCalculateSubscriptionsPrice(t *testing.T) {
	t.Parallel()

	t.Run("Ok", func(t *testing.T) {
		t.Parallel()

		ts := newTestService(t)

		userUid := uuid.New()
		req := &ds.CalculateSubscriptionsPriceRequest{
			Period: ds.Period{
				From: ds.DateType(time.Now()),
				To:   ds.DateType(time.Now().AddDate(0, 1, 0)),
			},
			UserUid: &userUid,
		}

		expectedResp := &ds.CalculateSubscriptionsPriceResponse{
			TotalPrice: 150,
			Status:     ds.Status{Message: ds.StatusSuccess},
		}

		ts.subsStorageMock.EXPECT().CalculateSubscriptionsPrice(req).Return(expectedResp, nil)

		resp := ts.service.CalculateSubscriptionsPrice(req)

		require.NotNil(t, resp)
		require.Equal(t, ds.StatusSuccess, resp.Status.Message)
		require.Equal(t, int64(150), resp.TotalPrice)
	})

	t.Run("Zero Price", func(t *testing.T) {
		t.Parallel()

		ts := newTestService(t)

		userUid := uuid.New()
		req := &ds.CalculateSubscriptionsPriceRequest{
			Period: ds.Period{
				From: ds.DateType(time.Now()),
				To:   ds.DateType(time.Now().AddDate(0, 1, 0)),
			},
			UserUid: &userUid,
		}

		expectedResp := &ds.CalculateSubscriptionsPriceResponse{
			TotalPrice: 0,
			Status:     ds.Status{Message: ds.StatusSuccess},
		}

		ts.subsStorageMock.EXPECT().CalculateSubscriptionsPrice(req).Return(expectedResp, nil)

		resp := ts.service.CalculateSubscriptionsPrice(req)

		require.NotNil(t, resp)
		require.Equal(t, ds.StatusSuccess, resp.Status.Message)
		require.Equal(t, int64(0), resp.TotalPrice)
	})

	t.Run("Error", func(t *testing.T) {
		t.Parallel()

		ts := newTestService(t)

		userUid := uuid.New()
		req := &ds.CalculateSubscriptionsPriceRequest{
			Period: ds.Period{
				From: ds.DateType(time.Now()),
				To:   ds.DateType(time.Now().AddDate(0, 1, 0)),
			},
			UserUid: &userUid,
		}

		ts.subsStorageMock.EXPECT().CalculateSubscriptionsPrice(req).Return(nil, ds.ErrTest)
		ts.loggerMock.EXPECT().ErrorKV("CalculateSubscriptionsPrice.CalculateSubscriptionsPrice", "error", ds.ErrTest.Error())

		resp := ts.service.CalculateSubscriptionsPrice(req)

		require.Nil(t, resp)
	})
}
