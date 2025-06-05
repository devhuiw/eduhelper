CREATE TABLE
    `user` (
        user_id BIGINT AUTO_INCREMENT PRIMARY KEY,
        created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
        updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
        first_name VARCHAR(100) NOT NULL,
        last_name VARCHAR(100) NOT NULL,
        middle_name VARCHAR(100),
        email VARCHAR(350) NOT NULL UNIQUE,
        password VARCHAR(64) NOT NULL,
        CHECK (CHAR_LENGTH(first_name) >= 2),
        CHECK (CHAR_LENGTH(last_name) >= 2),
        CHECK (
            middle_name IS NULL
            OR CHAR_LENGTH(middle_name) >= 2
        ),
        CHECK (CHAR_LENGTH(email) >= 5)
    );

CREATE TABLE
    `roles` (
        role_id BIGINT AUTO_INCREMENT PRIMARY KEY,
        created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
        updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
        role_name VARCHAR(150) NOT NULL,
        CHECK (CHAR_LENGTH(role_name) >= 3)
    );

CREATE TABLE
    `permissions` (
        permission_id BIGINT AUTO_INCREMENT PRIMARY KEY,
        created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
        updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
        permission_name VARCHAR(150) NOT NULL,
        CHECK (CHAR_LENGTH(permission_name) >= 6)
    );

CREATE TABLE
    `role_permissions` (
        created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
        updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
        role_id BIGINT NOT NULL,
        permission_id BIGINT NOT NULL,
        PRIMARY KEY (role_id, permission_id),
        FOREIGN KEY (role_id) REFERENCES roles (role_id),
        FOREIGN KEY (permission_id) REFERENCES permissions (permission_id)
    );

CREATE TABLE
    `user_roles` (
        created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
        updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
        role_id BIGINT NOT NULL,
        user_id BIGINT NOT NULL,
        PRIMARY KEY (role_id, user_id),
        FOREIGN KEY (role_id) REFERENCES roles (role_id),
        FOREIGN KEY (user_id) REFERENCES user (user_id)
    );

CREATE TABLE
    `teacher` (
        user_id BIGINT PRIMARY KEY,
        created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
        updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
        phone VARCHAR(100) NOT NULL,
        working_experience TEXT,
        education TEXT,
        FOREIGN KEY (user_id) REFERENCES user (user_id),
        CHECK (CHAR_LENGTH(phone) >= 2)
    );

CREATE TABLE
    `academic_year` (
        academic_year_id BIGINT AUTO_INCREMENT PRIMARY KEY,
        name_academic_year VARCHAR(155) NOT NULL,
        created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
        updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
        start_with DATE NOT NULL,
        ends_with DATE NOT NULL,
        CHECK (start_with <= '2024-01-01'),
        CHECK (ends_with >= '2024-01-01')
    );

CREATE TABLE
    `student_group` (
        student_group_id BIGINT AUTO_INCREMENT PRIMARY KEY,
        created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
        updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
        student_group_name VARCHAR(150) NOT NULL,
        curator_id BIGINT NOT NULL,
        academic_year_id BIGINT NOT NULL,
        FOREIGN KEY (curator_id) REFERENCES user (user_id),
        FOREIGN KEY (academic_year_id) REFERENCES academic_year (academic_year_id),
        CHECK (CHAR_LENGTH(student_group_name) >= 3)
    );

CREATE TABLE
    `student` (
        user_id BIGINT PRIMARY KEY,
        phone VARCHAR(100) NOT NULL,
        birtday DATE NOT NULL,
        created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
        updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
        student_group_id BIGINT NOT NULL,
        FOREIGN KEY (user_id) REFERENCES user (user_id),
        FOREIGN KEY (student_group_id) REFERENCES student_group (student_group_id),
        CHECK (CHAR_LENGTH(phone) >= 2),
        CHECK (birtday >= '1920-01-01')
    );

CREATE TABLE
    `semester` (
        semester_id BIGINT AUTO_INCREMENT PRIMARY KEY,
        created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
        updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
        start_with DATE NOT NULL,
        ends_with DATE NOT NULL,
        academic_year_id BIGINT NOT NULL,
        FOREIGN KEY (academic_year_id) REFERENCES academic_year (academic_year_id),
        CHECK (start_with <= '2024-01-01'),
        CHECK (ends_with >= '2024-01-01')
    );

CREATE TABLE
    `discipline` (
        discipline_id BIGINT AUTO_INCREMENT PRIMARY KEY,
        created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
        updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
        discipline_name VARCHAR(155) NOT NULL,
        teacher_id BIGINT NOT NULL,
        student_group_id BIGINT NOT NULL,
        FOREIGN KEY (teacher_id) REFERENCES teacher (user_id),
        FOREIGN KEY (student_group_id) REFERENCES student_group (student_group_id),
        CHECK (CHAR_LENGTH(discipline_name) >= 3)
    );

CREATE TABLE
    `grade_journal` (
        grade_journal_id BIGINT AUTO_INCREMENT PRIMARY KEY,
        created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
        updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
        student_id BIGINT NOT NULL,
        grade SMALLINT NOT NULL,
        comment TEXT,
        discipline_id BIGINT NOT NULL,
        FOREIGN KEY (student_id) REFERENCES student (user_id),
        FOREIGN KEY (discipline_id) REFERENCES discipline (discipline_id),
        CHECK (grade BETWEEN 1 AND 10)
    );

CREATE TABLE
    `attendance` (
        attendance_id BIGINT AUTO_INCREMENT PRIMARY KEY,
        created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
        visit BOOLEAN NOT NULL DEFAULT TRUE,
        comment TEXT,
        updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
        student_id BIGINT NOT NULL,
        discipline_id BIGINT NOT NULL,
        FOREIGN KEY (student_id) REFERENCES student (user_id),
        FOREIGN KEY (discipline_id) REFERENCES discipline (discipline_id)
    );

CREATE TABLE
    `curriculum` (
        curriculum_id BIGINT AUTO_INCREMENT PRIMARY KEY,
        created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
        updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
        subject_name VARCHAR(150) NOT NULL,
        subject_description TEXT,
        semester_id BIGINT,
        discipline_id BIGINT NOT NULL,
        FOREIGN KEY (semester_id) REFERENCES semester (semester_id),
        FOREIGN KEY (discipline_id) REFERENCES discipline (discipline_id)
    );

CREATE TABLE
    `audit_log` (
        audit_id BIGINT AUTO_INCREMENT PRIMARY KEY,
        created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
        user_id BIGINT,
        table_name VARCHAR(100) NOT NULL,
        row_id BIGINT NOT NULL,
        action_type ENUM ('INSERT', 'UPDATE', 'DELETE') NOT NULL,
        old_data JSON,
        new_data JSON,
        comment TEXT,
        FOREIGN KEY (user_id) REFERENCES user (user_id)
    );