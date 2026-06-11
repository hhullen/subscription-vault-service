-- +goose Up
-- +goose StatementBegin

CREATE INDEX idx_subscriptions_cover 
ON subscriptions ("from", "to", user_uid, service_name) INCLUDE (price);

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin

DROP INDEX IF EXISTS idx_subscriptions_cover;

-- +goose StatementEnd