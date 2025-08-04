package repository

import (
	"database/sql"
	"strconv"
	"strings"

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
func (r *UserRepository) GetProfile(userID int) (*models.UserProfile, error) {
	query := `
        SELECT id, username, email, created_at, last_login, bio, avatar_url 
        FROM users WHERE id = $1
    `

	profile := &models.UserProfile{}
	err := r.db.QueryRow(query, userID).Scan(
		&profile.ID,
		&profile.Username,
		&profile.Email,
		&profile.CreatedAt,
		&profile.LastLogin,
		&profile.Bio,
		&profile.AvatarURL,
	)

	if err != nil {
		return nil, err
	}
	return profile, nil
}

func (r *UserRepository) UpdateProfile(userID int, updates map[string]interface{}) error {
	query := `UPDATE users SET `
	params := []interface{}{}
	i := 1

	for field, value := range updates {
		query += field + " = $" + strconv.Itoa(i) + ", "
		params = append(params, value)
		i++
	}

	query = strings.TrimSuffix(query, ", ") + " WHERE id = $" + strconv.Itoa(i)
	params = append(params, userID)

	_, err := r.db.Exec(query, params...)
	return err
}
