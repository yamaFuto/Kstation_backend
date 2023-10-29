package dbrepo

import (
	"database/sql"
	"errors"
	"time"
	"kstation_backend/internal/models"
)

type TestDBRepo struct{}

func (m *TestDBRepo) Connection() *sql.DB {
	return nil
}

func (m *TestDBRepo) GetUserByID(id int) (*models.User, error) {
	var user = models.User{}
	if id == 1 {
		user = models.User{
			ID: 1,
			FirstName: "Admin",
			LastName: "User",
			Email: "admin@example.com",
		}
		return &user, nil
	}

	return nil, errors.New("user not found")
}

func (m *TestDBRepo) GetUserByEmail(email string) (*models.User, error) {
	if email == "admin@example.com" {
		user := models.User{
			ID: 1,
			FirstName: "Admin",
			LastName: "User",
			Email: "admin@example.com",
			Password: "$2a$14$ajq8Q7fbtFRQvXpdCq7Jcuy.Rx1h/L4J60Otx.gyNLbAYctGMJ9tK",
			IsAdmin: 1,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}
		return &user, nil
	}
	return nil, errors.New("not found")
}

func (m *TestDBRepo) UpdateUser(u models.User) error {
	if u.ID == 1 {
		return nil
	}
	return errors.New("update failed - no user found")
}

func (m *TestDBRepo) InsertUser(user models.User) (int, error) {
	return 2, nil
}

func (m *TestDBRepo) ResetPassword(id int, password string) error {
	if id == 1 {
		return nil
	}
	return errors.New("user not found")
}

func (m *TestDBRepo) InsertLesson(lesson models.Lesson) (int, error) {
	return 2, nil
}

func (m *TestDBRepo) UpdateLesson(l models.Lesson) error {
	if l.ID == 1 {
		return nil
	}
	return errors.New("lesson not found")
}

func (m *TestDBRepo) GetLessonByID(id int) (*models.Lesson, error) {
	if id == 1 {
		lesson := models.Lesson{
			ID: 1,
			LessonName: "Math",
			TeacherName: "Yamada",
			AvgStar: 0.0,
			CommentNumbers: 0,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}
		return &lesson, nil
	}

	return nil, errors.New("not found")
}

func (m *TestDBRepo) AllLessons() ([]*models.Lesson, error) {
	var lessons []*models.Lesson

	return lessons, nil
}

func (m *TestDBRepo) InsertComment(comment models.Comment) (int, error) {
	return 2, nil
}

func (m *TestDBRepo) GetCommentByID(id int) (*models.Comment, error) {
	if id == 1 {
		comment := models.Comment{
			ID: 1,
			LessonId: 1,
			UserId: 1,
			Year: 2023,
			Term: "former",
			Comment: "this is a test",
			TestOrReport: "report",
			Star: 3,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}

		return &comment, nil
	}

	return nil, errors.New("comment not found")
}

func (m *TestDBRepo) AllCommentsByLessonId(LessonId int) ([]*models.Comment, error) {
	if LessonId == 1 {
		var comments []*models.Comment

		return comments, nil
	}

	return nil, errors.New("comments are not found")
}

func (m *TestDBRepo) AllCommentsByUserId(UserId int) ([]*models.Comment, error) {
	if UserId == 1 {
		var comments []*models.Comment

		return comments, nil
	}

	return nil, errors.New("comments are not found")
}

func (m *TestDBRepo) UpdateComment(c models.Comment) error {
	if c.ID == 1 {
		return nil
	}

	return errors.New("comment not found")
}

func (m *TestDBRepo) DeleteComment(id int) error {
	if id == 1 {
		return nil
	}

	return errors.New("commet not found")
}