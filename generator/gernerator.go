package generator

import (
	"context"
	"log"

	db "github.com/evenboee/ams-generator/db/sqlc"

	"github.com/schollz/progressbar/v3"
)

func Generate(ctx context.Context, store db.Store, config *Config) error {

	log.Println("Adding users")
	userIDs := make([]string, config.UserCount)
	for i := int32(0); i < config.UserCount; i++ {
		user := RandomUser()
		err := store.CreateUser(ctx, db.CreateUserParams{
			ID:          user.ID,
			DisplayName: user.DisplayName,
		})
		if err != nil {
			return err
		}
		userIDs[i] = user.ID
	}

	var err error

	log.Println("Adding courses")
	for courseIdx := int32(0); courseIdx < config.CourseCount; courseIdx++ {
		course := RandomCourse()
		course.ID, err = store.CreateCourse(ctx, db.CreateCourseParams{
			Year: course.Year,
			Code: course.Code,
		})
		if err != nil {
			return err
		}

		assignmentCount := config.AssignmentsPerCourse
		submissionCount := config.SubmissionsPerUser * config.AssignmentsPerCourse * config.UserCount
		reviewCount := submissionCount * config.ReviewsPerSubmission
		// feedbackCount := submissionCount * config.QuestionsPerAssignment * config.ReviewsPerSubmission
		log.Printf("Adding %d submissions %d to assignment for course %d\n", submissionCount, assignmentCount, course.ID)
		progressBar := progressbar.Default(int64(reviewCount))

		for assignmentIdx := int32(0); assignmentIdx < config.AssignmentsPerCourse; assignmentIdx++ {
			assignment := RandomAssignment(course.ID, config.ReviewsPerSubmission)
			assignment.ID, err = store.CreateAssignment(ctx, db.CreateAssignmentParams{
				// ID:                   assignment.ID,
				Course:               assignment.Course,
				TimeDue:              assignment.TimeDue,
				Name:                 assignment.Name,
				ReviewsPerSubmission: assignment.ReviewsPerSubmission,
			})
			if err != nil {
				return err
			}

			questions := make([]int32, config.QuestionsPerAssignment)
			for questionIdx := int32(0); questionIdx < config.QuestionsPerAssignment; questionIdx++ {
				question := RandomQuestion(assignment.ID, questionIdx+1)
				question.ID, err = store.CreateQuestion(ctx, db.CreateQuestionParams{
					Assignment: question.Assignment,
					Prompt:     question.Prompt,
					Order:      question.Order,
				})
				if err != nil {
					return err
				}

				questions[questionIdx] = question.ID
			}

			submissions := make([]int32, 0, len(userIDs)) // add: * config.SubmissionsPerUser
			for submissionIdx, userID := range userIDs {
				submission := RandomSubmission(assignment.ID, userID)
				submission.ID, err = store.CreateSubmission(ctx, db.CreateSubmissionParams{
					User:       submission.User,
					Assignment: submission.Assignment,
				})
				if err != nil {
					return err
				}
				submissions = append(submissions, submission.ID)

				answers := make([]int32, len(questions))
				for answerIdx, question := range questions {
					answer := RandomAnswer(question, submission.ID)
					answer.ID, err = store.CreateAnswer(ctx, db.CreateAnswerParams{
						Question:   answer.Question,
						Submission: answer.Submission,
						Answer:     answer.Answer,
					})
					if err != nil {
						return err
					}

					answers[answerIdx] = answer.ID
				}

				// TODO: make the ones that are missing some reviews go back and add them after all reviews are added
				for submissionIdxOffset := 1; submissionIdxOffset <= int(config.ReviewsPerSubmission); submissionIdxOffset++ {
					reviewedSubmissionIdx := submissionIdx - submissionIdxOffset
					if reviewedSubmissionIdx < 0 {
						break
					}

					review := RandomReview(submissions[reviewedSubmissionIdx], userID)
					review.ID, err = store.CreateReview(ctx, db.CreateReviewParams{
						// ID:         review.ID,
						Submission: review.Submission,
						ReviewerID: review.ReviewerID,
						FinishedAt: review.FinishedAt,
						CreatedAt:  review.CreatedAt,
					})
					if err != nil {
						return err
					}

					for feedbackIdx := int32(0); feedbackIdx < config.QuestionsPerAssignment; feedbackIdx++ {
						feedback := RandomFeedback(review.ID, answers[feedbackIdx])
						_, err := store.CreateFeedback(ctx, db.CreateFeedbackParams{
							Review:   feedback.Review,
							Answer:   feedback.Answer,
							Rating:   feedback.Rating,
							Feedback: feedback.Feedback,
						})
						if err != nil {
							return err
						}
					}
					progressBar.Add(1)
				}
			}
		}
	}

	return nil
}
