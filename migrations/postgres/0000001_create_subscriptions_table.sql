-- +goose Up
-- +goose StatementBegin

CREATE TABLE IF NOT EXISTS subscriptions (
    id SERIAL PRIMARY KEY,
    user_uid UUID NOT NULL,
    service_name TEXT NOT NULL,
    price BIGINT NOT NULL,
    "from" TIMESTAMPTZ DEFAULT NULL,
    "to" TIMESTAMPTZ DEFAULT NULL,

    UNIQUE(user_uid, service_name)
);

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin

DROP TABLE IF EXISTS subscriptions;

-- +goose StatementEnd