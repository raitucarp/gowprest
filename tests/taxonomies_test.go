package tests

import (
	"os"
	"testing"

	_ "github.com/joho/godotenv/autoload"
	"github.com/raitucarp/gowprest"
	"github.com/stretchr/testify/assert"
)

func TestListTaxonomies(t *testing.T) {
	client := gowprest.NewClient(blogUrl)
	defer client.Close()

	t.Run("default taxonomy list", func(t *testing.T) {
		taxonomies, err := client.Taxonomies().List().Do()

		assert.Equal(t, nil, err)
		assert.GreaterOrEqual(t, len(taxonomies), 2, "Should return at least category and post_tag")
		assert.Contains(t, taxonomies, "category")
		assert.Contains(t, taxonomies, "post_tag")
	})

	t.Run("list with context edit", func(t *testing.T) {
		clientWithAuth := gowprest.NewClient(blogUrl).
			WithBasicAuth(
				os.Getenv("BLOG_USERNAME"),
				os.Getenv("BLOG_APP_PASSWORD"),
			)
		defer clientWithAuth.Close()

		taxonomies, err := clientWithAuth.Taxonomies().List().ContextEdit().Do()

		assert.Equal(t, nil, err)
		assert.GreaterOrEqual(t, len(taxonomies), 2)
	})

	t.Run("list by type", func(t *testing.T) {
		taxonomies, err := client.Taxonomies().List().Type("post").Do()

		assert.Equal(t, nil, err)
		assert.GreaterOrEqual(t, len(taxonomies), 2)
	})
}

func TestRetrieveTaxonomy(t *testing.T) {
	client := gowprest.NewClient(blogUrl)
	defer client.Close()

	t.Run("retrieve category", func(t *testing.T) {
		taxonomy, err := client.Taxonomies().Retrieve("category").Do()

		assert.Equal(t, nil, err)
		assert.Equal(t, "category", taxonomy.Slug)
		assert.Equal(t, "Categories", taxonomy.Name)
	})

	t.Run("retrieve non-existent taxonomy", func(t *testing.T) {
		taxonomy, err := client.Taxonomies().Retrieve("non-existent").Do()

		assert.NotNil(t, err)
		assert.Nil(t, taxonomy)
	})

	t.Run("retrieve with context edit", func(t *testing.T) {
		clientWithAuth := gowprest.NewClient(blogUrl).
			WithBasicAuth(
				os.Getenv("BLOG_USERNAME"),
				os.Getenv("BLOG_APP_PASSWORD"),
			)
		defer clientWithAuth.Close()

		taxonomy, err := clientWithAuth.Taxonomies().Retrieve("category").ContextEdit().Do()

		assert.Equal(t, nil, err)
		assert.Equal(t, "category", taxonomy.Slug)
	})
}
