CREATE TABLE IF NOT EXISTS warehouse_items (
    id BIGSERIAL PRIMARY KEY,

    name VARCHAR(255) NOT NULL,
    sku VARCHAR(100) NOT NULL,
    manufacturer VARCHAR(150) NOT NULL DEFAULT '',

    purchase_price_cents BIGINT NOT NULL DEFAULT 0,
    sale_price_cents BIGINT NOT NULL DEFAULT 0,

    quantity INTEGER NOT NULL DEFAULT 0,
    min_quantity INTEGER NOT NULL DEFAULT 0,

    location VARCHAR(100) NOT NULL DEFAULT '',
    note TEXT NOT NULL DEFAULT '',

    is_active BOOLEAN NOT NULL DEFAULT TRUE,

    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    CONSTRAINT warehouse_items_name_not_empty
        CHECK (TRIM(name) <> ''),

    CONSTRAINT warehouse_items_sku_not_empty
        CHECK (TRIM(sku) <> ''),

    CONSTRAINT warehouse_items_purchase_price_not_negative
        CHECK (purchase_price_cents >= 0),

    CONSTRAINT warehouse_items_sale_price_not_negative
        CHECK (sale_price_cents >= 0),

    CONSTRAINT warehouse_items_quantity_not_negative
        CHECK (quantity >= 0),

    CONSTRAINT warehouse_items_min_quantity_not_negative
        CHECK (min_quantity >= 0)
);

CREATE UNIQUE INDEX IF NOT EXISTS warehouse_items_sku_unique
    ON warehouse_items (LOWER(TRIM(sku)));

CREATE INDEX IF NOT EXISTS warehouse_items_active_index
    ON warehouse_items (is_active);

CREATE INDEX IF NOT EXISTS warehouse_items_name_index
    ON warehouse_items (LOWER(name));
