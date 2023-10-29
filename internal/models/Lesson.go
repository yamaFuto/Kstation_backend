package models

import "time"

type Lesson struct {
	ID             int       `json:"id"`
	LessonName     string    `json:"lesson_name"`
	TeacherName    string    `json:"teacher_name"`
	AvgStar        float32   `json:"avg_star"`
	AboutAvgStar   int       `json:"about_avg_star"`
	CommentNumbers int       `json:"comment_numbers"`
	CreatedAt      time.Time `json:"-"`
	UpdatedAt      time.Time `json:"-"`
}
