// Package gowprest provides a modern, fluent, strongly-typed Go client library
// for the WordPress REST API.
//
// It supports all 40 WordPress core endpoints including Posts, Pages, Media,
// Comments, Taxonomies, Users, Settings, Themes, Gutenberg Blocks, Block Patterns,
// Full Site Editing (FSE) Navigation & Templates, Application Passwords, and
// Site Health Diagnostics.
//
// Basic Usage:
//
//	client := gowprest.NewClient("https://example.com").
//		WithBasicAuth("username", "app_password")
//	defer client.Close()
//
//	posts, err := client.Posts().List().
//		Status("publish").
//		PerPage(10).
//		Do()
package gowprest

import (
	"strings"

	"resty.dev/v3"
)

// BlogInfo represents WordPress site discovery information returned by
// the index endpoint (/wp-json/).
type BlogInfo struct {
	// Name is the site title.
	Name string `json:"name"`
	// Description is the site tagline/description.
	Description string `json:"description"`
	// URL is the home URL of the site.
	URL string `json:"url"`
	// Home is the home URL of the site.
	Home string `json:"home"`
	// GmtOffset is the site's GMT offset in hours.
	GmtOffset string `json:"gmt_offset"`
	// TimezoneString is the site's PHP timezone string.
	TimezoneString string `json:"timezone_string"`
	// PageForPosts is the ID of the page set as the blog posts index.
	PageForPosts int `json:"page_for_posts"`
	// PageOnFront is the ID of the page set as the front page.
	PageOnFront int `json:"page_on_front"`
	// ShowOnFront specifies whether the front page displays latest posts or a static page.
	ShowOnFront string `json:"show_on_front"`
	// Namespaces lists the available REST API route namespaces on the site.
	Namespaces []string `json:"namespaces"`
	// Authentication contains authentication scheme metadata exposed by the site.
	Authentication any `json:"authentication"`
	// SiteLogo is the attachment ID of the site logo.
	SiteLogo int `json:"site_logo"`
	// SiteIcon is the attachment ID of the site icon (favicon).
	SiteIcon int `json:"site_icon"`
	// SiteIconURL is the full URL to the site icon.
	SiteIconURL string `json:"site_icon_url"`
}

// Authentication stores user credentials (typically username and Application Password)
// used for authenticating WordPress REST API requests.
type Authentication struct {
	Username string
	Password string
}

// RestClient is the primary client for interacting with the WordPress REST API.
// It manages the HTTP transport, authentication, base endpoints, and provides
// factory methods for all WordPress resource endpoints.
type RestClient struct {
	baseURL  string
	endpoint string
	auth     Authentication

	httpClient *resty.Client
}

// Close terminates idle connections and cleans up the underlying HTTP client.
func (api *RestClient) Close() {
	api.httpClient.Close()
}

// WithBasicAuth configures HTTP Basic Authentication credentials on the client.
// In modern WordPress (WP 5.6+), this is commonly used with Application Passwords.
func (api *RestClient) WithBasicAuth(username, password string) *RestClient {
	api.auth.Username = username
	api.auth.Password = password
	return api
}

// Discover queries the WordPress REST API index endpoint (/wp-json/) to retrieve
// general site metadata, available namespaces, timezone settings, and authentication schemes.
func (api *RestClient) Discover() (info BlogInfo, err error) {
	_, err = api.httpClient.R().
		SetHeader("Accept", "application/json").
		SetResult(&info).
		Get(api.endpoint)

	if err != nil {
		return
	}

	return
}

// NewClient initializes a new WordPress REST API client targeting the given baseURL.
// The baseURL can be provided with or without a trailing slash (e.g. "https://example.com").
func NewClient(baseURL string) *RestClient {
	client := resty.New()
	return &RestClient{
		baseURL:    baseURL,
		endpoint:   strings.Trim(baseURL, "/") + "/wp-json",
		httpClient: client,
	}
}
