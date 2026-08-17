CREATE TABLE role_permissions (
    role_id BIGINT NOT NULL,
    permission_id BIGINT NOT NULL,

    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    CONSTRAINT role_permissions_role_fk
        FOREIGN KEY (role_id)
        REFERENCES roles(id)
        ON DELETE CASCADE,

    CONSTRAINT role_permissions_permission_fk
        FOREIGN KEY (permission_id)
        REFERENCES permissions(id)
        ON DELETE CASCADE,

    CONSTRAINT role_permissions_unique
        UNIQUE (role_id, permission_id)
);

CREATE INDEX role_permissions_role_id_index
ON role_permissions (role_id);

CREATE INDEX role_permissions_permission_id_index
ON role_permissions (permission_id);
