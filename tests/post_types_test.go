package tests

import (
	"os"
	"testing"

	_ "github.com/joho/godotenv/autoload"
	"github.com/raitucarp/gowprest"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestPostTypes(t *testing.T) {
	blogUrl := os.Getenv("BLOG_URL")
	username := os.Getenv("BLOG_USERNAME")
	password := os.Getenv("BLOG_APP_PASSWORD")

	if blogUrl == "" || username == "" || password == "" {
		t.Skip("Skipping test; BLOG_URL, BLOG_USERNAME, or BLOG_APP_PASSWORD not set")
	}

	client := gowprest.NewClient(blogUrl).WithBasicAuth(username, password)
	defer client.Close()

	postTypesAPI := client.PostTypes()

	t.Run("TypesAlias", func(t *testing.T) {
		assert.NotNil(t, client.Types())
		assert.NotNil(t, client.PostTypes())
	})

	t.Run("ListPostTypes", func(t *testing.T) {
		t.Run("list_default", func(t *testing.T) {
			types, err := postTypesAPI.List().Do()
			require.NoError(t, err)
			assert.NotEmpty(t, types)
			assert.Contains(t, types, gowprest.PostTypePost)
			assert.Contains(t, types, gowprest.PostTypePage)
			assert.Contains(t, types, gowprest.PostTypeAttachment)

			postType := types[gowprest.PostTypePost]
			assert.Equal(t, "post", postType.Slug)
			assert.Equal(t, "Posts", postType.Name)
			assert.Equal(t, "posts", postType.RestBase)
		})

		t.Run("list_contexts", func(t *testing.T) {
			viewTypes, err := postTypesAPI.List().ContextView().Do()
			require.NoError(t, err)
			assert.NotEmpty(t, viewTypes)

			editTypes, err := postTypesAPI.List().ContextEdit().Do()
			require.NoError(t, err)
			assert.NotEmpty(t, editTypes)
			postEdit := editTypes[gowprest.PostTypePost]
			assert.NotEmpty(t, postEdit.Capabilities.EditPosts)
			assert.NotEmpty(t, postEdit.Labels.Name)

			embedTypes, err := postTypesAPI.List().ContextEmbed().Do()
			require.NoError(t, err)
			assert.NotEmpty(t, embedTypes)
		})

		t.Run("list_global_parameters", func(t *testing.T) {
			embedTypes, err := postTypesAPI.List().Embed().Do()
			require.NoError(t, err)
			assert.NotEmpty(t, embedTypes)
			assert.Contains(t, embedTypes, gowprest.PostTypePost)
		})
	})

	t.Run("RetrievePostType", func(t *testing.T) {
		t.Run("retrieve_post", func(t *testing.T) {
			pt, err := postTypesAPI.Retrieve(gowprest.PostTypePost).Do()
			require.NoError(t, err)
			require.NotNil(t, pt)
			assert.Equal(t, "post", pt.Slug)
			assert.Equal(t, "Posts", pt.Name)
			assert.Equal(t, "posts", pt.RestBase)
			assert.False(t, pt.Hierarchical)
			assert.Contains(t, pt.Taxonomies, "category")
			assert.Contains(t, pt.Taxonomies, "post_tag")
		})

		t.Run("retrieve_page", func(t *testing.T) {
			pt, err := postTypesAPI.Retrieve(gowprest.PostTypePage).Do()
			require.NoError(t, err)
			require.NotNil(t, pt)
			assert.Equal(t, "page", pt.Slug)
			assert.Equal(t, "Pages", pt.Name)
			assert.Equal(t, "pages", pt.RestBase)
			assert.True(t, pt.Hierarchical)
		})

		t.Run("retrieve_attachment", func(t *testing.T) {
			pt, err := postTypesAPI.Retrieve(gowprest.PostTypeAttachment).Do()
			require.NoError(t, err)
			require.NotNil(t, pt)
			assert.Equal(t, "attachment", pt.Slug)
			assert.Equal(t, "Media", pt.Name)
			assert.Equal(t, "media", pt.RestBase)
		})

		t.Run("retrieve_context_edit", func(t *testing.T) {
			pt, err := postTypesAPI.Retrieve(gowprest.PostTypePost).
				ContextEdit().
				Do()
			require.NoError(t, err)
			require.NotNil(t, pt)
			assert.Equal(t, "post", pt.Slug)
			assert.Equal(t, "edit_posts", pt.Capabilities.EditPosts)
			assert.Equal(t, "edit_post", pt.Capabilities.EditPost)
			assert.Equal(t, "Posts", pt.Labels.Name)
			assert.Equal(t, "Post", pt.Labels.SingularName)
			assert.True(t, pt.Visibility.ShowUI)
			assert.True(t, pt.Viewable)
			assert.NotEmpty(t, pt.Supports)
			assert.Equal(t, true, pt.Supports["title"])
		})

		t.Run("retrieve_global_parameters", func(t *testing.T) {
			pt, err := postTypesAPI.Retrieve(gowprest.PostTypePost).
				Fields("name", "slug", "rest_base").
				Embed().
				Do()
			require.NoError(t, err)
			require.NotNil(t, pt)
			assert.Equal(t, "post", pt.Slug)
			assert.Equal(t, "Posts", pt.Name)
			assert.Equal(t, "posts", pt.RestBase)
			assert.Empty(t, pt.Description)
		})

		t.Run("retrieve_not_found", func(t *testing.T) {
			pt, err := postTypesAPI.Retrieve("non_existent_type_xyz").Do()
			require.Error(t, err)
			assert.Nil(t, pt)
			wpErr, ok := err.(*gowprest.WPRestError)
			require.True(t, ok, "expected error to be *gowprest.WPRestError")
			assert.Equal(t, 404, wpErr.Data.Status)
			assert.Equal(t, "rest_type_invalid", wpErr.Code)
		})
	})
}
