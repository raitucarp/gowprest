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

func TestBlockPatterns(t *testing.T) {
	blogUrl := os.Getenv("BLOG_URL")
	username := os.Getenv("BLOG_USERNAME")
	password := os.Getenv("BLOG_APP_PASSWORD")

	if blogUrl == "" || username == "" || password == "" {
		t.Skip("Skipping test; BLOG_URL, BLOG_USERNAME, or BLOG_APP_PASSWORD not set")
	}

	client := gowprest.NewClient(blogUrl).WithBasicAuth(username, password)
	defer client.Close()

	t.Run("ServiceAliases", func(t *testing.T) {
		assert.NotNil(t, client.BlockPatterns())
		assert.NotNil(t, client.BlockPatterns().Patterns())
		assert.NotNil(t, client.BlockPatterns().List())
		assert.NotNil(t, client.BlockPatterns().Retrieve("core/query-standard-posts"))
	})

	t.Run("ListPatterns", func(t *testing.T) {
		t.Run("Default", func(t *testing.T) {
			patterns, err := client.BlockPatterns().List().Do()
			require.NoError(t, err)
			require.NotEmpty(t, patterns)

			var foundQueryStandard bool
			for _, p := range patterns {
				assert.NotEmpty(t, p.Name)
				assert.NotEmpty(t, p.Title)
				if p.Name == "core/query-standard-posts" {
					foundQueryStandard = true
					assert.Equal(t, "Standard", p.Title)
					assert.Equal(t, gowprest.BlockPatternSourceCore, p.Source)
					assert.NotEmpty(t, p.Content)
					assert.Contains(t, p.Categories, "posts")
					assert.Contains(t, p.BlockTypes, "core/query")
				}
			}
			assert.True(t, foundQueryStandard, "expected 'core/query-standard-posts' pattern to be registered")
		})

		t.Run("Contexts", func(t *testing.T) {
			for _, ctx := range []string{"view", "edit", "embed"} {
				patterns, err := client.BlockPatterns().
					List().
					Context(ctx).
					Do()
				require.NoError(t, err)
				assert.NotEmpty(t, patterns)
			}

			// Fluent context helpers
			patView, err := client.BlockPatterns().List().ContextView().Do()
			require.NoError(t, err)
			assert.NotEmpty(t, patView)

			patEdit, err := client.BlockPatterns().List().ContextEdit().Do()
			require.NoError(t, err)
			assert.NotEmpty(t, patEdit)

			patEmbed, err := client.BlockPatterns().List().ContextEmbed().Do()
			require.NoError(t, err)
			assert.NotEmpty(t, patEmbed)
		})

		t.Run("Fields", func(t *testing.T) {
			patterns, err := client.BlockPatterns().
				List().
				Fields("name", "title").
				Do()
			require.NoError(t, err)
			require.NotEmpty(t, patterns)

			for _, p := range patterns {
				assert.NotEmpty(t, p.Name)
				assert.NotEmpty(t, p.Title)
				assert.Empty(t, p.Content, "content should be empty when only name,title are requested")
			}
		})
	})

	t.Run("RetrievePattern", func(t *testing.T) {
		t.Run("ExistingPattern", func(t *testing.T) {
			p, err := client.BlockPatterns().Retrieve("core/query-standard-posts").Do()
			require.NoError(t, err)
			require.NotNil(t, p)
			assert.Equal(t, "core/query-standard-posts", p.Name)
			assert.Equal(t, "Standard", p.Title)
			assert.NotEmpty(t, p.Content)
			assert.Equal(t, gowprest.BlockPatternSourceCore, p.Source)

			// Contexts on Retrieve
			pView, err := client.BlockPatterns().Retrieve("core/query-standard-posts").ContextView().Do()
			require.NoError(t, err)
			assert.Equal(t, "core/query-standard-posts", pView.Name)

			pEdit, err := client.BlockPatterns().Retrieve("core/query-standard-posts").ContextEdit().Do()
			require.NoError(t, err)
			assert.Equal(t, "core/query-standard-posts", pEdit.Name)

			pEmbed, err := client.BlockPatterns().Retrieve("core/query-standard-posts").ContextEmbed().Do()
			require.NoError(t, err)
			assert.Equal(t, "core/query-standard-posts", pEmbed.Name)
		})

		t.Run("WithFields", func(t *testing.T) {
			p, err := client.BlockPatterns().
				Retrieve("core/query-standard-posts").
				Fields("name", "title").
				Do()
			require.NoError(t, err)
			require.NotNil(t, p)
			assert.Equal(t, "core/query-standard-posts", p.Name)
			assert.Equal(t, "Standard", p.Title)
			assert.Empty(t, p.Content)
		})

		t.Run("NotFound", func(t *testing.T) {
			p, err := client.BlockPatterns().Retrieve("non-existent-pattern-name-xyz").Do()
			require.Error(t, err)
			assert.Nil(t, p)

			var wpErr *gowprest.WPRestError
			if errors.As(err, &wpErr) {
				assert.Equal(t, "rest_block_pattern_not_found", wpErr.Code)
				assert.Equal(t, 404, wpErr.Data.Status)
			} else {
				t.Fatalf("expected WPRestError, got: %v", err)
			}
		})
	})

	t.Run("UnauthorizedAccess", func(t *testing.T) {
		unauthClient := gowprest.NewClient(blogUrl)
		defer unauthClient.Close()

		_, err := unauthClient.BlockPatterns().List().Do()
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
