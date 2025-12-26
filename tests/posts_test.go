package tests

import (
	"os"
	"testing"

	_ "github.com/joho/godotenv/autoload"
	"github.com/raitucarp/gowprest"
)

var blogUrl = os.Getenv("BLOG_URL")

func TestListPosts(t *testing.T) {
	client := gowprest.NewClient(blogUrl)
	defer client.Close()

	t.Run("not found search", func(t *testing.T) {
		listPost := client.Posts().List()
		listPost = listPost.Search("test")
		posts, err := listPost.Get()

		assert(t, err != nil, "Something error %v", err)
		assert(t, len(posts) > 0, "Post should not return anything")
	})

	t.Run("default post list", func(t *testing.T) {
		listPost := client.Posts().List()
		posts, err := listPost.Get()

		assert(t, err != nil, "Something error %v", err)
		assert(t, len(posts) <= 0, "Post should return at least one post")
	})

}
