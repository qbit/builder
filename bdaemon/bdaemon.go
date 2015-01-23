package main

import (
	"../../builder"
	"encoding/json"
	"flag"
	"log"
	"net/http"
)

func GetJobs(url *string) builder.Jobs {
	var jobs = builder.Jobs{}

	resp, err := http.Get(*url)
	if err != nil {
		log.Fatalf("Can't get jobs: %v", err)
	}
	log.Printf("%s", resp.Body)
	if err := json.NewDecoder(resp.Body).Decode(jobs); err != nil {
		log.Fatalf("Invalid response from server! %v", err)
	}
	resp.Body.Close()

	return jobs
}

func main() {
	var url = flag.String("url", "http://localhost:8001/jobs", "URL of build server")
	flag.Parse()

	GetJobs(url)
}
