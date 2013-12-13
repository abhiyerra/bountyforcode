package bountyforcode

import (
	"bytes"
	"fmt"
	. "github.com/abhiyerra/bountyforcode/app/models"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"text/template"
)

func ProjectRootHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	project_identifier := vars["subdomain"]

	type ProjectPage struct {
		Page
		Discover string
	}

	issues := FindProjectIssues(project_identifier)
	page := &ProjectPage{
		Page: Page{
			Title:    project_identifier,
			ViewFile: GetView("project_index.tmpl"),
		},
		Discover: RenderPartial(GetView("partials/_discover.html"), issues),
	}

	fmt.Printf("%v, %v\n", page.Title, issues)

	t, err := template.ParseFiles(page.ViewFile)
	if err != nil {
		log.Printf("%v\n", err)
	}

	buffer := new(bytes.Buffer)

	t.Execute(buffer, page)
	page.Content = buffer.String()

	page.RenderLayout(w)
}
