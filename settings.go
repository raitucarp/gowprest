package gowprest

import (
	"encoding/json"
	"strings"
)

// Settings represents the WordPress site settings.
type Settings struct {
	Title                string           `json:"title,omitempty"`
	Description          string           `json:"description,omitempty"`
	URL                  string           `json:"url,omitempty"`
	Email                string           `json:"email,omitempty"`
	Timezone             string           `json:"timezone,omitempty"`
	DateFormat           string           `json:"date_format,omitempty"`
	TimeFormat           string           `json:"time_format,omitempty"`
	StartOfWeek          int              `json:"start_of_week,omitempty"`
	Language             string           `json:"language,omitempty"`
	UseSmilies           bool             `json:"use_smilies,omitempty"`
	DefaultCategory      int              `json:"default_category,omitempty"`
	DefaultPostFormat    string           `json:"default_post_format,omitempty"`
	PostsPerPage         int              `json:"posts_per_page,omitempty"`
	ShowOnFront          string           `json:"show_on_front,omitempty"`
	PageOnFront          int              `json:"page_on_front,omitempty"`
	PageForPosts         int              `json:"page_for_posts,omitempty"`
	DefaultPingStatus    OpenClosedStatus `json:"default_ping_status,omitempty"`
	DefaultCommentStatus OpenClosedStatus `json:"default_comment_status,omitempty"`
	SiteLogo             *int             `json:"site_logo,omitempty"`
	SiteIcon             int              `json:"site_icon,omitempty"`
}

// SettingsData represents input payload for updating site settings.
type SettingsData struct {
	Title                *string           `json:"title,omitempty"`
	Description          *string           `json:"description,omitempty"`
	URL                  *string           `json:"url,omitempty"`
	Email                *string           `json:"email,omitempty"`
	Timezone             *string           `json:"timezone,omitempty"`
	DateFormat           *string           `json:"date_format,omitempty"`
	TimeFormat           *string           `json:"time_format,omitempty"`
	StartOfWeek          *int              `json:"start_of_week,omitempty"`
	Language             *string           `json:"language,omitempty"`
	UseSmilies           *bool             `json:"use_smilies,omitempty"`
	DefaultCategory      *int              `json:"default_category,omitempty"`
	DefaultPostFormat    *string           `json:"default_post_format,omitempty"`
	PostsPerPage         *int              `json:"posts_per_page,omitempty"`
	ShowOnFront          *string           `json:"show_on_front,omitempty"`
	PageOnFront          *int              `json:"page_on_front,omitempty"`
	PageForPosts         *int              `json:"page_for_posts,omitempty"`
	DefaultPingStatus    *OpenClosedStatus `json:"default_ping_status,omitempty"`
	DefaultCommentStatus *OpenClosedStatus `json:"default_comment_status,omitempty"`
	SiteLogo             *int              `json:"site_logo,omitempty"`
	SiteIcon             *int              `json:"site_icon,omitempty"`
}

// SettingsService handles requests to the WordPress settings API.
type SettingsService struct {
	client *RestClient
}

// Settings returns a SettingsService instance.
func (c *RestClient) Settings() *SettingsService {
	return &SettingsService{client: c}
}

// Retrieve returns a RetrieveSettings builder to get site settings.
func (api *SettingsService) Retrieve() *RetrieveSettings {
	return &RetrieveSettings{
		endpoint:  "/wp/v2/settings",
		client:    api.client,
		arguments: make(map[string]string),
	}
}

// Get is a convenience method to directly retrieve site settings.
func (api *SettingsService) Get() (*Settings, error) {
	return api.Retrieve().Do()
}

// Update returns an UpdateSettings builder to modify site settings.
func (api *SettingsService) Update(data ...SettingsData) *UpdateSettings {
	req := &UpdateSettings{
		endpoint: "/wp/v2/settings",
		client:   api.client,
		payload:  make(map[string]any),
	}
	if len(data) > 0 {
		d := data[0]
		if d.Title != nil {
			req.payload["title"] = *d.Title
		}
		if d.Description != nil {
			req.payload["description"] = *d.Description
		}
		if d.URL != nil {
			req.payload["url"] = *d.URL
		}
		if d.Email != nil {
			req.payload["email"] = *d.Email
		}
		if d.Timezone != nil {
			req.payload["timezone"] = *d.Timezone
		}
		if d.DateFormat != nil {
			req.payload["date_format"] = *d.DateFormat
		}
		if d.TimeFormat != nil {
			req.payload["time_format"] = *d.TimeFormat
		}
		if d.StartOfWeek != nil {
			req.payload["start_of_week"] = *d.StartOfWeek
		}
		if d.Language != nil {
			req.payload["language"] = *d.Language
		}
		if d.UseSmilies != nil {
			req.payload["use_smilies"] = *d.UseSmilies
		}
		if d.DefaultCategory != nil {
			req.payload["default_category"] = *d.DefaultCategory
		}
		if d.DefaultPostFormat != nil {
			req.payload["default_post_format"] = *d.DefaultPostFormat
		}
		if d.PostsPerPage != nil {
			req.payload["posts_per_page"] = *d.PostsPerPage
		}
		if d.ShowOnFront != nil {
			req.payload["show_on_front"] = *d.ShowOnFront
		}
		if d.PageOnFront != nil {
			req.payload["page_on_front"] = *d.PageOnFront
		}
		if d.PageForPosts != nil {
			req.payload["page_for_posts"] = *d.PageForPosts
		}
		if d.DefaultPingStatus != nil {
			req.payload["default_ping_status"] = *d.DefaultPingStatus
		}
		if d.DefaultCommentStatus != nil {
			req.payload["default_comment_status"] = *d.DefaultCommentStatus
		}
		if d.SiteLogo != nil {
			req.payload["site_logo"] = *d.SiteLogo
		}
		if d.SiteIcon != nil {
			req.payload["site_icon"] = *d.SiteIcon
		}
	}
	return req
}

// RetrieveSettings handles retrieving site settings.
type RetrieveSettings struct {
	endpoint  string
	client    *RestClient
	arguments map[string]string
}

func (api *RetrieveSettings) Fields(fields ...string) *RetrieveSettings {
	api.arguments["_fields"] = strings.Join(fields, ",")
	return api
}

func (api *RetrieveSettings) Do() (*Settings, error) {
	restyClient := api.client.httpClient.R()
	if api.client.auth.Username != "" && api.client.auth.Password != "" {
		restyClient.SetBasicAuth(api.client.auth.Username, api.client.auth.Password)
	}

	var settings Settings
	resp, err := restyClient.
		SetHeader("Accept", "application/json").
		SetResult(&settings).
		SetQueryParams(api.arguments).
		Get(api.client.endpoint + api.endpoint)

	if err != nil {
		return nil, err
	}

	if resp.IsError() {
		var wpError WPRestError
		err = json.Unmarshal(resp.Bytes(), &wpError)
		if err != nil {
			return nil, err
		}
		return nil, &wpError
	}

	return &settings, nil
}

// UpdateSettings handles modifying site settings.
type UpdateSettings struct {
	endpoint string
	client   *RestClient
	payload  map[string]any
}

func (api *UpdateSettings) Title(title string) *UpdateSettings {
	api.payload["title"] = title
	return api
}

func (api *UpdateSettings) Description(description string) *UpdateSettings {
	api.payload["description"] = description
	return api
}

func (api *UpdateSettings) URL(url string) *UpdateSettings {
	api.payload["url"] = url
	return api
}

func (api *UpdateSettings) Email(email string) *UpdateSettings {
	api.payload["email"] = email
	return api
}

func (api *UpdateSettings) Timezone(timezone string) *UpdateSettings {
	api.payload["timezone"] = timezone
	return api
}

func (api *UpdateSettings) DateFormat(format string) *UpdateSettings {
	api.payload["date_format"] = format
	return api
}

func (api *UpdateSettings) TimeFormat(format string) *UpdateSettings {
	api.payload["time_format"] = format
	return api
}

func (api *UpdateSettings) StartOfWeek(day int) *UpdateSettings {
	api.payload["start_of_week"] = day
	return api
}

func (api *UpdateSettings) Language(language string) *UpdateSettings {
	api.payload["language"] = language
	return api
}

func (api *UpdateSettings) UseSmilies(use bool) *UpdateSettings {
	api.payload["use_smilies"] = use
	return api
}

func (api *UpdateSettings) DefaultCategory(categoryID int) *UpdateSettings {
	api.payload["default_category"] = categoryID
	return api
}

func (api *UpdateSettings) DefaultPostFormat(format string) *UpdateSettings {
	api.payload["default_post_format"] = format
	return api
}

func (api *UpdateSettings) PostsPerPage(postsPerPage int) *UpdateSettings {
	api.payload["posts_per_page"] = postsPerPage
	return api
}

func (api *UpdateSettings) ShowOnFront(show string) *UpdateSettings {
	api.payload["show_on_front"] = show
	return api
}

func (api *UpdateSettings) PageOnFront(pageID int) *UpdateSettings {
	api.payload["page_on_front"] = pageID
	return api
}

func (api *UpdateSettings) PageForPosts(pageID int) *UpdateSettings {
	api.payload["page_for_posts"] = pageID
	return api
}

func (api *UpdateSettings) DefaultPingStatus(status OpenClosedStatus) *UpdateSettings {
	api.payload["default_ping_status"] = status
	return api
}

func (api *UpdateSettings) DefaultCommentStatus(status OpenClosedStatus) *UpdateSettings {
	api.payload["default_comment_status"] = status
	return api
}

func (api *UpdateSettings) SiteLogo(logoID int) *UpdateSettings {
	api.payload["site_logo"] = logoID
	return api
}

func (api *UpdateSettings) SiteIcon(iconID int) *UpdateSettings {
	api.payload["site_icon"] = iconID
	return api
}

func (api *UpdateSettings) Set(key string, value any) *UpdateSettings {
	api.payload[key] = value
	return api
}

func (api *UpdateSettings) Do() (*Settings, error) {
	restyClient := api.client.httpClient.R()
	if api.client.auth.Username != "" && api.client.auth.Password != "" {
		restyClient.SetBasicAuth(api.client.auth.Username, api.client.auth.Password)
	}

	var settings Settings
	resp, err := restyClient.
		SetHeader("Content-Type", "application/json").
		SetHeader("Accept", "application/json").
		SetBody(api.payload).
		SetResult(&settings).
		Post(api.client.endpoint + api.endpoint)

	if err != nil {
		return nil, err
	}

	if resp.IsError() {
		var wpError WPRestError
		err = json.Unmarshal(resp.Bytes(), &wpError)
		if err != nil {
			return nil, err
		}
		return nil, &wpError
	}

	return &settings, nil
}
