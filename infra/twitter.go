package infra

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"path/filepath"

	"github.com/garyburd/go-oauth/oauth"
	"github.com/tomocy/smoothie/domain"
	"github.com/tomocy/smoothie/infra/twitter"
)

func NewTwitter(id, secret string) *Twitter {
	return &Twitter{
		oauthClient: oauth.Client{
			TemporaryCredentialRequestURI: "https://api.twitter.com/oauth/request_token",
			ResourceOwnerAuthorizationURI: "https://api.twitter.com/oauth/authenticate",
			TokenRequestURI:               "https://api.twitter.com/oauth/access_token",
			Credentials: oauth.Credentials{
				Token:  id,
				Secret: secret,
			},
		},
	}
}

type Twitter struct {
	oauthClient oauth.Client
}

func (t *Twitter) FetchPosts() (domain.Posts, error) {
	ts, err := t.fetchTweets(t.endpoint("/statuses/home_timeline.json"), nil)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch posts: %s", err)
	}

	return ts.Adapt(), nil
}

func (t *Twitter) fetchTweets(dst string, params url.Values) (twitter.Tweets, error) {
	cred, err := t.retreiveAuthorization()
	if err != nil {
		return nil, err
	}

	assured := params
	if assured == nil {
		assured = make(url.Values)
	}
	assured.Set("tweet_mode", "extended")
	var ts twitter.Tweets
	if err := t.do(oauthReq{
		cred: cred, method: http.MethodGet, url: dst, params: assured,
	}, &ts); err != nil {
		return nil, err
	}

	return ts, nil
}

func (t *Twitter) retreiveAuthorization() (*oauth.Credentials, error) {
	temp, err := t.oauthClient.RequestTemporaryCredentials(http.DefaultClient, "", nil)
	if err != nil {
		return nil, err
	}

	return t.requestClientAuthorization(temp)
}

func (t *Twitter) requestClientAuthorization(temp *oauth.Credentials) (*oauth.Credentials, error) {
	url := t.oauthClient.AuthorizationURL(temp, nil)
	fmt.Println("open this url: ", url)

	fmt.Print("PIN: ")
	var pin string
	fmt.Scanln(&pin)

	token, _, err := t.oauthClient.RequestToken(http.DefaultClient, temp, pin)

	return token, err
}

func (t *Twitter) do(r oauthReq, dst interface{}) error {
	resp, err := r.do(t.oauthClient)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	return json.NewDecoder(resp.Body).Decode(dst)
}

func (t *Twitter) endpoint(ps ...string) string {
	ss := append([]string{"https://api.twitter.com/1.1"}, ps...)
	return filepath.Join(ss...)
}