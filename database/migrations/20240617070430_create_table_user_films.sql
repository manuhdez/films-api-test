-- +goose Up
-- +goose StatementBegin
CREATE TABLE user_films
(
    user_id uuid NOT NULL,
    films int default 0
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS user_films;
-- +goose StatementEnd
