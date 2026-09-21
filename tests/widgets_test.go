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

func TestWidgets(t *testing.T) {
	blogUrl := os.Getenv("BLOG_URL")
	username := os.Getenv("BLOG_USERNAME")
	password := os.Getenv("BLOG_APP_PASSWORD")

	if blogUrl == "" || username == "" || password == "" {
		t.Skip("Skipping test; BLOG_URL, BLOG_USERNAME, or BLOG_APP_PASSWORD not set")
	}

	client := gowprest.NewClient(blogUrl).WithBasicAuth(username, password)
	defer client.Close()

	t.Run("WidgetsAliases", func(t *testing.T) {
		assert.NotNil(t, client.Widgets())
		assert.NotNil(t, client.Widget())

		assert.NotNil(t, client.Widgets().List())
		assert.NotNil(t, client.Widgets().Retrieve("block-2"))
		assert.NotNil(t, client.Widgets().Create())
		assert.NotNil(t, client.Widgets().Update("block-2"))
		assert.NotNil(t, client.Widgets().Delete("block-2"))
	})

	t.Run("ListWidgets", func(t *testing.T) {
		widgets, err := client.Widgets().List().
			ContextEdit().
			Do()
		require.NoError(t, err)
		assert.NotEmpty(t, widgets)

		firstWidget := widgets[0]
		assert.NotEmpty(t, firstWidget.ID)
		assert.NotEmpty(t, firstWidget.IDBase)
		assert.NotEmpty(t, firstWidget.Sidebar)

		// ContextView & ContextEmbed
		_, err = client.Widgets().List().ContextView().Do()
		require.NoError(t, err)

		_, err = client.Widgets().List().ContextEmbed().Do()
		require.NoError(t, err)

		// Filter by Sidebar
		inactiveWidgets, err := client.Widgets().List().
			Sidebar("wp_inactive_widgets").
			Do()
		require.NoError(t, err)
		assert.NotEmpty(t, inactiveWidgets)
		for _, w := range inactiveWidgets {
			assert.Equal(t, "wp_inactive_widgets", w.Sidebar)
		}

		// Sparse fields
		sparseList, err := client.Widgets().List().
			Fields("id", "id_base", "sidebar").
			Do()
		require.NoError(t, err)
		assert.NotEmpty(t, sparseList)
		assert.NotEmpty(t, sparseList[0].ID)
		assert.NotEmpty(t, sparseList[0].IDBase)
		assert.NotEmpty(t, sparseList[0].Sidebar)
		assert.Nil(t, sparseList[0].Instance)
	})

	t.Run("WidgetCRUDWorkflow", func(t *testing.T) {
		initialContent := "<!-- wp:paragraph --><p>Custom Test Block Widget</p><!-- /wp:paragraph -->"
		updatedContent := "<!-- wp:paragraph --><p>Updated Test Block Widget</p><!-- /wp:paragraph -->"

		// 1. Create a block widget
		created, err := client.Widgets().Create().
			IDBase("block").
			Sidebar("wp_inactive_widgets").
			Instance(map[string]any{
				"raw": map[string]any{
					"content": initialContent,
				},
			}).
			Do()
		require.NoError(t, err)
		require.NotNil(t, created)
		assert.NotEmpty(t, created.ID)
		assert.Equal(t, "block", created.IDBase)
		assert.Equal(t, "wp_inactive_widgets", created.Sidebar)

		widgetID := created.ID
		defer func() {
			_, _ = client.Widgets().Delete(widgetID).Force().Do()
		}()

		// 2. Retrieve the widget
		retrieved, err := client.Widgets().Retrieve(widgetID).
			ContextEdit().
			Do()
		require.NoError(t, err)
		require.NotNil(t, retrieved)
		assert.Equal(t, widgetID, retrieved.ID)
		assert.Equal(t, "block", retrieved.IDBase)
		require.NotNil(t, retrieved.Instance)
		assert.Equal(t, initialContent, retrieved.Instance.Raw["content"])

		// Sparse retrieve
		sparseWidget, err := client.Widgets().Retrieve(widgetID).
			Fields("id", "id_base", "sidebar").
			Do()
		require.NoError(t, err)
		require.NotNil(t, sparseWidget)
		assert.Equal(t, widgetID, sparseWidget.ID)
		assert.Equal(t, "block", sparseWidget.IDBase)
		assert.Equal(t, "wp_inactive_widgets", sparseWidget.Sidebar)
		assert.Nil(t, sparseWidget.Instance)

		// 3. Update the widget
		updated, err := client.Widgets().Update(widgetID).
			Instance(map[string]any{
				"raw": map[string]any{
					"content": updatedContent,
				},
			}).
			Do()
		require.NoError(t, err)
		require.NotNil(t, updated)
		assert.Equal(t, widgetID, updated.ID)
		require.NotNil(t, updated.Instance)
		assert.Equal(t, updatedContent, updated.Instance.Raw["content"])

		// 4. Delete the widget with Force
		deleted, err := client.Widgets().Delete(widgetID).
			Force().
			Do()
		require.NoError(t, err)
		require.NotNil(t, deleted)
		assert.Equal(t, widgetID, deleted.ID)

		// Verify deletion
		_, err = client.Widgets().Retrieve(widgetID).Do()
		require.Error(t, err)
		var wpErr *gowprest.WPRestError
		require.True(t, errors.As(err, &wpErr))
		assert.Equal(t, 404, wpErr.Data.Status)
		assert.Equal(t, "rest_widget_not_found", wpErr.Code)
	})

	t.Run("WidgetsErrorHandling", func(t *testing.T) {
		// Nonexistent widget
		_, err := client.Widgets().Retrieve("nonexistent_widget_id").Do()
		require.Error(t, err)
		var wpErr *gowprest.WPRestError
		require.True(t, errors.As(err, &wpErr))
		assert.Equal(t, 404, wpErr.Data.Status)
		assert.Equal(t, "rest_widget_not_found", wpErr.Code)

		// Unauthorized request
		unauthClient := gowprest.NewClient(blogUrl)
		defer unauthClient.Close()

		_, err = unauthClient.Widgets().List().Do()
		require.Error(t, err)
		var unauthErr *gowprest.WPRestError
		require.True(t, errors.As(err, &unauthErr))
		assert.Equal(t, 401, unauthErr.Data.Status)
		assert.Equal(t, "rest_cannot_manage_widgets", unauthErr.Code)
	})
}
