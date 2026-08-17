INSERT INTO role_permissions (
    role_id,
    permission_id
)
SELECT
    r.id,
    p.id
FROM roles r
JOIN permissions p
    ON p.code IN (
        -- Orders
        'orders.view',
        'orders.create',
        'orders.edit',
        'orders.change_status',
        'orders.assign_employee',
        'orders.add_note',
        'orders.change_price',

        -- Clients
        'clients.view',
        'clients.create',
        'clients.edit',

        -- Cars
        'cars.view',
        'cars.create',
        'cars.edit',

        -- Employees
        'employees.view',

        -- Warehouse
        'warehouse.view',

        -- Works
        'works.view',
        'works.create',
        'works.edit',

        -- Parts / materials
        'parts.view',
        'parts.create',
        'parts.edit',

        -- Reports
        'reports.view'
    )
WHERE r.code = 'service_admin'

ON CONFLICT (role_id, permission_id) DO NOTHING;
