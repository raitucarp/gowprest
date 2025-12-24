package tests

import (
	"os"
	"testing"

	_ "github.com/joho/godotenv/autoload"
	"github.com/raitucarp/gowprest"
)

func TestListPosts(t *testing.T) {
	blogUrl := os.Getenv("BLOG_URL")

	t.Run("not found search", func(t *testing.T) {
		client := gowprest.NewClient(blogUrl)
		listPost := client.Posts().List()
		listPost = listPost.Search("test")
		posts, err := listPost.Get()

		if err != nil {
			t.Errorf("Something error %v", err)
		}

		if len(posts) > 0 {
			t.Errorf("Post should not return anything")
		}
	})

	t.Run("default post list", func(t *testing.T) {
		client := gowprest.NewClient(blogUrl)
		listPost := client.Posts().List()
		posts, err := listPost.Get()

		if err != nil {
			t.Errorf("Something error %v", err)
		}

		if len(posts) <= 0 {
			t.Errorf("Post should return at least one post")
		}
	})

}
