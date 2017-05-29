package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	"github.com/zmb3/spotify"
	"golang.org/x/oauth2"
)

const (
	redirectURI = "http://localhost:8080/callback"
	state       = "spotiplaylistanalyzer"
	tokenFile   = "/tmp/spotify-playlist-token"
)

var opts struct {
	listenPort int
}

var (
	scopes = []string{
		spotify.ScopePlaylistReadPrivate,
		spotify.ScopePlaylistReadCollaborative,
		spotify.ScopePlaylistModifyPublic,
	}
	auth     = spotify.NewAuthenticator(redirectURI, scopes...)
	clientCh = make(chan *spotify.Client)
)

func main() {
	flag.IntVar(&opts.listenPort, "port", 8080, "Port for HTTP connections")
	flag.Parse()

	http.HandleFunc("/callback", completeAuth)
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "To authenticate visit: %s", auth.AuthURL(state))
	})
	listenURL := fmt.Sprintf(":%d", opts.listenPort)
	go http.ListenAndServe(listenURL, nil)
	log.Printf("Listening on %s", listenURL)

	// Wait for the authentication to finish, from cache or callback.
	go readTokenCache(tokenFile)
	c := <-clientCh

	u, err := c.CurrentUser()
	if err != nil {
		log.Fatalf("could not get current user: %s", err)
		return
	}

	ps := newPlaylistStats(c)
	ps.analyzeUser(u.ID)
	ps.printStats()
}

func completeAuth(w http.ResponseWriter, r *http.Request) {
	token, err := auth.Token(state, r)
	if err != nil {
		http.Error(w, "Couldn't get token in auth callback", http.StatusBadRequest)
		log.Fatal(err)
	}
	if st := r.FormValue("state"); st != state {
		log.Fatalf("State mismatch: %s != %s", st, state)
	}

	// Use the token to get an authenticated client.
	client := auth.NewClient(token)
	fmt.Fprint(w, "Login complete, cached token to file.")
	writeTokenCache(tokenFile, token)

	clientCh <- &client
}

func readTokenCache(filename string) error {
	f, err := os.Open(filename)
	if err != nil {
		return err
	}
	b, err := ioutil.ReadAll(f)
	if err != nil {
		return err
	}

	var token *oauth2.Token
	if err := json.Unmarshal(b, &token); err != nil {
		return err
	}

	client := auth.NewClient(token)
	clientCh <- &client
	return nil
}

func writeTokenCache(filename string, token *oauth2.Token) {
	f, err := os.OpenFile(filename, os.O_RDWR|os.O_CREATE, 0755)
	if err != nil {
		log.Fatalf("Unable to open token cache file: %s", err)
		return
	}

	b, err := json.Marshal(token)
	if err != nil {
		log.Fatalf("Unable to marshal oauth token: %s", err)
		return
	}

	if _, err := f.Write(b); err != nil {
		log.Fatalf("Could not write to token file: %s", err)
	}

}
