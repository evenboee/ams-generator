-- name: CreateAssignment :one
INSERT INTO "assignments" ("course", "name", "reviews_per_submission", "time_due") 
    VALUES ($1, $2, $3, $4) RETURNING "id";


-- name: GetAssignment :one
SELECT *, get_assignment_rating(sqlc.arg(id)) AS rating
    FROM assignments
    WHERE id = sqlc.arg(id);
