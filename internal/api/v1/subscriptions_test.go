package api

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	ds "subscription-vault-service/internal/datastruct"

	gomock "github.com/golang/mock/gomock"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
)

func TestCreateSubscription(t *testing.T) {
	t.Parallel()

	t.Run("Ok", func(t *testing.T) {
		t.Parallel()

		ta := newTestAPI(t)

		reqS := ds.CreateSubscriptionRequest{
			Subscription: ds.Subscription{
				Price: 100,
				SubscriptionID: ds.SubscriptionID{
					ServiceName: "name",
					UserUid:     uuid.New(),
				},
				Period: ds.Period{
					From: ds.DateType(time.Now()),
					To:   ds.DateType(time.Now()),
				},
			},
		}

		ta.subsMock.EXPECT().AddSubscription(gomock.Any()).Return(&ds.CreateSubscriptionResponse{
			Status: ds.Status{Message: ds.StatusSuccess},
		})

		var buf bytes.Buffer
		_ = json.NewEncoder(&buf).Encode(&reqS)

		r := httptest.NewRequest(http.MethodPost, "/test", bytes.NewReader(buf.Bytes()))
		w := httptest.NewRecorder()
		ta.a.CreateSubscription(w, r)

		require.Equal(t, http.StatusOK, w.Code)

		respS := ds.CreateSubscriptionResponse{}
		require.Nil(t, json.Unmarshal(w.Body.Bytes(), &respS))

		require.Equal(t, ds.StatusSuccess, respS.Status.Message)
	})

	t.Run("No required field", func(t *testing.T) {
		t.Parallel()

		ta := newTestAPI(t)

		reqS := ds.CreateSubscriptionRequest{}

		ta.loggerMock.EXPECT().ErrorKV(gomock.Any(), gomock.All())

		var buf bytes.Buffer
		_ = json.NewEncoder(&buf).Encode(&reqS)
		body := buf.Bytes()

		r := httptest.NewRequest(http.MethodPost, "/test", bytes.NewReader(body))
		w := httptest.NewRecorder()
		ta.a.CreateSubscription(w, r)

		require.Equal(t, http.StatusBadRequest, w.Code)

		respS := ds.CreateSubscriptionResponse{}
		require.Nil(t, json.Unmarshal(w.Body.Bytes(), &respS))

		require.Equal(t, ds.StatusFailedValidatingRequest, respS.Status.Message)
	})

	t.Run("No response from service", func(t *testing.T) {
		t.Parallel()

		ta := newTestAPI(t)

		reqS := ds.CreateSubscriptionRequest{
			Subscription: ds.Subscription{
				Price: 100,
				SubscriptionID: ds.SubscriptionID{
					ServiceName: "name",
					UserUid:     uuid.New(),
				},
				Period: ds.Period{
					From: ds.DateType(time.Now()),
					To:   ds.DateType(time.Now()),
				},
			},
		}

		ta.subsMock.EXPECT().AddSubscription(gomock.Any()).Return(nil)
		ta.loggerMock.EXPECT().ErrorKV(gomock.Any(), gomock.All())

		var buf bytes.Buffer
		_ = json.NewEncoder(&buf).Encode(&reqS)
		body := buf.Bytes()

		r := httptest.NewRequest(http.MethodPost, "/test", bytes.NewReader(body))
		w := httptest.NewRecorder()
		ta.a.CreateSubscription(w, r)

		require.Equal(t, http.StatusInternalServerError, w.Code)

		respS := ds.CreateSubscriptionResponse{}
		require.Nil(t, json.Unmarshal(w.Body.Bytes(), &respS))

		require.Equal(t, ds.StatusServiceError, respS.Status.Message)
	})
}

func TestGetSubscription(t *testing.T) {
	t.Parallel()

	t.Run("Ok", func(t *testing.T) {
		t.Parallel()

		ta := newTestAPI(t)

		reqS := ds.GetSubscriptionRequest{
			SubscriptionID: ds.SubscriptionID{
				ServiceName: "name",
				UserUid:     uuid.New(),
			},
		}

		ta.subsMock.EXPECT().GetSubscription(gomock.Any()).Return(&ds.GetSubscriptionResponse{
			Status: ds.Status{Message: ds.StatusSuccess},
		})

		r := httptest.NewRequest(http.MethodGet, "/test", bytes.NewReader([]byte{}))
		q := r.URL.Query()
		q.Add("service_name", reqS.ServiceName)
		q.Add("user_uid", reqS.UserUid.String())
		r.URL.RawQuery = q.Encode()
		w := httptest.NewRecorder()
		ta.a.GetSubscription(w, r)

		require.Equal(t, http.StatusOK, w.Code)

		respS := ds.GetSubscriptionResponse{}
		require.Nil(t, json.Unmarshal(w.Body.Bytes(), &respS))

		require.Equal(t, ds.StatusSuccess, respS.Status.Message)
	})

	t.Run("No required field", func(t *testing.T) {
		t.Parallel()

		ta := newTestAPI(t)

		ta.loggerMock.EXPECT().ErrorKV(gomock.Any(), gomock.All())

		r := httptest.NewRequest(http.MethodGet, "/test", bytes.NewReader([]byte{}))
		w := httptest.NewRecorder()
		ta.a.GetSubscription(w, r)

		require.Equal(t, http.StatusBadRequest, w.Code)

		respS := ds.GetSubscriptionResponse{}
		require.Nil(t, json.Unmarshal(w.Body.Bytes(), &respS))

		require.Equal(t, ds.StatusFailedValidatingRequest, respS.Status.Message)
	})

	t.Run("No response from service", func(t *testing.T) {
		t.Parallel()

		ta := newTestAPI(t)

		reqS := ds.GetSubscriptionRequest{
			SubscriptionID: ds.SubscriptionID{
				ServiceName: "name",
				UserUid:     uuid.New(),
			},
		}

		ta.subsMock.EXPECT().GetSubscription(gomock.Any()).Return(nil)
		ta.loggerMock.EXPECT().ErrorKV(gomock.Any(), gomock.All())

		r := httptest.NewRequest(http.MethodGet, "/test", bytes.NewReader([]byte{}))
		q := r.URL.Query()
		q.Add("service_name", reqS.ServiceName)
		q.Add("user_uid", reqS.UserUid.String())
		r.URL.RawQuery = q.Encode()
		w := httptest.NewRecorder()
		ta.a.GetSubscription(w, r)

		require.Equal(t, http.StatusInternalServerError, w.Code)

		respS := ds.GetSubscriptionResponse{}
		require.Nil(t, json.Unmarshal(w.Body.Bytes(), &respS))

		require.Equal(t, ds.StatusServiceError, respS.Status.Message)
	})
}

func TestUpdateSubscription(t *testing.T) {
	t.Parallel()

	t.Run("Ok", func(t *testing.T) {
		t.Parallel()

		ta := newTestAPI(t)

		reqS := ds.UpdateSubscriptionRequest{
			Subscription: ds.Subscription{
				Price: 100,
				SubscriptionID: ds.SubscriptionID{
					ServiceName: "name",
					UserUid:     uuid.New(),
				},
				Period: ds.Period{
					From: ds.DateType(time.Now()),
					To:   ds.DateType(time.Now()),
				},
			},
		}

		ta.subsMock.EXPECT().UpdateSubscription(gomock.Any()).Return(&ds.UpdateSubscriptionResponse{
			Status: ds.Status{Message: ds.StatusSuccess},
		})

		var buf bytes.Buffer
		_ = json.NewEncoder(&buf).Encode(&reqS)

		r := httptest.NewRequest(http.MethodPatch, "/test", bytes.NewReader(buf.Bytes()))
		w := httptest.NewRecorder()
		ta.a.UpdateSubscription(w, r)

		require.Equal(t, http.StatusOK, w.Code)

		respS := ds.UpdateSubscriptionResponse{}
		require.Nil(t, json.Unmarshal(w.Body.Bytes(), &respS))

		require.Equal(t, ds.StatusSuccess, respS.Status.Message)
	})

	t.Run("No required field", func(t *testing.T) {
		t.Parallel()

		ta := newTestAPI(t)

		reqS := ds.UpdateSubscriptionRequest{}

		ta.loggerMock.EXPECT().ErrorKV(gomock.Any(), gomock.All())

		var buf bytes.Buffer
		_ = json.NewEncoder(&buf).Encode(&reqS)
		body := buf.Bytes()

		r := httptest.NewRequest(http.MethodPatch, "/test", bytes.NewReader(body))
		w := httptest.NewRecorder()
		ta.a.UpdateSubscription(w, r)

		require.Equal(t, http.StatusBadRequest, w.Code)

		respS := ds.UpdateSubscriptionResponse{}
		require.Nil(t, json.Unmarshal(w.Body.Bytes(), &respS))

		require.Equal(t, ds.StatusFailedValidatingRequest, respS.Status.Message)
	})

	t.Run("No response from service", func(t *testing.T) {
		t.Parallel()

		ta := newTestAPI(t)

		reqS := ds.UpdateSubscriptionRequest{
			Subscription: ds.Subscription{
				Price: 100,
				SubscriptionID: ds.SubscriptionID{
					ServiceName: "name",
					UserUid:     uuid.New(),
				},
				Period: ds.Period{
					From: ds.DateType(time.Now()),
					To:   ds.DateType(time.Now()),
				},
			},
		}

		ta.subsMock.EXPECT().UpdateSubscription(gomock.Any()).Return(nil)
		ta.loggerMock.EXPECT().ErrorKV(gomock.Any(), gomock.All())

		var buf bytes.Buffer
		_ = json.NewEncoder(&buf).Encode(&reqS)
		body := buf.Bytes()

		r := httptest.NewRequest(http.MethodPatch, "/test", bytes.NewReader(body))
		w := httptest.NewRecorder()
		ta.a.UpdateSubscription(w, r)

		require.Equal(t, http.StatusInternalServerError, w.Code)

		respS := ds.UpdateSubscriptionResponse{}
		require.Nil(t, json.Unmarshal(w.Body.Bytes(), &respS))

		require.Equal(t, ds.StatusServiceError, respS.Status.Message)
	})
}

func TestDeleteSubscription(t *testing.T) {
	t.Parallel()

	t.Run("Ok", func(t *testing.T) {
		t.Parallel()

		ta := newTestAPI(t)

		reqS := ds.DeleteSubscriptionRequest{
			SubscriptionID: ds.SubscriptionID{
				ServiceName: "name",
				UserUid:     uuid.New(),
			},
		}

		ta.subsMock.EXPECT().DeleteSubscription(gomock.Any()).Return(&ds.DeleteSubscriptionResponse{
			Status: ds.Status{Message: ds.StatusSuccess},
		})

		var buf bytes.Buffer
		_ = json.NewEncoder(&buf).Encode(&reqS)

		r := httptest.NewRequest(http.MethodDelete, "/test", bytes.NewReader(buf.Bytes()))
		w := httptest.NewRecorder()
		ta.a.DeleteSubscription(w, r)

		require.Equal(t, http.StatusOK, w.Code)

		respS := ds.DeleteSubscriptionResponse{}
		require.Nil(t, json.Unmarshal(w.Body.Bytes(), &respS))

		require.Equal(t, ds.StatusSuccess, respS.Status.Message)
	})

	t.Run("No required field", func(t *testing.T) {
		t.Parallel()

		ta := newTestAPI(t)

		reqS := ds.DeleteSubscriptionRequest{}

		ta.loggerMock.EXPECT().ErrorKV(gomock.Any(), gomock.All())

		var buf bytes.Buffer
		_ = json.NewEncoder(&buf).Encode(&reqS)
		body := buf.Bytes()

		r := httptest.NewRequest(http.MethodDelete, "/test", bytes.NewReader(body))
		w := httptest.NewRecorder()
		ta.a.DeleteSubscription(w, r)

		require.Equal(t, http.StatusBadRequest, w.Code)

		respS := ds.DeleteSubscriptionResponse{}
		require.Nil(t, json.Unmarshal(w.Body.Bytes(), &respS))

		require.Equal(t, ds.StatusFailedValidatingRequest, respS.Status.Message)
	})

	t.Run("No response from service", func(t *testing.T) {
		t.Parallel()

		ta := newTestAPI(t)

		reqS := ds.DeleteSubscriptionRequest{
			SubscriptionID: ds.SubscriptionID{
				ServiceName: "name",
				UserUid:     uuid.New(),
			},
		}

		ta.subsMock.EXPECT().DeleteSubscription(gomock.Any()).Return(nil)
		ta.loggerMock.EXPECT().ErrorKV(gomock.Any(), gomock.All())

		var buf bytes.Buffer
		_ = json.NewEncoder(&buf).Encode(&reqS)
		body := buf.Bytes()

		r := httptest.NewRequest(http.MethodDelete, "/test", bytes.NewReader(body))
		w := httptest.NewRecorder()
		ta.a.DeleteSubscription(w, r)

		require.Equal(t, http.StatusInternalServerError, w.Code)

		respS := ds.DeleteSubscriptionResponse{}
		require.Nil(t, json.Unmarshal(w.Body.Bytes(), &respS))

		require.Equal(t, ds.StatusServiceError, respS.Status.Message)
	})
}

func TestListSubscriptions(t *testing.T) {
	t.Parallel()

	t.Run("Ok", func(t *testing.T) {
		t.Parallel()

		ta := newTestAPI(t)

		reqS := ds.ListSubscriptionsRequest{
			UserUid: uuid.New(),
		}

		ta.subsMock.EXPECT().ListSubscriptions(gomock.Any()).Return(&ds.ListSubscriptionsResponse{
			Status: ds.Status{Message: ds.StatusSuccess},
		})

		r := httptest.NewRequest(http.MethodGet, "/test", bytes.NewReader([]byte{}))
		q := r.URL.Query()
		q.Add("user_uid", reqS.UserUid.String())
		r.URL.RawQuery = q.Encode()
		w := httptest.NewRecorder()
		ta.a.ListSubscriptions(w, r)

		require.Equal(t, http.StatusOK, w.Code)

		respS := ds.ListSubscriptionsResponse{}
		require.Nil(t, json.Unmarshal(w.Body.Bytes(), &respS))

		require.Equal(t, ds.StatusSuccess, respS.Status.Message)
	})

	t.Run("No required field", func(t *testing.T) {
		t.Parallel()

		ta := newTestAPI(t)

		ta.loggerMock.EXPECT().ErrorKV(gomock.Any(), gomock.All())

		r := httptest.NewRequest(http.MethodGet, "/test", bytes.NewReader([]byte{}))
		w := httptest.NewRecorder()
		ta.a.ListSubscriptions(w, r)

		require.Equal(t, http.StatusBadRequest, w.Code)

		respS := ds.ListSubscriptionsResponse{}
		require.Nil(t, json.Unmarshal(w.Body.Bytes(), &respS))

		require.Equal(t, ds.StatusFailedValidatingRequest, respS.Status.Message)
	})

	t.Run("No response from service", func(t *testing.T) {
		t.Parallel()

		ta := newTestAPI(t)

		reqS := ds.ListSubscriptionsRequest{
			UserUid: uuid.New(),
		}

		ta.subsMock.EXPECT().ListSubscriptions(gomock.Any()).Return(nil)
		ta.loggerMock.EXPECT().ErrorKV(gomock.Any(), gomock.All())

		r := httptest.NewRequest(http.MethodGet, "/test", bytes.NewReader([]byte{}))
		q := r.URL.Query()
		q.Add("user_uid", reqS.UserUid.String())
		r.URL.RawQuery = q.Encode()
		w := httptest.NewRecorder()
		ta.a.ListSubscriptions(w, r)

		require.Equal(t, http.StatusInternalServerError, w.Code)

		respS := ds.ListSubscriptionsResponse{}
		require.Nil(t, json.Unmarshal(w.Body.Bytes(), &respS))

		require.Equal(t, ds.StatusServiceError, respS.Status.Message)
	})
}

func TestCalculateSubscriptionsPrice(t *testing.T) {
	t.Parallel()

	t.Run("Ok", func(t *testing.T) {
		t.Parallel()

		ta := newTestAPI(t)

		ta.subsMock.EXPECT().CalculateSubscriptionsPrice(gomock.Any()).Return(&ds.CalculateSubscriptionsPriceResponse{
			Status: ds.Status{Message: ds.StatusSuccess},
		})

		r := httptest.NewRequest(http.MethodGet, "/test", bytes.NewReader([]byte{}))
		q := r.URL.Query()
		q.Add("user_uid", uuid.New().String())
		q.Add("start_date", time.Now().Format(time.DateOnly))
		q.Add("end_date", time.Now().Format(time.DateOnly))
		r.URL.RawQuery = q.Encode()
		w := httptest.NewRecorder()
		ta.a.CalculateSubscriptionsPrice(w, r)

		require.Equal(t, http.StatusOK, w.Code)

		respS := ds.CalculateSubscriptionsPriceResponse{}
		require.Nil(t, json.Unmarshal(w.Body.Bytes(), &respS))

		require.Equal(t, ds.StatusSuccess, respS.Status.Message)
	})

	t.Run("No required field", func(t *testing.T) {
		t.Parallel()

		ta := newTestAPI(t)

		ta.loggerMock.EXPECT().ErrorKV(gomock.Any(), gomock.All())

		r := httptest.NewRequest(http.MethodGet, "/test", bytes.NewReader([]byte{}))
		w := httptest.NewRecorder()
		ta.a.CalculateSubscriptionsPrice(w, r)

		require.Equal(t, http.StatusBadRequest, w.Code)

		respS := ds.CalculateSubscriptionsPriceResponse{}
		require.Nil(t, json.Unmarshal(w.Body.Bytes(), &respS))

		require.Equal(t, ds.StatusFailedValidatingRequest, respS.Status.Message)
	})

	t.Run("No response from service", func(t *testing.T) {
		t.Parallel()

		ta := newTestAPI(t)

		ta.subsMock.EXPECT().CalculateSubscriptionsPrice(gomock.Any()).Return(nil)
		ta.loggerMock.EXPECT().ErrorKV(gomock.Any(), gomock.All())

		r := httptest.NewRequest(http.MethodGet, "/test", bytes.NewReader([]byte{}))
		q := r.URL.Query()
		q.Add("user_uid", uuid.New().String())
		q.Add("start_date", time.Now().Format(time.DateOnly))
		q.Add("end_date", time.Now().Format(time.DateOnly))
		r.URL.RawQuery = q.Encode()
		w := httptest.NewRecorder()
		ta.a.CalculateSubscriptionsPrice(w, r)

		require.Equal(t, http.StatusInternalServerError, w.Code)

		respS := ds.CalculateSubscriptionsPriceResponse{}
		require.Nil(t, json.Unmarshal(w.Body.Bytes(), &respS))

		require.Equal(t, ds.StatusServiceError, respS.Status.Message)
	})
}
