-- +goose Up
-- +goose StatementBegin
CREATE TABLE users (
  id uuid not null,
  username text not null unique,
  password text not null,
  primary key (id)
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
drop table if exists users;
-- +goose StatementEnd
