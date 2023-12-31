// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.20.0

package db

import (
	"context"
	"database/sql"
	"fmt"
)

type DBTX interface {
	ExecContext(context.Context, string, ...interface{}) (sql.Result, error)
	PrepareContext(context.Context, string) (*sql.Stmt, error)
	QueryContext(context.Context, string, ...interface{}) (*sql.Rows, error)
	QueryRowContext(context.Context, string, ...interface{}) *sql.Row
}

func New(db DBTX) *Queries {
	return &Queries{db: db}
}

func Prepare(ctx context.Context, db DBTX) (*Queries, error) {
	q := Queries{db: db}
	var err error
	if q.createAnswerStmt, err = db.PrepareContext(ctx, createAnswer); err != nil {
		return nil, fmt.Errorf("error preparing query CreateAnswer: %w", err)
	}
	if q.createAssignmentStmt, err = db.PrepareContext(ctx, createAssignment); err != nil {
		return nil, fmt.Errorf("error preparing query CreateAssignment: %w", err)
	}
	if q.createCourseStmt, err = db.PrepareContext(ctx, createCourse); err != nil {
		return nil, fmt.Errorf("error preparing query CreateCourse: %w", err)
	}
	if q.createFeedbackStmt, err = db.PrepareContext(ctx, createFeedback); err != nil {
		return nil, fmt.Errorf("error preparing query CreateFeedback: %w", err)
	}
	if q.createQuestionStmt, err = db.PrepareContext(ctx, createQuestion); err != nil {
		return nil, fmt.Errorf("error preparing query CreateQuestion: %w", err)
	}
	if q.createReviewStmt, err = db.PrepareContext(ctx, createReview); err != nil {
		return nil, fmt.Errorf("error preparing query CreateReview: %w", err)
	}
	if q.createSubmissionStmt, err = db.PrepareContext(ctx, createSubmission); err != nil {
		return nil, fmt.Errorf("error preparing query CreateSubmission: %w", err)
	}
	if q.createUserStmt, err = db.PrepareContext(ctx, createUser); err != nil {
		return nil, fmt.Errorf("error preparing query CreateUser: %w", err)
	}
	if q.createUserCourseEnrollmentStmt, err = db.PrepareContext(ctx, createUserCourseEnrollment); err != nil {
		return nil, fmt.Errorf("error preparing query CreateUserCourseEnrollment: %w", err)
	}
	if q.getAssignmentStmt, err = db.PrepareContext(ctx, getAssignment); err != nil {
		return nil, fmt.Errorf("error preparing query GetAssignment: %w", err)
	}
	if q.getUserStmt, err = db.PrepareContext(ctx, getUser); err != nil {
		return nil, fmt.Errorf("error preparing query GetUser: %w", err)
	}
	if q.testStmt, err = db.PrepareContext(ctx, test); err != nil {
		return nil, fmt.Errorf("error preparing query Test: %w", err)
	}
	return &q, nil
}

func (q *Queries) Close() error {
	var err error
	if q.createAnswerStmt != nil {
		if cerr := q.createAnswerStmt.Close(); cerr != nil {
			err = fmt.Errorf("error closing createAnswerStmt: %w", cerr)
		}
	}
	if q.createAssignmentStmt != nil {
		if cerr := q.createAssignmentStmt.Close(); cerr != nil {
			err = fmt.Errorf("error closing createAssignmentStmt: %w", cerr)
		}
	}
	if q.createCourseStmt != nil {
		if cerr := q.createCourseStmt.Close(); cerr != nil {
			err = fmt.Errorf("error closing createCourseStmt: %w", cerr)
		}
	}
	if q.createFeedbackStmt != nil {
		if cerr := q.createFeedbackStmt.Close(); cerr != nil {
			err = fmt.Errorf("error closing createFeedbackStmt: %w", cerr)
		}
	}
	if q.createQuestionStmt != nil {
		if cerr := q.createQuestionStmt.Close(); cerr != nil {
			err = fmt.Errorf("error closing createQuestionStmt: %w", cerr)
		}
	}
	if q.createReviewStmt != nil {
		if cerr := q.createReviewStmt.Close(); cerr != nil {
			err = fmt.Errorf("error closing createReviewStmt: %w", cerr)
		}
	}
	if q.createSubmissionStmt != nil {
		if cerr := q.createSubmissionStmt.Close(); cerr != nil {
			err = fmt.Errorf("error closing createSubmissionStmt: %w", cerr)
		}
	}
	if q.createUserStmt != nil {
		if cerr := q.createUserStmt.Close(); cerr != nil {
			err = fmt.Errorf("error closing createUserStmt: %w", cerr)
		}
	}
	if q.createUserCourseEnrollmentStmt != nil {
		if cerr := q.createUserCourseEnrollmentStmt.Close(); cerr != nil {
			err = fmt.Errorf("error closing createUserCourseEnrollmentStmt: %w", cerr)
		}
	}
	if q.getAssignmentStmt != nil {
		if cerr := q.getAssignmentStmt.Close(); cerr != nil {
			err = fmt.Errorf("error closing getAssignmentStmt: %w", cerr)
		}
	}
	if q.getUserStmt != nil {
		if cerr := q.getUserStmt.Close(); cerr != nil {
			err = fmt.Errorf("error closing getUserStmt: %w", cerr)
		}
	}
	if q.testStmt != nil {
		if cerr := q.testStmt.Close(); cerr != nil {
			err = fmt.Errorf("error closing testStmt: %w", cerr)
		}
	}
	return err
}

func (q *Queries) exec(ctx context.Context, stmt *sql.Stmt, query string, args ...interface{}) (sql.Result, error) {
	switch {
	case stmt != nil && q.tx != nil:
		return q.tx.StmtContext(ctx, stmt).ExecContext(ctx, args...)
	case stmt != nil:
		return stmt.ExecContext(ctx, args...)
	default:
		return q.db.ExecContext(ctx, query, args...)
	}
}

func (q *Queries) query(ctx context.Context, stmt *sql.Stmt, query string, args ...interface{}) (*sql.Rows, error) {
	switch {
	case stmt != nil && q.tx != nil:
		return q.tx.StmtContext(ctx, stmt).QueryContext(ctx, args...)
	case stmt != nil:
		return stmt.QueryContext(ctx, args...)
	default:
		return q.db.QueryContext(ctx, query, args...)
	}
}

func (q *Queries) queryRow(ctx context.Context, stmt *sql.Stmt, query string, args ...interface{}) *sql.Row {
	switch {
	case stmt != nil && q.tx != nil:
		return q.tx.StmtContext(ctx, stmt).QueryRowContext(ctx, args...)
	case stmt != nil:
		return stmt.QueryRowContext(ctx, args...)
	default:
		return q.db.QueryRowContext(ctx, query, args...)
	}
}

type Queries struct {
	db                             DBTX
	tx                             *sql.Tx
	createAnswerStmt               *sql.Stmt
	createAssignmentStmt           *sql.Stmt
	createCourseStmt               *sql.Stmt
	createFeedbackStmt             *sql.Stmt
	createQuestionStmt             *sql.Stmt
	createReviewStmt               *sql.Stmt
	createSubmissionStmt           *sql.Stmt
	createUserStmt                 *sql.Stmt
	createUserCourseEnrollmentStmt *sql.Stmt
	getAssignmentStmt              *sql.Stmt
	getUserStmt                    *sql.Stmt
	testStmt                       *sql.Stmt
}

func (q *Queries) WithTx(tx *sql.Tx) *Queries {
	return &Queries{
		db:                             tx,
		tx:                             tx,
		createAnswerStmt:               q.createAnswerStmt,
		createAssignmentStmt:           q.createAssignmentStmt,
		createCourseStmt:               q.createCourseStmt,
		createFeedbackStmt:             q.createFeedbackStmt,
		createQuestionStmt:             q.createQuestionStmt,
		createReviewStmt:               q.createReviewStmt,
		createSubmissionStmt:           q.createSubmissionStmt,
		createUserStmt:                 q.createUserStmt,
		createUserCourseEnrollmentStmt: q.createUserCourseEnrollmentStmt,
		getAssignmentStmt:              q.getAssignmentStmt,
		getUserStmt:                    q.getUserStmt,
		testStmt:                       q.testStmt,
	}
}
