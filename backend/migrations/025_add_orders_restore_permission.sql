INSERT INTO permissions (
    code,
    name,
    description,
    module
)
VALUES (
    'orders.restore',
    'Restore orders',
    'Allows restoring closed orders',
    'orders'
)
ON CONFLICT (code) DO NOTHING;
