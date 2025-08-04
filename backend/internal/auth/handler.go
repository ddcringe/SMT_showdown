package auth

import (
	"net/http"
	"time"

	"github.com/ddcringe/SMT_showdown/internal/models"
	"github.com/ddcringe/SMT_showdown/internal/repository"
	"github.com/ddcringe/SMT_showdown/pkg/jwt"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

type Handler struct {
	userRepo *repository.UserRepository
}

func NewHandler(userRepo *repository.UserRepository) *Handler {
	return &Handler{userRepo: userRepo}
}

// ... (остальные структуры RegisterRequest, LoginRequest, AuthResponse остаются такими же)
type RegisterRequest struct {
	Username string `json:"username" binding:"required,min=3,max=30"`
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=8"`
}

// LoginRequest структура запроса входа
type LoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

// AuthResponse структура ответа с токеном
type AuthResponse struct {
	Token     string `json:"token"`
	ExpiresAt int64  `json:"expires_at"` // Unix timestamp
}

func (h *Handler) Register(c *gin.Context) {
	var req RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	exists, err := h.userRepo.UserExists(req.Username, req.Email)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "database error"})
		return
	}
	if exists {
		c.JSON(http.StatusConflict, gin.H{"error": "user already exists"})
		return
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "could not hash password"})
		return
	}

	user := models.User{
		Username:     req.Username,
		Email:        req.Email,
		PasswordHash: string(hashedPassword),
		CreatedAt:    time.Now(),
	}

	if err := h.userRepo.CreateUser(&user); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "could not create user"})
		return
	}

	token, err := jwt.GenerateToken(user.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "could not generate token"})
		return
	}

	c.JSON(http.StatusCreated, AuthResponse{
		Token:     token,
		ExpiresAt: time.Now().Add(24 * time.Hour).Unix(),
	})
}

func (h *Handler) Login(c *gin.Context) {
	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user, err := h.userRepo.GetUserByEmail(req.Email)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "database error"})
		return
	}
	if user == nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid credentials"})
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.Password)); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid credentials"})
		return
	}

	token, err := jwt.GenerateToken(user.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "could not generate token"})
		return
	}

	c.JSON(http.StatusOK, AuthResponse{
		Token:     token,
		ExpiresAt: time.Now().Add(24 * time.Hour).Unix(),
	})
}
