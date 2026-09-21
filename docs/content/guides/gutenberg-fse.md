---
title: "Gutenberg, Block Patterns, & Full Site Editing"
description: "Managing Gutenberg block types, block patterns, templates, template parts, and block navigation."
---

WordPress modern themes utilize **Full Site Editing (FSE)** and the **Gutenberg block editor**. `gowprest` provides dedicated APIs for the entire block ecosystem.

## Block Patterns & Categories

### Listing Block Pattern Categories
```go
cats, err := client.BlockPatternCategories().List().Do()
for _, c := range cats {
	fmt.Printf("Pattern Category: %s (%s)\n", c.Label, c.Name)
}
```

### Listing Registered Block Patterns
```go
patterns, err := client.BlockPatterns().List().Do()
for _, p := range patterns {
	fmt.Printf("Pattern: %s - Categories: %v\n", p.Title, p.Categories)
}
```

---

## Server-Side Block Rendering

Render block HTML directly on the server without needing a browser or Node runtime:

```go
rendered, err := client.BlockRenderer().
	Render("core/paragraph").
	Attributes(map[string]any{
		"content": "Rendered dynamically with gowprest!",
	}).
	Do()

fmt.Println("Rendered HTML:", rendered.Rendered)
```

---

## Full Site Editing: Templates & Template Parts

Block themes (such as Twenty Twenty-Four) define structural layouts in templates and template parts:

### Inspecting Block Templates
```go
templates, err := client.Templates().List().Do()
for _, t := range templates {
	fmt.Printf("Template: %s (%s) - Theme: %s\n", t.Slug, t.Title.Rendered, t.Theme)
}
```

### Managing Template Parts (Headers, Footers, Sidebars)
```go
parts, err := client.TemplateParts().List().Do()
for _, part := range parts {
	fmt.Printf("Part: %s - Area: %s\n", part.Slug, part.Area)
}
```

### Gutenberg Navigation Blocks
Query and manage the block-based navigation menus (`wp_navigation`):
```go
navs, err := client.Navigations().List().Do()
for _, n := range navs {
	fmt.Printf("FSE Navigation: %s (Status: %s)\n", n.Title.Rendered, n.Status)
}
```
