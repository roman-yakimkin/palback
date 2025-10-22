-- +goose Up
-- +goose StatementBegin
alter table countries
add column weight integer not null default 0;

create index countries_weight_idx on countries(weight);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
drop index countries_weight_idx;
alter table countries drop column weight;
-- +goose StatementEnd
