package db

type credential struct {
	Password  string
	User_type string
}

var users = map[string]credential{
	"2020UIN3356": {"Karan@123", "alumni"},
	"2020UIN3357": {"Sana@123", "alumni"},
	"2020UIN3358": {"Piyush@123", "alumni"},
	"staff101":    {"ProfDSingh@123", "staff"},
	"staff102":    {"ProfAJohn@123", "staff"},
	"staff103":    {"ProfRShah@123", "staff"},
}

const insertCredentials = `
INSERT OR REPLACE INTO credentials (username, password, user_type) VALUES (?,?,?);
`

const insertAcademicHistory = `
INSERT OR REPLACE INTO Academic_History (roll_no, SGPA, semester) VALUES
('2020UIN3356', 8.5, 1),('2020UIN3356', 8.3, 2),('2020UIN3356', 8.7, 3),('2020UIN3356', 8.8, 4),('2020UIN3356', 8.9, 5),('2020UIN3356', 9.0, 6),('2020UIN3356', 9.1, 7),('2020UIN3356', 9.2, 8),

('2020UIN3357', 7.9, 1),('2020UIN3357', 8.1, 2),('2020UIN3357', 8.0, 3),('2020UIN3357', 7.8, 4),('2020UIN3357', 8.2, 5),('2020UIN3357', 8.0, 6),('2020UIN3357', 8.3, 7),('2020UIN3357', 8.2, 8),

('2020UIN3358', 8.0, 1),('2020UIN3358', 8.2, 2),('2020UIN3358', 8.1, 3),('2020UIN3358', 8.4, 4),('2020UIN3358', 8.3, 5),('2020UIN3358', 8.2, 6),('2020UIN3358', 8.5, 7),('2020UIN3358', 8.7, 8)
`

const insertAlumni = `
INSERT OR REPLACE INTO Alumni (roll_no, name, email, phone_no, department_id) VALUES
('2020UIN3356', 'Karan Patel', 'karan.p@domain.com', '9123456780', 3),
('2020UIN3357', 'Sana Verma', 'sana.v@domain.com', '9123456781', 3),
('2020UIN3358', 'Piyush Das', 'piyush.d@domain.com', '9123456782', 3),

('2019UCA2260', 'Ananya Singh', 'ananya.s@domain.com', '9123456783', 1),
('2019UCA2261', 'Rohit Sharma', 'rohit.s@domain.com', '9123456784', 1),
('2019UCA2262', 'Ishita Mehta', 'ishita.m@domain.com', '9123456785', 1),
('2019UCA2263', 'Vivek Nair', 'vivek.n@domain.com', '9123456786', 1),

('2018UME4464', 'Neha Kapoor', 'neha.k@domain.com', '9123456787', 2),
('2018UME4465', 'Arjun Reddy', 'arjun.r@domain.com', '9123456788', 2),
('2018UME4466', 'Tanya Bhatt', 'tanya.b@domain.com', '9123456789', 2),
('2018UME4467', 'Manish Gupta', 'manish.g@domain.com', '9123456790', 2);
`

const insertDepartments = `
INSERT OR REPLACE INTO Department (id, name) VALUES
(1, 'Computer Science'),
(2, 'Mechanical'),
(3, 'Information Technology')
`

const insertDonation = `
INSERT OR REPLACE INTO Donation (id, roll_no, amount) VALUES
(1, '2020UIN3356', 5000),
(2, '2020UIN3357', 2500)
;`

const insertEmploymentHistory = `
INSERT OR REPLACE INTO Employment_History (id, roll_no, starting_year, ending_year, company, designation, location) VALUES
(1, '2020UIN3356', 2020, 2025, 'Infosys', 'Developer', 'New Delhi'),
(2, '2020UIN3357', 2020, 2025, 'Tata Steel', 'Engineer', 'Mumbai'),
(3, '2020UIN3358', 2020, 2025, 'L&T', 'Site Manager', 'Chennai')
`

const insertEventData = `
INSERT OR REPLACE INTO Event (id, title, description, date, location) VALUES
(1, 'Alumni Meet', 'Annual alumni gathering', '2025-01-15 18:00:00', 'Seminar Hall'),
(2, 'Tech Talk', 'AI trends lecture', '2025-03-10 14:30:00', 'Main Hall'),
(3, 'Sports Fest', 'Inter-dept tournament', '2025-02-20 09:00:00', 'Sports Field')
`
const insertNoticeData = `
INSERT OR REPLACE INTO Notice (id, title, description) VALUES
(1, 'Alumni Meet 2025', 'Join us for the Annual Alumni Meet on 20th December 2025 at NSUT Auditorium.'),
(2, 'Placement Drive Announcement', 'Top companies are visiting campus next month. Register before 25th November 2025.'),
(3, 'Scholarship Applications Open', 'Applications for Merit-based scholarships are now open. Apply through the student portal.'),
(4, 'Website Maintenance', 'The alumni portal will be under maintenance from 10 PM to 2 AM on 15th November 2025.'),
(5, 'New Research Grants', 'Faculty and alumni are invited to apply for the NSUT Innovation Grant 2025.'),
(6, 'Internship Opportunities', 'Tech startups are offering internship positions for 3rd-year students. Check the career section.'),
(7, 'Convocation Ceremony 2025', 'The annual convocation will be held on 10th January 2026. Graduates are requested to confirm attendance.'),
(8, 'Blood Donation Camp', 'Join the Blood Donation Camp organized by the NSS unit on 18th November 2025 at the main hall.'),
(9, 'Guest Lecture Series', 'Attend the AI and Robotics guest lecture by Dr. Meera Iyer on 22nd November 2025.'),
(10, 'Hostel Allotment Notice', 'Hostel allotment results for the 2025 batch are now available on the official website.');
`
const insertEventParticipation = `
INSERT OR REPLACE INTO Event_Participation (id, roll_no, rsvp) VALUES
(1, '2020UIN3356', 'Accept'),
(1, '2020UIN3357', 'Maybe'),
(2, '2020UIN3356', 'Decline'),
(3, '2020UIN3358', 'Accept')
`

const insertFaculty = `
INSERT OR REPLACE INTO Faculty (id, name, department_id) VALUES
('staff101', 'Prof. D. Singh', 1),
('staff102', 'Prof. A. John', 2),
('staff103', 'Prof. R. Shah', 3)
`

var insertData = []string{
	insertDepartments,
	insertFaculty,
	insertAlumni,
	insertAcademicHistory,
	insertEmploymentHistory,
	insertNoticeData,
	insertEventData,
	insertEventParticipation,
	insertDonation,
}
