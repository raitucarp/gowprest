---
title: "Taxonomies, Categories, & Tags"
description: "Reference for Categories, Tags, and Custom Taxonomies."
---

## Categories (`/wp/v2/categories`)

```go
cats := client.Categories()

// List
items, err := cats.List().
	Parent(parentID).
	HideEmpty(true).
	Do()

// Create
c, err := cats.Create().
	Name("Web Development").
	Slug("web-development").
	Description("Tutorials and news").
	Parent(0).
	Do()

// Retrieve
c, err = cats.Retrieve(categoryID).Do()

// Update
c, err = cats.Update().
	ID(categoryID).
	Name("Frontend & Web Dev").
	Do()

// Delete
deleted, err := cats.Delete(categoryID).Force().Do()
```

## Tags (`/wp/v2/tags`)

```go
tags := client.Tags()

// List
tList, err := tags.List().Do()

// Create
t, err := tags.Create().
	Name("Docker").
	Slug("docker").
	Do()

// Delete
deleted, err := tags.Delete(tagID).Do()
```

## Taxonomies Registry (`/wp/v2/taxonomies`)

```go
taxonomies, err := client.Taxonomies().List().Do()
tax, err := client.Taxonomies().Retrieve("category").Do()
```
