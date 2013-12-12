package bountyforcode

import (
	"database/sql"
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

	err := Db.QueryRow(`SELECT id FROM users WHERE github_access_token = $1`, access_token).Scan(&u.Id)
	switch {
	case err == sql.ErrNoRows:
		err := Db.QueryRow("INSERT INTO users (github_access_token) VALUES ($1) RETURNING id;", u.AccessToken).Scan(&u.Id)
		if err != nil {
			log.Printf("Couldn't add access_token %s", u.AccessToken)
			log.Fatal(err)
			return nil
		}
	case err != nil:
		log.Fatal(err)
	default:
		fmt.Printf("UserId is %s\n", u.Id)
	}

	return u
}
