---
title: "API Reference Overview"
description: "Comprehensive mapping of all 40 WordPress REST API endpoints supported by gowprest."
---

`gowprest` maps all 40 WordPress core REST API endpoints according to the official [WordPress Developer Handbook](https://developer.wordpress.org/rest-api/reference/).

| Module | Endpoints Covered | Client Entry Method |
|---|---|---|
| [**Posts & Pages**](/api-reference/content/) | Posts, Post Revisions, Pages, Page Revisions, Autosaves | `client.Posts()`, `client.Pages()` |
| [**Taxonomies**](/api-reference/taxonomies/) | Categories, Tags, Taxonomies, Types, Statuses | `client.Categories()`, `client.Tags()`, `client.Taxonomies()` |
| [**Users & Auth**](/api-reference/users-auth/) | Users, User `/me`, Application Passwords | `client.Users()`, `client.ApplicationPasswords()` |
| [**Media & Comments**](/api-reference/media-comments/) | Media attachments, Comments, Nested Replies | `client.Media()`, `client.Comments()` |
| [**Gutenberg & FSE**](/api-reference/gutenberg/) | Blocks, Block Types, Block Patterns, Templates, Navigation | `client.Blocks()`, `client.Templates()`, `client.Navigations()` |
| [**Administration**](/api-reference/administration/) | Settings, Themes, Global Styles, Plugins, Sidebars, Widgets, Site Health | `client.Settings()`, `client.Themes()`, `client.SiteHealth()` |
