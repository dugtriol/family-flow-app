BEGIN;
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE IF NOT EXISTS "families"
(
    id         UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    name       VARCHAR(255) NOT NULL,
    created_at TIMESTAMP        DEFAULT CURRENT_TIMESTAMP
);


DROP TYPE IF EXISTS user_role CASCADE;
CREATE TYPE user_role AS ENUM (
    'Parent',
    'Child'
    );

CREATE TABLE IF NOT EXISTS "users"
(
    id        UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    name      VARCHAR(255)        NOT NULL,
    email     VARCHAR(255) UNIQUE NOT NULL,
    password  VARCHAR(255)        NOT NULL,
    role      user_role,
    family_id UUID REFERENCES families (id) ON DELETE SET NULL
);

DROP TYPE IF EXISTS task_status CASCADE;
CREATE TYPE task_status AS ENUM (
    'Active',
    'Completed',
    'Overdue'
    );

CREATE TABLE IF NOT EXISTS "tasks"
(
    id          UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    title       VARCHAR(255) NOT NULL,
    description TEXT,
    status      task_status,
    deadline    TIMESTAMP,
    assigned_to UUID REFERENCES users (id) ON DELETE CASCADE,
    created_by  UUID REFERENCES users (id) ON DELETE CASCADE,
    reward      INT
);

COMMIT;