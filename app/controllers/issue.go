package bountyforcode

import (
	"bytes"
	"fmt"
	. "github.com/abhiyerra/bountyforcode/app/models"
	"github.com/abhiyerra/scalpy"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"text/template"
)

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

type IssuePage struct {
	Page
	Issue *Issue
}

func GetIssue(r *http.Request) (issue *Issue) {
	vars := mux.Vars(r)
	issue_id := vars["id"]
	
	log.Printf("%v\n", vars)
	issue = FindIssue(issue_id)
	if issue == nil {
		log.Printf("Can't find issue %s\n", issue_id)
		return nil
	}
	
	return issue
}

func ShowIssueHandler(w http.ResponseWriter, r *http.Request) {
	issue := GetIssue(r)
	if issue == nil {
		fmt.Fprintf(w, "Can't find issue")
		return
	}
	
	page := &IssuePage{
		Page: Page{
			Title:    "New Bounty",
			ViewFile: "views/root_index.tmpl",
		},
		Issue: issue,
	}

	t, err := template.ParseFiles(page.ViewFile)
	if err != nil {
		log.Printf("Error %v\n", err)
	}

	buffer := new(bytes.Buffer)

	t.Execute(buffer, page)
	page.Content = buffer.String()

	page.RenderLayout(w)
}

func ContributeIssueHandler(w http.ResponseWriter, r *http.Request) {
	if GetSessionUserId(r) == "" {
		http.Redirect(w, r, "/register", 302)
	}

	issue := GetIssue(r)
	if issue == nil {
		fmt.Fprintf(w, "Can't find issue")
		return
	}

	// TODO: Generate the CoinbaseButton here so we have an identifier for who donated
	// TODO: Insert into bounties (coinbase_button, user_id, issue_id) 
	page := &IssuePage{
		Page: Page{
			Title:    "Contribute Bounty",
			ViewFile: "views/issues/contribute.tmpl",
		},
		Issue: issue,
	}

	t, err := template.ParseFiles(page.ViewFile)
	if err != nil {
		log.Printf("Error %v\n", err)
	}

	buffer := new(bytes.Buffer)

	t.Execute(buffer, page)
	page.Content = buffer.String()

	page.RenderLayout(w)
}

func CoinbaseCallbackHandler(w http.ResponseWriter, r *http.Request) {

	fmt.Fprintf(w, "ok")
}
