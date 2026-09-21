package gowprest

import (
	"encoding/json"
	"strings"
)

// SiteHealthTestStatus constants representing test status values.
const (
	SiteHealthStatusGood        = "good"
	SiteHealthStatusRecommended = "recommended"
	SiteHealthStatusCritical    = "critical"
)

// SiteHealthTest constants for standard test slugs.
const (
	SiteHealthTestBackgroundUpdates   = "background-updates"
	SiteHealthTestLoopbackRequests    = "loopback-requests"
	SiteHealthTestHttpsStatus         = "https-status"
	SiteHealthTestDotorgCommunication = "dotorg-communication"
	SiteHealthTestAuthorizationHeader = "authorization-header"
	SiteHealthTestPageCache           = "page-cache"
)

// SiteHealthBadge represents the badge associated with a site health test.
type SiteHealthBadge struct {
	Label string `json:"label,omitempty"`
	Color string `json:"color,omitempty"`
}

// SiteHealthTest represents the result of a WordPress site health test (/wp-site-health/v1/tests/<test>).
type SiteHealthTest struct {
	Test        string          `json:"test,omitempty"`
	Label       string          `json:"label,omitempty"`
	Status      string          `json:"status,omitempty"` // "good", "recommended", "critical"
	Badge       SiteHealthBadge `json:"badge,omitempty"`
	Description string          `json:"description,omitempty"`
	Actions     string          `json:"actions,omitempty"`
}

// DirectorySizeItem represents size details for a specific directory or database.
type DirectorySizeItem struct {
	Size  string `json:"size,omitempty"`
	Debug string `json:"debug,omitempty"`
	Raw   int64  `json:"raw,omitempty"`
}

// DirectorySizes represents WordPress directory sizes information (/wp-site-health/v1/directory-sizes).
type DirectorySizes struct {
	Raw           int64             `json:"raw,omitempty"`
	WordPressSize DirectorySizeItem `json:"wordpress_size,omitempty"`
	ThemesSize    DirectorySizeItem `json:"themes_size,omitempty"`
	PluginsSize   DirectorySizeItem `json:"plugins_size,omitempty"`
	UploadsSize   DirectorySizeItem `json:"uploads_size,omitempty"`
	FontsSize     DirectorySizeItem `json:"fonts_size,omitempty"`
	DatabaseSize  DirectorySizeItem `json:"database_size,omitempty"`
	TotalSize     DirectorySizeItem `json:"total_size,omitempty"`
}

// SiteHealth provides access to WordPress Site Health endpoints (/wp-site-health/v1).
type SiteHealth struct {
	client *RestClient
}

// SiteHealth returns a SiteHealth service instance.
func (c *RestClient) SiteHealth() *SiteHealth {
	return &SiteHealth{client: c}
}

// SiteHealthTests is an ergonomic alias for SiteHealth.
func (c *RestClient) SiteHealthTests() *SiteHealth {
	return c.SiteHealth()
}

// Test returns a RetrieveSiteHealthTest builder for a given test slug.
func (api *SiteHealth) Test(testSlug string) *RetrieveSiteHealthTest {
	slug := strings.Trim(testSlug, "/")
	return &RetrieveSiteHealthTest{
		endpoint:  "/wp-site-health/v1/tests/" + slug,
		client:    api.client,
		arguments: make(map[string]string),
	}
}

// BackgroundUpdates runs the background updates site health test.
func (api *SiteHealth) BackgroundUpdates() *RetrieveSiteHealthTest {
	return api.Test(SiteHealthTestBackgroundUpdates)
}

// LoopbackRequests runs the loopback requests site health test.
func (api *SiteHealth) LoopbackRequests() *RetrieveSiteHealthTest {
	return api.Test(SiteHealthTestLoopbackRequests)
}

// HttpsStatus runs the HTTPS status site health test.
func (api *SiteHealth) HttpsStatus() *RetrieveSiteHealthTest {
	return api.Test(SiteHealthTestHttpsStatus)
}

// DotorgCommunication runs the WordPress.org communication site health test.
func (api *SiteHealth) DotorgCommunication() *RetrieveSiteHealthTest {
	return api.Test(SiteHealthTestDotorgCommunication)
}

// AuthorizationHeader runs the authorization header site health test.
func (api *SiteHealth) AuthorizationHeader() *RetrieveSiteHealthTest {
	return api.Test(SiteHealthTestAuthorizationHeader)
}

// PageCache runs the page cache site health test.
func (api *SiteHealth) PageCache() *RetrieveSiteHealthTest {
	return api.Test(SiteHealthTestPageCache)
}

// DirectorySizes returns a RetrieveDirectorySizes builder to get directory size metrics.
func (api *SiteHealth) DirectorySizes() *RetrieveDirectorySizes {
	return &RetrieveDirectorySizes{
		endpoint:  "/wp-site-health/v1/directory-sizes",
		client:    api.client,
		arguments: make(map[string]string),
	}
}

// RetrieveSiteHealthTest handles querying a single site health test.
type RetrieveSiteHealthTest struct {
	endpoint  string
	client    *RestClient
	arguments map[string]string
}

// Fields restricts the response to specific fields.
func (api *RetrieveSiteHealthTest) Fields(fields ...string) *RetrieveSiteHealthTest {
	if len(fields) > 0 {
		api.arguments["_fields"] = strings.Join(fields, ",")
	}
	return api
}

// Execute performs the GET request to retrieve the site health test result.
func (api *RetrieveSiteHealthTest) Execute() (*SiteHealthTest, error) {
	var result SiteHealthTest
	req := api.client.httpClient.R().
		SetHeader("Accept", "application/json").
		SetQueryParams(api.arguments).
		SetResult(&result)

	if api.client.auth.Username != "" && api.client.auth.Password != "" {
		req.SetBasicAuth(api.client.auth.Username, api.client.auth.Password)
	}

	url := api.client.endpoint + api.endpoint
	resp, err := req.Get(url)
	if err != nil {
		return nil, err
	}

	if resp.IsError() {
		var restErr WPRestError
		if err := json.Unmarshal(resp.Bytes(), &restErr); err == nil && restErr.Code != "" {
			return nil, &restErr
		}
		restErr.Code = "http_error"
		restErr.Message = resp.Status()
		restErr.Data.Status = resp.StatusCode()
		return nil, &restErr
	}

	return &result, nil
}

// Get is an alias for Execute.
func (api *RetrieveSiteHealthTest) Get() (*SiteHealthTest, error) {
	return api.Execute()
}

// RetrieveDirectorySizes handles querying WordPress directory sizes.
type RetrieveDirectorySizes struct {
	endpoint  string
	client    *RestClient
	arguments map[string]string
}

// Fields restricts the response to specific fields.
func (api *RetrieveDirectorySizes) Fields(fields ...string) *RetrieveDirectorySizes {
	if len(fields) > 0 {
		api.arguments["_fields"] = strings.Join(fields, ",")
	}
	return api
}

// Execute performs the GET request to retrieve directory sizes.
func (api *RetrieveDirectorySizes) Execute() (*DirectorySizes, error) {
	var result DirectorySizes
	req := api.client.httpClient.R().
		SetHeader("Accept", "application/json").
		SetQueryParams(api.arguments).
		SetResult(&result)

	if api.client.auth.Username != "" && api.client.auth.Password != "" {
		req.SetBasicAuth(api.client.auth.Username, api.client.auth.Password)
	}

	url := api.client.endpoint + api.endpoint
	resp, err := req.Get(url)
	if err != nil {
		return nil, err
	}

	if resp.IsError() {
		var restErr WPRestError
		if err := json.Unmarshal(resp.Bytes(), &restErr); err == nil && restErr.Code != "" {
			return nil, &restErr
		}
		restErr.Code = "http_error"
		restErr.Message = resp.Status()
		restErr.Data.Status = resp.StatusCode()
		return nil, &restErr
	}

	return &result, nil
}

// Get is an alias for Execute.
func (api *RetrieveDirectorySizes) Get() (*DirectorySizes, error) {
	return api.Execute()
}
