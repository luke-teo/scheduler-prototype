-- +goose Up
-- +goose StatementBegin
CREATE TABLE events (
    id UUID DEFAULT gen_random_uuid() PRIMARY KEY,
    user_id BIGINT NOT NULL,
    ical_uid VARCHAR(255) NOT NULL,
    event_id VARCHAR(255) NOT NULL,
    title VARCHAR(255) NOT NULL,
    description TEXT NOT NULL,
    locations_count INT NOT NULL,
    start_time TIMESTAMP WITH TIME ZONE NOT NULL,
    end_time TIMESTAMP WITH TIME ZONE NOT NULL,
    is_online BOOL NOT NULL,
    is_all_day BOOL NOT NULL,
    is_cancelled BOOL NOT NULL,
    organizer_user_id BIGINT NOT NULL,
    created_time TIMESTAMP WITH TIME ZONE NOT NULL,
    updated_time TIMESTAMP WITH TIME ZONE NOT NULL,
    timezone VARCHAR(255) NOT NULL,
    platform_url VARCHAR(255) NOT NULL,
    meeting_url VARCHAR(255) NOT NULL
);

CREATE TABLE locations (
    id UUID DEFAULT gen_random_uuid() PRIMARY KEY,
    ical_uid VARCHAR(255) NOT NULL,
    display_name VARCHAR(255),
    location_uri VARCHAR(255),
    address VARCHAR(255)
);

CREATE TABLE attendees (
    id UUID DEFAULT gen_random_uuid() PRIMARY KEY,
    user_id BIGINT 
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE attendees;
DROP TABLE locations;
DROP TABLE events;
-- +goose StatementEnd
