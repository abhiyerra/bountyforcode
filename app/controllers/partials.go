package bountyforcode

import (
	"bytes"
	. "github.com/abhiyerra/bountyforcode/app/models"
	"log"
	"text/template"
)

func DiscoverPartial(issues []Issue) string {
	t, err := template.ParseFiles("views/_discover.tmpl")
	if err != nil {
		log.Printf("%v\n", err)
	}

	buffer := new(bytes.Buffer)

	t.Execute(buffer, issues)
	return buffer.String()
}
