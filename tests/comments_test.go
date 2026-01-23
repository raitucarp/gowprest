package tests

import (
	"os"
	"testing"

	"github.com/go-faker/faker/v4"
	_ "github.com/joho/godotenv/autoload"
	"github.com/raitucarp/gowprest"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestComments(t *testing.T) {
	blogUrl := os.Getenv("BLOG_URL")
	username := os.Getenv("BLOG_USERNAME")
	password := os.Getenv("BLOG_APP_PASSWORD")

	if blogUrl == "" || username == "" || password == "" {
		t.Skip("Skipping test; BLOG_URL, BLOG_USERNAME, or BLOG_APP_PASSWORD not set")
	}

	client := gowprest.NewClient(blogUrl).WithBasicAuth(username, password)
	defer client.Close()

	commentsAPI := client.Comments()
	postsAPI := client.Posts()

	// 1. Create a post to comment on
	post, err := postsAPI.Create(gowprest.PostData{
		Title:   faker.Sentence(),
		Content: faker.Paragraph(),
		Status:  gowprest.StatusPublished,
	}).Do()
	require.NoError(t, err)
	// defer postsAPI.Delete(post.ID).Force().Do()

	// 2. Create a comment
	commentContent := faker.Paragraph()
	comment, err := commentsAPI.Create(gowprest.CommentData{
		Post:        post.ID,
		Content:     commentContent,
		AuthorName:  "Test User",
		AuthorEmail: "test@example.com",
	}).Do()
	require.NoError(t, err)
	assert.Equal(t, post.ID, comment.Post)
	assert.Contains(t, comment.Content.Rendered, commentContent)
	// defer commentsAPI.Delete(comment.ID).Force().Do()

	// 3. Retrieve the comment
	retrievedComment, err := commentsAPI.Retrieve(comment.ID).ContextEdit().Do()
	require.NoError(t, err)
	assert.Equal(t, comment.ID, retrievedComment.ID)

	// 4. Update the comment
	newContent := faker.Paragraph()
	updatedComment, err := commentsAPI.Update(gowprest.CommentData{
		ID:      comment.ID,
		Content: newContent,
	}).Do()
	require.NoError(t, err)
	assert.Contains(t, updatedComment.Content.Rendered, newContent)

	comments, err := commentsAPI.List().Post(post.ID).Status("any").ContextEdit().Do()
	require.NoError(t, err)

	assert.GreaterOrEqual(t, len(comments), 1, "Should have at least one comment for the post")

	found := false
	for _, c := range comments {
		if c.ID == comment.ID {
			found = true
			break
		}
	}
	assert.True(t, found, "Should find the created comment in the list")

	// 6. Delete the comment
	deletedComment, err := commentsAPI.Delete(comment.ID).Force().Do()

	require.NoError(t, err)
	assert.Equal(t, comment.ID, deletedComment.Previous.ID)

	// Verify deletion
	_, err = commentsAPI.Retrieve(comment.ID).Do()
	assert.Error(t, err, "Should error when retrieving deleted comment")
}
