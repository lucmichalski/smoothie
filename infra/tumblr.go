package infra

import (
	"context"
	"fmt"
	"net/http"

	"github.com/garyburd/go-oauth/oauth"
	"github.com/tomocy/deverr"

	"github.com/tomocy/smoothie/domain"
)

func NewTumblr(id, secret string) *Tumblr {
	return &Tumblr{
		oauth: oauthManager{
			client: oauth.Client{
				TemporaryCredentialRequestURI: "https://www.tumblr.com/oauth/request_token",
				ResourceOwnerAuthorizationURI: "https://www.tumblr.com/oauth/authorize",
				TokenRequestURI:               "https://www.tumblr.com/oauth/access_token",
				Credentials: oauth.Credentials{
					Token: id, Secret: secret,
				},
			},
		},
	}
}

type Tumblr struct {
	oauth oauthManager
}

func (t *Tumblr) StreamPosts(ctx context.Context) (<-chan domain.Posts, <-chan error) {
	psCh, errCh := make(chan domain.Posts), make(chan error)
	go func() {
		defer func() {
			close(psCh)
			close(errCh)
		}()
		select {
		case <-ctx.Done():
			return
		default:
			errCh <- deverr.NotImplemented
		}
	}()

	return psCh, errCh
}

func (t *Tumblr) FetchPosts() (domain.Posts, error) {
	_, err := t.retreiveAuthorization()
	if err != nil {
		return nil, err
	}

	return nil, deverr.NotImplemented
}

func (t *Tumblr) retreiveAuthorization() (*oauth.Credentials, error) {
	url, err := t.authorizationURL()
	if err != nil {
		return nil, err
	}
	fmt.Printf("open this url: %s\n", url)

	return t.handleAuthorizationRedirect()
}

func (t *Tumblr) authorizationURL() (string, error) {
	temp, err := t.oauth.client.RequestTemporaryCredentials(http.DefaultClient, "", nil)
	if err != nil {
		return "", err
	}
	t.oauth.temp = temp

	return t.oauth.client.AuthorizationURL(temp, nil), nil
}

func (t *Tumblr) handleAuthorizationRedirect() (*oauth.Credentials, error) {
	return t.oauth.handleRedirect(context.Background(), "/smoothie/tumblr/authorization")
}
