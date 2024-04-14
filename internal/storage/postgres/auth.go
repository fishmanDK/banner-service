package postgres

import (
	"fmt"
	"github.com/fishmanDK/avito_test_task/models"
)

func (p *Postgres) CreateUser(newUser models.NewUser) error {
	const op = "postgres.CreateUser"

	query := "INSERT INTO users (email, hash_password, role) VALUES ($1, $2, $3)"
	_, err := p.db.Query(query, newUser.Email, newUser.Password, newUser.Role)
	if err != nil {
		return fmt.Errorf("%s:%d", op, err)
	}

	return nil
}

type UserRole struct {
	Role string `db:"role"`
}

func (p *Postgres) GetUserRole(user models.User) (string, error) {
	const op = "postgres.GetUserRole"
	query := "SELECT users.role FROM users WHERE email = $1 AND hash_password = $2"

	var userRole UserRole
	err := p.db.Get(&userRole, query, user.Email, user.Password)
	if err != nil {
		return "", fmt.Errorf("%s: %w", op, err)
	}

	return userRole.Role, nil
}
