-- Общие таблицы: владелец всех данных и состояние диалогов.

CREATE TABLE users (
    id         BIGINT PRIMARY KEY, -- Telegram user id
    chat_id    BIGINT      NOT NULL,
    username   TEXT        NOT NULL DEFAULT '',
    first_name TEXT        NOT NULL DEFAULT '',
    timezone   TEXT        NOT NULL DEFAULT 'Europe/Moscow',
    created_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE TABLE bot_sessions (
    user_id    BIGINT PRIMARY KEY REFERENCES users (id) ON DELETE CASCADE,
    state      TEXT        NOT NULL,
    data       JSONB       NOT NULL DEFAULT '{}',
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now()
);
