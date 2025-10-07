-- +goose Up
-- +goose StatementBegin
create table city_types (
    id serial primary key,
    name varchar not null unique,
    short_name varchar,
    weight int
);

create index city_types_weight_idx on city_types(weight);

insert into city_types values (1, 'город', '', 10),
                             (2, 'село', 'с.', 20),
                             (3, 'посёлок', 'пос.', 30),
                             (4, 'хутор', 'х.', 40),
                             (5, 'станица', 'ст.', 50),
                             (6, 'аул', 'аул', 60),
                             (7, 'кишлак', 'к.', 70),
                             (8, 'местечко', 'мест.', 80),
                             (9, 'коммуна', 'комм.', 90),
                             (10, 'улус', 'ул.', 100),
                             (11, 'выселки', 'выс.', 110),
                             (12, 'погост', 'пог.', 120);

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin


drop index city_types_weight_idx;

drop table city_types;
-- +goose StatementEnd
