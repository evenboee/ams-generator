-- name: CreateQuestion :one
INSERT INTO "questions" ("assignment", "prompt", "order") VALUES ($1, $2, $3) RETURNING "id";
