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
	_ "github.com/lib/pq"
	"log"
	"text/template"
	//	"github.com/abhiyerra/scalpy"
	"net/http"
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

type Page struct {
	Title   string
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

	page.RenderView(w)
}

func main() {
	parseEnv()

	m := martini.Classic()
	m.Get("/", RootPathHandler)
	m.Get("/register", RootPathHandler)
	m.Get("/bounties", RootPathHandler)
	m.Post("/bounties", NewBountyHandler)
	m.Post("/search", RootPathHandler)
	m.Run()
}
