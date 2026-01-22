package tests

import (
	"os"
	"testing"
	"time"

	"github.com/go-faker/faker/v4"
	_ "github.com/joho/godotenv/autoload"
	"github.com/raitucarp/gowprest"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestPageRevisions(t *testing.T) {
	blogUrl := os.Getenv("BLOG_URL")
	username := os.Getenv("BLOG_USERNAME")
	password := os.Getenv("BLOG_APP_PASSWORD")

	if blogUrl == "" || username == "" || password == "" {
		t.Skip("Skipping test; BLOG_URL, BLOG_USERNAME, or BLOG_APP_PASSWORD not set")
	}

	client := gowprest.NewClient(blogUrl).WithBasicAuth(username, password)
	defer client.Close()

	pageAPI := client.Pages()

	// 1. Create a page
	page, err := pageAPI.Create(gowprest.PageData{
		Title:   faker.Sentence(),
		Content: faker.Paragraph(),
		Status:  gowprest.StatusPublished,
	}).Do()
	require.NoError(t, err)
	defer pageAPI.Delete(page.ID).Force().Do()

	// 2. Update the page to create a revision
	newTitle := faker.Sentence()
	newContent := faker.Paragraph()
	updatedPage, err := pageAPI.Update(gowprest.PageData{
		ID:      page.ID,
		Title:   newTitle,
		Content: newContent,
	}).Do()
	require.NoError(t, err)
	assert.Equal(t, page.ID, updatedPage.ID)
	assert.Equal(t, newTitle, updatedPage.Title.Rendered)
	assert.Contains(t, updatedPage.Content.Rendered, newContent)

	revisionsAPI := pageAPI.Revisions(page.ID)

	// 3. List revisions
	// WordPress might take a moment or need some specific state to create a revision.
	// We'll retry a few times if needed.
	var revisions []gowprest.Revision
	for i := 0; i < 10; i++ {
		revisions, err = revisionsAPI.List().Do()
		require.NoError(t, err)
		if len(revisions) >= 1 {
			break
		}
		t.Logf("No revisions found yet, retrying... (%d/10)", i+1)

		// Trigger another update just in case and sleep
		_, _ = pageAPI.Update(gowprest.PageData{
			ID:      page.ID,
			Content: faker.Paragraph(),
		}).Do()
		time.Sleep(2 * time.Second)
	}

	assert.GreaterOrEqual(t, len(revisions), 1, "Should have at least one revision")

	if len(revisions) > 0 {
		revisionID := revisions[0].ID

		// 4. Retrieve a specific revision
		revision, err := revisionsAPI.Retrieve(revisionID).Do()
		require.NoError(t, err)
		assert.Equal(t, revisionID, revision.ID)
		assert.Equal(t, page.ID, revision.Parent)

		// 5. Delete a revision
		deletedRevision, err := revisionsAPI.Delete(revisionID).Force().Do()
		require.NoError(t, err)
		assert.Equal(t, revisionID, deletedRevision.ID)

		// Verify deletion
		_, err = revisionsAPI.Retrieve(revisionID).Do()
		assert.Error(t, err, "Should error when retrieving deleted revision")
	}
}
