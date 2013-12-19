package bountyforcode

import (
	"fmt"
	. "github.com/abhiyerra/bountyforcode/app/models"
	"github.com/abhiyerra/scalpy"
	"github.com/gorilla/mux"
	"log"
	"net/http"
)

func IssuesHandler(w http.ResponseWriter, r *http.Request) {
	issues := FindAllIssues()
	RenderJson(w, issues)
}

func CreateIssueHandler(w http.ResponseWriter, r *http.Request) {
	issue_url := r.FormValue("issue-url")
	log.Println("New Issue", issue_url)

	scalp := scalpy.ScalpUrl(issue_url)
	if scalp == nil {
		fmt.Fprintf(w, "Issue doesn't exist!")
	} else {
		issue := NewIssue(scalp)
		log.Printf("%v\n", issue)

		http.Redirect(w, r, fmt.Sprintf("/issues/%s/contribute", issue.Id), 302)
	}
}

func ShowIssueHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	issue_id := vars["id"]

	issue := FindIssue(issue_id)

	if issue == nil {
		log.Printf("Can't find issue %s\n", issue_id)
		fmt.Fprintf(w, "Can't find issue")
		return
	}

	RenderJson(w, issue)
}
