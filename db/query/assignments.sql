-- name: CreateAssignment :one
INSERT INTO "assignments" ("course", "time_due") VALUES ($1, $2) RETURNING *;
