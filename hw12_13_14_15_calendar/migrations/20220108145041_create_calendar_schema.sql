-- +goose Up
-- +goose StatementBegin
CREATE SCHEMA IF NOT EXISTS calendar;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP SCHEMA IF EXISTS calendar;
-- +goose StatementEnd
