---
title: "Taxonomies, Terms, & User Management"
description: "Working with hierarchical categories, tags, user roles, and application passwords."
---

## Categories & Tags

### Creating Hierarchical Categories
```go
// 1. Create a parent category
techCat, err := client.Categories().Create().
	Name("Technology").
	Slug("technology").
	Description("Articles on modern engineering").
	Do()

// 2. Create a subcategory nested under Technology
goCat, err := client.Categories().Create().
	Name("Golang").
	Slug("golang").
	Parent(techCat.ID).
	Do()
```

### Creating Post Tags
Tags are flat, non-hierarchical taxonomy terms:
```go
tag, err := client.Tags().Create().
	Name("Open Source").
	Slug("open-source").
	Do()
```

---

## User Management

Manage WordPress users with administrative, editor, author, or subscriber privileges:

### Creating a User
```go
user, err := client.Users().Create().
	Username("jane_author").
	Email("jane@example.com").
	FirstName("Jane").
	LastName("Doe").
	Roles("author").
	Password("ComplexP@ssw0rd!2026").
	Do()
```

### Inspecting Current Authenticated User (`/me`)
To retrieve the credentials and capabilities of the token being used:
```go
me, err := client.Users().Me().Retrieve().
	ContextEdit(). // ContextEdit exposes roles and username
	Do()

fmt.Printf("Logged in as: %s (Roles: %v)\n", me.Username, me.Roles)
```

### Generating Application Passwords
Generate API tokens for automated background bots or microservices:
```go
pass, err := client.ApplicationPasswords(user.ID).Create().
	Name("Nightly Backup Worker").
	Do()

fmt.Println("Generated token:", pass.Password)
```
