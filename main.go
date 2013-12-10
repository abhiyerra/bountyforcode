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
	"database/sql"
	"flag"
	"fmt"
	"github.com/codegangsta/martini"
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
	Title string
	Body  []byte
}

func renderPage(w http.ResponseWriter, p *Page) {
	t, err := template.ParseFiles("views/layouts/application.html")

	if err != nil {
		log.Printf("%v\n", err)
	}

	t.Execute(w, p)
}

func RootPathHandler(w http.ResponseWriter, r *http.Request) {
	renderPage(w, &Page{
		Title: "Hello",
	})
}

func main() {
	parseEnv()

	m := martini.Classic()
	m.Get("/", RootPathHandler)
	m.Get("/bounties", RootPathHandler)
	m.Post("/bounties", RootPathHandler)
	m.Post("/search", RootPathHandler)
	m.Run()
}
