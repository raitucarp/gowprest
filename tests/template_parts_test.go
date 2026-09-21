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

func TestTemplateParts(t *testing.T) {
	blogUrl := os.Getenv("BLOG_URL")
	username := os.Getenv("BLOG_USERNAME")
	password := os.Getenv("BLOG_APP_PASSWORD")

	if blogUrl == "" || username == "" || password == "" {
		t.Skip("Skipping test; BLOG_URL, BLOG_USERNAME, or BLOG_APP_PASSWORD not set")
	}

	client := gowprest.NewClient(blogUrl).WithBasicAuth(username, password)
	defer client.Close()

	t.Run("TemplatePartsAliases", func(t *testing.T) {
		assert.NotNil(t, client.TemplateParts())
		assert.NotNil(t, client.TemplatePart())

		assert.NotNil(t, client.TemplateParts().List())
		assert.NotNil(t, client.TemplateParts().Retrieve("twentytwentyfive//footer-newsletter"))
		assert.NotNil(t, client.TemplateParts().Create())
		assert.NotNil(t, client.TemplateParts().Update("twentytwentyfive//footer-newsletter"))
		assert.NotNil(t, client.TemplateParts().Delete("twentytwentyfive//footer-newsletter"))
	})

	t.Run("ListTemplateParts", func(t *testing.T) {
		parts, err := client.TemplateParts().List().
			ContextEdit().
			Do()
		require.NoError(t, err)
		assert.NotEmpty(t, parts)

		firstPart := parts[0]
		assert.NotEmpty(t, firstPart.ID)
		assert.NotEmpty(t, firstPart.Slug)
		assert.NotEmpty(t, firstPart.Theme)
		assert.NotEmpty(t, firstPart.Type)
		assert.NotEmpty(t, firstPart.TitleString())
		assert.NotEmpty(t, firstPart.ContentString())

		// ContextView & ContextEmbed
		_, err = client.TemplateParts().List().ContextView().Do()
		require.NoError(t, err)

		_, err = client.TemplateParts().List().ContextEmbed().Do()
		require.NoError(t, err)

		// Filter by Area
		footerParts, err := client.TemplateParts().List().
			Area("footer").
			Do()
		require.NoError(t, err)
		assert.NotEmpty(t, footerParts)
		for _, fp := range footerParts {
			assert.Equal(t, "footer", fp.Area)
		}

		// Sparse fields
		sparseList, err := client.TemplateParts().List().
			Fields("id", "slug", "area").
			Do()
		require.NoError(t, err)
		assert.NotEmpty(t, sparseList)
		assert.NotEmpty(t, sparseList[0].ID)
		assert.NotEmpty(t, sparseList[0].Slug)
		assert.NotEmpty(t, sparseList[0].Area)
		assert.Empty(t, sparseList[0].Description)
	})

	t.Run("TemplatePartCRUDWorkflow", func(t *testing.T) {
		slug := "gowprest-test-footer-part"

		// 1. Create custom template part
		created, err := client.TemplateParts().Create().
			Slug(slug).
			Title("Gowprest Test Footer Part").
			Content("<!-- wp:paragraph --><p>Custom Footer Content</p><!-- /wp:paragraph -->").
			Area("footer").
			Description("Created by gowprest integration tests").
			StatusPublish().
			Do()
		require.NoError(t, err)
		require.NotNil(t, created)
		assert.NotEmpty(t, created.ID)
		assert.Equal(t, slug, created.Slug)
		assert.Equal(t, "footer", created.Area)
		assert.Equal(t, "custom", created.Source)
		assert.Equal(t, "Gowprest Test Footer Part", created.TitleString())
		assert.Contains(t, created.ContentString(), "Custom Footer Content")

		partID := created.ID
		defer func() {
			_, _ = client.TemplateParts().Delete(partID).Force().Do()
		}()

		// 2. Retrieve template part
		retrieved, err := client.TemplateParts().Retrieve(partID).
			ContextEdit().
			Do()
		require.NoError(t, err)
		require.NotNil(t, retrieved)
		assert.Equal(t, partID, retrieved.ID)
		assert.Equal(t, slug, retrieved.Slug)
		assert.Equal(t, "footer", retrieved.Area)
		assert.Equal(t, "Created by gowprest integration tests", retrieved.Description)

		// Sparse retrieve
		sparsePart, err := client.TemplateParts().Retrieve(partID).
			Fields("id", "slug", "area").
			Do()
		require.NoError(t, err)
		require.NotNil(t, sparsePart)
		assert.Equal(t, partID, sparsePart.ID)
		assert.Equal(t, slug, sparsePart.Slug)
		assert.Equal(t, "footer", sparsePart.Area)
		assert.Empty(t, sparsePart.Description)

		// 3. Update template part
		updated, err := client.TemplateParts().Update(partID).
			Title("Updated Gowprest Footer Title").
			Description("Updated template part description").
			Do()
		require.NoError(t, err)
		require.NotNil(t, updated)
		assert.Equal(t, partID, updated.ID)
		assert.Equal(t, "Updated Gowprest Footer Title", updated.TitleString())
		assert.Equal(t, "Updated template part description", updated.Description)

		// 4. Delete template part with Force
		deleted, err := client.TemplateParts().Delete(partID).
			Force().
			Do()
		require.NoError(t, err)
		require.NotNil(t, deleted)
		assert.Equal(t, partID, deleted.ID)

		// Verify deletion
		_, err = client.TemplateParts().Retrieve(partID).Do()
		require.Error(t, err)
		var wpErr *gowprest.WPRestError
		require.True(t, errors.As(err, &wpErr))
		assert.Equal(t, 404, wpErr.Data.Status)
		assert.Equal(t, "rest_template_not_found", wpErr.Code)
	})

	t.Run("TemplatePartsErrorHandling", func(t *testing.T) {
		// Nonexistent template part
		_, err := client.TemplateParts().Retrieve("nonexistent//template-part-slug").Do()
		require.Error(t, err)
		var wpErr *gowprest.WPRestError
		require.True(t, errors.As(err, &wpErr))
		assert.Equal(t, 404, wpErr.Data.Status)
		assert.Equal(t, "rest_template_not_found", wpErr.Code)

		// Unauthorized request
		unauthClient := gowprest.NewClient(blogUrl)
		defer unauthClient.Close()

		_, err = unauthClient.TemplateParts().List().Do()
		require.Error(t, err)
		var unauthErr *gowprest.WPRestError
		require.True(t, errors.As(err, &unauthErr))
		assert.Equal(t, 401, unauthErr.Data.Status)
		assert.Equal(t, "rest_cannot_manage_templates", unauthErr.Code)
	})
}
