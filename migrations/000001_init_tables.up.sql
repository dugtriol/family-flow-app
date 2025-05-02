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
    id        UUID PRIMARY KEY                                                DEFAULT uuid_generate_v4(),
    name      VARCHAR(255)        NOT NULL,
    email     VARCHAR(255) UNIQUE NOT NULL,
    password  VARCHAR(255)        NOT NULL,
    role      user_role,
    family_id UUID                REFERENCES families (id) ON DELETE SET NULL default null
);

DROP TYPE IF EXISTS item_visibility CASCADE;
CREATE TYPE item_visibility AS ENUM (
    'Private',
    'Public'
    );

DROP TYPE IF EXISTS shopping_item_status CASCADE;
CREATE TYPE shopping_item_status AS ENUM (
    'Active',
    'Reserved',
    'Completed'
    );

CREATE TABLE IF NOT EXISTS "shopping_items"
(
    id          UUID PRIMARY KEY                             DEFAULT uuid_generate_v4(),
    family_id   UUID REFERENCES families (id) ON DELETE CASCADE,
    title       VARCHAR(255) NOT NULL,
    description TEXT,
    status      shopping_item_status                         default 'Active',
    visibility  item_visibility,
    created_by  UUID REFERENCES users (id) ON DELETE CASCADE,
    reserved_by UUID REFERENCES users (id) ON DELETE CASCADE default null,
    buyer_id    UUID REFERENCES users (id) ON DELETE CASCADE default null,
    is_archived BOOLEAN                                      DEFAULT FALSE,
    created_at  TIMESTAMP                                    DEFAULT CURRENT_TIMESTAMP,
    updated_at  TIMESTAMP                                    DEFAULT CURRENT_TIMESTAMP
);

DROP TYPE IF EXISTS todos_item_status CASCADE;
CREATE TYPE todos_item_status AS ENUM (
    'Active',
    'Completed'
    );

CREATE TABLE IF NOT EXISTS "todo_items"
(
    id          UUID PRIMARY KEY  DEFAULT uuid_generate_v4(),
    family_id   UUID REFERENCES families (id) ON DELETE CASCADE,
    title       VARCHAR(255) NOT NULL,
    description TEXT,
    status      todos_item_status default 'Active',
    deadline    TIMESTAMP,
    assigned_to UUID REFERENCES users (id) ON DELETE CASCADE,
    created_by  UUID REFERENCES users (id) ON DELETE CASCADE,
    is_archived BOOLEAN           DEFAULT FALSE,
    created_at  TIMESTAMP         DEFAULT CURRENT_TIMESTAMP,
    updated_at  TIMESTAMP         DEFAULT CURRENT_TIMESTAMP
);

DROP TYPE IF EXISTS wishlist_item_status CASCADE;
CREATE TYPE wishlist_item_status AS ENUM (
    'Active',
    'Reserved',
    'Completed'
    );

CREATE TABLE IF NOT EXISTS "wishlist_items"
(
    id          UUID PRIMARY KEY                             DEFAULT uuid_generate_v4(),
    name        VARCHAR(255) NOT NULL,
    description TEXT,
    link        VARCHAR(255),
    status      wishlist_item_status                         default 'Active',
    created_by  UUID REFERENCES users (id) ON DELETE CASCADE,
    reserved_by UUID REFERENCES users (id) ON DELETE CASCADE default null,
    is_archived BOOLEAN                                      DEFAULT FALSE,
    created_at  TIMESTAMP                                    DEFAULT CURRENT_TIMESTAMP,
    updated_at  TIMESTAMP                                    DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS "diary_items"
(
    id          UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    title       VARCHAR(255) NOT NULL,
    description TEXT,
    emoji       VARCHAR(255),
    created_by  UUID REFERENCES users (id) ON DELETE CASCADE,
    created_at  TIMESTAMP        DEFAULT CURRENT_TIMESTAMP,
    updated_at  TIMESTAMP        DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS "locations"
(
    id         UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    latitude   FLOAT8,
    longitude  FLOAT8,
    created_by UUID REFERENCES users (id) ON DELETE CASCADE,
    created_at TIMESTAMP        DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP        DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS "notifications"
(
    id         UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_id    UUID REFERENCES users (id) ON DELETE CASCADE,
    title      TEXT,
    body       TEXT,
    data       TEXT,
    is_read    BOOLEAN          DEFAULT FALSE,
    created_at TIMESTAMP        DEFAULT CURRENT_TIMESTAMP
);

-- общие эвенты, несколько пользователей таблица пользователь - id эвента
-- диалог
-- Сообщения
-- документы

COMMIT;