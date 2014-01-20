package bountyforcode

import (
	. "github.com/abhiyerra/bountyforcode/app/models"
	"github.com/gorilla/mux"
	"log"
	"net/http"
)

func NewBountyHandler(w http.ResponseWriter, r *http.Request) {
	// Get the User
	user := GetSessionUser(r)
	log.Println("user", user.Id)
	if user == nil {
		RenderJson(w, StatusResponse{
			Status:  false,
			Message: "Need to register/login first",
		})
		return
	}

	// Get the Issue
	vars := mux.Vars(r)
	log.Println(vars)
	issue_id := vars["id"]
	issue := FindIssue(issue_id)
	if issue == nil {
		log.Printf("Can't find issue %s\n", issue_id)
		RenderJson(w, StatusResponse{
			Status:  false,
			Message: "Can't find issue",
		})
		return
	}

	// Make a bounty for the user
	bounty := NewBounty(issue, user)
	if bounty == nil {
		RenderJson(w, StatusResponse{
			Status:  false,
			Message: "Can't create bounty",
		})
		return
	}

	RenderJson(w, bounty)
}

func BountiesHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	log.Println(vars)
	issue_id := vars["id"]

	log.Println(issue_id)
	bounties := FindBounties(issue_id)

	RenderJson(w, bounties)
}
