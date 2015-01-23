package main

import (
	"../../builder"
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"flag"
	"io"
	"log"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

func MakePatch(job *builder.Job) string {
	// Turns out I don't need to do this.
	// I will keep it around for possible use later
	randBytes := make([]byte, 16)
	rand.Read(randBytes)
	filename := filepath.Join(os.TempDir(), hex.EncodeToString(randBytes)+".diff")

	file, err := os.Create(filename)
	if err != nil {
		log.Fatalf("%v", err)
	}

	_, err = io.WriteString(file, job.Diffdata)
	if err != nil {
		log.Fatalf("%v", err)
	}

	file.Close()
	return filename
}

func ApplyPatch(dir string, patch string) bool {
	wd, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}
	err = os.Chdir(dir)
	if err != nil {
		log.Fatal(err)
	}
	defer os.Chdir(wd)

	if err != nil {
		log.Fatalf("%v", err)
	}
	cmd := exec.Command("/usr/bin/patch", "-p0", "-E")
	cmd.Stdin = strings.NewReader(patch)
	err = cmd.Start()
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("Waiting for patching to finish...")
	err = cmd.Wait()
	if err != nil {
		log.Fatalf("Patching errored: %v", err)
	}

	return true
}

func GetJobs(url *string) builder.Jobs {
	var jobs = builder.Jobs{}

	resp, err := http.Get(*url)
	if err != nil {
		log.Fatalf("Can't get jobs: %v", err)
	}

	if err := json.NewDecoder(resp.Body).Decode(&jobs); err != nil {
		log.Fatalf("Invalid response from server! %v", err)
	}
	resp.Body.Close()

	return jobs
}

func main() {
	var url = flag.String("url", "http://localhost:8001/jobs", "URL of build server")
	var pdir = flag.String("pdir", "/usr/ports", "PORTSDIR to apply diffs in")

	flag.Parse()

	jobs := GetJobs(url)
	for job := range jobs {
		//fn := MakePatch(jobs[job])
		//log.Printf("New Job: %s:%s:%s", jobs[job].Title, jobs[job].Descr, fn)
		success := ApplyPatch(filepath.Join(*pdir, jobs[job].Port), jobs[job].Diffdata)
		if success {
			//log.Printf("New Job: %s:%s:%s", jobs[job].Title, jobs[job].Descr, fn)
			log.Printf("New Job: %s:%s", jobs[job].Title, jobs[job].Descr)
		} else {
			log.Printf("Failed to apply diff for: %s", jobs[job].Title)
		}
	}
}
