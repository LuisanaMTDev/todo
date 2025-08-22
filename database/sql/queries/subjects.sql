-- name: AddSubject :exec
INSERT INTO subjects (
	id, name
) VALUES ( $1, $2 );

-- name: ModifySubjectName :execrows
UPDATE subjects SET
name = $2,
updated_at = CURRENT_TIMESTAMP
WHERE id = $1;

-- name: GetAllSubjects :many
SELECT * FROM subjects;

-- name: GetLastSubjectId :one
SELECT id FROM subjects
ORDER BY created_at DESC
LIMIT 1;

-- name: DeleteSubject :execrows
DELETE FROM subjects WHERE id = $1;
