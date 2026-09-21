---
title: "Content APIs: Posts & Pages"
description: "Reference for Posts, Pages, Revisions, and Autosaves."
---

## Posts (`/wp/v2/posts`)

```go
// Access service
posts := client.Posts()

// List posts
list := posts.List().
	Context("view|edit|embed").
	Page(1).
	PerPage(10).
	Search("keyword").
	After(time.Now().Add(-24 * time.Hour)).
	Author(authorID).
	Categories(catID).
	Tags(tagID).
	Status("publish|draft|pending|private|future").
	Sticky(true).
	Order("desc").
	OrderBy("date").
	Embed().
	Fields("id", "title", "content")
postsList, err := list.Do()

// Retrieve post
p, err := posts.Retrieve(postID).
	ContextEdit().
	Password("secret").
	Do()

// Create post
created, err := posts.Create().
	Title("Title").
	Content("<p>Content</p>").
	Excerpt("Summary").
	StatusPublish().
	Categories(catID).
	Tags(tagID).
	FeaturedMedia(mediaID).
	Do()

// Update post
updated, err := posts.Update().
	ID(postID).
	Title("New Title").
	Do()

// Delete post
deleted, err := posts.Delete(postID).
	Force(). // Bypass trash
	Do()
```

## Post Revisions (`/wp/v2/posts/<parent>/revisions`)

```go
revs, err := client.PostRevisions(postID).List().Do()
rev, err := client.PostRevisions(postID).Retrieve(revisionID).Do()
deleted, err := client.PostRevisions(postID).Delete(revisionID).Do()
```

## Pages (`/wp/v2/pages`)

```go
pages := client.Pages()

// List pages
pagesList, err := pages.List().
	Parent(parentPageID).
	Order("asc").
	OrderBy("menu_order").
	Do()

// Create page
p, err := pages.Create().
	Title("Page Title").
	Content("<p>Body</p>").
	Parent(parentID).
	MenuOrder(1).
	StatusPublish().
	Do()
```
