package bountyforcode

import (
	"github.com/gorilla/sessions"
	"log"
)

var (
	Store          *sessions.CookieStore
	SecretStoreKey string
)

func InitSessionStore() {
	if SecretStoreKey == "" {
		log.Fatal("SecretStoreKey isn't set")
	}

	Store = sessions.NewCookieStore([]byte("something-very-secret"))
}
