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

func TestPatternDirectoryItems(t *testing.T) {
	blogUrl := os.Getenv("BLOG_URL")
	username := os.Getenv("BLOG_USERNAME")
	password := os.Getenv("BLOG_APP_PASSWORD")

	if blogUrl == "" || username == "" || password == "" {
		t.Skip("Skipping test; BLOG_URL, BLOG_USERNAME, or BLOG_APP_PASSWORD not set")
	}

	client := gowprest.NewClient(blogUrl).WithBasicAuth(username, password)
	defer client.Close()

	t.Run("PatternDirectoryAliases", func(t *testing.T) {
		assert.NotNil(t, client.PatternDirectory())
		assert.NotNil(t, client.PatternDirectoryItems())

		assert.NotNil(t, client.PatternDirectory().List())
		assert.NotNil(t, client.PatternDirectory().Search("hero"))
	})

	t.Run("ListPatternDirectoryItems", func(t *testing.T) {
		// 1. Default list
		items, err := client.PatternDirectory().List().
			PerPage(2).
			Do()
		require.NoError(t, err)
		assert.NotEmpty(t, items)
		assert.NotEmpty(t, items[0].ID)
		assert.NotEmpty(t, items[0].Title)
		assert.NotEmpty(t, items[0].Content)

		// 2. Search
		searchItems, err := client.PatternDirectory().Search("hero").
			PerPage(2).
			Do()
		require.NoError(t, err)
		assert.NotEmpty(t, searchItems)

		// 3. Pagination
		pageItems, err := client.PatternDirectory().List().
			Page(1).
			PerPage(1).
			Do()
		require.NoError(t, err)
		assert.Len(t, pageItems, 1)

		// 4. Ordering
		orderedItems, err := client.PatternDirectory().List().
			OrderDesc().
			OrderBy("date").
			PerPage(2).
			Do()
		require.NoError(t, err)
		assert.NotEmpty(t, orderedItems)

		// 5. Contexts
		editItems, err := client.PatternDirectory().List().
			ContextEdit().
			PerPage(1).
			Do()
		require.NoError(t, err)
		assert.NotEmpty(t, editItems)

		embedItems, err := client.PatternDirectory().List().
			ContextEmbed().
			PerPage(1).
			Do()
		require.NoError(t, err)
		assert.NotEmpty(t, embedItems)

		// 6. Fields
		fieldsItems, err := client.PatternDirectory().List().
			Fields("id", "title").
			PerPage(1).
			Do()
		require.NoError(t, err)
		require.NotEmpty(t, fieldsItems)
		assert.NotEmpty(t, fieldsItems[0].ID)
		assert.NotEmpty(t, fieldsItems[0].Title)
		assert.Empty(t, fieldsItems[0].Content)
	})

	t.Run("PatternDirectoryUnauthorized", func(t *testing.T) {
		unauthClient := gowprest.NewClient(blogUrl)
		defer unauthClient.Close()

		_, err := unauthClient.PatternDirectory().List().
			PerPage(1).
			Do()
		require.Error(t, err)
		var wpErr *gowprest.WPRestError
		require.True(t, errors.As(err, &wpErr))
		assert.Equal(t, 401, wpErr.Data.Status)
		assert.Equal(t, "rest_pattern_directory_cannot_view", wpErr.Code)
	})
}
