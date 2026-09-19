package model

type Patient struct {
	ID           int    `json:"patient_id"`
	NationalID   string `json:"national_id,omitempty"`
	PassportID   string `json:"passport_id,omitempty"`
	FirstNameTH  string `json:"first_name_th,omitempty"`
	MiddleNameTH string `json:"middle_name_th,omitempty"`
	LastNameTH   string `json:"last_name_th,omitempty"`
	FirstNameEN  string `json:"first_name_en,omitempty"`
	MiddleNameEN string `json:"middle_name_en,omitempty"`
	LastNameEN   string `json:"last_name_en,omitempty"`
	DateOfBirth  string `json:"date_of_birth,omitempty"`
	PatientHN    string `json:"patient_hn"`
	PhoneNumber  string `json:"phone_number,omitempty"`
	Email        string `json:"email,omitempty"`
	Gender       string `json:"gender"`
	HospitalID   int    `json:"hospital_id"`
}

type PatientSearchRequest struct {
	NationalID  string `form:"national_id"`
	PassportID  string `form:"passport_id"`
	FirstName   string `form:"first_name"`
	MiddleName  string `form:"middle_name"`
	LastName    string `form:"last_name"`
	DateOfBirth string `form:"date_of_birth"`
	PhoneNumber string `form:"phone_number"`
	Email       string `form:"email"`
	HospitalName string `form:"hospital_name"`
}
