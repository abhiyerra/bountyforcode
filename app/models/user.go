package bountyforcode

import (
	"fmt"
	"log"
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
