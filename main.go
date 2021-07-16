package main

import (
	"fmt"
	"github.com/zmb3/spotify"
	"golang.org/x/oauth2"
	"libman/control"
	"log"
	"net/http"
	"os"
)

var (
	redirectURI  = os.Getenv("LIBMAN_REDIRECT_URI")
	clientID     = os.Getenv("LIBMAN_ID")
	clientSecret = os.Getenv("LIBMAN_SECRET")
	state        = "xyz987"
	ch           = make(chan *spotify.Client)

	// all the permissions
	auth = spotify.NewAuthenticator(redirectURI,
		spotify.ScopeImageUpload,
		spotify.ScopePlaylistReadPrivate,
		spotify.ScopePlaylistModifyPublic,
		spotify.ScopePlaylistModifyPrivate,
		spotify.ScopePlaylistReadCollaborative,
		spotify.ScopeUserFollowModify,
		spotify.ScopeUserFollowRead,
		spotify.ScopeUserLibraryModify,
		spotify.ScopeUserLibraryRead,
		spotify.ScopeUserReadPrivate,
		spotify.ScopeUserReadEmail,
		spotify.ScopeUserReadCurrentlyPlaying,
		spotify.ScopeUserReadPlaybackState,
		spotify.ScopeUserModifyPlaybackState,
		spotify.ScopeUserReadRecentlyPlayed,
		spotify.ScopeUserTopRead,
		spotify.ScopeStreaming,
	)
)

func completeAuth(w http.ResponseWriter, r *http.Request) {
	tok, err := auth.Token(state, r)
	if err != nil {
		http.Error(w, "Couldn't get token", http.StatusForbidden)
		log.Fatal(err)
	}
	if st := r.FormValue("state"); st != state {
		http.NotFound(w, r)
		log.Fatalf("State mismatch: %s != %s\n", st, state)
	}
	// use the token to get an authenticated client
	client := auth.NewClient(tok)
	fmt.Fprintf(w, "Login Completed!")
	ch <- &client
}

func authorize() (*spotify.Client, *spotify.PrivateUser, *oauth2.Token) {

	http.HandleFunc("/callback", completeAuth)
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		//log.Println("Got request for:", r.URL.String())
	})
	go http.ListenAndServe(":8080", nil)

	url := auth.AuthURL(state)
	fmt.Println("Please log in to Spotify by visiting the following page in your browser:", url)
	// wait for auth to complete
	clt := <-ch

	// use the client to make calls that require authorization
	usr, err := clt.CurrentUser()
	if err != nil {
		log.Fatalf("error authorizing: %s\n", err)
	}
	fmt.Println("logged in!")

	token, err := clt.Token()
	if err != nil {
		log.Fatal(err)
	}
	return clt, usr, token
}

func main() {
	log.SetFlags(0)
	log.SetPrefix("")

	if clientID == "" {
		log.Fatal("error: the LIBMAN_ID env variable must be set")
	}
	if clientSecret == "" {
		log.Fatal("error: the LIBMAN_SECRET env variable must be set")
	}
	if redirectURI == "" {
		log.Fatal("error: the LIBMAN_REDIRECT_URI env variable must be set")
	}
	auth.SetAuthInfo(clientID, clientSecret)

	client, user, _ := authorize()
	control.SetHandlers(control.DefaultHandlers())
	control.Start(client, user, "@libman>")
}
