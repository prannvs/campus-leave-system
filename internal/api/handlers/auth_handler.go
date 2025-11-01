package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/prannvs/campus-leave-system/internal/auth"
	"github.com/prannvs/campus-leave-system/internal/core"
	"github.com/prannvs/campus-leave-system/internal/models"
	"github.com/prannvs/campus-leave-system/internal/services"
)

type AuthHandler struct {
	userService *services.UserService
	jwtService  *auth.JWTService
}

func NewAuthHandler(userService *services.UserService, jwtService *auth.JWTService) *AuthHandler {
	return &AuthHandler{
		userService: userService,
		jwtService:  jwtService,
	}
}

func (h *AuthHandler) Register(c *gin.Context) {
	var req models.RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		core.ErrorResponse(c, http.StatusBadRequest, err, nil)
		return
	}

	user, err := h.userService.Register(req)
	if err != nil {
		core.ErrorResponse(c, http.StatusBadRequest, err, nil)
		return
	}

	token, err := h.jwtService.GenerateToken(user)
	if err != nil {
		core.ErrorResponse(c, http.StatusInternalServerError, err, nil)
		return
	}

	core.SuccessResponse(c, http.StatusCreated, "User registered successfully", gin.H{
		"user":  user,
		"token": token,
	})
}

func (h *AuthHandler) Login(c *gin.Context) {
	var req models.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		core.ErrorResponse(c, http.StatusBadRequest, err, nil)
		return
	}

	user, err := h.userService.Login(req.Email, req.Password)
	if err != nil {
		core.ErrorResponse(c, http.StatusUnauthorized, err, nil)
		return
	}

	token, err := h.jwtService.GenerateToken(user)
	if err != nil {
		core.ErrorResponse(c, http.StatusInternalServerError, err, nil)
		return
	}

	core.SuccessResponse(c, http.StatusOK, "Login successful", gin.H{
		"user":  user,
		"token": token,
	})
}
