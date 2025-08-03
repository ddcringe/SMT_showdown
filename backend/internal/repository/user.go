package repository

import (
	"database/sql"

	"github.com/ddcringe/SMT_showdown/internal/models"
)

type UserRepository struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) *UserRepository {
	return &UserRepository{db: db}
}

func (r *UserRepository) CreateUser(user *models.User) error {
	query := `
		INSERT INTO users (username, email, password_hash, created_at)
		VALUES ($1, $2, $3, $4)
		RETURNING id
	`

	return r.db.QueryRow(
		query,
		user.Username,
		user.Email,
		user.PasswordHash,
		user.CreatedAt,
	).Scan(&user.ID)
}

func (r *UserRepository) UserExists(username, email string) (bool, error) {
	query := `
		SELECT EXISTS(
			SELECT 1 FROM users 
			WHERE username = $1 OR email = $2
		)
	`

	var exists bool
	err := r.db.QueryRow(query, username, email).Scan(&exists)
	return exists, err
}

func (r *UserRepository) GetUserByEmail(email string) (*models.User, error) {
	query := `
		SELECT id, username, email, password_hash, created_at, last_login
		FROM users 
		WHERE email = $1
	`

	user := &models.User{}
	err := r.db.QueryRow(query, email).Scan(
		&user.ID,
		&user.Username,
		&user.Email,
		&user.PasswordHash,
		&user.CreatedAt,
		&user.LastLogin,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil // Пользователь не найден - это не ошибка
		}
		return nil, err
	}

	return user, nil
}
