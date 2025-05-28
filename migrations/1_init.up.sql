create table
    "user" (
        user_id bigint generated always as identity primary key,
        created_at timestamp not null default current_timestamp,
        update_at timestamp not null default current_timestamp,
        first_name varchar(100) not null check (length (first_name) >= 2),
        last_name varchar(100) not null check (length (last_name) >= 2),
        middle_name varchar(100) check (length (middle_name) >= 2),
        email varchar(350) not null check (length (email) >= 5) unique,
        password varchar(64) not null
    );

create table
    "role" (
        role_id bigint generated always as identity primary key,
        created_at timestamp not null default current_timestamp,
        update_at timestamp not null default current_timestamp,
        role_name varchar(150) not null check (length (role_name) >= 3)
    );

create table
    "permission" (
        permission_id bigint generated always as identity primary key,
        created_at timestamp not null default current_timestamp,
        update_at timestamp not null default current_timestamp,
        permission_name varchar(150) not null check (length (permission_name) >= 6)
    );

create table
    "role_permission" (
        created_at timestamp not null default current_timestamp,
        update_at timestamp not null default current_timestamp,
        role_id bigint references role (role_id),
        permission_id bigint references permission (permission_id),
        primary key (role_id, permission_id)
    );

create table
    "user_role" (
        created_at timestamp not null default current_timestamp,
        update_at timestamp not null default current_timestamp,
        role_id bigint references role (role_id),
        user_id bigint references users (user_id),
        primary key (role_id, user_id)
    );

create table
    "teacher" (
        user_id bigint references user (user_id) primary key,
        created_at timestamp not null default current_timestamp,
        update_at timestamp not null default current_timestamp,
        phone varchar(100) not null check (length (phone) >= 2),
        working_experience text,
        education text
    );

create table
    "student" (
        user_id bigint references user (user_id) primary key,
        phone varchar(100) not null check (length (phone) >= 2),
        birtday date not null check (birtday >= to_date ('1920-01-01', 'YYYY-MM-DD')),
        created_at timestamp not null default current_timestamp,
        update_at timestamp not null default current_timestamp,
        student_group_id bigint references student_group (student_group_id) not null
    );

create table
    "academic_year" (
        academic_year_id bigint generated always as identity primary key,
        created_at timestamp not null default current_timestamp,
        update_at timestamp not null default current_timestamp,
        start_with date not null check (
            start_with <= to_date ('2024-01-01', 'YYYY-MM-DD')
        ),
        ends_with date not null check (ends_with >= to_date ('2024-01-01', 'YYYY-MM-DD'))
    );

create table
    "semester" (
        semester_id bigint generated always as identity primary key,
        created_at timestamp not null default current_timestamp,
        update_at timestamp not null default current_timestamp,
        start_with date not null check (
            start_with <= to_date ('2024-01-01', 'YYYY-MM-DD')
        ),
        ends_with date not null check (ends_with >= to_date ('2024-01-01', 'YYYY-MM-DD')),
        academic_year_id bigint not null references academic_year (academic_year_id)
    );

create table
    "student_group" (
        student_group_id bigint generated always as identity primary key,
        created_at timestamp not null default current_timestamp,
        update_at timestamp not null default current_timestamp,
        student_group_name varchar(150) not null check (length (student_group_name) >= 3),
        curator_id bigint not null references user (user_id),
        academic_year_id bigint not null references academic_year (academic_year_id)
    );

create table
    "discipline" (
        discipline_id bigint generated always as identity primary key,
        created_at timestamp not null default current_timestamp,
        update_at timestamp not null default current_timestamp,
        discipline_name varchar(155) not null check (length (discipline_name) >= 3),
        teacher_id bigint not null references user (user_id),
        student_group_id bigint not null references student_group (student_group_id)
    );

create table
    "grade_journal" (
        grade_journal_id bigint generated always as identity primary key,
        created_at timestamp not null default current_timestamp,
        update_at timestamp not null default current_timestamp,
        student_id bigint not null references student (user_id),
        grade smallint not null check (grade between 1 and 10),
        comment text,
        discipline_id bigint not null references discipline (discipline_id)
    );

create table
    "attendance" (
        attendance_id bigint generated always as identity primary key,
        created_at timestamp not null default current_timestamp,
        visit boolean not null default 'true',
        comment text,
        update_at timestamp not null default current_timestamp,
        student_id bigint not null references student (user_id),
        discipline_id bigint not null references discipline (discipline_id)
    );

create table
    "curriculum" (
        curriculum_id bigint generated always as identity primary key,
        created_at timestamp not null default current_timestamp,
        update_at timestamp not null default current_timestamp,
        subject_name varchar(150) not null,
        subject_description text,
        semester_id bigint references semester (semester_id),
        discipline_id bigint not null references discipline (discipline_id)
    );