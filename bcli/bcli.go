package main

import (
	"../../builder"
	"bufio"
	"flag"
	"log"
	"os"
	"strings"
)

func main() {
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

	log.Printf("Title: %s, Desc: %s", *title, *desc)

	db := builder.Connect()

	diffid, err := builder.CreateDiff(db, strings.Join(lines, "\n"))

	job := builder.Job{*title, *desc, *port, diffid}

	jobid, err := builder.CreateJob(db, &job)

	if err != nil {
		log.Fatalf("%v", err)
	}
	log.Printf("Created jobid %d", jobid)
}
