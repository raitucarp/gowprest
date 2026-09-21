package tests

import (
	"fmt"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/go-faker/faker/v4"
	_ "github.com/joho/godotenv/autoload"
	"github.com/raitucarp/gowprest"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var blogUrl = os.Getenv("BLOG_URL")

func getAuthenticatedClient() *gowprest.RestClient {
	return gowprest.NewClient(blogUrl).
		WithBasicAuth(
			os.Getenv("BLOG_USERNAME"),
			os.Getenv("BLOG_APP_PASSWORD"),
		)
}

func TestCreatePost(t *testing.T) {
	client := getAuthenticatedClient()
	defer client.Close()

	t.Run("create with post data struct", func(t *testing.T) {
		title := "Struct " + faker.Sentence()
		content := faker.Paragraph()
		excerpt := faker.Sentence()

		created, err := client.Posts().Create(gowprest.PostData{
			Title:   title,
			Content: content,
			Excerpt: excerpt,
			Status:  gowprest.StatusPublished,
		}).Do()

		require.NoError(t, err)
		assert.NotEmpty(t, created.ID)
		assert.Equal(t, title, created.Title.Rendered)
		assert.Equal(t, gowprest.StatusPublished, created.Status)

		// Cleanup
		_, _ = client.Posts().Delete(created.ID).Force().Do()
	})

	t.Run("create with fluent method chaining", func(t *testing.T) {
		title := "Chaining " + faker.Sentence()
		content := faker.Paragraph()
		excerpt := faker.Sentence()

		created, err := client.Posts().Create().
			Title(title).
			Content(content).
			Excerpt(excerpt).
			StatusPublish().
			Sticky(true).
			Do()

		require.NoError(t, err)
		assert.NotEmpty(t, created.ID)
		assert.Equal(t, title, created.Title.Rendered)
		assert.True(t, created.Sticky)
		assert.Equal(t, gowprest.StatusPublished, created.Status)

		// Cleanup
		_, _ = client.Posts().Delete(created.ID).Force().Do()
	})

	t.Run("create draft post", func(t *testing.T) {
		title := "Draft " + faker.Sentence()

		created, err := client.Posts().Create().
			Title(title).
			Content(faker.Paragraph()).
			StatusDraft().
			Do()

		require.NoError(t, err)
		assert.Equal(t, gowprest.StatusDraft, created.Status)

		// Cleanup
		_, _ = client.Posts().Delete(created.ID).Force().Do()
	})

	t.Run("create password protected post", func(t *testing.T) {
		title := "Protected " + faker.Sentence()
		password := "secret123"

		created, err := client.Posts().Create().
			Title(title).
			Content("Classified content here.").
			Password(password).
			StatusPublish().
			Do()

		require.NoError(t, err)
		assert.NotEmpty(t, created.ID)

		// Cleanup
		_, _ = client.Posts().Delete(created.ID).Force().Do()
	})

	t.Run("create with custom slug", func(t *testing.T) {
		customSlug := fmt.Sprintf("custom-slug-%d", time.Now().UnixNano())

		created, err := client.Posts().Create().
			Title("Custom Slug " + faker.Word()).
			Slug(customSlug).
			StatusPublish().
			Do()

		require.NoError(t, err)
		assert.Equal(t, customSlug, created.Slug)

		// Cleanup
		_, _ = client.Posts().Delete(created.ID).Force().Do()
	})

	t.Run("create post with format", func(t *testing.T) {
		created, err := client.Posts().Create().
			Title("Aside Format " + faker.Word()).
			Content(faker.Paragraph()).
			Format(gowprest.FormatAside).
			StatusPublish().
			Do()

		require.NoError(t, err)
		assert.Equal(t, gowprest.FormatAside, created.Format)

		// Cleanup
		_, _ = client.Posts().Delete(created.ID).Force().Do()
	})

	t.Run("create future scheduled post", func(t *testing.T) {
		futureTime := time.Now().Add(24 * time.Hour).UTC()

		created, err := client.Posts().Create().
			Title("Future " + faker.Sentence()).
			Content(faker.Paragraph()).
			Date(futureTime).
			StatusFuture().
			Do()

		require.NoError(t, err)
		assert.Equal(t, gowprest.StatusFuture, created.Status)

		// Cleanup
		_, _ = client.Posts().Delete(created.ID).Force().Do()
	})

	t.Run("create with categories and tags", func(t *testing.T) {
		// Category 1 is default "Uncategorized" in fresh WordPress
		created, err := client.Posts().Create().
			Title("Cat and Tag " + faker.Word()).
			Content(faker.Paragraph()).
			Categories(1).
			StatusPublish().
			Do()

		require.NoError(t, err)
		assert.Contains(t, created.Categories, 1)

		// Cleanup
		_, _ = client.Posts().Delete(created.ID).Force().Do()
	})

	t.Run("create validation error invalid author", func(t *testing.T) {
		_, err := client.Posts().Create().
			Title("Invalid Author").
			Author(999999).
			StatusPublish().
			Do()

		assert.Error(t, err)
		var wpErr *gowprest.WPRestError
		assert.ErrorAs(t, err, &wpErr)
		assert.Equal(t, 400, wpErr.Data.Status)
	})
}

func TestListPosts(t *testing.T) {
	client := getAuthenticatedClient()
	defer client.Close()

	// Seed test posts
	p1, err := client.Posts().Create().
		Title("Alpha " + faker.Word()).
		Content("Alpha content about golang development.").
		StatusPublish().
		Do()
	require.NoError(t, err)
	defer func() { _, _ = client.Posts().Delete(p1.ID).Force().Do() }()

	p2, err := client.Posts().Create().
		Title("Beta " + faker.Word()).
		Content("Beta content about wordpress rest api.").
		StatusPublish().
		Do()
	require.NoError(t, err)
	defer func() { _, _ = client.Posts().Delete(p2.ID).Force().Do() }()

	pDraft, err := client.Posts().Create().
		Title("Draft " + faker.Word()).
		Content("Draft content.").
		StatusDraft().
		Do()
	require.NoError(t, err)
	defer func() { _, _ = client.Posts().Delete(pDraft.ID).Force().Do() }()

	t.Run("default post list", func(t *testing.T) {
		posts, err := client.Posts().List().Do()
		require.NoError(t, err)
		assert.GreaterOrEqual(t, len(posts), 2)
	})

	t.Run("not found search", func(t *testing.T) {
		posts, err := client.Posts().List().Search("nonexistentqueryxyz999").Do()
		require.NoError(t, err)
		assert.Empty(t, posts)
	})

	t.Run("search matching keyword", func(t *testing.T) {
		posts, err := client.Posts().List().Search("golang development").Do()
		require.NoError(t, err)
		assert.NotEmpty(t, posts)
		assert.Contains(t, posts[0].Content.Rendered, "golang")
	})

	t.Run("search columns", func(t *testing.T) {
		posts, err := client.Posts().List().
			Search("Alpha").
			SearchColumns("post_title").
			Do()
		require.NoError(t, err)
		assert.NotEmpty(t, posts)
	})

	t.Run("pagination page and per_page", func(t *testing.T) {
		posts1, err := client.Posts().List().Page(1).PerPage(1).Do()
		require.NoError(t, err)
		assert.Len(t, posts1, 1)

		posts2, err := client.Posts().List().Page(2).PerPage(1).Do()
		require.NoError(t, err)
		assert.Len(t, posts2, 1)
		assert.NotEqual(t, posts1[0].ID, posts2[0].ID)
	})

	t.Run("pagination offset", func(t *testing.T) {
		allPosts, err := client.Posts().List().PerPage(10).Do()
		require.NoError(t, err)
		if len(allPosts) >= 2 {
			offsetPosts, err := client.Posts().List().Offset(1).PerPage(1).Do()
			require.NoError(t, err)
			assert.Len(t, offsetPosts, 1)
			assert.Equal(t, allPosts[1].ID, offsetPosts[0].ID)
		}
	})

	t.Run("order asc and desc", func(t *testing.T) {
		ascPosts, err := client.Posts().List().OrderAsc().Do()
		require.NoError(t, err)

		descPosts, err := client.Posts().List().OrderDesc().Do()
		require.NoError(t, err)

		if len(ascPosts) >= 2 && len(descPosts) >= 2 {
			assert.NotEqual(t, ascPosts[0].ID, descPosts[0].ID)
		}
	})

	t.Run("orderby title and id", func(t *testing.T) {
		byTitle, err := client.Posts().List().OrderByTitle().OrderAsc().Do()
		require.NoError(t, err)
		assert.NotEmpty(t, byTitle)

		byId, err := client.Posts().List().OrderById().OrderAsc().Do()
		require.NoError(t, err)
		assert.NotEmpty(t, byId)
	})

	t.Run("filter by status", func(t *testing.T) {
		drafts, err := client.Posts().List().StatusDraft().Do()
		require.NoError(t, err)
		assert.NotEmpty(t, drafts)
		for _, post := range drafts {
			assert.Equal(t, gowprest.StatusDraft, post.Status)
		}

		published, err := client.Posts().List().StatusPublish().Do()
		require.NoError(t, err)
		assert.NotEmpty(t, published)
	})

	t.Run("filter by slug", func(t *testing.T) {
		posts, err := client.Posts().List().Slug(p1.Slug).Do()
		require.NoError(t, err)
		require.Len(t, posts, 1)
		assert.Equal(t, p1.ID, posts[0].ID)
	})

	t.Run("filter by author and author_exclude", func(t *testing.T) {
		posts, err := client.Posts().List().Author(p1.Author).Do()
		require.NoError(t, err)
		assert.NotEmpty(t, posts)

		excluded, err := client.Posts().List().AuthorExclude(p1.Author).Do()
		require.NoError(t, err)
		for _, post := range excluded {
			assert.NotEqual(t, p1.Author, post.Author)
		}
	})

	t.Run("filter by include and exclude", func(t *testing.T) {
		included, err := client.Posts().List().Include(p1.ID, p2.ID).Do()
		require.NoError(t, err)
		assert.Len(t, included, 2)

		excluded, err := client.Posts().List().Exclude(p1.ID).Do()
		require.NoError(t, err)
		for _, post := range excluded {
			assert.NotEqual(t, p1.ID, post.ID)
		}
	})

	t.Run("filter by date after and before", func(t *testing.T) {
		past := time.Now().Add(-1 * time.Hour)
		afterPosts, err := client.Posts().List().After(past).Do()
		require.NoError(t, err)
		assert.NotEmpty(t, afterPosts)

		ancient := time.Now().Add(-1000 * time.Hour)
		beforePosts, err := client.Posts().List().Before(ancient).Do()
		require.NoError(t, err)
		assert.Empty(t, beforePosts)
	})

	t.Run("filter by sticky", func(t *testing.T) {
		stickyPost, err := client.Posts().Create().
			Title("Sticky " + faker.Word()).
			Content(faker.Paragraph()).
			Sticky(true).
			StatusPublish().
			Do()
		require.NoError(t, err)
		defer func() { _, _ = client.Posts().Delete(stickyPost.ID).Force().Do() }()

		posts, err := client.Posts().List().Sticky(true).Do()
		require.NoError(t, err)
		assert.NotEmpty(t, posts)
		assert.True(t, posts[0].Sticky)
	})

	t.Run("global parameter fields", func(t *testing.T) {
		posts, err := client.Posts().List().Fields("id", "title").Do()
		require.NoError(t, err)
		require.NotEmpty(t, posts)
		assert.NotEmpty(t, posts[0].ID)
		assert.NotNil(t, posts[0].Title)
		assert.Nil(t, posts[0].Excerpt)
	})

	t.Run("global parameter embed", func(t *testing.T) {
		posts, err := client.Posts().List().Embed().Do()
		require.NoError(t, err)
		require.NotEmpty(t, posts)
		assert.NotEmpty(t, posts[0].Links)
		assert.NotEmpty(t, posts[0].Embedded)
	})

	t.Run("context edit", func(t *testing.T) {
		posts, err := client.Posts().List().ContextEdit().Do()
		require.NoError(t, err)
		require.NotEmpty(t, posts)
		assert.NotEmpty(t, posts[0].Title.Raw)
	})
}

func TestRetrievePost(t *testing.T) {
	client := getAuthenticatedClient()
	defer client.Close()

	created, err := client.Posts().Create().
		Title("Retrieve Test " + faker.Word()).
		Content("Body for retrieve test.").
		Excerpt("Excerpt for retrieve test.").
		StatusPublish().
		Do()
	require.NoError(t, err)
	defer func() { _, _ = client.Posts().Delete(created.ID).Force().Do() }()

	t.Run("retrieve existing post", func(t *testing.T) {
		post, err := client.Posts().Retrieve(created.ID).Do()
		require.NoError(t, err)
		assert.Equal(t, created.ID, post.ID)
		assert.Equal(t, created.Title.Rendered, post.Title.Rendered)
	})

	t.Run("retrieve with context edit", func(t *testing.T) {
		post, err := client.Posts().Retrieve(created.ID).ContextEdit().Do()
		require.NoError(t, err)
		assert.Equal(t, created.ID, post.ID)
		assert.NotEmpty(t, post.Title.Raw)
		assert.NotEmpty(t, post.Content.Raw)
	})

	t.Run("retrieve with context embed", func(t *testing.T) {
		post, err := client.Posts().Retrieve(created.ID).ContextEmbed().Do()
		require.NoError(t, err)
		assert.Equal(t, created.ID, post.ID)
	})

	t.Run("retrieve with embed", func(t *testing.T) {
		post, err := client.Posts().Retrieve(created.ID).Embed().Do()
		require.NoError(t, err)
		assert.NotEmpty(t, post.Links)
		assert.NotEmpty(t, post.Embedded)
	})

	t.Run("retrieve with fields", func(t *testing.T) {
		post, err := client.Posts().Retrieve(created.ID).Fields("id", "title").Do()
		require.NoError(t, err)
		assert.Equal(t, created.ID, post.ID)
		assert.NotNil(t, post.Title)
		assert.Nil(t, post.Excerpt)
	})

	t.Run("retrieve password protected post", func(t *testing.T) {
		pwd := "supersecret"
		protected, err := client.Posts().Create().
			Title("Secret Post " + faker.Word()).
			Content("Confidential data.").
			Password(pwd).
			StatusPublish().
			Do()
		require.NoError(t, err)
		defer func() { _, _ = client.Posts().Delete(protected.ID).Force().Do() }()

		// Without password using an unauthenticated client
		publicClient := gowprest.NewClient(blogUrl)
		defer publicClient.Close()

		withoutPwd, err := publicClient.Posts().Retrieve(protected.ID).Do()
		require.NoError(t, err)
		assert.True(t, withoutPwd.Content.Protected)
		assert.Empty(t, withoutPwd.Content.Rendered)

		// With password
		withPwd, err := publicClient.Posts().Retrieve(protected.ID).Password(pwd).Do()
		require.NoError(t, err)
		assert.Contains(t, withPwd.Content.Rendered, "Confidential data")
	})

	t.Run("retrieve non-existent post", func(t *testing.T) {
		_, err := client.Posts().Retrieve(999999).Do()
		assert.Error(t, err)
		var wpErr *gowprest.WPRestError
		assert.ErrorAs(t, err, &wpErr)
		assert.Equal(t, 404, wpErr.Data.Status)
		assert.Equal(t, "rest_post_invalid_id", wpErr.Code)
	})
}

func TestUpdatePost(t *testing.T) {
	client := getAuthenticatedClient()
	defer client.Close()

	initial, err := client.Posts().Create().
		Title("Original Title " + faker.Word()).
		Content("Original content.").
		StatusPublish().
		Do()
	require.NoError(t, err)
	defer func() { _, _ = client.Posts().Delete(initial.ID).Force().Do() }()

	t.Run("update with post data struct", func(t *testing.T) {
		newTitle := "Updated Title " + faker.Word()
		newContent := "Updated content " + faker.Paragraph()

		updated, err := client.Posts().Update(gowprest.PostData{
			ID:      initial.ID,
			Title:   newTitle,
			Content: newContent,
		}).Do()

		require.NoError(t, err)
		assert.Equal(t, initial.ID, updated.ID)
		assert.Equal(t, newTitle, updated.Title.Rendered)
		assert.Contains(t, updated.Content.Rendered, strings.TrimSpace(newContent))
	})

	t.Run("update with fluent method chaining", func(t *testing.T) {
		chainedTitle := "Chained Update " + faker.Word()
		chainedExcerpt := "Chained excerpt " + faker.Sentence()

		updated, err := client.Posts().Update().
			ID(initial.ID).
			Title(chainedTitle).
			Excerpt(chainedExcerpt).
			StatusDraft().
			Sticky(true).
			Do()

		require.NoError(t, err)
		assert.Equal(t, initial.ID, updated.ID)
		assert.Equal(t, chainedTitle, updated.Title.Rendered)
		assert.Equal(t, gowprest.StatusDraft, updated.Status)
		assert.True(t, updated.Sticky)
	})

	t.Run("update comment status", func(t *testing.T) {
		updated, err := client.Posts().Update().
			ID(initial.ID).
			CommentStatus(gowprest.StatusClosed).
			PingStatus(gowprest.StatusClosed).
			Do()

		require.NoError(t, err)
		assert.Equal(t, gowprest.StatusClosed, updated.CommentStatus)
		assert.Equal(t, gowprest.StatusClosed, updated.PingStatus)
	})

	t.Run("update non-existent post", func(t *testing.T) {
		_, err := client.Posts().Update().
			ID(999999).
			Title("Should Fail").
			Do()

		assert.Error(t, err)
		var wpErr *gowprest.WPRestError
		assert.ErrorAs(t, err, &wpErr)
		assert.Equal(t, 404, wpErr.Data.Status)
	})
}

func TestDeletePost(t *testing.T) {
	client := getAuthenticatedClient()
	defer client.Close()

	t.Run("soft delete to trash and restore/verify", func(t *testing.T) {
		created, err := client.Posts().Create().
			Title("To Trash " + faker.Word()).
			Content(faker.Paragraph()).
			StatusPublish().
			Do()
		require.NoError(t, err)

		// Soft delete (move to trash)
		trashed, err := client.Posts().Delete(created.ID).Do()
		require.NoError(t, err)
		assert.Equal(t, created.ID, trashed.ID)
		assert.Equal(t, gowprest.StatusTrash, trashed.Status)

		// Retrieve trashed post in edit context
		inTrash, err := client.Posts().Retrieve(created.ID).ContextEdit().Do()
		require.NoError(t, err)
		assert.Equal(t, created.ID, inTrash.ID)
		assert.Equal(t, gowprest.StatusTrash, inTrash.Status)

		// Permanently delete from trash
		forceDeleted, err := client.Posts().Delete(created.ID).Force().Do()
		require.NoError(t, err)
		assert.Equal(t, created.ID, forceDeleted.ID)

		// Verify it no longer exists
		_, err = client.Posts().Retrieve(created.ID).Do()
		assert.Error(t, err)
	})

	t.Run("force delete directly", func(t *testing.T) {
		created, err := client.Posts().Create().
			Title("Direct Force Delete " + faker.Word()).
			Content(faker.Paragraph()).
			StatusPublish().
			Do()
		require.NoError(t, err)

		deleted, err := client.Posts().Delete(created.ID).Force().Do()
		require.NoError(t, err)
		assert.Equal(t, created.ID, deleted.ID)

		// Verify it no longer exists
		_, err = client.Posts().Retrieve(created.ID).Do()
		assert.Error(t, err)
	})

	t.Run("delete non-existent post", func(t *testing.T) {
		_, err := client.Posts().Delete(999999).Do()
		assert.Error(t, err)
		var wpErr *gowprest.WPRestError
		assert.ErrorAs(t, err, &wpErr)
		assert.Equal(t, 404, wpErr.Data.Status)
	})
}
