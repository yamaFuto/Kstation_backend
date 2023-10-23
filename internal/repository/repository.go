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
}