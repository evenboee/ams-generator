-- name: GetUser :one
SELECT * FROM "users" WHERE "id" = $1;

-- name: CreateUser :exec
INSERT INTO "users" ("id", "display_name") VALUES ($1, $2);

-- name: CreateUserCourseEnrollment :exec
INSERT INTO "user_course_enrollments" ("user", "course") VALUES ($1, $2);
