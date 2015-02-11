package main

import (
	"../../builder"
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"html/template"
	"io"
	"io/ioutil"
	"log"
	"net/http"
)

var templ = template.Must(template.New("builder").Parse(templateStr))

const templateStr = `
<html>
<head>
<title>Builder</title>
<script src="//code.jquery.com/jquery-2.1.3.min.js"></script>
<link rel="stylesheet" href="https://maxcdn.bootstrapcdn.com/bootstrap/3.3.2/css/bootstrap.min.css">
<link rel="stylesheet" href="https://maxcdn.bootstrapcdn.com/bootstrap/3.3.2/css/bootstrap-theme.min.css">
<script src="https://maxcdn.bootstrapcdn.com/bootstrap/3.3.2/js/bootstrap.min.js"></script>
<style>
.code {
    display: none
}
</style>
</head>
<body>
<table id="jobs" class="table table-striped">
<thead>
  <th>ID</th>
  <th>Title</th>
  <th>Descriptiton</th>
  <th>Port</th>
  <th>Created</th>
  <th>Status</th>
  <th>Diff</th>
</thead>
{{range .}}
<tr>
   <td>{{.ID}}</td>
   <td>{{.Title}}</td>
   <td>{{.Descr}}</td>
   <td>{{.Port}}</td>
   <td>{{.Created}}</td>
   <td>{{.Status}}</td>
   <td class="codeParent">diff<pre class="code">{{.Diffdata}}</pre></td>
</tr>
{{end}}
</table>
<script>
$('.codeParent').click(function() {
    $(this).find('.code').toggle();
});
</script>
</body>
</html>
`

func statusUpdate(res http.ResponseWriter, req *http.Request) {
	vars := mux.Vars(req)
	job := vars["job"]
	status := vars["status"]

	db, err := builder.Connect()
	defer db.Close()
	if err != nil {
		log.Fatalf("Can't connect to DB: %v", err)
	}

	row, err := db.Query(`update jobs set status = $1 where id = $2`, status, job)
	if err != nil {
		log.Fatalf("Can't update status")
	}

	if err := json.NewEncoder(res).Encode(row); err != nil {
		panic(err)
	}
}

func showJobs(w http.ResponseWriter, req *http.Request) {
	db, err := builder.Connect()
	defer db.Close()
	builder.LogFail(err, "Can't connect to DB: %v")

	jobs, err := builder.GetJobs(db)
	builder.LogFail(err, "Can't get jobs: %v")

	err = templ.Execute(w, jobs)
	if err != nil {
		panic(err)
	}
}

func sendWork(res http.ResponseWriter, req *http.Request) {
	db, err := builder.Connect()
	defer db.Close()
	builder.LogFail(err, "Can't connect to DB: %v")

	jobs, err := builder.GetJobs(db)
	builder.LogFail(err, "Can't get jobs: %v")

	res.Header().Set("Content-Type", "application/json")

	if err := json.NewEncoder(res).Encode(jobs); err != nil {
		panic(err)
	}
}

func newJob(w http.ResponseWriter, req *http.Request) {
	var resp = builder.Resp{}
	body, err := ioutil.ReadAll(io.LimitReader(req.Body, 1048576))
	if err := req.Body.Close(); err != nil {
		resp.Error = err.Error()
	}

	var job = builder.Job{}
	if err := json.Unmarshal(body, &job); err != nil {
		panic(err)
	}

	db, err := builder.Connect()
	defer db.Close()
	if err != nil {
		panic(err)
	}

	diffid, err := builder.CreateDiff(db, string(job.Diffdata))
	if err != nil {
		resp.Error = err.Error()
	}

	job.Diff = diffid
	jobid, err := builder.CreateJob(db, &job)
	if err != nil {
		resp.Error = err.Error()
	}

	resp.JobID = jobid

	err = json.NewEncoder(w).Encode(resp)
	if err != nil {
		panic(err)
	}
}

func main() {
	r := mux.NewRouter()

	r.HandleFunc("/status/{job}/{status}", statusUpdate)
	r.HandleFunc("/new", newJob).Methods("POST")
	r.HandleFunc("/jobs", sendWork)
	r.HandleFunc("/", showJobs)

	http.Handle("/", r)
	fmt.Println("Listening on :8001")
	http.ListenAndServe(":8001", nil)
}
