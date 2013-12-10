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
	_ "github.com/lib/pq"
	"github.com/codegangsta/martini"
	"log"
	"fmt"
	"flag"
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

func RootPathHandler(w http.ResponseWriter, r *http.Request) {
	t, _ := template.New("foo").Parse(`{{define "T"}}Hello, {{.}}!{{end}}`)
	t.ExecuteTemplate(w, "T", "<script>alert('you have been pwned')</script>")
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
