package bountyforcode

import (
	"bytes"
	"fmt"
	. "github.com/abhiyerra/bountyforcode/app/models"
	"github.com/abhiyerra/scalpy"
	"log"
	"net/http"
	"text/template"
)

func CreateBountyHandler(w http.ResponseWriter, r *http.Request) {
	page := &Page{
		Title:    "New Bounty",
		ViewFile: "views/bounty_confirm.tmpl",
	}

	var vals struct {
		Repo string
	}

	fmt.Println("Parse", r.FormValue("issue-url"))
	scalp := scalpy.ScalpUrl(r.FormValue("issue-url"))
	if scalp == nil {
		vals.Repo = "Issue doesn't exist"
	} else {
		log.Printf("%v\n", NewIssue(scalp))
	}

	t, err := template.ParseFiles(page.ViewFile)
	if err != nil {
		log.Printf("%v\n", err)
	}

	buffer := new(bytes.Buffer)

	t.Execute(buffer, vals)
	page.Content = buffer.String()

	page.RenderLayout(w)
}

func ShowBountyHandler(w http.ResponseWriter, r *http.Request) {
	page := &Page{
		Title:    "New Bounty",
		ViewFile: "views/root_index.tmpl",
	}

	var vals struct {
		Repo        string
		OriginalUrl string
	}

	vals.OriginalUrl = "123"

	t, err := template.ParseFiles(page.ViewFile)
	if err != nil {
		log.Printf("%v\n", err)
	}

	buffer := new(bytes.Buffer)

	t.Execute(buffer, vals)
	page.Content = buffer.String()

	page.RenderLayout(w)
}

func CoinbaseCallbackHandler(w http.ResponseWriter, r *http.Request) {

	fmt.Fprintf(w, "ok")
}
