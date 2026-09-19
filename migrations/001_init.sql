CREATE TABLE hospital (
    hospital_id INTEGER GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    hospital_name VARCHAR(255) NOT NULL UNIQUE
);

CREATE TABLE staff (
    staff_id INTEGER GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    username VARCHAR(255) NOT NULL,
    password_hash VARCHAR(255) NOT NULL,
    hospital_id INTEGER NOT NULL REFERENCES hospital(hospital_id),
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    UNIQUE (hospital_id, username)
);

CREATE TABLE patient (
    patient_id INTEGER GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    national_id VARCHAR(50),
    passport_id VARCHAR(50),
    first_name_th VARCHAR(255),
    middle_name_th VARCHAR(255),
    last_name_th VARCHAR(255),
    first_name_en VARCHAR(255),
    middle_name_en VARCHAR(255),
    last_name_en VARCHAR(255),
    date_of_birth DATE NOT NULL,
    patient_hn VARCHAR(100) NOT NULL,
    phone_number VARCHAR(50),
    email VARCHAR(255),
    gender VARCHAR(1) NOT NULL CHECK (gender IN ('M', 'F')),
    hospital_id INTEGER NOT NULL REFERENCES hospital(hospital_id),
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    UNIQUE (hospital_id, patient_hn),
    UNIQUE (hospital_id, national_id),
    UNIQUE (hospital_id, passport_id)
);

CREATE INDEX idx_patient_national_id
ON patient (national_id);

CREATE INDEX idx_patient_passport_id
ON patient (passport_id);