// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.20.0
// source: users.sql

package db

import (
	"context"
)

const createUser = `-- name: CreateUser :one
INSERT INTO "users" ("id", "display_name") VALUES ($1, $2) RETURNING id, display_name
`

type CreateUserParams struct {
	ID          string `json:"id"`
	DisplayName string `json:"display_name"`
}

func (q *Queries) CreateUser(ctx context.Context, arg CreateUserParams) (User, error) {
	row := q.db.QueryRowContext(ctx, createUser, arg.ID, arg.DisplayName)
	var i User
	err := row.Scan(&i.ID, &i.DisplayName)
	return i, err
}

const createUserCourseEnrollment = `-- name: CreateUserCourseEnrollment :one
INSERT INTO "user_course_enrollments" ("user", "course") VALUES ($1, $2) RETURNING "user", course
`

type CreateUserCourseEnrollmentParams struct {
	User   string `json:"user"`
	Course int32  `json:"course"`
}

func (q *Queries) CreateUserCourseEnrollment(ctx context.Context, arg CreateUserCourseEnrollmentParams) (UserCourseEnrollment, error) {
	row := q.db.QueryRowContext(ctx, createUserCourseEnrollment, arg.User, arg.Course)
	var i UserCourseEnrollment
	err := row.Scan(&i.User, &i.Course)
	return i, err
}

const getUser = `-- name: GetUser :one
SELECT id, display_name FROM "users" WHERE "id" = $1
`

func (q *Queries) GetUser(ctx context.Context, id string) (User, error) {
	row := q.db.QueryRowContext(ctx, getUser, id)
	var i User
	err := row.Scan(&i.ID, &i.DisplayName)
	return i, err
}
