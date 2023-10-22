package models

import "time"

type Lesson struct {
	ID             int       `json:"id"`
	LessonName     string    `json:"lesson_name"`
	TeacherNsme    string    `json:"teacher_name"`
	AvgStar        float32   `json:"avg_star"`
	CommentNumbers string    `json:"comment_numbers"`
	CreatedAt      time.Time `json:"-"`
	UpdatedAt      time.Time `json:"-"`
}
