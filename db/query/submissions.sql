-- name: CreateSubmission :one
INSERT INTO "submissions" ("id", "user", "assignment") VALUES ($1, $2, $3) RETURNING *;
