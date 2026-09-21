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

func TestGlobalStyles(t *testing.T) {
	blogUrl := os.Getenv("BLOG_URL")
	username := os.Getenv("BLOG_USERNAME")
	password := os.Getenv("BLOG_APP_PASSWORD")

	if blogUrl == "" || username == "" || password == "" {
		t.Skip("Skipping test; BLOG_URL, BLOG_USERNAME, or BLOG_APP_PASSWORD not set")
	}

	client := gowprest.NewClient(blogUrl).WithBasicAuth(username, password)
	defer client.Close()

	// 805 is the active user global styles post ID in the local Docker setup.
	activeGlobalStylesID := 805

	t.Run("ServiceAliases", func(t *testing.T) {
		assert.NotNil(t, client.GlobalStyles())
		assert.NotNil(t, client.GlobalStyles().Retrieve(activeGlobalStylesID))
		assert.NotNil(t, client.GlobalStyles().Update(activeGlobalStylesID))
		assert.NotNil(t, client.GlobalStyles().Themes("twentytwentyfive"))
		assert.NotNil(t, client.GlobalStyles().Themes("twentytwentyfive").Retrieve())
		assert.NotNil(t, client.GlobalStyles().Themes("twentytwentyfive").Variations())
	})

	t.Run("RetrieveGlobalStyles", func(t *testing.T) {
		t.Run("ActiveGlobalStyles", func(t *testing.T) {
			gs, err := client.GlobalStyles().Retrieve(activeGlobalStylesID).Do()
			require.NoError(t, err)
			require.NotNil(t, gs)
			assert.Equal(t, activeGlobalStylesID, gs.ID)
			require.NotNil(t, gs.Title)
			assert.NotEmpty(t, gs.Title.Rendered)
			assert.NotEmpty(t, gs.Links)
		})

		t.Run("Contexts", func(t *testing.T) {
			for _, ctx := range []string{"view", "edit", "embed"} {
				gs, err := client.GlobalStyles().
					Retrieve(activeGlobalStylesID).
					Context(ctx).
					Do()
				require.NoError(t, err)
				require.NotNil(t, gs)
				assert.Equal(t, activeGlobalStylesID, gs.ID)
			}

			// Context helper methods
			gsView, err := client.GlobalStyles().Retrieve(activeGlobalStylesID).ContextView().Do()
			require.NoError(t, err)
			assert.Equal(t, activeGlobalStylesID, gsView.ID)

			gsEdit, err := client.GlobalStyles().Retrieve(activeGlobalStylesID).ContextEdit().Do()
			require.NoError(t, err)
			assert.Equal(t, activeGlobalStylesID, gsEdit.ID)

			gsEmbed, err := client.GlobalStyles().Retrieve(activeGlobalStylesID).ContextEmbed().Do()
			require.NoError(t, err)
			assert.Equal(t, activeGlobalStylesID, gsEmbed.ID)
		})

		t.Run("Fields", func(t *testing.T) {
			gs, err := client.GlobalStyles().
				Retrieve(activeGlobalStylesID).
				Fields("id", "title").
				Do()
			require.NoError(t, err)
			require.NotNil(t, gs)
			assert.Equal(t, activeGlobalStylesID, gs.ID)
			require.NotNil(t, gs.Title)
			assert.Empty(t, gs.Styles)
		})

		t.Run("NotFound", func(t *testing.T) {
			gs, err := client.GlobalStyles().Retrieve(999999).Do()
			require.Error(t, err)
			assert.Nil(t, gs)

			var wpErr *gowprest.WPRestError
			if errors.As(err, &wpErr) {
				assert.Equal(t, "rest_global_styles_not_found", wpErr.Code)
				assert.Equal(t, 404, wpErr.Data.Status)
			} else {
				t.Fatalf("expected WPRestError, got: %v", err)
			}
		})
	})

	t.Run("UpdateGlobalStyles", func(t *testing.T) {
		updatedTitle := "Custom Styles Test Update"

		updated, err := client.GlobalStyles().
			Update(activeGlobalStylesID).
			Title(updatedTitle).
			Do()
		require.NoError(t, err)
		require.NotNil(t, updated)
		require.NotNil(t, updated.Title)
		assert.Equal(t, updatedTitle, updated.Title.Raw)

		// Revert back
		reverted, err := client.GlobalStyles().
			Update(activeGlobalStylesID).
			Title("Custom Styles").
			Do()
		require.NoError(t, err)
		require.NotNil(t, reverted)
		assert.Equal(t, "Custom Styles", reverted.Title.Raw)
	})

	t.Run("ThemeStylesAndVariations", func(t *testing.T) {
		t.Run("ThemeStyles", func(t *testing.T) {
			themeStyles, err := client.GlobalStyles().Themes("twentytwentyfive").Retrieve().Do()
			require.NoError(t, err)
			require.NotNil(t, themeStyles)
			assert.NotEmpty(t, themeStyles.Styles)
			assert.NotEmpty(t, themeStyles.Settings)
			assert.NotEmpty(t, themeStyles.Links)
		})

		t.Run("ThemeVariations", func(t *testing.T) {
			variations, err := client.GlobalStyles().Themes("twentytwentyfive").Variations().Do()
			require.NoError(t, err)
			require.NotEmpty(t, variations)

			for _, v := range variations {
				assert.Greater(t, v.Version, 0)
				require.NotNil(t, v.Title)
				assert.NotEmpty(t, v.Title.Rendered)
			}
		})
	})

	t.Run("UnauthorizedAccess", func(t *testing.T) {
		unauthClient := gowprest.NewClient(blogUrl)
		defer unauthClient.Close()

		_, err := unauthClient.GlobalStyles().Retrieve(activeGlobalStylesID).Do()
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
