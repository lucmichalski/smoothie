package app

import (
	"fmt"
	"testing"
	"time"

	"github.com/tomocy/smoothie/domain"
)

func TestFetchPostsOfDrivers(t *testing.T) {
	expectedDate := time.Date(2019, 8, 13, 0, 0, 0, 0, time.Local)
	expecteds := domain.Posts{
		{ID: "1", Driver: "a", Text: "one", CreatedAt: expectedDate.Add(2 * time.Hour)},
		{ID: "1", Driver: "b", Text: "one", CreatedAt: expectedDate.Add(2 * time.Hour)},
		{ID: "2", Driver: "a", Text: "two", CreatedAt: expectedDate.Add(1 * time.Hour)},
		{ID: "2", Driver: "b", Text: "two", CreatedAt: expectedDate.Add(1 * time.Hour)},
		{ID: "3", Driver: "a", Text: "three", CreatedAt: expectedDate},
		{ID: "3", Driver: "b", Text: "three", CreatedAt: expectedDate},
	}
	u := newMockPostUsecase()
	actuals, err := u.FetchPostsOfDrivers("a", "b")
	if err != nil {
		t.Errorf("unexpected error by (*PostUsecase).FetchPostsOfDrivers: got %s, expect <nil>\n", err)
	}
	if err := assertPosts(actuals, expecteds); err != nil {
		t.Errorf("unexpected posts by (*PostUsecase).FetchPostsOfDrivers: %s\n", err)
	}
}

func newMockPostUsecase() *PostUsecase {
	ds := [...]string{"a", "b", "c"}
	repoA, repoB, repoC := newMock(ds[0]), newMock(ds[1]), newMock(ds[2])

	return NewPostUsecase(map[string]domain.PostRepo{
		ds[0]: repoA, ds[1]: repoB, ds[2]: repoC,
	})
}

func newMock(d string) *mock {
	date := time.Date(2019, 8, 13, 0, 0, 0, 0, time.Local)
	return &mock{
		ps: domain.Posts{
			{ID: "1", Driver: d, Text: "one", CreatedAt: date.Add(2 * time.Hour)},
			{ID: "2", Driver: d, Text: "two", CreatedAt: date.Add(1 * time.Hour)},
			{ID: "3", Driver: d, Text: "three", CreatedAt: date},
		},
	}
}

type mock struct {
	ps domain.Posts
}

func (m *mock) FetchPosts() (domain.Posts, error) {
	return m.ps, nil
}

func assertPosts(actuals, expecteds domain.Posts) error {
	if len(actuals) != len(expecteds) {
		return reportUnexpected("len of posts", len(actuals), len(expecteds))
	}
	for i, expected := range expecteds {
		if err := assertPost(actuals[i], expected); err != nil {
			return fmt.Errorf("unexpected posts[%d]: %s", i, err)
		}
	}

	return nil
}

func assertPost(actual, expected *domain.Post) error {
	if actual.ID != expected.ID {
		return reportUnexpected("id of post", actual.ID, expected.ID)
	}
	if actual.Driver != expected.Driver {
		return reportUnexpected("driver of post", actual.Driver, expected.Driver)
	}
	if actual.Text != expected.Text {
		return reportUnexpected("text of post", actual.Text, expected.Text)
	}
	if !actual.CreatedAt.Equal(expected.CreatedAt) {
		return reportUnexpected("created at of post", actual.CreatedAt, expected.CreatedAt)
	}

	return nil
}

func reportUnexpected(name string, actual, expected interface{}) error {
	return fmt.Errorf("unexpected %s: got %v, expect %v", name, actual, expected)
}