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

func TestWidgetTypes(t *testing.T) {
	blogUrl := os.Getenv("BLOG_URL")
	username := os.Getenv("BLOG_USERNAME")
	password := os.Getenv("BLOG_APP_PASSWORD")

	if blogUrl == "" || username == "" || password == "" {
		t.Skip("Skipping test; BLOG_URL, BLOG_USERNAME, or BLOG_APP_PASSWORD not set")
	}

	client := gowprest.NewClient(blogUrl).WithBasicAuth(username, password)
	defer client.Close()

	t.Run("WidgetTypesAliases", func(t *testing.T) {
		assert.NotNil(t, client.WidgetTypes())
		assert.NotNil(t, client.WidgetType())

		assert.NotNil(t, client.WidgetTypes().List())
		assert.NotNil(t, client.WidgetTypes().Retrieve("archives"))
	})

	t.Run("ListWidgetTypes", func(t *testing.T) {
		widgetTypes, err := client.WidgetTypes().List().
			ContextEdit().
			Do()
		require.NoError(t, err)
		assert.NotEmpty(t, widgetTypes)

		firstWT := widgetTypes[0]
		assert.NotEmpty(t, firstWT.ID)
		assert.NotEmpty(t, firstWT.Name)

		// ContextView & ContextEmbed
		_, err = client.WidgetTypes().List().ContextView().Do()
		require.NoError(t, err)

		_, err = client.WidgetTypes().List().ContextEmbed().Do()
		require.NoError(t, err)

		// Sparse fields
		sparseList, err := client.WidgetTypes().List().
			Fields("id", "name").
			Do()
		require.NoError(t, err)
		assert.NotEmpty(t, sparseList)
		assert.NotEmpty(t, sparseList[0].ID)
		assert.NotEmpty(t, sparseList[0].Name)
		assert.Empty(t, sparseList[0].Description)
	})

	t.Run("RetrieveWidgetType", func(t *testing.T) {
		wt, err := client.WidgetTypes().Retrieve("archives").
			ContextEdit().
			Do()
		require.NoError(t, err)
		require.NotNil(t, wt)
		assert.Equal(t, "archives", wt.ID)
		assert.Equal(t, "Archives", wt.Name)

		// Sparse retrieve
		sparseWT, err := client.WidgetTypes().Retrieve("archives").
			Fields("id", "name").
			Do()
		require.NoError(t, err)
		require.NotNil(t, sparseWT)
		assert.Equal(t, "archives", sparseWT.ID)
		assert.Equal(t, "Archives", sparseWT.Name)
		assert.Empty(t, sparseWT.Description)
	})

	t.Run("WidgetTypesErrorHandling", func(t *testing.T) {
		// Nonexistent widget type
		_, err := client.WidgetTypes().Retrieve("nonexistent_widget_type").Do()
		require.Error(t, err)
		var wpErr *gowprest.WPRestError
		require.True(t, errors.As(err, &wpErr))
		assert.Equal(t, 404, wpErr.Data.Status)
		assert.Equal(t, "rest_widget_type_invalid", wpErr.Code)

		// Unauthorized request
		unauthClient := gowprest.NewClient(blogUrl)
		defer unauthClient.Close()

		_, err = unauthClient.WidgetTypes().List().Do()
		require.Error(t, err)
		var unauthErr *gowprest.WPRestError
		require.True(t, errors.As(err, &unauthErr))
		assert.Equal(t, 401, unauthErr.Data.Status)
		assert.Equal(t, "rest_cannot_manage_widgets", unauthErr.Code)
	})
}
