-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS calendar.user
(user_id   SERIAL PRIMARY KEY NOT NULL,
 first_name VARCHAR(50),
 last_name  VARCHAR(50),
 email      VARCHAR(50) NOT NULL UNIQUE
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS calendar.user
-- +goose StatementEnd
