-- Модуль planner: домашка и задачи в одной таблице.

CREATE TABLE items (
    id          BIGINT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    user_id     BIGINT      NOT NULL REFERENCES users (id) ON DELETE CASCADE,
    kind        TEXT        NOT NULL,
    subject_id  BIGINT REFERENCES subjects (id) ON DELETE SET NULL,
    text        TEXT        NOT NULL,
    due_date    DATE        NOT NULL, -- день, в который запись видна в списках
    deadline_at TIMESTAMPTZ,          -- точный срок, у ДЗ необязателен
    remind_at   TIMESTAMPTZ,
    reminded_at TIMESTAMPTZ,
    done_at     TIMESTAMPTZ,
    created_at  TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at  TIMESTAMPTZ NOT NULL DEFAULT now(),

    CONSTRAINT items_kind_valid CHECK (kind IN ('homework', 'task')),
    CONSTRAINT items_text_not_blank CHECK (btrim(text) <> ''),
    CONSTRAINT items_subject_only_for_homework CHECK (kind = 'homework' OR subject_id IS NULL)
);

-- Списки на день и неделю.
CREATE INDEX items_user_due_date_idx ON items (user_id, due_date);

-- Фоновая задача напоминаний: только то, что ещё предстоит отправить.
CREATE INDEX items_pending_reminders_idx ON items (remind_at)
    WHERE remind_at IS NOT NULL AND reminded_at IS NULL AND done_at IS NULL;

CREATE TABLE planner_settings (
    user_id                BIGINT PRIMARY KEY REFERENCES users (id) ON DELETE CASCADE,
    homework_remind_before INTERVAL NOT NULL DEFAULT '1 day',
    task_remind_before     INTERVAL NOT NULL DEFAULT '1 hour',

    CONSTRAINT planner_settings_homework_remind_not_negative CHECK (homework_remind_before >= INTERVAL '0'),
    CONSTRAINT planner_settings_task_remind_not_negative CHECK (task_remind_before >= INTERVAL '0')
);
