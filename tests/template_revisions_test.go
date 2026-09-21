package tests

import (
	"errors"
	"os"
	"testing"

	_ "github.com/joho/godotenv/autoload"
	"github.com/raitucarp/gowprest"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestTemplateRevisions(t *testing.T) {
	blogUrl := os.Getenv("BLOG_URL")
	username := os.Getenv("BLOG_USERNAME")
	password := os.Getenv("BLOG_APP_PASSWORD")

	if blogUrl == "" || username == "" || password == "" {
		t.Skip("Skipping test; BLOG_URL, BLOG_USERNAME, or BLOG_APP_PASSWORD not set")
	}

	client := gowprest.NewClient(blogUrl).WithBasicAuth(username, password)
	defer client.Close()

	parentSlug := "gowprest-test-rev-tmpl"

	t.Run("TemplateRevisionsAliases", func(t *testing.T) {
		assert.NotNil(t, client.TemplateRevisions(parentSlug))
		assert.NotNil(t, client.Templates().Revisions(parentSlug))
		assert.NotNil(t, client.TemplateAutosaves(parentSlug))
		assert.NotNil(t, client.Templates().Autosaves(parentSlug))
		assert.NotNil(t, client.TemplateRevisions(parentSlug).Autosaves())

		assert.NotNil(t, client.TemplateRevisions(parentSlug).List())
		assert.NotNil(t, client.TemplateRevisions(parentSlug).Retrieve(1))
		assert.NotNil(t, client.TemplateRevisions(parentSlug).Delete(1))

		assert.NotNil(t, client.TemplateAutosaves(parentSlug).List())
		assert.NotNil(t, client.TemplateAutosaves(parentSlug).Retrieve(1))
		assert.NotNil(t, client.TemplateAutosaves(parentSlug).Create())
	})

	t.Run("TemplateRevisionsWorkflow", func(t *testing.T) {
		// 1. Create a parent custom template
		parentTemplate, err := client.Templates().Create().
			Slug(parentSlug).
			Title("Revision Parent Template").
			Content("<!-- wp:paragraph --><p>Initial Content</p><!-- /wp:paragraph -->").
			StatusPublish().
			Do()
		require.NoError(t, err)
		require.NotNil(t, parentTemplate)
		parentID := parentTemplate.ID
		defer func() {
			_, _ = client.Templates().Delete(parentID).Force().Do()
		}()

		// 2. Perform updates to produce revisions
		_, err = client.Templates().Update(parentID).
			Title("Revision Parent Template Updated 1").
			Content("<!-- wp:paragraph --><p>Updated Content 1</p><!-- /wp:paragraph -->").
			Do()
		require.NoError(t, err)

		_, err = client.Templates().Update(parentID).
			Title("Revision Parent Template Updated 2").
			Content("<!-- wp:paragraph --><p>Updated Content 2</p><!-- /wp:paragraph -->").
			Do()
		require.NoError(t, err)

		// 3. List revisions
		revisions, err := client.TemplateRevisions(parentID).List().
			ContextEdit().
			Page(1).
			PerPage(10).
			OrderDesc().
			OrderBy("date").
			Do()
		require.NoError(t, err)
		require.NotEmpty(t, revisions)

		rev := revisions[0]
		assert.NotEmpty(t, rev.ID)
		assert.NotZero(t, rev.WPID)
		assert.NotZero(t, rev.Parent)
		assert.Equal(t, 1, rev.Author)
		assert.NotEmpty(t, rev.TitleString())
		assert.NotEmpty(t, rev.ContentString())

		// ContextView & ContextEmbed
		_, err = client.TemplateRevisions(parentID).List().ContextView().Do()
		require.NoError(t, err)

		_, err = client.TemplateRevisions(parentID).List().ContextEmbed().Do()
		require.NoError(t, err)

		// Sparse fields
		sparseRevs, err := client.TemplateRevisions(parentID).List().
			Fields("id", "wp_id", "title").
			Do()
		require.NoError(t, err)
		require.NotEmpty(t, sparseRevs)
		assert.NotEmpty(t, sparseRevs[0].ID)
		assert.NotZero(t, sparseRevs[0].WPID)
		assert.NotEmpty(t, sparseRevs[0].TitleString())
		assert.Empty(t, sparseRevs[0].Description)

		revID := rev.WPID

		// 4. Retrieve single revision
		retrieved, err := client.TemplateRevisions(parentID).Retrieve(revID).
			ContextEdit().
			Do()
		require.NoError(t, err)
		require.NotNil(t, retrieved)
		assert.Equal(t, revID, retrieved.WPID)
		assert.Equal(t, rev.ID, retrieved.ID)
		assert.Equal(t, rev.TitleString(), retrieved.TitleString())

		// Sparse retrieve
		sparseRev, err := client.TemplateRevisions(parentID).Retrieve(revID).
			Fields("id", "wp_id").
			Do()
		require.NoError(t, err)
		require.NotNil(t, sparseRev)
		assert.Equal(t, revID, sparseRev.WPID)
		assert.Empty(t, sparseRev.Description)

		// 5. Delete revision without force -> verify 501 rest_trash_not_supported
		_, err = client.TemplateRevisions(parentID).Delete(revID).Do()
		require.Error(t, err)
		var trashErr *gowprest.WPRestError
		require.True(t, errors.As(err, &trashErr))
		assert.Equal(t, 501, trashErr.Data.Status)
		assert.Equal(t, "rest_trash_not_supported", trashErr.Code)

		// 6. Delete revision with Force
		deleted, err := client.TemplateRevisions(parentID).Delete(revID).
			Force().
			Do()
		require.NoError(t, err)
		require.NotNil(t, deleted)
		assert.Equal(t, revID, deleted.WPID)

		// Verify deletion
		_, err = client.TemplateRevisions(parentID).Retrieve(revID).Do()
		require.Error(t, err)
		var notFoundErr *gowprest.WPRestError
		require.True(t, errors.As(err, &notFoundErr))
		assert.Equal(t, 404, notFoundErr.Data.Status)
	})

	t.Run("TemplateAutosavesWorkflow", func(t *testing.T) {
		autosaveParentSlug := "gowprest-test-autosave-tmpl"

		// 1. Create a parent template
		parentTemplate, err := client.Templates().Create().
			Slug(autosaveParentSlug).
			Title("Autosave Parent Template").
			Content("<!-- wp:paragraph --><p>Autosave Initial</p><!-- /wp:paragraph -->").
			StatusPublish().
			Do()
		require.NoError(t, err)
		require.NotNil(t, parentTemplate)
		parentID := parentTemplate.ID
		defer func() {
			_, _ = client.Templates().Delete(parentID).Force().Do()
		}()

		// 2. Create an autosave
		autosave, err := client.TemplateAutosaves(parentID).Create().
			Title("Autosaved Template Title").
			Content("<!-- wp:paragraph --><p>Autosaved Template Content</p><!-- /wp:paragraph -->").
			Do()
		require.NoError(t, err)
		require.NotNil(t, autosave)
		assert.NotZero(t, autosave.WPID)

		autosaveID := autosave.WPID

		// 3. List autosaves
		autosaves, err := client.TemplateAutosaves(parentID).List().
			ContextEdit().
			Do()
		require.NoError(t, err)
		assert.NotEmpty(t, autosaves)
		assert.Equal(t, autosaveID, autosaves[0].WPID)

		// 4. Retrieve single autosave
		retrievedAutosave, err := client.TemplateAutosaves(parentID).Retrieve(autosaveID).
			ContextEdit().
			Do()
		require.NoError(t, err)
		require.NotNil(t, retrievedAutosave)
		assert.Equal(t, autosaveID, retrievedAutosave.WPID)
		assert.Contains(t, retrievedAutosave.TitleString(), "Autosaved Template Title")
	})

	t.Run("TemplateRevisionsUnauthorizedAndInvalid", func(t *testing.T) {
		unauthParentSlug := "gowprest-test-unauth-tmpl"

		// Create a custom template first so it exists
		created, err := client.Templates().Create().
			Slug(unauthParentSlug).
			Title("Unauth Test Parent").
			Content("<!-- wp:paragraph --><p>Content</p><!-- /wp:paragraph -->").
			StatusPublish().
			Do()
		require.NoError(t, err)
		require.NotNil(t, created)
		defer func() {
			_, _ = client.Templates().Delete(created.ID).Force().Do()
		}()

		unauthClient := gowprest.NewClient(blogUrl)
		defer unauthClient.Close()

		_, err = unauthClient.TemplateRevisions(created.ID).List().Do()
		require.Error(t, err)
		var unauthErr *gowprest.WPRestError
		require.True(t, errors.As(err, &unauthErr))
		assert.Equal(t, 401, unauthErr.Data.Status)
		assert.Equal(t, "rest_cannot_read", unauthErr.Code)

		// Theme template cannot have revisions (400 rest_invalid_template)
		_, err = client.TemplateRevisions("twentytwentyfive//index").List().Do()
		require.Error(t, err)
		var invalidErr *gowprest.WPRestError
		require.True(t, errors.As(err, &invalidErr))
		assert.Equal(t, 400, invalidErr.Data.Status)
		assert.Equal(t, "rest_invalid_template", invalidErr.Code)
	})
}
