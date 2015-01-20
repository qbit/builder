package main

import (
	"../../builder"
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"html/template"
	"log"
	"net/http"
)

var templ = template.Must(template.New("builder").Parse(templateStr))

const templateStr = `
<html>
<head>
<title>Builder</title>
</head>
<body>
<table id="jobs">

</table>
</body>
</html>
`

func StatusUpdate(res http.ResponseWriter, req *http.Request) {
	vars := mux.Vars(req)
	job := vars["job"]

	fmt.Println("Status for '%s'", job)
}

func ShowJobs(res http.ResponseWriter, req *http.Request) {
	db, err := builder.Connect()
	if err != nil {
		log.Fatalf("Can't connect to DB: %v", err)
	}

	jobs, err := builder.GetJobs(db)
	if err != nil {
		log.Fatalf("Can't get jobs: %v", err)
	}

	templ.Execute(res, jobs)
}

func SendWork(res http.ResponseWriter, req *http.Request) {
	db, err := builder.Connect()
	if err != nil {
		log.Fatalf("Can't connect tot DB: %v", err)
	}

	jobs, err := builder.GetJobs(db)
	if err != nil {
		log.Fatalf("Can't get jobs: %v", err)
	}

	if err := json.NewEncoder(res).Encode(jobs); err != nil {
		panic(err)
	}
}

func main() {
	r := mux.NewRouter()

	r.HandleFunc("/status/{job}", StatusUpdate)
	r.HandleFunc("/jobs", SendWork)
	r.HandleFunc("/", ShowJobs)

	http.Handle("/", r)
	fmt.Println("Listening on :8001")
	http.ListenAndServe(":8001", nil)
}
