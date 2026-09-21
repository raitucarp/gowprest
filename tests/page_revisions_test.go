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

	// Setup: Create a page and update it multiple times to generate standard revisions
	page, err := pageAPI.Create().
		Title("Revision Page " + faker.Sentence()).
		Content("<p>Initial page content: " + faker.Paragraph() + "</p>").
		StatusPublish().
		Do()
	require.NoError(t, err)
	defer pageAPI.Delete(page.ID).Force().Do()

	// Update 1
	_, err = pageAPI.Update().
		ID(page.ID).
		Title("Revision Page " + faker.Sentence()).
		Content("<p>Update 1: " + faker.Paragraph() + "</p>").
		Do()
	require.NoError(t, err)

	// Update 2
	_, err = pageAPI.Update().
		ID(page.ID).
		Content("<p>Update 2: " + faker.Paragraph() + "</p>").
		Do()
	require.NoError(t, err)

	// Update 3
	_, err = pageAPI.Update().
		ID(page.ID).
		Content("<p>Update 3: " + faker.Paragraph() + "</p>").
		Do()
	require.NoError(t, err)

	revisionsAPI := pageAPI.Revisions(page.ID)

	t.Run("ListPageRevisions_Default", func(t *testing.T) {
		revisions, err := revisionsAPI.List().Do()
		require.NoError(t, err)
		assert.GreaterOrEqual(t, len(revisions), 2, "Should have at least 2 revisions generated from updates")
		for _, r := range revisions {
			assert.Equal(t, page.ID, r.Parent)
			assert.NotEmpty(t, r.ID)
		}
	})

	t.Run("ListPageRevisions_Pagination", func(t *testing.T) {
		revisionsPage1, err := revisionsAPI.List().Page(1).PerPage(1).Do()
		require.NoError(t, err)
		assert.Len(t, revisionsPage1, 1)

		revisionsPage2, err := revisionsAPI.List().Page(2).PerPage(1).Do()
		require.NoError(t, err)
		assert.Len(t, revisionsPage2, 1)
		assert.NotEqual(t, revisionsPage1[0].ID, revisionsPage2[0].ID)

		// With offset
		revisionsOffset, err := revisionsAPI.List().Offset(1).PerPage(1).Do()
		require.NoError(t, err)
		assert.Len(t, revisionsOffset, 1)
		assert.Equal(t, revisionsPage2[0].ID, revisionsOffset[0].ID)
	})

	t.Run("ListPageRevisions_Ordering", func(t *testing.T) {
		revsAsc, err := revisionsAPI.List().OrderAsc().OrderByDate().Do()
		require.NoError(t, err)
		require.GreaterOrEqual(t, len(revsAsc), 2)

		revsDesc, err := revisionsAPI.List().OrderDesc().OrderBy("date").Do()
		require.NoError(t, err)
		require.GreaterOrEqual(t, len(revsDesc), 2)

		// Asc first element ID should differ from Desc first element ID
		assert.NotEqual(t, revsAsc[0].ID, revsDesc[0].ID)
		assert.Equal(t, revsAsc[0].ID, revsDesc[len(revsDesc)-1].ID)

		// Order by ID
		revsById, err := revisionsAPI.List().OrderAsc().OrderById().Do()
		require.NoError(t, err)
		require.GreaterOrEqual(t, len(revsById), 2)
		assert.True(t, revsById[0].ID < revsById[1].ID)
	})

	t.Run("ListPageRevisions_IncludeAndExclude", func(t *testing.T) {
		allRevs, err := revisionsAPI.List().Do()
		require.NoError(t, err)
		require.GreaterOrEqual(t, len(allRevs), 2)

		id1 := allRevs[0].ID
		id2 := allRevs[1].ID

		// Include
		included, err := revisionsAPI.List().Include(id1, id2).Do()
		require.NoError(t, err)
		assert.Len(t, included, 2)

		// Exclude
		excluded, err := revisionsAPI.List().Exclude(id1).Do()
		require.NoError(t, err)
		for _, r := range excluded {
			assert.NotEqual(t, id1, r.ID)
		}
	})

	t.Run("ListPageRevisions_ContextAndGlobalParams", func(t *testing.T) {
		// Context edit
		editRevs, err := revisionsAPI.List().ContextEdit().Do()
		require.NoError(t, err)
		require.NotEmpty(t, editRevs)
		if editRevs[0].Title != nil {
			assert.NotEmpty(t, editRevs[0].Title.Raw)
		}

		// Embed
		embedRevs, err := revisionsAPI.List().Embed().Do()
		require.NoError(t, err)
		require.NotEmpty(t, embedRevs)

		// Fields filter
		fieldsRevs, err := revisionsAPI.List().Fields("id", "parent").Do()
		require.NoError(t, err)
		require.NotEmpty(t, fieldsRevs)
		assert.NotZero(t, fieldsRevs[0].ID)
		assert.Equal(t, page.ID, fieldsRevs[0].Parent)
		assert.Nil(t, fieldsRevs[0].Date)
	})

	t.Run("RetrievePageRevision_SuccessAndContexts", func(t *testing.T) {
		allRevs, err := revisionsAPI.List().Do()
		require.NoError(t, err)
		require.NotEmpty(t, allRevs)

		targetID := allRevs[0].ID

		// Default / View context
		revView, err := revisionsAPI.Retrieve(targetID).ContextView().Do()
		require.NoError(t, err)
		assert.Equal(t, targetID, revView.ID)
		assert.Equal(t, page.ID, revView.Parent)
		if revView.Title != nil {
			assert.NotEmpty(t, revView.Title.Rendered)
		}

		// Edit context
		revEdit, err := revisionsAPI.Retrieve(targetID).Context("edit").Embed().Do()
		require.NoError(t, err)
		assert.Equal(t, targetID, revEdit.ID)
		if revEdit.Title != nil {
			assert.NotEmpty(t, revEdit.Title.Raw)
		}

		// Fields
		revFields, err := revisionsAPI.Retrieve(targetID).Fields("id", "parent").Do()
		require.NoError(t, err)
		assert.Equal(t, targetID, revFields.ID)
		assert.Equal(t, page.ID, revFields.Parent)
		assert.Nil(t, revFields.Date)
	})

	t.Run("RetrievePageRevision_NotFound", func(t *testing.T) {
		_, err := revisionsAPI.Retrieve(999999999).Do()
		assert.Error(t, err)
		var wpErr *gowprest.WPRestError
		if assert.ErrorAs(t, err, &wpErr) {
			assert.Equal(t, 404, wpErr.Data.Status)
		}
	})

	var autosaveID int
	t.Run("CreatePageRevision_Autosave_MethodChaining", func(t *testing.T) {
		autosaveContent := "<p>Autosaved content: " + faker.Sentence() + "</p>"
		autosave, err := revisionsAPI.Create().
			Content(autosaveContent).
			Do()
		require.NoError(t, err)
		assert.Equal(t, page.ID, autosave.Parent)
		assert.NotZero(t, autosave.ID)
		autosaveID = autosave.ID
	})

	t.Run("PageAutosaves_ListAndRetrieve", func(t *testing.T) {
		if autosaveID == 0 {
			t.Skip("No autosave created to test autosaves endpoint")
		}

		autosavesAPI := revisionsAPI.Autosaves()

		// List autosaves
		autosaves, err := autosavesAPI.List().ContextView().Embed().Do()
		require.NoError(t, err)
		assert.NotEmpty(t, autosaves)

		// Retrieve autosave
		retrieved, err := autosavesAPI.Retrieve(autosaveID).ContextEdit().Do()
		require.NoError(t, err)
		assert.Equal(t, autosaveID, retrieved.ID)
		assert.Equal(t, page.ID, retrieved.Parent)
	})

	t.Run("DeletePageRevision_SuccessAndNotFound", func(t *testing.T) {
		allRevs, err := revisionsAPI.List().Do()
		require.NoError(t, err)
		require.NotEmpty(t, allRevs)

		// Pick the oldest revision to delete so the page maintains its latest revisions
		deleteTarget := allRevs[len(allRevs)-1]

		deleted, err := revisionsAPI.Delete(deleteTarget.ID).Do()
		require.NoError(t, err)
		assert.Equal(t, deleteTarget.ID, deleted.ID)

		// Verify retrieval fails after delete
		_, err = revisionsAPI.Retrieve(deleteTarget.ID).Do()
		assert.Error(t, err)
	})
}
