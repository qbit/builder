package builder

import (
	"database/sql"
	_ "github.com/lib/pq"
	"time"
)

type Job struct {
	Id       int
	Created  time.Time
	Title    string
	Descr    string
	Port     string
	Diff     int
	Active   bool
	Diffdata string
}

type Diff struct {
	Diffdata string
}

type Jobs []*Job

/*
func (jobs *Jobs) New() interface{} {
	j := &Job{}
	*jobs = append(*jobs, j)
	return j
}
*/

func Connect() (*sql.DB, error) {
	db, err := sql.Open("postgres", "user=postgres dbname=qbit sslmode=disable")

	if err != nil {
		return nil, err
	}

	return db, nil
}

func CreateJob(db *sql.DB, job *Job) (int, error) {
	// insert diff first, then do this below
	var id int

	err := db.QueryRow(`INSERT INTO jobs (title, descr, port, diff) values ($1, $2, $3, $4) returning id`, job.Title, job.Descr, job.Port, job.Diff).Scan(&id)

	return id, err
}

func CreateDiff(db *sql.DB, diff string) (int, error) {
	var id int

	err := db.QueryRow(`INSERT INTO diffs (id, diffdata) values (DEFAULT, $1) RETURNING id`, diff).Scan(&id)

	return id, err
}

//func GetJobs(db *sql.DB) (*sql.Rows, error) {
func GetJobs(db *sql.DB) (Jobs, error) {

	var job = Job{}
	var jobs = Jobs{}

	rows, err := db.Query(`SELECT jobs.id, created, title, descr, port, diffdata FROM jobs left join diffs on (diffs.id = jobs.diff) where active = true`)

	defer rows.Close()

	for rows.Next() {
		err := rows.Scan(&job.Id, &job.Created, &job.Title, &job.Descr, &job.Port, &job.Diffdata)

		if err != nil {
			return nil, err
		}
		rows.Scan(&job.Id, &job.Created, &job.Title, &job.Descr, &job.Port, &job.Diffdata)
		jobs = append(jobs, &job)
	}

	if err != nil {
		return nil, err
	}

	return jobs, nil
}
