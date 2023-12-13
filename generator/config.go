package generator

import (
	"os"

	"gopkg.in/yaml.v3"
)

type Config struct {
	// Scaling factor. The other resources (except users) will be multiplied by this number
	CourseCount int32 `yaml:"course_count"`
	// How many users to generate
	UserCount int32 `yaml:"user_count"`
	// How many questions per course
	AssignmentsPerCourse int32 `yaml:"assignments_per_course"`
	// How many questions per assignment
	QuestionsPerAssignment int32 `yaml:"questions_per_assignment"`
	// How many reviews per submission
	ReviewsPerSubmission int32 `yaml:"reviews_per_submission"`
	// How many submissions per user per assignment
	SubmissionsPerUser int32 `yaml:"submissions_per_user"`
}

// // How many submissions per assignment
// SubmissionsPerAssignment int32 `yaml:"submissions_per_assignment"`

func defaultConfig() *Config {
	return &Config{
		CourseCount:            1,
		UserCount:              100,
		AssignmentsPerCourse:   5,
		QuestionsPerAssignment: 10,
		ReviewsPerSubmission:   3,
		SubmissionsPerUser:     1,
	}
}

func LoadConfig(path string) (*Config, error) {
	content, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	config := defaultConfig()
	err = yaml.Unmarshal(content, config)
	return config, err
}
