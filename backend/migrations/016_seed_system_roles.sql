INSERT INTO roles (
    code,
    name,
    description,
    is_system,
    is_active
)
VALUES (
    'service_admin',
    'Service Administrator',
    'Default administrator role for GarageHub service staff',
    TRUE,
    TRUE
)
ON CONFLICT (code) DO NOTHING;
