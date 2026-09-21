---
title: "Media & Asset Uploads"
description: "Uploading images, documents, and media attachments to the WordPress media library."
---

WordPress accepts media uploads either via multipart forms or binary payloads with raw bytes and a `Content-Disposition` header. `gowprest` simplifies both approaches through its `Media()` builder.

## Uploading from Memory (Raw Bytes)

When generating images in memory or reading files:

```go
package main

import (
	"os"
	"log"
	"github.com/raitucarp/gowprest"
)

func main() {
	client := gowprest.NewClient("https://example.com").
		WithBasicAuth("admin", "app_password")
	defer client.Close()

	// Read binary file from disk
	imgData, err := os.ReadFile("hero-banner.png")
	if err != nil {
		log.Fatal(err)
	}

	// Upload asset
	media, err := client.Media().Create().
		FileBytes("hero-banner.png", imgData, "image/png").
		Title("Hero Banner 2026").
		AltText("Modern digital illustration of cloud architecture").
		Caption("Uploaded programmatically via gowprest").
		Description("Full-width vector graphic for the homepage hero").
		Do()
	if err != nil {
		log.Fatalf("Media upload failed: %v", err)
	}

	log.Printf("Media uploaded! ID: %d, Source URL: %s\n", media.ID, media.SourceURL)
}
```

## Assigning Featured Media to a Post

Once uploaded, use the returned media ID as the `FeaturedMedia` property when creating or updating a post:

```go
post, err := client.Posts().Create().
	Title("Cloud Computing Trends").
	Content("<p>Analysis of distributed computing in 2026.</p>").
	FeaturedMedia(media.ID). // Links the uploaded image as featured banner
	StatusPublish().
	Do()
```

## Querying the Media Library

Filter media items by type (`image`, `video`, `audio`, etc.) and date:

```go
images, err := client.Media().List().
	MediaType("image").
	PerPage(10).
	Do()

for _, img := range images {
	log.Printf("- [%d] %s (%s)\n", img.ID, img.Title.Rendered, img.SourceURL)
}
```
