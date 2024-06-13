-- +goose Up
-- +goose StatementBegin
SELECT 'up SQL query';
-- +goose StatementEnd
ALTER TABLE calculations
ALTER COLUMN a TYPE BIGINT USING a::BIGINT,
ALTER COLUMN b TYPE BIGINT USING a::BIGINT;
-- +goose Down
-- +goose StatementBegin
SELECT 'down SQL query';

ALTER TABLE calculations
ALTER COLUMN a TYPE INTEGER USING a::INTEGER,
ALTER COLUMN b TYPE INTEGER USING b::INTEGER;
-- +goose StatementEnd
