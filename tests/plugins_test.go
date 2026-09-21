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

func TestPlugins(t *testing.T) {
	blogUrl := os.Getenv("BLOG_URL")
	username := os.Getenv("BLOG_USERNAME")
	password := os.Getenv("BLOG_APP_PASSWORD")

	if blogUrl == "" || username == "" || password == "" {
		t.Skip("Skipping test; BLOG_URL, BLOG_USERNAME, or BLOG_APP_PASSWORD not set")
	}

	client := gowprest.NewClient(blogUrl).WithBasicAuth(username, password)
	defer client.Close()

	t.Run("PluginsAliases", func(t *testing.T) {
		assert.NotNil(t, client.Plugins())
		assert.NotNil(t, client.Plugins().List())
	})

	t.Run("ListPlugins", func(t *testing.T) {
		t.Run("Default", func(t *testing.T) {
			plugins, err := client.Plugins().List().Do()
			require.NoError(t, err)
			require.NotEmpty(t, plugins)

			var foundHello bool
			for _, p := range plugins {
				assert.NotEmpty(t, p.Plugin)
				assert.NotEmpty(t, p.Name)
				assert.NotEmpty(t, p.Status)
				if p.Plugin == "hello" {
					foundHello = true
				}
			}
			assert.True(t, foundHello, "expected 'hello' plugin to be installed by default")
		})

		t.Run("StatusFilter", func(t *testing.T) {
			inactivePlugins, err := client.Plugins().List().StatusInactive().Do()
			require.NoError(t, err)
			for _, p := range inactivePlugins {
				assert.Equal(t, gowprest.PluginStatusInactive, p.Status)
				assert.True(t, p.IsInactive())
				assert.False(t, p.IsActive())
			}
		})

		t.Run("Search", func(t *testing.T) {
			results, err := client.Plugins().List().Search("hello").Do()
			require.NoError(t, err)
			require.NotEmpty(t, results)
			assert.Equal(t, "hello", results[0].Plugin)
		})

		t.Run("Contexts", func(t *testing.T) {
			viewPlugins, err := client.Plugins().List().ContextView().Do()
			require.NoError(t, err)
			assert.NotEmpty(t, viewPlugins)

			editPlugins, err := client.Plugins().List().ContextEdit().Do()
			require.NoError(t, err)
			assert.NotEmpty(t, editPlugins)

			embedPlugins, err := client.Plugins().List().ContextEmbed().Do()
			require.NoError(t, err)
			assert.NotEmpty(t, embedPlugins)
		})

		t.Run("Fields", func(t *testing.T) {
			fieldsPlugins, err := client.Plugins().List().
				Fields("plugin", "status", "name").
				Do()
			require.NoError(t, err)
			require.NotEmpty(t, fieldsPlugins)

			for _, p := range fieldsPlugins {
				assert.NotEmpty(t, p.Plugin)
				assert.NotEmpty(t, p.Status)
				assert.NotEmpty(t, p.Name)
				assert.Empty(t, p.Version)
				assert.Empty(t, p.PluginURI)
			}
		})
	})

	t.Run("RetrievePlugin", func(t *testing.T) {
		t.Run("Retrieve_SingleFile", func(t *testing.T) {
			p, err := client.Plugins().Retrieve("hello").Do()
			require.NoError(t, err)
			require.NotNil(t, p)
			assert.Equal(t, "hello", p.Plugin)
			assert.Equal(t, "Hello Dolly", p.Name)
		})

		t.Run("Retrieve_SubdirectoryCombined", func(t *testing.T) {
			p, err := client.Plugins().Retrieve("akismet/akismet").Do()
			require.NoError(t, err)
			require.NotNil(t, p)
			assert.Equal(t, "akismet/akismet", p.Plugin)
			assert.Contains(t, p.Name, "Akismet")
		})

		t.Run("Retrieve_SubdirectorySeparate", func(t *testing.T) {
			p, err := client.Plugins().Retrieve("akismet", "akismet").Do()
			require.NoError(t, err)
			require.NotNil(t, p)
			assert.Equal(t, "akismet/akismet", p.Plugin)
		})

		t.Run("Retrieve_ContextAndFields", func(t *testing.T) {
			p, err := client.Plugins().Retrieve("hello").
				ContextEdit().
				Fields("plugin", "name").
				Do()
			require.NoError(t, err)
			require.NotNil(t, p)
			assert.Equal(t, "hello", p.Plugin)
			assert.Equal(t, "Hello Dolly", p.Name)
			assert.Empty(t, p.Version)
		})

		t.Run("Retrieve_NotFound", func(t *testing.T) {
			_, err := client.Plugins().Retrieve("nonexistent-plugin-xyz").Do()
			require.Error(t, err)

			var wpErr *gowprest.WPRestError
			if assert.True(t, errors.As(err, &wpErr)) {
				assert.Equal(t, "rest_plugin_not_found", wpErr.Code)
				assert.Equal(t, 404, wpErr.Data.Status)
			}
		})
	})

	t.Run("UpdatePlugin_ActivationAndDeactivation", func(t *testing.T) {
		// Ensure initial state is inactive
		_, _ = client.Plugins().Update("hello").Deactivate().Do()

		// Activate
		activated, err := client.Plugins().Update("hello").Activate().Do()
		require.NoError(t, err)
		require.NotNil(t, activated)
		assert.Equal(t, gowprest.PluginStatusActive, activated.Status)
		assert.True(t, activated.IsActive())

		// Verify via retrieve
		retrieved, err := client.Plugins().Retrieve("hello").Do()
		require.NoError(t, err)
		assert.Equal(t, gowprest.PluginStatusActive, retrieved.Status)

		// Deactivate
		deactivated, err := client.Plugins().Update("hello").Deactivate().Do()
		require.NoError(t, err)
		require.NotNil(t, deactivated)
		assert.Equal(t, gowprest.PluginStatusInactive, deactivated.Status)
		assert.True(t, deactivated.IsInactive())

		// Verify via retrieve
		retrieved2, err := client.Plugins().Retrieve("hello").Do()
		require.NoError(t, err)
		assert.Equal(t, gowprest.PluginStatusInactive, retrieved2.Status)
	})

	t.Run("CreatePlugin_InvalidSlug", func(t *testing.T) {
		_, err := client.Plugins().Create().
			Slug("nonexistent-plugin-slug-xyz-98765").
			Do()
		require.Error(t, err)

		var wpErr *gowprest.WPRestError
		assert.True(t, errors.As(err, &wpErr))
	})

	t.Run("PluginLifecycle_InstallUpdateDelete", func(t *testing.T) {
		pluginSlug := "classic-widgets"

		// Pre-cleanup in case it was left over from a previous run
		_, _ = client.Plugins().Delete(pluginSlug + "/classic-widgets").Do()
		_, _ = client.Plugins().Delete(pluginSlug).Do()

		// 1. Install (Create) plugin from WordPress.org
		installed, err := client.Plugins().Create().
			Slug(pluginSlug).
			StatusInactive().
			Do()
		if err != nil {
			t.Logf("Plugin install failed: %v", err)
		}
		require.NoError(t, err)
		require.NotNil(t, installed)
		t.Logf("Installed plugin: %s, status: %s", installed.Plugin, installed.Status)
		assert.Contains(t, installed.Plugin, pluginSlug)

		pluginFile := installed.Plugin

		// 2. Activate (Update)
		activated, err := client.Plugins().Update(pluginFile).Activate().Do()
		require.NoError(t, err)
		assert.True(t, activated.IsActive())

		// 3. Deactivate (Update)
		deactivated, err := client.Plugins().Update(pluginFile).Deactivate().Do()
		require.NoError(t, err)
		assert.True(t, deactivated.IsInactive())

		// 4. Delete
		deleted, err := client.Plugins().Delete(pluginFile).Do()
		require.NoError(t, err)
		require.NotNil(t, deleted)

		// 5. Verify 404 after deletion
		_, err = client.Plugins().Retrieve(pluginFile).Do()
		require.Error(t, err)
		var notFoundErr *gowprest.WPRestError
		require.True(t, errors.As(err, &notFoundErr))
		assert.Equal(t, 404, notFoundErr.Data.Status)
	})
}

