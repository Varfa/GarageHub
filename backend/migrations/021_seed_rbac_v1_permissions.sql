INSERT INTO permissions (
    code,
    name,
    description,
    module
)
VALUES
    (
        'warehouse.writeoff',
        'Write off warehouse stock',
        'Allows writing off warehouse stock and recording stock losses',
        'warehouse'
    ),
    (
        'reports.financial',
        'View financial reports',
        'Allows viewing financial reports and financial summaries',
        'reports'
    ),
    (
        'permissions.manage',
        'Manage permissions',
        'Allows managing role permission assignments',
        'permissions'
    ),
    (
        'parts.manage',
        'Manage parts catalog',
        'Allows managing parts catalog, stock-related part data and part configuration',
        'parts'
    )

ON CONFLICT (code) DO NOTHING;
