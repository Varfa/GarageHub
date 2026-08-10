ALTER TABLE employee_positions
ADD COLUMN code VARCHAR(100);

UPDATE employee_positions
SET code = LOWER(
    REGEXP_REPLACE(
        name,
        '[^a-zA-Z0-9]+',
        '_',
        'g'
    )
)
WHERE code IS NULL;

ALTER TABLE employee_positions
ALTER COLUMN code SET NOT NULL;

ALTER TABLE employee_positions
ADD CONSTRAINT employee_positions_code_unique
UNIQUE (code);
