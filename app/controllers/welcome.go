package bountyforcode

import (
	"bytes"
	"fmt"
	. "github.com/abhiyerra/bountyforcode/app/models"
	"log"
	"net/http"
	"text/template"
)

func RootHandler(w http.ResponseWriter, r *http.Request) {
	type RootPage struct {
		Page
		Discover string
	}

	issues := FindProjectIssues("")

	page := &RootPage{
		Page: Page{
			Title:    "Welcome",
			ViewFile: GetView("root_index.tmpl"),
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
