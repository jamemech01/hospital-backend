package auth

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"

	"hospital-backend/internal/repository"
)

func TestGenerateTokenRequiresSecret(t *testing.T) {
	_, err := GenerateToken(repository.StaffRecord{ID: 1, HospitalID: 1}, "")
	if err == nil {
		t.Fatal("expected an error when JWT secret is empty")
	}
}

func TestMiddlewareAcceptsValidTokenAndSetsClaims(t *testing.T) {
	gin.SetMode(gin.TestMode)
	secret := "test-secret"
	token, err := GenerateToken(repository.StaffRecord{ID: 7, HospitalID: 3}, secret)
	if err != nil {
		t.Fatal(err)
	}

	router := gin.New()
	router.GET("/protected", Middleware(secret), func(context *gin.Context) {
		staffID, _ := context.Get("staff_id")
		hospitalID, _ := context.Get("hospital_id")
		context.JSON(http.StatusOK, gin.H{"staff_id": staffID, "hospital_id": hospitalID})
	})

	request := httptest.NewRequest(http.MethodGet, "/protected", nil)
	request.Header.Set("Authorization", "Bearer "+token)
	recorder := httptest.NewRecorder()
	router.ServeHTTP(recorder, request)

	if recorder.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d", recorder.Code)
	}
}

func TestMiddlewareRejectsMissingToken(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.New()
	router.GET("/protected", Middleware("test-secret"), func(context *gin.Context) {
		context.Status(http.StatusOK)
	})

	recorder := httptest.NewRecorder()
	router.ServeHTTP(recorder, httptest.NewRequest(http.MethodGet, "/protected", nil))

	if recorder.Code != http.StatusUnauthorized {
		t.Fatalf("expected status 401, got %d", recorder.Code)
	}
}

func TestMiddlewareRejectsExpiredToken(t *testing.T) {
	gin.SetMode(gin.TestMode)
	secret := "test-secret"
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, Claims{
		StaffID:    1,
		HospitalID: 1,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(-time.Hour)),
		},
	})
	signedToken, err := token.SignedString([]byte(secret))
	if err != nil {
		t.Fatal(err)
	}

	router := gin.New()
	router.GET("/protected", Middleware(secret), func(context *gin.Context) {
		context.Status(http.StatusOK)
	})
	recorder := httptest.NewRecorder()
	request := httptest.NewRequest(http.MethodGet, "/protected", nil)
	request.Header.Set("Authorization", "Bearer "+signedToken)
	router.ServeHTTP(recorder, request)

	if recorder.Code != http.StatusUnauthorized {
		t.Fatalf("expected status 401, got %d", recorder.Code)
	}
}
