package repository

import (
	"database/sql"
	"kstation_backend/internal/models"
)

type DatabaseRepo interface {
	Connection() *sql.DB
	InsertUser(user models.User) (int, error)
	UpdateUser(u models.User) error
	GetUserByEmail(email string) (*models.User, error)
	GetUserByID(id int) (*models.User, error)
	ResetPassword(id int, password string) error
	InsertLesson(lesson models.Lesson) (int, error)
	UpdateLesson(l models.Lesson) error
	GetLessonByID(id int) (*models.Lesson, error)
	AllLessons(how int) ([]*models.Lesson, error)
	AllLessonsByUser(id int, how int) ([]*models.Lesson, error)
	InsertComment(comment models.Comment) (int, error)
	GetCommentByID(id int) (*models.Comment, error)
	AllCommentsByLessonId(LessonId int) ([]*models.Comment, error)
	AllCommentsByUserId(UserId int) ([]*models.Comment, error)
	UpdateComment(c models.Comment) error
	DeleteComment(id int) error
}