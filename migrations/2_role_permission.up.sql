-- Права
INSERT INTO
    permissions (permission_name)
VALUES
    -- Общие права на permissions и роли
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
    -- Права на пользователей
    ('user:create'),
    ('user:view'),
    ('user:update'),
    ('user:delete'),
    ('user:list'),
    -- Права на учителей
    ('teacher:create'),
    ('teacher:view'),
    ('teacher:view_self'),
    ('teacher:update'),
    ('teacher:update_self'),
    ('teacher:delete'),
    ('teacher:list'),
    -- Права на студентов
    ('student:create'),
    ('student:view'),
    ('student:view_public'),
    ('student:update'),
    ('student:delete'),
    ('student:list'),
    ('student:list_public'),
    -- Права на группы
    ('studentgroup:create'),
    ('studentgroup:view'),
    ('studentgroup:view_public'),
    ('studentgroup:update'),
    ('studentgroup:delete'),
    ('studentgroup:list'),
    ('studentgroup:list_public'),
    -- Права на дисциплины
    ('discipline:create'),
    ('discipline:view'),
    ('discipline:view_public'),
    ('discipline:update'),
    ('discipline:delete'),
    ('discipline:list'),
    ('discipline:list_public'),
    -- Права на посещаемость
    ('attendance:create'),
    ('attendance:view'),
    ('attendance:update'),
    ('attendance:delete'),
    ('attendance:list'),
    -- Права на журнал оценок
    ('gradejournal:create'),
    ('gradejournal:view'),
    ('gradejournal:list'),
    ('gradejournal:list_public'),
    ('gradejournal:avg'),
    ('gradejournal:update'),
    ('gradejournal:delete'),
    -- Права на семестры
    ('semester:create'),
    ('semester:view'),
    ('semester:update'),
    ('semester:delete'),
    ('semester:list'),
    -- Права на учебные года
    ('academicyear:create'),
    ('academicyear:view'),
    ('academicyear:update'),
    ('academicyear:delete'),
    ('academicyear:list'),
    -- Права на учебные планы
    ('curriculum:create'),
    ('curriculum:view'),
    ('curriculum:update'),
    ('curriculum:delete'),
    ('curriculum:list');

INSERT INTO
    roles (role_name)
VALUES
    ('admin'),
    ('admin-teacher'),
    ('teacher'),
    ('student');

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
    r.role_name = 'admin-teacher'
    AND p.permission_name NOT IN (
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
        'userrole:view',
        'rolepermission:assign',
        'rolepermission:remove',
        'rolepermission:view',
        'user:create',
        'user:update',
        'user:delete'
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
    AND p.permission_name IN (
        'teacher:view_self',
        'teacher:update_self',
        'student:view',
        'student:list',
        'student:view_public',
        'student:list_public',
        'studentgroup:view',
        'studentgroup:list',
        'studentgroup:view_public',
        'studentgroup:list_public',
        'discipline:view',
        'discipline:list',
        'discipline:view_public',
        'discipline:list_public',
        'gradejournal:create',
        'gradejournal:view',
        'gradejournal:list',
        'gradejournal:list_public',
        'gradejournal:avg',
        'attendance:create',
        'attendance:view',
        'attendance:list',
        'curriculum:view',
        'curriculum:list',
        'semester:view',
        'semester:list',
        'academicyear:view',
        'academicyear:list'
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
    r.role_name = 'student'
    AND p.permission_name IN (
        'student:view',
        'student:view_public',
        'student:list_public',
        'studentgroup:view_public',
        'studentgroup:list_public',
        'teacher:view_public',
        'discipline:view_public',
        'discipline:list_public',
        'gradejournal:view',
        'gradejournal:list_public',
        'gradejournal:avg',
        'attendance:view',
        'attendance:list',
        'curriculum:view',
        'curriculum:list',
        'semester:view',
        'semester:list',
        'academicyear:view',
        'academicyear:list'
    );