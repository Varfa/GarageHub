CREATE UNIQUE INDEX order_employees_active_unique
ON order_employees (order_id, employee_id)
WHERE unassigned_at IS NULL;
