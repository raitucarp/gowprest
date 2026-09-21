package tests

import (
	"errors"
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/go-faker/faker/v4"
	_ "github.com/joho/godotenv/autoload"
	"github.com/raitucarp/gowprest"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestBlockRevisions(t *testing.T) {
	blogUrl := os.Getenv("BLOG_URL")
	username := os.Getenv("BLOG_USERNAME")
	password := os.Getenv("BLOG_APP_PASSWORD")

	if blogUrl == "" || username == "" || password == "" {
		t.Skip("Skipping test; BLOG_URL, BLOG_USERNAME, or BLOG_APP_PASSWORD not set")
	}

	client := gowprest.NewClient(blogUrl).WithBasicAuth(username, password)
	defer client.Close()

	// Setup: Create a block and update it multiple times to generate revisions
	unique := fmt.Sprintf("BlkRev%d", time.Now().UnixNano()%1000000)
	block, err := client.Blocks().Create().
		Title("Parent Block " + unique).
		Content("<!-- wp:paragraph --><p>Initial block content " + faker.Sentence() + "</p><!-- /wp:paragraph -->").
		StatusPublish().
		Do()
	require.NoError(t, err)
	defer client.Blocks().Delete(block.ID).Force().Do()

	// Update 1
	_, err = client.Blocks().Update(block.ID).
		Title("Parent Block Update 1 " + unique).
		Content("<!-- wp:paragraph --><p>Update 1 content " + faker.Sentence() + "</p><!-- /wp:paragraph -->").
		Do()
	require.NoError(t, err)

	// Update 2
	_, err = client.Blocks().Update(block.ID).
		Content("<!-- wp:paragraph --><p>Update 2 content " + faker.Sentence() + "</p><!-- /wp:paragraph -->").
		Do()
	require.NoError(t, err)

	// Update 3
	_, err = client.Blocks().Update(block.ID).
		Content("<!-- wp:paragraph --><p>Update 3 content " + faker.Sentence() + "</p><!-- /wp:paragraph -->").
		Do()
	require.NoError(t, err)

	revisionsAPI := client.Blocks().Revisions(block.ID)

	t.Run("BlockRevisionsAliases", func(t *testing.T) {
		assert.NotNil(t, revisionsAPI)
		assert.NotNil(t, client.BlockRevisions(block.ID))
		assert.NotNil(t, revisionsAPI.List())
		assert.NotNil(t, revisionsAPI.Autosaves())
	})

	t.Run("ListBlockRevisions", func(t *testing.T) {
		revisions, err := revisionsAPI.List().Do()
		require.NoError(t, err)

		t.Logf("Found %d block revisions", len(revisions))

		t.Run("Default", func(t *testing.T) {
			for _, r := range revisions {
				assert.Equal(t, block.ID, r.Parent)
				assert.NotZero(t, r.ID)
			}
		})

		t.Run("Contexts", func(t *testing.T) {
			viewRevs, err := revisionsAPI.List().ContextView().Do()
			require.NoError(t, err)
			assert.Equal(t, len(revisions), len(viewRevs))

			editRevs, err := revisionsAPI.List().ContextEdit().Do()
			require.NoError(t, err)
			assert.Equal(t, len(revisions), len(editRevs))

			embedRevs, err := revisionsAPI.List().ContextEmbed().Do()
			require.NoError(t, err)
			assert.Equal(t, len(revisions), len(embedRevs))
		})

		t.Run("Fields", func(t *testing.T) {
			fieldsRevs, err := revisionsAPI.List().
				Fields("id", "parent").
				Do()
			require.NoError(t, err)
			for _, r := range fieldsRevs {
				assert.NotZero(t, r.ID)
				assert.Equal(t, block.ID, r.Parent)
				assert.Nil(t, r.Content)
			}
		})

		if len(revisions) >= 2 {
			t.Run("Pagination", func(t *testing.T) {
				page1, err := revisionsAPI.List().Page(1).PerPage(1).Do()
				require.NoError(t, err)
				require.Len(t, page1, 1)

				page2, err := revisionsAPI.List().Page(2).PerPage(1).Do()
				require.NoError(t, err)
				require.Len(t, page2, 1)

				assert.NotEqual(t, page1[0].ID, page2[0].ID)
			})

			t.Run("Ordering", func(t *testing.T) {
				asc, err := revisionsAPI.List().OrderAsc().Do()
				require.NoError(t, err)

				desc, err := revisionsAPI.List().OrderDesc().Do()
				require.NoError(t, err)

				assert.NotEqual(t, asc[0].ID, desc[0].ID)
			})
		}
	})

	t.Run("Autosaves", func(t *testing.T) {
		autosaveTitle := "Autosave Title " + unique
		autosaveContent := "<p>Autosave Content " + faker.Sentence() + "</p>"

		autosave, err := revisionsAPI.Autosaves().Create().
			Title(autosaveTitle).
			Content(autosaveContent).
			Do()
		require.NoError(t, err)
		require.NotNil(t, autosave)
		assert.NotZero(t, autosave.ID)
		assert.Equal(t, block.ID, autosave.Parent)

		// List autosaves
		autosaves, err := revisionsAPI.Autosaves().List().Do()
		require.NoError(t, err)
		assert.NotEmpty(t, autosaves)

		// Retrieve autosave
		retrieved, err := revisionsAPI.Autosaves().Retrieve(autosave.ID).Do()
		require.NoError(t, err)
		require.NotNil(t, retrieved)
		assert.Equal(t, autosave.ID, retrieved.ID)
	})

	t.Run("RetrieveBlockRevision", func(t *testing.T) {
		revisions, err := revisionsAPI.List().Do()
		require.NoError(t, err)

		if len(revisions) > 0 {
			target := revisions[0]
			retrieved, err := revisionsAPI.Retrieve(target.ID).Do()
			require.NoError(t, err)
			require.NotNil(t, retrieved)
			assert.Equal(t, target.ID, retrieved.ID)
			assert.Equal(t, block.ID, retrieved.Parent)

			// Retrieve with edit context and fields
			retrievedFields, err := revisionsAPI.Retrieve(target.ID).
				ContextEdit().
				Fields("id", "parent").
				Do()
			require.NoError(t, err)
			require.NotNil(t, retrievedFields)
			assert.Equal(t, target.ID, retrievedFields.ID)
			assert.Nil(t, retrievedFields.Content)
		}

		t.Run("NotFound", func(t *testing.T) {
			_, err := revisionsAPI.Retrieve(9999999).Do()
			require.Error(t, err)

			var wpErr *gowprest.WPRestError
			if assert.True(t, errors.As(err, &wpErr)) {
				assert.Equal(t, "rest_post_invalid_id", wpErr.Code)
				assert.Equal(t, 404, wpErr.Data.Status)
			}
		})
	})

	t.Run("DeleteBlockRevision", func(t *testing.T) {
		// Create a separate block and revision to safely delete
		tmpBlock, err := client.Blocks().Create().
			Title("Delete Rev Parent " + unique).
			Content("<p>Content</p>").
			StatusPublish().
			Do()
		require.NoError(t, err)
		defer client.Blocks().Delete(tmpBlock.ID).Force().Do()

		// Create an autosave revision
		as, err := client.Blocks().Revisions(tmpBlock.ID).Autosaves().Create().
			Title("Autosave to Delete").
			Content("<p>Autosave body</p>").
			Do()
		require.NoError(t, err)
		require.NotNil(t, as)

		// Delete the revision
		deleted, err := client.Blocks().Revisions(tmpBlock.ID).Delete(as.ID).Force().Do()
		require.NoError(t, err)
		require.NotNil(t, deleted)
		assert.Equal(t, as.ID, deleted.ID)

		// Verify deletion
		_, err = client.Blocks().Revisions(tmpBlock.ID).Retrieve(as.ID).Do()
		require.Error(t, err)
	})
}
