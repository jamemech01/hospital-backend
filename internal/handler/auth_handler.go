package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"hospital-backend/internal/auth"
	"hospital-backend/internal/model"
	"hospital-backend/internal/service"
)

type AuthHandler struct {
	authService *service.AuthService
	jwtSecret   string
}

func NewAuthHandler(authService *service.AuthService, jwtSecret string) *AuthHandler {
	return &AuthHandler{
		authService: authService,
		jwtSecret:   jwtSecret,
	}
}

func (handler *AuthHandler) CreateStaff(context *gin.Context) {
	var request model.CreateStaffRequest

	if err := context.ShouldBindJSON(&request); err != nil {
		context.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid request",
		})
		return
	}

	staffID, err := handler.authService.CreateStaff(
		context.Request.Context(),
		request.Username,
		request.Password,
		request.Hospital,
	)
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	context.JSON(http.StatusCreated, gin.H{
		"staff_id": staffID,
	})
}

func (handler *AuthHandler) Login(context *gin.Context) {
	var request model.LoginRequest

	if err := context.ShouldBindJSON(&request); err != nil {
		context.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid request",
		})
		return
	}

	staff, err := handler.authService.Login(
		context.Request.Context(),
		request.Username,
		request.Password,
		request.Hospital,
	)
	if err != nil {
		context.JSON(http.StatusUnauthorized, gin.H{
			"error": "invalid credentials",
		})
		return
	}

	token, err := auth.GenerateToken(staff, handler.jwtSecret)
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{
			"error": "could not create token",
		})
		return
	}

	context.JSON(http.StatusOK, gin.H{
		"staff_id":    staff.ID,
		"hospital_id": staff.HospitalID,
		"token":       token,
		"message":     "login successful",
	})
}
