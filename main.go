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
	. "bountyforcode/app"
	"bytes"
	"database/sql"
	"flag"
	"fmt"
	"github.com/abhiyerra/scalpy"
	//	"github.com/abhiyerra/coinbase"
	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
	"log"
	"net/http"
	"text/template"
)

var (
	postgresHost string
	postgresDb   string

	domain string
)

func dbConnect() *sql.DB {
	log.Printf("Connecting to DB: %s", postgresDb)

	db, err := sql.Open("postgres", fmt.Sprintf("dbname=%s sslmode=disable", postgresDb))
	if err != nil {
		log.Fatal(err)
	}

	return db
}

func initConfig() {
	flag.StringVar(&postgresHost, "pghost", "", "the host for postgres")
	flag.StringVar(&postgresDb, "dbname", "", "the db for postgres")
	flag.StringVar(&GithubClientId, "github_client_id", "", "github client id")
	flag.StringVar(&GithubClientSecret, "github_client_secret", "", "github client secret")
	flag.StringVar(&GithubRedirectUrl, "github_redirect_url", "", "github redirect url")
	flag.StringVar(&domain, "domain", "", "domain this is running on")

	flag.Parse()

	InitGithub()
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
		layout = "views/layouts/application.html"
	}

	t, err := template.ParseFiles(layout)
	if err != nil {
		log.Printf("%v\n", err)
	}

	t.Execute(w, p)
}

func RootHandler(w http.ResponseWriter, r *http.Request) {
	page := &Page{
		Title:    "Welcome",
		ViewFile: "views/bounties/_new.html",
	}

	page.RenderView(w)
}

func NewBountyHandler(w http.ResponseWriter, r *http.Request) {
	page := &Page{
		Title:    "New Bounty",
		ViewFile: "views/bounties/_confirm.html",
	}

	var vals struct {
		Repo string
	}

	issue := scalpy.ScalpUrl(r.FormValue("issue-url"))
	if issue == nil {
		vals.Repo = "Issue doesn't exist"
	} else {
		vals.Repo = issue.Repo
	}

	IssueFind("123", "123")

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
		ViewFile: "views/bounties/show.html",
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
		Projects []Project
	}

	page := &ProjectPage{
		Page: Page{
			Title:    project_identifier,
			ViewFile: "views/projects/show.html",
		},
	}

	page.Projects = FindProjects(project_identifier)
	fmt.Printf("%v, %v\n", page.Title, page.Projects)

	t, err := template.ParseFiles(page.ViewFile)
	if err != nil {
		log.Printf("%v\n", err)
	}

	buffer := new(bytes.Buffer)

	t.Execute(buffer, page)
	page.Content = buffer.String()

	page.RenderLayout(w)
}

func main() {
	initConfig()

	Db = dbConnect()

	fmt.Println(domain)
	subdom := fmt.Sprintf("{subdomain:[a-z]+}.%s", domain)
	fmt.Println(subdom)

	m := mux.NewRouter()
	m.HandleFunc("/", ProjectRootHandler).Host(subdom).Methods("GET")
	m.HandleFunc("/", RootHandler).Methods("GET")
	m.HandleFunc("/register/activate", RegisterAuthorizeHandler).Methods("GET") // TODO Should be authorize
	m.HandleFunc("/register", RegisterHandler).Methods("GET")

	m.HandleFunc("/bounties", NewBountyHandler).Methods("POST")
	m.HandleFunc("/bounties/{id}", ShowBountyHandler).Methods("POST")

	m.HandleFunc("/search", RootHandler).Methods("GET")
	m.HandleFunc("/discover", DiscoverHandler).Methods("GET")
	m.HandleFunc("/pricing", RootHandler).Methods("GET")

	http.Handle("/", m)
	http.ListenAndServe(":3000", nil)
}
