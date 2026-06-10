package datastruct

import (
	"errors"
	"time"

	"github.com/google/uuid"
)

var (
	ErrTest = errors.New("test")
)

const (
	StatusSuccess                  = "ok"
	StatusResurceNotFound          = "resource not found"
	StatusUserNotFound             = "user not found"
	StatusResourceAlreadyExists    = "resource already exists"
	StatusServiceError             = "service failed exec request"
	StatusDataTooLong              = "some data too long"
	StatusFailedExctractingRequest = "failed extracting request"
	StatusFailedValidatingRequest  = "failed validating request"
)

// type DateType time.Time

// func (d *DateType) UnmarshalJSON(b []byte) error {
// 	s := strings.Trim(string(b), "\"")

// 	dt, err := supports.ParseDate(s)
// 	if err != nil {
// 		return fmt.Errorf("incorrect date value: '%s'", s)
// 	}

// 	*d = DateType(dt)
// 	return nil
// }

// func (g *DateType) MarshalJSON() ([]byte, error) {
// 	return fmt.Appendf(nil, "\"%s\"", time.Time(*g).Format(time.DateOnly)), nil
// }

type Status struct {
	Message string `json:"status,omitempty" example:"status message"`
}

func (s Status) GetStatus() string {
	return s.Message
}

type SubscriptionID struct {
	ServiceName string    `json:"service_name" validate:"required" example:"Poople"`
	UserUid     uuid.UUID `json:"user_uid" validate:"required" example:"4988150e-1c82-490f-8c07-ee74ace2dd14"`
}

type Period struct {
	From time.Time `json:"from" validate:"required" example:"31.12.2006"`
	To   time.Time `json:"to" validate:"required" example:"31.12.2007"`
}

type Subscription struct {
	Period
	SubscriptionID
	Price int64 `json:"price" validate:"required" example:"399"`
}

type CreateSubscriptionRequest struct {
	Subscription
}

type CreateSubscriptionResponse struct {
	Status
}

type GetSubscriptionRequest struct {
	SubscriptionID
}

type GetSubscriptionResponse struct {
	Subscriptions Subscription
	Status
}

type UpdateSubscriptionRequest struct {
	Subscription
}

type UpdateSubscriptionResponse struct {
	Status
}

type DeleteSubscriptionRequest struct {
	SubscriptionID
}

type DeleteSubscriptionResponse struct {
	Status
}

type ListSubscriptionsRequest struct {
	UserUid uuid.UUID `json:"user_uid" example:"4988150e-1c82-490f-8c07-ee74ace2dd14"`
}

type ListSubscriptionsResponse struct {
	Subscriptions []Subscription
	Status
}

type CalculateSubscriptionsPriceRequest struct {
	Period
	ServiceName *string    `json:"service_name" example:"Poople"`
	UserUid     *uuid.UUID `json:"user_uid" example:"4988150e-1c82-490f-8c07-ee74ace2dd14"`
}

type CalculateSubscriptionsPriceResponse struct {
	Status
	TotalPrice int64
}
