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

func TestNavMenus(t *testing.T) {
	blogUrl := os.Getenv("BLOG_URL")
	username := os.Getenv("BLOG_USERNAME")
	password := os.Getenv("BLOG_APP_PASSWORD")

	if blogUrl == "" || username == "" || password == "" {
		t.Skip("Skipping test; BLOG_URL, BLOG_USERNAME, or BLOG_APP_PASSWORD not set")
	}

	client := gowprest.NewClient(blogUrl).WithBasicAuth(username, password)
	defer client.Close()

	t.Run("NavMenusAliases", func(t *testing.T) {
		assert.NotNil(t, client.NavMenus())
		assert.NotNil(t, client.NavMenu())
		assert.NotNil(t, client.Menus())
		assert.NotNil(t, client.Menu())

		assert.NotNil(t, client.NavMenus().List())
		assert.NotNil(t, client.NavMenus().Retrieve(1))
		assert.NotNil(t, client.NavMenus().Create())
		assert.NotNil(t, client.NavMenus().Update(1))
		assert.NotNil(t, client.NavMenus().Delete(1))
	})

	t.Run("NavMenusCRUDWorkflow", func(t *testing.T) {
		unique := fmt.Sprintf("Menu%d", time.Now().UnixNano()%1000000)
		name := "Nav Menu " + unique
		slug := "nav-menu-" + fmt.Sprintf("%d", time.Now().UnixNano()%1000000)
		description := "Primary navigation menu description " + unique

		// 1. Create Menu
		created, err := client.NavMenus().Create().
			Name(name).
			Description(description).
			Slug(slug).
			Locations("primary").
			AutoAdd(false).
			Do()
		require.NoError(t, err)
		require.NotNil(t, created)
		assert.NotEmpty(t, created.ID)
		assert.Equal(t, name, created.Name)
		assert.Equal(t, description, created.Description)
		assert.NotEmpty(t, created.Slug)
		assert.Contains(t, created.Locations, "primary")
		assert.False(t, created.AutoAdd)

		menuID := created.ID

		// Ensure cleanup
		defer func() {
			_, _ = client.NavMenus().Delete(menuID).Force().Do()
		}()

		// 2. Retrieve Menu
		retrieved, err := client.NavMenus().Retrieve(menuID).
			ContextEdit().
			Do()
		require.NoError(t, err)
		require.NotNil(t, retrieved)
		assert.Equal(t, menuID, retrieved.ID)
		assert.Equal(t, name, retrieved.Name)
		assert.Equal(t, description, retrieved.Description)
		assert.Equal(t, created.Slug, retrieved.Slug)
		assert.Contains(t, retrieved.Locations, "primary")

		// 3. Update Menu
		updatedName := "Updated " + name
		updatedDesc := "Updated description for " + unique
		updated, err := client.NavMenus().Update(menuID).
			Name(updatedName).
			Description(updatedDesc).
			Locations("footer").
			Do()
		require.NoError(t, err)
		require.NotNil(t, updated)
		assert.Equal(t, updatedName, updated.Name)
		assert.Equal(t, updatedDesc, updated.Description)
		assert.Contains(t, updated.Locations, "footer")

		// 4. List Menus with filters
		list, err := client.NavMenus().List().
			ContextEdit().
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
			if item.ID == menuID {
				found = true
				break
			}
		}
		assert.True(t, found, "Expected created menu to be in list results")

		// List with Include
		incList, err := client.NavMenus().List().
			Include(menuID).
			Do()
		require.NoError(t, err)
		assert.Len(t, incList, 1)
		assert.Equal(t, menuID, incList[0].ID)

		// List with Exclude
		excList, err := client.NavMenus().List().
			Exclude(menuID).
			PerPage(1).
			Do()
		require.NoError(t, err)
		for _, item := range excList {
			assert.NotEqual(t, menuID, item.ID)
		}

		// List with Fields
		fieldsList, err := client.NavMenus().List().
			Include(menuID).
			Fields("id", "name").
			Do()
		require.NoError(t, err)
		require.Len(t, fieldsList, 1)
		assert.Equal(t, menuID, fieldsList[0].ID)
		assert.NotEmpty(t, fieldsList[0].Name)

		// Retrieve with Fields
		fieldRetrieved, err := client.NavMenus().Retrieve(menuID).
			Fields("id", "name", "slug").
			Do()
		require.NoError(t, err)
		require.NotNil(t, fieldRetrieved)
		assert.Equal(t, menuID, fieldRetrieved.ID)
		assert.Equal(t, updatedName, fieldRetrieved.Name)
		assert.Equal(t, updated.Slug, fieldRetrieved.Slug)

		// 5. Delete Menu
		deleted, err := client.NavMenus().Delete(menuID).
			Force().
			Do()
		require.NoError(t, err)
		require.NotNil(t, deleted)
		assert.Equal(t, menuID, deleted.ID)

		// 6. Verify 404 after deletion
		_, err = client.NavMenus().Retrieve(menuID).Do()
		require.Error(t, err)
		var wpErr *gowprest.WPRestError
		require.True(t, errors.As(err, &wpErr))
		assert.Equal(t, "rest_term_invalid", wpErr.Code)
		assert.Equal(t, 404, wpErr.Data.Status)
	})

	t.Run("NavMenusUnauthorized", func(t *testing.T) {
		unauthClient := gowprest.NewClient(blogUrl)
		defer unauthClient.Close()

		_, err := unauthClient.NavMenus().Create().
			Name("Unauthorized Menu").
			Do()
		require.Error(t, err)
		var wpErr *gowprest.WPRestError
		require.True(t, errors.As(err, &wpErr))
		assert.Equal(t, 401, wpErr.Data.Status)
	})
}
