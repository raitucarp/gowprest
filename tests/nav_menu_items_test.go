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

func TestNavMenuItems(t *testing.T) {
	blogUrl := os.Getenv("BLOG_URL")
	username := os.Getenv("BLOG_USERNAME")
	password := os.Getenv("BLOG_APP_PASSWORD")

	if blogUrl == "" || username == "" || password == "" {
		t.Skip("Skipping test; BLOG_URL, BLOG_USERNAME, or BLOG_APP_PASSWORD not set")
	}

	client := gowprest.NewClient(blogUrl).WithBasicAuth(username, password)
	defer client.Close()

	t.Run("NavMenuItemsAliases", func(t *testing.T) {
		assert.NotNil(t, client.NavMenuItems())
		assert.NotNil(t, client.NavMenuItem())
		assert.NotNil(t, client.MenuItems())
		assert.NotNil(t, client.MenuItem())

		api := client.NavMenuItems()
		assert.NotNil(t, api.List())
		assert.NotNil(t, api.Retrieve(1))
		assert.NotNil(t, api.Create())
		assert.NotNil(t, api.Update(1))
		assert.NotNil(t, api.Delete(1))
	})

	t.Run("NavMenuItemsCRUDWorkflow", func(t *testing.T) {
		unique := fmt.Sprintf("Item%d", time.Now().UnixNano()%1000000)

		// 1. Setup parent menu
		menu, err := client.NavMenus().Create().
			Name("Menu " + unique).
			Do()
		require.NoError(t, err)
		require.NotNil(t, menu)
		defer func() {
			_, _ = client.NavMenus().Delete(menu.ID).Force().Do()
		}()

		title := "Custom Link " + unique
		url := "https://example.com/" + unique
		description := "Link description " + unique
		attrTitle := "Link title attr " + unique

		// 2. Create Nav Menu Item
		created, err := client.NavMenuItems().Create().
			Title(title).
			URL(url).
			Menus(menu.ID).
			Status("publish").
			Description(description).
			AttrTitle(attrTitle).
			Target("_blank").
			Classes("custom-class", "highlight").
			XFN("me").
			Do()
		require.NoError(t, err)
		require.NotNil(t, created)
		assert.NotEmpty(t, created.ID)
		assert.Equal(t, title, created.TitleString())
		assert.Equal(t, url, created.URL)
		assert.Equal(t, menu.ID, created.Menus)
		assert.Equal(t, "publish", created.Status)
		assert.Equal(t, description, created.Description)
		assert.Equal(t, attrTitle, created.AttrTitle)
		assert.Equal(t, "_blank", created.Target)
		assert.Contains(t, created.Classes, "custom-class")
		assert.Contains(t, created.XFN, "me")

		itemID := created.ID

		// Ensure cleanup
		defer func() {
			_, _ = client.NavMenuItems().Delete(itemID).Force().Do()
		}()

		// 3. Retrieve Nav Menu Item
		retrieved, err := client.NavMenuItems().Retrieve(itemID).
			ContextEdit().
			Do()
		require.NoError(t, err)
		require.NotNil(t, retrieved)
		assert.Equal(t, itemID, retrieved.ID)
		assert.Equal(t, title, retrieved.TitleString())
		assert.Equal(t, url, retrieved.URL)
		assert.Equal(t, menu.ID, retrieved.Menus)
		assert.Equal(t, description, retrieved.Description)

		// 4. Update Nav Menu Item
		updatedTitle := "Updated Link " + unique
		updatedURL := "https://example.com/updated-" + unique
		updatedDesc := "Updated description " + unique
		updated, err := client.NavMenuItems().Update(itemID).
			Title(updatedTitle).
			URL(updatedURL).
			Description(updatedDesc).
			Do()
		require.NoError(t, err)
		require.NotNil(t, updated)
		assert.Equal(t, updatedTitle, updated.TitleString())
		assert.Equal(t, updatedURL, updated.URL)
		assert.Equal(t, updatedDesc, updated.Description)

		// 5. List Nav Menu Items with filters
		list, err := client.NavMenuItems().List().
			ContextEdit().
			Menus(menu.ID).
			Search(unique).
			PerPage(10).
			Page(1).
			OrderDesc().
			OrderBy("id").
			Do()
		require.NoError(t, err)
		assert.NotEmpty(t, list)
		found := false
		for _, item := range list {
			if item.ID == itemID {
				found = true
				break
			}
		}
		assert.True(t, found, "Expected created item to be in list results")

		// List with Include
		incList, err := client.NavMenuItems().List().
			Include(itemID).
			Do()
		require.NoError(t, err)
		assert.Len(t, incList, 1)
		assert.Equal(t, itemID, incList[0].ID)

		// List with Exclude
		excList, err := client.NavMenuItems().List().
			Exclude(itemID).
			PerPage(1).
			Do()
		require.NoError(t, err)
		for _, item := range excList {
			assert.NotEqual(t, itemID, item.ID)
		}

		// Retrieve with Fields
		fieldRetrieved, err := client.NavMenuItems().Retrieve(itemID).
			Fields("id", "title", "url").
			Do()
		require.NoError(t, err)
		require.NotNil(t, fieldRetrieved)
		assert.Equal(t, itemID, fieldRetrieved.ID)
		assert.Equal(t, updatedTitle, fieldRetrieved.TitleString())
		assert.Equal(t, updatedURL, fieldRetrieved.URL)

		// 6. Delete Nav Menu Item
		deleted, err := client.NavMenuItems().Delete(itemID).
			Force().
			Do()
		require.NoError(t, err)
		require.NotNil(t, deleted)
		assert.Equal(t, itemID, deleted.ID)

		// 7. Verify 404 after deletion
		_, err = client.NavMenuItems().Retrieve(itemID).Do()
		require.Error(t, err)
		var wpErr *gowprest.WPRestError
		require.True(t, errors.As(err, &wpErr))
		assert.Equal(t, "rest_post_invalid_id", wpErr.Code)
		assert.Equal(t, 404, wpErr.Data.Status)
	})

	t.Run("NavMenuItemsUnauthorized", func(t *testing.T) {
		unauthClient := gowprest.NewClient(blogUrl)
		defer unauthClient.Close()

		_, err := unauthClient.NavMenuItems().Create().
			Title("Unauthorized Item").
			URL("https://example.com/unauth").
			Menus(1).
			Do()
		require.Error(t, err)
		var wpErr *gowprest.WPRestError
		require.True(t, errors.As(err, &wpErr))
		assert.Equal(t, 401, wpErr.Data.Status)
		assert.Equal(t, "rest_cannot_create", wpErr.Code)
	})
}
