-- name: CreateSubmission :one
INSERT INTO "submissions" ("user", "assignment") VALUES ($1, $2) RETURNING "id";
