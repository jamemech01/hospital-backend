package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"

	"hospital-backend/internal/auth"
	"hospital-backend/internal/handler"
	"hospital-backend/internal/repository"
	"hospital-backend/internal/service"
)

func main() {
	_ = godotenv.Load()

	databaseURL := fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s",
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_HOST"),
		os.Getenv("DB_PORT"),
		os.Getenv("DB_NAME"),
	)

	pool, err := pgxpool.New(context.Background(), databaseURL)
	if err != nil {
		log.Fatal(err)
	}
	defer pool.Close()

	if err := pool.Ping(context.Background()); err != nil {
		log.Fatal("database ping failed:", err)
	}

	jwtSecret := os.Getenv("JWT_SECRET")
	if jwtSecret == "" {
		log.Fatal("JWT_SECRET is required")
	}

	staffRepository := repository.NewStaffRepository(pool)
	authService := service.NewAuthService(staffRepository)
	authHandler := handler.NewAuthHandler(authService, jwtSecret)

	patientRepository := repository.NewPatientRepository(pool)
	patientService := service.NewPatientService(patientRepository)
	patientHandler := handler.NewPatientHandler(patientService)

	router := gin.Default()

	router.GET("/health", func(context *gin.Context) {
		context.JSON(200, gin.H{
			"status": "ok",
		})
	})

	router.POST("/staff/create", authHandler.CreateStaff)
	router.POST("/staff/login", authHandler.Login)
	router.GET("/patient/search", auth.Middleware(jwtSecret), patientHandler.SearchPatients)

	router.Run(":8080")
}
