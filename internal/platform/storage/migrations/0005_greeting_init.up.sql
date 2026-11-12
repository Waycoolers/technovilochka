-- Модуль greeting: настроение и время приветствия по дням недели.

CREATE TABLE greeting_settings (
    user_id BIGINT PRIMARY KEY REFERENCES users (id) ON DELETE CASCADE,
    enabled BOOLEAN NOT NULL DEFAULT TRUE,
    mood    TEXT    NOT NULL DEFAULT 'cute'
);

CREATE TABLE greeting_times (
    user_id      BIGINT   NOT NULL REFERENCES users (id) ON DELETE CASCADE,
    weekday      SMALLINT NOT NULL, -- 1 = понедельник, 7 = воскресенье
    send_at      TIME     NOT NULL, -- местное время пользователя
    last_sent_on DATE,              -- защита от повторной отправки после рестарта

    PRIMARY KEY (user_id, weekday),
    CONSTRAINT greeting_times_weekday_range CHECK (weekday BETWEEN 1 AND 7)
);
