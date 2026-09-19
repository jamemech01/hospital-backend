package service

import (
	"context"
	"errors"
	"testing"

	"hospital-backend/internal/model"
)

type fakePatientStore struct {
	called     bool
	hospitalID int
	request    model.PatientSearchRequest
	patients   []model.Patient
	err        error
}

func (store *fakePatientStore) Search(
	ctx context.Context,
	hospitalID int,
	request model.PatientSearchRequest,
) ([]model.Patient, error) {
	store.called = true
	store.hospitalID = hospitalID
	store.request = request
	return store.patients, store.err
}

func TestPatientServiceRejectsInvalidHospitalID(t *testing.T) {
	store := &fakePatientStore{}
	patientService := NewPatientService(store)

	patients, err := patientService.Search(context.Background(), 0, model.PatientSearchRequest{})
	if err == nil {
		t.Fatal("expected an error for invalid hospital ID")
	}
	if patients != nil {
		t.Fatalf("expected nil patients, got %#v", patients)
	}
	if store.called {
		t.Fatal("repository should not be called")
	}
}

func TestPatientServiceDelegatesSearch(t *testing.T) {
	store := &fakePatientStore{
		patients: []model.Patient{{ID: 1, HospitalID: 3}},
	}
	patientService := NewPatientService(store)
	request := model.PatientSearchRequest{NationalID: "1100000000001"}

	patients, err := patientService.Search(context.Background(), 3, request)
	if err != nil {
		t.Fatal(err)
	}
	if !store.called || store.hospitalID != 3 || store.request != request {
		t.Fatal("service did not delegate the expected search arguments")
	}
	if len(patients) != 1 || patients[0].ID != 1 {
		t.Fatalf("unexpected patients: %#v", patients)
	}
}

func TestPatientServiceReturnsRepositoryError(t *testing.T) {
	expectedErr := errors.New("database unavailable")
	store := &fakePatientStore{err: expectedErr}
	patientService := NewPatientService(store)

	_, err := patientService.Search(context.Background(), 1, model.PatientSearchRequest{})
	if !errors.Is(err, expectedErr) {
		t.Fatalf("expected repository error, got %v", err)
	}
}
