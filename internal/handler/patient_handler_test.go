package handler

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"

	"hospital-backend/internal/model"
	"hospital-backend/internal/service"
)

type fakePatientRepository struct {
	patients []model.Patient
	err      error
}

func (f fakePatientRepository) Search(
	ctx context.Context,
	hospitalID int,
	request model.PatientSearchRequest,
) ([]model.Patient, error) {
	if f.err != nil {
		return nil, f.err
	}
	return f.patients, nil
}

func TestPatientSearchRequiresHospitalID(t *testing.T) {
	gin.SetMode(gin.TestMode)

	handler := NewPatientHandler(service.NewPatientService(fakePatientRepository{patients: []model.Patient{{ID: 1, PatientHN: "HN-001", HospitalID: 7}}}))
	router := gin.New()
	router.GET("/patient/search", handler.SearchPatients)

	request := httptest.NewRequest(http.MethodGet, "/patient/search?first_name=Tom", nil)
	recorder := httptest.NewRecorder()

	router.ServeHTTP(recorder, request)

	if recorder.Code != http.StatusUnauthorized {
		t.Fatalf("expected status 401, got %d", recorder.Code)
	}
}

func TestPatientSearchReturnsPatientsForHospitalScope(t *testing.T) {
	gin.SetMode(gin.TestMode)

	handler := NewPatientHandler(service.NewPatientService(fakePatientRepository{patients: []model.Patient{{ID: 1, PatientHN: "HN-001", HospitalID: 7}}}))
	router := gin.New()
	router.GET("/patient/search", func(c *gin.Context) {
		c.Set("hospital_id", 7)
		c.Next()
	}, handler.SearchPatients)

	request := httptest.NewRequest(http.MethodGet, "/patient/search?first_name=Tom", nil)
	recorder := httptest.NewRecorder()

	router.ServeHTTP(recorder, request)

	if recorder.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d", recorder.Code)
	}

	var response map[string]any
	if err := json.Unmarshal(recorder.Body.Bytes(), &response); err != nil {
		t.Fatalf("invalid json response: %v", err)
	}

	if response["count"] == nil {
		t.Fatal("expected count in response")
	}
}

func TestPatientSearchReturnsInternalServerError(t *testing.T) {
	gin.SetMode(gin.TestMode)
	handler := NewPatientHandler(service.NewPatientService(fakePatientRepository{
		err: errors.New("database unavailable"),
	}))
	router := gin.New()
	router.GET("/patient/search", func(context *gin.Context) {
		context.Set("hospital_id", 1)
		context.Next()
	}, handler.SearchPatients)

	request := httptest.NewRequest(http.MethodGet, "/patient/search?national_id=1100000000001", nil)
	recorder := httptest.NewRecorder()
	router.ServeHTTP(recorder, request)

	if recorder.Code != http.StatusInternalServerError {
		t.Fatalf("expected status 500, got %d", recorder.Code)
	}
}
