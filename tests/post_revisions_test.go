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

	// Setup: Create a post and update it multiple times to generate standard revisions
	post, err := postAPI.Create().
		Title("Revision Test " + faker.Sentence()).
		Content("<p>Initial content: " + faker.Paragraph() + "</p>").
		StatusPublish().
		Do()
	require.NoError(t, err)
	defer postAPI.Delete(post.ID).Force().Do()

	// Update 1
	_, err = postAPI.Update().
		ID(post.ID).
		Title("Revision Test " + faker.Sentence()).
		Content("<p>Update 1: " + faker.Paragraph() + "</p>").
		Do()
	require.NoError(t, err)

	// Update 2
	_, err = postAPI.Update().
		ID(post.ID).
		Content("<p>Update 2: " + faker.Paragraph() + "</p>").
		Do()
	require.NoError(t, err)

	// Update 3
	_, err = postAPI.Update().
		ID(post.ID).
		Content("<p>Update 3: " + faker.Paragraph() + "</p>").
		Do()
	require.NoError(t, err)

	revisionsAPI := postAPI.Revisions(post.ID)

	t.Run("ListPostRevisions_Default", func(t *testing.T) {
		revisions, err := revisionsAPI.List().Do()
		require.NoError(t, err)
		assert.GreaterOrEqual(t, len(revisions), 2, "Should have at least 2 revisions generated from updates")
		for _, r := range revisions {
			assert.Equal(t, post.ID, r.Parent)
			assert.NotEmpty(t, r.ID)
		}
	})

	t.Run("ListPostRevisions_Pagination", func(t *testing.T) {
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

	t.Run("ListPostRevisions_Ordering", func(t *testing.T) {
		revsAsc, err := revisionsAPI.List().OrderAsc().OrderByDate().Do()
		require.NoError(t, err)
		require.GreaterOrEqual(t, len(revsAsc), 2)

		revsDesc, err := revisionsAPI.List().OrderDesc().OrderBy("date").Do()
		require.NoError(t, err)
		require.GreaterOrEqual(t, len(revsDesc), 2)

		// Asc first element ID should differ from Desc first element ID
		assert.NotEqual(t, revsAsc[0].ID, revsDesc[0].ID)
		assert.Equal(t, revsAsc[0].ID, revsDesc[len(revsDesc)-1].ID)

		// Other OrderBy variants
		_, err = revisionsAPI.List().OrderById().Do()
		require.NoError(t, err)

		_, err = revisionsAPI.List().OrderByTitle().Do()
		require.NoError(t, err)

		_, err = revisionsAPI.List().OrderBySlug().Do()
		require.NoError(t, err)
	})

	t.Run("ListPostRevisions_IncludeAndExclude", func(t *testing.T) {
		allRevs, err := revisionsAPI.List().Do()
		require.NoError(t, err)
		require.GreaterOrEqual(t, len(allRevs), 2)

		targetID := allRevs[0].ID
		incRevs, err := revisionsAPI.List().Include(targetID).Do()
		require.NoError(t, err)
		assert.Len(t, incRevs, 1)
		assert.Equal(t, targetID, incRevs[0].ID)

		excRevs, err := revisionsAPI.List().Exclude(targetID).Do()
		require.NoError(t, err)
		for _, r := range excRevs {
			assert.NotEqual(t, targetID, r.ID)
		}
	})

	t.Run("ListPostRevisions_ContextAndGlobalParams", func(t *testing.T) {
		// ContextEdit
		revsEdit, err := revisionsAPI.List().ContextEdit().Do()
		require.NoError(t, err)
		require.NotEmpty(t, revsEdit)
		assert.NotEmpty(t, revsEdit[0].Content.Rendered)

		// ContextEmbed
		revsEmbed, err := revisionsAPI.List().ContextEmbed().Do()
		require.NoError(t, err)
		require.NotEmpty(t, revsEmbed)

		// Fields filter
		revsFields, err := revisionsAPI.List().Fields("id", "parent").Do()
		require.NoError(t, err)
		require.NotEmpty(t, revsFields)
		assert.NotZero(t, revsFields[0].ID)
		assert.Equal(t, post.ID, revsFields[0].Parent)
		assert.Nil(t, revsFields[0].Content) // Content was not requested
	})

	t.Run("RetrievePostRevision_SuccessAndContexts", func(t *testing.T) {
		allRevs, err := revisionsAPI.List().Do()
		require.NoError(t, err)
		require.NotEmpty(t, allRevs)

		revID := allRevs[0].ID

		// ContextView
		revView, err := revisionsAPI.Retrieve(revID).ContextView().Do()
		require.NoError(t, err)
		assert.Equal(t, revID, revView.ID)
		assert.Equal(t, post.ID, revView.Parent)
		assert.NotNil(t, revView.Title)
		assert.NotNil(t, revView.Content)

		// ContextEdit
		revEdit, err := revisionsAPI.Retrieve(revID).ContextEdit().Embed().Do()
		require.NoError(t, err)
		assert.Equal(t, revID, revEdit.ID)

		// Fields filter
		revFields, err := revisionsAPI.Retrieve(revID).Fields("id", "slug").Do()
		require.NoError(t, err)
		assert.Equal(t, revID, revFields.ID)
		assert.NotEmpty(t, revFields.Slug)
		assert.Nil(t, revFields.Title)
	})

	t.Run("RetrievePostRevision_NotFound", func(t *testing.T) {
		_, err := revisionsAPI.Retrieve(99999999).Do()
		require.Error(t, err)
		var wpErr *gowprest.WPRestError
		require.ErrorAs(t, err, &wpErr)
		assert.Equal(t, 404, wpErr.Data.Status)
	})

	t.Run("CreatePostRevision_Autosave_MethodChaining", func(t *testing.T) {
		autosaveContent := "<p>Autosave chained: " + faker.Paragraph() + "</p>"
		autosave, err := revisionsAPI.Create().
			Title("Autosave Title " + faker.Word()).
			Content(autosaveContent).
			Excerpt("Autosave excerpt").
			Do()
		require.NoError(t, err)
		assert.Equal(t, post.ID, autosave.Parent)
		assert.NotZero(t, autosave.ID)

		// Check autosaves endpoint
		autosavesList, err := revisionsAPI.Autosaves().List().Do()
		require.NoError(t, err)
		assert.NotEmpty(t, autosavesList)

		// Retrieve specific autosave
		retrievedAutosave, err := revisionsAPI.Autosaves().Retrieve(autosave.ID).Do()
		require.NoError(t, err)
		assert.Equal(t, autosave.ID, retrievedAutosave.ID)
	})

	t.Run("DeletePostRevision_SuccessAndNotFound", func(t *testing.T) {
		// Get standard revisions
		allRevs, err := revisionsAPI.List().Do()
		require.NoError(t, err)
		require.GreaterOrEqual(t, len(allRevs), 2)

		// Pick the oldest deletable standard revision
		targetRev := allRevs[len(allRevs)-1]

		deletedRev, err := revisionsAPI.Delete(targetRev.ID).Force().Do()
		require.NoError(t, err)
		assert.Equal(t, targetRev.ID, deletedRev.ID)
		assert.Equal(t, post.ID, deletedRev.Parent)

		// Ensure it is deleted
		_, err = revisionsAPI.Retrieve(targetRev.ID).Do()
		require.Error(t, err)

		// Delete again should return 404
		_, err = revisionsAPI.Delete(targetRev.ID).Force().Do()
		require.Error(t, err)
		var wpErr *gowprest.WPRestError
		require.ErrorAs(t, err, &wpErr)
		assert.Equal(t, 404, wpErr.Data.Status)
	})
}
