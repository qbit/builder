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
   <td>{{.Id}}</td>
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

func StatusUpdate(res http.ResponseWriter, req *http.Request) {
	vars := mux.Vars(req)
	job := vars["job"]
	status := vars["status"]

	db, err := builder.Connect()
	if err != nil {
		log.Fatalf("Can't connect to DB: %v", err)
	}

	row, err := db.Query(`update jobs set status = $1 where id = $2`, status, job)
	if err != nil {
		log.Fatalf("Can't update status")
	}

	log.Printf("%v", row)
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

	r.HandleFunc("/status/{job}/{status}", StatusUpdate)
	r.HandleFunc("/jobs", SendWork)
	r.HandleFunc("/", ShowJobs)

	http.Handle("/", r)
	fmt.Println("Listening on :8001")
	http.ListenAndServe(":8001", nil)
}
