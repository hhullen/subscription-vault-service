package postgres

import (
	"subscription-vault-service/internal/clients/postgres/sqlc"
	ds "subscription-vault-service/internal/datastruct"
)

func (c *Client) AddSubscription(req *ds.CreateSubscriptionRequest) (*ds.CreateSubscriptionResponse, error) {
	ctx, cancel := c.db.CtxWithCancel()
	cancel()

	err := c.db.Querier().CreateSubscription(ctx, sqlc.CreateSubscriptionParams{
		UserUid:     req.UserUid,
		ServiceName: req.ServiceName,
		Price:       req.Price,
		From:        nullTime(&req.From),
		To:          nullTime(&req.To),
	})
	if err != nil {
		if isDuplicate(err) {
			return &ds.CreateSubscriptionResponse{
				Status: ds.Status{Message: ds.StatusResourceAlreadyExists},
			}, nil
		}
		return nil, err
	}

	return &ds.CreateSubscriptionResponse{
		Status: ds.Status{Message: ds.StatusSuccess},
	}, nil
}

func (c *Client) GetSubscription(req *ds.GetSubscriptionRequest) (*ds.GetSubscriptionResponse, error) {
	ctx, cancel := c.db.CtxWithCancel()
	cancel()

	res, err := c.db.Querier().GetSubscription(ctx, sqlc.GetSubscriptionParams{
		UserUid:     req.UserUid,
		ServiceName: req.ServiceName,
	})
	if err != nil {
		if isNoRows(err) {
			return &ds.GetSubscriptionResponse{
				Status: ds.Status{Message: ds.StatusResurceNotFound},
			}, nil
		}
		return nil, err
	}

	return &ds.GetSubscriptionResponse{
		Subscriptions: ds.Subscription{
			Period: ds.Period{
				From: res.From.Time,
				To:   res.To.Time,
			},
			SubscriptionID: ds.SubscriptionID{
				ServiceName: res.ServiceName,
				UserUid:     res.UserUid,
			},
			Price: res.Price,
		},
		Status: ds.Status{Message: ds.StatusSuccess},
	}, nil
}

func (c *Client) UpdateSubscription(req *ds.UpdateSubscriptionRequest) (*ds.UpdateSubscriptionResponse, error) {
	ctx, cancel := c.db.CtxWithCancel()
	cancel()

	err := c.db.Querier().UpdateSubscription(ctx, sqlc.UpdateSubscriptionParams{
		UserUid:     req.UserUid,
		ServiceName: req.ServiceName,
		Price:       req.Price,
		From:        nullTime(&req.From),
		To:          nullTime(&req.To),
	})
	if err != nil {
		if isNoRows(err) {
			return &ds.UpdateSubscriptionResponse{
				Status: ds.Status{Message: ds.StatusResurceNotFound},
			}, nil
		}
		return nil, err
	}

	return &ds.UpdateSubscriptionResponse{
		Status: ds.Status{Message: ds.StatusSuccess},
	}, nil
}

func (c *Client) DeleteSubscription(req *ds.DeleteSubscriptionRequest) (*ds.DeleteSubscriptionResponse, error) {
	ctx, cancel := c.db.CtxWithCancel()
	cancel()

	err := c.db.Querier().DeleteSubscription(ctx, sqlc.DeleteSubscriptionParams{
		UserUid:     req.UserUid,
		ServiceName: req.ServiceName,
	})
	if err != nil {
		if isNoRows(err) {
			return &ds.DeleteSubscriptionResponse{
				Status: ds.Status{Message: ds.StatusResurceNotFound},
			}, nil
		}
		return nil, err
	}

	return &ds.DeleteSubscriptionResponse{
		Status: ds.Status{Message: ds.StatusSuccess},
	}, nil
}

func (c *Client) ListSubscriptions(req *ds.ListSubscriptionsRequest) (*ds.ListSubscriptionsResponse, error) {
	ctx, cancel := c.db.CtxWithCancel()
	cancel()

	res, err := c.db.Querier().ListSubscriptions(ctx, req.UserUid)
	if err != nil {
		return nil, err
	}

	ret := &ds.ListSubscriptionsResponse{}
	ret.Subscriptions = make([]ds.Subscription, len(res))

	for i := range len(res) {
		ret.Subscriptions[i].From = res[i].From.Time
		ret.Subscriptions[i].To = res[i].From.Time
		ret.Subscriptions[i].Price = res[i].Price
		ret.Subscriptions[i].ServiceName = res[i].ServiceName
		ret.Subscriptions[i].UserUid = res[i].UserUid
	}

	return ret, nil
}

func (c *Client) CalculateSubscriptionsPrice(req *ds.CalculateSubscriptionsPriceRequest) (*ds.CalculateSubscriptionsPriceResponse, error) {
	ctx, cancel := c.db.CtxWithCancel()
	cancel()

	sum, err := c.db.Querier().CalculateSubscriptionsPrice(ctx, sqlc.CalculateSubscriptionsPriceParams{
		From:        nullTime(&req.From),
		To:          nullTime(&req.To),
		ServiceName: nullString(req.ServiceName),
		UserUid:     nullUUID(req.UserUid),
	})
	if err != nil {
		return nil, err
	}

	return &ds.CalculateSubscriptionsPriceResponse{
		TotalPrice: sum,
		Status:     ds.Status{Message: ds.StatusSuccess},
	}, nil
}
