package tests

import (
	"os"
	"testing"

	_ "github.com/joho/godotenv/autoload"
	"github.com/raitucarp/gowprest"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestTaxonomies(t *testing.T) {
	blogUrl := os.Getenv("BLOG_URL")
	username := os.Getenv("BLOG_USERNAME")
	password := os.Getenv("BLOG_APP_PASSWORD")

	if blogUrl == "" || username == "" || password == "" {
		t.Skip("Skipping test; BLOG_URL, BLOG_USERNAME, or BLOG_APP_PASSWORD not set")
	}

	client := gowprest.NewClient(blogUrl).WithBasicAuth(username, password)
	defer client.Close()

	taxonomiesAPI := client.Taxonomies()

	t.Run("ListTaxonomies_Default", func(t *testing.T) {
		taxonomies, err := taxonomiesAPI.List().Do()
		require.NoError(t, err)
		assert.GreaterOrEqual(t, len(taxonomies), 2, "Should return at least category and post_tag")

		assert.Contains(t, taxonomies, gowprest.TaxonomyCategory)
		assert.Contains(t, taxonomies, gowprest.TaxonomyPostTag)

		cat := taxonomies[gowprest.TaxonomyCategory]
		assert.Equal(t, "Categories", cat.Name)
		assert.Equal(t, "categories", cat.RestBase)
		assert.True(t, cat.Hierarchical)
		assert.Contains(t, cat.Types, "post")

		tag := taxonomies[gowprest.TaxonomyPostTag]
		assert.Equal(t, "Tags", tag.Name)
		assert.Equal(t, "tags", tag.RestBase)
		assert.False(t, tag.Hierarchical)
		assert.Contains(t, tag.Types, "post")
	})

	t.Run("ListTaxonomies_FilterByType", func(t *testing.T) {
		postTaxonomies, err := taxonomiesAPI.List().Type("post").Do()
		require.NoError(t, err)
		assert.Contains(t, postTaxonomies, gowprest.TaxonomyCategory)
		assert.Contains(t, postTaxonomies, gowprest.TaxonomyPostTag)

		pageTaxonomies, err := taxonomiesAPI.List().Type("page").Do()
		require.NoError(t, err)
		// Default WordPress pages do not have tags
		assert.NotContains(t, pageTaxonomies, gowprest.TaxonomyPostTag)
	})

	t.Run("ListTaxonomies_Contexts", func(t *testing.T) {
		// Context View
		viewTaxes, err := taxonomiesAPI.List().ContextView().Do()
		require.NoError(t, err)
		assert.NotEmpty(t, viewTaxes)
		assert.Empty(t, viewTaxes[gowprest.TaxonomyCategory].Capabilities.ManageTerms)

		// Context Edit
		editTaxes, err := taxonomiesAPI.List().ContextEdit().Do()
		require.NoError(t, err)
		assert.NotEmpty(t, editTaxes)
		assert.NotEmpty(t, editTaxes[gowprest.TaxonomyCategory].Capabilities.ManageTerms)
		assert.NotEmpty(t, editTaxes[gowprest.TaxonomyCategory].Labels.SingularName)
		assert.True(t, editTaxes[gowprest.TaxonomyCategory].Visibility.Public)

		// Context Embed
		embedTaxes, err := taxonomiesAPI.List().ContextEmbed().Do()
		require.NoError(t, err)
		assert.NotEmpty(t, embedTaxes)
		assert.Equal(t, "Categories", embedTaxes[gowprest.TaxonomyCategory].Name)
	})

	t.Run("ListTaxonomies_GlobalParameters", func(t *testing.T) {
		// Embed
		embedTaxes, err := taxonomiesAPI.List().Embed().Do()
		require.NoError(t, err)
		assert.NotEmpty(t, embedTaxes)
		assert.Contains(t, embedTaxes, gowprest.TaxonomyCategory)
	})

	t.Run("RetrieveTaxonomy_Success", func(t *testing.T) {
		cat, err := taxonomiesAPI.Retrieve(gowprest.TaxonomyCategory).Do()
		require.NoError(t, err)
		assert.Equal(t, "category", cat.Slug)
		assert.Equal(t, "Categories", cat.Name)
		assert.True(t, cat.Hierarchical)
		assert.Equal(t, "categories", cat.RestBase)
		assert.Contains(t, cat.Types, "post")

		tag, err := taxonomiesAPI.Retrieve(gowprest.TaxonomyPostTag).Do()
		require.NoError(t, err)
		assert.Equal(t, "post_tag", tag.Slug)
		assert.Equal(t, "Tags", tag.Name)
		assert.False(t, tag.Hierarchical)
		assert.Equal(t, "tags", tag.RestBase)
		assert.Contains(t, tag.Types, "post")
	})

	t.Run("RetrieveTaxonomy_ContextEdit", func(t *testing.T) {
		cat, err := taxonomiesAPI.Retrieve(gowprest.TaxonomyCategory).Context("edit").Do()
		require.NoError(t, err)
		assert.Equal(t, "category", cat.Slug)
		assert.NotEmpty(t, cat.Capabilities.ManageTerms)
		assert.NotEmpty(t, cat.Capabilities.EditTerms)
		assert.NotEmpty(t, cat.Capabilities.DeleteTerms)
		assert.NotEmpty(t, cat.Capabilities.AssignTerms)
		assert.NotEmpty(t, cat.Labels.Name)
		assert.NotEmpty(t, cat.Labels.SingularName)
		assert.True(t, cat.Visibility.Public)
	})

	t.Run("RetrieveTaxonomy_GlobalParameters", func(t *testing.T) {
		cat, err := taxonomiesAPI.Retrieve(gowprest.TaxonomyCategory).
			Fields("name", "slug").
			Embed().
			Do()
		require.NoError(t, err)
		assert.Equal(t, "category", cat.Slug)
		assert.Equal(t, "Categories", cat.Name)
		assert.Empty(t, cat.RestBase)
	})

	t.Run("RetrieveTaxonomy_NotFound", func(t *testing.T) {
		tax, err := taxonomiesAPI.Retrieve("non_existent_tax_xyz").Do()
		assert.Error(t, err)
		assert.Nil(t, tax)
		var wpErr *gowprest.WPRestError
		if assert.ErrorAs(t, err, &wpErr) {
			assert.Equal(t, 404, wpErr.Data.Status)
			assert.Equal(t, "rest_taxonomy_invalid", wpErr.Code)
		}
	})
}
