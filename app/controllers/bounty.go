package bountyforcode

import (
	"fmt"
	. "github.com/abhiyerra/bountyforcode/app/models"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"strconv"
)

func NewBountyHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	issue_id := vars["id"]

	user_id := GetSessionUserId(r)
	if user_id == "" {
		http.Redirect(w, r, "/register", 302)
		return
	}

	issue := FindIssue(issue_id)
	if issue == nil {
		log.Printf("Can't find issue %s\n", issue_id)
		fmt.Fprintf(w, "Can't find issue")
		return
	}

	user_id_int, _ := strconv.Atoi(user_id)
	bounty := NewBounty(issue, user_id_int)
	if bounty == nil {
		fmt.Fprintf(w, "Can't create bounty")
		return
	}

	RenderJson(w, bounty)
}
