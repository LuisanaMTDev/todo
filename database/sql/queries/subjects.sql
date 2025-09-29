-- name: AddSubject :exec
INSERT INTO subjects (
	id, name
) VALUES ( ?, ? );

-- name: ModifySubjectName :execrows
UPDATE subjects SET
name = ?,
updated_at = datetime('now', 'localtime')
WHERE id = ?;

-- name: GetAllSubjects :many
SELECT * FROM subjects;

-- name: GetLastSubjectId :one
SELECT id FROM subjects
ORDER BY created_at DESC
LIMIT 1;

-- name: DeleteSubject :execrows
DELETE FROM subjects WHERE id = ?;
