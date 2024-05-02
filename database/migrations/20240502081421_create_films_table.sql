-- +goose Up
-- +goose StatementBegin
CREATE TABLE films
(
    id           uuid    NOT NULL,
    title        text    NOT NULl unique,
    director     text    NOT NULL,
    release_date integer NOT NULL,
    casting      text[] DEFAULT '{}'::text[],
    genre        text,
    synopsis     text,
    created_by   uuid    NOT NULL REFERENCES users (id),
    PRIMARY KEY (id)
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS films;
-- +goose StatementEnd
