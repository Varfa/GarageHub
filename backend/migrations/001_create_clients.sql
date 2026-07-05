CREATE TABLE clients (
    id INTEGER GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    number INTEGER NOT NULL,

    name VARCHAR(100) NOT NULL,
    phone VARCHAR(30) NOT NULL,
    email VARCHAR(255),

    address TEXT,
    note TEXT,

    last_visit_at TIMESTAMP,

    created_at TIMESTAMP,
    updated_at TIMESTAMP
);
