package datastruct

import (
	"errors"
	"fmt"
	"reflect"
	"strings"
	"subscription-vault-service/internal/supports"
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

	DefaultSecretsDir          = "./secrets/"
	DefaultContainerSecretsDir = "/run/secrets/"
)

type DateType time.Time

func (d *DateType) Time() time.Time {
	return time.Time(*d)
}

func (d *DateType) UnmarshalJSON(b []byte) error {
	s := strings.Trim(string(b), "\"")

	dt, err := supports.ParseDate(s)
	if err != nil {
		return fmt.Errorf("incorrect date value: '%s'", s)
	}

	*d = DateType(dt)
	return nil
}

func (g *DateType) MarshalJSON() ([]byte, error) {
	return fmt.Appendf(nil, "\"%s\"", time.Time(*g).Format(time.DateOnly)), nil
}

func ParseSchemaDateType(value string) reflect.Value {
	if value == "" {
		return reflect.ValueOf(DateType{})
	}

	dt, err := supports.ParseDate(value)
	if err != nil {
		return reflect.ValueOf(DateType{})
	}

	return reflect.ValueOf(DateType(dt))
}

type Status struct {
	Message string `json:"status,omitempty" example:"status message"`
}

func (s Status) GetStatus() string {
	return s.Message
}

type SubscriptionID struct {
	ServiceName string    `json:"service_name" schema:"service_name" validate:"required" example:"Poople"`
	UserUid     uuid.UUID `json:"user_uid" schema:"user_uid" validate:"required" example:"4988150e-1c82-490f-8c07-ee74ace2dd14"`
}

type Period struct {
	From DateType `json:"start_date" schema:"start_date" validate:"required" example:"31.12.2006"`
	To   DateType `json:"end_date" schema:"end_date" validate:"required" example:"31.12.2007"`
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
	Subscription *Subscription `json:"subscription,omitempty"`
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
	UserUid uuid.UUID `schema:"user_uid" example:"4988150e-1c82-490f-8c07-ee74ace2dd14"`
}

type ListSubscriptionsResponse struct {
	Subscriptions []Subscription `json:"subscriptions,omitempty"`
	Status
}

type CalculateSubscriptionsPriceRequest struct {
	Period
	ServiceName *string    `schema:"service_name" example:"Poople"`
	UserUid     *uuid.UUID `schema:"user_uid" example:"4988150e-1c82-490f-8c07-ee74ace2dd14"`
}

type CalculateSubscriptionsPriceResponse struct {
	Status
	TotalPrice int64 `json:"total_price"`
}
