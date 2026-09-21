package tests

import (
	"os"
	"testing"
	"time"

	"github.com/go-faker/faker/v4"
	_ "github.com/joho/godotenv/autoload"
	"github.com/raitucarp/gowprest"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestComments(t *testing.T) {
	blogUrl := os.Getenv("BLOG_URL")
	username := os.Getenv("BLOG_USERNAME")
	password := os.Getenv("BLOG_APP_PASSWORD")

	if blogUrl == "" || username == "" || password == "" {
		t.Skip("Skipping test; BLOG_URL, BLOG_USERNAME, or BLOG_APP_PASSWORD not set")
	}

	client := gowprest.NewClient(blogUrl).WithBasicAuth(username, password)
	defer client.Close()

	commentsAPI := client.Comments()
	postsAPI := client.Posts()

	// Setup: Create a post for testing comments
	post, err := postsAPI.Create().
		Title("Comment Test Post " + faker.Sentence()).
		Content("<p>Post content for comment tests: " + faker.Paragraph() + "</p>").
		StatusPublish().
		Do()
	require.NoError(t, err)
	defer postsAPI.Delete(post.ID).Force().Do()

	var comment1 gowprest.Comment
	var comment2 gowprest.Comment
	var replyComment gowprest.Comment
	uniqueKeyword := "UniqueKw" + faker.Word()

	t.Run("CreateComment", func(t *testing.T) {
		t.Run("create_with_struct", func(t *testing.T) {
			c1Content := "Comment 1: " + uniqueKeyword + " " + faker.Paragraph()
			c, err := commentsAPI.Create(gowprest.CommentData{
				Post:        post.ID,
				Content:     c1Content,
				AuthorName:  "Author One",
				AuthorEmail: "author1@example.com",
				Status:      string(gowprest.CommentStatusApprove),
			}).Do()
			require.NoError(t, err)
			assert.Equal(t, post.ID, c.Post)
			assert.Equal(t, "Author One", c.AuthorName)
			assert.Contains(t, c.Content.Rendered, c1Content)
			comment1 = c
		})

		t.Run("create_with_method_chaining", func(t *testing.T) {
			c2Content := "Comment 2: " + faker.Paragraph()
			c, err := commentsAPI.Create().
				Post(post.ID).
				Content(c2Content).
				AuthorName("Author Two").
				AuthorEmail("author2@example.com").
				AuthorURL("https://example.com/author2").
				StatusApprove().
				Do()
			require.NoError(t, err)
			assert.Equal(t, post.ID, c.Post)
			assert.Equal(t, "Author Two", c.AuthorName)
			assert.Equal(t, "https://example.com/author2", c.AuthorURL)
			assert.Contains(t, c.Content.Rendered, c2Content)
			comment2 = c
		})

		t.Run("create_child_reply_comment", func(t *testing.T) {
			replyContent := "Reply to Comment 1: " + faker.Paragraph()
			c, err := commentsAPI.Create().
				Post(post.ID).
				Parent(comment1.ID).
				Content(replyContent).
				AuthorName("Replier").
				AuthorEmail("replier@example.com").
				StatusApprove().
				Do()
			require.NoError(t, err)
			assert.Equal(t, post.ID, c.Post)
			assert.Equal(t, comment1.ID, c.Parent)
			assert.Contains(t, c.Content.Rendered, replyContent)
			replyComment = c
		})
	})

	t.Run("ListComments", func(t *testing.T) {
		t.Run("default_list", func(t *testing.T) {
			comments, err := commentsAPI.List().Post(post.ID).Do()
			require.NoError(t, err)
			assert.GreaterOrEqual(t, len(comments), 3)
			for _, c := range comments {
				assert.Equal(t, post.ID, c.Post)
			}
		})

		t.Run("pagination_page_and_per_page", func(t *testing.T) {
			page1, err := commentsAPI.List().Post(post.ID).Page(1).PerPage(1).Do()
			require.NoError(t, err)
			assert.Len(t, page1, 1)

			page2, err := commentsAPI.List().Post(post.ID).Page(2).PerPage(1).Do()
			require.NoError(t, err)
			assert.Len(t, page2, 1)
			assert.NotEqual(t, page1[0].ID, page2[0].ID)

			// Offset
			offsetList, err := commentsAPI.List().Post(post.ID).Offset(1).PerPage(1).Do()
			require.NoError(t, err)
			assert.Len(t, offsetList, 1)
			assert.Equal(t, page2[0].ID, offsetList[0].ID)
		})

		t.Run("filter_by_date", func(t *testing.T) {
			afterDate := time.Now().Add(-24 * time.Hour)
			beforeDate := time.Now().Add(24 * time.Hour)
			dateList, err := commentsAPI.List().
				Post(post.ID).
				After(afterDate).
				Before(beforeDate).
				Do()
			require.NoError(t, err)
			assert.NotEmpty(t, dateList)
		})

		t.Run("ordering_asc_and_desc", func(t *testing.T) {
			ascList, err := commentsAPI.List().Post(post.ID).OrderAsc().OrderByID().Do()
			require.NoError(t, err)
			require.GreaterOrEqual(t, len(ascList), 2)

			descList, err := commentsAPI.List().Post(post.ID).OrderDesc().OrderBy("id").Do()
			require.NoError(t, err)
			require.GreaterOrEqual(t, len(descList), 2)

			assert.NotEqual(t, ascList[0].ID, descList[0].ID)
			assert.Equal(t, ascList[0].ID, descList[len(descList)-1].ID)

			// OrderByDateGMT
			gmtList, err := commentsAPI.List().Post(post.ID).OrderDesc().OrderByDateGMT().Do()
			require.NoError(t, err)
			assert.NotEmpty(t, gmtList)
		})

		t.Run("filter_by_post_and_parent", func(t *testing.T) {
			// By parent
			replies, err := commentsAPI.List().Post(post.ID).Parent(comment1.ID).Do()
			require.NoError(t, err)
			assert.Len(t, replies, 1)
			assert.Equal(t, replyComment.ID, replies[0].ID)

			// By parent exclude
			noReplies, err := commentsAPI.List().Post(post.ID).ParentExclude(comment1.ID).Do()
			require.NoError(t, err)
			for _, c := range noReplies {
				assert.NotEqual(t, comment1.ID, c.Parent)
			}
		})

		t.Run("filter_by_include_and_exclude", func(t *testing.T) {
			// Include
			included, err := commentsAPI.List().Include(comment1.ID, comment2.ID).Do()
			require.NoError(t, err)
			assert.Len(t, included, 2)

			// Exclude
			excluded, err := commentsAPI.List().Post(post.ID).Exclude(comment1.ID).Do()
			require.NoError(t, err)
			for _, c := range excluded {
				assert.NotEqual(t, comment1.ID, c.ID)
			}
		})

		t.Run("filter_by_author_and_email", func(t *testing.T) {
			byEmail, err := commentsAPI.List().Post(post.ID).AuthorEmail("author1@example.com").Do()
			require.NoError(t, err)
			assert.Len(t, byEmail, 1)
			assert.Equal(t, comment1.ID, byEmail[0].ID)
		})

		t.Run("filter_by_status", func(t *testing.T) {
			approved, err := commentsAPI.List().Post(post.ID).StatusApprove().Do()
			require.NoError(t, err)
			assert.GreaterOrEqual(t, len(approved), 3)

			allStatus, err := commentsAPI.List().Post(post.ID).Status("all").Do()
			require.NoError(t, err)
			assert.GreaterOrEqual(t, len(allStatus), 3)
		})

		t.Run("search_query", func(t *testing.T) {
			searchRes, err := commentsAPI.List().Post(post.ID).Search(uniqueKeyword).Do()
			require.NoError(t, err)
			assert.NotEmpty(t, searchRes)
			assert.Equal(t, comment1.ID, searchRes[0].ID)
		})

		t.Run("context_and_global_parameters", func(t *testing.T) {
			// Context edit
			editList, err := commentsAPI.List().Post(post.ID).ContextEdit().Do()
			require.NoError(t, err)
			require.NotEmpty(t, editList)
			assert.NotEmpty(t, editList[0].AuthorEmail)

			// Context embed
			embedList, err := commentsAPI.List().Post(post.ID).ContextEmbed().Embed().Do()
			require.NoError(t, err)
			require.NotEmpty(t, embedList)

			// Fields filter
			fieldsList, err := commentsAPI.List().Post(post.ID).Fields("id", "post").Do()
			require.NoError(t, err)
			require.NotEmpty(t, fieldsList)
			assert.NotZero(t, fieldsList[0].ID)
			assert.Equal(t, post.ID, fieldsList[0].Post)
			assert.Nil(t, fieldsList[0].Content)
		})
	})

	t.Run("RetrieveComment", func(t *testing.T) {
		t.Run("retrieve_existing_comment", func(t *testing.T) {
			c, err := commentsAPI.Retrieve(comment1.ID).Do()
			require.NoError(t, err)
			assert.Equal(t, comment1.ID, c.ID)
			assert.Equal(t, post.ID, c.Post)
			if c.Content != nil {
				assert.NotEmpty(t, c.Content.Rendered)
			}
		})

		t.Run("retrieve_with_contexts", func(t *testing.T) {
			// Context view
			viewComment, err := commentsAPI.Retrieve(comment1.ID).ContextView().Do()
			require.NoError(t, err)
			assert.Equal(t, comment1.ID, viewComment.ID)
			assert.Empty(t, viewComment.AuthorEmail, "view context should not expose author email")

			// Context edit
			editComment, err := commentsAPI.Retrieve(comment1.ID).ContextEdit().Do()
			require.NoError(t, err)
			assert.Equal(t, comment1.ID, editComment.ID)
			assert.NotEmpty(t, editComment.AuthorEmail, "edit context should expose author email")
			assert.NotEmpty(t, editComment.Content.Raw, "edit context should expose content.raw")
		})

		t.Run("retrieve_with_fields_and_embed", func(t *testing.T) {
			fieldsComment, err := commentsAPI.Retrieve(comment1.ID).
				Fields("id", "post", "author_name").
				Embed().
				Do()
			require.NoError(t, err)
			assert.Equal(t, comment1.ID, fieldsComment.ID)
			assert.Equal(t, post.ID, fieldsComment.Post)
			assert.Equal(t, "Author One", fieldsComment.AuthorName)
			assert.Nil(t, fieldsComment.Content)
		})

		t.Run("retrieve_not_found", func(t *testing.T) {
			_, err := commentsAPI.Retrieve(999999999).Do()
			assert.Error(t, err)
			var wpErr *gowprest.WPRestError
			if assert.ErrorAs(t, err, &wpErr) {
				assert.Equal(t, 404, wpErr.Data.Status)
			}
		})
	})

	t.Run("UpdateComment", func(t *testing.T) {
		t.Run("update_with_struct", func(t *testing.T) {
			updatedContent := "Updated via struct: " + faker.Paragraph()
			c, err := commentsAPI.Update(gowprest.CommentData{
				ID:      comment1.ID,
				Content: updatedContent,
			}).Do()
			require.NoError(t, err)
			assert.Equal(t, comment1.ID, c.ID)
			assert.Contains(t, c.Content.Rendered, updatedContent)
		})

		t.Run("update_with_method_chaining", func(t *testing.T) {
			updatedContent := "Updated via method chaining: " + faker.Paragraph()
			c, err := commentsAPI.Update().
				ID(comment2.ID).
				Content(updatedContent).
				AuthorName("Author Two Updated").
				StatusHold().
				Do()
			require.NoError(t, err)
			assert.Equal(t, comment2.ID, c.ID)
			assert.Equal(t, "Author Two Updated", c.AuthorName)
			assert.Contains(t, c.Content.Rendered, updatedContent)
			assert.Equal(t, string(gowprest.CommentStatusHold), c.Status)
		})

		t.Run("update_not_found", func(t *testing.T) {
			_, err := commentsAPI.Update().
				ID(999999999).
				Content("Does not exist").
				Do()
			assert.Error(t, err)
		})
	})

	t.Run("DeleteComment", func(t *testing.T) {
		t.Run("trash_comment", func(t *testing.T) {
			// Soft delete (move to trash)
			deleted, err := commentsAPI.Delete(replyComment.ID).Do()
			require.NoError(t, err)
			assert.Equal(t, replyComment.ID, deleted.ID)
			assert.Equal(t, string(gowprest.CommentStatusTrash), deleted.Status)
		})

		t.Run("force_delete_success", func(t *testing.T) {
			// Permanently delete comment2
			deleted, err := commentsAPI.Delete(comment2.ID).Force().Do()
			require.NoError(t, err)
			assert.True(t, deleted.Deleted)
			assert.Equal(t, comment2.ID, deleted.ID)

			// Verify comment no longer retrievable
			_, err = commentsAPI.Retrieve(comment2.ID).Do()
			assert.Error(t, err)
		})

		t.Run("delete_not_found", func(t *testing.T) {
			_, err := commentsAPI.Delete(999999999).Force().Do()
			assert.Error(t, err)
			var wpErr *gowprest.WPRestError
			if assert.ErrorAs(t, err, &wpErr) {
				assert.Equal(t, 404, wpErr.Data.Status)
			}
		})
	})
}
