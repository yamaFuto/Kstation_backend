package dbrepo

import (
	"database/sql"
	"fmt"
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

func TestPostgresDBRepoGetUser(t *testing.T) {
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
	if user.FirstName != "Jane" || user.Email != "jane@smith.com" {
		t.Errorf("expected updated record to have first name Jane and email jane@smith.com, but get %s %s", user.FirstName, user.Email)
	}
}

func TestPostgresDBRepoInsertLesson(t *testing.T) {
	testLesson := models.Lesson{
		LessonName: "Math",
		TeacherName: "User",
		AvgStar: 0.0,
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

func TestPostgresDBRepoAllLessons(t *testing.T) {
	users, err := testRepo.AllLessons()
	if err != nil {
		t.Errorf("all lessons reports an error: %s", err)
	}

	if len(users) != 1 {
		t.Errorf("all lessons reports wrong size; expected 1, but got %d", len(users))
	}

	testLesson := models.Lesson{
		LessonName: "English",
		TeacherName: "Smith",
		AvgStar: 0,
		CommentNumbers: 0,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	_, _ = testRepo.InsertLesson(testLesson)

	users, err = testRepo.AllLessons()
	if err != nil {
		t.Errorf("all lessons reports an error: %s", err)
	}

	if len(users) != 2 {
		t.Errorf("all lessons reports wrong size after insert; expected 2, but got %d", len(users))
	}
}