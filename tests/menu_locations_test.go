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

func TestMenuLocations(t *testing.T) {
	blogUrl := os.Getenv("BLOG_URL")
	username := os.Getenv("BLOG_USERNAME")
	password := os.Getenv("BLOG_APP_PASSWORD")

	if blogUrl == "" || username == "" || password == "" {
		t.Skip("Skipping test; BLOG_URL, BLOG_USERNAME, or BLOG_APP_PASSWORD not set")
	}

	client := gowprest.NewClient(blogUrl).WithBasicAuth(username, password)
	defer client.Close()

	t.Run("ServiceAliases", func(t *testing.T) {
		assert.NotNil(t, client.MenuLocations())
		assert.NotNil(t, client.NavMenuLocations())
		assert.NotNil(t, client.MenuLocations().List())
		assert.NotNil(t, client.MenuLocations().Retrieve("primary"))
	})

	t.Run("ListMenuLocations", func(t *testing.T) {
		t.Run("Default", func(t *testing.T) {
			locations, err := client.MenuLocations().List().Do()
			require.NoError(t, err)
			require.NotEmpty(t, locations)

			primary, found := locations["primary"]
			require.True(t, found, "expected 'primary' location to exist")
			assert.Equal(t, "primary", primary.Name)
			assert.Equal(t, "Primary Navigation", primary.Description)

			footer, found := locations["footer"]
			require.True(t, found, "expected 'footer' location to exist")
			assert.Equal(t, "footer", footer.Name)
			assert.Equal(t, "Footer Navigation", footer.Description)
		})

		t.Run("Contexts", func(t *testing.T) {
			for _, ctx := range []string{"view", "edit", "embed"} {
				locations, err := client.MenuLocations().
					List().
					Context(ctx).
					Do()
				require.NoError(t, err)
				assert.NotEmpty(t, locations)
			}

			// Fluent context helpers
			locView, err := client.MenuLocations().List().ContextView().Do()
			require.NoError(t, err)
			assert.NotEmpty(t, locView)

			locEdit, err := client.MenuLocations().List().ContextEdit().Do()
			require.NoError(t, err)
			assert.NotEmpty(t, locEdit)

			locEmbed, err := client.MenuLocations().List().ContextEmbed().Do()
			require.NoError(t, err)
			assert.NotEmpty(t, locEmbed)
		})
	})

	t.Run("RetrieveMenuLocation", func(t *testing.T) {
		t.Run("ExistingLocation", func(t *testing.T) {
			loc, err := client.MenuLocations().Retrieve("primary").Do()
			require.NoError(t, err)
			require.NotNil(t, loc)
			assert.Equal(t, "primary", loc.Name)
			assert.Equal(t, "Primary Navigation", loc.Description)
			assert.NotEmpty(t, loc.Links)

			// Context helpers on Retrieve
			locView, err := client.MenuLocations().Retrieve("primary").ContextView().Do()
			require.NoError(t, err)
			assert.Equal(t, "primary", locView.Name)

			locEdit, err := client.MenuLocations().Retrieve("primary").ContextEdit().Do()
			require.NoError(t, err)
			assert.Equal(t, "primary", locEdit.Name)

			locEmbed, err := client.MenuLocations().Retrieve("primary").ContextEmbed().Do()
			require.NoError(t, err)
			assert.Equal(t, "primary", locEmbed.Name)
		})

		t.Run("WithFields", func(t *testing.T) {
			loc, err := client.MenuLocations().
				Retrieve("primary").
				Fields("name").
				Do()
			require.NoError(t, err)
			require.NotNil(t, loc)
			assert.Equal(t, "primary", loc.Name)
			assert.Empty(t, loc.Description)
			assert.Empty(t, loc.Links)
		})

		t.Run("NotFound", func(t *testing.T) {
			loc, err := client.MenuLocations().Retrieve("invalid-location-xyz").Do()
			require.Error(t, err)
			assert.Nil(t, loc)

			var wpErr *gowprest.WPRestError
			if errors.As(err, &wpErr) {
				assert.Equal(t, "rest_menu_location_invalid", wpErr.Code)
				assert.Equal(t, 404, wpErr.Data.Status)
			} else {
				t.Fatalf("expected WPRestError, got: %v", err)
			}
		})
	})

	t.Run("UnauthorizedAccess", func(t *testing.T) {
		unauthClient := gowprest.NewClient(blogUrl)
		defer unauthClient.Close()

		_, err := unauthClient.MenuLocations().List().Do()
		require.Error(t, err)

		var wpErr *gowprest.WPRestError
		if errors.As(err, &wpErr) {
			assert.Equal(t, "rest_cannot_view", wpErr.Code)
			assert.Equal(t, 401, wpErr.Data.Status)
		} else {
			t.Fatalf("expected WPRestError, got: %v", err)
		}

		_, err = unauthClient.MenuLocations().Retrieve("primary").Do()
		require.Error(t, err)
		if errors.As(err, &wpErr) {
			assert.Equal(t, "rest_cannot_view", wpErr.Code)
			assert.Equal(t, 401, wpErr.Data.Status)
		} else {
			t.Fatalf("expected WPRestError, got: %v", err)
		}
	})
}
