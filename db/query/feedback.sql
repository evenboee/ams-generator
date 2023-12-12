-- name: CreateFeedback :one
INSERT INTO "feedback" ("id", "review", "answer", "rating", "feedback") 
    VALUES ($1, $2, $3, $4, $5) RETURNING *;
