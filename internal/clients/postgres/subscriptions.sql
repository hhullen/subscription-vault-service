-- name: CreateSubscription :exec
INSERT INTO subscriptions 
(user_uid, service_name, price, "from", "to") VALUES
($1, $2, $3, $4, $5);

-- name: GetSubscription :one
SELECT user_uid, service_name, price, "from", "to"
FROM subscriptions
WHERE user_uid = $1
    AND service_name = $2;

-- name: UpdateSubscription :execresult
UPDATE subscriptions
SET price = $1,
    "from" = $2,
    "to" = $3
WHERE user_uid = $4
    AND service_name = $5;

-- name: DeleteSubscription :execresult
DELETE FROM subscriptions
WHERE user_uid = $1
    AND service_name = $2;

-- name: ListSubscriptions :many
SELECT user_uid, service_name, price, "from", "to"
FROM subscriptions
WHERE user_uid = $1;

-- name: CalculateSubscriptionsPrice :one
SELECT COALESCE(SUM(price), 0)::BIGINT AS total_sum
FROM subscriptions
WHERE "from" >= sqlc.arg(from_date)
    AND "from" <= sqlc.arg(to_date)
    AND (user_uid = sqlc.narg(user_uid) OR sqlc.narg(user_uid) IS NULL)
    AND (service_name = sqlc.narg(service_name) OR sqlc.narg(service_name) IS NULL);
