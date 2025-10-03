-- +goose Up
-- +goose StatementBegin
create table countries (
    id varchar(6) not null unique check ( id ~ '^[a-z]+$'),
    name varchar not null
);

create index countries_name_idx on countries(name);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
drop index if exists countries_lang_name_idx;
drop table countries;
-- +goose StatementEnd
