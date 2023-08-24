-- +goose Up
-- +goose StatementBegin
ALTER TABLE users -- Added TABLE keyword here
ADD COLUMN subscription_id VARCHAR(255),
ADD COLUMN subscription_expires_at TIMESTAMP WITH TIME ZONE;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE users -- Added TABLE keyword here
DROP COLUMN subscription_id,
DROP COLUMN subscription_expires_at;
-- +goose StatementEnd

