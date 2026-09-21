package tests

import (
	"errors"
	"fmt"
	"os"
	"testing"
	"time"

	_ "github.com/joho/godotenv/autoload"
	"github.com/raitucarp/gowprest"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// 1x1 transparent PNG for media upload simulation
var simPNG = []byte{
	0x89, 0x50, 0x4E, 0x47, 0x0D, 0x0A, 0x1A, 0x0A, 0x00, 0x00, 0x00, 0x0D,
	0x49, 0x48, 0x44, 0x52, 0x00, 0x00, 0x00, 0x01, 0x00, 0x00, 0x00, 0x01,
	0x08, 0x06, 0x00, 0x00, 0x00, 0x1F, 0x15, 0xC4, 0x89, 0x00, 0x00, 0x00,
	0x0A, 0x49, 0x44, 0x41, 0x54, 0x78, 0x9C, 0x63, 0x00, 0x01, 0x00, 0x00,
	0x05, 0x00, 0x01, 0x0D, 0x0A, 0x2D, 0xB4, 0x00, 0x00, 0x00, 0x00, 0x49,
	0x45, 0x4E, 0x44, 0xAE, 0x42, 0x60, 0x82,
}

func TestEndToEndBlogSimulation(t *testing.T) {
	blogUrl := os.Getenv("BLOG_URL")
	username := os.Getenv("BLOG_USERNAME")
	password := os.Getenv("BLOG_APP_PASSWORD")

	if blogUrl == "" || username == "" || password == "" {
		t.Skip("Skipping simulation; BLOG_URL, BLOG_USERNAME, or BLOG_APP_PASSWORD not set")
	}

	client := gowprest.NewClient(blogUrl).WithBasicAuth(username, password)
	defer client.Close()

	uniqueSuffix := fmt.Sprintf("%d", time.Now().UnixNano()%1000000)

	// Simulation State shared across lifecycle steps
	var (
		createdUser         *gowprest.User
		createdAppPassword  *gowprest.ApplicationPassword
		createdCategory     gowprest.Category
		createdTag          gowprest.Tag
		createdMedia        gowprest.Media
		createdPost         gowprest.Post
		createdComment      gowprest.Comment
		createdReplyComment gowprest.Comment
		createdPage         gowprest.Page
		createdMenu         *gowprest.NavMenu
		createdMenuItem     *gowprest.NavMenuItem
		createdNav          *gowprest.Navigation
		createdWidget       *gowprest.Widget
		createdTemplatePart *gowprest.TemplatePart
	)

	// =========================================================================
	// STEP 1: Site Discovery & Health Diagnostics
	// =========================================================================
	t.Run("Step1_SiteDiscoveryAndHealth", func(t *testing.T) {
		// 1. Discover blog info
		info, err := client.Discover()
		require.NoError(t, err)
		assert.NotEmpty(t, info.Name)
		assert.NotEmpty(t, info.URL)
		assert.NotEmpty(t, info.Namespaces)

		// 2. Run Site Health tests
		sh := client.SiteHealth()
		bg, err := sh.BackgroundUpdates().Execute()
		require.NoError(t, err)
		assert.Equal(t, "background_updates", bg.Test)

		https, err := sh.HttpsStatus().Execute()
		require.NoError(t, err)
		assert.Equal(t, "https_status", https.Test)

		dotorg, err := sh.DotorgCommunication().Execute()
		require.NoError(t, err)
		assert.Equal(t, "dotorg_communication", dotorg.Test)

		authHdr, err := sh.AuthorizationHeader().Execute()
		require.NoError(t, err)
		assert.Equal(t, "authorization_header", authHdr.Test)

		dirSizes, err := sh.DirectorySizes().Execute()
		require.NoError(t, err)
		require.NotNil(t, dirSizes)
		assert.True(t, dirSizes.WordPressSize.Size != "" || dirSizes.DatabaseSize.Size != "")
	})

	// =========================================================================
	// STEP 2: Site Settings & Theme Inspection
	// =========================================================================
	t.Run("Step2_SiteSettingsAndTheme", func(t *testing.T) {
		settingsAPI := client.Settings()
		origSettings, err := settingsAPI.Get()
		require.NoError(t, err)
		require.NotNil(t, origSettings)

		// Temporarily update description (tagline)
		tempTagline := "Empowering Go developers with WordPress " + uniqueSuffix
		updatedSettings, err := settingsAPI.Update().
			Description(tempTagline).
			Do()
		require.NoError(t, err)
		assert.Equal(t, tempTagline, updatedSettings.Description)

		// Revert back
		_, err = settingsAPI.Update().
			Description(origSettings.Description).
			Do()
		require.NoError(t, err)

		// Inspect themes
		themes, err := client.Themes().List().StatusActive().Do()
		require.NoError(t, err)
		require.NotEmpty(t, themes)
		activeTheme := themes[0]
		assert.NotEmpty(t, activeTheme.Stylesheet)

		// Inspect Global Styles for theme
		themeStyles, err := client.GlobalStyles().Themes(activeTheme.Stylesheet).Retrieve().Do()
		require.NoError(t, err)
		require.NotNil(t, themeStyles)
	})

	// =========================================================================
	// STEP 3: Taxonomies (Categories & Tags) Setup
	// =========================================================================
	t.Run("Step3_TaxonomiesSetup", func(t *testing.T) {
		// Verify registered taxonomies
		taxList, err := client.Taxonomies().List().Do()
		require.NoError(t, err)
		assert.Contains(t, taxList, "category")
		assert.Contains(t, taxList, "post_tag")

		// Create Category
		catName := "Technology " + uniqueSuffix
		cat, err := client.Categories().Create().
			Name(catName).
			Description("Tech articles and tutorials").
			Slug("tech-" + uniqueSuffix).
			Do()
		require.NoError(t, err)
		assert.Equal(t, catName, cat.Name)
		createdCategory = cat

		// Create Tag
		tagName := "Golang " + uniqueSuffix
		tag, err := client.Tags().Create().
			Name(tagName).
			Description("Posts regarding Go development").
			Slug("golang-" + uniqueSuffix).
			Do()
		require.NoError(t, err)
		assert.Equal(t, tagName, tag.Name)
		createdTag = tag
	})

	// =========================================================================
	// STEP 4: User Onboarding & Application Passwords
	// =========================================================================
	t.Run("Step4_UserOnboardingAndAppPassword", func(t *testing.T) {
		uName := "author" + uniqueSuffix
		email := fmt.Sprintf("%s@example.com", uName)

		user, err := client.Users().Create().
			Username(uName).
			Email(email).
			Password("SuperSecretP@ssw0rd!#99").
			Name("Jane Doe " + uniqueSuffix).
			FirstName("Jane").
			LastName("Doe").
			Roles("editor").
			Do()
		require.NoError(t, err)
		require.NotNil(t, user)
		assert.Equal(t, uName, user.Username)
		createdUser = user

		// Issue Application Password for the new editor
		appPass, err := client.ApplicationPasswords(user.ID).
			Create("Mobile Editor App " + uniqueSuffix).
			Do()
		require.NoError(t, err)
		require.NotNil(t, appPass)
		assert.NotEmpty(t, appPass.UUID)
		assert.NotEmpty(t, appPass.Password)
		createdAppPassword = appPass

		// Retrieve the created application password
		retrievedPass, err := client.ApplicationPasswords(user.ID).
			Retrieve(appPass.UUID).
			Do()
		require.NoError(t, err)
		assert.Equal(t, appPass.UUID, retrievedPass.UUID)

		// Verify Current Authenticated User (Admin)
		me, err := client.Users().Me().Retrieve().
			ContextEdit().
			Do()
		require.NoError(t, err)
		require.NotNil(t, me)
		assert.Equal(t, username, me.Username)
		assert.Contains(t, me.Roles, "administrator")
	})

	// =========================================================================
	// STEP 5: Media Asset Upload
	// =========================================================================
	t.Run("Step5_MediaAssetUpload", func(t *testing.T) {
		m, err := client.Media().Create().
			FileBytes("featured-banner-"+uniqueSuffix+".png", simPNG, "image/png").
			Title("Featured Banner " + uniqueSuffix).
			AltText("A descriptive alt text for banner").
			Caption("Banner caption").
			Do()
		require.NoError(t, err)
		assert.NotZero(t, m.ID)
		assert.Equal(t, "image/png", m.MimeType)
		createdMedia = m

		// Update media metadata
		updatedMedia, err := client.Media().Update().
			ID(m.ID).
			Description("Updated media description via gowprest").
			Do()
		require.NoError(t, err)
		assert.Contains(t, updatedMedia.Description.Rendered, "Updated media description")
	})

	// =========================================================================
	// STEP 6: Content Authoring (Post, Draft, Revisions, Publishing)
	// =========================================================================
	t.Run("Step6_ContentAuthoringAndPublishing", func(t *testing.T) {
		require.NotZero(t, createdCategory.ID)
		require.NotZero(t, createdTag.ID)
		require.NotZero(t, createdMedia.ID)
		require.NotNil(t, createdUser)

		postTitle := "Building Modern Web Services with Go " + uniqueSuffix
		initialContent := "<p>Go is an incredible language for high concurrency web backends.</p>"

		// 1. Create Draft Post assigned to the new author, category, tag, and featured media
		post, err := client.Posts().Create().
			Title(postTitle).
			Content(initialContent).
			Author(createdUser.ID).
			Categories(createdCategory.ID).
			Tags(createdTag.ID).
			FeaturedMedia(createdMedia.ID).
			StatusDraft().
			Do()
		require.NoError(t, err)
		assert.Equal(t, gowprest.StatusDraft, post.Status)
		assert.Equal(t, postTitle, post.Title.Rendered)
		createdPost = post

		// 2. Create an Autosave for the draft
		autosaveContent := "<p>Work in progress autosave content for post " + uniqueSuffix + "</p>"
		autosave, err := client.Posts().Revisions(post.ID).Autosaves().Create().
			Content(autosaveContent).
			Do()
		require.NoError(t, err)
		require.NotNil(t, autosave)
		assert.Equal(t, post.ID, autosave.Parent)

		// 3. Update the post with second revision and Publish
		updatedContent := "<p>Comprehensive guide to building production WordPress integrations with gowprest.</p>"
		publishedPost, err := client.Posts().Update().
			ID(post.ID).
			Content(updatedContent).
			StatusPublish().
			Do()
		require.NoError(t, err)
		assert.Equal(t, gowprest.StatusPublished, publishedPost.Status)

		// 4. Update a third time to create another revision
		finalContent := "<p>Final revised edition with updated REST patterns and examples.</p>"
		_, err = client.Posts().Update().
			ID(post.ID).
			Content(finalContent).
			Do()
		require.NoError(t, err)

		// 5. Query Revisions
		revisions, err := client.Posts().Revisions(post.ID).List().Do()
		require.NoError(t, err)
		assert.NotEmpty(t, revisions)
		for _, rev := range revisions {
			assert.Equal(t, post.ID, rev.Parent)
		}
	})

	// =========================================================================
	// STEP 7: Audience Engagement (Comments & Nested Replies)
	// =========================================================================
	t.Run("Step7_AudienceEngagementComments", func(t *testing.T) {
		require.NotZero(t, createdPost.ID)

		// 1. Visitor leaves a top-level comment
		comment1, err := client.Comments().Create().
			Post(createdPost.ID).
			AuthorName("Alice Visitor").
			AuthorEmail("alice@example.org").
			Content("Outstanding post! I loved the clear architecture.").
			StatusApprove().
			Do()
		require.NoError(t, err)
		assert.Equal(t, createdPost.ID, comment1.Post)
		createdComment = comment1

		// 2. Author replies to the comment (nested reply)
		reply, err := client.Comments().Create().
			Post(createdPost.ID).
			Parent(comment1.ID).
			AuthorName("Jane Doe (Author)").
			AuthorEmail(createdUser.Email).
			Content("Thank you Alice! More tutorials on gowprest coming soon.").
			StatusApprove().
			Do()
		require.NoError(t, err)
		assert.Equal(t, comment1.ID, reply.Parent)
		createdReplyComment = reply

		// 3. Verify comments listed under this post
		commentsList, err := client.Comments().List().
			Post(createdPost.ID).
			Do()
		require.NoError(t, err)
		assert.GreaterOrEqual(t, len(commentsList), 2)
	})

	// =========================================================================
	// STEP 8: Static Pages (Pages & Hierarchy)
	// =========================================================================
	t.Run("Step8_StaticPagesSetup", func(t *testing.T) {
		pageTitle := "About Us " + uniqueSuffix
		pageContent := "<p>Welcome to our tech publication powered by WordPress & gowprest.</p>"

		page, err := client.Pages().Create().
			Title(pageTitle).
			Content(pageContent).
			StatusPublish().
			Do()
		require.NoError(t, err)
		assert.Equal(t, gowprest.StatusPublished, page.Status)
		assert.Equal(t, pageTitle, page.Title.Rendered)
		createdPage = page

		// Update page to verify revision generation
		updatedPage, err := client.Pages().Update().
			ID(page.ID).
			Content("<p>Updated About Us page content with team information.</p>").
			Do()
		require.NoError(t, err)
		assert.NotEmpty(t, updatedPage.ID)

		// List page revisions
		pageRevs, err := client.Pages().Revisions(page.ID).List().Do()
		require.NoError(t, err)
		assert.NotEmpty(t, pageRevs)
	})

	// =========================================================================
	// STEP 9: Navigation & Menu Management
	// =========================================================================
	t.Run("Step9_NavigationAndMenuManagement", func(t *testing.T) {
		require.NotZero(t, createdPage.ID)

		// 1. Create Nav Menu
		menuName := "Main Site Nav " + uniqueSuffix
		menu, err := client.NavMenus().Create().
			Name(menuName).
			Description("Primary navigation menu for the site").
			Do()
		require.NoError(t, err)
		require.NotNil(t, menu)
		createdMenu = menu

		// 2. Add Menu Item pointing to Home
		item1, err := client.NavMenuItems().Create().
			Menus(menu.ID).
			Title("Home").
			URL("/").
			Status("publish").
			Do()
		require.NoError(t, err)
		assert.Equal(t, "Home", item1.TitleString())

		// 3. Add Menu Item pointing to About page
		item2, err := client.NavMenuItems().Create().
			Menus(menu.ID).
			Title("About").
			URL("/about-" + uniqueSuffix).
			Status("publish").
			Do()
		require.NoError(t, err)
		assert.Equal(t, "About", item2.TitleString())
		createdMenuItem = item2

		// 4. Create Navigation post type (FSE Navigation)
		navBlock, err := client.Navigations().Create().
			Title("FSE Header Navigation " + uniqueSuffix).
			Content(`<!-- wp:navigation-link {"label":"Home","url":"/"} /--><!-- wp:navigation-link {"label":"Blog","url":"/blog"} /-->`).
			StatusPublish().
			Do()
		require.NoError(t, err)
		require.NotNil(t, navBlock)
		createdNav = navBlock
	})

	// =========================================================================
	// STEP 10: Sidebars & Widgets
	// =========================================================================
	t.Run("Step10_SidebarsAndWidgets", func(t *testing.T) {
		// Inspect registered widget types
		types, err := client.WidgetTypes().List().Do()
		require.NoError(t, err)
		assert.NotEmpty(t, types)

		// Create a block widget
		widgetContent := "<!-- wp:paragraph --><p>Welcome to our tech publication! " + uniqueSuffix + "</p><!-- /wp:paragraph -->"
		widget, err := client.Widgets().Create().
			IDBase("block").
			Instance(map[string]any{
				"raw": map[string]any{
					"content": widgetContent,
				},
			}).
			Do()
		require.NoError(t, err)
		require.NotNil(t, widget)
		assert.Equal(t, "block", widget.IDBase)
		createdWidget = widget

		// Update widget
		updatedWidget, err := client.Widgets().Update(widget.ID).
			Instance(map[string]any{
				"raw": map[string]any{
					"content": "<!-- wp:paragraph --><p>Updated sidebar widget content</p><!-- /wp:paragraph -->",
				},
			}).
			Do()
		require.NoError(t, err)
		assert.Equal(t, widget.ID, updatedWidget.ID)
	})

	// =========================================================================
	// STEP 11: Full Site Editing (Templates & Template Parts)
	// =========================================================================
	t.Run("Step11_TemplatesAndTemplateParts", func(t *testing.T) {
		// List templates
		templates, err := client.Templates().List().Do()
		require.NoError(t, err)
		assert.NotEmpty(t, templates)

		// Create a custom template part (Header)
		partSlug := "custom-header-" + uniqueSuffix
		part, err := client.TemplateParts().Create().
			Slug(partSlug).
			Title("Custom Header " + uniqueSuffix).
			Area("header").
			Content("<!-- wp:group --><header><h1>My Tech Blog</h1></header><!-- /wp:group -->").
			Do()
		require.NoError(t, err)
		require.NotNil(t, part)
		assert.Equal(t, "header", part.Area)
		createdTemplatePart = part

		// Update template part
		updatedPart, err := client.TemplateParts().Update(part.ID).
			Title("Custom Header Updated " + uniqueSuffix).
			Do()
		require.NoError(t, err)
		assert.Equal(t, "Custom Header Updated "+uniqueSuffix, updatedPart.TitleString())
	})

	// =========================================================================
	// STEP 12: Global Search & Cross-Resource Discovery
	// =========================================================================
	t.Run("Step12_GlobalSearchAndDiscovery", func(t *testing.T) {
		require.NotZero(t, createdPost.ID)

		// Perform global search for unique suffix
		searchResults, err := client.Search().
			Query(uniqueSuffix).
			PerPage(10).
			Do()
		require.NoError(t, err)
		require.NotEmpty(t, searchResults, "Expected search to discover entities containing the unique suffix")

		var foundPost bool
		for _, item := range searchResults {
			if item.IntID() == createdPost.ID {
				foundPost = true
				assert.Equal(t, "post", item.Type)
				assert.Equal(t, "post", item.Subtype)
				break
			}
		}
		assert.True(t, foundPost, "Newly published post should be indexed in global search results")
	})

	// =========================================================================
	// STEP 13: Gutenberg Ecosystem, Block Renderer & Directories
	// =========================================================================
	t.Run("Step13_GutenbergEcosystem", func(t *testing.T) {
		// 1. Dynamic Block Renderer
		rendered, err := client.BlockRenderer().Render("core/calendar").Do()
		require.NoError(t, err)
		require.NotNil(t, rendered)
		assert.Contains(t, rendered.Rendered, "wp-block-calendar")

		// 2. Block Types
		bTypes, err := client.BlockTypes().List().Namespace("core").Do()
		require.NoError(t, err)
		assert.NotEmpty(t, bTypes)

		// 3. Block Patterns & Categories
		patterns, err := client.BlockPatterns().List().Do()
		require.NoError(t, err)
		assert.NotEmpty(t, patterns)

		categories, err := client.BlockPatternCategories().List().Do()
		require.NoError(t, err)
		assert.NotEmpty(t, categories)

		// 4. Pattern Directory
		patDir, err := client.PatternDirectory().Search("hero").PerPage(2).Do()
		require.NoError(t, err)
		assert.NotEmpty(t, patDir)

		// 5. Block Directory
		blkDir, err := client.BlockDirectory().Search("slider").PerPage(2).Do()
		require.NoError(t, err)
		assert.NotEmpty(t, blkDir)
	})

	// =========================================================================
	// STEP 14: Graceful Cleanup & Teardown
	// =========================================================================
	t.Run("Step14_GracefulTeardown", func(t *testing.T) {
		// Clean up Template Part
		if createdTemplatePart != nil {
			delPart, err := client.TemplateParts().Delete(createdTemplatePart.ID).Force().Do()
			assert.NoError(t, err)
			assert.NotNil(t, delPart)
		}

		// Clean up Widget
		if createdWidget != nil {
			delWidget, err := client.Widgets().Delete(createdWidget.ID).Force().Do()
			assert.NoError(t, err)
			assert.NotNil(t, delWidget)
		}

		// Clean up Navigations
		if createdNav != nil {
			_, err := client.Navigations().Delete(createdNav.ID).Force().Do()
			assert.NoError(t, err)
		}

		// Clean up Menu Items & Menu
		if createdMenuItem != nil {
			_, _ = client.NavMenuItems().Delete(createdMenuItem.ID).Force().Do()
		}
		if createdMenu != nil {
			_, err := client.NavMenus().Delete(createdMenu.ID).Force().Do()
			assert.NoError(t, err)
		}

		// Clean up Comments
		if createdReplyComment.ID != 0 {
			_, _ = client.Comments().Delete(createdReplyComment.ID).Force().Do()
		}
		if createdComment.ID != 0 {
			_, _ = client.Comments().Delete(createdComment.ID).Force().Do()
		}

		// Clean up Page
		if createdPage.ID != 0 {
			_, err := client.Pages().Delete(createdPage.ID).Force().Do()
			assert.NoError(t, err)
		}

		// Clean up Post
		if createdPost.ID != 0 {
			_, err := client.Posts().Delete(createdPost.ID).Force().Do()
			assert.NoError(t, err)
		}

		// Clean up Media
		if createdMedia.ID != 0 {
			_, err := client.Media().Delete(createdMedia.ID).Force().Do()
			assert.NoError(t, err)
		}

		// Clean up Tag & Category
		if createdTag.ID != 0 {
			_, err := client.Tags().Delete(createdTag.ID).Force().Do()
			assert.NoError(t, err)
		}
		if createdCategory.ID != 0 {
			_, err := client.Categories().Delete(createdCategory.ID).Force().Do()
			assert.NoError(t, err)
		}

		// Clean up Application Password & User
		if createdUser != nil {
			if createdAppPassword != nil {
				_, _ = client.ApplicationPasswords(createdUser.ID).Delete(createdAppPassword.UUID).Do()
			}
			_, err := client.Users().Delete(createdUser.ID).Force().Reassign(1).Do()
			assert.NoError(t, err)
		}

		// Verify 404 on deleted post
		if createdPost.ID != 0 {
			_, err := client.Posts().Retrieve(createdPost.ID).Do()
			require.Error(t, err)
			var restErr *gowprest.WPRestError
			require.True(t, errors.As(err, &restErr))
			assert.Equal(t, 404, restErr.Data.Status)
		}
	})
}
