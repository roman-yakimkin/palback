-- +goose Up
-- +goose StatementBegin
CREATE TABLE languages (
    id char(2) unique
);

INSERT INTO languages VALUES ('ru');
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DELETE FROM languages;
DROP TABLE languages;
-- +goose StatementEnd
