-- name: CreateTask :one
INSERT INTO tasks (
    title,
    text
) VALUES (
    $1,
    $2
) RETURNING id;

-- name: SetDone :exec
UPDATE tasks SET done = $1 WHERE id = $2;

-- name: GetTask :one
SELECT * FROM tasks WHERE id = $1;

-- name: GetTaskTitles :many
SELECT id, title, done FROM tasks LIMIT $1 OFFSET $2;

-- name: DeleteTask :exec
DELETE FROM tasks WHERE id = $1;

-- name: UpdateTaskTitle :exec
UPDATE tasks SET title = $1 WHERE id = $2;

-- name: UpdateTaskText :exec
UPDATE tasks SET text = $1 WHERE id = $2;

-- name: Search :many
SELECT id, title, done FROM tasks
                       WHERE to_tsvector(title) @@ plainto_tsquery($1)
   OR to_tsvector(text) @@ plainto_tsquery($1);