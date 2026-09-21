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

func TestNavigationRevisions(t *testing.T) {
	blogUrl := os.Getenv("BLOG_URL")
	username := os.Getenv("BLOG_USERNAME")
	password := os.Getenv("BLOG_APP_PASSWORD")

	if blogUrl == "" || username == "" || password == "" {
		t.Skip("Skipping test; BLOG_URL, BLOG_USERNAME, or BLOG_APP_PASSWORD not set")
	}

	client := gowprest.NewClient(blogUrl).WithBasicAuth(username, password)
	defer client.Close()

	// Setup: Create a navigation post and update it multiple times to generate revisions
	unique := fmt.Sprintf("NavRev%d", time.Now().UnixNano()%1000000)
	nav, err := client.Navigations().Create().
		Title("Parent Nav " + unique).
		Content(`<!-- wp:navigation-link {"label":"Home","url":"/"} /-->`).
		StatusPublish().
		Do()
	require.NoError(t, err)
	defer client.Navigations().Delete(nav.ID).Force().Do()

	// Update 1
	_, err = client.Navigations().Update(nav.ID).
		Title("Parent Nav Update 1 " + unique).
		Content(`<!-- wp:navigation-link {"label":"Home","url":"/"} /--><!-- wp:navigation-link {"label":"About","url":"/about"} /-->`).
		Do()
	require.NoError(t, err)

	// Update 2
	_, err = client.Navigations().Update(nav.ID).
		Content(`<!-- wp:navigation-link {"label":"Home","url":"/"} /--><!-- wp:navigation-link {"label":"About","url":"/about"} /--><!-- wp:navigation-link {"label":"Contact","url":"/contact"} /-->`).
		Do()
	require.NoError(t, err)

	revisionsAPI := client.Navigations().Revisions(nav.ID)

	t.Run("NavigationRevisionsAliases", func(t *testing.T) {
		assert.NotNil(t, client.Navigations().Revisions(nav.ID))
		assert.NotNil(t, client.NavigationRevisions(nav.ID))
		assert.NotNil(t, revisionsAPI.List())
		assert.NotNil(t, revisionsAPI.Retrieve(1))
		assert.NotNil(t, revisionsAPI.Delete(1))
	})

	var targetRevID int

	t.Run("ListNavigationRevisions", func(t *testing.T) {
		t.Run("Default", func(t *testing.T) {
			revisions, err := revisionsAPI.List().Do()
			require.NoError(t, err)
			require.NotEmpty(t, revisions)

			targetRevID = revisions[0].ID

			for _, rev := range revisions {
				assert.Equal(t, nav.ID, rev.Parent)
				assert.NotZero(t, rev.ID)
				assert.NotEmpty(t, rev.TitleString())
			}
		})

		t.Run("PaginationAndPerPage", func(t *testing.T) {
			revisions, err := revisionsAPI.List().
				Page(1).
				PerPage(1).
				Do()
			require.NoError(t, err)
			assert.Len(t, revisions, 1)
		})

		t.Run("Ordering", func(t *testing.T) {
			revisionsAsc, err := revisionsAPI.List().
				OrderAsc().
				OrderBy("date").
				Do()
			require.NoError(t, err)

			revisionsDesc, err := revisionsAPI.List().
				OrderDesc().
				OrderBy("date").
				Do()
			require.NoError(t, err)

			if len(revisionsAsc) > 1 && len(revisionsDesc) > 1 {
				assert.NotEqual(t, revisionsAsc[0].ID, revisionsDesc[0].ID)
			}
		})

		t.Run("Contexts", func(t *testing.T) {
			for _, ctx := range []string{"view", "edit", "embed"} {
				revisions, err := revisionsAPI.List().Context(ctx).Do()
				require.NoError(t, err)
				assert.NotEmpty(t, revisions)
			}

			revView, err := revisionsAPI.List().ContextView().Do()
			require.NoError(t, err)
			assert.NotEmpty(t, revView)

			revEdit, err := revisionsAPI.List().ContextEdit().Do()
			require.NoError(t, err)
			assert.NotEmpty(t, revEdit)

			revEmbed, err := revisionsAPI.List().ContextEmbed().Do()
			require.NoError(t, err)
			assert.NotEmpty(t, revEmbed)
		})

		t.Run("Fields", func(t *testing.T) {
			revisions, err := revisionsAPI.List().
				Fields("id", "parent").
				Do()
			require.NoError(t, err)
			require.NotEmpty(t, revisions)

			for _, rev := range revisions {
				assert.NotZero(t, rev.ID)
				assert.Equal(t, nav.ID, rev.Parent)
				assert.Empty(t, rev.ContentString())
			}
		})
	})

	t.Run("RetrieveNavigationRevision", func(t *testing.T) {
		require.NotZero(t, targetRevID, "targetRevID should have been captured from list")

		t.Run("ExistingRevision", func(t *testing.T) {
			rev, err := revisionsAPI.Retrieve(targetRevID).Do()
			require.NoError(t, err)
			require.NotNil(t, rev)
			assert.Equal(t, targetRevID, rev.ID)
			assert.Equal(t, nav.ID, rev.Parent)
			assert.NotEmpty(t, rev.TitleString())
			assert.NotEmpty(t, rev.ContentString())

			// Contexts on Retrieve
			revView, err := revisionsAPI.Retrieve(targetRevID).ContextView().Do()
			require.NoError(t, err)
			assert.Equal(t, targetRevID, revView.ID)

			revEdit, err := revisionsAPI.Retrieve(targetRevID).ContextEdit().Do()
			require.NoError(t, err)
			assert.Equal(t, targetRevID, revEdit.ID)

			revEmbed, err := revisionsAPI.Retrieve(targetRevID).ContextEmbed().Do()
			require.NoError(t, err)
			assert.Equal(t, targetRevID, revEmbed.ID)
		})

		t.Run("WithFields", func(t *testing.T) {
			rev, err := revisionsAPI.Retrieve(targetRevID).
				Fields("id", "parent").
				Do()
			require.NoError(t, err)
			require.NotNil(t, rev)
			assert.Equal(t, targetRevID, rev.ID)
			assert.Equal(t, nav.ID, rev.Parent)
			assert.Empty(t, rev.ContentString())
		})

		t.Run("NotFound", func(t *testing.T) {
			rev, err := revisionsAPI.Retrieve(999999).Do()
			require.Error(t, err)
			assert.Nil(t, rev)

			var wpErr *gowprest.WPRestError
			if errors.As(err, &wpErr) {
				assert.Equal(t, 404, wpErr.Data.Status)
			} else {
				t.Fatalf("expected WPRestError, got: %v", err)
			}
		})
	})

	t.Run("DeleteNavigationRevision", func(t *testing.T) {
		require.NotZero(t, targetRevID, "targetRevID should have been captured from list")

		deleted, err := revisionsAPI.Delete(targetRevID).Do()
		require.NoError(t, err)
		require.NotNil(t, deleted)
		assert.Equal(t, targetRevID, deleted.ID)

		// Subsequent retrieve should return 404
		_, err = revisionsAPI.Retrieve(targetRevID).Do()
		require.Error(t, err)

		var wpErr *gowprest.WPRestError
		if errors.As(err, &wpErr) {
			assert.Equal(t, 404, wpErr.Data.Status)
		} else {
			t.Fatalf("expected WPRestError, got: %v", err)
		}
	})

	t.Run("UnauthorizedAccess", func(t *testing.T) {
		unauthClient := gowprest.NewClient(blogUrl)
		defer unauthClient.Close()

		_, err := unauthClient.Navigations().Revisions(nav.ID).List().Do()
		require.Error(t, err)

		var wpErr *gowprest.WPRestError
		if errors.As(err, &wpErr) {
			assert.Equal(t, "rest_cannot_read", wpErr.Code)
			assert.Equal(t, 401, wpErr.Data.Status)
		} else {
			t.Fatalf("expected WPRestError, got: %v", err)
		}
	})
}
