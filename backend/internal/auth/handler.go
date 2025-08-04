package auth

import (
	"context"
	"database/sql"
	"errors"
	"strconv"
	"strings"
	"time"

	"github.com/ddcringe/SMT_showdown/internal/models"
	"github.com/jmoiron/sqlx"
)

type UserRepository struct {
	db *sqlx.DB
}

func NewUserRepository(db *sqlx.DB) *UserRepository {
	return &UserRepository{db: db}
}

// CreateUser создает нового пользователя
func (r *UserRepository) CreateUser(ctx context.Context, user *models.User) error {
	query := `
		INSERT INTO users (
			username, 
			email, 
			password_hash, 
			created_at,
			bio,
			avatar_url
		) VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING id
	`

	err := r.db.QueryRowContext(
		ctx,
		query,
		user.Username,
		user.Email,
		user.PasswordHash,
		time.Now(),
		user.Bio,
		user.AvatarURL,
	).Scan(&user.ID)

	if err != nil {
		return err
	}

	return nil
}

// GetUserByID возвращает пользователя по ID
func (r *UserRepository) GetUserByID(ctx context.Context, id int) (*models.User, error) {
	query := `
		SELECT 
			id,
			username,
			email,
			password_hash,
			created_at,
			last_login,
			bio,
			avatar_url
		FROM users
		WHERE id = $1
	`

	var user models.User
	err := r.db.GetContext(ctx, &user, query, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrNotFound
		}
		return nil, err
	}

	return &user, nil
}

// GetUserByEmail возвращает пользователя по email
func (r *UserRepository) GetUserByEmail(ctx context.Context, email string) (*models.User, error) {
	query := `
		SELECT 
			id,
			username,
			email,
			password_hash,
			created_at,
			last_login,
			bio,
			avatar_url
		FROM users
		WHERE email = $1
	`

	var user models.User
	err := r.db.GetContext(ctx, &user, query, email)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrNotFound
		}
		return nil, err
	}

	return &user, nil
}

// UpdateUser обновляет данные пользователя
func (r *UserRepository) UpdateUser(ctx context.Context, user *models.User) error {
	query := `
		UPDATE users SET
			username = $1,
			email = $2,
			bio = $3,
			avatar_url = $4,
			last_login = $5
		WHERE id = $6
	`

	result, err := r.db.ExecContext(
		ctx,
		query,
		user.Username,
		user.Email,
		user.Bio,
		user.AvatarURL,
		user.LastLogin,
		user.ID,
	)

	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return ErrNotFound
	}

	return nil
}

// UpdatePassword обновляет пароль пользователя
func (r *UserRepository) UpdatePassword(ctx context.Context, userID int, newHash string) error {
	query := `
		UPDATE users SET
			password_hash = $1
		WHERE id = $2
	`

	result, err := r.db.ExecContext(ctx, query, newHash, userID)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return ErrNotFound
	}

	return nil
}

// UserExists проверяет существование пользователя
func (r *UserRepository) UserExists(ctx context.Context, username, email string) (bool, error) {
	query := `
		SELECT EXISTS(
			SELECT 1 FROM users 
			WHERE username = $1 OR email = $2
		)
	`

	var exists bool
	err := r.db.QueryRowContext(ctx, query, username, email).Scan(&exists)
	if err != nil {
		return false, err
	}

	return exists, nil
}

// GetUserProfile возвращает профиль пользователя
func (r *UserRepository) GetUserProfile(ctx context.Context, userID int) (*models.UserProfile, error) {
	query := `
		SELECT 
			id,
			username,
			email,
			created_at,
			bio,
			avatar_url
		FROM users
		WHERE id = $1
	`

	var profile models.UserProfile
	err := r.db.GetContext(ctx, &profile, query, userID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrNotFound
		}
		return nil, err
	}

	return &profile, nil
}

// UpdateUserProfile обновляет профиль пользователя
func (r *UserRepository) UpdateUserProfile(ctx context.Context, userID int, updates map[string]interface{}) error {
	if len(updates) == 0 {
		return nil
	}

	query := `UPDATE users SET `
	params := []interface{}{}
	paramCounter := 1

	for field, value := range updates {
		query += field + " = $" + strconv.Itoa(paramCounter) + ", "
		params = append(params, value)
		paramCounter++
	}

	query = strings.TrimSuffix(query, ", ")
	query += " WHERE id = $" + strconv.Itoa(paramCounter)
	params = append(params, userID)

	result, err := r.db.ExecContext(ctx, query, params...)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return ErrNotFound
	}

	return nil
}

var ErrNotFound = errors.New("user not found")
