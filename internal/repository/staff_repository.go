package repository

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
)

type StaffRepository struct {
	pool *pgxpool.Pool
}

type StaffRecord struct {
	ID           int
	PasswordHash string
	HospitalID   int
}

func NewStaffRepository(pool *pgxpool.Pool) *StaffRepository {
	return &StaffRepository{pool: pool}
}

func (repository *StaffRepository) Create(
	ctx context.Context,
	username string,
	passwordHash string,
	hospitalName string,
) (int, error) {
	var staffID int

	err := repository.pool.QueryRow(
		ctx,
		`INSERT INTO staff (username, password_hash, hospital_id)
         SELECT $1, $2, hospital_id
         FROM hospital
         WHERE hospital_name = $3
         RETURNING staff_id`,
		username,
		passwordHash,
		hospitalName,
	).Scan(&staffID)

	return staffID, err
}

func (repository *StaffRepository) FindByUsernameAndHospital(
	ctx context.Context,
	username string,
	hospitalName string,
) (StaffRecord, error) {
	var staff StaffRecord

	err := repository.pool.QueryRow(
		ctx,
		`SELECT s.staff_id, s.password_hash, s.hospital_id
         FROM staff s
         JOIN hospital h ON h.hospital_id = s.hospital_id
         WHERE s.username = $1
           AND h.hospital_name = $2`,
		username,
		hospitalName,
	).Scan(&staff.ID, &staff.PasswordHash, &staff.HospitalID)

	return staff, err
}
