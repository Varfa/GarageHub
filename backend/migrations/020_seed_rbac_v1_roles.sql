INSERT INTO roles (
    code,
    name,
    description,
    is_system,
    is_active
)
VALUES
    (
        'service_admin',
        'Service Administrator',
        'Operational administrator role for GarageHub service staff',
        TRUE,
        TRUE
    ),
    (
        'mechanic',
        'Mechanic',
        'Default role for mechanics working with clients, cars and service orders',
        TRUE,
        TRUE
    ),
    (
        'auto_electrician',
        'Auto Electrician',
        'Default role for automotive electricians and diagnostic technicians',
        TRUE,
        TRUE
    ),
    (
        'tire_technician',
        'Tire Technician',
        'Default role for tire service employees',
        TRUE,
        TRUE
    ),
    (
        'warehouse_manager',
        'Warehouse Manager',
        'Default role for employees responsible for warehouse and parts management',
        TRUE,
        TRUE
    ),
    (
        'accountant',
        'Accountant',
        'Default role for accounting, payroll and financial reporting',
        TRUE,
        TRUE
    )

ON CONFLICT (code) DO NOTHING;
