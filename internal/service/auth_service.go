package service

import (
	"context"
	"errors"

	"golang.org/x/crypto/bcrypt"

	"hospital-backend/internal/repository"
)

type StaffStore interface {
	Create(
		ctx context.Context,
		username string,
		passwordHash string,
		hospitalName string,
	) (int, error)

	FindByUsernameAndHospital(
		ctx context.Context,
		username string,
		hospitalName string,
	) (repository.StaffRecord, error)
}

type AuthService struct {
	staffRepository StaffStore
}

func NewAuthService(staffRepository StaffStore) *AuthService {
	return &AuthService{
		staffRepository: staffRepository,
	}
}

func (service *AuthService) CreateStaff(
	ctx context.Context,
	username string,
	password string,
	hospitalName string,
) (int, error) {
	if username == "" || password == "" || hospitalName == "" {
		return 0, errors.New("username, password and hospital are required")
	}

	passwordHash, err := bcrypt.GenerateFromPassword(
		[]byte(password),
		bcrypt.DefaultCost,
	)
	if err != nil {
		return 0, err
	}

	return service.staffRepository.Create(
		ctx,
		username,
		string(passwordHash),
		hospitalName,
	)
}

func (service *AuthService) Login(
	ctx context.Context,
	username string,
	password string,
	hospitalName string,
) (repository.StaffRecord, error) {
	staff, err := service.staffRepository.FindByUsernameAndHospital(
		ctx,
		username,
		hospitalName,
	)
	if err != nil {
		return repository.StaffRecord{}, errors.New("invalid credentials")
	}

	if err := bcrypt.CompareHashAndPassword(
		[]byte(staff.PasswordHash),
		[]byte(password),
	); err != nil {
		return repository.StaffRecord{}, errors.New("invalid credentials")
	}

	return staff, nil
}
