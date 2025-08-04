package auth

import (
	"context"
	"errors"
	"time"

	"github.com/ddcringe/SMT_showdown/internal/models"
	"github.com/ddcringe/SMT_showdown/internal/repository"
	"golang.org/x/crypto/bcrypt"
)

var (
	ErrUserExists         = errors.New("user already exists")
	ErrInvalidCredentials = errors.New("invalid credentials")
)

type AuthService struct {
	userRepo *repository.UserRepository
}

func NewAuthService(userRepo *repository.UserRepository) *AuthService {
	return &AuthService{userRepo: userRepo}
}

func (s *AuthService) Register(ctx context.Context, email, username, password string) (int, error) {
	// Проверка существования пользователя
	exists, err := s.userRepo.UserExists(ctx, username, email)
	if err != nil {
		return 0, err
	}
	if exists {
		return 0, ErrUserExists
	}

	// Хеширование пароля
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return 0, err
	}

	// Создание пользователя
	user := &models.User{
		Username:     username,
		Email:        email,
		PasswordHash: string(hashedPassword),
		CreatedAt:    time.Now(),
	}

	if err := s.userRepo.CreateUser(ctx, user); err != nil {
		return 0, err
	}

	return user.ID, nil
}

func (s *AuthService) Login(ctx context.Context, email, password string) (int, error) {
	user, err := s.userRepo.GetUserByEmail(ctx, email)
	if err != nil {
		return 0, err
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password)); err != nil {
		return 0, ErrInvalidCredentials
	}

	return user.ID, nil
}
