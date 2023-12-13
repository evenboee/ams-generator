-- name: CreateFeedback :one
INSERT INTO "feedback" ("review", "answer", "rating", "feedback") 
    VALUES ($1, $2, $3, $4) RETURNING "id";
