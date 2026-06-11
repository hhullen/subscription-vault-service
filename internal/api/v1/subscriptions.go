package api

import (
	"net/http"

	ds "subscription-vault-service/internal/datastruct"
)

const (
	subscriptionsPrefix          = apiPrefix + "/subscriptions"
	listSubscriptionsPrefix      = subscriptionsPrefix + "/list"
	calculateSubscriptionsPrefix = subscriptionsPrefix + "/total_price"
)

func (a *API) setupSubscriptionHandlers() {
	a.router.Handle(pattern(http.MethodPost, subscriptionsPrefix), http.HandlerFunc(a.CreateSubscription))
	a.router.Handle(pattern(http.MethodGet, subscriptionsPrefix), http.HandlerFunc(a.GetSubscription))
	a.router.Handle(pattern(http.MethodPatch, subscriptionsPrefix), http.HandlerFunc(a.UpdateSubscription))
	a.router.Handle(pattern(http.MethodDelete, subscriptionsPrefix), http.HandlerFunc(a.DeleteSubscription))
	a.router.Handle(pattern(http.MethodGet, listSubscriptionsPrefix), http.HandlerFunc(a.ListSubscriptions))
	a.router.Handle(pattern(http.MethodGet, calculateSubscriptionsPrefix), http.HandlerFunc(a.CalculateSubscriptionsPrice))
}

// CreateSubscription Create new subscription
// @Summary      Create new subscription
// @Description  Create new subscription.
// @Tags         Subscriptions
// @Accept       json
// @Produce      json
// @Param        input body      ds.CreateSubscriptionRequest  true "Subscription"
// @Success      200   {object}  ds.CreateSubscriptionResponse
// @Failure      400   {object}  ds.Status
// @Failure      500   {object}  ds.Status
// @Router       /subscriptions [post]
func (a *API) CreateSubscription(w http.ResponseWriter, r *http.Request) {
	Exec(ExecArgs[ds.CreateSubscriptionRequest, ds.CreateSubscriptionResponse]{
		serviceFunc:      a.subsService.AddSubscription,
		responseWriter:   writeJsonResponse[*ds.CreateSubscriptionResponse],
		requestExtractor: extractJsonBody[*ds.CreateSubscriptionRequest],
		httpResponse:     w,
		httpRequest:      r,
		api:              a,
	})
}

// GetSubscription Get subscription
// @Summary      Get subscription
// @Description  Get subscription.
// @Tags         Subscriptions
// @Produce      json
// @Param        service_name   query     string                     true "service_name" example(Poople)
// @Param        user_uid       query     string                     true "user_uid"     example(4988150e-1c82-490f-8c07-ee74ace2dd14)
// @Success      200            {object}  ds.GetSubscriptionResponse
// @Failure      400            {object}  ds.Status
// @Failure      500            {object}  ds.Status
// @Router       /subscriptions  [get]
func (a *API) GetSubscription(w http.ResponseWriter, r *http.Request) {
	Exec(ExecArgs[ds.GetSubscriptionRequest, ds.GetSubscriptionResponse]{
		serviceFunc:      a.subsService.GetSubscription,
		responseWriter:   writeJsonResponse[*ds.GetSubscriptionResponse],
		requestExtractor: extractSchemaQuery[*ds.GetSubscriptionRequest],
		httpResponse:     w,
		httpRequest:      r,
		api:              a,
	})
}

// UpdateSubscription Update subscription
// @Summary      Update subscription
// @Description  Update subscription.
// @Tags         Subscriptions
// @Accept       json
// @Produce      json
// @Param        input body      ds.UpdateSubscriptionRequest  true "Subscription"
// @Success      200   {object}  ds.UpdateSubscriptionResponse
// @Failure      400   {object}  ds.Status
// @Failure      500   {object}  ds.Status
// @Router       /subscriptions [patch]
func (a *API) UpdateSubscription(w http.ResponseWriter, r *http.Request) {
	Exec(ExecArgs[ds.UpdateSubscriptionRequest, ds.UpdateSubscriptionResponse]{
		serviceFunc:      a.subsService.UpdateSubscription,
		responseWriter:   writeJsonResponse[*ds.UpdateSubscriptionResponse],
		requestExtractor: extractJsonBody[*ds.UpdateSubscriptionRequest],
		httpResponse:     w,
		httpRequest:      r,
		api:              a,
	})
}

// DeleteSubscription Delete subscription
// @Summary      Delete subscription
// @Description  Delete subscription.
// @Tags         Subscriptions
// @Accept       json
// @Produce      json
// @Param        input body      ds.DeleteSubscriptionRequest  true "Subscription"
// @Success      200   {object}  ds.DeleteSubscriptionResponse
// @Failure      400   {object}  ds.Status
// @Failure      500   {object}  ds.Status
// @Router       /subscriptions [delete]
func (a *API) DeleteSubscription(w http.ResponseWriter, r *http.Request) {
	Exec(ExecArgs[ds.DeleteSubscriptionRequest, ds.DeleteSubscriptionResponse]{
		serviceFunc:      a.subsService.DeleteSubscription,
		responseWriter:   writeJsonResponse[*ds.DeleteSubscriptionResponse],
		requestExtractor: extractJsonBody[*ds.DeleteSubscriptionRequest],
		httpResponse:     w,
		httpRequest:      r,
		api:              a,
	})
}

// ListSubscriptions List subscriptions
// @Summary      List subscriptions
// @Description  List subscriptions.
// @Tags         Subscriptions
// @Produce      json
// @Param        user_uid       query     string                       true "user_uid" example(4988150e-1c82-490f-8c07-ee74ace2dd14)
// @Success      200            {object}  ds.ListSubscriptionsResponse
// @Failure      400            {object}  ds.Status
// @Failure      500            {object}  ds.Status
// @Router       /subscriptions/list  [get]
func (a *API) ListSubscriptions(w http.ResponseWriter, r *http.Request) {
	Exec(ExecArgs[ds.ListSubscriptionsRequest, ds.ListSubscriptionsResponse]{
		serviceFunc:      a.subsService.ListSubscriptions,
		responseWriter:   writeJsonResponse[*ds.ListSubscriptionsResponse],
		requestExtractor: extractSchemaQuery[*ds.ListSubscriptionsRequest],
		httpResponse:     w,
		httpRequest:      r,
		api:              a,
	})
}

// CalculateSubscriptionsPrice Calculate price of subscriptions
// @Summary      Calculate price of subscriptions
// @Description  Calculate price of subscriptions.
// @Tags         Subscriptions
// @Produce      json
// @Param        start_date     query     string                                  true  "start_date"   example(31.12.2006)
// @Param        end_date       query     string                                  true  "end_date"     example(31.12.2008)
// @Param        service_name   query     string                                  false "service_name" example(Poople)
// @Param        user_uid       query     string                                  false "user_uid"     example(4988150e-1c82-490f-8c07-ee74ace2dd14)
// @Success      200            {object}  ds.CalculateSubscriptionsPriceResponse
// @Failure      400            {object}  ds.Status
// @Failure      500            {object}  ds.Status
// @Router       /subscriptions/total_price  [get]
func (a *API) CalculateSubscriptionsPrice(w http.ResponseWriter, r *http.Request) {
	Exec(ExecArgs[ds.CalculateSubscriptionsPriceRequest, ds.CalculateSubscriptionsPriceResponse]{
		serviceFunc:      a.subsService.CalculateSubscriptionsPrice,
		responseWriter:   writeJsonResponse[*ds.CalculateSubscriptionsPriceResponse],
		requestExtractor: extractSchemaQuery[*ds.CalculateSubscriptionsPriceRequest],
		httpResponse:     w,
		httpRequest:      r,
		api:              a,
	})
}
