package service

import (
    "context"
    "errors"

    "hospital-backend/internal/model"
)

type PatientStore interface {
    Search(
        ctx context.Context,
        hospitalID int,
        request model.PatientSearchRequest,
    ) ([]model.Patient, error)
}

type PatientService struct {
    patientRepository PatientStore
}

func NewPatientService(patientRepository PatientStore) *PatientService {
    return &PatientService{
        patientRepository: patientRepository,
    }
}

func (service *PatientService) Search(
    ctx context.Context,
    hospitalID int,
    request model.PatientSearchRequest,
) ([]model.Patient, error) {
    if hospitalID <= 0 {
        return nil, errors.New("hospital_id is required")
    }

    return service.patientRepository.Search(ctx, hospitalID, request)
}