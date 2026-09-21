---
title: "Gutenberg, Blocks, & FSE"
description: "Reference for Gutenberg Blocks, Patterns, Templates, and Navigation."
---

## Blocks (`/wp/v2/blocks`)

```go
// Reusable editor blocks
blocks, err := client.Blocks().List().Do()
block, err := client.Blocks().Create().
	Title("Call to Action Banner").
	Content("<!-- wp:paragraph --><p>Join now!</p><!-- /wp:paragraph -->").
	Do()
```

## Block Types (`/wp/v2/block-types`)

```go
types, err := client.BlockTypes().List().Do()
pType, err := client.BlockTypes().Retrieve("core/paragraph").Do()
```

## Block Patterns & Categories

```go
// Patterns
patterns, err := client.BlockPatterns().List().Do()

// Pattern Categories
patternCats, err := client.BlockPatternCategories().List().Do()
```

## Block Renderer (`/wp/v2/block-renderer/<name>`)

```go
rendered, err := client.BlockRenderer().
	Render("core/heading").
	Attributes(map[string]any{"content": "Dynamic Heading"}).
	Do()
```

## Block Templates & Template Parts (`/wp/v2/templates`, `/wp/v2/template-parts`)

```go
templates, err := client.Templates().List().Do()
templateParts, err := client.TemplateParts().List().Do()
```

## Navigation Menus & Locations

```go
// Classic Menus
menus, err := client.NavMenus().List().Do()

// Menu Items
items, err := client.NavMenuItems().List().Menus(menuID).Do()

// Menu Locations
locs, err := client.MenuLocations().List().Do()

// FSE Navigation Blocks
navs, err := client.Navigations().List().Do()
```
