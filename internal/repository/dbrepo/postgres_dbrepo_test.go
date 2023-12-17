package dbrepo

import (
	"database/sql"
	"fmt"
	"math"
	"kstation_backend/internal/models"
	"kstation_backend/internal/repository"
	"log"
	"os"
	"testing"
	"time"

	_ "github.com/jackc/pgconn"
	_ "github.com/jackc/pgx/v4"
	_ "github.com/jackc/pgx/v4/stdlib"
	"github.com/ory/dockertest/v3"
	"github.com/ory/dockertest/v3/docker"
)

var (
	host     = "localhost"
	user     = "postgres"
	password = "postgres"
	dbName   = "kstation"
	port     = "5435"
	dsn      = "host=%s port=%s user=%s password=%s dbname=%s sslmode=disable timezone=UTC connect_timeout=5"
)

var resource *dockertest.Resource
var pool *dockertest.Pool
var testDB *sql.DB
var testRepo repository.DatabaseRepo

func TestMain(m *testing.M) {
	p, err := dockertest.NewPool("")
	if err != nil {
		log.Fatalf("could not connect to docker; is it running? %s", err)
	}

	pool = p

	opts := dockertest.RunOptions{
		Repository: "postgres",
		Tag: "14.5",
		Env: []string{
			"POSTGRES_USER=" + user,
			"POSTGRES_PASSWORD=" + password,
			"POSTGRES_DB=" + dbName,
		},
		ExposedPorts: []string{"5432"},
		PortBindings: map[docker.Port][]docker.PortBinding{
			"5432": {
				{HostIP:"0.0.0.0", HostPort: port},
			},
		},
	}

	resource, err = pool.RunWithOptions(&opts)
	if err != nil {
		_ = pool.Purge(resource)
		log.Fatalf("could not start resource: %s", err)
	}

	if err := pool.Retry(func() error {
		var err error
		testDB, err = sql.Open("pgx", fmt.Sprintf(dsn, host, port, user, password, dbName))
		if err != nil {
			log.Println("Error:", err)
			return err
		}
		return testDB.Ping()
	}); err != nil {
		_ = pool.Purge(resource)
		log.Fatalf("could not connect to database: %s", err)
	}

	err = createTables()
	if err != nil {
		log.Fatalf("error creating tables: %s", err)
	}

	testRepo = &PostgresDBRepo{DB: testDB}

	code := m.Run()

	if err := pool.Purge(resource); err != nil {
		log.Fatalf("could not purge resource: %s", err)
	}

	os.Exit(code)
}

func createTables() error {
	tableSQL, err := os.ReadFile("./testdata/create_tables.sql")
	if err != nil {
		fmt.Println(err)
		return err
	}

	_, err = testDB.Exec(string(tableSQL))
	if err != nil {
		fmt.Println(err)
		return err
	}

	return nil
}

func Test_pingDB(t *testing.T) {
	err := testDB.Ping()
	if err != nil {
		t.Error("can't ping database")
	}
}

func TestPostgresDBRepoInsertUser(t *testing.T) {
	testUser := models.User{
		FirstName: "Admin",
		LastName: "User",
		Email: "admin@example.com",
		Password: "secret",
		Image: "test",
		IsAdmin: 1,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	id, err := testRepo.InsertUser(testUser)
	if err != nil {
		t.Errorf("insert user returned an error %s", err)
	}

	if id !=1 {
		t.Errorf("insert user returned wrong id; expected 1, butgot %d", id)
	}

}

func TestPostgresDBRepoGetUserById(t *testing.T) {
	user, err := testRepo.GetUserByID(1)
	if err != nil {
		t.Errorf("error getting user by id: %s", err)
	}

	if user.Email != "admin@example.com" {
		t.Errorf("wrong email returned by GetUser: expected admin@example.com but got %s", user.Email)
	}

	_, err = testRepo.GetUserByID(3)
	if err == nil {
		t.Errorf("no error reported when gettig non existent user by id")
	}
}

func TestPostgresDBRepoGetUserByEmail(t *testing.T) {
	user, err := testRepo.GetUserByEmail("admin@example.com")
	if err != nil {
		t.Errorf("error getting user by id: %s", err)
	}

	if user.ID != 1 {
		t.Errorf("wrong email returned by GetUserByEmail: expected 1 but got %d", user.ID)
	}
}

func TestPostgresDBRepoUpdateUser(t *testing.T) {
	user, _ := testRepo.GetUserByID(1)
	user.FirstName = "Jane"
	user.Email = "jane@smith.com"

	err := testRepo.UpdateUser(*user)
	if err != nil {
		t.Errorf("error updating user %d: %s", 2, err)
	}

	user, _= testRepo.GetUserByID(1)
	if user.FirstName != "Jane" || user.Email != "jane@smith.com" || user.IsAdmin != 1 {
		t.Errorf("expected updated record to have first name Jane and email jane@smith.com, but get %s %s", user.FirstName, user.Email)
	}
}

func TestPostgresDBRepoResetPassword(t *testing.T) {
	err := testRepo.ResetPassword(1, "password")
	if err != nil {
		t.Error("error resetting user's a password", err)
	}

	user, _ := testRepo.GetUserByID(1)
	matches, err := user.PasswordMatches("password")
	if err != nil {
		t.Error(err)
	}

	if !matches {
		t.Errorf("password should match 'password', but does not")
	}
}

func TestPostgresDBRepoInsertLesson(t *testing.T) {
	testLesson := models.Lesson{
		UserId: 1,
		LessonName: "Math",
		TeacherName: "User",
		AvgStar: 7.0,
		AboutAvgStar: int(math.Round(7.0)),
		CommentNumbers: 0,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	id, err := testRepo.InsertLesson(testLesson)
	if err != nil {
		t.Errorf("insert lesson returned an error %s", err)
	}

	if id !=1 {
		t.Errorf("insert lesson returned wrong id; expected 1, but got %d", id)
	}

}

func TestPostgresDBRepoGetLessonById(t *testing.T) {
	lesson, err := testRepo.GetLessonByID(1)
	if err != nil {
		t.Errorf("error getting lesson by id: %s", err)
	}

	if lesson.LessonName != "Math" {
		t.Errorf("wrong lessonName returned by GetLessnoByID: expected Math but got %s", lesson.LessonName)
	}

	_, err = testRepo.GetLessonByID(3)
	if err == nil {
		t.Errorf("no error reported when gettig non existent lesson by id")
	}
}

func TestPostgresDBRepoUpdateLesson(t *testing.T) {
	lesson, _ := testRepo.GetLessonByID(1)
	lesson.AvgStar = 1.5
	lesson.AboutAvgStar = int(math.Round(1.5))
	lesson.CommentNumbers++

	err := testRepo.UpdateLesson(*lesson)
	if err != nil {
		t.Errorf("error updating lesson %d: %s",  1, err)
	}

	lesson, _= testRepo.GetLessonByID(1)
	if lesson.AvgStar != 1.5 || lesson.CommentNumbers != 1 || lesson.LessonName != "Math" {
		t.Errorf("expected updated record to have average star 1.5 and comment numbers 2, but get %f %d", lesson.AvgStar, lesson.CommentNumbers)
	}
}

func TestPostgresDBRepoAllLessons(t *testing.T) {
	lessons, err := testRepo.AllLessons(0)
	if err != nil {
		t.Errorf("0 all lessons reports an error: %s", err)
	}

	if len(lessons) != 1 {
		t.Errorf("0 all lessons reports wrong size; expected 1, but got %d", len(lessons))
	}

	testLesson := models.Lesson{
		UserId: 1,
		LessonName: "English",
		TeacherName: "Smith",
		AvgStar: 3.8,
		AboutAvgStar: int(math.Round(3.8)),
		CommentNumbers: 0,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	testUser := models.User{
		FirstName: "Yamamoto",
		LastName: "Futo",
		Email: "yamamoto@example.com",
		Password: "secret",
		Image: "test",
		IsAdmin: 1,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	_, _ = testRepo.InsertUser(testUser)

	testLesson2 := models.Lesson{
		UserId: 2,
		LessonName: "Science",
		TeacherName: "Yamada",
		AvgStar: 3.2,
		AboutAvgStar: int(math.Round(3.2)),
		CommentNumbers: 0,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	_, _ = testRepo.InsertLesson(testLesson)

	_, _ = testRepo.InsertLesson(testLesson2)

	lessons, err = testRepo.AllLessons(1)
	if err != nil {
		t.Errorf("1 all lessons reports an error: %s", err)
	}

	if len(lessons) != 3 {
		t.Errorf("1 all lessons reports wrong size after insert; expected 2, but got %d", len(lessons))
	}

	if lessons[0].LessonName != "Math" ||  lessons[1].LessonName != "English" || lessons[2].LessonName != "Science" {
		t.Errorf("wrong order 1")
	}

	lessons, err = testRepo.AllLessons(2)
	if err != nil {
		t.Errorf("2 all lessons reports an error: %s", err)
	}

	if len(lessons) != 3 {
		t.Errorf("2 all lessons reports wrong size after insert; expected 2, but got %d", len(lessons))
	}

	if lessons[0].LessonName != "Science" ||  lessons[1].LessonName != "English" || lessons[2].LessonName != "Math" {
		t.Errorf("wrong order 2")
	}

	lessons, err = testRepo.AllLessons(3)
	if err != nil {
		t.Errorf("3 all lessons reports an error: %s", err)
	}

	if len(lessons) != 3 {
		t.Errorf("3 all lessons reports wrong size after insert; expected 2, but got %d", len(lessons))
	}

	if lessons[0].LessonName != "English" ||  lessons[1].LessonName != "Science" || lessons[2].LessonName != "Math" {
		t.Errorf("wrong order 3 %s %s %s", lessons[0].LessonName, lessons[1].LessonName, lessons[2].LessonName)
	}
}

func TestPostgresDBRepoAllLessonsByUser(t *testing.T) {
	lessons, err := testRepo.AllLessonsByUser(1, 0)
	if err != nil {
		t.Errorf("0 all lessons reports an error: %s", err)
	}

	if len(lessons) != 2 {
		t.Errorf("0 all lessons reports wrong size; expected 1, but got %d", len(lessons))
	}

	testLesson3 := models.Lesson{
		UserId: 2,
		LessonName: "English",
		TeacherName: "Smith",
		AvgStar: 3.8,
		AboutAvgStar: int(math.Round(3.8)),
		CommentNumbers: 0,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	testUser := models.User{
		FirstName: "Yamamoto",
		LastName: "Futo",
		Email: "yamamoto@example.com",
		Password: "secret",
		Image: "test",
		IsAdmin: 1,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	_, _ = testRepo.InsertUser(testUser)

	testLesson4 := models.Lesson{
		UserId: 2,
		LessonName: "PE",
		TeacherName: "Yamada",
		AvgStar: 2,
		AboutAvgStar: int(math.Round(2)),
		CommentNumbers: 0,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	testLesson5 := models.Lesson{
		UserId: 5,
		LessonName: "PE",
		TeacherName: "Yamada",
		AvgStar: 2,
		AboutAvgStar: int(math.Round(2)),
		CommentNumbers: 0,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	_, _ = testRepo.InsertLesson(testLesson3)

	_, _ = testRepo.InsertLesson(testLesson4)

	_, err = testRepo.InsertLesson(testLesson5)

	if err == nil {
		t.Error("foreign key not functioned", err)
	}

	lessons, err = testRepo.AllLessonsByUser(2, 1)
	if err != nil {
		t.Errorf("1 all lessons reports an error: %s", err)
	}

	if len(lessons) != 3 {
		t.Errorf("1 all lessons reports wrong size after insert; expected 2, but got %d", len(lessons))
	}

	if lessons[0].LessonName != "Science" ||  lessons[1].LessonName != "English" || lessons[2].LessonName != "PE" {
		t.Errorf("wrong order 1")
	}

	lessons, err = testRepo.AllLessonsByUser(2, 2)
	if err != nil {
		t.Errorf("2 all lessons reports an error: %s", err)
	}

	if len(lessons) != 3 {
		t.Errorf("2 all lessons reports wrong size after insert; expected 2, but got %d", len(lessons))
	}

	if lessons[0].LessonName != "PE" ||  lessons[1].LessonName != "English" || lessons[2].LessonName != "Science" {
		t.Errorf("wrong order 2")
	}

	lessons, err = testRepo.AllLessonsByUser(2, 3)
	if err != nil {
		t.Errorf("3 all lessons reports an error: %s", err)
	}

	if len(lessons) != 3 {
		t.Errorf("3 all lessons reports wrong size after insert; expected 2, but got %d", len(lessons))
	}

	if lessons[0].LessonName != "English" ||  lessons[1].LessonName != "Science" || lessons[2].LessonName != "PE" {
		t.Errorf("wrong order 3 %s %s %s", lessons[0].LessonName, lessons[1].LessonName, lessons[2].LessonName)
	}
}

func TestPostgresDBRepoInsertComment(t *testing.T) {
	testComment := models.Comment{
		LessonId: 1,
		UserId: 1,
		Year: 2023,
		Term: "test",
		Comment: "this is a test",
		TestOrReport: "Report",
		Star: 3,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	id, err := testRepo.InsertComment(testComment)
	if err != nil {
		t.Errorf("insert comment returned an error %s", err)
	}

	if id !=1 {
		t.Errorf("insert comment returned wrong id; expected 1, but got %d", id)
	}

}

func TestPostgresDBRepoGetCommentById(t *testing.T) {
	comment, err := testRepo.GetCommentByID(1)
	if err != nil {
		t.Errorf("error getting user by id: %s", err)
	}

	if comment.Comment != "this is a test" {
		t.Errorf("wrong comment returned by GetUser: expected this is a test but got %s", comment.Comment)
	}

	_, err = testRepo.GetCommentByID(3)
	if err == nil {
		t.Errorf("no error reported when gettig non existent comment by id")
	}
}

func TestPostgresDBRepoAllCommentsByLessonId(t *testing.T) {

	comments, err := testRepo.AllCommentsByLessonId(1)
	if err != nil {
		t.Errorf("all comments reports an error: %s", err)
	}

	if len(comments) != 1 {
		t.Errorf("all comments reports wrong size; expected 1, but got %d", len(comments))
	}

	testComment := models.Comment{
		LessonId: 2,
		UserId: 1,
		Year: 2023,
		Term: "test2",
		Comment: "this is a test",
		TestOrReport: "Test",
		Star: 4,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	_, _ = testRepo.InsertComment(testComment)

	comments, err = testRepo.AllCommentsByLessonId(2)
	if err != nil {
		t.Errorf("all comments reports an error: %s", err)
	}

	if len(comments) != 1 {
		t.Errorf("all comments reports wrong size after insert; expected 2, but got %d", len(comments))
	}
}

func TestPostgresDBRepoAllCommentsByUserId(t *testing.T) {

	comments, err := testRepo.AllCommentsByUserId(1)
	if err != nil {
		t.Errorf("all comments reports an error: %s", err)
	}

	if len(comments) != 2 {
		t.Errorf("all comments reports wrong size; expected 2, but got %d", len(comments))
	}

	testUser := models.User{
		FirstName: "Yamada",
		LastName: "Taro",
		Email: "Futo@example.com",
		Password: "secret",
		Image: "test2",
		IsAdmin: 1,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	_, _ = testRepo.InsertUser(testUser)

	testComment := models.Comment{
		LessonId: 2,
		UserId: 2,
		Year: 2023,
		Term: "test2",
		Comment: "this is a test",
		TestOrReport: "Test",
		Star: 4,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	_, _ = testRepo.InsertComment(testComment)

	comments, err = testRepo.AllCommentsByUserId(2)
	if err != nil {
		t.Errorf("all comments reports an error: %s", err)
	}

	if len(comments) != 1 {
		t.Errorf("all comments reports wrong size after insert; expected 1, but got %d", len(comments))
	}
}

func TestPostgresDBRepoUpdateComment(t *testing.T) {
	comment, _ := testRepo.GetCommentByID(1)
	comment.Year = 2020
	comment.Comment = "Test succeeded"

	err := testRepo.UpdateComment(*comment)
	if err != nil {
		t.Errorf("error updating lesson %d: %s", 2, err)
	}

	comment, _= testRepo.GetCommentByID(1)
	if comment.Year != 2020 || comment.Comment != "Test succeeded" || comment.Star != 3 {
		t.Errorf("expected updated record to have Year 2020 and comment Test is succeeded, but get %d %s", comment.Year, comment.Comment)
	}
}

func TestPostgresDBRepoDeleteUser(t *testing.T) {
	err := testRepo.DeleteComment(2)
	if err != nil{
		t.Errorf("error deleting comment id 2: %s", err)
	}

	_, err = testRepo.GetCommentByID(2)
	if err == nil {
		t.Error("retrieved user id 2, who should have been deleted")
	}
}