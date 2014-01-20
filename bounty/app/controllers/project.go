package bountyforcode

import (
	. "github.com/abhiyerra/bountyforcode/app/models"
	"github.com/gorilla/mux"
	"log"
	"net/http"
)

func ProjectIssuesHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	project := vars["subdomain"]

	log.Println("Project Page Loaded for", project)
	issues := FindProjectIssues(project)

	RenderJson(w, issues)
}
