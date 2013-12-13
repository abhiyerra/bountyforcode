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
	"flag"
	"fmt"
	. "github.com/abhiyerra/bountyforcode/app/controllers"
	. "github.com/abhiyerra/bountyforcode/app/models"
	"github.com/abhiyerra/coinbase"
	"github.com/gorilla/mux"
	"log"
	"net/http"
)

var (
	domain string
)

func initConfig() {
	flag.StringVar(&PostgresHost, "pghost", "", "the host for postgres")
	flag.StringVar(&PostgresDb, "dbname", "", "the db for postgres")
	flag.StringVar(&GithubClientId, "github_client_id", "", "github client id")
	flag.StringVar(&GithubClientSecret, "github_client_secret", "", "github client secret")
	flag.StringVar(&GithubRedirectUrl, "github_redirect_url", "", "github redirect url")
	flag.StringVar(&domain, "domain", "", "domain this is running on")
	flag.StringVar(&SecretStoreKey, "secret_store_key", "", "Secret session store key")
	flag.StringVar(&coinbase.CoinbaseApiKey, "coinbase_api_key", "", "Coinbase API Key")
	flag.StringVar(&coinbase.CoinbaseCallbackSecret, "coinbase_callback_secret", "", "Coinbase Callback Secret")

	flag.Parse()
}

func main() {
	initConfig()

	InitDb()
	InitSessionStore()

	InitGithub()

	log.Printf("Server running on %s", domain)

	m := mux.NewRouter()

	subdom := fmt.Sprintf("{subdomain:[a-z]+}.%s", domain)
	m.HandleFunc("/", ProjectRootHandler).Host(subdom).Methods("GET")

	m.HandleFunc("/", RootHandler).Methods("GET")

	m.HandleFunc("/register/activate", RegisterAuthorizeHandler).Methods("GET") // TODO Should be authorize
	m.HandleFunc("/register", RegisterHandler).Methods("GET")

	m.HandleFunc("/bounties", CreateBountyHandler).Methods("POST")
	m.HandleFunc("/bounties/{id}", ShowBountyHandler).Methods("POST")

	if coinbase.CoinbaseCallbackSecret == "" {
		log.Fatal("CoinbaseCallbackSecret can't be empty")
	}
	coinbase_callback_path := fmt.Sprintf("/coinbase/callback/%s", coinbase.CoinbaseCallbackSecret)
	m.HandleFunc(coinbase_callback_path, CoinbaseCallbackHandler).Methods("POST")
	m.HandleFunc("/pricing", RootHandler).Methods("GET")

	m.HandleFunc("/admin", AdminHandler).Methods("GET")

	http.Handle("/", m)
	http.ListenAndServe(":3000", nil)
}
