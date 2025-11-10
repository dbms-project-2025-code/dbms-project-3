-- name: InsertDepartments :exec
INSERT INTO Department (department_id, name, hod_name) VALUES
(1, 'Computer Eng', 'Dr. Anita Rao'),
(2, 'Mechanical', 'Dr. S. Mehra'),
(3, 'Civil', 'Dr. R. Yadav');

-- name: InsertFaculty :exec
INSERT INTO Faculty (faculty_id, name, department_id) VALUES
(101, 'Prof. D. Singh', 1),
(102, 'Prof. A. John', 2),
(103, 'Prof. R. Shah', 3);

-- name: InsertAlumni :exec
INSERT INTO Alumni (roll_no, name, email_id, phone_no, batch, department_id) VALUES
('2024UIN3356', 'Karan Patel', 'karan.p@domain.com', '9123456780', 2024, 1),
('2024UIN3357', 'Sana Verma', 'sana.v@domain.com', '9123456781', 2024, 2),
('2024UIN3358', 'Piyush Das', 'piyush.d@domain.com', '9123456782', 2024, 3);

-- name: InsertAcademicHistory :exec
INSERT INTO Academic_History (roll_no, cg1, cg2, cg3, cg4, cg5, cg6, cg7, cg8) VALUES
('2024UIN3356', 8.5, 8.3, 8.7, 8.8, 8.9, 9.0, 9.1, 9.2),
('2024UIN3357', 7.9, 8.1, 8.0, 7.8, 8.2, 8.0, 8.3, 8.2),
('2024UIN3358', 8.0, 8.2, 8.1, 8.4, 8.3, 8.2, 8.5, 8.7);

-- name: InsertEventData :exec
INSERT INTO Event_Data (event_id, event_name, description, event_datetime, location) VALUES
(1, 'Alumni Meet', 'Annual alumni gathering', '2025-01-15 18:00:00', 'Seminar Hall'),
(2, 'Tech Talk', 'AI trends lecture', '2025-03-10 14:30:00', 'Main Hall'),
(3, 'Sports Fest', 'Inter-dept tournament', '2025-02-20 09:00:00', 'Sports Field');

-- name: InsertEventParticipation :exec
INSERT INTO Event_Participation (event_id, roll_no, rsvp) VALUES
(1, '2024UIN3356', 'Yes'),
(1, '2024UIN3357', 'Maybe'),
(2, '2024UIN3356', 'No'),
(3, '2024UIN3358', 'Yes');

-- name: InsertDonation :exec
INSERT INTO Donation (roll_no, amount, timestamp) VALUES
('2024UIN3356', 5000, '2025-02-17 11:30:00'),
('2024UIN3357', 2500, '2025-03-05 16:45:00');

-- name: InsertEmploymentHistory :exec
INSERT INTO Employment_History (roll_no, starting_year, ending_year, company, designation) VALUES
('2024UIN3356', 2024, 2025, 'Infosys', 'Developer'),
('2024UIN3357', 2024, 2024, 'Tata Steel', 'Engineer'),
('2024UIN3358', 2024, 2025, 'L&T', 'Site Manager');

-- name: InsertAlumniCredentials :exec
INSERT INTO credentials (username, password, user_type) VALUES
('2024UIN3356', 'Karan@123', 'alumni'),
('2024UIN3357', 'Sana@123', 'alumni'),
('2024UIN3358', 'Piyush@123', 'alumni');

-- name: InsertStaffCredentials :exec
INSERT INTO credentials (username, password, user_type) VALUES
('staff101', 'ProfDSingh@123', 'staff'),
('staff102', 'ProfAJohn@123', 'staff'),
('staff103', 'ProfRShah@123', 'staff');

