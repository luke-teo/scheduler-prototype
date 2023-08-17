-- +goose Up
-- +goose StatementBegin
ALTER TABLE attendees
ADD COLUMN name varchar(255),
ADD COLUMN email_address varchar(255) NOT NULL,
ADD COLUMN ical_uid varchar(255) NOT NULL;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE attendees
DROP COLUMN name,
DROP COLUMN email_address,
DROP COLUMN ical_uid;
-- +goose StatementEnd
