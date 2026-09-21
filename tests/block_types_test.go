package tests

import (
	"errors"
	"os"
	"strings"
	"testing"

	_ "github.com/joho/godotenv/autoload"
	"github.com/raitucarp/gowprest"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestBlockTypes(t *testing.T) {
	blogUrl := os.Getenv("BLOG_URL")
	username := os.Getenv("BLOG_USERNAME")
	password := os.Getenv("BLOG_APP_PASSWORD")

	if blogUrl == "" || username == "" || password == "" {
		t.Skip("Skipping test; BLOG_URL, BLOG_USERNAME, or BLOG_APP_PASSWORD not set")
	}

	client := gowprest.NewClient(blogUrl).WithBasicAuth(username, password)
	defer client.Close()

	t.Run("ListBlockTypes_Default", func(t *testing.T) {
		blocks, err := client.BlockTypes().List().Do()
		require.NoError(t, err)
		assert.NotEmpty(t, blocks)

		var foundParagraph bool
		for _, b := range blocks {
			if b.Name == "core/paragraph" {
				foundParagraph = true
				assert.NotEmpty(t, b.Title)
				assert.Equal(t, gowprest.BlockCategoryText, b.Category)
				break
			}
		}
		assert.True(t, foundParagraph, "expected core/paragraph block to be in the list")
	})

	t.Run("ListBlockTypes_NamespaceFilter", func(t *testing.T) {
		blocks, err := client.BlockTypes().List().Namespace("core").Do()
		require.NoError(t, err)
		assert.NotEmpty(t, blocks)

		for _, b := range blocks {
			assert.True(t, strings.HasPrefix(b.Name, "core/"), "block name %q should start with 'core/'", b.Name)
		}

		// Also test ByNamespace helper method
		blocksByNs, err := client.BlockTypes().ByNamespace("core").Do()
		require.NoError(t, err)
		assert.Equal(t, len(blocks), len(blocksByNs))
	})

	t.Run("ListBlockTypes_Contexts", func(t *testing.T) {
		viewBlocks, err := client.BlockTypes().List().ContextView().Do()
		require.NoError(t, err)
		assert.NotEmpty(t, viewBlocks)

		editBlocks, err := client.BlockTypes().List().ContextEdit().Do()
		require.NoError(t, err)
		assert.NotEmpty(t, editBlocks)

		embedBlocks, err := client.BlockTypes().List().ContextEmbed().Do()
		require.NoError(t, err)
		assert.NotEmpty(t, embedBlocks)
	})

	t.Run("ListBlockTypes_Fields", func(t *testing.T) {
		blocks, err := client.BlockTypes().List().
			Fields("name", "title").
			Do()
		require.NoError(t, err)
		require.NotEmpty(t, blocks)

		for _, b := range blocks {
			assert.NotEmpty(t, b.Name)
			if b.Title == "" {
				t.Logf("Block with empty title: %s", b.Name)
			}
			assert.Empty(t, b.Description)
			assert.Empty(t, b.Category)
		}
		// Ensure at least core/paragraph had title populated
		var paragraphTitle string
		for _, b := range blocks {
			if b.Name == "core/paragraph" {
				paragraphTitle = b.Title
			}
		}
		assert.NotEmpty(t, paragraphTitle)
	})

	t.Run("RetrieveBlockType_CombinedIdentifier", func(t *testing.T) {
		block, err := client.BlockTypes().Retrieve("core/paragraph").Do()
		require.NoError(t, err)
		require.NotNil(t, block)

		assert.Equal(t, "core/paragraph", block.Name)
		assert.NotEmpty(t, block.Title)
		assert.Equal(t, gowprest.BlockCategoryText, block.Category)
		assert.NotEmpty(t, block.Description)
	})

	t.Run("RetrieveBlockType_SeparateNamespaceAndName", func(t *testing.T) {
		block, err := client.BlockTypes().Retrieve("core", "paragraph").Do()
		require.NoError(t, err)
		require.NotNil(t, block)

		assert.Equal(t, "core/paragraph", block.Name)
		assert.NotEmpty(t, block.Title)
		assert.Equal(t, gowprest.BlockCategoryText, block.Category)
	})

	t.Run("RetrieveBlockType_ContextAndFields", func(t *testing.T) {
		block, err := client.BlockTypes().Retrieve("core/paragraph").
			ContextEdit().
			Fields("name", "title").
			Do()
		require.NoError(t, err)
		require.NotNil(t, block)

		assert.Equal(t, "core/paragraph", block.Name)
		assert.NotEmpty(t, block.Title)
		assert.Empty(t, block.Description)
		assert.Empty(t, block.Category)
	})

	t.Run("RetrieveBlockType_NotFound", func(t *testing.T) {
		block, err := client.BlockTypes().Retrieve("nonexistent/block-xyz").Do()
		require.Error(t, err)
		assert.Nil(t, block)

		var wpErr *gowprest.WPRestError
		if assert.True(t, errors.As(err, &wpErr)) {
			assert.Equal(t, "rest_block_type_invalid", wpErr.Code)
			assert.Equal(t, 404, wpErr.Data.Status)
		}
	})
}
