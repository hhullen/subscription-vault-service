package postgres

import (
	ds "subscription-vault-service/internal/datastruct"
)

func (c *Client) AddSubscription(*ds.CreateSubscriptionRequest) (*ds.CreateSubscriptionResponse, error) {
}

func (c *Client) GetSubscription(*ds.GetSubscriptionRequest) (*ds.GetSubscriptionResponse, error) {

}

func (c *Client) UpdateSubscription(*ds.UpdateSubscriptionRequest) (*ds.UpdateSubscriptionResponse, error) {
}

func (c *Client) DeleteSubscription(*ds.DeleteSubscriptionRequest) (*ds.DeleteSubscriptionResponse, error) {
}

func (c *Client) ListSubscriptions(*ds.ListSubscriptionsRequest) (*ds.ListSubscriptionsResponse, error) {
}

func (c *Client) CalculateSubscriptionsPrice(*ds.CalculateSubscriptionsPriceRequest) (*ds.CalculateSubscriptionsPriceResponse, error) {
}
