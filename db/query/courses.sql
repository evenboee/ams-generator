-- name: CreateCourse :one
INSERT INTO "courses" ("code", "year") VALUES ($1, $2) RETURNING "id";
