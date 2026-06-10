package service

import (
	"context"

	ds "subscription-vault-service/internal/datastruct"
)

//go:generate mockgen -source=service.go -destination=service_mock.go -package=service ILogger,ISubscriptionsStorage

type ILogger interface {
	InfoKV(message string, argsKV ...any)
	WarnKV(message string, argsKV ...any)
	ErrorKV(message string, argsKV ...any)
	FatalKV(message string, argsKV ...any)
}

type ISubscriptionsStorage interface {
	AddSubscription(*ds.CreateSubscriptionRequest) (*ds.CreateSubscriptionResponse, error)
	GetSubscription(*ds.GetSubscriptionRequest) (*ds.GetSubscriptionResponse, error)
	UpdateSubscription(*ds.UpdateSubscriptionRequest) (*ds.UpdateSubscriptionResponse, error)
	DeleteSubscription(*ds.DeleteSubscriptionRequest) (*ds.DeleteSubscriptionResponse, error)
	ListSubscriptions(*ds.ListSubscriptionsRequest) (*ds.ListSubscriptionsResponse, error)
	CalculateSubscriptionsPrice(*ds.CalculateSubscriptionsPriceRequest) (*ds.CalculateSubscriptionsPriceResponse, error)
}

type Service struct {
	ctx       context.Context
	subsStore ISubscriptionsStorage
	log       ILogger
}

func NewService(ctx context.Context, store ISubscriptionsStorage, log ILogger) *Service {
	return buildService(ctx, store, log)
}

func buildService(ctx context.Context, store ISubscriptionsStorage, log ILogger) *Service {
	return &Service{
		ctx:       ctx,
		subsStore: store,
		log:       log,
	}
}

func (s *Service) AddSubscription(*ds.CreateSubscriptionRequest) *ds.CreateSubscriptionResponse {

}

func (s *Service) GetSubscription(*ds.GetSubscriptionRequest) *ds.GetSubscriptionResponse {

}

func (s *Service) UpdateSubscription(*ds.UpdateSubscriptionRequest) *ds.UpdateSubscriptionResponse {

}

func (s *Service) DeleteSubscription(*ds.DeleteSubscriptionRequest) *ds.DeleteSubscriptionResponse {

}

func (s *Service) ListSubscriptions(*ds.ListSubscriptionsRequest) *ds.ListSubscriptionsResponse {

}

func (s *Service) CalculateSubscriptionsPrice(*ds.CalculateSubscriptionsPriceRequest) *ds.CalculateSubscriptionsPriceResponse {

}
