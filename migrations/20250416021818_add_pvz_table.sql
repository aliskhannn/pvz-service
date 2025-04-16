-- +goose Up
-- +goose StatementBegin
DO $$
BEGIN
    IF NOT EXISTS (
        SELECT 1 FROM pg_type WHERE typname = 'city'
    ) THEN
        CREATE TYPE city AS ENUM ('Москва', 'Санкт-Петербург', 'Казань');
    END IF;
END$$;

CREATE TABLE IF NOT EXISTS pvz
(
    id                UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    registration_date TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    city              city NOT NULL
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS pvz;

DROP TYPE IF EXISTS city;
-- +goose StatementEnd
