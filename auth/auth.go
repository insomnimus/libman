package auth

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/insomnimus/libman/config"
	"github.com/zmb3/spotify"
	"golang.org/x/oauth2"
)

const state = "xyz987"

type Creds struct {
	Client *spotify.Client
	User   *spotify.PrivateUser
	Token  *oauth2.Token
}

func authorize(auth *spotify.Authenticator) (*Creds, error) {
	ch := make(chan *spotify.Client)
	completeAuth := func(w http.ResponseWriter, r *http.Request) {
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
	http.HandleFunc("/callback", completeAuth)
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		// log.Println("Got request for:", r.URL.String())
	})
	go http.ListenAndServe(":8080", nil)

	url := auth.AuthURL(state)
	fmt.Println("Please log in to Spotify by visiting the following page in your browser:", url)
	// wait for auth to complete
	clt := <-ch
	// use the client to make calls that require authorization
	usr, err := clt.CurrentUser()
	if err != nil {
		return nil, fmt.Errorf("error auhorizing: %s", err)
	}
	fmt.Println("logged in!")

	token, err := clt.Token()
	if err != nil {
		return nil, fmt.Errorf("error fetching user info: %s", err)
	}

	return &Creds{
		Client: clt,
		User:   usr,
		Token:  token,
	}, nil
}

func Login(c *config.Config) (*Creds, error) {
	// all the permissions
	auth := spotify.NewAuthenticator(c.RedirectURI,
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
	auth.SetAuthInfo(c.ID, c.Secret)

	// check if the token is cached
	if c.CacheFile != "" {
		if _, err := os.Stat(c.CacheFile); err == nil {
			data, err := os.ReadFile(c.CacheFile)
			if err != nil {
				return nil, fmt.Errorf("could not read the cache file: %s", err)
			}
			var token oauth2.Token
			err = json.Unmarshal(data, &token)
			if err != nil {
				return nil, fmt.Errorf("could not unmarshal the cached token: %s", err)
			}
			return login(&auth, &token)
		}
	}

	// no cache file, prompt for access grant
	return authorize(&auth)
}

func login(auth *spotify.Authenticator, token *oauth2.Token) (*Creds, error) {
	clt := auth.NewClient(token)
	usr, err := clt.CurrentUser()
	if err != nil {
		return nil, fmt.Errorf("error fetching user details: %s", err)
	}
	fmt.Printf("welcome %s\n", usr.DisplayName)
	return &Creds{
		Client: &clt,
		User:   usr,
		Token:  token,
	}, nil
}
