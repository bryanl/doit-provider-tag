package godoext

import (
	"github.com/digitalocean/godo"
	"golang.org/x/oauth2"
)

// Client wraps godo.Client to add tagging endpoints.
type Client struct {
	*godo.Client
	Tags TagService
}

// New creates a Client.
func New(token string) *Client {
	ts := &tokenSource{
		AccessToken: token,
	}

	oauthClient := oauth2.NewClient(oauth2.NoContext, ts)
	godoClient := godo.NewClient(oauthClient)

	c := &Client{Client: godoClient}
	tagServ := &tagsService{client: c}

	c.Tags = tagServ

	return c
}

type tokenSource struct {
	AccessToken string
}

func (t *tokenSource) Token() (*oauth2.Token, error) {
	token := &oauth2.Token{
		AccessToken: t.AccessToken,
	}
	return token, nil
}
