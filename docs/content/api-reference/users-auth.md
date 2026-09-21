---
title: "Users & Authentication"
description: "Reference for User management and Application Passwords."
---

## Users (`/wp/v2/users`)

```go
users := client.Users()

// List users
userList, err := users.List().
	Roles("administrator", "editor").
	Search("john").
	Do()

// Retrieve user
u, err := users.Retrieve(userID).ContextEdit().Do()

// Create user
u, err = users.Create().
	Username("john_doe").
	Email("john@example.com").
	Roles("editor").
	Password("SecretP@ssword123").
	Do()

// Update user
u, err = users.Update().
	ID(userID).
	FirstName("John").
	LastName("Doe").
	Do()

// Delete user (requires reassigning posts)
deleted, err := users.Delete(userID).
	Reassign(1). // Reassign content to user ID 1
	Do()
```

## Current User (`/wp/v2/users/me`)

```go
me, err := client.Users().Me().Retrieve().ContextEdit().Do()
updatedMe, err := client.Users().Me().Update().
	Description("Updated bio").
	Do()
```

## Application Passwords (`/wp/v2/users/<user>/application-passwords`)

```go
appPassAPI := client.ApplicationPasswords(userID) // or "me"

// List passwords
passwords, err := appPassAPI.List().Do()

// Create password
pass, err := appPassAPI.Create().
	Name("Mobile App Client").
	Do()

// Introspect password UUID
info, err := appPassAPI.Introspect(pass.UUID).Do()

// Delete password
err = appPassAPI.Delete(pass.UUID).Do()

// Bulk delete all passwords for user
err = appPassAPI.DeleteAll().Do()
```
