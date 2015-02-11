package builder

import (
	"database/sql"
	"log"
	"time"

	// pull in postgres
	_ "github.com/lib/pq"
	// pgenv to parse env vars to get pg info
	"github.com/qbit/pgenv"
)

// Job represents the row in the db
type Job struct {
	ID       int
	Created  time.Time
	Title    string
	Descr    string
	Port     string
	Diff     int
	Statid   int
	Active   bool
	Diffdata string
	Status   string
}

// Jobs for when you need more than one job!
type Jobs []*Job

// Diff is also a row in the db. Not really needed yet.
type Diff struct {
	Diffdata string
}

// Resp is used for sending job id or errors to the client
type Resp struct {
	JobID int
	Error string
}

// LogFail helper function for failing
func LogFail(err error, msg string) {
	if err != nil {
		log.Fatal(msg, err)
	}
}

// Connect to the db
func Connect() (*sql.DB, error) {
	var s = pgenv.ConnStr{}
	s.SetDefaults()
	db, err := sql.Open("postgres", s.ToString())

	return db, err
}

// CreateJob takes a Job and inserts it into the db.
// TODO need a interface or something to allow for
// more dynamic assignment of values to struct
func CreateJob(db *sql.DB, job *Job) (int, error) {
	// insert diff first, then do this below
	var id int
	err := db.QueryRow(`INSERT INTO jobs (title, descr, port, diff) values ($1, $2, $3, $4) returning id`, job.Title, job.Descr, job.Port, job.Diff).Scan(&id)
	return id, err
}

// CreateDiff inserts a diff into the db
func CreateDiff(db *sql.DB, diff string) (int, error) {
	var id int
	err := db.QueryRow(`INSERT INTO diffs (id, diffdata) values (DEFAULT, $1) RETURNING id`, diff).Scan(&id)
	return id, err
}

// GetJobs get a list of active jobs from the db
func GetJobs(db *sql.DB) (Jobs, error) {
	var jobs = Jobs{}

	rows, err := db.Query(`
SELECT
 jobs.id,
 created,
 title,
 descr,
 port,
 diffdata,
 stat.status
FROM
 jobs
 left join diffs on (diffs.id = jobs.diff)
 left join stat on (stat.id = jobs.status)
where
active = true`)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	for rows.Next() {
		var job = Job{}
		err := rows.Scan(&job.ID, &job.Created, &job.Title, &job.Descr, &job.Port, &job.Diffdata, &job.Status)
		if err != nil {
			return nil, err
		}
		jobs = append(jobs, &job)
	}

	return jobs, nil
}
