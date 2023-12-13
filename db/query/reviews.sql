-- name: CreateReview :one
INSERT INTO "reviews" ("submission", "reviewer_id", "finished_at", "created_at") 
    VALUES ($1, $2, $3, $4) RETURNING "id";
