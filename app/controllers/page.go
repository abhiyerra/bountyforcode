package bountyforcode

import (
	"log"
	"net/http"
	"text/template"
)

type Page struct {
	Title string
	Body  string

	Content string

	ViewFile string
	Layout   string
}

func (p *Page) RenderLayout(w http.ResponseWriter) {
	var layout string = p.Layout

	if layout == "" {
		layout = "views/layout.tmpl"
	}

	t, err := template.ParseFiles(layout)
	if err != nil {
		log.Printf("%v\n", err)
	}

	t.Execute(w, p)
}
