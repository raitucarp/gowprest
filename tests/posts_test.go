package tests

import (
	"os"
	"testing"

	"github.com/go-faker/faker/v4"
	_ "github.com/joho/godotenv/autoload"
	"github.com/raitucarp/gowprest"
	"github.com/stretchr/testify/assert"
)

var blogUrl = os.Getenv("BLOG_URL")

func TestListPosts(t *testing.T) {
	client := gowprest.NewClient(blogUrl)
	defer client.Close()

	t.Run("not found search", func(t *testing.T) {
		listPost := client.Posts().List()
		listPost = listPost.Search("test")
		posts, err := listPost.Get()

		assert.Equal(t, err, nil, "Something error %v", err)
		assert.Equal(t, len(posts), 0, "Post should not return anything")
	})

	t.Run("default post list", func(t *testing.T) {
		listPost := client.Posts().List()
		posts, err := listPost.Get()

		assert.Equal(t, err, nil, "Something error %v", err)
		assert.Greater(t, len(posts), 0, "Post should return at least one post")
	})

}

func TestCreatePost(t *testing.T) {
	client := gowprest.NewClient(blogUrl).
		WithBasicAuth(
			os.Getenv("BLOG_USERNAME"),
			os.Getenv("BLOG_APP_PASSWORD"),
		)

	defer client.Close()

	postAPI := client.Posts()

	listAPI := postAPI.List()
	posts, err := listAPI.Get()

	postCountBefore := len(posts)

	assert.Equal(t, err, nil, "Something error %v", err)

	_, err = postAPI.Create(gowprest.NewPost{
		Title:   faker.Sentence(),
		Content: faker.Paragraph(),
		Excerpt: faker.Sentence(),
		Status:  gowprest.StatusPublished,
	}).Post()

	assert.Equal(t, err, nil, err)
	posts, err = listAPI.Get()

	postCountAfter := len(posts)

	assert.LessOrEqual(t, postCountBefore, postCountAfter,
		"New post count should greater than or equal to the previous count. Expected %d, got %d",
		postCountBefore, postCountAfter)
}
