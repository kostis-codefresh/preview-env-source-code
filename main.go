package main

import (
	"fmt"
	"net/http"
	"os"
)

type gitInformation struct {
	branch string
	hash   string
}

func main() {

	port := os.Getenv("PORT")
	if len(port) == 0 {
		port = "8080"
	}

	branch := os.Getenv("GIT_BRANCH")
	if len(branch) == 0 {
		branch = "unknown"
	}

	hash := os.Getenv("GIT_HASH")
	if len(hash) == 0 {
		hash = "unknown"
	}

	// Kubernetes check if app is ok
	http.HandleFunc("/health/live", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "up")
	})

	// Kubernetes check if app can serve requests
	http.HandleFunc("/health/ready", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "yes")
	})

	gitInfo := gitInformation{}
	gitInfo.branch = branch
	gitInfo.hash = hash

	fmt.Printf("Preview app for branch %s is listening now at port %s\n", branch, port)
	http.Handle("/", &gitInfo)
	err := http.ListenAndServe(":"+port, nil)
	if err != nil {
		fmt.Printf("Failed to start server at port 8080: %v", err)
	}
}

func (g *gitInformation) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	fmt.Fprintf(w, "<h2>I am an application inside a preview environment<h2> <h3>My properties are:</h3><ul>")
	fmt.Fprintf(w, "<li>branch: "+g.branch+"</li>")
	fmt.Fprintf(w, "<li>git hash: "+g.hash+"</li>")
	fmt.Fprintf(w, "</ol>")

}
