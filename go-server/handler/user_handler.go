package handler

import (
	jwt "go-server/jwt"
	"go-server/repository"
	"net/http"

	"github.com/gin-gonic/gin"
)

type UserHandler struct {
	repo *repository.UserRepo
}

// Разделили структуры запросов, так как для регистрации нужно имя
type RegisterRequest struct {
	Userlogin string `json:"userlogin" binding:"required"`
	Username  string `json:"username" binding:"required"`
	Password  string `json:"password" binding:"required"`
}

type LoginRequest struct {
	Userlogin string `json:"userlogin" binding:"required"`
	Password  string `json:"password" binding:"required"`
}

func NewUserHandler(repo *repository.UserRepo) (*UserHandler, error) {
	return &UserHandler{repo: repo}, nil
}

func (h *UserHandler) CreateUser(c *gin.Context) {
	var req RegisterRequest // Используем правильную структуру с Username
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid json input"})
		return
	}

	if err := h.repo.Create(req.Userlogin, req.Username, req.Password); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create user"})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"status": "created"})
}

func (h *UserHandler) GetUser(c *gin.Context) {
	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid json input"})
		return
	}

	user, err := h.repo.GetUserByCredentials(req.Userlogin, req.Password)
	if err != nil {
		// Ошибка 401 вместо 400, если логин/пароль не подошли
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid login or password"})
		return
	}

	token, err := jwt.GenerateToken(int(user.ID), user.Username)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to generate token"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"token": token,
	})
}

func WsUserValid() {

}
