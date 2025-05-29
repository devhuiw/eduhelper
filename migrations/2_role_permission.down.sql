DELETE rp
FROM
    role_permissions rp
    JOIN roles r ON rp.role_id = r.role_id
WHERE
    r.role_name IN ('admin', 'moderator', 'teacher');

DELETE FROM roles
WHERE
    role_name IN ('admin', 'moderator', 'teacher', 'user');

DELETE FROM permissions
WHERE
    permission_name IN (
        'permission:create',
        'permission:update',
        'permission:delete',
        'permission:view',
        'permission:list',
        'role:create',
        'role:update',
        'role:delete',
        'role:view',
        'role:list',
        'userrole:assign',
        'userrole:remove',
        'rolepermission:assign',
        'rolepermission:remove',
        'user:list',
        'user:view',
        'user:update',
        'user:delete',
        'teacher:create',
        'teacher:view',
        'teacher:view_self',
        'teacher:update',
        'teacher:update_self',
        'teacher:delete'
    );