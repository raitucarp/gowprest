package gowprest

import (
	"encoding/json"
	"errors"
	"strings"
)

// ThemeTags describes the raw and rendered tags of a WordPress theme.
type ThemeTags struct {
	Raw      []string `json:"raw,omitempty"`
	Rendered string   `json:"rendered,omitempty"`
}

// DefaultTemplateType describes a default block template type provided by a theme.
type DefaultTemplateType struct {
	Slug        string `json:"slug,omitempty"`
	Title       string `json:"title,omitempty"`
	Description string `json:"description,omitempty"`
}

// DefaultTemplatePartArea describes an allowed template part area in a block theme.
type DefaultTemplatePartArea struct {
	Area        string `json:"area,omitempty"`
	Label       string `json:"label,omitempty"`
	Description string `json:"description,omitempty"`
	Icon        string `json:"icon,omitempty"`
	AreaTag     string `json:"area_tag,omitempty"`
}

// Theme represents a WordPress theme.
type Theme struct {
	Stylesheet               string                    `json:"stylesheet,omitempty"`
	Template                 string                    `json:"template,omitempty"`
	RequiresPHP              string                    `json:"requires_php,omitempty"`
	RequiresWP               string                    `json:"requires_wp,omitempty"`
	Textdomain               string                    `json:"textdomain,omitempty"`
	Version                  string                    `json:"version,omitempty"`
	Screenshot               string                    `json:"screenshot,omitempty"`
	Author                   Object                    `json:"author,omitempty"`
	AuthorURI                Object                    `json:"author_uri,omitempty"`
	Description              Object                    `json:"description,omitempty"`
	Name                     Object                    `json:"name,omitempty"`
	Tags                     ThemeTags                 `json:"tags,omitempty"`
	ThemeURI                 Object                    `json:"theme_uri,omitempty"`
	Status                   string                    `json:"status,omitempty"`
	ThemeSupports            map[string]any            `json:"theme_supports,omitempty"`
	IsBlockTheme             bool                      `json:"is_block_theme,omitempty"`
	StylesheetURI            string                    `json:"stylesheet_uri,omitempty"`
	TemplateURI              string                    `json:"template_uri,omitempty"`
	DefaultTemplateTypes     []DefaultTemplateType     `json:"default_template_types,omitempty"`
	DefaultTemplatePartAreas []DefaultTemplatePartArea `json:"default_template_part_areas,omitempty"`
	Links                    map[string]any            `json:"_links,omitempty"`
}

// Themes handles requests to the WordPress themes API.
type Themes struct {
	client *RestClient
}

// Themes returns a Themes service instance.
func (c *RestClient) Themes() *Themes {
	return &Themes{client: c}
}

// List returns a ListThemes builder to query installed themes.
func (api *Themes) List() *ListThemes {
	return &ListThemes{
		endpoint:  "/wp/v2/themes",
		client:    api.client,
		arguments: make(map[string]string),
	}
}

// Retrieve returns a RetrieveTheme builder to get a specific theme by stylesheet slug.
func (api *Themes) Retrieve(stylesheet string) *RetrieveTheme {
	return &RetrieveTheme{
		endpoint:  "/wp/v2/themes/" + stylesheet,
		client:    api.client,
		arguments: make(map[string]string),
	}
}

// Active is a convenience method that retrieves the currently active theme.
func (api *Themes) Active() (*Theme, error) {
	themes, err := api.List().StatusActive().Do()
	if err != nil {
		return nil, err
	}
	if len(themes) == 0 {
		return nil, errors.New("no active theme found")
	}
	return &themes[0], nil
}

// ListThemes handles querying theme collections.
type ListThemes struct {
	endpoint  string
	client    *RestClient
	arguments map[string]string
}

func (api *ListThemes) Status(status string) *ListThemes {
	api.arguments["status"] = status
	return api
}

func (api *ListThemes) StatusActive() *ListThemes {
	return api.Status("active")
}

func (api *ListThemes) StatusInactive() *ListThemes {
	return api.Status("inactive")
}

func (api *ListThemes) Fields(fields ...string) *ListThemes {
	api.arguments["_fields"] = strings.Join(fields, ",")
	return api
}

func (api *ListThemes) Do() (themes []Theme, err error) {
	endpoint := api.client.endpoint + api.endpoint

	restyClient := api.client.httpClient.R()
	if api.client.auth.Username != "" && api.client.auth.Password != "" {
		restyClient.SetBasicAuth(api.client.auth.Username, api.client.auth.Password)
	}

	resp, err := restyClient.
		SetHeader("Accept", "application/json").
		SetResult(&themes).
		SetQueryParams(api.arguments).
		Get(endpoint)

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

	return themes, nil
}

// RetrieveTheme handles querying a single theme.
type RetrieveTheme struct {
	endpoint  string
	client    *RestClient
	arguments map[string]string
}

func (api *RetrieveTheme) Fields(fields ...string) *RetrieveTheme {
	api.arguments["_fields"] = strings.Join(fields, ",")
	return api
}

func (api *RetrieveTheme) Do() (theme *Theme, err error) {
	endpoint := api.client.endpoint + api.endpoint

	restyClient := api.client.httpClient.R()
	if api.client.auth.Username != "" && api.client.auth.Password != "" {
		restyClient.SetBasicAuth(api.client.auth.Username, api.client.auth.Password)
	}

	resp, err := restyClient.
		SetHeader("Accept", "application/json").
		SetResult(&theme).
		SetQueryParams(api.arguments).
		Get(endpoint)

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

	return theme, nil
}
