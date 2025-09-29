-- +goose Up
CREATE TABLE IF NOT EXISTS subjects (
	id TEXT NOT NULL UNIQUE PRIMARY KEY, -- Pattern: sb-#
  name TEXT NOT NULL UNIQUE,
	created_at TEXT NOT NULL DEFAULT (datetime('now', 'localtime')),
	updated_at TEXT NOT NULL DEFAULT (datetime('now', 'localtime'))
);

INSERT INTO subjects(id, name) VALUES ('sb-0', 'Negocio');
-- +goose Down
DELETE FROM subjects;
DROP TABLE subjects;
