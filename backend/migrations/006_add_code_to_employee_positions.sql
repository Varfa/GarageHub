ALTER TABLE employee_positions
ADD COLUMN code VARCHAR(100);

UPDATE employee_positions
SET code = CASE
    WHEN name = 'Механик' THEN 'mechanic'
    WHEN name = 'Автоэлектрик' THEN 'auto_electrician'
    WHEN name = 'Мастер-приёмщик' THEN 'service_advisor'
    WHEN name = 'Администратор' THEN 'administrator'
    ELSE 'position_' || id
END
WHERE code IS NULL;

ALTER TABLE employee_positions
ALTER COLUMN code SET NOT NULL;

ALTER TABLE employee_positions
ADD CONSTRAINT employee_positions_code_unique
UNIQUE (code);
