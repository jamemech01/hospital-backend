package handler

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"

	"hospital-backend/internal/repository"
	"hospital-backend/internal/service"
)

type authHandlerStore struct {
	staff repository.StaffRecord
}

func (store *authHandlerStore) Create(
	ctx context.Context,
	username string,
	passwordHash string,
	hospitalName string,
) (int, error) {
	return 12, nil
}

func (store *authHandlerStore) FindByUsernameAndHospital(
	ctx context.Context,
	username string,
	hospitalName string,
) (repository.StaffRecord, error) {
	return store.staff, nil
}

func TestCreateStaffReturnsBadRequestForInvalidJSON(t *testing.T) {
	gin.SetMode(gin.TestMode)

	router := gin.New()
	router.POST("/staff/create", (&AuthHandler{}).CreateStaff)

	request := httptest.NewRequest(
		http.MethodPost,
		"/staff/create",
		strings.NewReader(`{"username":`),
	)
	request.Header.Set("Content-Type", "application/json")

	recorder := httptest.NewRecorder()
	router.ServeHTTP(recorder, request)

	if recorder.Code != http.StatusBadRequest {
		t.Fatalf("expected status 400, got %d", recorder.Code)
	}
}

func TestCreateStaffReturnsCreatedID(t *testing.T) {
	gin.SetMode(gin.TestMode)

	handler := NewAuthHandler(
		service.NewAuthService(&authHandlerStore{}),
		"test-secret",
	)

	router := gin.New()
	router.POST("/staff/create", handler.CreateStaff)

	request := httptest.NewRequest(
		http.MethodPost,
		"/staff/create",
		strings.NewReader(`{"username":"admin","password":"secret1234","hospital":"Hospital A"}`),
	)
	request.Header.Set("Content-Type", "application/json")

	recorder := httptest.NewRecorder()
	router.ServeHTTP(recorder, request)

	if recorder.Code != http.StatusCreated {
		t.Fatalf("expected status 201, got %d", recorder.Code)
	}

	if !strings.Contains(recorder.Body.String(), `"staff_id":12`) {
		t.Fatalf("expected staff ID in response, got %s", recorder.Body.String())
	}
}

func TestLoginReturnsJWT(t *testing.T) {
	gin.SetMode(gin.TestMode)

	passwordHash, err := bcrypt.GenerateFromPassword(
		[]byte("secret1234"),
		bcrypt.MinCost,
	)
	if err != nil {
		t.Fatal(err)
	}

	store := &authHandlerStore{
		staff: repository.StaffRecord{
			ID:           4,
			PasswordHash: string(passwordHash),
			HospitalID:   3,
		},
	}

	handler := NewAuthHandler(
		service.NewAuthService(store),
		"test-secret",
	)

	router := gin.New()
	router.POST("/staff/login", handler.Login)

	request := httptest.NewRequest(
		http.MethodPost,
		"/staff/login",
		strings.NewReader(`{"username":"admin","password":"secret1234","hospital":"Hospital B"}`),
	)
	request.Header.Set("Content-Type", "application/json")

	recorder := httptest.NewRecorder()
	router.ServeHTTP(recorder, request)

	if recorder.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d", recorder.Code)
	}

	var response struct {
		Token      string `json:"token"`
		HospitalID int    `json:"hospital_id"`
	}

	if err := json.Unmarshal(recorder.Body.Bytes(), &response); err != nil {
		t.Fatal(err)
	}

	if response.Token == "" {
		t.Fatal("expected JWT token")
	}

	if response.HospitalID != 3 {
		t.Fatalf("expected hospital ID 3, got %d", response.HospitalID)
	}
}

func TestLoginRejectsInvalidCredentials(t *testing.T) {
	gin.SetMode(gin.TestMode)

	passwordHash, err := bcrypt.GenerateFromPassword(
		[]byte("secret1234"),
		bcrypt.MinCost,
	)
	if err != nil {
		t.Fatal(err)
	}

	store := &authHandlerStore{
		staff: repository.StaffRecord{
			ID:           1,
			PasswordHash: string(passwordHash),
			HospitalID:   1,
		},
	}

	handler := NewAuthHandler(
		service.NewAuthService(store),
		"test-secret",
	)

	router := gin.New()
	router.POST("/staff/login", handler.Login)

	request := httptest.NewRequest(
		http.MethodPost,
		"/staff/login",
		strings.NewReader(`{"username":"admin","password":"wrong-password","hospital":"Hospital A"}`),
	)
	request.Header.Set("Content-Type", "application/json")

	recorder := httptest.NewRecorder()
	router.ServeHTTP(recorder, request)

	if recorder.Code != http.StatusUnauthorized {
		t.Fatalf("expected status 401, got %d", recorder.Code)
	}
}
