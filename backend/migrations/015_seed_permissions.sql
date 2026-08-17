INSERT INTO permissions (
    code,
    name,
    description,
    module
)
VALUES
    -- Orders
    (
        'orders.view',
        'View orders',
        'Allows viewing order lists and order details',
        'orders'
    ),
    (
        'orders.create',
        'Create orders',
        'Allows creating new orders',
        'orders'
    ),
    (
        'orders.edit',
        'Edit orders',
        'Allows editing order information',
        'orders'
    ),
    (
        'orders.change_status',
        'Change order status',
        'Allows changing order status',
        'orders'
    ),
    (
        'orders.assign_employee',
        'Assign employees to orders',
        'Allows assigning and unassigning employees on orders',
        'orders'
    ),
    (
        'orders.add_note',
        'Add order notes',
        'Allows adding notes to orders',
        'orders'
    ),
    (
        'orders.delete',
        'Delete orders',
        'Allows deleting orders',
        'orders'
    ),
    (
        'orders.change_price',
        'Change order price',
        'Allows changing the final order price',
        'orders'
    ),

    -- Clients
    (
        'clients.view',
        'View clients',
        'Allows viewing client information',
        'clients'
    ),
    (
        'clients.create',
        'Create clients',
        'Allows creating new clients',
        'clients'
    ),
    (
        'clients.edit',
        'Edit clients',
        'Allows editing client information',
        'clients'
    ),
    (
        'clients.delete',
        'Delete clients',
        'Allows deleting clients',
        'clients'
    ),

    -- Cars
    (
        'cars.view',
        'View cars',
        'Allows viewing vehicle information',
        'cars'
    ),
    (
        'cars.create',
        'Create cars',
        'Allows adding vehicles',
        'cars'
    ),
    (
        'cars.edit',
        'Edit cars',
        'Allows editing vehicle information',
        'cars'
    ),
    (
        'cars.delete',
        'Delete cars',
        'Allows deleting vehicles',
        'cars'
    ),

    -- Employees
    (
        'employees.view',
        'View employees',
        'Allows viewing employee information',
        'employees'
    ),
    (
        'employees.create',
        'Create employees',
        'Allows creating employees',
        'employees'
    ),
    (
        'employees.edit',
        'Edit employees',
        'Allows editing employee information',
        'employees'
    ),
    (
        'employees.archive',
        'Archive employees',
        'Allows archiving and restoring employees',
        'employees'
    ),
    (
        'employees.manage_positions',
        'Manage employee positions',
        'Allows creating and editing employee positions',
        'employees'
    ),

    -- Users and roles
    (
        'users.view',
        'View users',
        'Allows viewing GarageHub user accounts',
        'users'
    ),
    (
        'users.create',
        'Create users',
        'Allows creating GarageHub user accounts',
        'users'
    ),
    (
        'users.edit',
        'Edit users',
        'Allows editing user accounts',
        'users'
    ),
    (
        'users.deactivate',
        'Deactivate users',
        'Allows deactivating and restoring user accounts',
        'users'
    ),
    (
        'roles.view',
        'View roles',
        'Allows viewing roles and their permissions',
        'roles'
    ),
    (
        'roles.create',
        'Create roles',
        'Allows creating custom roles',
        'roles'
    ),
    (
        'roles.edit',
        'Edit roles',
        'Allows editing custom roles and permissions',
        'roles'
    ),
    (
        'roles.delete',
        'Delete roles',
        'Allows deleting custom roles',
        'roles'
    ),

    -- Warehouse
    (
        'warehouse.view',
        'View warehouse',
        'Allows viewing warehouse items and stock levels',
        'warehouse'
    ),
    (
        'warehouse.manage',
        'Manage warehouse',
        'Allows adding, editing and adjusting warehouse items',
        'warehouse'
    ),

    -- Settings / modules
    (
        'settings.view',
        'View settings',
        'Allows viewing application settings',
        'settings'
    ),
    (
        'settings.edit',
        'Edit settings',
        'Allows changing application settings',
        'settings'
    ),
    (
        'settings.manage_modules',
        'Manage modules',
        'Allows enabling and disabling GarageHub modules',
        'settings'
    ),

    -- Future order works
    (
        'works.view',
        'View order works',
        'Allows viewing work items in orders',
        'works'
    ),
    (
        'works.create',
        'Create order works',
        'Allows adding work items to orders',
        'works'
    ),
    (
        'works.edit',
        'Edit order works',
        'Allows editing work items',
        'works'
    ),
    (
        'works.delete',
        'Delete order works',
        'Allows deleting work items',
        'works'
    ),

    -- Future parts/materials
    (
        'parts.view',
        'View order parts',
        'Allows viewing parts and materials used in orders',
        'parts'
    ),
    (
        'parts.create',
        'Add order parts',
        'Allows adding parts and materials to orders',
        'parts'
    ),
    (
        'parts.edit',
        'Edit order parts',
        'Allows editing parts and materials in orders',
        'parts'
    ),
    (
        'parts.delete',
        'Delete order parts',
        'Allows deleting parts and materials from orders',
        'parts'
    ),

    -- Reports
    (
        'reports.view',
        'View reports',
        'Allows viewing GarageHub reports',
        'reports'
    ),

    -- Payroll
    (
        'payroll.view',
        'View payroll',
        'Allows viewing payroll information',
        'payroll'
    ),
    (
        'payroll.manage',
        'Manage payroll',
        'Allows configuring and managing payroll',
        'payroll'
    ),

    -- Accounting
    (
        'accounting.view',
        'View accounting',
        'Allows viewing accounting information',
        'accounting'
    ),
    (
        'accounting.manage',
        'Manage accounting',
        'Allows managing accounting information',
        'accounting'
    )

ON CONFLICT (code) DO NOTHING;
