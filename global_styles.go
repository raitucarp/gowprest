package gowprest

import (
	"encoding/json"
	"strconv"
	"strings"
)

// GlobalStylesTitle represents the title of a global styles configuration.
type GlobalStylesTitle struct {
	Raw      string `json:"raw,omitempty"`
	Rendered string `json:"rendered,omitempty"`
}

// UnmarshalJSON handles both object form {"raw": "...", "rendered": "..."} and plain string form "...".
func (t *GlobalStylesTitle) UnmarshalJSON(b []byte) error {
	trimmed := strings.TrimSpace(string(b))
	if trimmed == "" || trimmed == "null" {
		return nil
	}
	if len(trimmed) > 0 && trimmed[0] == '"' {
		var s string
		if err := json.Unmarshal(b, &s); err != nil {
			return err
		}
		t.Raw = s
		t.Rendered = s
		return nil
	}
	type Alias GlobalStylesTitle
	var aux Alias
	if err := json.Unmarshal(b, &aux); err != nil {
		return err
	}
	*t = GlobalStylesTitle(aux)
	return nil
}

// GlobalStyles represents a WordPress global styles record (/wp/v2/global-styles/<id>).
type GlobalStyles struct {
	ID       int                `json:"id,omitempty"`
	Title    *GlobalStylesTitle `json:"title,omitempty"`
	Styles   map[string]any     `json:"styles,omitempty"`
	Settings map[string]any     `json:"settings,omitempty"`
	Links    map[string]any     `json:"_links,omitempty"`
}

// GlobalStylesTheme represents theme global styles or variation.
type GlobalStylesTheme struct {
	Version  int                `json:"version,omitempty"`
	Title    *GlobalStylesTitle `json:"title,omitempty"`
	Settings map[string]any     `json:"settings,omitempty"`
	Styles   map[string]any     `json:"styles,omitempty"`
	Links    map[string]any     `json:"_links,omitempty"`
}

// GlobalStylesAPI handles requests to the WordPress Global Styles API.
type GlobalStylesAPI struct {
	client *RestClient
}

// GlobalStyles returns a GlobalStylesAPI service instance.
func (c *RestClient) GlobalStyles() *GlobalStylesAPI {
	return &GlobalStylesAPI{client: c}
}

// Retrieve returns a RetrieveGlobalStyles builder for a specific global styles config ID.
func (api *GlobalStylesAPI) Retrieve(id int) *RetrieveGlobalStyles {
	return &RetrieveGlobalStyles{
		endpoint:  "/wp/v2/global-styles/" + strconv.Itoa(id),
		client:    api.client,
		arguments: make(map[string]string),
	}
}

// Update returns an UpdateGlobalStyles builder for modifying a global styles config.
func (api *GlobalStylesAPI) Update(id int) *UpdateGlobalStyles {
	return &UpdateGlobalStyles{
		endpoint: "/wp/v2/global-styles/" + strconv.Itoa(id),
		client:   api.client,
		body:     make(map[string]any),
	}
}

// Themes returns a ThemeGlobalStylesService for querying theme styles or variations.
func (api *GlobalStylesAPI) Themes(stylesheet string) *ThemeGlobalStylesService {
	return &ThemeGlobalStylesService{
		stylesheet: stylesheet,
		client:     api.client,
	}
}

// RetrieveGlobalStyles handles retrieving a global styles record.
type RetrieveGlobalStyles struct {
	endpoint  string
	client    *RestClient
	arguments map[string]string
}

// Context sets the context parameter (e.g. "view", "edit", "embed").
func (api *RetrieveGlobalStyles) Context(ctx string) *RetrieveGlobalStyles {
	api.arguments["context"] = ctx
	return api
}

// ContextView sets the context parameter to "view".
func (api *RetrieveGlobalStyles) ContextView() *RetrieveGlobalStyles {
	return api.Context("view")
}

// ContextEdit sets the context parameter to "edit".
func (api *RetrieveGlobalStyles) ContextEdit() *RetrieveGlobalStyles {
	return api.Context("edit")
}

// ContextEmbed sets the context parameter to "embed".
func (api *RetrieveGlobalStyles) ContextEmbed() *RetrieveGlobalStyles {
	return api.Context("embed")
}

// Fields limits the response to specific fields.
func (api *RetrieveGlobalStyles) Fields(fields ...string) *RetrieveGlobalStyles {
	api.arguments["_fields"] = strings.Join(fields, ",")
	return api
}

// Do executes the retrieve request.
func (api *RetrieveGlobalStyles) Do() (styles *GlobalStyles, err error) {
	endpoint := api.client.endpoint + api.endpoint

	restyClient := api.client.httpClient.R()
	if api.client.auth.Username != "" && api.client.auth.Password != "" {
		restyClient = restyClient.SetBasicAuth(api.client.auth.Username, api.client.auth.Password)
	}

	resp, err := restyClient.
		SetHeader("Accept", "application/json").
		SetResult(&styles).
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

	return styles, nil
}

// UpdateGlobalStyles handles modifying a global styles record.
type UpdateGlobalStyles struct {
	endpoint string
	client   *RestClient
	body     map[string]any
}

// Title sets the title of the global styles config.
func (api *UpdateGlobalStyles) Title(title string) *UpdateGlobalStyles {
	api.body["title"] = title
	return api
}

// Styles sets the styles object.
func (api *UpdateGlobalStyles) Styles(styles map[string]any) *UpdateGlobalStyles {
	api.body["styles"] = styles
	return api
}

// Settings sets the settings object.
func (api *UpdateGlobalStyles) Settings(settings map[string]any) *UpdateGlobalStyles {
	api.body["settings"] = settings
	return api
}

// Do executes the update request and returns the updated global styles record.
func (api *UpdateGlobalStyles) Do() (styles *GlobalStyles, err error) {
	endpoint := api.client.endpoint + api.endpoint

	restyClient := api.client.httpClient.R()
	if api.client.auth.Username != "" && api.client.auth.Password != "" {
		restyClient = restyClient.SetBasicAuth(api.client.auth.Username, api.client.auth.Password)
	}

	resp, err := restyClient.
		SetHeader("Content-Type", "application/json").
		SetHeader("Accept", "application/json").
		SetBody(api.body).
		SetResult(&styles).
		Post(endpoint)

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

	return styles, nil
}

// ThemeGlobalStylesService provides operations on a specific theme's global styles.
type ThemeGlobalStylesService struct {
	stylesheet string
	client     *RestClient
}

// Retrieve returns a builder to get the theme's global styles.
func (api *ThemeGlobalStylesService) Retrieve() *RetrieveThemeGlobalStyles {
	return &RetrieveThemeGlobalStyles{
		endpoint: "/wp/v2/global-styles/themes/" + api.stylesheet,
		client:   api.client,
	}
}

// Variations returns a builder to get the theme's global styles variations.
func (api *ThemeGlobalStylesService) Variations() *ThemeGlobalStylesVariations {
	return &ThemeGlobalStylesVariations{
		endpoint: "/wp/v2/global-styles/themes/" + api.stylesheet + "/variations",
		client:   api.client,
	}
}

// RetrieveThemeGlobalStyles handles getting a theme's global styles.
type RetrieveThemeGlobalStyles struct {
	endpoint string
	client   *RestClient
}

// Do executes the request and returns the theme global styles.
func (api *RetrieveThemeGlobalStyles) Do() (themeStyles *GlobalStylesTheme, err error) {
	endpoint := api.client.endpoint + api.endpoint

	restyClient := api.client.httpClient.R()
	if api.client.auth.Username != "" && api.client.auth.Password != "" {
		restyClient = restyClient.SetBasicAuth(api.client.auth.Username, api.client.auth.Password)
	}

	resp, err := restyClient.
		SetHeader("Accept", "application/json").
		SetResult(&themeStyles).
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

	return themeStyles, nil
}

// ThemeGlobalStylesVariations handles getting variations for a theme.
type ThemeGlobalStylesVariations struct {
	endpoint string
	client   *RestClient
}

// Do executes the request and returns the list of theme global styles variations.
func (api *ThemeGlobalStylesVariations) Do() (variations []GlobalStylesTheme, err error) {
	endpoint := api.client.endpoint + api.endpoint

	restyClient := api.client.httpClient.R()
	if api.client.auth.Username != "" && api.client.auth.Password != "" {
		restyClient = restyClient.SetBasicAuth(api.client.auth.Username, api.client.auth.Password)
	}

	resp, err := restyClient.
		SetHeader("Accept", "application/json").
		SetResult(&variations).
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

	return variations, nil
}
