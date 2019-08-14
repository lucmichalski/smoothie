package runner

import (
	"fmt"
	"io"

	"github.com/tomocy/smoothie/domain"
)

type text struct {
	printed bool
}

func (t *text) PrintPosts(w io.Writer, ps domain.Posts) {
	vl := "----------"
	for _, p := range ps {
		if !t.printed {
			fmt.Fprintln(w, vl)
			t.printed = true
		}
		t.printPost(w, p)
		fmt.Fprintln(w, vl)
	}
}

func (t *text) printPost(w io.Writer, p *domain.Post) {
	fmt.Fprintf(w, "(%s) %s", p.Driver, p.User.Name)
	if p.User.Username != "" {
		fmt.Fprintf(w, " @%s", p.User.Username)
	}
	fmt.Fprintf(w, " %s\n%s\n", p.CreatedAt.Format("2006/01/02"), p.Text)
}
