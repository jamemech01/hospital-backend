INSERT INTO hospital (hospital_name)
VALUES
	('Hospital A'),
	('Hospital B')
ON CONFLICT (hospital_name) DO NOTHING;