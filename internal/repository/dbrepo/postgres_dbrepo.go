package dbrepo

import (
	"context"
	"database/sql"
	"kstation_backend/internal/models"
	"log"
	"time"

	"golang.org/x/crypto/bcrypt"
)

type PostgresDBRepo struct {
	DB *sql.DB
}

const dbTimeout = time.Second * 3

func (m *PostgresDBRepo) Connection() *sql.DB {
	return m.DB
}

func (m *PostgresDBRepo) InsertUser(user models.User) (int, error) {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), 12)
	if err != nil {
		return 0, err
	}

	var newID int
	stmt := `insert into users (email, first_name, last_name, password, image, is_admin, created_at, updated_at)
		values ($1, $2, $3, $4, $5, $6, $7, $8) returning id`

	err = m.DB.QueryRowContext(ctx, stmt,
		user.Email,
		user.FirstName,
		user.LastName,
		user.Image,
		hashedPassword,
		user.IsAdmin,
		time.Now(),
		time.Now(),
	).Scan(&newID)

	if err != nil {
		return 0, err
	}

	return newID, nil
}

func (m *PostgresDBRepo) GetUserByID(id int) (*models.User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	query := `
		select
			id, email, first_name, last_name, password, image, is_admin, created_at, updated_at
		from users
		where
		    id = $1`

	var user models.User
	row := m.DB.QueryRowContext(ctx, query, id)

	err := row.Scan(
		&user.ID,
		&user.Email,
		&user.FirstName,
		&user.LastName,
		&user.Password,
		&user.Image,
		&user.IsAdmin,
		&user.CreatedAt,
		&user.UpdatedAt,
	)

	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (m *PostgresDBRepo) UpdateUser(u models.User) error {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	stmt := `update users set
		email = $1,
		first_name = $2,
		last_name = $3,
		image = $4,
		is_admin = $5,
		updated_at = $6
		where id = $7
	`

	_, err := m.DB.ExecContext(ctx, stmt,
		u.Email,
		u.FirstName,
		u.LastName,
		u.Image,
		u.IsAdmin,
		time.Now(),
		u.ID,
	)

	if err != nil {
		return err
	}

	return nil
}

func (m *PostgresDBRepo) GetUserByEmail(email string) (*models.User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	query := `
		select
			id, email, first_name, last_name, password, image, is_admin, created_at, updated_at
		from users
		where
		    email = $1`

	var user models.User
	row := m.DB.QueryRowContext(ctx, query, email)

	err := row.Scan(
		&user.ID,
		&user.Email,
		&user.FirstName,
		&user.LastName,
		&user.Password,
		&user.Image,
		&user.IsAdmin,
		&user.CreatedAt,
		&user.UpdatedAt,
	)

	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (m *PostgresDBRepo) ResetPassword(id int, password string) error {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), 12)
	if err != nil {
		return err
	}

	stmt := `update users set password = $1 where id = $2`
	_, err = m.DB.ExecContext(ctx, stmt, hashedPassword, id)
	if err != nil {
		return err
	}

	return nil
}

func (m *PostgresDBRepo) InsertLesson(lesson models.Lesson) (int, error) {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	var newID int
	stmt := `insert into lessons (lesson_name, teacher_name, avg_star, about_avg_star, comment_numbers, created_at, updated_at)
		values ($1, $2, $3, $4, $5, $6, $7) returning id`

	err := m.DB.QueryRowContext(ctx, stmt,
		lesson.LessonName,
		lesson.TeacherName,
		lesson.AvgStar,
		lesson.AboutAvgStar,
		lesson.CommentNumbers,
		time.Now(),
		time.Now(),
	).Scan(&newID)

	if err != nil {
		return 0, err
	}

	return newID, nil
}

func (m *PostgresDBRepo) GetLessonByID(id int) (*models.Lesson, error) {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	query := `
		select
			id, lesson_name, teacher_name, avg_star, about_avg_star, comment_numbers, created_at, updated_at
		from lessons
		where
		    id = $1`

	var lesson models.Lesson
	row := m.DB.QueryRowContext(ctx, query, id)

	err := row.Scan(
		&lesson.ID,
		&lesson.LessonName,
		&lesson.TeacherName,
		&lesson.AvgStar,
		&lesson.AboutAvgStar,
		&lesson.CommentNumbers,
		&lesson.CreatedAt,
		&lesson.UpdatedAt,
	)

	if err != nil {
		return nil, err
	}

	return &lesson, nil
}

func (m *PostgresDBRepo) UpdateLesson(l models.Lesson) error {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	stmt := `update lessons set
		avg_star = $1,
		about_avg_star = $2,
		comment_numbers = $3,
		updated_at = $4
		where id = $5
	`

	_, err := m.DB.ExecContext(ctx, stmt,
		l.AvgStar,
		l.AboutAvgStar,
		l.CommentNumbers,
		time.Now(),
		l.ID,
	)

	if err != nil {
		return err
	}

	return nil
}

func (m *PostgresDBRepo) AllLessons(how int) ([]*models.Lesson, error) {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	query := `select id, lesson_name, teacher_name, avg_star, about_avg_star, comment_numbers, created_at, updated_at
	from lessons order by lesson_name`

	if how == 1 {
		query = `select id, lesson_name, teacher_name, avg_star, about_avg_star, comment_numbers, created_at, updated_at
	from lessons order by created_at`
	} else if how == 2 {
		query = `select id, lesson_name, teacher_name, avg_star, about_avg_star, comment_numbers, created_at, updated_at
	from lessons order by created_at desc`
	} else if how == 3 {
		query = `select id, lesson_name, teacher_name, avg_star, about_avg_star, comment_numbers, created_at, updated_at
	from lessons order by about_avg_star desc`
	} else if how == 0 {
		query = `select id, lesson_name, teacher_name, avg_star, about_avg_star, comment_numbers, created_at, updated_at
	from lessons order by lesson_name`
	}

	rows, err := m.DB.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var lessons []*models.Lesson

	for rows.Next() {
		var lesson models.Lesson
		err := rows.Scan(
			&lesson.ID,
			&lesson.LessonName,
			&lesson.TeacherName,
			&lesson.AvgStar,
			&lesson.AboutAvgStar,
			&lesson.CommentNumbers,
			&lesson.CreatedAt,
			&lesson.UpdatedAt,
		)
		if err != nil {
			log.Println("Error scanning", err)
			return nil, err
		}

		lessons = append(lessons, &lesson)
	}

	return lessons, nil
}

func (m *PostgresDBRepo) InsertComment(comment models.Comment) (int, error) {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	var newID int
	stmt := `insert into comments (lesson_id, user_id, year, term, comment, test_or_report, star, created_at, updated_at)
		values ($1, $2, $3, $4, $5, $6, $7, $8, $9) returning id`

	err := m.DB.QueryRowContext(ctx, stmt,
		comment.LessonId,
		comment.UserId,
		comment.Year,
		comment.Term,
		comment.Comment,
		comment.TestOrReport,
		comment.Star,
		time.Now(),
		time.Now(),
	).Scan(&newID)

	if err != nil {
		return 0, err
	}

	return newID, nil
}

func (m *PostgresDBRepo) GetCommentByID(id int) (*models.Comment, error) {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	query := `
		select
			id, lesson_id, user_id, year, term, comment, test_or_report, star, created_at, updated_at
		from comments
		where
		    id = $1`

	var comment models.Comment
	row := m.DB.QueryRowContext(ctx, query, id)

	err := row.Scan(
		&comment.ID,
		&comment.LessonId,
		&comment.UserId,
		&comment.Year,
		&comment.Term,
		&comment.Comment,
		&comment.TestOrReport,
		&comment.Star,
		&comment.CreatedAt,
		&comment.UpdatedAt,
	)

	if err != nil {
		return nil, err
	}

	return &comment, nil
}

func (m *PostgresDBRepo) AllCommentsByLessonId(LessonId int) ([]*models.Comment, error) {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	query := `select id, lesson_id, user_id, year, term, comment, test_or_report, star, created_at, updated_at
						from comments
						where lesson_id = $1
						order by id desc`

	rows, err := m.DB.QueryContext(ctx, query, LessonId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var comments []*models.Comment

	for rows.Next() {
		var comment models.Comment
		err := rows.Scan(
			&comment.ID,
			&comment.LessonId,
			&comment.UserId,
			&comment.Year,
			&comment.Term,
			&comment.Comment,
			&comment.TestOrReport,
			&comment.Star,
			&comment.CreatedAt,
			&comment.UpdatedAt,
		)
		if err != nil {
			log.Println("Error scanning", err)
			return nil, err
		}

		comments = append(comments, &comment)
	}

	return comments, nil
}

func (m *PostgresDBRepo) AllCommentsByUserId(UserId int) ([]*models.Comment, error) {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	query := `select id, lesson_id, user_id, year, term, comment, test_or_report, star, created_at, updated_at
						from comments
						where user_id = $1
						order by id`

	rows, err := m.DB.QueryContext(ctx, query, UserId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var comments []*models.Comment

	for rows.Next() {
		var comment models.Comment
		err := rows.Scan(
			&comment.ID,
			&comment.LessonId,
			&comment.UserId,
			&comment.Year,
			&comment.Term,
			&comment.Comment,
			&comment.TestOrReport,
			&comment.Star,
			&comment.CreatedAt,
			&comment.UpdatedAt,
		)
		if err != nil {
			log.Println("Error scanning", err)
			return nil, err
		}

		comments = append(comments, &comment)
	}

	return comments, nil
}

func (m *PostgresDBRepo) UpdateComment(c models.Comment) error {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	stmt := `update comments set
		comment = $1,
		year = $2,
		term = $3,
		test_or_report = $4,
		star = $5,
		updated_at = $6
		where id = $7
	`

	_, err := m.DB.ExecContext(ctx, stmt,
		c.Comment,
		c.Year,
		c.Term,
		c.TestOrReport,
		c.Star,
		time.Now(),
		c.ID,
	)

	if err != nil {
		return err
	}

	return nil
}

func (m *PostgresDBRepo) DeleteComment(id int) error {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	stmt := `delete from comments where id = $1`

	_, err := m.DB.ExecContext(ctx, stmt, id)
	if err != nil {
		return err
	}

	return nil
}