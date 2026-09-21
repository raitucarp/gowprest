package tests

import (
	"errors"
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/go-faker/faker/v4"
	_ "github.com/joho/godotenv/autoload"
	"github.com/raitucarp/gowprest"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestBlocks(t *testing.T) {
	blogUrl := os.Getenv("BLOG_URL")
	username := os.Getenv("BLOG_USERNAME")
	password := os.Getenv("BLOG_APP_PASSWORD")

	if blogUrl == "" || username == "" || password == "" {
		t.Skip("Skipping test; BLOG_URL, BLOG_USERNAME, or BLOG_APP_PASSWORD not set")
	}

	client := gowprest.NewClient(blogUrl).WithBasicAuth(username, password)
	defer client.Close()

	t.Run("BlocksAliases", func(t *testing.T) {
		assert.NotNil(t, client.Blocks())
		assert.NotNil(t, client.EditorBlocks())
		assert.NotNil(t, client.Blocks().List())
	})

	t.Run("CreateBlock_ChainedSetters", func(t *testing.T) {
		unique := fmt.Sprintf("Blk%d", time.Now().UnixNano()%1000000)
		title := "Test Reusable Block " + unique
		content := "<!-- wp:paragraph --><p>" + faker.Sentence() + "</p><!-- /wp:paragraph -->"

		block, err := client.Blocks().Create().
			Title(title).
			Content(content).
			StatusPublish().
			Do()
		require.NoError(t, err)
		require.NotNil(t, block)
		assert.NotZero(t, block.ID)
		assert.Equal(t, gowprest.StatusPublished, block.Status)
		assert.Equal(t, "wp_block", block.Type)

		assert.Contains(t, block.TitleString(), title)

		// Cleanup permanently
		defer client.Blocks().Delete(block.ID).Force().Do()
	})

	t.Run("CreateBlock_StructPayload", func(t *testing.T) {
		unique := fmt.Sprintf("BlkStruct%d", time.Now().UnixNano()%1000000)
		title := "Struct Block " + unique

		block, err := client.Blocks().Create(gowprest.BlockData{
			Title:   title,
			Content: "<!-- wp:paragraph --><p>Content via struct</p><!-- /wp:paragraph -->",
			Status:  gowprest.StatusDraft,
		}).Do()
		require.NoError(t, err)
		require.NotNil(t, block)
		assert.NotZero(t, block.ID)
		assert.Equal(t, gowprest.StatusDraft, block.Status)

		defer client.Blocks().Delete(block.ID).Force().Do()
	})

	t.Run("ListBlocks", func(t *testing.T) {
		// Create 2 test blocks
		unique := fmt.Sprintf("List%d", time.Now().UnixNano()%1000000)
		b1, err := client.Blocks().Create().
			Title("A List Block " + unique).
			Content("Content A").
			StatusPublish().
			Do()
		require.NoError(t, err)
		defer client.Blocks().Delete(b1.ID).Force().Do()

		b2, err := client.Blocks().Create().
			Title("Z List Block " + unique).
			Content("Content Z").
			StatusPublish().
			Do()
		require.NoError(t, err)
		defer client.Blocks().Delete(b2.ID).Force().Do()

		t.Run("Default", func(t *testing.T) {
			blocks, err := client.Blocks().List().Do()
			require.NoError(t, err)
			assert.NotEmpty(t, blocks)
		})

		t.Run("Pagination", func(t *testing.T) {
			page1, err := client.Blocks().List().
				Page(1).
				PerPage(1).
				Do()
			require.NoError(t, err)
			require.Len(t, page1, 1)

			page2, err := client.Blocks().List().
				Page(2).
				PerPage(1).
				Do()
			require.NoError(t, err)
			require.Len(t, page2, 1)

			assert.NotEqual(t, page1[0].ID, page2[0].ID)
		})

		t.Run("Search", func(t *testing.T) {
			results, err := client.Blocks().List().
				Search(unique).
				Do()
			require.NoError(t, err)
			assert.GreaterOrEqual(t, len(results), 2)
		})

		t.Run("Ordering", func(t *testing.T) {
			asc, err := client.Blocks().List().
				Search(unique).
				OrderByTitle().
				OrderAsc().
				Do()
			require.NoError(t, err)
			require.GreaterOrEqual(t, len(asc), 2)

			desc, err := client.Blocks().List().
				Search(unique).
				OrderByTitle().
				OrderDesc().
				Do()
			require.NoError(t, err)
			require.GreaterOrEqual(t, len(desc), 2)

			assert.NotEqual(t, asc[0].ID, desc[0].ID)
		})

		t.Run("IncludeAndExclude", func(t *testing.T) {
			inc, err := client.Blocks().List().
				Include(b1.ID).
				Do()
			require.NoError(t, err)
			require.Len(t, inc, 1)
			assert.Equal(t, b1.ID, inc[0].ID)

			exc, err := client.Blocks().List().
				Search(unique).
				Exclude(b1.ID).
				Do()
			require.NoError(t, err)
			for _, b := range exc {
				assert.NotEqual(t, b1.ID, b.ID)
			}
		})

		t.Run("Contexts", func(t *testing.T) {
			viewList, err := client.Blocks().List().ContextView().PerPage(2).Do()
			require.NoError(t, err)
			assert.NotEmpty(t, viewList)

			editList, err := client.Blocks().List().ContextEdit().PerPage(2).Do()
			require.NoError(t, err)
			assert.NotEmpty(t, editList)

			embedList, err := client.Blocks().List().ContextEmbed().PerPage(2).Do()
			require.NoError(t, err)
			assert.NotEmpty(t, embedList)
		})

		t.Run("Fields", func(t *testing.T) {
			fieldsList, err := client.Blocks().List().
				Search(unique).
				Fields("id", "slug").
				Do()
			require.NoError(t, err)
			require.NotEmpty(t, fieldsList)

			for _, b := range fieldsList {
				assert.NotZero(t, b.ID)
				assert.NotEmpty(t, b.Slug)
				assert.Nil(t, b.Title)
				assert.Nil(t, b.Content)
			}
		})
	})

	t.Run("RetrieveBlock", func(t *testing.T) {
		unique := fmt.Sprintf("Ret%d", time.Now().UnixNano()%1000000)
		created, err := client.Blocks().Create().
			Title("Retrieve Test " + unique).
			Content("<p>Content</p>").
			StatusPublish().
			Do()
		require.NoError(t, err)
		defer client.Blocks().Delete(created.ID).Force().Do()

		t.Run("Retrieve_Success", func(t *testing.T) {
			retrieved, err := client.Blocks().Retrieve(created.ID).Do()
			require.NoError(t, err)
			require.NotNil(t, retrieved)
			assert.Equal(t, created.ID, retrieved.ID)
			assert.Equal(t, "wp_block", retrieved.Type)
			assert.Contains(t, retrieved.TitleString(), unique)
		})

		t.Run("Retrieve_ContextAndFields", func(t *testing.T) {
			retrieved, err := client.Blocks().Retrieve(created.ID).
				ContextEdit().
				Fields("id", "title").
				Do()
			require.NoError(t, err)
			require.NotNil(t, retrieved)
			assert.Equal(t, created.ID, retrieved.ID)
			assert.NotNil(t, retrieved.Title)
			assert.Nil(t, retrieved.Content)
		})

		t.Run("Retrieve_NotFound", func(t *testing.T) {
			_, err := client.Blocks().Retrieve(9999999).Do()
			require.Error(t, err)

			var wpErr *gowprest.WPRestError
			if assert.True(t, errors.As(err, &wpErr)) {
				assert.Equal(t, "rest_post_invalid_id", wpErr.Code)
				assert.Equal(t, 404, wpErr.Data.Status)
			}
		})
	})

	t.Run("UpdateBlock", func(t *testing.T) {
		unique := fmt.Sprintf("Upd%d", time.Now().UnixNano()%1000000)
		created, err := client.Blocks().Create().
			Title("Original Title " + unique).
			Content("<p>Original Content</p>").
			StatusPublish().
			Do()
		require.NoError(t, err)
		defer client.Blocks().Delete(created.ID).Force().Do()

		updatedTitle := "Updated Title " + unique
		updatedContent := "<p>Updated Content</p>"

		updated, err := client.Blocks().Update(created.ID).
			Title(updatedTitle).
			Content(updatedContent).
			StatusDraft().
			Do()
		require.NoError(t, err)
		require.NotNil(t, updated)
		assert.Equal(t, created.ID, updated.ID)
		assert.Equal(t, gowprest.StatusDraft, updated.Status)
		assert.Contains(t, updated.TitleString(), updatedTitle)
	})

	t.Run("DeleteBlock", func(t *testing.T) {
		t.Run("DeleteToTrash", func(t *testing.T) {
			unique := fmt.Sprintf("Trash%d", time.Now().UnixNano()%1000000)
			created, err := client.Blocks().Create().
				Title("Trash Block " + unique).
				Content("<p>To be trashed</p>").
				StatusPublish().
				Do()
			require.NoError(t, err)

			trashed, err := client.Blocks().Delete(created.ID).Do()
			require.NoError(t, err)
			require.NotNil(t, trashed)
			assert.Equal(t, created.ID, trashed.ID)
			assert.Equal(t, gowprest.PostStatus("trash"), trashed.Status)

			// Clean up permanently
			defer client.Blocks().Delete(created.ID).Force().Do()
		})

		t.Run("ForceDelete", func(t *testing.T) {
			unique := fmt.Sprintf("ForceDel%d", time.Now().UnixNano()%1000000)
			created, err := client.Blocks().Create().
				Title("Force Delete Block " + unique).
				Content("<p>To be permanently deleted</p>").
				StatusPublish().
				Do()
			require.NoError(t, err)

			deleted, err := client.Blocks().Delete(created.ID).Force().Do()
			require.NoError(t, err)
			require.NotNil(t, deleted)
			assert.Equal(t, created.ID, deleted.ID)

			// Verify it is gone
			_, err = client.Blocks().Retrieve(created.ID).Do()
			require.Error(t, err)
		})
	})
}
