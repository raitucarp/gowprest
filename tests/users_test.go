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

func TestUsers(t *testing.T) {
	blogUrl := os.Getenv("BLOG_URL")
	username := os.Getenv("BLOG_USERNAME")
	password := os.Getenv("BLOG_APP_PASSWORD")

	if blogUrl == "" || username == "" || password == "" {
		t.Skip("Skipping test; BLOG_URL, BLOG_USERNAME, or BLOG_APP_PASSWORD not set")
	}

	client := gowprest.NewClient(blogUrl).WithBasicAuth(username, password)
	defer client.Close()

	usersAPI := client.Users()

	// Helper to generate a unique WP-compliant username
	generateUsername := func(prefix string) string {
		cleanWord := strings.ToLower(faker.Word())
		// Remove any non-alphanumeric characters
		cleanWord = strings.Map(func(r rune) rune {
			if (r >= 'a' && r <= 'z') || (r >= '0' && r <= '9') {
				return r
			}
			return -1
		}, cleanWord)
		return fmt.Sprintf("%s%s%d", prefix, cleanWord, time.Now().UnixNano()%100000)
	}

	var user1 *gowprest.User
	var user2 *gowprest.User

	t.Run("CreateUser", func(t *testing.T) {
		t.Run("create_via_chained_setters", func(t *testing.T) {
			uName := generateUsername("user")
			email := fmt.Sprintf("%s@example.com", uName)

			u, err := usersAPI.Create().
				Username(uName).
				Email(email).
				Password("StrongPassword123!@#").
				Name("Test Chained User").
				FirstName("John").
				LastName("Doe").
				Nickname("johndoe").
				Description("Author account created via chained setters").
				Roles("author").
				Do()

			require.NoError(t, err)
			require.NotNil(t, u)
			assert.NotZero(t, u.ID)
			assert.Equal(t, uName, u.Username)
			assert.Equal(t, email, u.Email)
			assert.Equal(t, "John", u.FirstName)
			assert.Equal(t, "Doe", u.LastName)
			assert.Equal(t, "johndoe", u.Nickname)
			assert.Contains(t, u.Roles, "author")
			assert.Equal(t, "Author account created via chained setters", u.Description)

			user1 = u
		})

		t.Run("create_via_struct", func(t *testing.T) {
			uName := generateUsername("structuser")
			email := fmt.Sprintf("%s@example.com", uName)

			userData := gowprest.UserData{
				Username:    uName,
				Email:       email,
				Password:    "AnotherStrongPassword456!@#",
				Name:        "Struct Test User",
				FirstName:   "Jane",
				LastName:    "Smith",
				Nickname:    "janesmith",
				Description: "Created via UserData struct",
				Roles:       []string{"contributor"},
			}

			u, err := usersAPI.Create(userData).Do()
			require.NoError(t, err)
			require.NotNil(t, u)
			assert.NotZero(t, u.ID)
			assert.Equal(t, uName, u.Username)
			assert.Equal(t, email, u.Email)
			assert.Equal(t, "Jane", u.FirstName)
			assert.Equal(t, "Smith", u.LastName)
			assert.Contains(t, u.Roles, "contributor")

			user2 = u
		})

		t.Run("create_validation_error", func(t *testing.T) {
			// Missing email and password
			_, err := usersAPI.Create().
				Username(generateUsername("invalid")).
				Do()

			require.Error(t, err)
			wpErr, ok := err.(*gowprest.WPRestError)
			require.True(t, ok, "expected error to be *gowprest.WPRestError")
			assert.Equal(t, 400, wpErr.Data.Status)
		})
	})

	// Ensure cleanup of user2 at the end of the test suite
	defer func() {
		if user2 != nil && user2.ID != 0 {
			_, _ = usersAPI.Delete(user2.ID).Force().Reassign(1).Do()
		}
	}()

	t.Run("ListUsers", func(t *testing.T) {
		require.NotNil(t, user1, "user1 must exist for list tests")

		t.Run("default_list", func(t *testing.T) {
			users, err := usersAPI.List().Do()
			require.NoError(t, err)
			assert.NotEmpty(t, users)
		})

		t.Run("pagination_and_per_page", func(t *testing.T) {
			users, err := usersAPI.List().
				Page(1).
				PerPage(1).
				Do()
			require.NoError(t, err)
			assert.Len(t, users, 1)

			usersPage2, err := usersAPI.List().
				Page(2).
				PerPage(1).
				Do()
			require.NoError(t, err)
			assert.Len(t, usersPage2, 1)
			assert.NotEqual(t, users[0].ID, usersPage2[0].ID)

			usersOffset, err := usersAPI.List().
				Offset(1).
				PerPage(1).
				Do()
			require.NoError(t, err)
			assert.Len(t, usersOffset, 1)
			assert.Equal(t, usersPage2[0].ID, usersOffset[0].ID)
		})

		t.Run("ordering", func(t *testing.T) {
			ascUsers, err := usersAPI.List().
				OrderAsc().
				OrderByID().
				Do()
			require.NoError(t, err)
			require.True(t, len(ascUsers) >= 2)
			assert.True(t, ascUsers[0].ID <= ascUsers[1].ID)

			descUsers, err := usersAPI.List().
				OrderDesc().
				OrderByID().
				Do()
			require.NoError(t, err)
			require.True(t, len(descUsers) >= 2)
			assert.True(t, descUsers[0].ID >= descUsers[1].ID)
		})

		t.Run("filtering_by_roles", func(t *testing.T) {
			admins, err := usersAPI.List().
				ContextEdit().
				Roles("administrator").
				Do()
			require.NoError(t, err)
			require.NotEmpty(t, admins)
			for _, admin := range admins {
				assert.Contains(t, admin.Roles, "administrator")
			}
		})

		t.Run("filtering_by_who_authors", func(t *testing.T) {
			authors, err := usersAPI.List().
				WhoAuthors().
				Do()
			require.NoError(t, err)
			assert.NotEmpty(t, authors)
		})

		t.Run("filtering_by_include_and_exclude", func(t *testing.T) {
			included, err := usersAPI.List().
				Include(user1.ID).
				Do()
			require.NoError(t, err)
			require.Len(t, included, 1)
			assert.Equal(t, user1.ID, included[0].ID)

			excluded, err := usersAPI.List().
				Exclude(user1.ID).
				Do()
			require.NoError(t, err)
			for _, u := range excluded {
				assert.NotEqual(t, user1.ID, u.ID)
			}
		})

		t.Run("filtering_by_slug", func(t *testing.T) {
			slugUsers, err := usersAPI.List().
				Slug(user1.Slug).
				Do()
			require.NoError(t, err)
			require.Len(t, slugUsers, 1)
			assert.Equal(t, user1.ID, slugUsers[0].ID)
		})

		t.Run("search_query", func(t *testing.T) {
			found, err := usersAPI.List().
				Search(user1.Username).
				Do()
			require.NoError(t, err)
			require.NotEmpty(t, found)
			assert.Equal(t, user1.ID, found[0].ID)
		})

		t.Run("context_and_global_params", func(t *testing.T) {
			editUsers, err := usersAPI.List().
				ContextEdit().
				Include(user1.ID).
				Do()
			require.NoError(t, err)
			require.Len(t, editUsers, 1)
			assert.NotEmpty(t, editUsers[0].Roles)
			assert.NotEmpty(t, editUsers[0].Capabilities)
			assert.Equal(t, user1.Email, editUsers[0].Email)

			fieldsUsers, err := usersAPI.List().
				Fields("id", "name", "slug").
				Include(user1.ID).
				Do()
			require.NoError(t, err)
			require.Len(t, fieldsUsers, 1)
			assert.Equal(t, user1.ID, fieldsUsers[0].ID)
			assert.NotEmpty(t, fieldsUsers[0].Name)
			// Email should not be returned when fields are limited
			assert.Empty(t, fieldsUsers[0].Email)
		})
	})

	t.Run("RetrieveUser", func(t *testing.T) {
		require.NotNil(t, user1, "user1 must exist for retrieve tests")

		t.Run("retrieve_existing_user", func(t *testing.T) {
			u, err := usersAPI.Retrieve(user1.ID).Do()
			require.NoError(t, err)
			require.NotNil(t, u)
			assert.Equal(t, user1.ID, u.ID)
			assert.Equal(t, user1.Name, u.Name)
			assert.Equal(t, user1.Slug, u.Slug)
		})

		t.Run("retrieve_contexts", func(t *testing.T) {
			viewUser, err := usersAPI.Retrieve(user1.ID).
				ContextView().
				Do()
			require.NoError(t, err)
			// In view context without edit capability on field, email/roles may be empty
			assert.Equal(t, user1.ID, viewUser.ID)

			editUser, err := usersAPI.Retrieve(user1.ID).
				ContextEdit().
				Do()
			require.NoError(t, err)
			assert.Equal(t, user1.ID, editUser.ID)
			assert.Equal(t, user1.Email, editUser.Email)
			assert.NotEmpty(t, editUser.Roles)
			assert.NotEmpty(t, editUser.Capabilities)
		})

		t.Run("retrieve_with_fields_and_embed", func(t *testing.T) {
			u, err := usersAPI.Retrieve(user1.ID).
				Fields("id", "name").
				Embed().
				Do()
			require.NoError(t, err)
			require.NotNil(t, u)
			assert.Equal(t, user1.ID, u.ID)
			assert.NotEmpty(t, u.Name)
			assert.Empty(t, u.Description)
		})

		t.Run("retrieve_not_found", func(t *testing.T) {
			_, err := usersAPI.Retrieve(999999999).Do()
			require.Error(t, err)
			wpErr, ok := err.(*gowprest.WPRestError)
			require.True(t, ok, "expected error to be *gowprest.WPRestError")
			assert.Equal(t, 404, wpErr.Data.Status)
			assert.Equal(t, "rest_user_invalid_id", wpErr.Code)
		})
	})

	t.Run("UpdateUser", func(t *testing.T) {
		require.NotNil(t, user1, "user1 must exist for update tests")

		t.Run("update_chained", func(t *testing.T) {
			updatedFirst := "Johnny"
			updatedLast := "Updated"
			updatedDesc := "Updated bio for chained user"

			u, err := usersAPI.Update().
				ID(user1.ID).
				FirstName(updatedFirst).
				LastName(updatedLast).
				Description(updatedDesc).
				Do()

			require.NoError(t, err)
			require.NotNil(t, u)
			assert.Equal(t, user1.ID, u.ID)
			assert.Equal(t, updatedFirst, u.FirstName)
			assert.Equal(t, updatedLast, u.LastName)
			assert.Equal(t, updatedDesc, u.Description)
		})

		t.Run("update_with_struct", func(t *testing.T) {
			updateData := gowprest.UserData{
				ID:          user1.ID,
				Nickname:    "johnny_the_great",
				Description: "Updated via UserData struct",
			}

			u, err := usersAPI.Update(updateData).Do()
			require.NoError(t, err)
			require.NotNil(t, u)
			assert.Equal(t, user1.ID, u.ID)
			assert.Equal(t, "johnny_the_great", u.Nickname)
			assert.Equal(t, "Updated via UserData struct", u.Description)
		})

		t.Run("update_not_found", func(t *testing.T) {
			_, err := usersAPI.Update().
				ID(999999999).
				Description("Should fail").
				Do()
			require.Error(t, err)
			wpErr, ok := err.(*gowprest.WPRestError)
			require.True(t, ok, "expected error to be *gowprest.WPRestError")
			assert.Equal(t, 404, wpErr.Data.Status)
		})
	})

	t.Run("UserMe", func(t *testing.T) {
		t.Run("retrieve_me", func(t *testing.T) {
			me, err := client.Users().Me().Retrieve().
				ContextEdit().
				Do()

			require.NoError(t, err)
			require.NotNil(t, me)
			assert.Equal(t, 1, me.ID)
			assert.Equal(t, username, me.Username)
			assert.Contains(t, me.Roles, "administrator")
			assert.NotEmpty(t, me.Capabilities)
		})

		t.Run("update_me", func(t *testing.T) {
			newBio := "Admin bio updated at " + time.Now().Format(time.RFC3339)
			me, err := client.Users().Me().Update().
				Description(newBio).
				Do()

			require.NoError(t, err)
			require.NotNil(t, me)
			assert.Equal(t, 1, me.ID)
			assert.Equal(t, newBio, me.Description)
		})
	})

	t.Run("DeleteUser", func(t *testing.T) {
		require.NotNil(t, user1, "user1 must exist for delete tests")

		t.Run("delete_with_reassign", func(t *testing.T) {
			del, err := usersAPI.Delete(user1.ID).
				Force().
				Reassign(1).
				Do()

			require.NoError(t, err)
			require.NotNil(t, del)
			assert.True(t, del.Deleted)
			assert.Equal(t, user1.ID, del.ID)

			// Verify user no longer exists
			_, err = usersAPI.Retrieve(user1.ID).Do()
			require.Error(t, err)
			wpErr, ok := err.(*gowprest.WPRestError)
			require.True(t, ok, "expected error to be *gowprest.WPRestError")
			assert.Equal(t, 404, wpErr.Data.Status)

			user1 = nil
		})

		t.Run("delete_not_found", func(t *testing.T) {
			_, err := usersAPI.Delete(999999999).
				Force().
				Reassign(1).
				Do()

			require.Error(t, err)
			wpErr, ok := err.(*gowprest.WPRestError)
			require.True(t, ok, "expected error to be *gowprest.WPRestError")
			assert.Equal(t, 404, wpErr.Data.Status)
		})
	})
}
