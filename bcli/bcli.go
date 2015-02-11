package main

import (
	"../../builder"
	"bufio"
	"bytes"
	"encoding/json"
	"flag"
	"log"
	"net/http"
	"os"
	"strings"
)

func addRemote(url *string, job builder.Job) (int, error) {
	j, err := json.Marshal(job)
	if err != nil {
		builder.LogFail(err, "Can't marshal job: %v")
	}

	resp, err := http.Post(*url, "application/json", bytes.NewReader(j))
	if err != nil {
		log.Fatalf("Can't connect to '%s'", url)
	}
	res := &builder.Resp{}
	if err := json.NewDecoder(resp.Body).Decode(res); err != nil {
		log.Fatalf("Invalid response from server! %v", err)
	}

	resp.Body.Close()

	if res.Error != "" {
		log.Fatalf(res.Error)
	}

	return res.JobID, nil
}

func addLocal(job builder.Job, lines []string) (int, error) {
	dbuser, dbpass, dbname := builder.GetCreds()

	db, err := builder.Connect(dbuser, dbpass, dbname)
	builder.LogFail(err, "Can't connect to DB: %v")

	jobid, err := builder.CreateJob(db, &job)
	builder.LogFail(err, "Can't create job: %v")

	log.Printf("Created jobid %d", jobid)
	return jobid, nil
}

func main() {
	var url = flag.String("url", "http://localhost:8001/new", "URL of server")
	var dfile = flag.String("diff", "", "Diff to be tested")
	var title = flag.String("title", "", "Title of job")
	var desc = flag.String("desc", "", "Description of job")
	var port = flag.String("port", "", "Port being updated")
	var lines []string

	flag.Parse()
	file, err := os.Open(*dfile)

	defer file.Close()
	if err != nil {
		log.Fatalf("Can't open file! - %v", err)
	}

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}

	var job = builder.Job{Title: *title, Descr: *desc, Port: *port, Diffdata: strings.Join(lines, "\n")}
	jobid, err := addRemote(url, job)
	if err != nil {
		log.Fatalf("Can't blablab")
	}

	log.Printf("Title: %s, Desc: %s, JobID: %d", *title, *desc, jobid)
}
