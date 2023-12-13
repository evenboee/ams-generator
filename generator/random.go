package generator

import (
	"database/sql"
	"time"

	db "github.com/evenboee/ams-generator/db/sqlc"
	"github.com/evenboee/ams-generator/utils/random"
)

func RandomAnswer(questionID int32, submissionID int32) db.Answer {
	return db.Answer{
		ID:         random.PosInt32(),
		Question:   questionID,
		Submission: submissionID,
		Answer:     random.String(50),
	}
}

func RandomAssignment(courseID int32, submissionsPerAssignment int32) db.Assignment {
	return db.Assignment{
		ID:                   random.PosInt32(),
		Name:                 random.String(10),
		Course:               courseID,
		ReviewsPerSubmission: submissionsPerAssignment,
		TimeDue:              time.Now().Add(24 * time.Hour * time.Duration(random.Intr(5, 120))),
	}
}

func RandomCourse() db.Course {
	currYear := time.Now().Year()

	return db.Course{
		ID:   random.PosInt32(),
		Code: random.StringWith(3, random.Uppercase) + random.StringWith(4, random.Nums),
		Year: int32(random.Intr(currYear, currYear+2)),
	}
}

func RandomFeedback(reviewID int32, answerID int32) db.Feedback {
	return db.Feedback{
		ID:       random.PosInt32(),
		Review:   reviewID,
		Answer:   answerID,
		Rating:   int32(random.Intr(1, 5)),
		Feedback: random.String(50),
	}
}

// func RandomQuestions(assignmentID int32, count int32) []db.Question {
// 	questions := make([]db.Question, count)
// 	for i := int32(0); i < count; i++ {
// 		questions[i] = RandomQuestion(assignmentID, i+1)
// 	}
// 	return questions
// }

func RandomQuestion(assignmentID int32, order int32) db.Question {
	return db.Question{
		ID:         random.PosInt32(),
		Assignment: assignmentID,
		Prompt:     random.String(50),
		Order:      order,
	}
}

func RandomReview(submissionID int32, reviewerID string) db.Review {
	return db.Review{
		ID:         random.PosInt32(),
		Submission: submissionID,
		ReviewerID: reviewerID,
		FinishedAt: sql.NullTime{Time: time.Now(), Valid: true},
		CreatedAt:  time.Now(),
	}
}

func RandomSubmission(assignmentID int32, userID string) db.Submission {
	return db.Submission{
		ID:         random.PosInt32(),
		User:       userID,
		Assignment: assignmentID,
		CreatedAt:  time.Now(),
	}
}

func RandomUser() db.User {
	return db.User{
		ID:          random.String(8),
		DisplayName: random.String(6) + " " + random.String(6),
	}
}
