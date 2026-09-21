package tests

import (
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/go-faker/faker/v4"
	_ "github.com/joho/godotenv/autoload"
	"github.com/raitucarp/gowprest"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// sample 1x1 transparent PNG
var samplePNG = []byte{
	0x89, 0x50, 0x4E, 0x47, 0x0D, 0x0A, 0x1A, 0x0A, 0x00, 0x00, 0x00, 0x0D,
	0x49, 0x48, 0x44, 0x52, 0x00, 0x00, 0x00, 0x01, 0x00, 0x00, 0x00, 0x01,
	0x08, 0x06, 0x00, 0x00, 0x00, 0x1F, 0x15, 0xC4, 0x89, 0x00, 0x00, 0x00,
	0x0A, 0x49, 0x44, 0x41, 0x54, 0x78, 0x9C, 0x63, 0x00, 0x01, 0x00, 0x00,
	0x05, 0x00, 0x01, 0x0D, 0x0A, 0x2D, 0xB4, 0x00, 0x00, 0x00, 0x00, 0x49,
	0x45, 0x4E, 0x44, 0xAE, 0x42, 0x60, 0x82,
}

func TestMedia(t *testing.T) {
	blogUrl := os.Getenv("BLOG_URL")
	username := os.Getenv("BLOG_USERNAME")
	password := os.Getenv("BLOG_APP_PASSWORD")

	if blogUrl == "" || username == "" || password == "" {
		t.Skip("Skipping test; BLOG_URL, BLOG_USERNAME, or BLOG_APP_PASSWORD not set")
	}

	client := gowprest.NewClient(blogUrl).WithBasicAuth(username, password)
	defer client.Close()

	mediaAPI := client.Media()
	postsAPI := client.Posts()

	// Setup a post to test media attachment
	post, err := postsAPI.Create().
		Title("Media Test Post " + faker.Sentence()).
		Content("<p>" + faker.Paragraph() + "</p>").
		StatusPublish().
		Do()
	require.NoError(t, err)
	defer postsAPI.Delete(post.ID).Force().Do()

	var media1 gowprest.Media
	var media2 gowprest.Media
	var mediaAttached gowprest.Media
	uniqueKeyword := "MedKw" + faker.Word()

	t.Run("CreateMedia", func(t *testing.T) {
		t.Run("upload_from_bytes", func(t *testing.T) {
			m, err := mediaAPI.Create().
				FileBytes("sample_image.png", samplePNG, "image/png").
				Title("Test Image 1 " + uniqueKeyword).
				AltText("Alt text for image 1").
				Caption("Caption for image 1").
				Description("Description for image 1").
				Do()
			require.NoError(t, err)
			assert.NotZero(t, m.ID)
			assert.Contains(t, m.Title.Rendered, "Test Image 1")
			assert.Equal(t, "Alt text for image 1", m.AltText)
			assert.Equal(t, "image/png", m.MimeType)
			assert.Equal(t, "image", m.MediaType)
			assert.NotEmpty(t, m.SourceURL)
			media1 = m
		})

		t.Run("upload_from_file_path", func(t *testing.T) {
			tempDir := t.TempDir()
			tempFile := filepath.Join(tempDir, "temp_upload.png")
			err := os.WriteFile(tempFile, samplePNG, 0644)
			require.NoError(t, err)

			m, err := mediaAPI.Create().
				File(tempFile).
				Title("Test Image 2 " + faker.Sentence()).
				AltText("Alt text for image 2").
				Do()
			require.NoError(t, err)
			assert.NotZero(t, m.ID)
			assert.Equal(t, "image/png", m.MimeType)
			media2 = m
		})

		t.Run("upload_attached_to_post", func(t *testing.T) {
			m, err := mediaAPI.Create().
				FileBytes("attached_image.png", samplePNG, "image/png").
				Title("Attached Image").
				Post(post.ID).
				Do()
			require.NoError(t, err)
			assert.NotZero(t, m.ID)
			require.NotNil(t, m.Post)
			assert.Equal(t, post.ID, *m.Post)
			mediaAttached = m
		})
	})

	// Ensure cleanup of created media items
	defer func() {
		if media1.ID != 0 {
			mediaAPI.Delete(media1.ID).Force().Do()
		}
		if media2.ID != 0 {
			mediaAPI.Delete(media2.ID).Force().Do()
		}
		if mediaAttached.ID != 0 {
			mediaAPI.Delete(mediaAttached.ID).Force().Do()
		}
	}()

	t.Run("ListMedia", func(t *testing.T) {
		t.Run("default_list", func(t *testing.T) {
			items, err := mediaAPI.List().Do()
			require.NoError(t, err)
			assert.GreaterOrEqual(t, len(items), 2)
		})

		t.Run("pagination_page_and_per_page", func(t *testing.T) {
			page1, err := mediaAPI.List().Page(1).PerPage(1).Do()
			require.NoError(t, err)
			assert.Len(t, page1, 1)

			page2, err := mediaAPI.List().Page(2).PerPage(1).Do()
			require.NoError(t, err)
			assert.Len(t, page2, 1)
			assert.NotEqual(t, page1[0].ID, page2[0].ID)

			// Offset
			offsetList, err := mediaAPI.List().Offset(1).PerPage(1).Do()
			require.NoError(t, err)
			assert.Len(t, offsetList, 1)
			assert.Equal(t, page2[0].ID, offsetList[0].ID)
		})

		t.Run("filter_by_date", func(t *testing.T) {
			afterDate := time.Now().Add(-24 * time.Hour)
			beforeDate := time.Now().Add(24 * time.Hour)
			dateList, err := mediaAPI.List().
				After(afterDate).
				Before(beforeDate).
				Do()
			require.NoError(t, err)
			assert.NotEmpty(t, dateList)
		})

		t.Run("ordering_asc_and_desc", func(t *testing.T) {
			ascList, err := mediaAPI.List().OrderAsc().OrderByID().Do()
			require.NoError(t, err)
			require.GreaterOrEqual(t, len(ascList), 2)

			descList, err := mediaAPI.List().OrderDesc().OrderBy("id").Do()
			require.NoError(t, err)
			require.GreaterOrEqual(t, len(descList), 2)

			assert.NotEqual(t, ascList[0].ID, descList[0].ID)
			assert.True(t, ascList[0].ID < ascList[1].ID, "ascList should be sorted ascending by ID")
			assert.True(t, descList[0].ID > descList[1].ID, "descList should be sorted descending by ID")

			// OrderByDate
			dateList, err := mediaAPI.List().OrderDesc().OrderByDate().Do()
			require.NoError(t, err)
			assert.NotEmpty(t, dateList)
		})

		t.Run("filter_by_parent_and_media_type", func(t *testing.T) {
			// By parent
			attached, err := mediaAPI.List().Parent(post.ID).Do()
			require.NoError(t, err)
			assert.NotEmpty(t, attached)
			for _, item := range attached {
				require.NotNil(t, item.Post)
				assert.Equal(t, post.ID, *item.Post)
			}

			// By media_type
			images, err := mediaAPI.List().MediaTypeImage().Do()
			require.NoError(t, err)
			assert.NotEmpty(t, images)
			for _, item := range images {
				assert.Equal(t, "image", item.MediaType)
			}

			// By mime_type
			pngs, err := mediaAPI.List().MimeType("image/png").Do()
			require.NoError(t, err)
			assert.NotEmpty(t, pngs)
			for _, item := range pngs {
				assert.Equal(t, "image/png", item.MimeType)
			}
		})

		t.Run("filter_by_include_and_exclude", func(t *testing.T) {
			// Include
			included, err := mediaAPI.List().Include(media1.ID, media2.ID).Do()
			require.NoError(t, err)
			assert.Len(t, included, 2)

			// Exclude
			excluded, err := mediaAPI.List().Exclude(media1.ID).Do()
			require.NoError(t, err)
			for _, item := range excluded {
				assert.NotEqual(t, media1.ID, item.ID)
			}
		})

		t.Run("search_query", func(t *testing.T) {
			searchRes, err := mediaAPI.List().Search(uniqueKeyword).Do()
			require.NoError(t, err)
			assert.NotEmpty(t, searchRes)
			assert.Equal(t, media1.ID, searchRes[0].ID)
		})

		t.Run("context_and_global_parameters", func(t *testing.T) {
			// Context edit
			editList, err := mediaAPI.List().ContextEdit().Do()
			require.NoError(t, err)
			assert.NotEmpty(t, editList)

			// Context embed
			embedList, err := mediaAPI.List().ContextEmbed().Embed().Do()
			require.NoError(t, err)
			assert.NotEmpty(t, embedList)

			// Fields
			fieldsList, err := mediaAPI.List().Fields("id", "source_url").Do()
			require.NoError(t, err)
			assert.NotEmpty(t, fieldsList)
			assert.NotZero(t, fieldsList[0].ID)
			assert.NotEmpty(t, fieldsList[0].SourceURL)
			assert.Nil(t, fieldsList[0].Title)
		})
	})

	t.Run("RetrieveMedia", func(t *testing.T) {
		t.Run("retrieve_existing_media", func(t *testing.T) {
			m, err := mediaAPI.Retrieve(media1.ID).Do()
			require.NoError(t, err)
			assert.Equal(t, media1.ID, m.ID)
			assert.Equal(t, "image/png", m.MimeType)
			assert.NotEmpty(t, m.SourceURL)
		})

		t.Run("retrieve_with_contexts", func(t *testing.T) {
			viewMedia, err := mediaAPI.Retrieve(media1.ID).ContextView().Do()
			require.NoError(t, err)
			assert.Equal(t, media1.ID, viewMedia.ID)

			editMedia, err := mediaAPI.Retrieve(media1.ID).ContextEdit().Do()
			require.NoError(t, err)
			assert.Equal(t, media1.ID, editMedia.ID)
			if editMedia.Title != nil {
				assert.NotEmpty(t, editMedia.Title.Raw)
			}
		})

		t.Run("retrieve_with_fields_and_embed", func(t *testing.T) {
			fieldsMedia, err := mediaAPI.Retrieve(media1.ID).
				Fields("id", "source_url").
				Embed().
				Do()
			require.NoError(t, err)
			assert.Equal(t, media1.ID, fieldsMedia.ID)
			assert.NotEmpty(t, fieldsMedia.SourceURL)
			assert.Nil(t, fieldsMedia.Title)
		})

		t.Run("retrieve_not_found", func(t *testing.T) {
			_, err := mediaAPI.Retrieve(999999999).Do()
			assert.Error(t, err)
			var wpErr *gowprest.WPRestError
			if assert.ErrorAs(t, err, &wpErr) {
				assert.Equal(t, 404, wpErr.Data.Status)
			}
		})
	})

	t.Run("UpdateMedia", func(t *testing.T) {
		t.Run("update_metadata", func(t *testing.T) {
			updatedTitle := "Updated Title " + faker.Sentence()
			updatedAlt := "Updated Alt Text"
			m, err := mediaAPI.Update().
				ID(media1.ID).
				Title(updatedTitle).
				AltText(updatedAlt).
				Do()
			require.NoError(t, err)
			assert.Equal(t, media1.ID, m.ID)
			assert.Contains(t, m.Title.Rendered, updatedTitle)
			assert.Equal(t, updatedAlt, m.AltText)
		})

		t.Run("update_with_struct", func(t *testing.T) {
			updatedCaption := "Updated caption via struct: " + faker.Sentence()
			m, err := mediaAPI.Update(gowprest.MediaData{
				ID:      media1.ID,
				Caption: updatedCaption,
			}).Do()
			require.NoError(t, err)
			assert.Equal(t, media1.ID, m.ID)
		})

		t.Run("update_not_found", func(t *testing.T) {
			_, err := mediaAPI.Update().
				ID(999999999).
				Title("Non existent").
				Do()
			assert.Error(t, err)
		})
	})

	t.Run("DeleteMedia", func(t *testing.T) {
		t.Run("force_delete_success", func(t *testing.T) {
			deleted, err := mediaAPI.Delete(media2.ID).Force().Do()
			require.NoError(t, err)
			assert.True(t, deleted.Deleted)
			assert.Equal(t, media2.ID, deleted.ID)

			// Verify media is no longer retrievable
			_, err = mediaAPI.Retrieve(media2.ID).Do()
			assert.Error(t, err)
			media2.ID = 0 // prevent double delete in defer
		})

		t.Run("delete_not_found", func(t *testing.T) {
			_, err := mediaAPI.Delete(999999999).Force().Do()
			assert.Error(t, err)
			var wpErr *gowprest.WPRestError
			if assert.ErrorAs(t, err, &wpErr) {
				assert.Equal(t, 404, wpErr.Data.Status)
			}
		})
	})
}
