-- +goose Up
-- +goose StatementBegin
alter table countries add column has_regions bool not null default false;

create table regions (
    id serial primary key,
    country_id varchar(6) not null,
    name varchar not null,
    constraint fk_country foreign key (country_id) references countries(id) on update cascade,
    constraint fk_unique_country_and_name unique (country_id, name)
);

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
drop table regions;

alter table countries drop column has_regions;
-- +goose StatementEnd
