package infra

import (
	"fmt"
	"net/http"
	"net/url"
	"os"
	"path/filepath"

	"github.com/garyburd/go-oauth/oauth"
)

func init() {
	must(createWorkspace())
}

func must(errs ...error) {
	for _, err := range errs {
		if err != nil {
			panic(fmt.Errorf("development error in infra: %s", err))
		}
	}
}

func createWorkspace() error {
	name := configFilename()
	if _, err := os.Stat(name); err == nil {
		return nil
	}

	dir := workspaceName()
	if err := os.RemoveAll(dir); err != nil {
		return err
	}
	if err := os.MkdirAll(dir, 0700); err != nil {
		return err
	}

	f, err := os.Create(name)
	if err != nil {
		return err
	}

	return f.Close()
}

func configFilename() string {
	return filepath.Join(workspaceName(), "config.json")
}

func workspaceName() string {
	return filepath.Join(os.Getenv("HOME"), ".smoothie")
}

type oauthReq struct {
	cred        *oauth.Credentials
	method, url string
	params      url.Values
}

func (r *oauthReq) do(client oauth.Client) (*http.Response, error) {
	if r.method != http.MethodGet {
		return client.Post(http.DefaultClient, r.cred, r.url, r.params)
	}

	return client.Get(http.DefaultClient, r.cred, r.url, r.params)
}