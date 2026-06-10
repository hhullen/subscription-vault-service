-- name: CreateSubscription :one
INSERT INTO subscriptions 
(user_uid, service_name, price, "from", "to") VALUES
($1, $2, $3, $4, $5)
ON CONFLICT (user_uid, service_name)
DO NOTHING
RETURNING user_uid;