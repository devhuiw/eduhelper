INSERT INTO
    permissions (permission_name)
VALUES
    ('permission:create'),
    ('permission:update'),
    ('permission:delete'),
    ('permission:view'),
    ('permission:list'),
    ('role:create'),
    ('role:update'),
    ('role:delete'),
    ('role:view'),
    ('role:list'),
    ('userrole:assign'),
    ('userrole:remove'),
    ('userrole:view'),
    ('rolepermission:assign'),
    ('rolepermission:remove'),
    ('rolepermission:view'),
    ('user:list'),
    ('user:view'),
    ('user:update'),
    ('user:delete'),
    ('teacher:create'),
    ('teacher:view'),
    ('teacher:view_self'),
    ('teacher:update'),
    ('teacher:update_self'),
    ('teacher:delete');

INSERT INTO
    roles (role_name)
VALUES
    ('admin'),
    ('moderator'),
    ('teacher'),
    ('user');

INSERT INTO
    role_permissions (role_id, permission_id)
SELECT
    r.role_id,
    p.permission_id
FROM
    roles r,
    permissions p
WHERE
    r.role_name = 'admin';

INSERT INTO
    role_permissions (role_id, permission_id)
SELECT
    r.role_id,
    p.permission_id
FROM
    roles r,
    permissions p
WHERE
    r.role_name = 'moderator'
    AND p.permission_name IN (
        'user:list',
        'user:view',
        'user:update',
        'user:delete',
        'teacher:view',
        'teacher:list'
    );

INSERT INTO
    role_permissions (role_id, permission_id)
SELECT
    r.role_id,
    p.permission_id
FROM
    roles r,
    permissions p
WHERE
    r.role_name = 'teacher'
    AND p.permission_name IN ('teacher:view_self', 'teacher:update_self');