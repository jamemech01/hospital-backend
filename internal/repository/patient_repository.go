package repository

import (
	"context"
	"fmt"
	"strings"

	"github.com/jackc/pgx/v5/pgxpool"

	"hospital-backend/internal/model"
)

type PatientRepository struct {
	pool *pgxpool.Pool
}

func NewPatientRepository(pool *pgxpool.Pool) *PatientRepository {
	return &PatientRepository{pool: pool}
}

func (repository *PatientRepository) Search(
	ctx context.Context,
	hospitalID int,
	request model.PatientSearchRequest,
) ([]model.Patient, error) {
	conditions := []string{"p.hospital_id = $1"}
	args := []any{hospitalID}

	addExactCondition := func(column string, value string) {
		if value == "" {
			return
		}

		placeholder := len(args) + 1
		conditions = append(conditions, fmt.Sprintf("%s = $%d", column, placeholder))
		args = append(args, value)
	}

	addNameCondition := func(thaiColumn string, englishColumn string, value string) {
		if value == "" {
			return
		}

		placeholder := len(args) + 1
		conditions = append(conditions, fmt.Sprintf(
			"(%s ILIKE '%%' || $%d || '%%' OR %s ILIKE '%%' || $%d || '%%')",
			thaiColumn,
			placeholder,
			englishColumn,
			placeholder,
		))
		args = append(args, value)
	}

	addExactCondition("p.national_id", request.NationalID)
	addExactCondition("p.passport_id", request.PassportID)
	addNameCondition("p.first_name_th", "p.first_name_en", request.FirstName)
	addNameCondition("p.middle_name_th", "p.middle_name_en", request.MiddleName)
	addNameCondition("p.last_name_th", "p.last_name_en", request.LastName)
	addExactCondition("p.phone_number", request.PhoneNumber)
	addExactCondition("p.email", request.Email)

	query := `
		SELECT p.patient_id,
		       COALESCE(p.national_id, ''), COALESCE(p.passport_id, ''),
		       COALESCE(p.first_name_th, ''), COALESCE(p.middle_name_th, ''), COALESCE(p.last_name_th, ''),
		       COALESCE(p.first_name_en, ''), COALESCE(p.middle_name_en, ''), COALESCE(p.last_name_en, ''),
		       TO_CHAR(p.date_of_birth, 'YYYY-MM-DD'), p.patient_hn,
		       COALESCE(p.phone_number, ''), COALESCE(p.email, ''),
		       p.gender, p.hospital_id
		FROM patient p
		WHERE ` + strings.Join(conditions, " AND ") + `
		ORDER BY p.patient_id`

	rows, err := repository.pool.Query(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	patients := make([]model.Patient, 0)
	for rows.Next() {
		var patient model.Patient
		if err := rows.Scan(
			&patient.ID,
			&patient.NationalID,
			&patient.PassportID,
			&patient.FirstNameTH,
			&patient.MiddleNameTH,
			&patient.LastNameTH,
			&patient.FirstNameEN,
			&patient.MiddleNameEN,
			&patient.LastNameEN,
			&patient.DateOfBirth,
			&patient.PatientHN,
			&patient.PhoneNumber,
			&patient.Email,
			&patient.Gender,
			&patient.HospitalID,
		); err != nil {
			return nil, err
		}

		patients = append(patients, patient)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return patients, nil
}
