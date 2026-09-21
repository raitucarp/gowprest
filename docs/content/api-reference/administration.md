---
title: "Administration, Themes, & Site Health"
description: "Reference for Settings, Themes, Global Styles, Plugins, Sidebars, Widgets, and Site Health."
---

## Site Settings (`/wp/v2/settings`)

```go
settings, err := client.Settings().Get()

updated, err := client.Settings().Update().
	Title("My New Blog").
	Description("Fresh perspective").
	Timezone("UTC").
	Do()
```

## Themes & Global Styles (`/wp/v2/themes`, `/wp/v2/global-styles`)

```go
themes, err := client.Themes().List().Status("active").Do()
activeTheme, err := client.Themes().Active()

styles, err := client.GlobalStyles().Themes(activeTheme.Stylesheet).Retrieve().Do()
```

## Plugins (`/wp/v2/plugins`)

```go
plugins, err := client.Plugins().List().Do()

// Install plugin from WordPress.org
installed, err := client.Plugins().Install("classic-widgets").Do()

// Activate
activated, err := client.Plugins().Activate(installed.Plugin).Do()

// Deactivate
deactivated, err := client.Plugins().Deactivate(installed.Plugin).Do()

// Delete
deleted, err := client.Plugins().Delete(installed.Plugin).Do()
```

## Sidebars & Widgets (`/wp/v2/sidebars`, `/wp/v2/widgets`)

```go
sidebars, err := client.Sidebars().List().Do()

// Create widget
widget, err := client.Widgets().Create().
	Sidebar("sidebar-1").
	IDBase("block").
	Instance(map[string]any{
		"raw": "<!-- wp:paragraph --><p>Custom widget</p><!-- /wp:paragraph -->",
	}).
	Do()
```

## Site Health (`/wp-site-health/v1/`)

```go
siteHealth := client.SiteHealth()

// Dedicated test runners
loopback, err := siteHealth.LoopbackRequests().Execute()
bgUpdates, err := siteHealth.BackgroundUpdates().Execute()
httpsTest, err := siteHealth.HttpsStatus().Execute()
dotorg, err := siteHealth.DotorgCommunication().Execute()
authHeader, err := siteHealth.AuthorizationHeader().Execute()

// Directory sizing metrics
sizes, err := siteHealth.DirectorySizes().Execute()
```
