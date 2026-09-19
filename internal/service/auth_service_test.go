package service

import (
	"context"
	"errors"
	"testing"

	"golang.org/x/crypto/bcrypt"

	"hospital-backend/internal/repository"
)

func TestCreateStaffRejectsMissingFields(t *testing.T) {
	testCases := []struct {
		name         string
		username     string
		password     string
		hospitalName string
	}{
		{
			name:         "missing username",
			password:     "admin1234",
			hospitalName: "Hospital A",
		},
		{
			name:         "missing password",
			username:     "admin1",
			hospitalName: "Hospital A",
		},
		{
			name:     "missing hospital",
			username: "admin1",
			password: "admin1234",
		},
	}

	authService := &AuthService{}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			_, err := authService.CreateStaff(
				context.Background(),
				testCase.username,
				testCase.password,
				testCase.hospitalName,
			)

			if err == nil {
				t.Fatal("expected error, got nil")
			}
		})
	}
}

type fakeStaffStore struct {
	createdPasswordHash string
	staff               repository.StaffRecord
	findError           error
}

func (fake *fakeStaffStore) Create(
	ctx context.Context,
	username string,
	passwordHash string,
	hospitalName string,
) (int, error) {
	fake.createdPasswordHash = passwordHash
	return 10, nil
}

func (fake *fakeStaffStore) FindByUsernameAndHospital(
	ctx context.Context,
	username string,
	hospitalName string,
) (repository.StaffRecord, error) {
	return fake.staff, fake.findError
}

func TestCreateStaffHashesPassword(t *testing.T) {
	fakeStore := &fakeStaffStore{}
	authService := NewAuthService(fakeStore)

	staffID, err := authService.CreateStaff(
		context.Background(),
		"admin1",
		"admin1234",
		"Hospital A",
	)
	if err != nil {
		t.Fatal(err)
	}

	if staffID != 10 {
		t.Fatalf("expected staff ID 10, got %d", staffID)
	}

	if err := bcrypt.CompareHashAndPassword(
		[]byte(fakeStore.createdPasswordHash),
		[]byte("admin1234"),
	); err != nil {
		t.Fatal("password was not hashed correctly")
	}
}

func TestLogin(t *testing.T) {
	passwordHash, err := bcrypt.GenerateFromPassword(
		[]byte("admin1234"),
		bcrypt.DefaultCost,
	)
	if err != nil {
		t.Fatal(err)
	}

	testCases := []struct {
		name        string
		password    string
		findError   error
		expectError bool
	}{
		{
			name:     "valid password",
			password: "admin1234",
		},
		{
			name:        "invalid password",
			password:    "wrong-password",
			expectError: true,
		},
		{
			name:        "staff not found",
			password:    "admin1234",
			findError:   errors.New("staff not found"),
			expectError: true,
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			fakeStore := &fakeStaffStore{
				staff: repository.StaffRecord{
					ID:           1,
					PasswordHash: string(passwordHash),
					HospitalID:   1,
				},
				findError: testCase.findError,
			}
			authService := NewAuthService(fakeStore)

			_, err := authService.Login(
				context.Background(),
				"admin1",
				testCase.password,
				"Hospital A",
			)

			if testCase.expectError && err == nil {
				t.Fatal("expected error, got nil")
			}

			if !testCase.expectError && err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
		})
	}
}
