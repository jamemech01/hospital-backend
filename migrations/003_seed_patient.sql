INSERT INTO patient (
    national_id,
    passport_id,
    first_name_th,
    middle_name_th,
    last_name_th,
    first_name_en,
    middle_name_en,
    last_name_en,
    date_of_birth,
    patient_hn,
    phone_number,
    email,
    gender,
    hospital_id
)
VALUES
    ('1100000000001', NULL, 'สมชาย', NULL, 'ใจดี', 'Somchai', NULL, 'Jaidee', '1985-01-15', 'A-HN-0001', '0810000001', 'somchai.a@example.com', 'M', (SELECT hospital_id FROM hospital WHERE hospital_name = 'Hospital A')),
    ('1100000000002', NULL, 'สุดา', NULL, 'รักดี', 'Suda', NULL, 'Rakdee', '1990-03-22', 'A-HN-0002', '0810000002', 'suda.a@example.com', 'F', (SELECT hospital_id FROM hospital WHERE hospital_name = 'Hospital A')),
    ('1100000000003', NULL, 'วิชัย', NULL, 'มั่นคง', 'Wichai', NULL, 'Mankhong', '1978-07-09', 'A-HN-0003', '0810000003', 'wichai.a@example.com', 'M', (SELECT hospital_id FROM hospital WHERE hospital_name = 'Hospital A')),
    ('1100000000004', NULL, 'พรทิพย์', NULL, 'แสงทอง', 'Porntip', NULL, 'Saengthong', '1995-11-30', 'A-HN-0004', '0810000004', 'porntip.a@example.com', 'F', (SELECT hospital_id FROM hospital WHERE hospital_name = 'Hospital A')),
    ('1100000000005', NULL, 'ณัฐวุฒิ', NULL, 'สุขใจ', 'Nattawut', NULL, 'Sukjai', '1988-06-18', 'A-HN-0005', '0810000005', 'nattawut.a@example.com', 'M', (SELECT hospital_id FROM hospital WHERE hospital_name = 'Hospital A')),
    ('2200000000001', NULL, 'กิตติ', NULL, 'เจริญสุข', 'Kitti', NULL, 'Charoensuk', '1982-02-10', 'B-HN-0001', '0820000001', 'kitti.b@example.com', 'M', (SELECT hospital_id FROM hospital WHERE hospital_name = 'Hospital B')),
    ('2200000000002', NULL, 'มาลี', NULL, 'ใจงาม', 'Malee', NULL, 'Jaingam', '1992-04-25', 'B-HN-0002', '0820000002', 'malee.b@example.com', 'F', (SELECT hospital_id FROM hospital WHERE hospital_name = 'Hospital B')),
    ('2200000000003', NULL, 'ธนา', NULL, 'วัฒนะ', 'Thana', NULL, 'Watthana', '1975-08-14', 'B-HN-0003', '0820000003', 'thana.b@example.com', 'M', (SELECT hospital_id FROM hospital WHERE hospital_name = 'Hospital B')),
    ('2200000000004', NULL, 'อรทัย', NULL, 'บุญมี', 'Orathai', NULL, 'Boonmee', '1987-10-03', 'B-HN-0004', '0820000004', 'orathai.b@example.com', 'F', (SELECT hospital_id FROM hospital WHERE hospital_name = 'Hospital B')),
    ('2200000000005', NULL, 'เอกชัย', NULL, 'รุ่งเรือง', 'Ekachai', NULL, 'Rungrueang', '1998-12-19', 'B-HN-0005', '0820000005', 'ekachai.b@example.com', 'M', (SELECT hospital_id FROM hospital WHERE hospital_name = 'Hospital B'))
ON CONFLICT (hospital_id, patient_hn) DO NOTHING;
