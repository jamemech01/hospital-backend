package model

type Staff struct {
	ID           int    `json:"staff_id"`
	Username     string `json:"username"`
	PasswordHash string `json:"-"`
	HospitalID   int    `json:"hospital_id"`
}

type CreateStaffRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required,min=8"`
	Hospital string `json:"hospital" binding:"required"`
}

type LoginRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
	Hospital string `json:"hospital" binding:"required"`
}
