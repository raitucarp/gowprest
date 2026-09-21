package tests

import (
	"os"
	"testing"

	_ "github.com/joho/godotenv/autoload"
	"github.com/raitucarp/gowprest"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestThemes(t *testing.T) {
	blogUrl := os.Getenv("BLOG_URL")
	username := os.Getenv("BLOG_USERNAME")
	password := os.Getenv("BLOG_APP_PASSWORD")

	if blogUrl == "" || username == "" || password == "" {
		t.Skip("Skipping test; BLOG_URL, BLOG_USERNAME, or BLOG_APP_PASSWORD not set")
	}

	client := gowprest.NewClient(blogUrl).WithBasicAuth(username, password)
	defer client.Close()

	themesAPI := client.Themes()

	var activeStylesheet string

	t.Run("ListThemes", func(t *testing.T) {
		t.Run("list_default", func(t *testing.T) {
			themes, err := themesAPI.List().Do()
			require.NoError(t, err)
			assert.NotEmpty(t, themes)

			first := themes[0]
			assert.NotEmpty(t, first.Stylesheet)
			assert.NotEmpty(t, first.Name.Rendered)
			assert.NotEmpty(t, first.Version)
		})

		t.Run("filter_by_status_active", func(t *testing.T) {
			activeThemes, err := themesAPI.List().StatusActive().Do()
			require.NoError(t, err)
			require.Len(t, activeThemes, 1)
			assert.Equal(t, "active", activeThemes[0].Status)
			assert.NotEmpty(t, activeThemes[0].Stylesheet)

			activeStylesheet = activeThemes[0].Stylesheet
		})

		t.Run("filter_by_status_inactive", func(t *testing.T) {
			inactiveThemes, err := themesAPI.List().StatusInactive().Do()
			require.NoError(t, err)
			for _, th := range inactiveThemes {
				assert.Equal(t, "inactive", th.Status)
			}
		})

		t.Run("global_parameter_fields", func(t *testing.T) {
			fieldsThemes, err := themesAPI.List().
				Fields("stylesheet", "version").
				Do()
			require.NoError(t, err)
			assert.NotEmpty(t, fieldsThemes)
			assert.NotEmpty(t, fieldsThemes[0].Stylesheet)
			assert.NotEmpty(t, fieldsThemes[0].Version)
			// Template should be empty when limited by _fields
			assert.Empty(t, fieldsThemes[0].Template)
		})

		t.Run("unauthorized_access", func(t *testing.T) {
			unauthClient := gowprest.NewClient(blogUrl)
			defer unauthClient.Close()

			themes, err := unauthClient.Themes().List().Do()
			require.Error(t, err)
			assert.Nil(t, themes)
			wpErr, ok := err.(*gowprest.WPRestError)
			require.True(t, ok, "expected error to be *gowprest.WPRestError")
			assert.Equal(t, 401, wpErr.Data.Status)
			assert.Equal(t, "rest_cannot_view_themes", wpErr.Code)
		})
	})

	t.Run("RetrieveTheme", func(t *testing.T) {
		require.NotEmpty(t, activeStylesheet, "activeStylesheet must be set from list tests")

		t.Run("retrieve_active_theme", func(t *testing.T) {
			theme, err := themesAPI.Retrieve(activeStylesheet).Do()
			require.NoError(t, err)
			require.NotNil(t, theme)
			assert.Equal(t, activeStylesheet, theme.Stylesheet)
			assert.Equal(t, "active", theme.Status)
			assert.True(t, theme.IsBlockTheme)
			assert.NotEmpty(t, theme.Screenshot)
			assert.NotEmpty(t, theme.RequiresPHP)
			assert.NotEmpty(t, theme.RequiresWP)
		})

		t.Run("active_convenience_method", func(t *testing.T) {
			theme, err := themesAPI.Active()
			require.NoError(t, err)
			require.NotNil(t, theme)
			assert.Equal(t, activeStylesheet, theme.Stylesheet)
			assert.Equal(t, "active", theme.Status)
		})

		t.Run("retrieve_with_fields", func(t *testing.T) {
			theme, err := themesAPI.Retrieve(activeStylesheet).
				Fields("stylesheet", "name").
				Do()
			require.NoError(t, err)
			require.NotNil(t, theme)
			assert.Equal(t, activeStylesheet, theme.Stylesheet)
			assert.NotEmpty(t, theme.Name.Rendered)
			assert.Empty(t, theme.Version)
		})

		t.Run("retrieve_not_found", func(t *testing.T) {
			theme, err := themesAPI.Retrieve("non_existent_theme_xyz").Do()
			require.Error(t, err)
			assert.Nil(t, theme)
			wpErr, ok := err.(*gowprest.WPRestError)
			require.True(t, ok, "expected error to be *gowprest.WPRestError")
			assert.Equal(t, 404, wpErr.Data.Status)
			assert.Equal(t, "rest_theme_not_found", wpErr.Code)
		})
	})
}
