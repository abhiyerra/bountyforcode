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
	_ "github.com/lib/pq"
	"log"
)

var (
	Db *sql.DB
)

type User struct {
	Id          string
	AccessToken string
}

func NewUser(access_token string) (u *User) {
	if access_token == "" {
		log.Printf("No access_token to add to db\n")
		return nil
	}

	u = &User{
		AccessToken: access_token,
	}

	err := Db.QueryRow("INSERT INTO users (github_access_token) VALUES ($1) RETURNING id;", u.AccessToken).Scan(&u.Id)
	if err != nil {
		log.Printf("Couldn't add access_token %s", u.AccessToken)
		log.Fatal(err)
		return nil
	}

	fmt.Printf("asd id %s", u.Id)
	return u
}

func (i *Issue) Bounties() []Bounty {
	// SELECT * FROM bounties WHERE issue_id = ?
	return nil
}

type Bounty struct {
}

func NewBounty() {

}

func BountiesOpen() {
	// SELECT * FROM issues INNER JOIN bounties ON bounties.issue_id = issue.id WHERE bounty_state in ('open', 'paid')
}

func BountiesRecent() {
	// SELECT * FROM bounties ORDER BY created_at DESC limit 10
}
