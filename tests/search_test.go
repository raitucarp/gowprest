package tests

import (
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

func TestSearch(t *testing.T) {
	blogUrl := os.Getenv("BLOG_URL")
	username := os.Getenv("BLOG_USERNAME")
	password := os.Getenv("BLOG_APP_PASSWORD")

	if blogUrl == "" || username == "" || password == "" {
		t.Skip("Skipping test; BLOG_URL, BLOG_USERNAME, or BLOG_APP_PASSWORD not set")
	}

	client := gowprest.NewClient(blogUrl).WithBasicAuth(username, password)
	defer client.Close()

	t.Run("SearchAliases", func(t *testing.T) {
		assert.NotNil(t, client.Search())
		assert.NotNil(t, client.SearchResults())
		assert.NotNil(t, client.Search().List())
	})

	t.Run("Search_Default", func(t *testing.T) {
		results, err := client.Search().
			PerPage(5).
			Do()
		require.NoError(t, err)
		assert.NotEmpty(t, results)

		first := results[0]
		assert.NotZero(t, first.IntID())
		assert.NotEmpty(t, first.StringID())
		assert.NotEmpty(t, first.Title)
		assert.NotEmpty(t, first.URL)
		assert.NotEmpty(t, first.Type)
		assert.NotEmpty(t, first.Subtype)
	})

	t.Run("Search_ByQuery", func(t *testing.T) {
		uniqueWord := fmt.Sprintf("SrchKw%d", time.Now().UnixNano()%1000000)
		postTitle := "Search Test Post " + uniqueWord

		// Create a test post to search for
		post, err := client.Posts().Create().
			Title(postTitle).
			Content("<p>Search content body " + faker.Paragraph() + "</p>").
			StatusPublish().
			Do()
		require.NoError(t, err)
		defer client.Posts().Delete(post.ID).Force().Do()

		// Search for the unique word
		results, err := client.Search().
			Query(uniqueWord).
			Do()
		require.NoError(t, err)
		require.NotEmpty(t, results)

		found := false
		for _, r := range results {
			if r.IntID() == post.ID {
				found = true
				assert.Contains(t, r.Title, uniqueWord)
				assert.Equal(t, "post", r.Type)
				assert.Equal(t, "post", r.Subtype)
				break
			}
		}
		assert.True(t, found, "expected newly created post to be found in search results")
	})

	t.Run("Search_Pagination", func(t *testing.T) {
		page1, err := client.Search().
			Page(1).
			PerPage(1).
			Do()
		require.NoError(t, err)
		require.Len(t, page1, 1)

		page2, err := client.Search().
			Page(2).
			PerPage(1).
			Do()
		require.NoError(t, err)
		require.Len(t, page2, 1)

		assert.NotEqual(t, page1[0].IntID(), page2[0].IntID())
	})

	t.Run("Search_ByType", func(t *testing.T) {
		t.Run("type_post", func(t *testing.T) {
			results, err := client.Search().
				TypePost().
				PerPage(5).
				Do()
			require.NoError(t, err)
			assert.NotEmpty(t, results)
			for _, r := range results {
				assert.Equal(t, gowprest.SearchTypePost, r.Type)
			}
		})

		t.Run("type_post_format", func(t *testing.T) {
			results, err := client.Search().
				TypePostFormat().
				Do()
			require.NoError(t, err)
			assert.NotEmpty(t, results)
			for _, r := range results {
				assert.Equal(t, gowprest.SearchTypePostFormat, r.Type)
				assert.NotEmpty(t, r.StringID())
			}
		})
	})

	t.Run("Search_BySubtype", func(t *testing.T) {
		t.Run("subtype_post", func(t *testing.T) {
			results, err := client.Search().
				SubtypePost().
				PerPage(5).
				Do()
			require.NoError(t, err)
			assert.NotEmpty(t, results)
			for _, r := range results {
				assert.Equal(t, gowprest.SearchSubtypePost, r.Subtype)
			}
		})

		t.Run("subtype_page", func(t *testing.T) {
			results, err := client.Search().
				SubtypePage().
				PerPage(5).
				Do()
			require.NoError(t, err)
			assert.NotEmpty(t, results)
			for _, r := range results {
				assert.Equal(t, gowprest.SearchSubtypePage, r.Subtype)
			}
		})
	})

	t.Run("Search_IncludeAndExclude", func(t *testing.T) {
		initial, err := client.Search().
			TypePost().
			PerPage(2).
			Do()
		require.NoError(t, err)
		require.True(t, len(initial) >= 2)

		firstID := initial[0].IntID()

		// Exclude first item
		excluded, err := client.Search().
			TypePost().
			PerPage(5).
			Exclude(firstID).
			Do()
		require.NoError(t, err)
		for _, r := range excluded {
			assert.NotEqual(t, firstID, r.IntID())
		}

		// Include only first item
		included, err := client.Search().
			Include(firstID).
			Do()
		require.NoError(t, err)
		require.Len(t, included, 1)
		assert.Equal(t, firstID, included[0].IntID())
	})

	t.Run("Search_ContextAndGlobalParams", func(t *testing.T) {
		viewResults, err := client.Search().
			ContextView().
			PerPage(2).
			Do()
		require.NoError(t, err)
		assert.NotEmpty(t, viewResults)

		embedResults, err := client.Search().
			ContextEmbed().
			PerPage(2).
			Do()
		require.NoError(t, err)
		assert.NotEmpty(t, embedResults)

		fieldsResults, err := client.Search().
			Fields("id", "title").
			PerPage(2).
			Embed().
			Do()
		require.NoError(t, err)
		assert.NotEmpty(t, fieldsResults)
		assert.NotEmpty(t, fieldsResults[0].Title)
		assert.NotZero(t, fieldsResults[0].IntID())
		// URL should be omitted
		assert.Empty(t, fieldsResults[0].URL)
	})
}
