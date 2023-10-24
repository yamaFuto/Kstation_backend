package models

import "time"

type Comment struct {
	ID           int       `json:"id"`
	LessonId     int       `json:"lesson_id"`
	UserId       int       `json:"user_id"`
	Year         int       `json:"year"`
	Term         string    `json:"Term"`
	Comment      string    `json:"comment"`
	TestOrReport string    `json:"test_or_report"`
	Star         int       `json:"star"`
	CreatedAt    time.Time `json:"-"`
	UpdatedAt    time.Time `json:"-"`
}
