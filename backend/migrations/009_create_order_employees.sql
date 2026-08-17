CREATE TABLE order_employees (
    id BIGSERIAL PRIMARY KEY,

    order_id BIGINT NOT NULL
        REFERENCES orders(id)
        ON DELETE CASCADE,

    employee_id BIGINT NOT NULL
        REFERENCES employees(id),

    assigned_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    unassigned_at TIMESTAMPTZ
);

CREATE INDEX idx_order_employees_order_id
    ON order_employees(order_id);

CREATE INDEX idx_order_employees_employee_id
    ON order_employees(employee_id);
