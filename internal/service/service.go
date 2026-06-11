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

func (s *Service) AddSubscription(req *ds.CreateSubscriptionRequest) *ds.CreateSubscriptionResponse {
	resp, err := s.subsStore.AddSubscription(req)
	if err != nil {
		s.log.ErrorKV("AddSubscription.AddSubscription", "error", err.Error())
		return nil
	}

	return resp
}

func (s *Service) GetSubscription(req *ds.GetSubscriptionRequest) *ds.GetSubscriptionResponse {
	resp, err := s.subsStore.GetSubscription(req)
	if err != nil {
		s.log.ErrorKV("GetSubscription.GetSubscription", "error", err.Error())
		return nil
	}

	return resp
}

func (s *Service) UpdateSubscription(req *ds.UpdateSubscriptionRequest) *ds.UpdateSubscriptionResponse {
	resp, err := s.subsStore.UpdateSubscription(req)
	if err != nil {
		s.log.ErrorKV("UpdateSubscription.UpdateSubscription", "error", err.Error())
		return nil
	}

	return resp
}

func (s *Service) DeleteSubscription(req *ds.DeleteSubscriptionRequest) *ds.DeleteSubscriptionResponse {
	resp, err := s.subsStore.DeleteSubscription(req)
	if err != nil {
		s.log.ErrorKV("DeleteSubscription.DeleteSubscription", "error", err.Error())
		return nil
	}

	return resp
}

func (s *Service) ListSubscriptions(req *ds.ListSubscriptionsRequest) *ds.ListSubscriptionsResponse {
	resp, err := s.subsStore.ListSubscriptions(req)
	if err != nil {
		s.log.ErrorKV("ListSubscriptions.ListSubscriptions", "error", err.Error())
		return nil
	}

	return resp
}

func (s *Service) CalculateSubscriptionsPrice(req *ds.CalculateSubscriptionsPriceRequest) *ds.CalculateSubscriptionsPriceResponse {
	resp, err := s.subsStore.CalculateSubscriptionsPrice(req)
	if err != nil {
		s.log.ErrorKV("CalculateSubscriptionsPrice.CalculateSubscriptionsPrice", "error", err.Error())
		return nil
	}

	return resp
}
