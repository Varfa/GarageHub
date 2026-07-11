CREATE TABLE cars (
    id SERIAL PRIMARY KEY,
    client_id INTEGER NOT NULL REFERENCES clients(id) ON DELETE CASCADE,

    brand VARCHAR(100) NOT NULL,
    model VARCHAR(100) NOT NULL,
    year INTEGER,

    vin VARCHAR(50),
    plate_number VARCHAR(30),

    engine VARCHAR(100),
    power_kw INTEGER,

    color VARCHAR(50),
    mileage INTEGER,
    note TEXT,

    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW()
);
