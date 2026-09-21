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

func TestBlockPatternCategories(t *testing.T) {
	blogUrl := os.Getenv("BLOG_URL")
	username := os.Getenv("BLOG_USERNAME")
	password := os.Getenv("BLOG_APP_PASSWORD")

	if blogUrl == "" || username == "" || password == "" {
		t.Skip("Skipping test; BLOG_URL, BLOG_USERNAME, or BLOG_APP_PASSWORD not set")
	}

	client := gowprest.NewClient(blogUrl).WithBasicAuth(username, password)
	defer client.Close()

	t.Run("ServiceAliases", func(t *testing.T) {
		assert.NotNil(t, client.BlockPatternCategories())
		assert.NotNil(t, client.BlockPatterns())
		assert.NotNil(t, client.BlockPatterns().Categories())
		assert.NotNil(t, client.BlockPatternCategories().List())
		assert.NotNil(t, client.BlockPatternCategories().Retrieve("buttons"))
	})

	t.Run("ListCategories", func(t *testing.T) {
		t.Run("Default", func(t *testing.T) {
			categories, err := client.BlockPatternCategories().List().Do()
			require.NoError(t, err)
			require.NotEmpty(t, categories)

			var foundButtons bool
			for _, cat := range categories {
				assert.NotEmpty(t, cat.Name)
				assert.NotEmpty(t, cat.Label)
				if cat.Name == "buttons" {
					foundButtons = true
					assert.Equal(t, "Buttons", cat.Label)
					assert.NotEmpty(t, cat.Description)
				}
			}
			assert.True(t, foundButtons, "expected 'buttons' category to be registered")
		})

		t.Run("Contexts", func(t *testing.T) {
			for _, ctx := range []string{"view", "edit", "embed"} {
				categories, err := client.BlockPatternCategories().
					List().
					Context(ctx).
					Do()
				require.NoError(t, err)
				assert.NotEmpty(t, categories)
			}

			// Fluent context helpers
			catView, err := client.BlockPatternCategories().List().ContextView().Do()
			require.NoError(t, err)
			assert.NotEmpty(t, catView)

			catEdit, err := client.BlockPatternCategories().List().ContextEdit().Do()
			require.NoError(t, err)
			assert.NotEmpty(t, catEdit)

			catEmbed, err := client.BlockPatternCategories().List().ContextEmbed().Do()
			require.NoError(t, err)
			assert.NotEmpty(t, catEmbed)
		})

		t.Run("Fields", func(t *testing.T) {
			categories, err := client.BlockPatternCategories().
				List().
				Fields("name", "label").
				Do()
			require.NoError(t, err)
			require.NotEmpty(t, categories)

			for _, cat := range categories {
				assert.NotEmpty(t, cat.Name)
				assert.NotEmpty(t, cat.Label)
				assert.Empty(t, cat.Description, "description should be empty when only name,label are requested")
			}
		})
	})

	t.Run("RetrieveCategory", func(t *testing.T) {
		t.Run("ExistingCategory", func(t *testing.T) {
			cat, err := client.BlockPatternCategories().Retrieve("buttons").Do()
			require.NoError(t, err)
			require.NotNil(t, cat)
			assert.Equal(t, "buttons", cat.Name)
			assert.Equal(t, "Buttons", cat.Label)
			assert.NotEmpty(t, cat.Description)

			// Contexts on Retrieve
			catView, err := client.BlockPatternCategories().Retrieve("buttons").ContextView().Do()
			require.NoError(t, err)
			assert.Equal(t, "buttons", catView.Name)

			catEdit, err := client.BlockPatternCategories().Retrieve("buttons").ContextEdit().Do()
			require.NoError(t, err)
			assert.Equal(t, "buttons", catEdit.Name)

			catEmbed, err := client.BlockPatternCategories().Retrieve("buttons").ContextEmbed().Do()
			require.NoError(t, err)
			assert.Equal(t, "buttons", catEmbed.Name)
		})

		t.Run("WithFields", func(t *testing.T) {
			cat, err := client.BlockPatternCategories().
				Retrieve("buttons").
				Fields("name").
				Do()
			require.NoError(t, err)
			require.NotNil(t, cat)
			assert.Equal(t, "buttons", cat.Name)
			assert.Empty(t, cat.Label)
			assert.Empty(t, cat.Description)
		})

		t.Run("NotFound", func(t *testing.T) {
			cat, err := client.BlockPatternCategories().Retrieve("non-existent-pattern-category-xyz").Do()
			require.Error(t, err)
			assert.Nil(t, cat)

			var wpErr *gowprest.WPRestError
			if errors.As(err, &wpErr) {
				assert.Equal(t, "rest_block_pattern_category_not_found", wpErr.Code)
				assert.Equal(t, 404, wpErr.Data.Status)
			} else {
				t.Fatalf("expected WPRestError, got: %v", err)
			}
		})
	})

	t.Run("AlternativeEntrypoint", func(t *testing.T) {
		categories, err := client.BlockPatterns().Categories().List().Do()
		require.NoError(t, err)
		require.NotEmpty(t, categories)

		cat, err := client.BlockPatterns().Categories().Retrieve("banner").Do()
		require.NoError(t, err)
		require.NotNil(t, cat)
		assert.Equal(t, "banner", cat.Name)
	})

	t.Run("UnauthorizedAccess", func(t *testing.T) {
		unauthClient := gowprest.NewClient(blogUrl)
		defer unauthClient.Close()

		_, err := unauthClient.BlockPatternCategories().List().Do()
		require.Error(t, err)

		var wpErr *gowprest.WPRestError
		if errors.As(err, &wpErr) {
			assert.Equal(t, "rest_cannot_view", wpErr.Code)
			assert.Equal(t, 401, wpErr.Data.Status)
		} else {
			t.Fatalf("expected WPRestError, got: %v", err)
		}
	})
}
