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
		posts, err := listPost.Do()

		assert.Equal(t, err, nil, "Something error %v", err)
		assert.Equal(t, len(posts), 0, "Post should not return anything")
	})

	t.Run("default post list", func(t *testing.T) {
		listPost := client.Posts().List()
		posts, err := listPost.Do()

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
	posts, err := listAPI.Do()

	postCountBefore := len(posts)

	assert.Equal(t, err, nil, "Something error %v", err)

	_, err = postAPI.Create(gowprest.PostData{
		Title:   faker.Sentence(),
		Content: faker.Paragraph(),
		Excerpt: faker.Sentence(),
		Status:  gowprest.StatusPublished,
	}).Do()

	assert.Equal(t, err, nil, err)
	posts, err = listAPI.Do()

	postCountAfter := len(posts)

	assert.LessOrEqual(t, postCountBefore, postCountAfter,
		"New post count should greater than or equal to the previous count. Expected %d, got %d",
		postCountBefore, postCountAfter)
}

func TestRetrievePost(t *testing.T) {
	client := gowprest.NewClient(blogUrl)
	defer client.Close()

	listAPI := client.Posts().List()
	posts, err := listAPI.Do()
	assert.Equal(t, err, nil, "Something error %v", err)

	singlePost, err := client.Posts().Retrieve(posts[0].ID).Do()
	assert.Equal(t, err, nil, "Something error %v", err)
	assert.Equal(t, posts[0].Title.Rendered, singlePost.Title.Rendered)
}

func TestUpdatePost(t *testing.T) {
	client := gowprest.NewClient(blogUrl).
		WithBasicAuth(
			os.Getenv("BLOG_USERNAME"),
			os.Getenv("BLOG_APP_PASSWORD"),
		)
	defer client.Close()

	listAPI := client.Posts().List()
	posts, err := listAPI.Do()
	assert.Equal(t, err, nil, "Something error %v", err)

	selectedPost := posts[0]
	postID := selectedPost.ID
	updateAPI := client.Posts().Update(gowprest.PostData{
		ID:      postID,
		Title:   faker.Sentence(),
		Content: faker.Paragraph(),
		Excerpt: faker.Sentence(),
		Status:  gowprest.StatusPublished,
	})

	updatedPost, err := updateAPI.Do()
	assert.Equal(t, err, nil, "Something error %v", err)
	assert.Equal(t, postID, updatedPost.ID)
	assert.NotEqual(t, selectedPost.Title.Rendered, updatedPost.Title.Rendered)
	assert.NotEqual(t, selectedPost.Content.Rendered, updatedPost.Content.Rendered)

	updateAPI = client.Posts().Update(gowprest.PostData{
		ID:     postID,
		Status: gowprest.StatusDraft,
	})

	updatedPost, err = updateAPI.Do()
	assert.Equal(t, err, nil, "Something error %v", err)
	assert.Equal(t, postID, updatedPost.ID)
	assert.NotEqual(t, selectedPost.Status, updatedPost.Status)
}
