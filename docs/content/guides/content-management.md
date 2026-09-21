---
title: "Content Management: Posts, Pages, Revisions, and Autosaves"
description: "How to draft, publish, update, and manage revision history for Posts and Pages."
---

## Managing Posts

### Creating Posts
Use the fluent builder on `client.Posts().Create()`:

```go
post, err := client.Posts().Create().
	Title("Building Microservices in Go").
	Content("<p>Microservices with Go and WordPress headless CMS...</p>").
	Excerpt("A quick summary of microservice architecture.").
	StatusPublish().
	Sticky(true).
	Categories(12).
	Tags(4, 7).
	Do()
```

### Updating Posts
```go
updated, err := client.Posts().Update().
	ID(post.ID).
	Title("Building Microservices in Go (Updated Edition)").
	Do()
```

### Retrieving a Post by ID
```go
p, err := client.Posts().Retrieve(post.ID).
	ContextEdit(). // Retrieve in edit context to view raw markup
	Do()
```

### Deleting Posts
By default, deleting moves the post to the trash:
```go
// Move to trash
deleted, err := client.Posts().Delete(post.ID).Do()

// Permanently delete (force)
deleted, err = client.Posts().Delete(post.ID).Force().Do()
```

---

## Working with Revisions & Autosaves

WordPress tracks every edit made to posts and pages as revisions.

### Inspecting Post Revisions
```go
revisions, err := client.PostRevisions(post.ID).List().Do()
for _, rev := range revisions {
	fmt.Printf("Revision ID: %d, Modified: %s\n", rev.ID, rev.Date.Format(time.RFC3339))
}
```

### Post Autosaves
```go
autosaves, err := client.Posts().Autosaves(post.ID).List().Do()
```

---

## Managing Pages

Pages follow the same fluent pattern as Posts, with page-specific properties such as `parent` and `menu_order`:

```go
// Create parent page
aboutPage, err := client.Pages().Create().
	Title("About Our Company").
	Content("<p>Our history and mission statement.</p>").
	StatusPublish().
	MenuOrder(1).
	Do()

// Create child sub-page
teamPage, err := client.Pages().Create().
	Title("Executive Team").
	Content("<p>Meet the team behind our platform.</p>").
	Parent(aboutPage.ID).
	MenuOrder(2).
	StatusPublish().
	Do()
```
