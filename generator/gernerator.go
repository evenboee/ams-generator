package generator

import (
	"context"

	db "github.com/evenboee/ams-generator/db/sqlc"
)

func Generate(ctx context.Context, store db.Store, config *Config) error {

	/*
		Generate
		1. users
		2. courses
			1. enrollments
			2. assignments

	*/

	usersIDs := make([]string, config.UserCount)
	for i := 0; i < config.UserCount; i++ {
		user := RandomUser()
		_, err := store.CreateUser(ctx, db.CreateUserParams{
			ID:          user.ID,
			DisplayName: user.DisplayName,
		})
		if err != nil {
			return err
		}
		usersIDs[i] = user.ID
	}

	for i := 0; i < config.CourseCount; i++ {
		course := RandomCourse()
		_, err := store.CreateCourse(ctx, db.CreateCourseParams{
			Code: course.Code,
			Year: course.Year,
		})
		if err != nil {
			return err
		}

		assignments := make([]int32, config.AssignmentsPerCourse)
		for j := 0; j < config.AssignmentsPerCourse; j++ {
			assignment := RandomAssignment(course.ID)
			assignment.Course = int32(i + 1)
			_, err := store.CreateAssignment(ctx, db.CreateAssignmentParams{
				Course:  assignment.Course,
				TimeDue: assignment.TimeDue,
			})
			if err != nil {
				return err
			}

			assignments[j] = assignment.Course
		}

		submissions := make([]int32, 0, config.AssignmentsPerCourse*config.SubmissionsPerUser*config.UserCount)
		for _, userID := range usersIDs {
			for _, assignmentID := range assignments {
				submission := RandomSubmission(assignmentID, userID)
				_, err := store.CreateSubmission(ctx, db.CreateSubmissionParams{
					ID:         submission.ID,
					User:       submission.User,
					Assignment: submission.Assignment,
				})
				if err != nil {
					return err
				}

				submissions = append(submissions, submission.ID)
			}
		}

	}

	return nil
}
