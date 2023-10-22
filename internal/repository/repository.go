package repository

import (
	"database/sql"
	// "kstation_backend/internal/models"
)

type DatabaseRepo interface {
	Connection() *sql.DB
}