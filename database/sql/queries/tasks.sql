-- name: AddTask :exec
INSERT INTO tasks (
	id, description, due_date, state, subject_id, academic_value
) VALUES ( ?, ?, datetime(?), ?, ?, ? );

-- name: ModifyTaskDescription :execrows
UPDATE tasks SET
description = ?,
updated_at = datetime('now', 'localtime')
WHERE id = ?;
-- name: ModifyTaskState :execrows
UPDATE tasks SET
state = ?,
updated_at = datetime('now', 'localtime')
WHERE id = ?;
-- name: ModifyTaskDueDate :execrows
UPDATE tasks SET
due_date = datetime(?),
updated_at = datetime('now', 'localtime')
WHERE id = ?;
-- name: ToggleTaskUrgentState :execrows
UPDATE tasks SET
is_urgent = CASE WHEN is_urgent = 1 THEN 0 ELSE 1 END,
updated_at = datetime('now', 'localtime')
WHERE id = ?;
-- name: ToggleTaskImportantState :execrows
UPDATE tasks SET
is_important = CASE WHEN is_important = 1 THEN 0 ELSE 1 END,
updated_at = datetime('now', 'localtime')
WHERE id = ?;
-- name: ModifyTaskAcademicValue :execrows
UPDATE tasks SET
academic_value = ?,
updated_at = datetime('now', 'localtime')
WHERE id = ?;

-- name: GetAllTasksBySubject :many
SELECT * FROM tasks
WHERE subject_id = ?;

-- name: GetLastTaskID :one
SELECT id FROM tasks
ORDER BY created_at DESC
LIMIT 1;

-- name: GetTaskByID :one
SELECT * FROM tasks
WHERE id = ?;

-- name: DeleteTask :exec
DELETE FROM tasks WHERE id = ?;
