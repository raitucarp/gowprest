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

func TestCategoriesAPI(t *testing.T) {
	client := gowprest.NewClient(blogUrl).
		WithBasicAuth(
			os.Getenv("BLOG_USERNAME"),
			os.Getenv("BLOG_APP_PASSWORD"),
		)
	defer client.Close()

	categoryName := faker.Word()
	categoryDescription := faker.Sentence()

	// 1. Create Category
	newCategory, err := client.Categories().Create(gowprest.CategoryData{
		Name:        categoryName,
		Description: categoryDescription,
	}).Do()

	require.Nil(t, err)
	assert.Equal(t, categoryName, newCategory.Name)
	assert.Equal(t, categoryDescription, newCategory.Description)
	categoryID := newCategory.ID

	// 2. List Categories
	categories, err := client.Categories().List().Do()
	require.Nil(t, err)
	assert.GreaterOrEqual(t, len(categories), 1)

	found := false
	for _, cat := range categories {
		if cat.ID == categoryID {
			found = true
			break
		}
	}
	assert.True(t, found, "Created category should be in the list")

	// 3. Retrieve Category
	retrievedCategory, err := client.Categories().Retrieve(categoryID).Do()
	require.Nil(t, err)
	assert.Equal(t, categoryID, retrievedCategory.ID)
	assert.Equal(t, categoryName, retrievedCategory.Name)

	// 4. Update Category
	updatedName := categoryName + " (Updated)"
	updatedCategory, err := client.Categories().Update(gowprest.CategoryData{
		ID:   categoryID,
		Name: updatedName,
	}).Do()

	require.Nil(t, err)
	assert.Equal(t, categoryID, updatedCategory.ID)
	assert.Equal(t, updatedName, updatedCategory.Name)

	// 5. Delete Category
	deletedCategory, err := client.Categories().Delete(categoryID).Force().Do()
	require.Nil(t, err)
	if deletedCategory.ID == 0 {
		t.Logf("Deleted category ID is 0. Full object: %+v", deletedCategory)
	}
	assert.Equal(t, categoryID, deletedCategory.ID)

	// 6. Verify Deletion
	_, err = client.Categories().Retrieve(categoryID).Do()
	assert.NotNil(t, err, "Category should not be found after deletion")
}
