-- name: GetUser :one
SELECT * FROM "users" WHERE "id" = $1;

-- name: CreateUser :one
INSERT INTO "users" ("id", "display_name") VALUES ($1, $2) RETURNING *;

-- name: CreateUserCourseEnrollment :one
INSERT INTO "user_course_enrollments" ("user", "course") VALUES ($1, $2) RETURNING *;
