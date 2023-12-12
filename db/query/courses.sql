-- name: CreateCourse :one
INSERT INTO "courses" ("id", "code", "year") VALUES ($1, $2, $3) RETURNING *;
