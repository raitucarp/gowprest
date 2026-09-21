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

func TestBlockRenderer(t *testing.T) {
	blogUrl := os.Getenv("BLOG_URL")
	username := os.Getenv("BLOG_USERNAME")
	password := os.Getenv("BLOG_APP_PASSWORD")

	if blogUrl == "" || username == "" || password == "" {
		t.Skip("Skipping test; BLOG_URL, BLOG_USERNAME, or BLOG_APP_PASSWORD not set")
	}

	client := gowprest.NewClient(blogUrl).WithBasicAuth(username, password)
	defer client.Close()

	t.Run("BlockRendererAliases", func(t *testing.T) {
		assert.NotNil(t, client.BlockRenderer())
		assert.NotNil(t, client.RenderedBlocks())
		assert.NotNil(t, client.BlockRenderer().Render("core/calendar"))
		assert.NotNil(t, client.BlockRenderer().Get("core/calendar"))
		assert.NotNil(t, client.BlockRenderer().Create("core/calendar"))
	})

	t.Run("RenderBlock_Default", func(t *testing.T) {
		res, err := client.BlockRenderer().Render("core/calendar").Do()
		require.NoError(t, err)
		require.NotNil(t, res)
		assert.NotEmpty(t, res.Rendered)
		assert.Contains(t, res.Rendered, "wp-block-calendar")
	})

	t.Run("RenderBlock_SeparateIdentifier", func(t *testing.T) {
		res, err := client.BlockRenderer().Render("core", "calendar").Do()
		require.NoError(t, err)
		require.NotNil(t, res)
		assert.NotEmpty(t, res.Rendered)
		assert.Contains(t, res.Rendered, "wp-block-calendar")
	})

	t.Run("RenderBlock_WithAttributes", func(t *testing.T) {
		res, err := client.BlockRenderer().Render("core/archives").
			SetAttribute("showPostCounts", true).
			SetAttribute("displayAsDropdown", false).
			Do()
		require.NoError(t, err)
		require.NotNil(t, res)
		assert.NotEmpty(t, res.Rendered)
		assert.Contains(t, res.Rendered, "wp-block-archives")
	})

	t.Run("RenderBlock_WithPostID", func(t *testing.T) {
		post, err := client.Posts().Create().
			Title("Renderer Context Post").
			Content("<p>Context Post Body</p>").
			StatusPublish().
			Do()
		require.NoError(t, err)
		defer client.Posts().Delete(post.ID).Force().Do()

		res, err := client.BlockRenderer().Render("core/calendar").
			PostID(post.ID).
			Do()
		require.NoError(t, err)
		require.NotNil(t, res)
		assert.NotEmpty(t, res.Rendered)
	})

	t.Run("RenderBlock_GetMethod", func(t *testing.T) {
		res, err := client.BlockRenderer().Get("core/calendar").Get()
		require.NoError(t, err)
		require.NotNil(t, res)
		assert.NotEmpty(t, res.Rendered)
		assert.Contains(t, res.Rendered, "wp-block-calendar")
	})

	t.Run("RenderBlock_NotFound", func(t *testing.T) {
		res, err := client.BlockRenderer().Render("nonexistent/invalid-block-xyz").Do()
		require.Error(t, err)
		assert.Nil(t, res)

		var wpErr *gowprest.WPRestError
		if assert.True(t, errors.As(err, &wpErr)) {
			assert.Equal(t, "block_invalid", wpErr.Code)
			assert.Equal(t, 404, wpErr.Data.Status)
		}
	})
}
