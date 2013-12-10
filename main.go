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
	"database/sql"
	"flag"
	"fmt"
	"github.com/codegangsta/martini"
	//	"code.google.com/p/goauth2/oauth" http://code.google.com/p/goauth2/source/browse/oauth/example/oauthreq.go
	"github.com/abhiyerra/scalpy"
	_ "github.com/lib/pq"
	"log"
	"net/http"
	"text/template"
)

var (
	postgresHost string
	postgresDb   string
)

func dbConnect() *sql.DB {
	db, err := sql.Open("postgres", fmt.Sprintf("dbname=%s sslmode=disable"))
	if err != nil {
		log.Println(err)
	}

	return db
}

func parseEnv() {
	flag.StringVar(&postgresHost, "pghost", "", "the host for postgres")
	flag.StringVar(&postgresDb, "dbname", "", "the db for postgres")

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
		layout = "views/layouts/application.html"
	}

	t, err := template.ParseFiles(layout)
	if err != nil {
		log.Printf("%v\n", err)
	}

	t.Execute(w, p)
}

func RootPathHandler(w http.ResponseWriter, r *http.Request) {
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

func main() {
	parseEnv()

	m := martini.Classic()
	m.Get("/", RootPathHandler)
	m.Post("/bounties", NewBountyHandler)
	m.Get("/bounties/:id", ShowBountyHandler)

	m.Get("/register", RootPathHandler)
	m.Post("/search", RootPathHandler)
	m.Run()
}

/*

CREATE TYPE hosting_provider AS ENUM ('github');

create table users (
   id serial primary key,
   username varchar(255),
   password varchar(255),
   salt varchar(255),
   github_identifier varchar(255)
   created_at datetime default now()
   created_at datetime default now()
)

create table issues (
   id serial primary key,
   original_url varchar(255),
   hoster hosting_provider,
   repo varchar(255)
   created_at datetime default now()
   created_at datetime default now()
);

CREATE TYPE bounty_state as ENUM ('open', 'paid', 'closed', 'cancelled');

create table bounties (
   user_id serial references users(id),
   issue_id serial references issues(id),
   amount float,
   transaction_status bounty_state,
   expires_at datetime
   created_at datetime default now()
   created_at datetime default now()
);

*/
