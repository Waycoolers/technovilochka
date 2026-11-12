-- Модуль schedule: предметы, шаблон пар на неделю и настройки расписания.

CREATE TABLE schedule_settings (
    user_id          BIGINT PRIMARY KEY REFERENCES users (id) ON DELETE CASCADE,
    odd_week_anchor  DATE, -- любая дата из нечётной недели
    group_code       TEXT, -- группа для импорта с сайта
    last_imported_at TIMESTAMPTZ
);

CREATE TABLE subjects (
    id      BIGINT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    user_id BIGINT NOT NULL REFERENCES users (id) ON DELETE CASCADE,
    name    TEXT   NOT NULL,

    CONSTRAINT subjects_name_not_blank CHECK (btrim(name) <> ''),
    CONSTRAINT subjects_user_name_key UNIQUE (user_id, name)
);

CREATE TABLE lessons (
    id               BIGINT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    user_id          BIGINT   NOT NULL REFERENCES users (id) ON DELETE CASCADE,
    subject_id       BIGINT   NOT NULL REFERENCES subjects (id) ON DELETE CASCADE,
    weekday          SMALLINT NOT NULL, -- 1 = понедельник, 7 = воскресенье
    week_parity      TEXT     NOT NULL,
    starts_at        TIME     NOT NULL, -- местное время пользователя
    ends_at          TIME     NOT NULL,
    room             TEXT,              -- значения с сайта, импорт их обновляет
    teacher          TEXT,
    room_override    TEXT,              -- ручные правки, импорт их не трогает
    teacher_override TEXT,
    source           TEXT     NOT NULL DEFAULT 'manual',
    is_hidden        BOOLEAN  NOT NULL DEFAULT FALSE, -- убрана вручную, импорт не должен вернуть

    CONSTRAINT lessons_weekday_range CHECK (weekday BETWEEN 1 AND 7),
    CONSTRAINT lessons_week_parity_valid CHECK (week_parity IN ('odd', 'even', 'both')),
    CONSTRAINT lessons_source_valid CHECK (source IN ('import', 'manual')),
    CONSTRAINT lessons_time_order CHECK (ends_at > starts_at),
    CONSTRAINT lessons_natural_key UNIQUE (user_id, weekday, week_parity, starts_at, subject_id)
);
