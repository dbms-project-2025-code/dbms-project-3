package db

// SQL constants for schema creation
const createCredentials = `
CREATE TABLE IF NOT EXISTS credentials (
    username TEXT PRIMARY KEY,
    password TEXT NOT NULL,
    user_type TEXT NOT NULL CHECK(user_type IN ('staff', 'alumni'))
);
`
const createAlumni = `
CREATE TABLE IF NOT EXISTS Alumni (
    roll_no TEXT PRIMARY KEY,
    name TEXT NOT NULL,
    email TEXT UNIQUE,
    phone_no TEXT,
    department_id INTEGER,
    FOREIGN KEY (department_id) REFERENCES Department(id)
);
`

const createEmploymentHistory = `
CREATE TABLE IF NOT EXISTS Employment_History (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    roll_no TEXT,
    starting_year INTEGER,
    ending_year INTEGER,
    company TEXT,
    designation TEXT,
    location TEXT,
    FOREIGN KEY (roll_no) REFERENCES Alumni(roll_no)
);
`


const createAcademicHistory = `
CREATE TABLE IF NOT EXISTS Academic_History (
    roll_no TEXT,
    SGPA REAL,
    semester INT CHECK(semester >= 1 and semester <= 8),
    PRIMARY KEY (roll_no ,semester),
    FOREIGN KEY (roll_no) REFERENCES Alumni(roll_no)
);
`

const createEvent = `
CREATE TABLE IF NOT EXISTS Event (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    title TEXT NOT NULL,
    description TEXT,
    date DATETIME,
    location TEXT,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP
);
`

const createDonation = `
CREATE TABLE IF NOT EXISTS Donation (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    roll_no TEXT,
    amount REAL NOT NULL check(amount > 0),
    message TEXT NULL,
    timestamp DATETIME DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (roll_no) REFERENCES Alumni(roll_no)
);
`

const createDepartment = `
CREATE TABLE IF NOT EXISTS Department (
    id INTEGER PRIMARY KEY,
    name TEXT NOT NULL
);
`

const createFaculty = `
CREATE TABLE IF NOT EXISTS Faculty (
    id TEXT PRIMARY KEY,
    name TEXT NOT NULL,
    department_id INTEGER,
    FOREIGN KEY (department_id) REFERENCES Department(id)
);
`

const createSessions = `
CREATE TABLE IF NOT EXISTS sessions (
    session_id TEXT PRIMARY KEY,
    username TEXT NOT NULL,
    expiry TIMESTAMP NOT NULL,
    FOREIGN KEY (username) REFERENCES credentials(username)
);
`

const createNotice = `
CREATE TABLE IF NOT EXISTS Notice (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    title TEXT NOT NULL,
    description TEXT,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP
);
`

const createEventParticipation = `
CREATE TABLE IF NOT EXISTS Event_Participation (
    id INTEGER,
    roll_no TEXT,
    rsvp TEXT CHECK(rsvp IN ('Accept','Decline','Maybe')),
    PRIMARY KEY (id, roll_no),
    FOREIGN KEY (id) REFERENCES Event(id),
    FOREIGN KEY (roll_no) REFERENCES Alumni(roll_no)
);
`
const createDonation_Backup = `
CREATE TABLE IF NOT EXISTS Donation_Backup (
    id INTEGER PRIMARY KEY,
    roll_no TEXT,
    amount REAL,
    message TEXT,
    timestamp DATETIME,
    backup_timestamp DATETIME DEFAULT CURRENT_TIMESTAMP
)
;`

const createBackupTrigger = `
CREATE TRIGGER IF NOT EXISTS backup_donation_after_insert
AFTER INSERT ON Donation
BEGIN
    INSERT OR REPLACE INTO Donation_Backup (id, roll_no, amount, message, timestamp)
    VALUES (NEW.id, NEW.roll_no, NEW.amount, NEW.message, NEW.timestamp);
END;
;`

var SchemaTables = []string{
	createDepartment,
	createFaculty,
	createAlumni,
	createAcademicHistory,
	createEvent,
	createNotice,
	createEventParticipation,
	createDonation,
	createEmploymentHistory,
	createCredentials,
	createSessions,
	createDonation_Backup,
	createBackupTrigger,
}
