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
	_ "github.com/lib/pq"
	"log"
	"strings"
)

var (
	Db *sql.DB
)

type Issue struct {
	Id          uint
	Repo        string
	OriginalUrl string
	CreatedAt   string
	BountyTotal float32
}

type Project struct {
	Repo string
}

func FindProjects(identifier string) (projects []Project) {
	rows, err := Db.Query("SELECT repo FROM projects WHERE identifier = $1", strings.ToLower(identifier))
	if err != nil {
		log.Println(err)
	}

	if rows == nil {
		// TODO go get projects
		return nil
	}

	for rows.Next() {
		var project Project
		rows.Scan(&project.Repo)

		projects = append(projects, project)
	}

	return projects
}

func IssueFind(project, repo string) {
	// SELECT * FROM issues WHERE project = ? AND repo = ?
}

func NewIssue() *Issue {

	return &Issue{}
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
