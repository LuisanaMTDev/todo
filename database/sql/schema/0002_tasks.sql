-- +goose Up
CREATE TABLE IF NOT EXISTS tasks (
	id TEXT NOT NULL UNIQUE PRIMARY KEY, -- Pattern: tk-#
    description TEXT NOT NULL UNIQUE,
    state TEXT NOT NULL, -- Valid options: not_started, started, done
    due_date TIMESTAMP NOT NULL,
    subject_id TEXT NOT NULL,
    is_urgent INTEGER NOT NULL DEFAULT 0 CHECK (is_urgent IN (0,1)),
    is_important INTEGER NOT NULL DEFAULT 0 CHECK (is_important IN (0,1)),
    academic_value NUMERIC,
	created_at TEXT NOT NULL DEFAULT (datetime('now', 'localtime')),
	updated_at TEXT NOT NULL DEFAULT (datetime('now', 'localtime')),
  CONSTRAINT fk_task_subject
    FOREIGN KEY (subject_id)
    REFERENCES subjects(id)
    ON DELETE CASCADE
);

INSERT INTO tasks(id, description, state, due_date, subject_id, academic_value) VALUES ('tk-0', 'Try', 'not_started', datetime('now', 'localtime'), 'sb-0', 5.5);
-- +goose Down
DELETE FROM tasks;
DROP TABLE tasks;
