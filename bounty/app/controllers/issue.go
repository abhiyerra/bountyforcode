package bountyforcode

import (
	"encoding/json"
	"fmt"
	. "github.com/abhiyerra/bountyforcode/bounty/app/models"
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
	decoder := json.NewDecoder(r.Body)
	var t struct {
		IssueUrl string `json:"issue_url"`
	}
	err := decoder.Decode(&t)
	if err != nil {
		log.Println(err)
	}
	issue_url := t.IssueUrl
	log.Println("New Issue", issue_url)

	scalp := scalpy.ScalpUrl(issue_url)
	if scalp == nil {
		RenderJson(w, StatusResponse{
			Status:  false,
			Message: "Issue doesn't exist!",
		})
	} else {
		issue := NewIssue(scalp)
		log.Printf("%v\n", issue)

		RenderJson(w, issue)
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
