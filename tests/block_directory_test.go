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

func TestBlockDirectory(t *testing.T) {
	blogUrl := os.Getenv("BLOG_URL")
	username := os.Getenv("BLOG_USERNAME")
	password := os.Getenv("BLOG_APP_PASSWORD")

	if blogUrl == "" || username == "" || password == "" {
		t.Skip("Skipping test; BLOG_URL, BLOG_USERNAME, or BLOG_APP_PASSWORD not set")
	}

	client := gowprest.NewClient(blogUrl).WithBasicAuth(username, password)
	defer client.Close()

	t.Run("BlockDirectoryAliases", func(t *testing.T) {
		assert.NotNil(t, client.BlockDirectory())
		assert.NotNil(t, client.BlockDirectoryItems())
		assert.NotNil(t, client.DirectorySearch())
		assert.NotNil(t, client.BlockDirectory().Search("slider"))
		assert.NotNil(t, client.BlockDirectory().List("slider"))
	})

	t.Run("SearchBlockDirectory_Default", func(t *testing.T) {
		items, err := client.BlockDirectory().Search("slider").Do()
		require.NoError(t, err)
		require.NotEmpty(t, items)

		first := items[0]
		assert.NotEmpty(t, first.Name)
		assert.NotEmpty(t, first.Title)
		assert.NotEmpty(t, first.ID)
		assert.NotEmpty(t, first.Author)
		assert.NotEmpty(t, first.Icon)
		assert.NotEmpty(t, first.LastUpdated)
		assert.NotEmpty(t, first.HumanizedUpdated)
	})

	t.Run("SearchBlockDirectory_Pagination", func(t *testing.T) {
		page1, err := client.BlockDirectory().
			Search("slider").
			Page(1).
			PerPage(2).
			Do()
		require.NoError(t, err)
		require.Len(t, page1, 2)

		page2, err := client.BlockDirectory().
			Search("slider").
			Page(2).
			PerPage(2).
			Do()
		require.NoError(t, err)
		require.Len(t, page2, 2)

		assert.NotEqual(t, page1[0].ID, page2[0].ID)
	})

	t.Run("SearchBlockDirectory_Context", func(t *testing.T) {
		items, err := client.BlockDirectory().
			Search("slider").
			ContextView().
			PerPage(3).
			Do()
		require.NoError(t, err)
		assert.NotEmpty(t, items)
	})

	t.Run("SearchBlockDirectory_Fields", func(t *testing.T) {
		items, err := client.BlockDirectory().
			Search("slider").
			Fields("name", "title", "id").
			PerPage(3).
			Do()
		require.NoError(t, err)
		require.NotEmpty(t, items)

		for _, item := range items {
			assert.NotEmpty(t, item.Name)
			assert.NotEmpty(t, item.Title)
			assert.NotEmpty(t, item.ID)
			assert.Empty(t, item.Description)
			assert.Empty(t, item.Author)
		}
	})

	t.Run("SearchBlockDirectory_MissingTerm", func(t *testing.T) {
		items, err := client.BlockDirectory().Search().Do()
		require.Error(t, err)
		assert.Nil(t, items)

		var wpErr *gowprest.WPRestError
		if assert.True(t, errors.As(err, &wpErr)) {
			assert.Equal(t, "rest_missing_callback_param", wpErr.Code)
			assert.Equal(t, 400, wpErr.Data.Status)
		}
	})
}
