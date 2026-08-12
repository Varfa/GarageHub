CREATE TABLE orders (
    id BIGSERIAL PRIMARY KEY,

    client_id INTEGER NOT NULL,
    car_id INTEGER NOT NULL,

    complaint TEXT NOT NULL,
    diagnosis TEXT NOT NULL DEFAULT '',
    note TEXT NOT NULL DEFAULT '',

    status VARCHAR(50) NOT NULL DEFAULT 'new',

    estimated_cost_cents BIGINT NOT NULL DEFAULT 0,
    final_cost_cents BIGINT NOT NULL DEFAULT 0,

    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    started_at TIMESTAMPTZ,
    completed_at TIMESTAMPTZ,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    CONSTRAINT orders_client_fk
        FOREIGN KEY (client_id)
        REFERENCES clients(id),

    CONSTRAINT orders_car_fk
        FOREIGN KEY (car_id)
        REFERENCES cars(id),

    CONSTRAINT orders_estimated_cost_not_negative
        CHECK (estimated_cost_cents >= 0),

    CONSTRAINT orders_final_cost_not_negative
        CHECK (final_cost_cents >= 0)
);

CREATE INDEX orders_client_id_index
ON orders (client_id);

CREATE INDEX orders_car_id_index
ON orders (car_id);

CREATE INDEX orders_status_index
ON orders (status);
