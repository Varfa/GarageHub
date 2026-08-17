CREATE TABLE users (
    id BIGSERIAL PRIMARY KEY,

    employee_id BIGINT,
    role_id BIGINT,

    email VARCHAR(255) NOT NULL UNIQUE,
    password_hash TEXT NOT NULL,

    is_owner BOOLEAN NOT NULL DEFAULT FALSE,
    is_active BOOLEAN NOT NULL DEFAULT TRUE,

    last_login_at TIMESTAMPTZ,

    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    CONSTRAINT users_employee_fk
        FOREIGN KEY (employee_id)
        REFERENCES employees(id)
        ON DELETE SET NULL,

    CONSTRAINT users_role_fk
        FOREIGN KEY (role_id)
        REFERENCES roles(id)
        ON DELETE SET NULL,

    CONSTRAINT users_employee_unique
        UNIQUE (employee_id),

    CONSTRAINT users_owner_role_check
        CHECK (
            is_owner = FALSE
            OR role_id IS NULL
        )
);

CREATE UNIQUE INDEX users_email_lower_unique
    ON users (LOWER(email));

CREATE INDEX users_role_id_index
    ON users (role_id);

CREATE INDEX users_active_index
    ON users (is_active);
