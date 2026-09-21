package tests

import (
	"errors"
	"fmt"
	"os"
	"testing"
	"time"

	_ "github.com/joho/godotenv/autoload"
	"github.com/raitucarp/gowprest"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNavMenuItemRevisions(t *testing.T) {
	blogUrl := os.Getenv("BLOG_URL")
	username := os.Getenv("BLOG_USERNAME")
	password := os.Getenv("BLOG_APP_PASSWORD")

	if blogUrl == "" || username == "" || password == "" {
		t.Skip("Skipping test; BLOG_URL, BLOG_USERNAME, or BLOG_APP_PASSWORD not set")
	}

	client := gowprest.NewClient(blogUrl).WithBasicAuth(username, password)
	defer client.Close()

	t.Run("NavMenuItemRevisionsAliases", func(t *testing.T) {
		assert.NotNil(t, client.NavMenuItemRevisions(1))
		assert.NotNil(t, client.NavMenuItemAutosaves(1))
		assert.NotNil(t, client.MenuItemRevisions(1))
		assert.NotNil(t, client.MenuItemAutosaves(1))

		assert.NotNil(t, client.NavMenuItems().Revisions(1))
		assert.NotNil(t, client.NavMenuItems().Autosaves(1))
		assert.NotNil(t, client.MenuItems().Revisions(1))
		assert.NotNil(t, client.MenuItems().Autosaves(1))

		assert.NotNil(t, client.NavMenuItem().Revisions(1))
		assert.NotNil(t, client.MenuItem().Revisions(1))

		api := client.NavMenuItemRevisions(1)
		assert.NotNil(t, api.List())
		assert.NotNil(t, api.Retrieve(1))
		assert.NotNil(t, api.Create())
	})

	t.Run("NavMenuItemRevisionsCRUDWorkflow", func(t *testing.T) {
		unique := fmt.Sprintf("ItemRev%d", time.Now().UnixNano()%1000000)

		// 1. Setup parent menu
		menu, err := client.NavMenus().Create().
			Name("Menu " + unique).
			Do()
		require.NoError(t, err)
		require.NotNil(t, menu)
		defer func() {
			_, _ = client.NavMenus().Delete(menu.ID).Force().Do()
		}()

		// 2. Setup parent menu item
		parentItem, err := client.NavMenuItems().Create().
			Title("Parent Link " + unique).
			URL("https://example.com/" + unique).
			Menus(menu.ID).
			Status("publish").
			Do()
		require.NoError(t, err)
		require.NotNil(t, parentItem)
		parentItemID := parentItem.ID
		defer func() {
			_, _ = client.NavMenuItems().Delete(parentItemID).Force().Do()
		}()

		revisionsAPI := client.NavMenuItemRevisions(parentItemID)

		// 3. Create Autosave / Revision
		autosaveTitle := "Autosaved Link " + unique
		autosaveURL := "https://example.com/autosaved-" + unique
		created, err := revisionsAPI.Create().
			Title(autosaveTitle).
			URL(autosaveURL).
			Menus(menu.ID).
			Status("draft").
			Do()
		require.NoError(t, err)
		require.NotNil(t, created)
		assert.NotEmpty(t, created.ID)
		assert.Equal(t, parentItemID, created.Parent)

		autosaveID := created.ID

		// 4. List Autosaves
		list, err := revisionsAPI.List().
			ContextEdit().
			Do()
		require.NoError(t, err)
		assert.NotEmpty(t, list)
		found := false
		for _, as := range list {
			if as.ID == autosaveID {
				found = true
				break
			}
		}
		assert.True(t, found, "Expected created autosave to be found in list")

		// List with Fields
		fieldsList, err := revisionsAPI.List().
			Fields("id", "parent").
			Do()
		require.NoError(t, err)
		assert.NotEmpty(t, fieldsList)
		for _, item := range fieldsList {
			assert.NotEmpty(t, item.ID)
			assert.NotEmpty(t, item.Parent)
		}

		// 5. Retrieve Autosave by ID
		retrieved, err := revisionsAPI.Retrieve(autosaveID).
			ContextEdit().
			Do()
		require.NoError(t, err)
		require.NotNil(t, retrieved)
		assert.Equal(t, autosaveID, retrieved.ID)
		assert.Equal(t, parentItemID, retrieved.Parent)
		assert.NotEmpty(t, retrieved.Slug)
		assert.NotNil(t, retrieved.Date)

		// Retrieve with Fields
		fieldRetrieved, err := revisionsAPI.Retrieve(autosaveID).
			Fields("id", "parent", "slug").
			Do()
		require.NoError(t, err)
		require.NotNil(t, fieldRetrieved)
		assert.Equal(t, autosaveID, fieldRetrieved.ID)
		assert.Equal(t, parentItemID, fieldRetrieved.Parent)
		assert.NotEmpty(t, fieldRetrieved.Slug)

		// 6. Verify 404 for invalid parent ID
		invalidParentAPI := client.NavMenuItemRevisions(999999)
		_, err = invalidParentAPI.Retrieve(1).Do()
		require.Error(t, err)
		var wpErr *gowprest.WPRestError
		require.True(t, errors.As(err, &wpErr))
		assert.Equal(t, 404, wpErr.Data.Status)
		assert.Equal(t, "rest_post_invalid_parent", wpErr.Code)

		// 7. Verify unauthorized access (401) using valid parent item
		unauthClient := gowprest.NewClient(blogUrl)
		defer unauthClient.Close()

		_, err = unauthClient.NavMenuItemRevisions(parentItemID).Create().
			Title("Unauthorized Autosave").
			Do()
		require.Error(t, err)
		require.True(t, errors.As(err, &wpErr))
		assert.Equal(t, 401, wpErr.Data.Status)
		assert.Equal(t, "rest_cannot_edit", wpErr.Code)
	})
}
