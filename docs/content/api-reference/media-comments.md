---
title: "Media & Comments"
description: "Reference for Media uploads and Comment discussions."
---

## Media (`/wp/v2/media`)

```go
mediaAPI := client.Media()

// Upload from raw bytes
m, err := mediaAPI.Create().
	FileBytes("cover.jpg", fileBytes, "image/jpeg").
	Title("Post Cover Image").
	AltText("Cover description").
	Caption("Photo caption").
	Do()

// List media
items, err := mediaAPI.List().
	MediaType("image").
	MimeType("image/jpeg").
	Do()

// Retrieve media
item, err := mediaAPI.Retrieve(mediaID).Do()

// Delete media
deleted, err := mediaAPI.Delete(mediaID).Force().Do()
```

## Comments (`/wp/v2/comments`)

```go
commentsAPI := client.Comments()

// List comments for a post
comments, err := commentsAPI.List().
	Post(postID).
	Order("asc").
	Do()

// Create top-level comment
c, err := commentsAPI.Create().
	Post(postID).
	Content("Great article! Very helpful.").
	AuthorName("Alice").
	AuthorEmail("alice@example.com").
	StatusApprove().
	Do()

// Create nested reply
reply, err := commentsAPI.Create().
	Post(postID).
	Parent(c.ID). // Nested under Alice's comment
	Content("Thanks for reading, Alice!").
	AuthorName("Post Author").
	StatusApprove().
	Do()

// Update comment status
updated, err := commentsAPI.Update().
	ID(c.ID).
	Content("Updated comment text").
	Do()

// Delete comment
deleted, err := commentsAPI.Delete(c.ID).Force().Do()
```
