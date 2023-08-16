-- +goose Up
-- +goose StatementBegin
CREATE TABLE events (
    id BIGINT PRIMARY KEY,
    user_id BIGINT NOT NULL,
    ical_uid VARCHAR(255) NOT NULL,
    event_id VARCHAR(255) NOT NULL,
    title VARCHAR(255) NOT NULL,
    description TEXT NOT NULL,
    locations_count INT NOT NULL,
    start_time TIMESTAMP NOT NULL,
    end_time TIMESTAMP NOT NULL,
    is_online BOOL NOT NULL,
    is_all_day BOOL NOT NULL,
    is_cancelled BOOL NOT NULL,
    organizer_user_id BIGINT NOT NULL,
    created_time TIMESTAMP NOT NULL,
    updated_time TIMESTAMP NOT NULL,
    timezone TIMESTAMP NOT NULL,
    platform_url VARCHAR(255) NOT NULL,
    meeting_url VARCHAR(255) NOT NULL
);

CREATE TABLE locations (
    id BIGINT PRIMARY KEY,
    ical_uid VARCHAR(255) NOT NULL,
    display_name VARCHAR(255),
    locationUri VARCHAR(255),
    address VARCHAR(255)
);

CREATE TABLE attendees (
    id BIGINT PRIMARY KEY,
    user_id BIGINT NOT NULL
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE attendees;
DROP TABLE locations;
DROP TABLE events;
-- +goose StatementEnd