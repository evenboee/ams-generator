// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.20.0
// source: assignments.sql

package db

import (
	"context"
	"time"
)

const createAssignment = `-- name: CreateAssignment :one
INSERT INTO "assignments" ("course", "time_due") VALUES ($1, $2) RETURNING course, time_due
`

type CreateAssignmentParams struct {
	Course  int32     `json:"course"`
	TimeDue time.Time `json:"time_due"`
}

func (q *Queries) CreateAssignment(ctx context.Context, arg CreateAssignmentParams) (Assignment, error) {
	row := q.db.QueryRowContext(ctx, createAssignment, arg.Course, arg.TimeDue)
	var i Assignment
	err := row.Scan(&i.Course, &i.TimeDue)
	return i, err
}