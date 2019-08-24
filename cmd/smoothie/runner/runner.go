package runner

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"io"
	"os"
	"os/signal"
	"sort"
	"strings"
	"syscall"

	"github.com/joho/godotenv"

	"github.com/tomocy/smoothie/app"
	"github.com/tomocy/smoothie/domain"
	"github.com/tomocy/smoothie/infra"
)

func init() {
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage of %s: [optinos] drivers...\n", os.Args[0])

		flag.PrintDefaults()
	}
}

func New() Runner {
	cnf, err := parseConfig()
	if err != nil {
		return &Help{
			err: err,
		}
	}

	godotenv.Load()

	return &Continue{
		cnf:       cnf,
		presenter: newPresenter(cnf.mode, cnf.format),
	}
}

type Runner interface {
	Run() error
}

func parseConfig() (config, error) {
	m, f := flag.String("m", modeCLI, "name of mode"), flag.String("f", formatText, "format")
	s := flag.Bool("s", false, "enable streaming")
	flag.Parse()

	return config{
		mode: *m, format: *f,
		isStreaming: *s,
		drivers:     flag.Args(),
	}, nil
}

type config struct {
	mode, format string
	isStreaming  bool
	drivers      []string
}

const (
	modeCLI  = "cli"
	modeHTTP = "http"

	formatText  = "text"
	formatColor = "color"
	formatHTML  = "html"
	formatJSON  = "json"
)

func newPresenter(mode, format string) presenter {
	switch mode {
	case modeCLI:
		return &cli{
			printer: newPrinter(format),
		}
	case modeHTTP:
		return &http{
			printer: newPrinter(format),
		}
	default:
		return nil
	}
}

type presenter interface {
	ShowPosts(domain.Posts)
}

func newPrinter(format string) printer {
	switch format {
	case formatText:
		return new(text)
	case formatColor:
		return new(color)
	case formatHTML:
		return new(html)
	case formatJSON:
		return new(json)
	default:
		return nil
	}
}

type printer interface {
	PrintPosts(io.Writer, domain.Posts)
}

type Continue struct {
	cnf       config
	presenter presenter
}

func (c *Continue) Run() error {
	if c.cnf.isStreaming {
		return c.streamPostsOfDrivers()
	}

	return c.fetchPostsOfDrivers()
}

func (c *Continue) streamPostsOfDrivers() error {
	u := newPostUsecase()
	ctx, cancelFn := context.WithCancel(context.Background())
	psCh, errCh := u.StreamPostsOfDrivers(ctx, c.cnf.drivers...)
	sigCh := make(chan os.Signal)
	defer close(sigCh)
	signal.Notify(sigCh, syscall.SIGINT)
	for {
		select {
		case ps := <-psCh:
			c.presenter.ShowPosts(ps)
		case err := <-errCh:
			cancelFn()
			return err
		case sig := <-sigCh:
			cancelFn()
			return errors.New(sig.String())
		}
	}
}

func (c *Continue) fetchPostsOfDrivers() error {
	u := newPostUsecase()
	ps, err := u.FetchPostsOfDrivers(c.cnf.drivers...)
	if err != nil {
		return err
	}

	c.presenter.ShowPosts(ps)

	return nil
}

func orderPostsByOldest(ps domain.Posts) domain.Posts {
	ordered := make(domain.Posts, len(ps))
	copy(ordered, ps)
	sort.Slice(ordered, func(i, j int) bool {
		return ordered[i].CreatedAt.Before(ordered[j].CreatedAt)
	})

	return ordered
}

type Help struct {
	err error
}

func (h *Help) Run() error {
	flag.Usage()
	return h.err
}

func newPostUsecase() *app.PostUsecase {
	rs := make(map[string]domain.PostRepo)
	rs["gmail"] = infra.NewGmail("840774204831-aqpeh6dt6d9v1jbks7iilgpkcfliigtd.apps.googleusercontent.com", "fVHGfRA8I5iubiIJPsI2SpZw")
	rs["tumblr"] = infra.NewTumblr(idAndSecret("tumblr"))
	rs["twitter"] = infra.NewTwitter(idAndSecret("twitter"))
	rs["reddit"] = infra.NewReddit(idAndSecret("reddit"))

	return app.NewPostUsecase(rs)
}

func idAndSecret(driver string) (string, string) {
	return os.Getenv(fmt.Sprintf("%s_CLIENT_ID", strings.ToUpper(driver))), os.Getenv(fmt.Sprintf("%s_CLIENT_SECRET", strings.ToUpper(driver)))
}
