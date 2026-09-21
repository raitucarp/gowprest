---
title: "Authentication"
description: "Configuring authentication with WordPress Application Passwords or Basic Auth."
---

WordPress REST API requires authentication for write operations (`POST`, `PUT`, `DELETE`), viewing drafts, or inspecting private settings.

## Application Passwords (Recommended)

WordPress 5.6+ natively supports **Application Passwords**. Application Passwords are user-specific tokens that can be easily revoked without revealing or changing your account's main WordPress login password.

### Generating an Application Password in WP Admin

1. Log into your WordPress admin dashboard (`/wp-admin/`).
2. Navigate to **Users** → **Profile** (or edit any user).
3. Scroll down to the **Application Passwords** section.
4. Enter a name (e.g., `gowprest-cli` or `ci-publisher`) and click **Add New Application Password**.
5. Copy the generated 24-character password (e.g. `abcd efgh ijkl mnop qrst uvwx`).

### Using Application Passwords in Code

```go
package main

import (
	"github.com/raitucarp/gowprest"
)

func main() {
	client := gowprest.NewClient("https://example.com").
		WithBasicAuth("admin_username", "abcd efgh ijkl mnop qrst uvwx")
	defer client.Close()

	// You can now execute authenticated calls
	me, err := client.Users().Me().Retrieve().ContextEdit().Do()
	if err != nil {
		panic(err)
	}
	println("Authenticated as:", me.Username)
}
```

## Application Passwords API via `gowprest`

`gowprest` also allows managing Application Passwords programmatically! You can create, list, introspect, and delete passwords:

```go
// Generate an application password for a user
newPass, err := client.ApplicationPasswords(targetUserID).Create().
	Name("Automated Pipeline").
	Do()
if err != nil {
	log.Fatal(err)
}

fmt.Printf("Generated UUID: %s, Password: %s\n", newPass.UUID, newPass.Password)
```

> **Security Note**: Always use **HTTPS** in production when transmitting Application Passwords to protect credentials in transit.
