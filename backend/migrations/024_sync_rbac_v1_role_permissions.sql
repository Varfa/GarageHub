-- =========================================================
-- Service Administrator
-- =========================================================

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
        -- Clients
        'clients.view',
        'clients.create',
        'clients.edit',

        -- Cars
        'cars.view',
        'cars.create',
        'cars.edit',

        -- Orders
        'orders.view',
        'orders.create',
        'orders.edit',
        'orders.change_status',
        'orders.assign_employee',
        'orders.add_note',
        'orders.change_price',

        -- Employees
        'employees.view',

        -- Warehouse
        'warehouse.view',
        'warehouse.manage',
        'warehouse.writeoff',

        -- Works
        'works.view',
        'works.create',
        'works.edit',

        -- Parts
        'parts.view',
        'parts.create',
        'parts.edit',
        'parts.manage',

        -- Reports
        'reports.view'
    )
WHERE r.code = 'service_admin'

ON CONFLICT (role_id, permission_id) DO NOTHING;


-- =========================================================
-- Mechanic
-- =========================================================

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
        'clients.view',
        'clients.create',
        'clients.edit',

        'cars.view',
        'cars.create',
        'cars.edit',

        'orders.view',
        'orders.create',
        'orders.edit',
        'orders.change_status',
        'orders.add_note',

        'employees.view',

        'warehouse.view',

        'works.view',
        'works.create',
        'works.edit',

        'parts.view',
        'parts.create'
    )
WHERE r.code = 'mechanic'

ON CONFLICT (role_id, permission_id) DO NOTHING;


-- =========================================================
-- Auto Electrician
-- =========================================================

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
        'clients.view',
        'clients.create',
        'clients.edit',

        'cars.view',
        'cars.create',
        'cars.edit',

        'orders.view',
        'orders.create',
        'orders.edit',
        'orders.change_status',
        'orders.add_note',

        'employees.view',

        'warehouse.view',

        'works.view',
        'works.create',
        'works.edit',

        'parts.view',
        'parts.create'
    )
WHERE r.code = 'auto_electrician'

ON CONFLICT (role_id, permission_id) DO NOTHING;


-- =========================================================
-- Tire Technician
-- =========================================================

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
        'clients.view',
        'clients.create',
        'clients.edit',

        'cars.view',
        'cars.create',
        'cars.edit',

        'orders.view',
        'orders.create',
        'orders.edit',
        'orders.change_status',
        'orders.add_note',

        'employees.view',

        'warehouse.view',

        'works.view',
        'works.create',
        'works.edit',

        'parts.view',
        'parts.create'
    )
WHERE r.code = 'tire_technician'

ON CONFLICT (role_id, permission_id) DO NOTHING;


-- =========================================================
-- Warehouse Manager
-- =========================================================

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
        'orders.view',

        'warehouse.view',
        'warehouse.manage',
        'warehouse.writeoff',

        'parts.view',
        'parts.create',
        'parts.edit',
        'parts.delete',
        'parts.manage',

        'employees.view'
    )
WHERE r.code = 'warehouse_manager'

ON CONFLICT (role_id, permission_id) DO NOTHING;


-- =========================================================
-- Accountant
-- =========================================================

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
        'orders.view',

        'warehouse.view',

        'reports.view',
        'reports.financial',

        'payroll.view',
        'payroll.manage',

        'accounting.view',
        'accounting.manage',

        'employees.view'
    )
WHERE r.code = 'accountant'

ON CONFLICT (role_id, permission_id) DO NOTHING;
