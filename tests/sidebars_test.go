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

func TestSidebars(t *testing.T) {
	blogUrl := os.Getenv("BLOG_URL")
	username := os.Getenv("BLOG_USERNAME")
	password := os.Getenv("BLOG_APP_PASSWORD")

	if blogUrl == "" || username == "" || password == "" {
		t.Skip("Skipping test; BLOG_URL, BLOG_USERNAME, or BLOG_APP_PASSWORD not set")
	}

	client := gowprest.NewClient(blogUrl).WithBasicAuth(username, password)
	defer client.Close()

	t.Run("SidebarsAliases", func(t *testing.T) {
		assert.NotNil(t, client.Sidebars())
		assert.NotNil(t, client.Sidebar())

		assert.NotNil(t, client.Sidebars().List())
		assert.NotNil(t, client.Sidebars().Retrieve("wp_inactive_widgets"))
		assert.NotNil(t, client.Sidebars().Update("wp_inactive_widgets"))
	})

	t.Run("SidebarsWorkflow", func(t *testing.T) {
		// 1. List Sidebars
		sidebars, err := client.Sidebars().List().
			ContextEdit().
			Do()
		require.NoError(t, err)
		assert.NotEmpty(t, sidebars)

		targetSidebarID := sidebars[0].ID
		assert.NotEmpty(t, targetSidebarID)

		// List with Fields
		fieldsList, err := client.Sidebars().List().
			Fields("id", "name", "status").
			Do()
		require.NoError(t, err)
		assert.NotEmpty(t, fieldsList)
		assert.NotEmpty(t, fieldsList[0].ID)
		assert.NotEmpty(t, fieldsList[0].Name)
		assert.NotEmpty(t, fieldsList[0].Status)
		assert.Empty(t, fieldsList[0].BeforeWidget)

		// 2. Retrieve Sidebar
		sidebar, err := client.Sidebars().Retrieve(targetSidebarID).
			ContextEdit().
			Do()
		require.NoError(t, err)
		require.NotNil(t, sidebar)
		assert.Equal(t, targetSidebarID, sidebar.ID)
		assert.NotEmpty(t, sidebar.Name)
		assert.NotEmpty(t, sidebar.Status)

		// ContextView & ContextEmbed
		_, err = client.Sidebars().Retrieve(targetSidebarID).ContextView().Do()
		require.NoError(t, err)

		_, err = client.Sidebars().Retrieve(targetSidebarID).ContextEmbed().Do()
		require.NoError(t, err)

		// Retrieve with Fields
		fieldSidebar, err := client.Sidebars().Retrieve(targetSidebarID).
			Fields("id", "status").
			Do()
		require.NoError(t, err)
		require.NotNil(t, fieldSidebar)
		assert.Equal(t, targetSidebarID, fieldSidebar.ID)
		assert.NotEmpty(t, fieldSidebar.Status)
		assert.Empty(t, fieldSidebar.Name)

		// 3. Update Sidebar widgets
		updated, err := client.Sidebars().Update(targetSidebarID).
			Widgets(sidebar.Widgets...).
			Do()
		require.NoError(t, err)
		require.NotNil(t, updated)
		assert.Equal(t, targetSidebarID, updated.ID)
		assert.Equal(t, len(sidebar.Widgets), len(updated.Widgets))

		// 4. Verify 404 for nonexistent sidebar
		_, err = client.Sidebars().Retrieve("nonexistent_sidebar_id").Do()
		require.Error(t, err)
		var wpErr *gowprest.WPRestError
		require.True(t, errors.As(err, &wpErr))
		assert.Equal(t, 404, wpErr.Data.Status)
		assert.Equal(t, "rest_sidebar_not_found", wpErr.Code)
	})

	t.Run("SidebarsUnauthorized", func(t *testing.T) {
		unauthClient := gowprest.NewClient(blogUrl)
		defer unauthClient.Close()

		_, err := unauthClient.Sidebars().List().Do()
		require.Error(t, err)
		var wpErr *gowprest.WPRestError
		require.True(t, errors.As(err, &wpErr))
		assert.Equal(t, 401, wpErr.Data.Status)
		assert.Equal(t, "rest_cannot_manage_widgets", wpErr.Code)
	})
}
