package db

const authenticateUser = `SELECT password from credentials where username = ? and user_type = ?`

const checkSession = `
select username, expiry from sessions where session_id = ?
and username in (select username from credentials where user_type = ?)
;`

const deleteSession = `DELETE FROM sessions where session_id = ?;`

const getAlumni = `
SELECT a.roll_no, a.name, a.email, a.phone_no, SUBSTR(a.roll_no, 1, 4), d.name
FROM alumni a join department d on a.department_id = d.id
WHERE a.roll_no = ?
;`
const getFaculty = `
SELECT f.id, f.name, d.name FROM faculty f join department d on f.department_id = d.id where f.id = ?
;`
const searchDirectory = `
SELECT a.name, a.roll_no, SUBSTR(a.roll_no, 1, 4), d.name 
FROM alumni a join department d on a.department_id = d.id
WHERE 
	((a.name LIKE '%' || ?1 || '%' OR ?1 = '') OR (a.roll_no LIKE '%' || ?1 || '%' OR ?1 = '')) AND
	(d.name = ?2 OR ?2 = '') AND
	(SUBSTR(a.roll_no, 1, 4) = ?3 OR ?3 = '')
;`

const getAcademicHistory = `
SELECT semester, SGPA
FROM Academic_History
WHERE roll_no = ?
ORDER BY semester ASC
;`

const getPrevDonations = `
SELECT a.name, d.amount, d.message, d.timestamp 
FROM donation d inner join alumni a on a.roll_no = d.roll_no
ORDER BY timestamp DESC
;`

const getPrev3Donations = `
SELECT a.name, d.amount 
FROM donation d inner join alumni a on a.roll_no = d.roll_no
ORDER BY timestamp DESC
;`

const getTotalDonationsByAlum = `
SELECT sum(d.amount) 
FROM donation d join alumni a on a.roll_no = d.roll_no
where a.roll_no = ?
;`

const addDonation = `
INSERT into DONATION(roll_no, amount, message)
VALUES (?,?,?)
;`

const getEmploymentHistory = `
SELECT id, roll_no, starting_year, ending_year, company, designation, location
FROM Employment_History where roll_no = ?
ORDER BY starting_year DESC, ending_year DESC
;`

const updateEmploymentHistory = `
INSERT OR REPLACE into Employment_History(id, roll_no, starting_year, ending_year, company, designation, location)
values (?,?,?,?,?,?,?)
;`

const addEmploymentHistory = `
INSERT into Employment_History(roll_no, starting_year, ending_year, company, designation, location)
values (?,?,?,?,?,?)
;`

const deleteEmploymentHistory = `
DELETE FROM Employment_History 
WHERE id = ?
	;`

const getNoticesAndEvents = `
SELECT 
    Notice.id,
    title,
    description,
    NULL AS date,
    NULL AS location,
    'Notice' AS type,
    NULL AS rsvp,
    created_at
FROM Notice

UNION ALL

SELECT 
    Event.id,
    title,
    description,
    DATETIME(date),
    location,
    'Event' AS type,
    rsvp,
    created_at
FROM Event LEFT OUTER JOIN (

	SELECT * from Event_Participation where roll_no = ?) a
ON Event.id = a.id

ORDER BY created_at DESC;
`

const getNoticesAndEventsAndResponses = `
SELECT 
    Notice.id,
    title,
    description,
    NULL AS date,
    NULL AS location,
    'Notice' AS type,
    created_at,
    NULL,
    NULL,
    NULL,
    NULL
FROM Notice

UNION ALL

SELECT 
    e.id,
    e.title,
    e.description,
    e.date,
    e.location,
    'Event' AS type,
    e.created_at,
    COUNT(CASE WHEN ep.rsvp = 'Accept' THEN 1 END),
    COUNT(CASE WHEN ep.rsvp = 'Maybe' THEN 1 END),
    COUNT(CASE WHEN ep.rsvp = 'Decline' THEN 1 END),
    (SELECT COUNT(*) 
     FROM Alumni a 
     WHERE a.roll_no NOT IN (
         SELECT ep2.roll_no 
         FROM Event_Participation ep2 
         WHERE ep2.id = e.id
     )
    )
FROM 
    Event e
LEFT JOIN 
    Event_Participation ep 
ON 
    e.id = ep.id
GROUP BY 
	e.id
	ORDER BY created_at;
`

const addRSVP = `
INSERT INTO Event_Participation(id, roll_no, rsvp)
VALUES (?,?,?)
;`

const addNotice = `
INSERT INTO Notice(title, description)
VALUES (?,?)
;`

const addEvent = `
INSERT INTO Event(title, description, date, location)
VALUES (?,?,?,?)
;`
