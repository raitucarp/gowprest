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

func TestSiteHealth(t *testing.T) {
	blogUrl := os.Getenv("BLOG_URL")
	username := os.Getenv("BLOG_USERNAME")
	password := os.Getenv("BLOG_APP_PASSWORD")

	if blogUrl == "" || username == "" || password == "" {
		t.Skip("Skipping test; BLOG_URL, BLOG_USERNAME, or BLOG_APP_PASSWORD not set")
	}

	client := gowprest.NewClient(blogUrl).WithBasicAuth(username, password)
	defer client.Close()

	t.Run("SiteHealthAliases", func(t *testing.T) {
		assert.NotNil(t, client.SiteHealth())
		assert.NotNil(t, client.SiteHealthTests())

		assert.NotNil(t, client.SiteHealth().Test("background-updates"))
		assert.NotNil(t, client.SiteHealth().BackgroundUpdates())
		assert.NotNil(t, client.SiteHealth().LoopbackRequests())
		assert.NotNil(t, client.SiteHealth().HttpsStatus())
		assert.NotNil(t, client.SiteHealth().DotorgCommunication())
		assert.NotNil(t, client.SiteHealth().AuthorizationHeader())
		assert.NotNil(t, client.SiteHealth().PageCache())
		assert.NotNil(t, client.SiteHealth().DirectorySizes())
	})

	t.Run("SiteHealthOfficialTests", func(t *testing.T) {
		testsToRun := []struct {
			name     string
			run      func() (*gowprest.SiteHealthTest, error)
			expected string
		}{
			{
				name: "background-updates",
				run: func() (*gowprest.SiteHealthTest, error) {
					return client.SiteHealth().BackgroundUpdates().Execute()
				},
				expected: "background_updates",
			},
			{
				name: "loopback-requests",
				run: func() (*gowprest.SiteHealthTest, error) {
					return client.SiteHealth().LoopbackRequests().Get()
				},
				expected: "loopback_requests",
			},
			{
				name: "https-status",
				run: func() (*gowprest.SiteHealthTest, error) {
					return client.SiteHealth().HttpsStatus().Execute()
				},
				expected: "https_status",
			},
			{
				name: "dotorg-communication",
				run: func() (*gowprest.SiteHealthTest, error) {
					return client.SiteHealth().DotorgCommunication().Execute()
				},
				expected: "dotorg_communication",
			},
			{
				name: "authorization-header",
				run: func() (*gowprest.SiteHealthTest, error) {
					return client.SiteHealth().AuthorizationHeader().Execute()
				},
				expected: "authorization_header",
			},
		}

		for _, tc := range testsToRun {
			t.Run(tc.name, func(t *testing.T) {
				result, err := tc.run()
				require.NoError(t, err)
				require.NotNil(t, result)

				assert.Equal(t, tc.expected, result.Test)
				assert.NotEmpty(t, result.Label)
				assert.Contains(t, []string{
					gowprest.SiteHealthStatusGood,
					gowprest.SiteHealthStatusRecommended,
					gowprest.SiteHealthStatusCritical,
				}, result.Status)
				assert.NotEmpty(t, result.Badge.Label)
			})
		}
	})

	t.Run("SiteHealthPageCacheAndDirectorySizes", func(t *testing.T) {
		// Page Cache test
		pcResult, err := client.SiteHealth().PageCache().Execute()
		require.NoError(t, err)
		require.NotNil(t, pcResult)
		assert.Equal(t, "page_cache", pcResult.Test)
		assert.NotEmpty(t, pcResult.Label)

		// Directory Sizes endpoint
		dirSizes, err := client.SiteHealth().DirectorySizes().Execute()
		require.NoError(t, err)
		require.NotNil(t, dirSizes)
		assert.True(t, dirSizes.WordPressSize.Size != "" || dirSizes.DatabaseSize.Size != "")
	})

	t.Run("SiteHealthFieldsProjection", func(t *testing.T) {
		res, err := client.SiteHealth().
			BackgroundUpdates().
			Fields("test", "status").
			Execute()
		require.NoError(t, err)
		require.NotNil(t, res)

		assert.Equal(t, "background_updates", res.Test)
		assert.NotEmpty(t, res.Status)
		assert.Empty(t, res.Description)
		assert.Empty(t, res.Actions)
	})

	t.Run("SiteHealthErrorHandling", func(t *testing.T) {
		// Nonexistent test
		_, err := client.SiteHealth().Test("nonexistent-test-identifier").Execute()
		require.Error(t, err)

		var restErr *gowprest.WPRestError
		require.True(t, errors.As(err, &restErr))
		assert.Equal(t, "rest_no_route", restErr.Code)
		assert.Equal(t, 404, restErr.Data.Status)

		// Unauthorized request
		unauthClient := gowprest.NewClient(blogUrl)
		defer unauthClient.Close()
		_, err = unauthClient.SiteHealth().BackgroundUpdates().Execute()
		require.Error(t, err)

		var unauthErr *gowprest.WPRestError
		require.True(t, errors.As(err, &unauthErr))
		assert.Equal(t, "rest_forbidden", unauthErr.Code)
		assert.Equal(t, 401, unauthErr.Data.Status)
	})
}
