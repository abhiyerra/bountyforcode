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
package bountyforcode

import (
	"database/sql"
	"fmt"
	"github.com/coopernurse/gorp"
	_ "github.com/lib/pq"
	"log"
)

var (
	PostgresHost string
	PostgresDb   string

	Db    *sql.DB
	DbMap *gorp.DbMap
)

func InitDb() {
	log.Printf("Connecting to DB: %s", PostgresDb)

	var err error
	Db, err = sql.Open("postgres", fmt.Sprintf("dbname=%s sslmode=disable", PostgresDb))
	if err != nil {
		log.Fatal(err)
	}

	DbMap := &gorp.DbMap{Db: Db, Dialect: gorp.PostgresDialect{}}
	DbMap.AddTableWithName(Bounty{}, "bounties").SetKeys(true, "Id")
}

func CloseDb() {
	DbMap.Db.Close()
}
