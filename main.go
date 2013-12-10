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
	"code.google.com/p/goauth2/oauth"
	"database/sql"
	"flag"
	"fmt"
	"github.com/abhiyerra/scalpy"
	"github.com/codegangsta/martini"
	_ "github.com/lib/pq"
	"log"
	"net/http"
	"text/template"
)

var (
	postgresHost string
	postgresDb   string

	githubClientId     string
	githubClientSecret string
	githubRedirectUrl  string
	githubAuthUrl      = "https://github.com/login/oauth/authorize"
	githubTokenUrl     = "https://github.com/login/oauth/access_token"
	githubConfig       *oauth.Config

	db *sql.DB
)

func dbConnect() *sql.DB {
	db, err := sql.Open("postgres", fmt.Sprintf("dbname=%s sslmode=disable", postgresDb))
	if err != nil {
		log.Fatal(err)
	}

	return db
}

func initConfig() {
	flag.StringVar(&postgresHost, "pghost", "", "the host for postgres")
	flag.StringVar(&postgresDb, "dbname", "", "the db for postgres")
	flag.StringVar(&githubClientId, "github_client_id", "", "github client id")
	flag.StringVar(&githubClientSecret, "github_client_secret", "", "github client secret")
	flag.StringVar(&githubRedirectUrl, "github_redirect_url", "", "github redirect url")

	flag.Parse()

	githubConfig = &oauth.Config{
		ClientId:     githubClientId,
		ClientSecret: githubClientSecret,
		AuthURL:      githubAuthUrl,
		TokenURL:     githubTokenUrl,
	}

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

func RegisterHandler(w http.ResponseWriter, r *http.Request) {
	// Make sure token still exists.

	url := githubConfig.AuthCodeURL("")
	http.Redirect(w, r, url, 302)
}

func RegisterAuthorizeHandler(w http.ResponseWriter, r *http.Request) {
	code := r.FormValue("code")
	if code == "" {
		fmt.Fprintf(w, "Failed to login")
		return
	}

	transport := &oauth.Transport{Config: githubConfig}
	token, _ := transport.Exchange(code)

	fmt.Fprintf(w, token.AccessToken)
}

func DiscoverHandler(w http.ResponseWriter, r *http.Request) {

}

func main() {
	initConfig()

	db = dbConnect()

	fmt.Println(githubClientId)

	m := martini.Classic()
	m.Get("/", RootHandler)
	m.Get("/register/activate", RegisterAuthorizeHandler) // TODO Should be authorize
	m.Get("/register", RegisterHandler)

	m.Post("/bounties", NewBountyHandler)
	m.Get("/bounties/:id", ShowBountyHandler)

	//	m.Post("/search", RootHandler)
	m.Post("/discover", RootHandler)
	m.Post("/pricing", RootHandler)
	m.Run()
}
