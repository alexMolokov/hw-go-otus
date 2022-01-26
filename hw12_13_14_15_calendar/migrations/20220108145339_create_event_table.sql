-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS calendar.event (
    event_id SERIAL PRIMARY KEY,
    title VARCHAR(100) NOT NULL,
    description VARCHAR,
    start_date TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    end_date TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    owner_id INT NOT NULL REFERENCES calendar.user(user_id),
    notify_for smallint default 1,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE
    );

create index owner_id_idx on calendar.event (owner_id);
create index start_end_idx on calendar.event (start_date, end_date);

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS calendar.event;
-- +goose StatementEnd