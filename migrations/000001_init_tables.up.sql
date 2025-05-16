BEGIN;

CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE IF NOT EXISTS "families" (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    name VARCHAR(255) NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    photo VARCHAR(255) DEFAULT NULL
);

DROP TYPE IF EXISTS user_role CASCADE;

CREATE TYPE user_role AS ENUM ('Parent', 'Child', 'Unknown');

CREATE TYPE user_gender AS ENUM ('Male', 'Female', 'Unknown');

CREATE TABLE IF NOT EXISTS "users" (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    name VARCHAR(255) NOT NULL,
    email VARCHAR(255) UNIQUE NOT NULL,
    password VARCHAR(255) NOT NULL,
    role user_role default 'Unknown',
    family_id UUID REFERENCES families (id) ON DELETE
    SET
        NULL default null,
        latitude FLOAT8 default null,
        longitude FLOAT8 default null,
        gender user_gender default 'Unknown',
        point INT default 0,
        birth_date DATE DEFAULT NULL,
        avatar VARCHAR(255) DEFAULT NULL
);

DROP TYPE IF EXISTS item_visibility CASCADE;

CREATE TYPE item_visibility AS ENUM ('Private', 'Public');

DROP TYPE IF EXISTS shopping_item_status CASCADE;

CREATE TYPE shopping_item_status AS ENUM ('Active', 'Reserved', 'Completed');

CREATE TABLE IF NOT EXISTS "shopping_items" (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    family_id UUID REFERENCES families (id) ON DELETE CASCADE,
    title VARCHAR(255) NOT NULL,
    description TEXT,
    status shopping_item_status default 'Active',
    visibility item_visibility,
    created_by UUID REFERENCES users (id) ON DELETE CASCADE,
    reserved_by UUID REFERENCES users (id) ON DELETE CASCADE default null,
    buyer_id UUID REFERENCES users (id) ON DELETE CASCADE default null,
    is_archived BOOLEAN DEFAULT FALSE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

DROP TYPE IF EXISTS todos_item_status CASCADE;

CREATE TYPE todos_item_status AS ENUM ('Active', 'Completed');

CREATE TABLE IF NOT EXISTS "todo_items" (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    family_id UUID REFERENCES families (id) ON DELETE CASCADE,
    title VARCHAR(255) NOT NULL,
    description TEXT,
    status todos_item_status default 'Active',
    deadline TIMESTAMP,
    assigned_to UUID REFERENCES users (id) ON DELETE CASCADE,
    created_by UUID REFERENCES users (id) ON DELETE CASCADE,
    is_archived BOOLEAN DEFAULT FALSE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    point INT default 0
);

DROP TYPE IF EXISTS wishlist_item_status CASCADE;

CREATE TYPE wishlist_item_status AS ENUM ('Active', 'Reserved', 'Completed');

CREATE TABLE IF NOT EXISTS "wishlist_items" (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    name VARCHAR(255) NOT NULL,
    description TEXT,
    link VARCHAR(255),
    status wishlist_item_status default 'Active',
    created_by UUID REFERENCES users (id) ON DELETE CASCADE,
    reserved_by UUID REFERENCES users (id) ON DELETE CASCADE default null,
    is_archived BOOLEAN DEFAULT FALSE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    photo VARCHAR(255) DEFAULT NULL
);

CREATE TABLE IF NOT EXISTS "diary_items" (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    title VARCHAR(255) NOT NULL,
    description TEXT,
    emoji VARCHAR(255),
    created_by UUID REFERENCES users (id) ON DELETE CASCADE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS "notifications" (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_id UUID REFERENCES users (id) ON DELETE CASCADE,
    title TEXT,
    body TEXT,
    data TEXT,
    is_read BOOLEAN DEFAULT FALSE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Таблица для хранения чатов
CREATE TABLE IF NOT EXISTS chats (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    name VARCHAR(255) NOT NULL,
    -- Название чата (например, "Семейный чат")
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Таблица для хранения сообщений
CREATE TABLE IF NOT EXISTS messages (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    chat_id UUID REFERENCES chats (id) ON DELETE CASCADE,
    -- ID чата, к которому относится сообщение
    sender_id UUID REFERENCES users (id) ON DELETE CASCADE,
    -- ID получателя
    content TEXT NOT NULL,
    -- Текст сообщения
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP -- Время отправки сообщения
);

-- Таблица для хранения участников чатов
CREATE TABLE IF NOT EXISTS chat_participants (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    chat_id UUID REFERENCES chats (id) ON DELETE CASCADE,
    -- ID чата
    user_id UUID REFERENCES users (id) ON DELETE CASCADE,
    -- ID пользователя
    joined_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP -- Время добавления в чат
);

CREATE TABLE IF NOT EXISTS rewards (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    family_id UUID REFERENCES families (id) ON DELETE CASCADE,
    title VARCHAR(255) NOT NULL,
    -- Название вознаграждения
    description TEXT,
    -- Описание вознаграждения
    cost INT NOT NULL,
    -- Стоимость вознаграждения в очках
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS reward_redemptions (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_id UUID REFERENCES users (id) ON DELETE CASCADE,
    reward_id UUID REFERENCES rewards (id) ON DELETE CASCADE,
    redeemed_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE fcm_tokens (
    id SERIAL PRIMARY KEY,
    user_id UUID NOT NULL UNIQUE,
    token TEXT NOT NULL,
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW()
);

COMMIT;