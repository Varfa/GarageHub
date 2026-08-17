INSERT INTO role_permissions (
    role_id,
    permission_id
)
SELECT
    r.id,
    p.id
FROM roles r
JOIN permissions p
    ON p.code = 'orders.restore'
WHERE r.code IN (
    'service_admin',
    'mechanic',
    'auto_electrician',
    'tire_technician'
)
ON CONFLICT (role_id, permission_id) DO NOTHING;
