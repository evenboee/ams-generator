package generator

import (
	"context"
	"database/sql"
	"log"
	"sync"
	"time"

	db "github.com/evenboee/ams-generator/db/sqlc"
	"github.com/evenboee/ams-generator/utils/random"

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

		for _, userID := range userIDs {
			err := store.CreateUserCourseEnrollment(ctx, db.CreateUserCourseEnrollmentParams{
				User:   userID,
				Course: course.ID,
			})
			if err != nil {
				return err
			}
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

func GenerateAsync(ctx context.Context, store db.Store, config *Config) error {

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
	var wg sync.WaitGroup

	assignmentCount := config.AssignmentsPerCourse * config.CourseCount
	submissionCount := assignmentCount * config.SubmissionsPerUser * config.UserCount

	progressBar := progressbar.Default(int64(submissionCount))

	log.Println("Adding courses")
	for courseIdx := int32(0); courseIdx < config.CourseCount; courseIdx++ {
		wg.Add(1)
		go func(courseIdx int32) {

			defer wg.Done()

			err := func() error {
				course := RandomCourse()
				course.ID, err = store.CreateCourse(ctx, db.CreateCourseParams{
					Year: course.Year,
					Code: course.Code,
				})
				if err != nil {
					return err
				}

				for _, userID := range userIDs {
					err := store.CreateUserCourseEnrollment(ctx, db.CreateUserCourseEnrollmentParams{
						User:   userID,
						Course: course.ID,
					})
					if err != nil {
						return err
					}
				}

				// assignmentCount := config.AssignmentsPerCourse
				// submissionCount := config.SubmissionsPerUser * config.AssignmentsPerCourse * config.UserCount
				// reviewCount := submissionCount * config.ReviewsPerSubmission
				// // feedbackCount := submissionCount * config.QuestionsPerAssignment * config.ReviewsPerSubmission
				// log.Printf("Adding %d submissions %d to assignment for course %d\n", submissionCount, assignmentCount, course.ID)
				// progressBar := progressbar.Default(int64(reviewCount))

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
							// progressBar.Add(1)
						}

						progressBar.Add(1)
					}
				}
				return nil
			}()

			if err != nil {
				log.Printf("Error adding course %d: %v\n", courseIdx, err)
			}
		}(courseIdx)
	}

	wg.Wait()

	return nil
}

type generatedAssignment struct {
	ID        int32
	Questions []int32
}

func GenerateAsync2(ctx context.Context, store db.Store, config *Config2) <-chan error {
	errChan := make(chan error, 1)

	go func() {
		defer close(errChan)

		log.Printf("Adding %d courses", config.CourseCount)
		baseStore := store
		// baseStore, storeCloser, err := genStore()
		// if err != nil {
		// 	errChan <- err
		// 	return
		// }
		// defer storeCloser()

		courseBar := progressbar.Default(int64(config.CourseCount))

		courses := make([]int32, 0, config.CourseCount*config.CourseYears)
		for courseIdx := int32(0); courseIdx < config.CourseCount; courseIdx++ {
			baseCourse := RandomCourse()
			for yearOffset := int32(0); yearOffset < config.CourseYears; yearOffset++ {
				courseID, err := baseStore.CreateCourse(ctx, db.CreateCourseParams{
					Code: baseCourse.Code,
					Year: baseCourse.Year + yearOffset,
				})
				if err != nil {
					errChan <- err
					return
				}
				courses = append(courses, courseID)
			}
			courseBar.Add(1)
		}

		// submission count
		submissionCount := config.CourseCount * config.CourseYears // course runs
		submissionCount *= config.CourseAssignments                // assignments per course
		submissionCount *= config.CourseUsers                      // users that have had the course

		log.Printf("Adding %d assignments", submissionCount)
		bar := progressbar.Default(int64(submissionCount))

		var wg sync.WaitGroup
		for courseIdx := int32(0); courseIdx < int32(len(courses)); courseIdx++ {
			wg.Add(1)
			go func(courseIdx int32) {
				// log.Printf("Starting %d ...", courseIdx)
				defer wg.Done()

				var err error
				courseStore := store
				// courseStore, courseStoreCloser, err := genStore()
				// if err != nil {
				// 	errChan <- err
				// 	return
				// }
				// defer courseStoreCloser()
				courseID := courses[courseIdx]

				// genereate users
				userIDs := make([]string, config.CourseUsers)
				for i := int32(0); i < config.CourseUsers; i++ {
					userID := random.String(8)
					err := courseStore.CreateUser(ctx, db.CreateUserParams{
						ID:          userID,
						DisplayName: random.String(12),
					})
					if err != nil {
						errChan <- err
						return
					}
					userIDs[i] = userID

					// enroll users
					if err := courseStore.CreateUserCourseEnrollment(ctx, db.CreateUserCourseEnrollmentParams{
						User:   userID,
						Course: courseID,
					}); err != nil {
						errChan <- err
						return
					}
				}
				lenUsers := int32(len(userIDs))

				// generate assignments
				assignments := make([]generatedAssignment, config.CourseAssignments)
				for assignmentIdx := int32(0); assignmentIdx < config.CourseAssignments; assignmentIdx++ {
					assignment := RandomAssignment(courseID, config.SubmissionReviews)
					assignment.ID, err = courseStore.CreateAssignment(ctx, db.CreateAssignmentParams{
						Course:               assignment.Course,
						TimeDue:              assignment.TimeDue,
						Name:                 assignment.Name,
						ReviewsPerSubmission: assignment.ReviewsPerSubmission,
					})
					if err != nil {
						errChan <- err
						return
					}

					questions := make([]int32, config.AssignmentQuestions)
					for questionIdx := int32(0); questionIdx < config.AssignmentQuestions; questionIdx++ {
						question := RandomQuestion(assignment.ID, questionIdx+1)
						question.ID, err = courseStore.CreateQuestion(ctx, db.CreateQuestionParams{
							Assignment: question.Assignment,
							Prompt:     question.Prompt,
							Order:      question.Order,
						})
						if err != nil {
							errChan <- err
							return
						}

						questions[questionIdx] = question.ID
					}

					assignments[assignmentIdx] = generatedAssignment{
						ID:        assignment.ID,
						Questions: questions,
					}
				}

				// generate submissions
				// var subWg sync.WaitGroup
				for _, assignment := range assignments {
					for userIdx, userID := range userIDs {
						// subWg.Add(1)
						// go func(userIdx int32, userID string, assignment generatedAssignment) {
						// 	defer subWg.Done()

						submission := RandomSubmission(assignment.ID, userID)
						submission.ID, err = courseStore.CreateSubmission(ctx, db.CreateSubmissionParams{
							User:       submission.User,
							Assignment: submission.Assignment,
						})
						if err != nil {
							errChan <- err
							return
						}

						reviews := make([]int32, config.SubmissionReviews)
						for reviewIdx := int32(0); reviewIdx < config.SubmissionReviews; reviewIdx++ {
							reviewerIdx := (int32(userIdx) + reviewIdx + 1) % lenUsers
							reviewerID := userIDs[reviewerIdx]
							n := time.Now()
							reviewID, err := courseStore.CreateReview(ctx, db.CreateReviewParams{
								Submission: submission.ID,
								ReviewerID: reviewerID,
								FinishedAt: sql.NullTime{Time: n, Valid: true},
								CreatedAt:  n,
							})
							if err != nil {
								errChan <- err
								return
							}

							reviews[reviewIdx] = reviewID
						}

						// generate answers
						answers := make([]int32, config.AssignmentQuestions)
						for answerIdx, question := range assignment.Questions {
							answer := RandomAnswer(question, submission.ID)
							answer.ID, err = courseStore.CreateAnswer(ctx, db.CreateAnswerParams{
								Question:   question,
								Submission: submission.ID,
								Answer:     random.String(50),
							})
							if err != nil {
								errChan <- err
								return
							}

							answers[answerIdx] = answer.ID
						}

						// time.Sleep(20 * time.Second)

						// generate feedback
						for _, reviewID := range reviews {
							for _, answerID := range answers {
								// feedback := RandomFeedback(reviewID, answerID)
								_, err := courseStore.CreateFeedback(ctx, db.CreateFeedbackParams{
									Review:   reviewID,
									Answer:   answerID,
									Rating:   int32(random.Intr(1, 5)),
									Feedback: random.String(20),
								})
								if err != nil {
									errChan <- err
									return
								}
							}
						}

						bar.Add(1)
						// }(int32(userIdx), userID, assignment)
					}
					// bar.Add(1)
				}

				// subWg.Wait()
			}(courseIdx)
		}

		wg.Wait()
	}()

	return errChan
}
