package tests

import (
	"os"
	"testing"

	_ "github.com/joho/godotenv/autoload"
	"github.com/raitucarp/gowprest"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSettings(t *testing.T) {
	blogUrl := os.Getenv("BLOG_URL")
	username := os.Getenv("BLOG_USERNAME")
	password := os.Getenv("BLOG_APP_PASSWORD")

	if blogUrl == "" || username == "" || password == "" {
		t.Skip("Skipping test; BLOG_URL, BLOG_USERNAME, or BLOG_APP_PASSWORD not set")
	}

	client := gowprest.NewClient(blogUrl).WithBasicAuth(username, password)
	defer client.Close()

	settingsAPI := client.Settings()

	t.Run("RetrieveSettings", func(t *testing.T) {
		t.Run("retrieve_default", func(t *testing.T) {
			s, err := settingsAPI.Retrieve().Do()
			require.NoError(t, err)
			require.NotNil(t, s)
			assert.NotEmpty(t, s.Title)
			assert.NotEmpty(t, s.URL)
			assert.NotEmpty(t, s.Email)
			assert.NotZero(t, s.PostsPerPage)

			// Also verify Get() convenience method
			sGet, err := settingsAPI.Get()
			require.NoError(t, err)
			require.NotNil(t, sGet)
			assert.Equal(t, s.Title, sGet.Title)
		})

		t.Run("retrieve_with_fields", func(t *testing.T) {
			s, err := settingsAPI.Retrieve().
				Fields("title", "description").
				Do()
			require.NoError(t, err)
			require.NotNil(t, s)
			assert.NotEmpty(t, s.Title)
			// Email should be omitted when restricted by _fields
			assert.Empty(t, s.Email)
		})

		t.Run("retrieve_unauthorized", func(t *testing.T) {
			unauthClient := gowprest.NewClient(blogUrl)
			defer unauthClient.Close()

			s, err := unauthClient.Settings().Retrieve().Do()
			require.Error(t, err)
			assert.Nil(t, s)
			wpErr, ok := err.(*gowprest.WPRestError)
			require.True(t, ok, "expected error to be *gowprest.WPRestError")
			assert.Equal(t, 401, wpErr.Data.Status)
			assert.Equal(t, "rest_forbidden", wpErr.Code)
		})
	})

	t.Run("UpdateSettings", func(t *testing.T) {
		// Fetch initial settings to restore later
		original, err := settingsAPI.Get()
		require.NoError(t, err)
		require.NotNil(t, original)

		defer func() {
			// Restore original description
			_, _ = settingsAPI.Update().Description(original.Description).Do()
		}()

		t.Run("update_chained", func(t *testing.T) {
			testTagline := "Test Tagline via Chaining"
			updated, err := settingsAPI.Update().
				Description(testTagline).
				Do()
			require.NoError(t, err)
			require.NotNil(t, updated)
			assert.Equal(t, testTagline, updated.Description)

			// Verify with a separate Get()
			fetched, err := settingsAPI.Get()
			require.NoError(t, err)
			assert.Equal(t, testTagline, fetched.Description)
		})

		t.Run("update_with_struct", func(t *testing.T) {
			newTagline := "Tagline updated via SettingsData"
			data := gowprest.SettingsData{
				Description: &newTagline,
			}
			updated, err := settingsAPI.Update(data).Do()
			require.NoError(t, err)
			require.NotNil(t, updated)
			assert.Equal(t, newTagline, updated.Description)
		})

		t.Run("update_invalid_param", func(t *testing.T) {
			_, err := settingsAPI.Update().
				Email("not-a-valid-email").
				Do()
			require.Error(t, err)
			wpErr, ok := err.(*gowprest.WPRestError)
			require.True(t, ok, "expected error to be *gowprest.WPRestError")
			assert.Equal(t, 400, wpErr.Data.Status)
			assert.Equal(t, "rest_invalid_param", wpErr.Code)
		})
	})
}
