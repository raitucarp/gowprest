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

func TestPostRevisions(t *testing.T) {
	blogUrl := os.Getenv("BLOG_URL")
	username := os.Getenv("BLOG_USERNAME")
	password := os.Getenv("BLOG_APP_PASSWORD")

	if blogUrl == "" || username == "" || password == "" {
		t.Skip("Skipping test; BLOG_URL, BLOG_USERNAME, or BLOG_APP_PASSWORD not set")
	}

	client := gowprest.NewClient(blogUrl).WithBasicAuth(username, password)
	defer client.Close()

	postAPI := client.Posts()

	// 1. Create a post
	post, err := postAPI.Create(gowprest.PostData{
		Title:   faker.Sentence(),
		Content: faker.Paragraph(),
		Status:  gowprest.StatusPublished,
	}).Do()
	require.NoError(t, err)
	defer postAPI.Delete(post.ID).Force().Do()

	// 2. Update the post to create a revision
	newTitle := faker.Sentence()
	newContent := faker.Paragraph()
	updatedPost, err := postAPI.Update(gowprest.PostData{
		ID:      post.ID,
		Title:   newTitle,
		Content: newContent,
	}).Do()
	require.NoError(t, err)
	assert.Equal(t, post.ID, updatedPost.ID)
	assert.Equal(t, newTitle, updatedPost.Title.Rendered)
	assert.Contains(t, updatedPost.Content.Rendered, newContent)

	revisionsAPI := postAPI.Revisions(post.ID)

	// 3. List revisions
	// WordPress might take a moment or need some specific state to create a revision.
	// We'll retry a few times if needed.
	var revisions []gowprest.Revision
	for i := 0; i < 5; i++ {
		revisions, err = revisionsAPI.List().Do()
		require.NoError(t, err)
		if len(revisions) >= 1 {
			break
		}
		t.Logf("No revisions found yet, retrying... (%d/5)", i+1)

		// Trigger another update just in case
		_, _ = postAPI.Update(gowprest.PostData{
			ID:      post.ID,
			Content: faker.Paragraph(),
		}).Do()
	}

	assert.GreaterOrEqual(t, len(revisions), 1, "Should have at least one revision")

	if len(revisions) > 0 {
		revisionID := revisions[0].ID

		// 4. Retrieve a specific revision
		revision, err := revisionsAPI.Retrieve(revisionID).Do()
		require.NoError(t, err)
		assert.Equal(t, revisionID, revision.ID)
		assert.Equal(t, post.ID, revision.Parent)

		// 5. Delete a revision (if supported/needed for testing)
		// Usually, revisions can be deleted.
		deletedRevision, err := revisionsAPI.Delete(revisionID).Force().Do()
		require.NoError(t, err)
		assert.Equal(t, revisionID, deletedRevision.ID)

		// Verify deletion
		_, err = revisionsAPI.Retrieve(revisionID).Do()
		assert.Error(t, err, "Should error when retrieving deleted revision")
	}
}
