-- +goose Up
CREATE TABLE IF NOT EXISTS tasks (
	id TEXT NOT NULL UNIQUE PRIMARY KEY, -- Pattern: tk-#
    description TEXT NOT NULL UNIQUE,
    state TEXT NOT NULL, -- Valid options: not_started, started, done
    due_date TIMESTAMP NOT NULL,
    subject_id TEXT NOT NULL,
    tags TEXT[] NOT NULL,
	created_at TIMESTAMPTZ NOT NULL UNIQUE,
	updated_at TIMESTAMPTZ NOT NULL UNIQUE
);

ALTER TABLE tasks
ADD CONSTRAINT fk_tasks_subject
FOREIGN KEY (subject_id)
REFERENCES subjects(id);

ALTER TABLE tasks
ALTER COLUMN created_at SET DEFAULT CURRENT_TIMESTAMP,
ALTER COLUMN updated_at SET DEFAULT CURRENT_TIMESTAMP + INTERVAL '1 second';

INSERT INTO tasks(id, description, state, due_date, subject_id, tags) VALUES ('tk-0', 'Try', 'not_started', CURRENT_TIMESTAMP, 'sb-0', '{"Urgente"}');
-- +goose Down
ALTER TABLE tasks
DROP CONSTRAINT fk_tasks_subject;
DELETE FROM tasks;
DROP TABLE tasks;
