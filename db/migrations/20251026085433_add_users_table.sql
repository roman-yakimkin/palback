-- +goose Up
-- +goose StatementBegin
create table users (
    id SERIAL PRIMARY KEY,
    role_id varchar(20) NOT NULL default 'user',
    username TEXT UNIQUE NOT NULL,
    email TEXT UNIQUE NOT NULL,
    password TEXT NOT NULL,
    email_verified BOOLEAN DEFAULT false,
    session_version BIGINT NOT NULL default 0,
    created_at TIMESTAMP DEFAULT NOW()
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
drop table users;
-- +goose StatementEnd
