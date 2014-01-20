package main

import (
	"encoding/json"
	"github.com/abhiyerra/scalpy"
	"github.com/octokit/go-octokit/octokit"
	"log"
	"net/http"
)

type Issue struct {
	Scalp     *octokit.Issue
	ExpiresAt int
}

type User struct {
	Scalp     *octokit.User
	ExpiresAt int
}

var (
	Issues map[string]Issue
	Users  map[string]User
)

func IssueUpdater() {
	// TODO: Run this code at some point.
}

func IssuesHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	// access_token := r.FormValue("access_token")
	// if access_token == "" {
	// 	log.Println("No AccessToken given")
	// }

	issue_url := r.FormValue("issue_url")

	issue, ok := Issues[issue_url]
	if !ok {
		log.Printf("Making call for url %s", issue_url)

		issue = Issue{}
		scalp := scalpy.ScalpUrl(issue_url)
		issue.Scalp = scalp.GithubIssue()
		Issues[issue_url] = issue
	}

	b, _ := json.Marshal(issue.Scalp)
	w.Write(b)
}

func UserUpdater() {
	// TODO: Run this code at some point.
}

func UsersHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	github_login := r.FormValue("github_login")

	log.Println(github_login)

	user, ok := Users[github_login]
	if !ok {
		log.Printf("Making call for url %s", github_login)

		client := octokit.NewClient(nil)
		url, _ := octokit.UserURL.Expand(octokit.M{"user": github_login})

		user = User{}
		user.Scalp, _ = client.Users(url).One()

		Users[github_login] = user
	}

	b, _ := json.Marshal(user.Scalp)
	w.Write(b)
}

func main() {
	Issues = make(map[string]Issue)
	Users = make(map[string]User)

	go IssueUpdater()
	go UserUpdater()

	http.HandleFunc("/v1/users", UsersHandler)
	http.HandleFunc("/v1/issues", IssuesHandler)
	http.ListenAndServe(":3001", nil)
}
