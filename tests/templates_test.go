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

func TestTemplates(t *testing.T) {
	blogUrl := os.Getenv("BLOG_URL")
	username := os.Getenv("BLOG_USERNAME")
	password := os.Getenv("BLOG_APP_PASSWORD")

	if blogUrl == "" || username == "" || password == "" {
		t.Skip("Skipping test; BLOG_URL, BLOG_USERNAME, or BLOG_APP_PASSWORD not set")
	}

	client := gowprest.NewClient(blogUrl).WithBasicAuth(username, password)
	defer client.Close()

	t.Run("TemplatesAliases", func(t *testing.T) {
		assert.NotNil(t, client.Templates())
		assert.NotNil(t, client.Template())

		assert.NotNil(t, client.Templates().List())
		assert.NotNil(t, client.Templates().Retrieve("twentytwentyfive//index"))
		assert.NotNil(t, client.Templates().Create())
		assert.NotNil(t, client.Templates().Update("twentytwentyfive//index"))
		assert.NotNil(t, client.Templates().Delete("twentytwentyfive//index"))
	})

	t.Run("ListTemplates", func(t *testing.T) {
		templates, err := client.Templates().List().
			ContextEdit().
			Do()
		require.NoError(t, err)
		assert.NotEmpty(t, templates)

		// Verify fields of the first template
		tmpl := templates[0]
		assert.NotEmpty(t, tmpl.ID)
		assert.NotEmpty(t, tmpl.Slug)
		assert.NotEmpty(t, tmpl.Theme)
		assert.NotEmpty(t, tmpl.Type)

		// ContextView & ContextEmbed
		_, err = client.Templates().List().ContextView().Do()
		require.NoError(t, err)

		_, err = client.Templates().List().ContextEmbed().Do()
		require.NoError(t, err)

		// Filter by PostType
		byPostType, err := client.Templates().List().
			PostType("post").
			Do()
		require.NoError(t, err)
		assert.NotNil(t, byPostType)

		// Sparse fields
		sparseList, err := client.Templates().List().
			Fields("id", "slug", "theme").
			Do()
		require.NoError(t, err)
		assert.NotEmpty(t, sparseList)
		assert.NotEmpty(t, sparseList[0].ID)
		assert.NotEmpty(t, sparseList[0].Slug)
		assert.NotEmpty(t, sparseList[0].Theme)
		assert.Empty(t, sparseList[0].Description)
	})

	t.Run("TemplateCRUDWorkflow", func(t *testing.T) {
		slug := "gowprest-test-template"

		// 1. Create a custom template
		created, err := client.Templates().Create().
			Slug(slug).
			Title("Gowprest Test Template").
			Content("<!-- wp:paragraph --><p>Custom template content for testing</p><!-- /wp:paragraph -->").
			Description("Created by gowprest integration tests").
			StatusPublish().
			Do()
		require.NoError(t, err)
		require.NotNil(t, created)
		assert.NotEmpty(t, created.ID)
		assert.Equal(t, slug, created.Slug)
		assert.Equal(t, "custom", created.Source)
		assert.True(t, created.IsCustom)
		assert.Equal(t, "publish", created.Status)
		require.NotNil(t, created.Title)
		assert.Equal(t, "Gowprest Test Template", created.Title.Raw)
		require.NotNil(t, created.Content)
		assert.Contains(t, created.Content.Raw, "Custom template content for testing")

		templateID := created.ID

		// Ensure cleanup if anything fails later
		defer func() {
			_, _ = client.Templates().Delete(templateID).Force().Do()
		}()

		// 2. Retrieve the template
		retrieved, err := client.Templates().Retrieve(templateID).
			ContextEdit().
			Do()
		require.NoError(t, err)
		require.NotNil(t, retrieved)
		assert.Equal(t, templateID, retrieved.ID)
		assert.Equal(t, slug, retrieved.Slug)
		assert.Equal(t, "Created by gowprest integration tests", retrieved.Description)

		// Sparse field retrieval
		sparseTmpl, err := client.Templates().Retrieve(templateID).
			Fields("id", "slug").
			Do()
		require.NoError(t, err)
		require.NotNil(t, sparseTmpl)
		assert.Equal(t, templateID, sparseTmpl.ID)
		assert.Equal(t, slug, sparseTmpl.Slug)
		assert.Empty(t, sparseTmpl.Description)

		// 3. Update the template
		updated, err := client.Templates().Update(templateID).
			Title("Updated Gowprest Template Title").
			Description("Updated template description").
			Do()
		require.NoError(t, err)
		require.NotNil(t, updated)
		assert.Equal(t, templateID, updated.ID)
		require.NotNil(t, updated.Title)
		assert.Equal(t, "Updated Gowprest Template Title", updated.Title.Raw)
		assert.Equal(t, "Updated template description", updated.Description)

		// 4. Delete the template with Force
		deleted, err := client.Templates().Delete(templateID).
			Force().
			Do()
		require.NoError(t, err)
		require.NotNil(t, deleted)
		assert.Equal(t, templateID, deleted.ID)

		// Verify deletion
		_, err = client.Templates().Retrieve(templateID).Do()
		require.Error(t, err)
		var wpErr *gowprest.WPRestError
		require.True(t, errors.As(err, &wpErr))
		assert.Equal(t, 404, wpErr.Data.Status)
		assert.Equal(t, "rest_template_not_found", wpErr.Code)
	})

	t.Run("TemplatesErrorHandling", func(t *testing.T) {
		// Nonexistent template
		_, err := client.Templates().Retrieve("nonexistent//template-slug").Do()
		require.Error(t, err)
		var wpErr *gowprest.WPRestError
		require.True(t, errors.As(err, &wpErr))
		assert.Equal(t, 404, wpErr.Data.Status)
		assert.Equal(t, "rest_template_not_found", wpErr.Code)

		// Unauthorized request
		unauthClient := gowprest.NewClient(blogUrl)
		defer unauthClient.Close()

		_, err = unauthClient.Templates().List().Do()
		require.Error(t, err)
		var unauthErr *gowprest.WPRestError
		require.True(t, errors.As(err, &unauthErr))
		assert.Equal(t, 401, unauthErr.Data.Status)
		assert.Equal(t, "rest_cannot_manage_templates", unauthErr.Code)
	})
}
