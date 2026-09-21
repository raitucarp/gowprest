---
title: "Quickstart"
description: "Start querying and publishing WordPress content with gowprest in under 5 minutes."
---

In this quickstart guide, you will:
1. Initialize the `gowprest` client.
2. Discover site metadata and available namespaces.
3. List recent published posts with pagination.
4. Create and publish a new blog post.

## Complete Example

```go
package main

import (
	"fmt"
	"log"

	"github.com/raitucarp/gowprest"
)

func main() {
	// 1. Initialize client
	client := gowprest.NewClient("https://example.com").
		WithBasicAuth("my_user", "xxxx xxxx xxxx xxxx")
	defer client.Close()

	// 2. Discover site metadata
	info, err := client.Discover()
	if err != nil {
		log.Fatalf("Failed site discovery: %v", err)
	}
	fmt.Printf("Connected to site: %s\nSite tagline: %s\nTimezone: %s\n\n",
		info.Name, info.Description, info.TimezoneString)

	// 3. Query the latest 5 published posts
	posts, err := client.Posts().List().
		Status("publish").
		PerPage(5).
		Order("desc").
		OrderBy("date").
		Embed(). // Includes author and featured media in response
		Do()
	if err != nil {
		log.Fatalf("Failed to fetch posts: %v", err)
	}

	fmt.Println("Recent Posts:")
	for _, p := range posts {
		fmt.Printf("- [%d] %s (%s)\n", p.ID, p.Title.Rendered, p.Date.Format("2006-01-02"))
	}

	// 4. Create and publish a new post
	newPost, err := client.Posts().Create().
		Title("Hello from Go!").
		Content("<p>This post was authored using the <strong>gowprest</strong> client library.</p>").
		StatusPublish().
		Do()
	if err != nil {
		log.Fatalf("Failed to create post: %v", err)
	}

	fmt.Printf("\nSuccessfully created post ID: %d!\nView post at: %s\n", newPost.ID, newPost.Link)
}
```

## Next Steps

- Check out the [Content Management Guide](/guides/content-management/) to learn about drafts, revisions, and autosaves.
- Learn how to upload images in the [Media Uploads Guide](/guides/media-uploads/).
- Read the [E2E Blog Simulation](/guides/simulation/) for an end-to-end walkthrough across all features.
