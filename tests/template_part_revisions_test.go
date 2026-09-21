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

func TestTemplatePartRevisions(t *testing.T) {
	blogUrl := os.Getenv("BLOG_URL")
	username := os.Getenv("BLOG_USERNAME")
	password := os.Getenv("BLOG_APP_PASSWORD")

	if blogUrl == "" || username == "" || password == "" {
		t.Skip("Skipping test; BLOG_URL, BLOG_USERNAME, or BLOG_APP_PASSWORD not set")
	}

	client := gowprest.NewClient(blogUrl).WithBasicAuth(username, password)
	defer client.Close()

	parentSlug := "gowprest-test-part-rev"

	t.Run("TemplatePartRevisionsAliases", func(t *testing.T) {
		assert.NotNil(t, client.TemplatePartRevisions(parentSlug))
		assert.NotNil(t, client.TemplateParts().Revisions(parentSlug))
		assert.NotNil(t, client.TemplatePartAutosaves(parentSlug))
		assert.NotNil(t, client.TemplateParts().Autosaves(parentSlug))
		assert.NotNil(t, client.TemplatePartRevisions(parentSlug).Autosaves())

		assert.NotNil(t, client.TemplatePartRevisions(parentSlug).List())
		assert.NotNil(t, client.TemplatePartRevisions(parentSlug).Retrieve(1))
		assert.NotNil(t, client.TemplatePartRevisions(parentSlug).Delete(1))

		assert.NotNil(t, client.TemplatePartAutosaves(parentSlug).List())
		assert.NotNil(t, client.TemplatePartAutosaves(parentSlug).Retrieve(1))
		assert.NotNil(t, client.TemplatePartAutosaves(parentSlug).Create())
	})

	t.Run("TemplatePartRevisionsWorkflow", func(t *testing.T) {
		// 1. Create a parent custom template part
		parentPart, err := client.TemplateParts().Create().
			Slug(parentSlug).
			Title("Revision Parent Part").
			Content("<!-- wp:paragraph --><p>Initial Part Content</p><!-- /wp:paragraph -->").
			Area("header").
			StatusPublish().
			Do()
		require.NoError(t, err)
		require.NotNil(t, parentPart)
		parentID := parentPart.ID
		defer func() {
			_, _ = client.TemplateParts().Delete(parentID).Force().Do()
		}()

		// 2. Perform updates to produce revisions
		_, err = client.TemplateParts().Update(parentID).
			Title("Revision Parent Part Updated 1").
			Content("<!-- wp:paragraph --><p>Updated Part Content 1</p><!-- /wp:paragraph -->").
			Do()
		require.NoError(t, err)

		_, err = client.TemplateParts().Update(parentID).
			Title("Revision Parent Part Updated 2").
			Content("<!-- wp:paragraph --><p>Updated Part Content 2</p><!-- /wp:paragraph -->").
			Do()
		require.NoError(t, err)

		// 3. List revisions
		revisions, err := client.TemplatePartRevisions(parentID).List().
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
		assert.Equal(t, "header", rev.Area)
		assert.NotEmpty(t, rev.TitleString())
		assert.NotEmpty(t, rev.ContentString())

		// ContextView & ContextEmbed
		_, err = client.TemplatePartRevisions(parentID).List().ContextView().Do()
		require.NoError(t, err)

		_, err = client.TemplatePartRevisions(parentID).List().ContextEmbed().Do()
		require.NoError(t, err)

		// Sparse fields
		sparseRevs, err := client.TemplatePartRevisions(parentID).List().
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
		retrieved, err := client.TemplatePartRevisions(parentID).Retrieve(revID).
			ContextEdit().
			Do()
		require.NoError(t, err)
		require.NotNil(t, retrieved)
		assert.Equal(t, revID, retrieved.WPID)
		assert.Equal(t, rev.ID, retrieved.ID)
		assert.Equal(t, "header", retrieved.Area)
		assert.Equal(t, rev.TitleString(), retrieved.TitleString())

		// Sparse retrieve
		sparseRev, err := client.TemplatePartRevisions(parentID).Retrieve(revID).
			Fields("id", "wp_id").
			Do()
		require.NoError(t, err)
		require.NotNil(t, sparseRev)
		assert.Equal(t, revID, sparseRev.WPID)
		assert.Empty(t, sparseRev.Description)

		// 5. Delete revision without force -> verify 501 rest_trash_not_supported
		_, err = client.TemplatePartRevisions(parentID).Delete(revID).Do()
		require.Error(t, err)
		var trashErr *gowprest.WPRestError
		require.True(t, errors.As(err, &trashErr))
		assert.Equal(t, 501, trashErr.Data.Status)
		assert.Equal(t, "rest_trash_not_supported", trashErr.Code)

		// 6. Delete revision with Force
		deleted, err := client.TemplatePartRevisions(parentID).Delete(revID).
			Force().
			Do()
		require.NoError(t, err)
		require.NotNil(t, deleted)
		assert.Equal(t, revID, deleted.WPID)

		// Verify deletion
		_, err = client.TemplatePartRevisions(parentID).Retrieve(revID).Do()
		require.Error(t, err)
		var notFoundErr *gowprest.WPRestError
		require.True(t, errors.As(err, &notFoundErr))
		assert.Equal(t, 404, notFoundErr.Data.Status)
	})

	t.Run("TemplatePartAutosavesWorkflow", func(t *testing.T) {
		autosaveParentSlug := "gowprest-test-part-autosave"

		// 1. Create a parent template part
		parentPart, err := client.TemplateParts().Create().
			Slug(autosaveParentSlug).
			Title("Autosave Parent Part").
			Content("<!-- wp:paragraph --><p>Autosave Part Initial</p><!-- /wp:paragraph -->").
			Area("footer").
			StatusPublish().
			Do()
		require.NoError(t, err)
		require.NotNil(t, parentPart)
		parentID := parentPart.ID
		defer func() {
			_, _ = client.TemplateParts().Delete(parentID).Force().Do()
		}()

		// 2. Create an autosave
		autosave, err := client.TemplatePartAutosaves(parentID).Create().
			Title("Autosaved Part Title").
			Content("<!-- wp:paragraph --><p>Autosaved Part Content</p><!-- /wp:paragraph -->").
			Area("footer").
			Do()
		require.NoError(t, err)
		require.NotNil(t, autosave)
		assert.NotZero(t, autosave.WPID)

		autosaveID := autosave.WPID

		// 3. List autosaves
		autosaves, err := client.TemplatePartAutosaves(parentID).List().
			ContextEdit().
			Do()
		require.NoError(t, err)
		assert.NotEmpty(t, autosaves)
		assert.Equal(t, autosaveID, autosaves[0].WPID)

		// 4. Retrieve single autosave
		retrievedAutosave, err := client.TemplatePartAutosaves(parentID).Retrieve(autosaveID).
			ContextEdit().
			Do()
		require.NoError(t, err)
		require.NotNil(t, retrievedAutosave)
		assert.Equal(t, autosaveID, retrievedAutosave.WPID)
		assert.Contains(t, retrievedAutosave.TitleString(), "Autosaved Part Title")
	})

	t.Run("TemplatePartRevisionsUnauthorizedAndInvalid", func(t *testing.T) {
		unauthParentSlug := "gowprest-test-unauth-part"

		// Create a custom template part first so it exists
		created, err := client.TemplateParts().Create().
			Slug(unauthParentSlug).
			Title("Unauth Test Parent Part").
			Content("<!-- wp:paragraph --><p>Content</p><!-- /wp:paragraph -->").
			Area("header").
			StatusPublish().
			Do()
		require.NoError(t, err)
		require.NotNil(t, created)
		defer func() {
			_, _ = client.TemplateParts().Delete(created.ID).Force().Do()
		}()

		unauthClient := gowprest.NewClient(blogUrl)
		defer unauthClient.Close()

		_, err = unauthClient.TemplatePartRevisions(created.ID).List().Do()
		require.Error(t, err)
		var unauthErr *gowprest.WPRestError
		require.True(t, errors.As(err, &unauthErr))
		assert.Equal(t, 401, unauthErr.Data.Status)
		assert.Equal(t, "rest_cannot_read", unauthErr.Code)

		// Theme template part cannot have revisions (400 rest_invalid_template)
		_, err = client.TemplatePartRevisions("twentytwentyfive//footer-newsletter").List().Do()
		require.Error(t, err)
		var invalidErr *gowprest.WPRestError
		require.True(t, errors.As(err, &invalidErr))
		assert.Equal(t, 400, invalidErr.Data.Status)
		assert.Equal(t, "rest_invalid_template", invalidErr.Code)
	})
}
