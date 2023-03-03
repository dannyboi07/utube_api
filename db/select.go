package db

import (
	"errors"
	"utube/models"

	"github.com/jackc/pgx/v5"
)

func SelectUserByEmail(email string) (models.Actor, bool, error) {
	var user models.Actor
	row := db.QueryRow(dbContext, "SELECT id, email, password, name, created_at, updated_at FROM actor WHERE email = $1", email)
	err := row.Scan(&user.Id, &user.Email, &user.Password, &user.Name, &user.CreatedAt, &user.UpdatedAt)
	if errors.Is(err, pgx.ErrNoRows) {
		return user, false, nil
	}

	return user, true, err
}
