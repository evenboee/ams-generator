-- name: CreateAnswer :one
INSERT INTO "answers" ("id", "question", "submission", "answer") VALUES ($1, $2, $3, $4) RETURNING *;
