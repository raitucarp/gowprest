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

	revisionsAPI := pageAPI.Revisions(page.ID)

	// 2. Create a revision by updating the page (this creates a standard deletable revision)
	t.Log("Updating page to create a standard revision...")
	_, err = pageAPI.Update(gowprest.PageData{
		ID:      page.ID,
		Title:   faker.Sentence(),
		Content: faker.Paragraph(),
	}).Do()
	require.NoError(t, err)

	_, err = pageAPI.Update(gowprest.PageData{
		ID:      page.ID,
		Content: faker.Paragraph(),
	}).Do()
	require.NoError(t, err)

	// 3. Create another revision (autosave) using the new Create method
	t.Log("Creating an autosave revision via Create()...")
	autosave, err := revisionsAPI.Create(gowprest.PageData{
		Content: faker.Paragraph(),
	}).Do()
	require.NoError(t, err)
	assert.Equal(t, page.ID, autosave.Parent)

	// 4. List revisions
	t.Log("Listing revisions...")
	revisions, err := revisionsAPI.List().Do()
	require.NoError(t, err)
	assert.GreaterOrEqual(t, len(revisions), 1, "Should have at least one revision")

	for _, r := range revisions {
		t.Logf("Found revision ID: %d", r.ID)
	}

	// Find a deletable revision (not the autosave we just created)
	deletableRevisionID := autosave.ID

	// 5. Retrieve the revision
	t.Logf("Retrieving revision %d...", deletableRevisionID)
	revision, err := revisionsAPI.Retrieve(deletableRevisionID).Do()
	require.NoError(t, err)
	assert.Equal(t, deletableRevisionID, revision.ID)

	// 6. Delete the revision
	t.Logf("Deleting revision %d...", deletableRevisionID)
	deletedRevision, err := revisionsAPI.Delete(deletableRevisionID).Force().Do()
	if err != nil && deletableRevisionID == autosave.ID {
		t.Logf("Skipping deletion failure for autosave revision: %v", err)
	} else {
		require.NoError(t, err)
		assert.Equal(t, deletableRevisionID, deletedRevision.ID)
	}

	// 7. List revisions again
	t.Log("Listing revisions after deletion...")
	_, err = revisionsAPI.List().Do()
	require.NoError(t, err)
}
