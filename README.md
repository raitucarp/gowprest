# gowprest

<div align="center">

[![Go Version](https://img.shields.io/github/go-mod/go-version/raitucarp/gowprest?style=flat-square)](https://go.dev/)
[![GoDoc](https://pkg.go.dev/badge/github.com/raitucarp/gowprest.svg)](https://pkg.go.dev/github.com/raitucarp/gowprest)
[![Documentation](https://img.shields.io/badge/docs-gowprest.raitucarp.name-blue?style=flat-square)](https://gowprest.raitucarp.name)
[![License: MIT](https://img.shields.io/badge/License-MIT-green.svg?style=flat-square)](LICENSE)
[![Build & Tests](https://img.shields.io/badge/tests-41%20suites%20passing-brightgreen?style=flat-square)](https://github.com/raitucarp/gowprest/actions)

**A modern, ergonomic, and fully-typed Go client library for the WordPress REST API.**

[Documentation](https://gowprest.raitucarp.name) • [Quickstart](#quickstart) • [Features](#key-features) • [Endpoints](#supported-endpoints) • [Examples](#usage-examples)

</div>

---

## Overview

`gowprest` provides an intuitive, fluent, and strongly-typed interface to interact with WordPress REST API endpoints. Whether you are building CLI tools, headless CMS frontends, migration pipelines, or automation bots, `gowprest` makes WordPress integration clean and idiomatic in Go.

### Key Features

- **Complete API Coverage**: Supports all 40 WordPress REST API reference endpoints (Posts, Pages, Media, Comments, Taxonomies, Users, Settings, Themes, Site Health, and more).
- **Gutenberg & Full Site Editing (FSE)**: Native support for Blocks, Block Types, Block Patterns, Template Parts, and FSE Navigation.
- **Fluent Builder Pattern**: Type-safe query and mutation builders with sparse field projection (`_fields`), embedding (`_embed`), and context switching (`view`, `edit`, `embed`).
- **Robust Authentication**: Built-in support for WordPress **Application Passwords** and Basic Auth.
- **Media & Binary Uploads**: Seamless multipart and raw byte image/asset uploads with automatic mime-type detection.
- **Revisions & Autosaves**: Full access to version history and drafts for Posts, Pages, Blocks, and Templates.
- **WordPress Site Health**: Diagnostic test runners (`loopback-requests`, `background-updates`, `https-status`, etc.) and directory size metrics.
- **Zero Heavyweight Dependencies**: Lightweight and fast, built on top of `resty/v3`.

---

## Installation

```bash
go get -u github.com/raitucarp/gowprest
```

Requires **Go 1.22+**.

---

## Quickstart

```go
package main

import (
	"context"
	"fmt"
	"log"

	"github.com/raitucarp/gowprest"
)

func main() {
	// Initialize client
	client := gowprest.NewClient("https://your-wordpress-site.com").
		WithBasicAuth("your_username", "xxxx xxxx xxxx xxxx") // Application Password
	defer client.Close()

	// 1. Discover site metadata & health
	siteInfo, err := client.Discover()
	if err != nil {
		log.Fatalf("Discovery failed: %v", err)
	}
	fmt.Printf("Connected to: %s (%s)\n", siteInfo.Name, siteInfo.Home)

	// 2. Query published posts with fluent filters
	posts, err := client.Posts().List().
		Status("publish").
		PerPage(5).
		Order("desc").
		OrderBy("date").
		Do()
	if err != nil {
		log.Fatalf("Failed to fetch posts: %v", err)
	}

	for _, post := range posts {
		fmt.Printf("- [%d] %s\n", post.ID, post.Title.Rendered)
	}

	// 3. Create a new draft post
	newPost, err := client.Posts().Create().
		Title("Announcing gowprest").
		Content("<p>Created with the fluent Go client for WordPress!</p>").
		Status("draft").
		Do()
	if err != nil {
		log.Fatalf("Failed to create post: %v", err)
	}
	fmt.Printf("Created draft post ID: %d\n", newPost.ID)
}
```

---

## Usage Examples

### 1. Uploading Media Assets

```go
imgBytes, _ := os.ReadFile("banner.png")

media, err := client.Media().Create().
	FileBytes("banner.png", imgBytes, "image/png").
	Title("Featured Post Banner").
	AltText("Banner illustration").
	Caption("Created via gowprest").
	Do()
if err != nil {
	log.Fatal(err)
}

fmt.Printf("Uploaded media ID: %d, URL: %s\n", media.ID, media.SourceURL)
```

### 2. Working with Taxonomies (Categories & Tags)

```go
// Create a new category
cat, err := client.Categories().Create().
	Name("Engineering").
	Slug("engineering").
	Description("Technical articles and tutorials").
	Do()

// Assign to post
post, err := client.Posts().Create().
	Title("Go & WordPress Architecture").
	Content("<p>Deep dive into headless WordPress...</p>").
	Categories(cat.ID).
	Status("publish").
	Do()
```

### 3. User & Application Password Management

```go
// Create a dedicated user
user, err := client.Users().Create().
	Username("content-bot").
	Email("bot@example.com").
	Roles("author").
	Password("SecurePassword!123").
	Do()

// Generate an application password for automated service workers
appPass, err := client.ApplicationPasswords(user.ID).Create().
	Name("CI Deployment Token").
	Do()

fmt.Printf("App Password created: %s (Password: %s)\n", appPass.UUID, appPass.Password)
```

### 4. Running WordPress Site Health Diagnostics

```go
// Check HTTPS configuration status
httpsTest, err := client.SiteHealth().HttpsStatus().Execute()
if err != nil {
	log.Fatal(err)
}
fmt.Printf("HTTPS Status: %s - %s\n", httpsTest.Status, httpsTest.Label)

// Inspect directory sizes
sizes, err := client.SiteHealth().DirectorySizes().Execute()
if err != nil {
	log.Fatal(err)
}
fmt.Printf("Total WordPress Size: %s (Uploads: %s)\n", sizes.TotalSize, sizes.UploadsSize)
```

---

## Supported Endpoints

`gowprest` covers the complete WordPress REST API specification:

| Module | Endpoints Covered | Key Client Methods |
|---|---|---|
| **Content & Publishing** | Posts, Pages, Revisions, Autosaves | `client.Posts()`, `client.Pages()`, `client.PostRevisions()`, `client.PageRevisions()` |
| **Media & Files** | Media Library, Uploads, Revisions | `client.Media()` |
| **Taxonomies** | Categories, Tags, Taxonomies, Types, Statuses | `client.Categories()`, `client.Tags()`, `client.Taxonomies()`, `client.Statuses()`, `client.Types()` |
| **Users & Security** | Users, Current User (`/me`), Application Passwords | `client.Users()`, `client.Users().Me()`, `client.ApplicationPasswords(userID)` |
| **Comments & Engagement** | Comments, Nested Replies, Moderation | `client.Comments()` |
| **Gutenberg & FSE** | Blocks, Block Types, Patterns, Pattern Categories, Pattern Directory | `client.Blocks()`, `client.BlockTypes()`, `client.BlockPatterns()`, `client.BlockRenderer()` |
| **Navigation & Menus** | Classic Nav Menus, Menu Items, Locations, FSE Navigation | `client.NavMenus()`, `client.NavMenuItems()`, `client.MenuLocations()`, `client.Navigations()` |
| **Templates & Widgets** | Block Templates, Template Parts, Sidebars, Widgets, Widget Types | `client.Templates()`, `client.TemplateParts()`, `client.Sidebars()`, `client.Widgets()` |
| **Administration** | Site Settings, Active Themes, Global Styles, Plugins, Site Health | `client.Settings()`, `client.Themes()`, `client.GlobalStyles()`, `client.Plugins()`, `client.SiteHealth()` |
| **Search & Discovery** | Global Cross-Type Search, Directory Search | `client.Search()`, `client.DirectorySearch()` |

---

## Documentation

Full documentation, architecture guides, and interactive cookbook examples are hosted at:

👉 **[https://gowprest.raitucarp.name](https://gowprest.raitucarp.name)**

---

## Running Tests

Integration tests run against a local WordPress instance configured with Application Passwords.

```bash
# Set environment credentials in tests/.env or shell:
export BLOG_URL="http://localhost:8080"
export BLOG_USERNAME="admin"
export BLOG_APP_PASSWORD="xxxx xxxx xxxx xxxx"

# Run all 41 test suites
go test -v -count=1 ./tests

# Run the complete end-to-end blog lifecycle simulation test
go test -v -count=1 ./tests -run TestEndToEndBlogSimulation
```

---

## Contributing

Contributions, bug reports, and suggestions are warmly welcome! Please feel free to open an issue or submit a pull request.

1. Fork the repository
2. Create your feature branch (`git checkout -b feat/my-feature`)
3. Commit your changes (`git commit -m 'feat: add support for custom endpoint'`)
4. Push to the branch (`git push origin feat/my-feature`)
5. Open a Pull Request

---

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.
