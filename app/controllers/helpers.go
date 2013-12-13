package bountyforcode

import (
	"github.com/gorilla/sessions"
	"log"
	"fmt"
	"net/http"
)

var (
	HtmlDir        string   
	Store          *sessions.CookieStore
	SecretStoreKey string
)

func InitSessionStore() {
	if SecretStoreKey == "" {
		log.Fatal("SecretStoreKey isn't set")
	}

	Store = sessions.NewCookieStore([]byte("something-very-secret"))
}

func GetSessionUserId(r *http.Request) string {
	session, _ := Store.Get(r, "user")
	user_id := session.Values["UserId"]

	if str, ok := user_id.(string); ok {
		return str
	} else {
		return ""
	}
}

func SetSessionUserId(w http.ResponseWriter, r *http.Request, user_id string) {
	session, _ := Store.Get(r, "user")
	session.Values["UserId"] = user_id
	session.Save(r, w)
}

func GetView(view string) string {
	fmt.Println(HtmlDir)
	return fmt.Sprintf("%s/%s", HtmlDir, view)
}
