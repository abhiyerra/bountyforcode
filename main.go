/*
   bountyforcode - A bounty system for bug and feature fixes
   Copyright (C) 2013 Abhi Yerra <abhi@berkeley.edu>

   This program is free software: you can redistribute it and/or modify
   it under the terms of the GNU Affero General Public License as
   published by the Free Software Foundation, either version 3 of the
   License, or (at your option) any later version.

   This program is distributed in the hope that it will be useful,
   but WITHOUT ANY WARRANTY; without even the implied warranty of
   MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
   GNU Affero General Public License for more details.

   You should have received a copy of the GNU Affero General Public License
   along with this program.  If not, see <http://www.gnu.org/licenses/>.
*/
package main

import (
	"bytes"
	"flag"
	"fmt"
	. "github.com/abhiyerra/bountyforcode/app/controllers"
	. "github.com/abhiyerra/bountyforcode/app/models"
	"github.com/abhiyerra/scalpy"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"text/template"
)

var (
	domain string
)

func initConfig() {
	flag.StringVar(&PostgresHost, "pghost", "", "the host for postgres")
	flag.StringVar(&PostgresDb, "dbname", "", "the db for postgres")
	flag.StringVar(&GithubClientId, "github_client_id", "", "github client id")
	flag.StringVar(&GithubClientSecret, "github_client_secret", "", "github client secret")
	flag.StringVar(&GithubRedirectUrl, "github_redirect_url", "", "github redirect url")
	flag.StringVar(&domain, "domain", "", "domain this is running on")

	flag.Parse()
}

// TODO Probably want to use this: https://github.com/codegangsta/martini-contrib/tree/master/render
type Page struct {
	Title string
	Body  string

	Content string

	ViewFile string
	Layout   string
}

func (p *Page) RenderView(w http.ResponseWriter) {
	t, err := template.ParseFiles(p.ViewFile)
	if err != nil {
		log.Printf("%v\n", err)
	}

	buffer := new(bytes.Buffer)

	t.Execute(buffer, p)
	p.Content = buffer.String()
	p.RenderLayout(w)
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

func AdminHandler(w http.ResponseWriter, r *http.Request) {
	//	session, _ := Store.Get(r, "user")
	// fmt.Fprintf(w, "%v", session.Values["UserId"])

	// button := coinbase.GetButton(&coinbase.ButtonRequest{
	// 	Name: "Abhi Yerra",
	// 	Type: "donation",
	// 	PriceString: "10.00",
	// 	PriceCurrencyIso: "USD",
	// })

	// fmt.Printf("%v", button.Response.Button.Code)
	fmt.Fprintf(w, "hi")
}

func RootHandler(w http.ResponseWriter, r *http.Request) {
	type RootPage struct {
		Page
		Discover string
	}

	issues := FindProjectIssues("abhiyerra")

	page := &RootPage{
		Page: Page{
			Title:    "Welcome",
			ViewFile: "views/root_index.tmpl",
		},
		Discover: renderDiscover(issues),
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

func DiscoverHandler(w http.ResponseWriter, r *http.Request) {

}

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
			ViewFile: "views/project_index.tmpl",
		},
		Discover: renderDiscover(issues),
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

func renderDiscover(issues []Issue) string {
	t, err := template.ParseFiles("views/_discover.tmpl")
	if err != nil {
		log.Printf("%v\n", err)
	}

	buffer := new(bytes.Buffer)

	t.Execute(buffer, issues)
	return buffer.String()
}

func main() {
	initConfig()
	InitDb()
	InitGithub()

	log.Printf("Server running on %s", domain)
	subdom := fmt.Sprintf("{subdomain:[a-z]+}.%s", domain)

	m := mux.NewRouter()
	m.HandleFunc("/", ProjectRootHandler).Host(subdom).Methods("GET")
	m.HandleFunc("/", RootHandler).Methods("GET")

	m.HandleFunc("/register/activate", RegisterAuthorizeHandler).Methods("GET") // TODO Should be authorize
	m.HandleFunc("/register", RegisterHandler).Methods("GET")

	m.HandleFunc("/bounties", CreateBountyHandler).Methods("POST")
	m.HandleFunc("/bounties/{id}", ShowBountyHandler).Methods("POST")

	m.HandleFunc("/search", RootHandler).Methods("GET")
	m.HandleFunc("/discover", DiscoverHandler).Methods("GET")
	m.HandleFunc("/pricing", RootHandler).Methods("GET")

	m.HandleFunc("/admin", AdminHandler).Methods("GET")

	http.Handle("/", m)
	http.ListenAndServe(":3000", nil)
}
