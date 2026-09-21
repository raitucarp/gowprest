package tests

import (
	"os"
	"strings"
	"testing"

	"github.com/go-faker/faker/v4"
	_ "github.com/joho/godotenv/autoload"
	"github.com/raitucarp/gowprest"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestTagsAPI(t *testing.T) {
	blogUrl := os.Getenv("BLOG_URL")
	username := os.Getenv("BLOG_USERNAME")
	password := os.Getenv("BLOG_APP_PASSWORD")

	if blogUrl == "" || username == "" || password == "" {
		t.Skip("Skipping test; BLOG_URL, BLOG_USERNAME, or BLOG_APP_PASSWORD not set")
	}

	client := gowprest.NewClient(blogUrl).WithBasicAuth(username, password)
	defer client.Close()

	tagAPI := client.Tags()

	t.Run("CreateTag", func(t *testing.T) {
		t.Run("create_with_struct", func(t *testing.T) {
			name := "Tag Struct " + faker.Word()
			desc := faker.Sentence()
			created, err := tagAPI.Create(gowprest.TagData{
				Name:        name,
				Description: desc,
			}).Do()
			require.NoError(t, err)
			assert.NotZero(t, created.ID)
			assert.Equal(t, name, created.Name)
			assert.Equal(t, desc, created.Description)
			assert.Equal(t, "post_tag", created.Taxonomy)
			defer tagAPI.Delete(created.ID).Do()
		})

		t.Run("create_with_method_chaining", func(t *testing.T) {
			name := "Tag Chain " + faker.Word()
			desc := faker.Sentence()
			slug := "tag-chain-" + strings.ToLower(faker.Word())
			created, err := tagAPI.Create().
				Name(name).
				Description(desc).
				Slug(slug).
				Do()
			require.NoError(t, err)
			assert.NotZero(t, created.ID)
			assert.Equal(t, name, created.Name)
			assert.Equal(t, desc, created.Description)
			assert.Equal(t, slug, created.Slug)
			defer tagAPI.Delete(created.ID).Do()
		})

		t.Run("create_validation_error_empty_name", func(t *testing.T) {
			_, err := tagAPI.Create().Do()
			require.Error(t, err)
			var wpErr *gowprest.WPRestError
			require.ErrorAs(t, err, &wpErr)
			assert.Equal(t, 400, wpErr.Data.Status)
		})
	})

	t.Run("ListTags", func(t *testing.T) {
		// Setup tags
		t1, err := tagAPI.Create().Name("Alpha " + faker.Word()).Do()
		require.NoError(t, err)
		defer tagAPI.Delete(t1.ID).Do()

		t2, err := tagAPI.Create().Name("Beta " + faker.Word()).Do()
		require.NoError(t, err)
		defer tagAPI.Delete(t2.ID).Do()

		t3, err := tagAPI.Create().Name("Gamma " + faker.Word()).Do()
		require.NoError(t, err)
		defer tagAPI.Delete(t3.ID).Do()

		t.Run("default_list", func(t *testing.T) {
			tags, err := tagAPI.List().Do()
			require.NoError(t, err)
			assert.GreaterOrEqual(t, len(tags), 3)
		})

		t.Run("pagination_page_and_per_page", func(t *testing.T) {
			page1, err := tagAPI.List().OrderById().OrderAsc().Page(1).PerPage(1).Do()
			require.NoError(t, err)
			assert.Len(t, page1, 1)

			page2, err := tagAPI.List().OrderById().OrderAsc().Page(2).PerPage(1).Do()
			require.NoError(t, err)
			assert.Len(t, page2, 1)
			assert.NotEqual(t, page1[0].ID, page2[0].ID)

			multi, err := tagAPI.List().OrderById().OrderAsc().Page(1).PerPage(2).Do()
			require.NoError(t, err)
			assert.Len(t, multi, 2)
			assert.Equal(t, page1[0].ID, multi[0].ID)
			assert.Equal(t, page2[0].ID, multi[1].ID)

			// Offset
			offsetList, err := tagAPI.List().OrderById().OrderAsc().Offset(1).PerPage(1).Do()
			require.NoError(t, err)
			assert.Len(t, offsetList, 1)
			assert.Equal(t, page2[0].ID, offsetList[0].ID)
		})

		t.Run("ordering_asc_and_desc", func(t *testing.T) {
			ascList, err := tagAPI.List().OrderAsc().OrderBy("name").Do()
			require.NoError(t, err)
			require.GreaterOrEqual(t, len(ascList), 2)

			descList, err := tagAPI.List().OrderDesc().OrderByName().Do()
			require.NoError(t, err)
			require.GreaterOrEqual(t, len(descList), 2)

			assert.NotEqual(t, ascList[0].ID, descList[0].ID)

			// OrderBy helpers
			_, err = tagAPI.List().OrderById().Do()
			require.NoError(t, err)
			_, err = tagAPI.List().OrderBySlug().Do()
			require.NoError(t, err)
			_, err = tagAPI.List().OrderByCount().Do()
			require.NoError(t, err)
			_, err = tagAPI.List().OrderByDescription().Do()
			require.NoError(t, err)
		})

		t.Run("filter_by_include_and_exclude", func(t *testing.T) {
			incList, err := tagAPI.List().Include(t1.ID, t2.ID).Do()
			require.NoError(t, err)
			assert.Len(t, incList, 2)

			excList, err := tagAPI.List().Exclude(t1.ID).Do()
			require.NoError(t, err)
			for _, tag := range excList {
				assert.NotEqual(t, t1.ID, tag.ID)
			}
		})

		t.Run("filter_by_slug", func(t *testing.T) {
			bySlug, err := tagAPI.List().Slug(t1.Slug).Do()
			require.NoError(t, err)
			require.Len(t, bySlug, 1)
			assert.Equal(t, t1.ID, bySlug[0].ID)
		})

		t.Run("search_query", func(t *testing.T) {
			searchList, err := tagAPI.List().Search("Alpha").Do()
			require.NoError(t, err)
			assert.NotEmpty(t, searchList)
		})

		t.Run("filter_by_post", func(t *testing.T) {
			post, err := client.Posts().Create().
				Title("Post with Tag " + faker.Word()).
				Tags(t1.ID).
				StatusPublish().
				Do()
			require.NoError(t, err)
			defer client.Posts().Delete(post.ID).Force().Do()

			postTags, err := tagAPI.List().Post(post.ID).Do()
			require.NoError(t, err)
			assert.Len(t, postTags, 1)
			assert.Equal(t, t1.ID, postTags[0].ID)
		})

		t.Run("context_and_global_parameters", func(t *testing.T) {
			editList, err := tagAPI.List().ContextEdit().Do()
			require.NoError(t, err)
			assert.NotEmpty(t, editList)

			embedList, err := tagAPI.List().ContextEmbed().Embed().Do()
			require.NoError(t, err)
			assert.NotEmpty(t, embedList)

			fieldsList, err := tagAPI.List().Fields("id", "name").Do()
			require.NoError(t, err)
			require.NotEmpty(t, fieldsList)
			assert.NotZero(t, fieldsList[0].ID)
			assert.NotEmpty(t, fieldsList[0].Name)
			assert.Empty(t, fieldsList[0].Slug)
		})
	})

	t.Run("RetrieveTag", func(t *testing.T) {
		tag, err := tagAPI.Create().
			Name("Retrieve " + faker.Word()).
			Description("For retrieve testing").
			Do()
		require.NoError(t, err)
		defer tagAPI.Delete(tag.ID).Do()

		t.Run("retrieve_existing_tag", func(t *testing.T) {
			retrieved, err := tagAPI.Retrieve(tag.ID).Do()
			require.NoError(t, err)
			assert.Equal(t, tag.ID, retrieved.ID)
			assert.Equal(t, tag.Name, retrieved.Name)
			assert.Equal(t, tag.Description, retrieved.Description)
			assert.Equal(t, "post_tag", retrieved.Taxonomy)
		})

		t.Run("retrieve_with_contexts", func(t *testing.T) {
			viewTag, err := tagAPI.Retrieve(tag.ID).ContextView().Do()
			require.NoError(t, err)
			assert.Equal(t, tag.ID, viewTag.ID)

			editTag, err := tagAPI.Retrieve(tag.ID).ContextEdit().Do()
			require.NoError(t, err)
			assert.Equal(t, tag.ID, editTag.ID)

			embedTag, err := tagAPI.Retrieve(tag.ID).ContextEmbed().Do()
			require.NoError(t, err)
			assert.Equal(t, tag.ID, embedTag.ID)
		})

		t.Run("retrieve_with_fields_and_embed", func(t *testing.T) {
			retrieved, err := tagAPI.Retrieve(tag.ID).Fields("id", "slug").Embed().Do()
			require.NoError(t, err)
			assert.Equal(t, tag.ID, retrieved.ID)
			assert.NotEmpty(t, retrieved.Slug)
			assert.Empty(t, retrieved.Name)
		})

		t.Run("retrieve_not_found", func(t *testing.T) {
			_, err := tagAPI.Retrieve(9999999).Do()
			require.Error(t, err)
			var wpErr *gowprest.WPRestError
			require.ErrorAs(t, err, &wpErr)
			assert.Equal(t, 404, wpErr.Data.Status)
		})
	})

	t.Run("UpdateTag", func(t *testing.T) {
		tag, err := tagAPI.Create().Name("Orig " + faker.Word()).Do()
		require.NoError(t, err)
		defer tagAPI.Delete(tag.ID).Do()

		t.Run("update_with_struct", func(t *testing.T) {
			updatedName := "Struct Updated " + faker.Word()
			updated, err := tagAPI.Update(gowprest.TagData{
				ID:   tag.ID,
				Name: updatedName,
			}).Do()
			require.NoError(t, err)
			assert.Equal(t, tag.ID, updated.ID)
			assert.Equal(t, updatedName, updated.Name)
		})

		t.Run("update_with_method_chaining", func(t *testing.T) {
			chainedName := "Chained Updated " + faker.Word()
			chainedDesc := faker.Sentence()
			updated, err := tagAPI.Update().
				ID(tag.ID).
				Name(chainedName).
				Description(chainedDesc).
				Do()
			require.NoError(t, err)
			assert.Equal(t, tag.ID, updated.ID)
			assert.Equal(t, chainedName, updated.Name)
			assert.Equal(t, chainedDesc, updated.Description)
		})

		t.Run("update_not_found", func(t *testing.T) {
			_, err := tagAPI.Update().
				ID(9999999).
				Name("Ghost Tag").
				Do()
			require.Error(t, err)
			var wpErr *gowprest.WPRestError
			require.ErrorAs(t, err, &wpErr)
			assert.Equal(t, 404, wpErr.Data.Status)
		})
	})

	t.Run("DeleteTag", func(t *testing.T) {
		tag, err := tagAPI.Create().Name("To Delete " + faker.Word()).Do()
		require.NoError(t, err)

		t.Run("force_delete_success", func(t *testing.T) {
			deleted, err := tagAPI.Delete(tag.ID).Force().Do()
			require.NoError(t, err)
			assert.Equal(t, tag.ID, deleted.ID)

			// Verify gone
			_, err = tagAPI.Retrieve(tag.ID).Do()
			require.Error(t, err)
		})

		t.Run("delete_not_found", func(t *testing.T) {
			_, err := tagAPI.Delete(tag.ID).Force().Do()
			require.Error(t, err)
			var wpErr *gowprest.WPRestError
			require.ErrorAs(t, err, &wpErr)
			assert.Equal(t, 404, wpErr.Data.Status)
		})
	})
}
