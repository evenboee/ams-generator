-- name: CreateAnswer :one
INSERT INTO "answers" ("question", "submission", "answer") VALUES ($1, $2, $3) RETURNING "id";
