-- name: AddTask :exec
INSERT INTO tasks (
	id, description, due_date, state, tags, subject_id
) VALUES ( $1, $2, TO_TIMESTAMP($3, 'YYYY-MM-DD"T"HH24:MI'), $4, $5, $6 );

-- name: AddTag :exec
UPDATE tasks
SET tags = array_append(tags, $2)
WHERE id = $1;

-- name: ModifyTaskDescription :execrows
UPDATE tasks SET
description = $2,
updated_at = CURRENT_TIMESTAMP
WHERE id = $1;
-- name: ModifyTaskState :execrows
UPDATE tasks SET
state = $2,
updated_at = CURRENT_TIMESTAMP
WHERE id = $1;
-- name: ModifyTaskDueDate :execrows
UPDATE tasks SET
due_date = TO_TIMESTAMP($2, 'YYYY-MM-DD"T"HH24:MI'),
updated_at = CURRENT_TIMESTAMP
WHERE id = $1;

-- name: GetAllTasksBySubject :many
SELECT * FROM tasks
WHERE subject_id = $1;

-- name: GetLastTaskID :one
SELECT id FROM tasks
ORDER BY created_at DESC
LIMIT 1;

-- name: GetTaskByID :one
SELECT * FROM tasks
WHERE id = $1;

-- name: DeleteTask :exec
DELETE FROM tasks WHERE id = $1;

-- name: DeleteTag :exec
UPDATE tasks
SET tags = array_remove(tags, $2)
WHERE id = $1;
