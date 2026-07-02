package db

import (
	"database/sql"
	"go-chat/internal/domain"
	port "go-chat/internal/port/db"
)

type UserRepository struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) port.UserPort {
	return &UserRepository{db: db}
}

func (r *UserRepository) GetAll() ([]domain.User, error) {
	rows, err := r.db.Query("SELECT id, name, email FROM users")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []domain.User
	for rows.Next() {
		var u domain.User
		rows.Scan(&u.ID, &u.Name, &u.Email)
		users = append(users, u)
	}
	return users, nil
}
