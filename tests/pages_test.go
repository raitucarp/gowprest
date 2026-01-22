package tests

import (
	"os"
	"testing"

	"github.com/go-faker/faker/v4"
	_ "github.com/joho/godotenv/autoload"
	"github.com/raitucarp/gowprest"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCreatePage(t *testing.T) {
	client := gowprest.NewClient(blogUrl).
		WithBasicAuth(
			os.Getenv("BLOG_USERNAME"),
			os.Getenv("BLOG_APP_PASSWORD"),
		)

	defer client.Close()

	pageAPI := client.Pages()

	listAPI := pageAPI.List()
	pages, err := listAPI.Do()

	pageCountBefore := len(pages)

	assert.Equal(t, nil, err)

	_, err = pageAPI.Create(gowprest.PageData{
		Title:   faker.Sentence(),
		Content: faker.Paragraph(),
		Excerpt: faker.Sentence(),
		Status:  gowprest.StatusPublished,
	}).Do()

	assert.Equal(t, nil, err, err)
	pages, err = listAPI.Do()

	pageCountAfter := len(pages)

	assert.LessOrEqual(t, pageCountBefore, pageCountAfter,
		"New page count should greater than or equal to the previous count. Expected %d, got %d",
		pageCountBefore, pageCountAfter)
}

func TestListPages(t *testing.T) {
	client := gowprest.NewClient(blogUrl)
	defer client.Close()

	t.Run("not found search", func(t *testing.T) {
		listPage := client.Pages().List()
		listPage = listPage.Search("non-existent-page-search-query-12345")
		pages, err := listPage.Do()

		assert.Equal(t, nil, err)
		assert.Equal(t, len(pages), 0, "Page should not return anything")
	})

	t.Run("default page list", func(t *testing.T) {
		listPage := client.Pages().List()
		pages, err := listPage.Do()

		assert.Equal(t, nil, err)
		assert.GreaterOrEqual(t, len(pages), 0, "Page list should be returned")
	})
}

func TestRetrievePage(t *testing.T) {
	client := gowprest.NewClient(blogUrl)
	defer client.Close()

	listAPI := client.Pages().List()
	pages, err := listAPI.Do()
	assert.Equal(t, nil, err)

	if len(pages) == 0 {
		t.Skip("No pages found to test Retrieve")
	}

	singlePage, err := client.Pages().Retrieve(pages[0].ID).Do()
	assert.Equal(t, nil, err)
	assert.Equal(t, pages[0].Title.Rendered, singlePage.Title.Rendered)
}

func TestUpdatePage(t *testing.T) {
	client := gowprest.NewClient(blogUrl).
		WithBasicAuth(
			os.Getenv("BLOG_USERNAME"),
			os.Getenv("BLOG_APP_PASSWORD"),
		)
	defer client.Close()

	listAPI := client.Pages().List()
	pages, err := listAPI.Do()
	assert.Equal(t, nil, err)

	if len(pages) == 0 {
		t.Skip("No pages found to test Update")
	}

	selectedPage := pages[0]
	pageID := selectedPage.ID
	updateAPI := client.Pages().Update(gowprest.PageData{
		ID:      pageID,
		Title:   faker.Sentence(),
		Content: faker.Paragraph(),
		Excerpt: faker.Sentence(),
		Status:  gowprest.StatusPublished,
	})

	updatedPage, err := updateAPI.Do()
	assert.Equal(t, nil, err)
	assert.Equal(t, pageID, updatedPage.ID)
	assert.NotEqual(t, selectedPage.Title.Rendered, updatedPage.Title.Rendered)

	updateAPI = client.Pages().Update(gowprest.PageData{
		ID:     pageID,
		Status: gowprest.StatusDraft,
	})

	updatedPage, err = updateAPI.Do()
	assert.Equal(t, nil, err)
	assert.Equal(t, pageID, updatedPage.ID)
	assert.NotEqual(t, selectedPage.Status, updatedPage.Status)
}

func TestDeletePage(t *testing.T) {
	client := gowprest.NewClient(blogUrl).
		WithBasicAuth(
			os.Getenv("BLOG_USERNAME"),
			os.Getenv("BLOG_APP_PASSWORD"),
		)
	defer client.Close()

	createdPage, err := client.Pages().Create(gowprest.PageData{
		Title:   faker.Sentence(),
		Content: faker.Paragraph(),
		Excerpt: faker.Sentence(),
		Status:  gowprest.StatusPublished,
	}).Do()

	require.Equal(t, nil, err)

	pageID := createdPage.ID
	deleteAPI := client.Pages().Delete(pageID)

	trashedPage, err := deleteAPI.Do()
	require.Equal(t, nil, err)

	singlePage, err := client.Pages().Retrieve(trashedPage.ID).
		ContextEdit().
		Do()

	require.Equal(t, nil, err)
	require.Equal(t, trashedPage.ID, singlePage.ID)

	_, err = client.Pages().Delete(pageID).Force().Do()
	require.Equal(t, nil, err)

	_, err = client.Pages().Retrieve(pageID).Do()
	assert.NotEqual(t, nil, err)
}
