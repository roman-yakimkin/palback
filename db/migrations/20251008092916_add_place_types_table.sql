-- +goose Up
-- +goose StatementBegin
create table place_types (
    id serial primary key,
    name varchar not null unique,
    weight int
);

create index place_types_weight_idx on place_types(weight);

insert into place_types values (1, 'храм', 10),
                               (2, 'собор', 20),
                               (3, 'часовня', 30),
                               (4, 'монастырь', 40),
                               (5, 'скит', 50),
                               (6, 'лавра', 60),
                               (7, 'гора', 70),
                               (8, 'могила', 80),
                               (9, 'пещера', 90),
                               (10, 'иные', 10000);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
drop index place_types_weight_idx;

drop table place_types;
-- +goose StatementEnd
