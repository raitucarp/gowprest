---
title: "WordPress Site Health Diagnostics"
description: "Executing WordPress core health checks and directory sizing metrics."
---

WordPress includes diagnostic APIs introduced in Site Health (`/wp-site-health/v1/`). `gowprest` provides dedicated runners for these official tests.

## Running Site Health Tests

Execute specific health diagnostics:

```go
package main

import (
	"fmt"
	"log"
	"github.com/raitucarp/gowprest"
)

func main() {
	client := gowprest.NewClient("https://example.com").
		WithBasicAuth("admin", "app_password")
	defer client.Close()

	// 1. Loopback requests test
	loopback, err := client.SiteHealth().LoopbackRequests().Execute()
	if err != nil {
		log.Fatalf("Loopback check failed: %v", err)
	}
	fmt.Printf("Loopback: %s (%s)\n", loopback.Status, loopback.Label)

	// 2. Background updates test
	bgUpdates, err := client.SiteHealth().BackgroundUpdates().Execute()
	if err != nil {
		log.Fatalf("Background updates check failed: %v", err)
	}
	fmt.Printf("Background Updates: %s (%s)\n", bgUpdates.Status, bgUpdates.Label)

	// 3. HTTPS Status test
	httpsTest, err := client.SiteHealth().HttpsStatus().Execute()
	if err != nil {
		log.Fatalf("HTTPS check failed: %v", err)
	}
	fmt.Printf("HTTPS: %s (%s)\n", httpsTest.Status, httpsTest.Label)

	// 4. Dotorg communication test
	dotorg, err := client.SiteHealth().DotorgCommunication().Execute()
	if err != nil {
		log.Fatalf("Dotorg communication check failed: %v", err)
	}
	fmt.Printf("WordPress.org Connectivity: %s (%s)\n", dotorg.Status, dotorg.Label)
}
```

## Inspecting Directory Sizes

Query the filesystem breakdown of your WordPress installation:

```go
sizes, err := client.SiteHealth().DirectorySizes().Execute()
if err != nil {
	log.Fatal(err)
}

fmt.Printf("Total WordPress Size: %s\n", sizes.TotalSize)
fmt.Printf("Uploads: %s\n", sizes.UploadsSize)
fmt.Printf("Themes: %s\n", sizes.ThemesSize)
fmt.Printf("Plugins: %s\n", sizes.PluginsSize)
fmt.Printf("Database: %s\n", sizes.DatabaseSize)
```
