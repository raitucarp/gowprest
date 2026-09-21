package tests

import (
	"os"
	"testing"

	_ "github.com/joho/godotenv/autoload"
	"github.com/raitucarp/gowprest"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestStatuses(t *testing.T) {
	blogUrl := os.Getenv("BLOG_URL")
	username := os.Getenv("BLOG_USERNAME")
	password := os.Getenv("BLOG_APP_PASSWORD")

	if blogUrl == "" || username == "" || password == "" {
		t.Skip("Skipping test; BLOG_URL, BLOG_USERNAME, or BLOG_APP_PASSWORD not set")
	}

	client := gowprest.NewClient(blogUrl).WithBasicAuth(username, password)
	defer client.Close()

	statusesAPI := client.Statuses()

	t.Run("StatusesAlias", func(t *testing.T) {
		assert.NotNil(t, client.Statuses())
		assert.NotNil(t, client.PostStatuses())
	})

	t.Run("ListStatuses", func(t *testing.T) {
		t.Run("list_default", func(t *testing.T) {
			statuses, err := statusesAPI.List().Do()
			require.NoError(t, err)
			assert.NotEmpty(t, statuses)

			assert.Contains(t, statuses, string(gowprest.StatusPublished))
			assert.Contains(t, statuses, string(gowprest.StatusDraft))
			assert.Contains(t, statuses, string(gowprest.StatusPending))
			assert.Contains(t, statuses, string(gowprest.StatusFuture))
			assert.Contains(t, statuses, string(gowprest.StatusPrivate))
			assert.Contains(t, statuses, string(gowprest.StatusTrash))

			publish := statuses[string(gowprest.StatusPublished)]
			assert.Equal(t, "Published", publish.Name)
			assert.Equal(t, "publish", publish.Slug)
			assert.True(t, publish.Public)
			assert.True(t, publish.Queryable)
		})

		t.Run("list_contexts", func(t *testing.T) {
			viewStatuses, err := statusesAPI.List().ContextView().Do()
			require.NoError(t, err)
			assert.NotEmpty(t, viewStatuses)

			editStatuses, err := statusesAPI.List().ContextEdit().Do()
			require.NoError(t, err)
			assert.NotEmpty(t, editStatuses)
			publishEdit := editStatuses[string(gowprest.StatusPublished)]
			assert.False(t, publishEdit.Private)
			assert.False(t, publishEdit.Protected)
			assert.True(t, publishEdit.ShowInList)

			embedStatuses, err := statusesAPI.List().ContextEmbed().Do()
			require.NoError(t, err)
			assert.NotEmpty(t, embedStatuses)
		})

		t.Run("list_global_parameters", func(t *testing.T) {
			embedStatuses, err := statusesAPI.List().Embed().Do()
			require.NoError(t, err)
			assert.NotEmpty(t, embedStatuses)
			assert.Contains(t, embedStatuses, string(gowprest.StatusPublished))
		})
	})

	t.Run("RetrieveStatus", func(t *testing.T) {
		t.Run("retrieve_publish", func(t *testing.T) {
			st, err := statusesAPI.Retrieve(string(gowprest.StatusPublished)).Do()
			require.NoError(t, err)
			require.NotNil(t, st)
			assert.Equal(t, "publish", st.Slug)
			assert.Equal(t, "Published", st.Name)
			assert.True(t, st.Public)
			assert.True(t, st.Queryable)
			assert.False(t, st.DateFloating)
		})

		t.Run("retrieve_draft", func(t *testing.T) {
			st, err := statusesAPI.Retrieve(string(gowprest.StatusDraft)).Do()
			require.NoError(t, err)
			require.NotNil(t, st)
			assert.Equal(t, "draft", st.Slug)
			assert.Equal(t, "Draft", st.Name)
			assert.False(t, st.Public)
			assert.False(t, st.Queryable)
			assert.True(t, st.DateFloating)
		})

		t.Run("retrieve_context_edit", func(t *testing.T) {
			st, err := statusesAPI.Retrieve(string(gowprest.StatusPublished)).
				ContextEdit().
				Do()
			require.NoError(t, err)
			require.NotNil(t, st)
			assert.Equal(t, "publish", st.Slug)
			assert.False(t, st.Private)
			assert.False(t, st.Protected)
			assert.True(t, st.ShowInList)
			assert.True(t, st.Public)
		})

		t.Run("retrieve_global_parameters", func(t *testing.T) {
			st, err := statusesAPI.Retrieve(string(gowprest.StatusPublished)).
				Fields("name", "slug").
				Embed().
				Do()
			require.NoError(t, err)
			require.NotNil(t, st)
			assert.Equal(t, "publish", st.Slug)
			assert.Equal(t, "Published", st.Name)
			// Public should be empty/zero when limited by _fields
			assert.False(t, st.Public)
		})

		t.Run("retrieve_not_found", func(t *testing.T) {
			st, err := statusesAPI.Retrieve("non_existent_status_xyz").Do()
			require.Error(t, err)
			assert.Nil(t, st)
			wpErr, ok := err.(*gowprest.WPRestError)
			require.True(t, ok, "expected error to be *gowprest.WPRestError")
			assert.Equal(t, 404, wpErr.Data.Status)
			assert.Equal(t, "rest_status_invalid", wpErr.Code)
		})
	})
}
