package tests

import (
	"errors"
	"os"
	"testing"

	_ "github.com/joho/godotenv/autoload"
	"github.com/raitucarp/gowprest"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestApplicationPasswords(t *testing.T) {
	blogUrl := os.Getenv("BLOG_URL")
	username := os.Getenv("BLOG_USERNAME")
	password := os.Getenv("BLOG_APP_PASSWORD")

	if blogUrl == "" || username == "" || password == "" {
		t.Skip("Skipping test; BLOG_URL, BLOG_USERNAME, or BLOG_APP_PASSWORD not set")
	}

	client := gowprest.NewClient(blogUrl).WithBasicAuth(username, password)
	defer client.Close()

	t.Run("ApplicationPasswordsAliases", func(t *testing.T) {
		assert.NotNil(t, client.ApplicationPasswords())
		assert.NotNil(t, client.ApplicationPasswords("me"))
		assert.NotNil(t, client.ApplicationPasswords(1))
		assert.NotNil(t, client.Users().ApplicationPasswords())
		assert.NotNil(t, client.Users().ApplicationPasswords("me"))
		assert.NotNil(t, client.Users().ApplicationPasswords(1))
		assert.NotNil(t, client.Users().Me().ApplicationPasswords())
	})

	t.Run("Introspect", func(t *testing.T) {
		currentAppPass, err := client.ApplicationPasswords().Introspect().Do()
		require.NoError(t, err)
		require.NotNil(t, currentAppPass)
		assert.NotEmpty(t, currentAppPass.UUID)
		assert.NotEmpty(t, currentAppPass.Name)

		// Test Introspect with Context and Fields builders
		introspectWithFields, err := client.ApplicationPasswords().
			Introspect().
			ContextView().
			ContextEdit().
			ContextEmbed().
			Fields("uuid", "name").
			Do()
		require.NoError(t, err)
		require.NotNil(t, introspectWithFields)
		assert.Equal(t, currentAppPass.UUID, introspectWithFields.UUID)
		assert.Equal(t, currentAppPass.Name, introspectWithFields.Name)
	})

	t.Run("CRUD_Workflow", func(t *testing.T) {
		tempName := "temp-gowprest-test-key"
		tempAppID := "e16b92f7-7b89-4d64-8f4b-0123456789ab"

		// 1. Create
		created, err := client.ApplicationPasswords().
			Create(tempName).
			AppID(tempAppID).
			Do()
		require.NoError(t, err)
		require.NotNil(t, created)
		assert.NotEmpty(t, created.UUID)
		assert.Equal(t, tempName, created.Name)
		assert.Equal(t, tempAppID, created.AppID)
		assert.NotEmpty(t, created.Password, "Password should be returned in plain text on creation")

		tempUUID := created.UUID

		// 2. List
		list, err := client.ApplicationPasswords().
			List().
			ContextView().
			ContextEdit().
			ContextEmbed().
			Do()
		require.NoError(t, err)
		require.NotEmpty(t, list)

		var found bool
		for _, item := range list {
			if item.UUID == tempUUID {
				found = true
				assert.Equal(t, tempName, item.Name)
				assert.Equal(t, tempAppID, item.AppID)
			}
		}
		assert.True(t, found, "Newly created application password should be present in list")

		// List with Fields
		listFields, err := client.ApplicationPasswords().
			List().
			Fields("uuid", "name").
			Do()
		require.NoError(t, err)
		require.NotEmpty(t, listFields)
		for _, item := range listFields {
			assert.NotEmpty(t, item.UUID)
		}

		// 3. Retrieve
		retrieved, err := client.ApplicationPasswords().
			Retrieve(tempUUID).
			ContextView().
			ContextEdit().
			ContextEmbed().
			Do()
		require.NoError(t, err)
		require.NotNil(t, retrieved)
		assert.Equal(t, tempUUID, retrieved.UUID)
		assert.Equal(t, tempName, retrieved.Name)
		assert.Equal(t, tempAppID, retrieved.AppID)

		// Retrieve with Fields
		retrievedFields, err := client.ApplicationPasswords().
			Retrieve(tempUUID).
			Fields("uuid", "name").
			Do()
		require.NoError(t, err)
		require.NotNil(t, retrievedFields)
		assert.Equal(t, tempUUID, retrievedFields.UUID)
		assert.Equal(t, tempName, retrievedFields.Name)

		// 4. Update
		updatedName := "temp-gowprest-test-key-updated"
		updated, err := client.ApplicationPasswords().
			Update(tempUUID).
			Name(updatedName).
			Do()
		require.NoError(t, err)
		require.NotNil(t, updated)
		assert.Equal(t, tempUUID, updated.UUID)
		assert.Equal(t, updatedName, updated.Name)
		assert.Equal(t, tempAppID, updated.AppID)

		// 5. Delete
		deleted, err := client.ApplicationPasswords().
			Delete(tempUUID).
			Do()
		require.NoError(t, err)
		require.NotNil(t, deleted)
		assert.Equal(t, tempUUID, deleted.UUID)

		// Confirm deleted by attempting to retrieve
		_, err = client.ApplicationPasswords().Retrieve(tempUUID).Do()
		require.Error(t, err)

		var wpErr *gowprest.WPRestError
		if errors.As(err, &wpErr) {
			assert.Equal(t, "rest_application_password_not_found", wpErr.Code)
			assert.Equal(t, 404, wpErr.Data.Status)
		} else {
			t.Fatalf("expected WPRestError, got: %v", err)
		}
	})

	t.Run("NotFoundHandling", func(t *testing.T) {
		nonExistentUUID := "00000000-0000-0000-0000-000000000000"

		// Retrieve non-existent
		_, err := client.ApplicationPasswords().Retrieve(nonExistentUUID).Do()
		require.Error(t, err)
		var wpErr *gowprest.WPRestError
		if errors.As(err, &wpErr) {
			assert.Equal(t, "rest_application_password_not_found", wpErr.Code)
			assert.Equal(t, 404, wpErr.Data.Status)
		} else {
			t.Fatalf("expected WPRestError, got: %v", err)
		}

		// Delete non-existent
		_, err = client.ApplicationPasswords().Delete(nonExistentUUID).Do()
		require.Error(t, err)
		if errors.As(err, &wpErr) {
			assert.Equal(t, "rest_application_password_not_found", wpErr.Code)
			assert.Equal(t, 404, wpErr.Data.Status)
		} else {
			t.Fatalf("expected WPRestError, got: %v", err)
		}
	})
}
