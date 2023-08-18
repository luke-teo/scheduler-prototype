-- +goose Up
-- +goose StatementBegin
ALTER TABLE events 
ADD COLUMN is_recurring BOOL NOT NULL,
ADD COLUMN series_master_id VARCHAR(255);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE events
DROP COLUMN is_recurring,
DROP COLUMN series_master_id;
-- +goose StatementEnd
