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

func TestNavigations(t *testing.T) {
	blogUrl := os.Getenv("BLOG_URL")
	username := os.Getenv("BLOG_USERNAME")
	password := os.Getenv("BLOG_APP_PASSWORD")

	if blogUrl == "" || username == "" || password == "" {
		t.Skip("Skipping test; BLOG_URL, BLOG_USERNAME, or BLOG_APP_PASSWORD not set")
	}

	client := gowprest.NewClient(blogUrl).WithBasicAuth(username, password)
	defer client.Close()

	t.Run("NavigationsAliases", func(t *testing.T) {
		assert.NotNil(t, client.Navigations())
		assert.NotNil(t, client.Navigation())
		assert.NotNil(t, client.Navigations().List())
		assert.NotNil(t, client.Navigations().Retrieve(1))
		assert.NotNil(t, client.Navigations().Create())
		assert.NotNil(t, client.Navigations().Update(1))
		assert.NotNil(t, client.Navigations().Delete(1))
	})

	t.Run("NavigationsCRUDWorkflow", func(t *testing.T) {
		unique := fmt.Sprintf("Nav%d", time.Now().UnixNano()%1000000)
		title := "Test Navigation " + unique
		slug := "test-nav-" + fmt.Sprintf("%d", time.Now().UnixNano()%1000000)
		content := `<!-- wp:navigation-link {"label":"Home","url":"/"} /-->`

		// 1. Create Navigation
		created, err := client.Navigations().Create().
			Title(title).
			Content(content).
			Slug(slug).
			StatusPublish().
			Do()
		require.NoError(t, err)
		require.NotNil(t, created)
		assert.NotEmpty(t, created.ID)
		assert.Equal(t, title, created.TitleString())
		assert.Contains(t, created.ContentString(), "wp:navigation-link")
		assert.Equal(t, gowprest.StatusPublished, created.Status)
		assert.Equal(t, "wp_navigation", created.Type)

		navID := created.ID

		// Ensure cleanup if test fails early
		defer func() {
			_, _ = client.Navigations().Delete(navID).Force().Do()
		}()

		// 2. Retrieve Navigation
		retrieved, err := client.Navigations().Retrieve(navID).
			ContextEdit().
			Do()
		require.NoError(t, err)
		require.NotNil(t, retrieved)
		assert.Equal(t, navID, retrieved.ID)
		assert.Equal(t, title, retrieved.TitleString())
		assert.Contains(t, retrieved.ContentString(), "wp:navigation-link")

		// 3. Update Navigation
		updatedTitle := "Updated " + title
		updatedContent := `<!-- wp:navigation-link {"label":"Home","url":"/"} /--><!-- wp:navigation-link {"label":"About","url":"/about"} /-->`
		updated, err := client.Navigations().Update(navID).
			Title(updatedTitle).
			Content(updatedContent).
			Do()
		require.NoError(t, err)
		require.NotNil(t, updated)
		assert.Equal(t, updatedTitle, updated.TitleString())
		assert.Contains(t, updated.ContentString(), "About")

		// 4. List Navigations with filters
		list, err := client.Navigations().List().
			ContextEdit().
			Search(unique).
			PerPage(10).
			Page(1).
			OrderDesc().
			OrderBy("id").
			Status(gowprest.StatusPublished).
			Do()
		require.NoError(t, err)
		assert.NotEmpty(t, list)
		found := false
		for _, item := range list {
			if item.ID == navID {
				found = true
				break
			}
		}
		assert.True(t, found, "Expected created navigation to be found in search list")

		// List with Include
		incList, err := client.Navigations().List().
			Include(navID).
			Do()
		require.NoError(t, err)
		assert.Len(t, incList, 1)
		assert.Equal(t, navID, incList[0].ID)

		// List with Exclude
		excList, err := client.Navigations().List().
			Exclude(navID).
			PerPage(1).
			Do()
		require.NoError(t, err)
		for _, item := range excList {
			assert.NotEqual(t, navID, item.ID)
		}

		// List with Fields
		fieldsList, err := client.Navigations().List().
			Include(navID).
			Fields("id", "title").
			Do()
		require.NoError(t, err)
		require.Len(t, fieldsList, 1)
		assert.Equal(t, navID, fieldsList[0].ID)
		assert.NotEmpty(t, fieldsList[0].TitleString())

		// Retrieve with Fields
		fieldRetrieved, err := client.Navigations().Retrieve(navID).
			Fields("id", "slug").
			Do()
		require.NoError(t, err)
		require.NotNil(t, fieldRetrieved)
		assert.Equal(t, navID, fieldRetrieved.ID)
		assert.NotEmpty(t, fieldRetrieved.Slug)

		// 5. Delete Navigation with Force
		deleted, err := client.Navigations().Delete(navID).
			Force().
			Do()
		require.NoError(t, err)
		require.NotNil(t, deleted)
		assert.Equal(t, navID, deleted.ID)

		// 6. Verify 404 after deletion
		_, err = client.Navigations().Retrieve(navID).Do()
		require.Error(t, err)
		var wpErr *gowprest.WPRestError
		require.True(t, errors.As(err, &wpErr))
		assert.Equal(t, "rest_post_invalid_id", wpErr.Code)
		assert.Equal(t, 404, wpErr.Data.Status)
	})

	t.Run("NavigationsUnauthorized", func(t *testing.T) {
		unauthClient := gowprest.NewClient(blogUrl)
		defer unauthClient.Close()

		_, err := unauthClient.Navigations().Create().
			Title("Unauthorized Nav").
			Content("Test").
			StatusPublish().
			Do()
		require.Error(t, err)
		var wpErr *gowprest.WPRestError
		require.True(t, errors.As(err, &wpErr))
		assert.Equal(t, 401, wpErr.Data.Status)
	})
}
