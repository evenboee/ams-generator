-- name: CreateAssignment :one
INSERT INTO "assignments" ("course", "name", "reviews_per_submission", "time_due") 
    VALUES ($1, $2, $3, $4) RETURNING "id";
