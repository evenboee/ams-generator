-- name: CrateQuestion :one
INSERT INTO "questions" ("id", "assignment", "prompt", "order") VALUES ($1, $2, $3, $4) RETURNING *;
