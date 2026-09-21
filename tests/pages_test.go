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

func TestCreatePage(t *testing.T) {
	blogUrl := os.Getenv("BLOG_URL")
	username := os.Getenv("BLOG_USERNAME")
	password := os.Getenv("BLOG_APP_PASSWORD")

	if blogUrl == "" || username == "" || password == "" {
		t.Skip("Skipping test; credentials not set")
	}

	client := gowprest.NewClient(blogUrl).WithBasicAuth(username, password)
	defer client.Close()

	pageAPI := client.Pages()

	t.Run("create_with_page_data_struct", func(t *testing.T) {
		title := "Struct Page " + faker.Sentence()
		content := "<p>" + faker.Paragraph() + "</p>"
		excerpt := faker.Sentence()

		created, err := pageAPI.Create(gowprest.PageData{
			Title:         title,
			Content:       content,
			Excerpt:       excerpt,
			Status:        gowprest.StatusPublished,
			CommentStatus: gowprest.StatusOpen,
			PingStatus:    gowprest.StatusClosed,
			MenuOrder:     5,
		}).Do()

		require.NoError(t, err)
		assert.NotZero(t, created.ID)
		assert.Equal(t, title, created.Title.Rendered)
		assert.Equal(t, gowprest.StatusPublished, created.Status)
		assert.Equal(t, 5, created.MenuOrder)
		defer pageAPI.Delete(created.ID).Force().Do()
	})

	t.Run("create_with_fluent_method_chaining", func(t *testing.T) {
		title := "Fluent Page " + faker.Sentence()
		content := "<p>" + faker.Paragraph() + "</p>"
		excerpt := faker.Sentence()
		slug := "fluent-page-" + strings.ToLower(faker.Word())

		created, err := pageAPI.Create().
			Title(title).
			Content(content).
			Excerpt(excerpt).
			Slug(slug).
			StatusPublish().
			MenuOrder(10).
			Do()

		require.NoError(t, err)
		assert.NotZero(t, created.ID)
		assert.Equal(t, title, created.Title.Rendered)
		assert.Equal(t, slug, created.Slug)
		assert.Equal(t, 10, created.MenuOrder)
		defer pageAPI.Delete(created.ID).Force().Do()
	})

	t.Run("create_child_page", func(t *testing.T) {
		parent, err := pageAPI.Create().
			Title("Parent Page " + faker.Word()).
			StatusPublish().
			Do()
		require.NoError(t, err)
		defer pageAPI.Delete(parent.ID).Force().Do()

		child, err := pageAPI.Create().
			Title("Child Page " + faker.Word()).
			Parent(parent.ID).
			StatusPublish().
			Do()
		require.NoError(t, err)
		assert.Equal(t, parent.ID, child.Parent)
		defer pageAPI.Delete(child.ID).Force().Do()
	})

	t.Run("create_draft_page", func(t *testing.T) {
		page, err := pageAPI.Create().
			Title("Draft Page " + faker.Word()).
			StatusDraft().
			Do()
		require.NoError(t, err)
		assert.Equal(t, gowprest.StatusDraft, page.Status)
		defer pageAPI.Delete(page.ID).Force().Do()
	})

	t.Run("create_password_protected_page", func(t *testing.T) {
		page, err := pageAPI.Create().
			Title("Protected Page " + faker.Word()).
			Content("<p>Top secret</p>").
			Password("pass123").
			StatusPublish().
			Do()
		require.NoError(t, err)
		assert.True(t, page.Content.Protected)
		defer pageAPI.Delete(page.ID).Force().Do()
	})

	t.Run("create_validation_error_invalid_author", func(t *testing.T) {
		_, err := pageAPI.Create().
			Title("Invalid Author Page").
			Author(999999).
			Do()
		require.Error(t, err)
		var wpErr *gowprest.WPRestError
		require.ErrorAs(t, err, &wpErr)
		assert.Equal(t, 400, wpErr.Data.Status)
	})
}

func TestListPages(t *testing.T) {
	blogUrl := os.Getenv("BLOG_URL")
	username := os.Getenv("BLOG_USERNAME")
	password := os.Getenv("BLOG_APP_PASSWORD")

	if blogUrl == "" || username == "" || password == "" {
		t.Skip("Skipping test; credentials not set")
	}

	client := gowprest.NewClient(blogUrl).WithBasicAuth(username, password)
	defer client.Close()

	pageAPI := client.Pages()

	// Setup pages
	p1, err := pageAPI.Create().Title("List Page Alpha " + faker.Word()).Content("<p>Alpha content</p>").MenuOrder(1).StatusPublish().Do()
	require.NoError(t, err)
	defer pageAPI.Delete(p1.ID).Force().Do()

	p2, err := pageAPI.Create().Title("List Page Beta " + faker.Word()).Content("<p>Beta content</p>").MenuOrder(2).StatusPublish().Do()
	require.NoError(t, err)
	defer pageAPI.Delete(p2.ID).Force().Do()

	p3, err := pageAPI.Create().Title("List Page Gamma " + faker.Word()).Content("<p>Gamma content</p>").Parent(p1.ID).MenuOrder(3).StatusPublish().Do()
	require.NoError(t, err)
	defer pageAPI.Delete(p3.ID).Force().Do()

	t.Run("default_page_list", func(t *testing.T) {
		pages, err := pageAPI.List().Do()
		require.NoError(t, err)
		assert.GreaterOrEqual(t, len(pages), 3)
	})

	t.Run("not_found_search", func(t *testing.T) {
		pages, err := pageAPI.List().Search("non-existent-search-keyword-xyz-987").Do()
		require.NoError(t, err)
		assert.Empty(t, pages)
	})

	t.Run("search_matching_keyword", func(t *testing.T) {
		pages, err := pageAPI.List().Search("Alpha").Do()
		require.NoError(t, err)
		assert.NotEmpty(t, pages)
	})

	t.Run("pagination_page_and_per_page", func(t *testing.T) {
		page1, err := pageAPI.List().OrderById().OrderAsc().Page(1).PerPage(1).Do()
		require.NoError(t, err)
		assert.Len(t, page1, 1)

		page2, err := pageAPI.List().OrderById().OrderAsc().Page(2).PerPage(1).Do()
		require.NoError(t, err)
		assert.Len(t, page2, 1)
		assert.NotEqual(t, page1[0].ID, page2[0].ID)

		offset, err := pageAPI.List().OrderById().OrderAsc().Offset(1).PerPage(1).Do()
		require.NoError(t, err)
		assert.Len(t, offset, 1)
		assert.Equal(t, page2[0].ID, offset[0].ID)
	})

	t.Run("order_asc_and_desc", func(t *testing.T) {
		ascList, err := pageAPI.List().OrderAsc().OrderByDate().Do()
		require.NoError(t, err)
		require.GreaterOrEqual(t, len(ascList), 2)

		descList, err := pageAPI.List().OrderDesc().OrderBy("date").Do()
		require.NoError(t, err)
		require.GreaterOrEqual(t, len(descList), 2)

		assert.NotEqual(t, ascList[0].ID, descList[0].ID)
	})

	t.Run("orderby_menu_order_and_title", func(t *testing.T) {
		_, err := pageAPI.List().OrderByMenuOrder().OrderAsc().Do()
		require.NoError(t, err)

		_, err = pageAPI.List().OrderByTitle().Do()
		require.NoError(t, err)

		_, err = pageAPI.List().OrderByAuthor().Do()
		require.NoError(t, err)
	})

	t.Run("filter_by_parent_and_parent_exclude", func(t *testing.T) {
		byParent, err := pageAPI.List().Parent(p1.ID).Do()
		require.NoError(t, err)
		require.Len(t, byParent, 1)
		assert.Equal(t, p3.ID, byParent[0].ID)

		byParentExc, err := pageAPI.List().ParentExclude(p1.ID).Do()
		require.NoError(t, err)
		for _, p := range byParentExc {
			assert.NotEqual(t, p3.ID, p.ID)
		}
	})

	t.Run("filter_by_include_and_exclude", func(t *testing.T) {
		incList, err := pageAPI.List().Include(p1.ID, p2.ID).Do()
		require.NoError(t, err)
		assert.Len(t, incList, 2)

		excList, err := pageAPI.List().Exclude(p1.ID).Do()
		require.NoError(t, err)
		for _, p := range excList {
			assert.NotEqual(t, p1.ID, p.ID)
		}
	})

	t.Run("filter_by_slug", func(t *testing.T) {
		bySlug, err := pageAPI.List().Slug(p1.Slug).Do()
		require.NoError(t, err)
		require.Len(t, bySlug, 1)
		assert.Equal(t, p1.ID, bySlug[0].ID)
	})

	t.Run("filter_by_status", func(t *testing.T) {
		pubPages, err := pageAPI.List().StatusPublish().Do()
		require.NoError(t, err)
		assert.NotEmpty(t, pubPages)

		allPages, err := pageAPI.List().StatusAny().Do()
		require.NoError(t, err)
		assert.GreaterOrEqual(t, len(allPages), len(pubPages))
	})

	t.Run("global_parameters_fields_and_embed", func(t *testing.T) {
		fieldsList, err := pageAPI.List().Fields("id", "title").Do()
		require.NoError(t, err)
		require.NotEmpty(t, fieldsList)
		assert.NotZero(t, fieldsList[0].ID)
		assert.NotNil(t, fieldsList[0].Title)
		assert.Nil(t, fieldsList[0].Content)

		embedList, err := pageAPI.List().Embed().Do()
		require.NoError(t, err)
		assert.NotEmpty(t, embedList)
	})

	t.Run("context_edit", func(t *testing.T) {
		editList, err := pageAPI.List().ContextEdit().Do()
		require.NoError(t, err)
		require.NotEmpty(t, editList)
		assert.NotEmpty(t, editList[0].Content.Rendered)
	})
}

func TestRetrievePage(t *testing.T) {
	blogUrl := os.Getenv("BLOG_URL")
	username := os.Getenv("BLOG_USERNAME")
	password := os.Getenv("BLOG_APP_PASSWORD")

	if blogUrl == "" || username == "" || password == "" {
		t.Skip("Skipping test; credentials not set")
	}

	client := gowprest.NewClient(blogUrl).WithBasicAuth(username, password)
	defer client.Close()

	pageAPI := client.Pages()

	page, err := pageAPI.Create().
		Title("Retrieve Test Page " + faker.Word()).
		Content("<p>Detailed page content</p>").
		StatusPublish().
		Do()
	require.NoError(t, err)
	defer pageAPI.Delete(page.ID).Force().Do()

	t.Run("retrieve_existing_page", func(t *testing.T) {
		retrieved, err := pageAPI.Retrieve(page.ID).Do()
		require.NoError(t, err)
		assert.Equal(t, page.ID, retrieved.ID)
		assert.Equal(t, page.Title.Rendered, retrieved.Title.Rendered)
		assert.Equal(t, "page", retrieved.Type)
	})

	t.Run("retrieve_with_context_edit", func(t *testing.T) {
		retrieved, err := pageAPI.Retrieve(page.ID).ContextEdit().Do()
		require.NoError(t, err)
		assert.Equal(t, page.ID, retrieved.ID)
		assert.NotEmpty(t, retrieved.Content.Raw)
	})

	t.Run("retrieve_with_context_embed", func(t *testing.T) {
		retrieved, err := pageAPI.Retrieve(page.ID).ContextEmbed().Do()
		require.NoError(t, err)
		assert.Equal(t, page.ID, retrieved.ID)
	})

	t.Run("retrieve_with_embed", func(t *testing.T) {
		retrieved, err := pageAPI.Retrieve(page.ID).Embed().Do()
		require.NoError(t, err)
		assert.Equal(t, page.ID, retrieved.ID)
		assert.NotEmpty(t, retrieved.Links)
	})

	t.Run("retrieve_with_fields", func(t *testing.T) {
		retrieved, err := pageAPI.Retrieve(page.ID).Fields("id", "slug").Do()
		require.NoError(t, err)
		assert.Equal(t, page.ID, retrieved.ID)
		assert.NotEmpty(t, retrieved.Slug)
		assert.Nil(t, retrieved.Title)
	})

	t.Run("retrieve_password_protected_page", func(t *testing.T) {
		pwPage, err := pageAPI.Create().
			Title("Secret Page").
			Content("<p>Super confidential data</p>").
			Password("pagekey").
			StatusPublish().
			Do()
		require.NoError(t, err)
		defer pageAPI.Delete(pwPage.ID).Force().Do()

		// Without password query param
		retrievedNoPass, err := client.Pages().Retrieve(pwPage.ID).ContextView().Do()
		require.NoError(t, err)
		assert.True(t, retrievedNoPass.Content.Protected)

		// With password query param
		retrievedWithPass, err := client.Pages().Retrieve(pwPage.ID).Password("pagekey").Do()
		require.NoError(t, err)
		assert.Contains(t, retrievedWithPass.Content.Rendered, "Super confidential data")
	})

	t.Run("retrieve_non_existent_page", func(t *testing.T) {
		_, err := pageAPI.Retrieve(9999999).Do()
		require.Error(t, err)
		var wpErr *gowprest.WPRestError
		require.ErrorAs(t, err, &wpErr)
		assert.Equal(t, 404, wpErr.Data.Status)
	})
}

func TestUpdatePage(t *testing.T) {
	blogUrl := os.Getenv("BLOG_URL")
	username := os.Getenv("BLOG_USERNAME")
	password := os.Getenv("BLOG_APP_PASSWORD")

	if blogUrl == "" || username == "" || password == "" {
		t.Skip("Skipping test; credentials not set")
	}

	client := gowprest.NewClient(blogUrl).WithBasicAuth(username, password)
	defer client.Close()

	pageAPI := client.Pages()

	page, err := pageAPI.Create().
		Title("Orig Page " + faker.Word()).
		Content("<p>Orig content</p>").
		StatusPublish().
		Do()
	require.NoError(t, err)
	defer pageAPI.Delete(page.ID).Force().Do()

	t.Run("update_with_page_data_struct", func(t *testing.T) {
		updatedTitle := "Struct Updated " + faker.Sentence()
		updated, err := pageAPI.Update(gowprest.PageData{
			ID:    page.ID,
			Title: updatedTitle,
		}).Do()
		require.NoError(t, err)
		assert.Equal(t, page.ID, updated.ID)
		assert.Equal(t, updatedTitle, updated.Title.Rendered)
	})

	t.Run("update_with_fluent_method_chaining", func(t *testing.T) {
		chainedTitle := "Chained Updated " + faker.Sentence()
		chainedContent := "<p>New chained content " + faker.Paragraph() + "</p>"
		updated, err := pageAPI.Update().
			ID(page.ID).
			Title(chainedTitle).
			Content(chainedContent).
			MenuOrder(99).
			Do()
		require.NoError(t, err)
		assert.Equal(t, page.ID, updated.ID)
		assert.Equal(t, chainedTitle, updated.Title.Rendered)
		assert.Equal(t, 99, updated.MenuOrder)
	})

	t.Run("update_parent_hierarchy", func(t *testing.T) {
		newParent, err := pageAPI.Create().
			Title("New Parent Page " + faker.Word()).
			StatusPublish().
			Do()
		require.NoError(t, err)
		defer pageAPI.Delete(newParent.ID).Force().Do()

		updated, err := pageAPI.Update().
			ID(page.ID).
			Parent(newParent.ID).
			Do()
		require.NoError(t, err)
		assert.Equal(t, newParent.ID, updated.Parent)
	})

	t.Run("update_non_existent_page", func(t *testing.T) {
		_, err := pageAPI.Update().
			ID(9999999).
			Title("Ghost Page").
			Do()
		require.Error(t, err)
		var wpErr *gowprest.WPRestError
		require.ErrorAs(t, err, &wpErr)
		assert.Equal(t, 404, wpErr.Data.Status)
	})
}

func TestDeletePage(t *testing.T) {
	blogUrl := os.Getenv("BLOG_URL")
	username := os.Getenv("BLOG_USERNAME")
	password := os.Getenv("BLOG_APP_PASSWORD")

	if blogUrl == "" || username == "" || password == "" {
		t.Skip("Skipping test; credentials not set")
	}

	client := gowprest.NewClient(blogUrl).WithBasicAuth(username, password)
	defer client.Close()

	pageAPI := client.Pages()

	t.Run("soft_delete_to_trash_and_restore_verify", func(t *testing.T) {
		page, err := pageAPI.Create().
			Title("Trash Page " + faker.Word()).
			StatusPublish().
			Do()
		require.NoError(t, err)

		// Soft delete
		trashed, err := pageAPI.Delete(page.ID).Do()
		require.NoError(t, err)
		assert.Equal(t, page.ID, trashed.ID)
		assert.Equal(t, gowprest.StatusTrash, trashed.Status)

		// Check status in trash
		retrieved, err := pageAPI.Retrieve(page.ID).ContextEdit().Do()
		require.NoError(t, err)
		assert.Equal(t, gowprest.StatusTrash, retrieved.Status)

		// Restore
		restored, err := pageAPI.Update().
			ID(page.ID).
			StatusPublish().
			Do()
		require.NoError(t, err)
		assert.Equal(t, gowprest.StatusPublished, restored.Status)

		// Final Force Delete
		_, err = pageAPI.Delete(page.ID).Force().Do()
		require.NoError(t, err)
	})

	t.Run("force_delete_directly", func(t *testing.T) {
		page, err := pageAPI.Create().
			Title("Force Delete Page " + faker.Word()).
			StatusPublish().
			Do()
		require.NoError(t, err)

		deleted, err := pageAPI.Delete(page.ID).Force().Do()
		require.NoError(t, err)
		assert.Equal(t, page.ID, deleted.ID)
		assert.Equal(t, page.Title.Rendered, deleted.Title.Rendered)

		// Verify gone
		_, err = pageAPI.Retrieve(page.ID).Do()
		require.Error(t, err)
	})

	t.Run("delete_non_existent_page", func(t *testing.T) {
		_, err := pageAPI.Delete(9999999).Force().Do()
		require.Error(t, err)
		var wpErr *gowprest.WPRestError
		require.ErrorAs(t, err, &wpErr)
		assert.Equal(t, 404, wpErr.Data.Status)
	})
}
