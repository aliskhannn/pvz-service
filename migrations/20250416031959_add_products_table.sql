-- +goose Up
-- +goose StatementBegin
DO$$
BEGIN
    IF NOT EXISTS (
       SELECT 1 FROM pg_type WHERE typename = 'product_type'
    ) THEN
       CREATE TYPE product_type AS ENUM ('электроника', 'одежда', 'обувь');
    END IF;
END$$;

CREATE TABLE IF NOT EXISTS products (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    date_time TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    type product_type NOT NULL,
    reception_id UUID NOT NULL REFERENCES receptions(id) ON DELETE CASCADE,
)
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS products;

DROP TYPE IF EXISTS product_type;
-- +goose StatementEnd
